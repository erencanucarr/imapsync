package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"imapsync/internal/ui"
)

// ParallelTransfer handles multiple transfer jobs in parallel
func ParallelTransfer(lang string) {
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
func ShowPerformanceStats(lang string) {
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
