package app

import (
	"fmt"
	"strings"
	"time"
)

// ProgressBar represents a simple progress bar
type ProgressBar struct {
	current     int
	total       int
	width       int
	description string
	startTime   time.Time
	lastUpdate  time.Time
}

// NewProgressBar creates a new progress bar
func NewProgressBar(total int) *ProgressBar {
	return &ProgressBar{
		current:     0,
		total:       total,
		width:       50,
		description: "Progress",
		startTime:   time.Now(),
		lastUpdate:  time.Now(),
	}
}

// SetDescription sets the description of the progress bar
func (pb *ProgressBar) SetDescription(desc string) {
	pb.description = desc
}

// Set sets the current progress value
func (pb *ProgressBar) Set(current int) {
	pb.current = current
	if pb.current > pb.total {
		pb.current = pb.total
	}

	// Only update display every 100ms to avoid too frequent updates
	if time.Since(pb.lastUpdate) > 100*time.Millisecond {
		pb.render()
		pb.lastUpdate = time.Now()
	}
}

// Add increments the progress by the given amount
func (pb *ProgressBar) Add(amount int) {
	pb.Set(pb.current + amount)
}

// render displays the progress bar
func (pb *ProgressBar) render() {
	if pb.total <= 0 {
		return
	}

	percentage := float64(pb.current) / float64(pb.total) * 100
	filled := int(float64(pb.width) * percentage / 100)

	bar := strings.Repeat("=", filled)
	if filled < pb.width {
		bar += ">"
		bar += strings.Repeat(" ", pb.width-filled-1)
	}

	elapsed := time.Since(pb.startTime)
	eta := time.Duration(0)
	if pb.current > 0 {
		eta = time.Duration(float64(elapsed) * float64(pb.total-pb.current) / float64(pb.current))
	}

	// Clear line and print progress
	fmt.Printf("\r[%s] %s %d/%d (%.1f%%) ETA: %s",
		bar,
		pb.description,
		pb.current,
		pb.total,
		percentage,
		formatDuration(eta))
}

// Finish completes the progress bar
func (pb *ProgressBar) Finish() {
	pb.Set(pb.total)
	fmt.Println() // New line after progress bar
}

// formatDuration formats a duration in a human-readable way
func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	} else if d < time.Hour {
		return fmt.Sprintf("%dm%ds", int(d.Minutes()), int(d.Seconds())%60)
	} else {
		return fmt.Sprintf("%dh%dm", int(d.Hours()), int(d.Minutes())%60)
	}
}
