package ui

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// SimpleTUI represents a simple TUI application
type SimpleTUI struct {
	reader *bufio.Reader
}

// NewSimpleTUI creates a new simple TUI
func NewSimpleTUI() *SimpleTUI {
	return &SimpleTUI{
		reader: bufio.NewReader(os.Stdin),
	}
}

// ShowBanner displays a beautiful banner
func (tui *SimpleTUI) ShowBanner() {
	tui.ClearScreen()
	fmt.Println("")
	fmt.Println(Cyan(Bold("         IMAPSYNC")))
	fmt.Println("")
	fmt.Println(Blue("  📧 Transfer emails between IMAP servers efficiently"))
	fmt.Println("")
	fmt.Println(Dim("  " + strings.Repeat("═", 50)))
	fmt.Println("")
}

// ShowMenu displays a beautiful menu
func (tui *SimpleTUI) ShowMenu(title string, items []string) int {
	tui.ShowBanner()
	fmt.Println(Bold(Yellow("  📋 " + title)))
	fmt.Println(Dim("  " + strings.Repeat("─", 70)))
	fmt.Println("")

	for i, item := range items {
		icon := "🔹"
		color := White
		switch i {
		case 0:
			icon, color = "📧", Blue
		case 1:
			icon, color = "⚡", Yellow
		case 2:
			icon, color = "📊", Green
		case 3:
			icon, color = "⚙️", Purple
		case 4:
			icon, color = "ℹ️", Cyan
		case 5:
			icon, color = "🖥️", Blue
		case 6:
			icon, color = "📜", Purple
		case 7:
			icon, color = "🚪", Red
		}
		fmt.Printf("  %s %2d. %s\n", color(icon), i+1, White(item))
	}
	fmt.Printf("  %s %2d. %s\n", Red("🚪"), len(items)+1, Red("Exit"))
	fmt.Println("")
	fmt.Println(Dim("  " + strings.Repeat("─", 70)))

	for {
		fmt.Print(Bold(Yellow("  🎯 Enter your choice (1-" + fmt.Sprintf("%d", len(items)+1) + "): ")))
		var choice int
		fmt.Scanln(&choice)
		if choice == len(items)+1 {
			return -1 // Exit
		}
		if choice > 0 && choice <= len(items) {
			return choice - 1
		}
		fmt.Println(Red("  ❌ Invalid choice. Please try again."))
	}
}

// ShowModal displays a beautiful modal
func (tui *SimpleTUI) ShowModal(title, content string, buttons []string) int {
	tui.ShowBanner()
	fmt.Println(Bold(Yellow("  📋 " + title)))
	fmt.Println(Dim("  " + strings.Repeat("─", 70)))
	fmt.Println("")
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		fmt.Printf("  %s\n", White(line))
	}
	fmt.Println("")
	fmt.Println(Dim("  " + strings.Repeat("─", 70)))
	for i, button := range buttons {
		icon := "🔘"
		color := White
		if strings.Contains(strings.ToLower(button), "ok") || strings.Contains(strings.ToLower(button), "yes") {
			icon, color = "✅", Green
		} else if strings.Contains(strings.ToLower(button), "cancel") || strings.Contains(strings.ToLower(button), "no") {
			icon, color = "❌", Red
		}
		if i > 0 {
			fmt.Print("  ")
		}
		fmt.Printf("%s [%d] %s", color(icon), i+1, White(button))
	}
	fmt.Println()
	fmt.Println("")
	for {
		fmt.Print(Bold(Yellow("  🎯 Enter your choice: ")))
		var choice int
		fmt.Scanln(&choice)
		if choice == 0 {
			return -1 // Cancel
		}
		if choice > 0 && choice <= len(buttons) {
			return choice - 1
		}
		fmt.Println(Red("  ❌ Invalid choice. Please try again."))
	}
}

// ShowForm displays a beautiful form
func (tui *SimpleTUI) ShowForm(title string, fields []string) map[string]string {
	tui.ShowBanner()
	fmt.Println(Bold(Yellow("  📝 " + title)))
	fmt.Println(Dim("  " + strings.Repeat("─", 70)))
	fmt.Println("")
	result := make(map[string]string)
	for _, field := range fields {
		icon := "📝"
		color := White
		if strings.Contains(strings.ToLower(field), "email") {
			icon, color = "📧", Blue
		} else if strings.Contains(strings.ToLower(field), "password") {
			icon, color = "🔒", Red
		} else if strings.Contains(strings.ToLower(field), "server") {
			icon, color = "🌐", Green
		} else if strings.Contains(strings.ToLower(field), "port") {
			icon, color = "🔌", Yellow
		}
		fmt.Printf("  %s %s:\n", color(icon), White(field))
		fmt.Print(Cyan("  └─ "))
		var value string
		fmt.Scanln(&value)
		result[field] = value
		fmt.Println("")
	}
	fmt.Println(Dim("  " + strings.Repeat("─", 70)))
	return result
}

// ClearScreen clears the screen
func (tui *SimpleTUI) ClearScreen() {
	fmt.Print("\033[2J")
	fmt.Print("\033[H")
}

// PrintSuccess prints a success message
func (tui *SimpleTUI) PrintSuccess(message string) {
	fmt.Printf("  %s %s\n", Green("✅"), Green(message))
}

// PrintError prints an error message
func (tui *SimpleTUI) PrintError(message string) {
	fmt.Printf("  %s %s\n", Red("❌"), Red(message))
}

// PrintInfo prints an info message
func (tui *SimpleTUI) PrintInfo(message string) {
	fmt.Printf("  %s %s\n", Cyan("ℹ️"), Cyan(message))
}

// PrintWarning prints a warning message
func (tui *SimpleTUI) PrintWarning(message string) {
	fmt.Printf("  %s %s\n", Yellow("⚠️"), Yellow(message))
}

// ShowProgress displays a progress bar
func (tui *SimpleTUI) ShowProgress(current, total int, description string) {
	percentage := float64(current) / float64(total) * 100
	barWidth := 50
	filled := int(float64(barWidth) * percentage / 100)
	bar := Green(strings.Repeat("█", filled)) + Dim(strings.Repeat("░", barWidth-filled))
	fmt.Printf("\r  %s %s [%s] %.1f%% (%d/%d)", Blue("📊"), White(description), bar, percentage, current, total)
	if current >= total {
		fmt.Println()
	}
}

// ShowLoading displays a loading animation
func (tui *SimpleTUI) ShowLoading(message string) {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	for i := 0; i < 10; i++ {
		fmt.Printf("\r  %s %s", Cyan(frames[i%len(frames)]), White(message))
		time.Sleep(100 * time.Millisecond)
	}
	fmt.Println()
}

// WaitForKey waits for user to press Enter
func (tui *SimpleTUI) WaitForKey() {
	fmt.Print(Yellow("  ⏸️  Press Enter to continue..."))
	tui.reader.ReadString('\n')
}

// ShowRealTimeStats displays real-time statistics
func (tui *SimpleTUI) ShowRealTimeStats(stats map[string]interface{}) {
	tui.ClearScreen()
	fmt.Println("")
	fmt.Println(Cyan(Bold("         IMAPSYNC")))
	fmt.Println("")
	fmt.Println(Bold(Yellow("  📊 Real-Time Statistics")))
	fmt.Println(Dim("  " + strings.Repeat("─", 70)))
	fmt.Println("")

	// Display stats in a nice format
	for key, value := range stats {
		icon := ""
		color := White

		switch key {
		case "Total Transfers":
			icon, color = "📧", Blue
		case "Successful":
			icon, color = "✅", Green
		case "Failed":
			icon, color = "❌", Red
		case "Success Rate":
			icon, color = "📈", Yellow
		case "Total Data":
			icon, color = "💾", Purple
		case "Average Speed":
			icon, color = "⚡", Cyan
		case "Memory Usage":
			icon, color = "💾", Blue
		case "Uptime":
			icon, color = "⏱️", Green
		}

		fmt.Printf("  %s %s: %s\n", color(icon), White(key), Cyan(fmt.Sprintf("%v", value)))
	}

	fmt.Println("")
	fmt.Println(Dim("  " + strings.Repeat("─", 70)))
	fmt.Println(Yellow("  🔄 Press Enter to refresh, 'q' to quit"))
}

// ShowLiveProgress displays live progress with real-time updates
func (tui *SimpleTUI) ShowLiveProgress(jobID string, current, total int, speed float64, eta time.Duration) {
	percentage := float64(current) / float64(total) * 100
	barWidth := 50
	filled := int(float64(barWidth) * percentage / 100)
	bar := Green(strings.Repeat("█", filled)) + Dim(strings.Repeat("░", barWidth-filled))

	fmt.Printf("\r  %s Job %s [%s] %.1f%% (%d/%d) %.2f KB/s ETA: %s",
		Blue("📊"),
		White(jobID),
		bar,
		percentage,
		current,
		total,
		speed/1024,
		eta.Round(time.Second))

	if current >= total {
		fmt.Println()
	}
}
