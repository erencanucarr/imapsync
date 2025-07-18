package app

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"imapsync/internal/ui"
)

// TransferJob represents a single transfer job
type TransferJob struct {
	ID               string
	SourceHost       string
	SourceEmail      string
	SourcePass       string
	DestHost         string
	DestEmail        string
	DestPass         string
	Status           TransferStatus
	Progress         float64
	Error            error
	StartTime        time.Time
	EndTime          time.Time
	BytesTransferred int64
}

// TransferStatus represents the status of a transfer job
type TransferStatus string

const (
	StatusPending   TransferStatus = "pending"
	StatusRunning   TransferStatus = "running"
	StatusCompleted TransferStatus = "completed"
	StatusFailed    TransferStatus = "failed"
	StatusCancelled TransferStatus = "cancelled"
)

// ParallelTransferManager manages parallel transfer operations
type ParallelTransferManager struct {
	jobs        map[string]*TransferJob
	mu          sync.RWMutex
	perfManager *PerformanceManager
	logger      *Logger
	ctx         context.Context
	cancel      context.CancelFunc
}

// NewParallelTransferManager creates a new parallel transfer manager
func NewParallelTransferManager(perfManager *PerformanceManager) *ParallelTransferManager {
	ctx, cancel := context.WithCancel(context.Background())

	logger := NewLogger()
	logger.SetLevel(LevelInfo)

	return &ParallelTransferManager{
		jobs:        make(map[string]*TransferJob),
		perfManager: perfManager,
		logger:      logger,
		ctx:         ctx,
		cancel:      cancel,
	}
}

// AddJob adds a new transfer job to the queue
func (ptm *ParallelTransferManager) AddJob(job *TransferJob) error {
	ptm.mu.Lock()
	defer ptm.mu.Unlock()

	if job.ID == "" {
		job.ID = fmt.Sprintf("job_%d", time.Now().UnixNano())
	}

	job.Status = StatusPending
	ptm.jobs[job.ID] = job

	ptm.logger.Info("Added transfer job: %s (%s -> %s)", job.ID, job.SourceEmail, job.DestEmail)
	return nil
}

// StartAllJobs starts all pending jobs in parallel
func (ptm *ParallelTransferManager) StartAllJobs() {
	ptm.mu.RLock()
	var pendingJobs []*TransferJob
	for _, job := range ptm.jobs {
		if job.Status == StatusPending {
			pendingJobs = append(pendingJobs, job)
		}
	}
	ptm.mu.RUnlock()

	ptm.logger.Info("Starting %d transfer jobs in parallel", len(pendingJobs))

	var wg sync.WaitGroup
	for _, job := range pendingJobs {
		wg.Add(1)
		go func(j *TransferJob) {
			defer wg.Done()
			ptm.executeJob(j)
		}(job)
	}

	wg.Wait()
	ptm.logger.Info("All transfer jobs completed")
}

// executeJob executes a single transfer job
func (ptm *ParallelTransferManager) executeJob(job *TransferJob) {
	// Acquire connection from pool
	if err := ptm.perfManager.AcquireConnection(ptm.ctx); err != nil {
		ptm.updateJobStatus(job, StatusFailed, err)
		return
	}
	defer ptm.perfManager.ReleaseConnection()

	ptm.updateJobStatus(job, StatusRunning, nil)
	job.StartTime = time.Now()

	// Execute transfer with retry logic
	err := ptm.perfManager.RetryWithBackoff(ptm.ctx, func() error {
		return ptm.runImapsync(job)
	})

	job.EndTime = time.Now()

	if err != nil {
		ptm.updateJobStatus(job, StatusFailed, err)
		ptm.perfManager.UpdateStats(false, job.BytesTransferred)
	} else {
		ptm.updateJobStatus(job, StatusCompleted, nil)
		ptm.perfManager.UpdateStats(true, job.BytesTransferred)
	}
}

// runImapsync runs the actual imapsync command for a job
func (ptm *ParallelTransferManager) runImapsync(job *TransferJob) error {
	// This is a simplified version - in real implementation, you'd parse imapsync output
	// and update progress in real-time

	args := []string{
		"--host1", job.SourceHost, "--ssl1",
		"--user1", job.SourceEmail, "--password1", job.SourcePass,
		"--host2", job.DestHost, "--ssl2",
		"--user2", job.DestEmail, "--password2", job.DestPass,
		"--exclude", "^Junk\\ E-Mail",
		"--exclude", "^Deleted\\ Items",
		"--exclude", "^Deleted",
		"--exclude", "^Trash",
		"--regextrans2", "s#^Sent$#Sent Items#",
		"--regextrans2", "s#^Spam$#Junk E-Mail#",
		"--useuid",
		"--usecache",
		"--tmpdir", fmt.Sprintf("./tmp_%s", job.ID),
		"--syncinternaldates",
		"--progress",
	}

	// Execute imapsync command
	cmd := exec.Command("imapsync", args...)

	// Set up output parsing for progress updates
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start imapsync: %w", err)
	}

	// Parse output for progress updates
	scanner := bufio.NewScanner(stdout)
	percentRe := regexp.MustCompile(`([0-9]{1,3}(?:\.[0-9]+)?)%`)

	for scanner.Scan() {
		line := scanner.Text()

		if m := percentRe.FindStringSubmatch(line); len(m) == 2 {
			if p, err := strconv.ParseFloat(m[1], 64); err == nil {
				ptm.updateJobProgress(job, p)
			}
		}

		// Check if context is cancelled
		select {
		case <-ptm.ctx.Done():
			cmd.Process.Kill()
			return ptm.ctx.Err()
		default:
		}
	}

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("imapsync failed: %w", err)
	}

	return nil
}

// updateJobStatus updates the status of a job
func (ptm *ParallelTransferManager) updateJobStatus(job *TransferJob, status TransferStatus, err error) {
	ptm.mu.Lock()
	defer ptm.mu.Unlock()

	job.Status = status
	job.Error = err

	ptm.logger.Info("Job %s status: %s", job.ID, status)
	if err != nil {
		ptm.logger.Error("Job %s error: %v", job.ID, err)
	}
}

// updateJobProgress updates the progress of a job
func (ptm *ParallelTransferManager) updateJobProgress(job *TransferJob, progress float64) {
	ptm.mu.Lock()
	defer ptm.mu.Unlock()

	job.Progress = progress
}

// GetJobStatus returns the status of a specific job
func (ptm *ParallelTransferManager) GetJobStatus(jobID string) (*TransferJob, bool) {
	ptm.mu.RLock()
	defer ptm.mu.RUnlock()

	job, exists := ptm.jobs[jobID]
	return job, exists
}

// GetAllJobs returns all jobs
func (ptm *ParallelTransferManager) GetAllJobs() map[string]*TransferJob {
	ptm.mu.RLock()
	defer ptm.mu.RUnlock()

	result := make(map[string]*TransferJob)
	for k, v := range ptm.jobs {
		result[k] = v
	}
	return result
}

// CancelJob cancels a specific job
func (ptm *ParallelTransferManager) CancelJob(jobID string) error {
	ptm.mu.Lock()
	defer ptm.mu.Unlock()

	job, exists := ptm.jobs[jobID]
	if !exists {
		return fmt.Errorf("job %s not found", jobID)
	}

	if job.Status == StatusRunning {
		job.Status = StatusCancelled
		ptm.logger.Info("Cancelled job: %s", jobID)
	}

	return nil
}

// CancelAllJobs cancels all running jobs
func (ptm *ParallelTransferManager) CancelAllJobs() {
	ptm.cancel()

	ptm.mu.Lock()
	defer ptm.mu.Unlock()

	for _, job := range ptm.jobs {
		if job.Status == StatusRunning {
			job.Status = StatusCancelled
		}
	}

	ptm.logger.Info("Cancelled all running jobs")
}

// GetJobSummary returns a summary of all jobs
func (ptm *ParallelTransferManager) GetJobSummary() map[TransferStatus]int {
	ptm.mu.RLock()
	defer ptm.mu.RUnlock()

	summary := make(map[TransferStatus]int)
	for _, job := range ptm.jobs {
		summary[job.Status]++
	}

	return summary
}

// PrintJobSummary prints a summary of all jobs
func (ptm *ParallelTransferManager) PrintJobSummary() {
	summary := ptm.GetJobSummary()

	fmt.Printf("\n=== Transfer Job Summary ===\n")
	fmt.Printf("Pending: %d\n", summary[StatusPending])
	fmt.Printf("Running: %d\n", summary[StatusRunning])
	fmt.Printf("Completed: %d\n", summary[StatusCompleted])
	fmt.Printf("Failed: %d\n", summary[StatusFailed])
	fmt.Printf("Cancelled: %d\n", summary[StatusCancelled])

	total := 0
	for _, count := range summary {
		total += count
	}
	fmt.Printf("Total: %d\n", total)
}

// ParallelTransfer handles multiple transfer jobs in parallel
func ParallelTransfer() {
	fmt.Println(ui.Cyan("=== Parallel Transfer Manager ==="))

	// Initialize managers
	perfManager := NewPerformanceManager(nil)
	parallelManager := NewParallelTransferManager(perfManager)

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\n1 - Add Transfer Job")
		fmt.Println("2 - Start All Jobs")
		fmt.Println("3 - View Job Status")
		fmt.Println("4 - Cancel Job")
		fmt.Println("5 - Show Summary")
		fmt.Println("6 - Back to Main Menu")

		fmt.Print("Choice: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			addTransferJob(parallelManager, reader)
		case "2":
			fmt.Println(ui.Cyan("Starting all pending jobs..."))
			parallelManager.StartAllJobs()
		case "3":
			showJobStatus(parallelManager)
		case "4":
			cancelJob(parallelManager, reader)
		case "5":
			parallelManager.PrintJobSummary()
		case "6":
			return
		default:
			fmt.Println(ui.Red("Invalid choice"))
		}
	}
}

// addTransferJob adds a new transfer job to the queue
func addTransferJob(ptm *ParallelTransferManager, reader *bufio.Reader) {
	fmt.Println(ui.Cyan("=== Add Transfer Job ==="))

	job := &TransferJob{}

	fmt.Print("Source IMAP host: ")
	srcHost, _ := reader.ReadString('\n')
	job.SourceHost = strings.TrimSpace(srcHost)

	fmt.Print("Source email: ")
	srcEmail, _ := reader.ReadString('\n')
	job.SourceEmail = strings.TrimSpace(srcEmail)

	fmt.Print("Source password: ")
	srcPass, _ := ReadPassword()
	job.SourcePass = srcPass
	fmt.Println()

	fmt.Print("Destination IMAP host: ")
	dstHost, _ := reader.ReadString('\n')
	job.DestHost = strings.TrimSpace(dstHost)

	fmt.Print("Destination email: ")
	dstEmail, _ := reader.ReadString('\n')
	job.DestEmail = strings.TrimSpace(dstEmail)

	fmt.Print("Destination password: ")
	dstPass, _ := ReadPassword()
	job.DestPass = dstPass
	fmt.Println()

	if err := ptm.AddJob(job); err != nil {
		fmt.Println(ui.Red("Failed to add job:"), err)
	} else {
		fmt.Println(ui.Green("Job added successfully!"))
	}
}

// showJobStatus displays the status of all jobs
func showJobStatus(ptm *ParallelTransferManager) {
	jobs := ptm.GetAllJobs()

	if len(jobs) == 0 {
		fmt.Println(ui.Yellow("No jobs found"))
		return
	}

	fmt.Println(ui.Cyan("=== Job Status ==="))
	for id, job := range jobs {
		statusColor := ui.Green
		switch job.Status {
		case StatusPending:
			statusColor = ui.Yellow
		case StatusRunning:
			statusColor = ui.Cyan
		case StatusFailed:
			statusColor = ui.Red
		case StatusCancelled:
			statusColor = ui.Red
		}

		fmt.Printf("ID: %s\n", id)
		fmt.Printf("  From: %s\n", job.SourceEmail)
		fmt.Printf("  To: %s\n", job.DestEmail)
		fmt.Printf("  Status: %s\n", statusColor(string(job.Status)))
		fmt.Printf("  Progress: %.1f%%\n", job.Progress)

		if job.StartTime != (time.Time{}) {
			fmt.Printf("  Started: %s\n", job.StartTime.Format("2006-01-02 15:04:05"))
		}

		if job.Error != nil {
			fmt.Printf("  Error: %v\n", job.Error)
		}
		fmt.Println()
	}
}

// cancelJob cancels a specific job
func cancelJob(ptm *ParallelTransferManager, reader *bufio.Reader) {
	fmt.Print("Enter job ID to cancel: ")
	jobID, _ := reader.ReadString('\n')
	jobID = strings.TrimSpace(jobID)

	if err := ptm.CancelJob(jobID); err != nil {
		fmt.Println(ui.Red("Failed to cancel job:"), err)
	} else {
		fmt.Println(ui.Green("Job cancelled successfully!"))
	}
}

// ShowPerformanceStats displays performance statistics
func ShowPerformanceStats() {
	fmt.Println(ui.Cyan("=== Performance Statistics ==="))

	// Create a new performance manager to show stats
	perfManager := NewPerformanceManager(nil)

	// Show memory usage
	memoryUsage := perfManager.MemoryUsage()
	fmt.Printf("Current Memory Usage: %.2f MB\n", memoryUsage)

	if perfManager.CheckMemoryLimit() {
		fmt.Println(ui.Green("Memory usage is within limits"))
	} else {
		fmt.Println(ui.Red("Memory usage is above limit"))
	}

	// Show cache information
	fmt.Printf("Cache Items: %d\n", perfManager.cache.ItemCount())

	// Show connection pool status
	activeConnections := perfManager.config.MaxConcurrentTransfers - int(perfManager.semaphore.Available())
	fmt.Printf("Active Connections: %d/%d\n", activeConnections, perfManager.config.MaxConcurrentTransfers)

	// Show transfer statistics
	perfManager.PrintStats()

	fmt.Println("\nPress Enter to continue...")
	fmt.Scanln()
}
