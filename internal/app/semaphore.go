package app

import (
	"context"
	"sync"
)

// Semaphore provides a simple semaphore implementation
type Semaphore struct {
	permits int64
	mu      sync.Mutex
	cond    *sync.Cond
}

// NewSemaphore creates a new semaphore with the given number of permits
func NewSemaphore(permits int64) *Semaphore {
	s := &Semaphore{
		permits: permits,
	}
	s.cond = sync.NewCond(&s.mu)
	return s
}

// Acquire acquires a permit from the semaphore
func (s *Semaphore) Acquire(ctx context.Context, n int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	for s.permits < n {
		// Check if context is cancelled
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		s.cond.Wait()
	}

	s.permits -= n
	return nil
}

// Release releases permits back to the semaphore
func (s *Semaphore) Release(n int64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.permits += n
	s.cond.Broadcast()
}

// TryAcquire attempts to acquire a permit without blocking
func (s *Semaphore) TryAcquire(n int64) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.permits >= n {
		s.permits -= n
		return true
	}
	return false
}

// Available returns the number of available permits
func (s *Semaphore) Available() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.permits
}
