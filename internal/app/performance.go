package app

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"time"
)

// PerformanceConfig holds performance-related configuration
type PerformanceConfig struct {
	MaxConcurrentTransfers int           // Maximum concurrent transfer operations
	ConnectionPoolSize     int           // Size of connection pool
	CacheExpiration        time.Duration // Cache expiration time
	CacheCleanupInterval   time.Duration // Cache cleanup interval
	MemoryLimitMB          int           // Memory limit in MB
	RetryAttempts          int           // Number of retry attempts
	RetryDelay             time.Duration // Delay between retries
}

// DefaultPerformanceConfig returns default performance settings
func DefaultPerformanceConfig() *PerformanceConfig {
	return &PerformanceConfig{
		MaxConcurrentTransfers: 3,                // Limit concurrent transfers
		ConnectionPoolSize:     5,                // Connection pool size
		CacheExpiration:        30 * time.Minute, // Cache for 30 minutes
		CacheCleanupInterval:   10 * time.Minute, // Cleanup every 10 minutes
		MemoryLimitMB:          512,              // 512MB memory limit
		RetryAttempts:          3,                // 3 retry attempts
		RetryDelay:             5 * time.Second,  // 5 second delay
	}
}

// PerformanceManager handles performance optimizations
type PerformanceManager struct {
	config    *PerformanceConfig
	cache     *Cache
	semaphore *Semaphore
	mu        sync.RWMutex
	logger    *Logger
	stats     *TransferStats
}

// TransferStats tracks transfer performance metrics
type TransferStats struct {
	mu                  sync.RWMutex
	TotalTransfers      int64
	SuccessfulTransfers int64
	FailedTransfers     int64
	TotalBytes          int64
	AverageSpeed        float64 // bytes per second
	StartTime           time.Time
	LastTransferTime    time.Time
}

// NewPerformanceManager creates a new performance manager
func NewPerformanceManager(config *PerformanceConfig) *PerformanceManager {
	if config == nil {
		config = DefaultPerformanceConfig()
	}

	logger := NewLogger()
	logger.SetLevel(LevelInfo)

	return &PerformanceManager{
		config:    config,
		cache:     NewCache(),
		semaphore: NewSemaphore(int64(config.MaxConcurrentTransfers)),
		logger:    logger,
		stats: &TransferStats{
			StartTime: time.Now(),
		},
	}
}

// AcquireConnection acquires a connection from the pool
func (pm *PerformanceManager) AcquireConnection(ctx context.Context) error {
	return pm.semaphore.Acquire(ctx, 1)
}

// ReleaseConnection releases a connection back to the pool
func (pm *PerformanceManager) ReleaseConnection() {
	pm.semaphore.Release(1)
}

// GetCachedData retrieves data from cache
func (pm *PerformanceManager) GetCachedData(key string) (interface{}, bool) {
	return pm.cache.Get(key)
}

// SetCachedData stores data in cache
func (pm *PerformanceManager) SetCachedData(key string, data interface{}) {
	pm.cache.Set(key, data, pm.config.CacheExpiration)
}

// InvalidateCache removes specific cache entries
func (pm *PerformanceManager) InvalidateCache(key string) {
	pm.cache.Delete(key)
}

// ClearCache clears all cache entries
func (pm *PerformanceManager) ClearCache() {
	pm.cache.Flush()
}

// UpdateStats updates transfer statistics
func (pm *PerformanceManager) UpdateStats(success bool, bytesTransferred int64) {
	pm.stats.mu.Lock()
	defer pm.stats.mu.Unlock()

	pm.stats.TotalTransfers++
	if success {
		pm.stats.SuccessfulTransfers++
	} else {
		pm.stats.FailedTransfers++
	}
	pm.stats.TotalBytes += bytesTransferred
	pm.stats.LastTransferTime = time.Now()

	// Calculate average speed
	duration := time.Since(pm.stats.StartTime).Seconds()
	if duration > 0 {
		pm.stats.AverageSpeed = float64(pm.stats.TotalBytes) / duration
	}
}

// GetStats returns current transfer statistics
func (pm *PerformanceManager) GetStats() TransferStats {
	pm.stats.mu.RLock()
	defer pm.stats.mu.RUnlock()

	return *pm.stats
}

// PrintStats prints performance statistics
func (pm *PerformanceManager) PrintStats() {
	stats := pm.GetStats()

	fmt.Printf("\n=== Performance Statistics ===\n")
	fmt.Printf("Total Transfers: %d\n", stats.TotalTransfers)
	fmt.Printf("Successful: %d\n", stats.SuccessfulTransfers)
	fmt.Printf("Failed: %d\n", stats.FailedTransfers)
	if stats.TotalTransfers > 0 {
		fmt.Printf("Success Rate: %.2f%%\n", float64(stats.SuccessfulTransfers)/float64(stats.TotalTransfers)*100)
	}
	fmt.Printf("Total Data: %.2f MB\n", float64(stats.TotalBytes)/(1024*1024))
	fmt.Printf("Average Speed: %.2f KB/s\n", stats.AverageSpeed/1024)
	fmt.Printf("Uptime: %s\n", time.Since(stats.StartTime).Round(time.Second))
	fmt.Printf("Last Transfer: %s\n", stats.LastTransferTime.Format("2006-01-02 15:04:05"))
	fmt.Printf("Cache Items: %d\n", pm.cache.ItemCount())
	fmt.Printf("Active Connections: %d/%d\n", pm.config.MaxConcurrentTransfers-int(pm.semaphore.Available()), pm.config.MaxConcurrentTransfers)
}

// RetryWithBackoff executes a function with retry logic
func (pm *PerformanceManager) RetryWithBackoff(ctx context.Context, operation func() error) error {
	var lastErr error

	for attempt := 0; attempt < pm.config.RetryAttempts; attempt++ {
		if err := operation(); err == nil {
			return nil
		} else {
			lastErr = err
			pm.logger.Warn("Attempt %d failed: %v", attempt+1, err)

			if attempt < pm.config.RetryAttempts-1 {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-time.After(pm.config.RetryDelay * time.Duration(attempt+1)):
					continue
				}
			}
		}
	}

	return fmt.Errorf("operation failed after %d attempts: %w", pm.config.RetryAttempts, lastErr)
}

// MemoryUsage returns current memory usage in MB
func (pm *PerformanceManager) MemoryUsage() float64 {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return float64(m.Alloc) / 1024 / 1024
}

// CheckMemoryLimit checks if memory usage is within limits
func (pm *PerformanceManager) CheckMemoryLimit() bool {
	usage := pm.MemoryUsage()
	return usage < float64(pm.config.MemoryLimitMB)
}

// OptimizeMemory performs memory optimization
func (pm *PerformanceManager) OptimizeMemory() {
	// Clear cache if memory usage is high
	if !pm.CheckMemoryLimit() {
		pm.logger.Info("Memory usage high, clearing cache")
		pm.ClearCache()
		runtime.GC()
	}
}
