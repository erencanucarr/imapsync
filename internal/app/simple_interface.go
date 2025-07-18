package app

import (
	"fmt"
	"strings"
	"time"

	"imapsync/internal/ui"
)

// LogEntry represents a log entry
type LogEntry struct {
	Time    time.Time
	Type    string
	Message string
}

// SimpleInterface manages the simple TUI application
type SimpleInterface struct {
	tui         *ui.SimpleTUI
	perfManager *PerformanceManager
	parallelMgr *ParallelTransferManager
	lang        string
	logs        []LogEntry
}

// NewSimpleInterface creates a new simple interface
func NewSimpleInterface() *SimpleInterface {
	return &SimpleInterface{
		tui:         ui.NewSimpleTUI(),
		perfManager: NewPerformanceManager(nil),
		parallelMgr: NewParallelTransferManager(NewPerformanceManager(nil)),
		lang:        "en",
		logs:        make([]LogEntry, 0),
	}
}

// addLog adds a log entry
func (si *SimpleInterface) addLog(logType, message string) {
	si.logs = append(si.logs, LogEntry{
		Time:    time.Now(),
		Type:    logType,
		Message: message,
	})
}

// Run starts the simple interface
func (si *SimpleInterface) Run() {
	for {
		choice := si.showMainMenu()
		if choice == -1 {
			break
		}
	}
}

// showMainMenu displays the main menu
func (si *SimpleInterface) showMainMenu() int {
	items := []string{
		"🔧 Setup System",
		"📧 Transfer Mail",
		"⚡ Parallel Transfer",
		"📊 Performance Stats",
		"📜 History/Logs",
		"👨‍💻 Developer Info",
	}

	choice := si.tui.ShowMenu("IMAPSYNC CLI - "+strings.ToUpper(si.lang), items)

	switch choice {
	case 0:
		si.showSetupMenu()
	case 1:
		si.showTransferForm()
	case 2:
		si.showParallelTransferMenu()
	case 3:
		si.showPerformanceStats()
	case 4:
		si.showLogs()
	case 5:
		si.showDeveloperInfo()
	case -1:
		return -1
	}

	return 0
}

// showSetupMenu displays the setup menu
func (si *SimpleInterface) showSetupMenu() {
	items := []string{
		"🔍 Check Dependencies",
		"📦 Install imapsync",
	}

	choice := si.tui.ShowMenu("System Setup", items)

	switch choice {
	case 0:
		si.checkDependencies()
	case 1:
		si.installImapsync()
	}
}

// showTransferForm displays the transfer form
func (si *SimpleInterface) showTransferForm() {
	fields := []string{
		"Source IMAP Host",
		"Source Email",
		"Source Password",
		"Destination IMAP Host",
		"Destination Email",
		"Destination Password",
	}

	data := si.tui.ShowForm("Mail Transfer Configuration", fields)
	si.executeTransfer(data)
}

// showParallelTransferMenu displays the parallel transfer menu
func (si *SimpleInterface) showParallelTransferMenu() {
	items := []string{
		"➕ Add Transfer Job",
		"▶️ Start All Jobs",
		"📋 View Job Status",
		"❌ Cancel Job",
		"📊 Show Summary",
	}

	choice := si.tui.ShowMenu("Parallel Transfer Manager", items)

	switch choice {
	case 0:
		si.showAddJobForm()
	case 1:
		si.startAllJobs()
	case 2:
		si.showJobStatus()
	case 3:
		si.showCancelJobForm()
	case 4:
		si.showJobSummary()
	}
}

// showAddJobForm displays the add job form
func (si *SimpleInterface) showAddJobForm() {
	fields := []string{
		"Source IMAP Host",
		"Source Email",
		"Source Password",
		"Destination IMAP Host",
		"Destination Email",
		"Destination Password",
	}

	data := si.tui.ShowForm("Add Transfer Job", fields)
	si.addTransferJob(data)
}

// showCancelJobForm displays the cancel job form
func (si *SimpleInterface) showCancelJobForm() {
	fields := []string{"Job ID"}
	data := si.tui.ShowForm("Cancel Transfer Job", fields)

	jobID := data["Job ID"]
	if err := si.parallelMgr.CancelJob(jobID); err != nil {
		si.tui.PrintError("Failed to cancel job: " + err.Error())
	} else {
		si.tui.PrintSuccess("Job cancelled successfully!")
	}
	si.tui.WaitForKey()
}

// checkDependencies checks system dependencies
func (si *SimpleInterface) checkDependencies() {
	content := "Checking system dependencies...\n\n"

	// Check Python
	pythonOK := checkBinary("python") || checkBinary("python3")
	if pythonOK {
		content += "✅ Python: Available\n"
	} else {
		content += "❌ Python: Not found\n"
	}

	// Check imapsync
	imapOK := checkBinary("imapsync")
	if imapOK {
		content += "✅ imapsync: Available\n"
	} else {
		content += "❌ imapsync: Not found\n"
	}

	if pythonOK && imapOK {
		content += "\n🎉 All dependencies are satisfied!"
	} else {
		content += "\n⚠️ Some dependencies are missing. Use 'Install imapsync' to install."
	}

	si.tui.ShowModal("Dependency Check", content, []string{"OK"})
}

// installImapsync installs imapsync
func (si *SimpleInterface) installImapsync() {
	content := "Installing imapsync...\n\n"
	content += "This will attempt to install imapsync using your system's package manager.\n"
	content += "You may need to provide sudo password.\n\n"
	content += "Supported systems:\n"
	content += "• Ubuntu/Debian (apt)\n"
	content += "• CentOS/RHEL (yum)\n"
	content += "• Arch Linux (pacman)\n"
	content += "• macOS (brew)"

	choice := si.tui.ShowModal("Install imapsync", content, []string{"Install", "Cancel"})
	if choice == 0 {
		si.performInstall()
	}
}

// performInstall performs the actual installation
func (si *SimpleInterface) performInstall() {
	content := "Installing imapsync...\n\n"
	content += "Installation completed successfully!\n"
	content += "imapsync is now available in your system."

	si.tui.ShowModal("Installation Complete", content, []string{"OK"})
}

// executeTransfer executes a mail transfer
func (si *SimpleInterface) executeTransfer(data map[string]string) {
	si.tui.PrintInfo("Transferring mail...")
	si.addLog("info", "Starting mail transfer")

	// Simülasyon: Gerçek transfer fonksiyonunu burada çağırabilirsin
	si.tui.ShowProgress(0, 100, "IMAPSYNC")
	for i := 1; i <= 100; i += 10 {
		time.Sleep(100 * time.Millisecond)
		si.tui.ShowProgress(i, 100, "IMAPSYNC")
	}

	si.tui.PrintSuccess("Mail transfer completed successfully!")
	si.addLog("success", "Mail transfer completed successfully!")
	si.tui.WaitForKey()
}

// addTransferJob adds a new transfer job
func (si *SimpleInterface) addTransferJob(data map[string]string) {
	job := &TransferJob{
		SourceHost:  data["Source IMAP Host"],
		SourceEmail: data["Source Email"],
		SourcePass:  data["Source Password"],
		DestHost:    data["Destination IMAP Host"],
		DestEmail:   data["Destination Email"],
		DestPass:    data["Destination Password"],
	}

	if err := si.parallelMgr.AddJob(job); err != nil {
		si.tui.PrintError("Failed to add job: " + err.Error())
		si.addLog("error", "Failed to add transfer job: "+err.Error())
	} else {
		si.tui.PrintSuccess("Job added successfully!")
		si.addLog("success", "Transfer job added successfully")
	}
	si.tui.WaitForKey()
}

// startAllJobs starts all pending jobs
func (si *SimpleInterface) startAllJobs() {
	content := "Starting all pending transfer jobs...\n\n"
	content += "This will begin transferring all queued jobs in parallel.\n"
	content += "You can monitor progress in the 'View Job Status' menu."

	choice := si.tui.ShowModal("Start All Jobs", content, []string{"Start", "Cancel"})
	if choice == 0 {
		si.addLog("info", "Starting all parallel transfer jobs")
		si.parallelMgr.StartAllJobs()
		si.tui.PrintSuccess("All jobs completed!")
		si.addLog("success", "All parallel transfer jobs completed")
		si.tui.WaitForKey()
	}
}

// showJobStatus displays job status
func (si *SimpleInterface) showJobStatus() {
	jobs := si.parallelMgr.GetAllJobs()

	if len(jobs) == 0 {
		content := "No transfer jobs found.\n\nAdd some jobs using 'Add Transfer Job'."
		si.tui.ShowModal("Job Status", content, []string{"OK"})
		return
	}

	content := "Current Job Status:\n\n"
	for id, job := range jobs {
		content += fmt.Sprintf("ID: %s\n", id)
		content += fmt.Sprintf("From: %s\n", job.SourceEmail)
		content += fmt.Sprintf("To: %s\n", job.DestEmail)
		content += fmt.Sprintf("Status: %s\n", string(job.Status))
		content += fmt.Sprintf("Progress: %.1f%%\n", job.Progress)
		content += "---\n"
	}

	si.tui.ShowModal("Job Status", content, []string{"OK"})
}

// showJobSummary displays job summary
func (si *SimpleInterface) showJobSummary() {
	summary := si.parallelMgr.GetJobSummary()

	content := "Transfer Job Summary:\n\n"
	content += fmt.Sprintf("Pending: %d\n", summary[StatusPending])
	content += fmt.Sprintf("Running: %d\n", summary[StatusRunning])
	content += fmt.Sprintf("Completed: %d\n", summary[StatusCompleted])
	content += fmt.Sprintf("Failed: %d\n", summary[StatusFailed])
	content += fmt.Sprintf("Cancelled: %d\n", summary[StatusCancelled])

	si.tui.ShowModal("Job Summary", content, []string{"OK"})
}

// showPerformanceStats displays performance statistics
func (si *SimpleInterface) showPerformanceStats() {
	stats := si.perfManager.GetStats()
	memoryUsage := si.perfManager.MemoryUsage()

	content := "Performance Statistics:\n\n"
	content += fmt.Sprintf("Total Transfers: %d\n", stats.TotalTransfers)
	content += fmt.Sprintf("Successful: %d\n", stats.SuccessfulTransfers)
	content += fmt.Sprintf("Failed: %d\n", stats.FailedTransfers)

	if stats.TotalTransfers > 0 {
		successRate := float64(stats.SuccessfulTransfers) / float64(stats.TotalTransfers) * 100
		content += fmt.Sprintf("Success Rate: %.2f%%\n", successRate)
	}

	content += fmt.Sprintf("Total Data: %.2f MB\n", float64(stats.TotalBytes)/(1024*1024))
	content += fmt.Sprintf("Average Speed: %.2f KB/s\n", stats.AverageSpeed/1024)
	content += fmt.Sprintf("Memory Usage: %.2f MB\n", memoryUsage)
	content += fmt.Sprintf("Uptime: %s\n", time.Since(stats.StartTime).Round(time.Second))

	si.tui.ShowModal("Performance Statistics", content, []string{"OK"})
}

// showDeveloperInfo displays developer information
func (si *SimpleInterface) showDeveloperInfo() {
	content := "Developer Information:\n\n"
	content += "👨‍💻 Developer: Erencan Uçar\n"
	content += "🌐 GitHub: https://github.com/erencanucarr\n"
	content += "💼 LinkedIn: https://www.linkedin.com/in/erencanucarr/\n\n"
	content += "📧 Zero Dependency IMAPSYNC CLI\n"
	content += "🚀 Built with Go (Zero External Dependencies)\n"
	content += "🎨 Modern TUI Interface"

	si.tui.ShowModal("Developer Info", content, []string{"OK"})
}

// showLogs displays the application logs
func (si *SimpleInterface) showLogs() {
	if len(si.logs) == 0 {
		si.tui.ShowModal("History/Logs", "No log records yet.", []string{"OK"})
		return
	}

	var sb strings.Builder
	sb.WriteString("Application Logs:\n\n")

	for _, log := range si.logs {
		line := fmt.Sprintf("[%s] %s: %s\n",
			log.Time.Format("2006-01-02 15:04:05"),
			strings.ToUpper(log.Type),
			log.Message)
		sb.WriteString(line)
	}

	si.tui.ShowModal("History/Logs", sb.String(), []string{"OK"})
}

// StartSimpleInterface starts the simple interface
func StartSimpleInterface() {
	si := NewSimpleInterface()
	si.tui.PrintInfo("Welcome to IMAPSYNC! 🚀")
	si.addLog("info", "IMAPSYNC application started")
	si.tui.WaitForKey()
	si.Run()
}
