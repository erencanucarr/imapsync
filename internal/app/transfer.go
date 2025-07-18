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
	"time"

	"imapsync/internal/ui"
)

// TransferMail runs imapsync and shows a progress bar.
// It parses stdout looking for "Transferred:" lines to update progress.
func TransferMail() {
	fmt.Println(ui.Cyan("Starting mail transfer..."))

	// Initialize performance manager
	perfManager := NewPerformanceManager(nil)
	defer perfManager.PrintStats()

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Source IMAP host: ")
	srcHost, _ := reader.ReadString('\n')
	srcHost = strings.TrimSpace(srcHost)

	fmt.Print("Source email: ")
	srcEmail, _ := reader.ReadString('\n')
	srcEmail = strings.TrimSpace(srcEmail)

	fmt.Print("Source password: ")
	srcPass, _ := ReadPassword()
	fmt.Println()

	fmt.Print("Destination IMAP host: ")
	dstHost, _ := reader.ReadString('\n')
	dstHost = strings.TrimSpace(dstHost)

	fmt.Print("Destination email: ")
	dstEmail, _ := reader.ReadString('\n')
	dstEmail = strings.TrimSpace(dstEmail)

	fmt.Print("Destination password: ")
	dstPass, _ := ReadPassword()
	fmt.Println()

	// Check cache for previous successful transfers
	cacheKey := fmt.Sprintf("%s_%s_%s", srcEmail, dstEmail, srcHost)
	if cachedData, found := perfManager.GetCachedData(cacheKey); found {
		fmt.Println(ui.Yellow("Found cached transfer data for this combination"))
		fmt.Printf("Last successful transfer: %v\n", cachedData)
	}

	fmt.Println(ui.Cyan("Testing credentials..."))

	// Use retry mechanism for credential testing
	ctx := context.Background()
	err := perfManager.RetryWithBackoff(ctx, func() error {
		testCmd := exec.Command("imapsync", "--justlogin", "--host1", srcHost, "--ssl1", "--user1", srcEmail, "--password1", srcPass, "--host2", dstHost, "--ssl2", "--user2", dstEmail, "--password2", dstPass)
		return testCmd.Run()
	})

	if err != nil {
		fmt.Println(ui.Red("Error:"), err)
		return
	}

	// Acquire connection from pool
	if err := perfManager.AcquireConnection(ctx); err != nil {
		fmt.Println(ui.Red("Failed to acquire connection from pool"), err)
		return
	}
	defer perfManager.ReleaseConnection()

	// Check memory usage before starting transfer
	if !perfManager.CheckMemoryLimit() {
		fmt.Println(ui.Yellow("Memory usage high, optimizing..."))
		perfManager.OptimizeMemory()
	}

	args := []string{
		"--host1", srcHost, "--ssl1",
		"--user1", srcEmail, "--password1", srcPass,
		"--host2", dstHost, "--ssl2",
		"--user2", dstEmail, "--password2", dstPass,
		"--exclude", "^Junk\\ E-Mail",
		"--exclude", "^Deleted\\ Items",
		"--exclude", "^Deleted",
		"--exclude", "^Trash",
		"--regextrans2", "s#^Sent$#Sent Items#",
		"--regextrans2", "s#^Spam$#Junk E-Mail#",
		"--useuid",
		"--usecache",
		"--tmpdir", "./tmp",
		"--syncinternaldates",
		"--progress",
	}

	startTime := time.Now()
	cmd := exec.Command("imapsync", args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(ui.Red("Error:"), err)
		return
	}

	if err := cmd.Start(); err != nil {
		fmt.Println(ui.Red("Error:"), err)
		return
	}

	bar := NewProgressBar(100)
	bar.SetDescription("IMAPSYNC")

	scanner := bufio.NewScanner(stdout)
	percentRe := regexp.MustCompile(`([0-9]{1,3}(?:\.[0-9]+)?)%`)
	ratioRe := regexp.MustCompile(`(?i)(\d+)/(\d+)`)
	for scanner.Scan() {
		line := scanner.Text()

		if m := percentRe.FindStringSubmatch(line); len(m) == 2 {
			p, _ := strconv.ParseFloat(m[1], 64)
			bar.Set(int(p))
			continue
		}
		if m := ratioRe.FindStringSubmatch(line); len(m) == 3 {
			current, _ := strconv.Atoi(m[1])
			total, _ := strconv.Atoi(m[2])
			if total > 0 {
				bar.Set(int(float64(current) / float64(total) * 100))
			}
		}
	}

	var bytesTransferred int64
	var transferSuccess bool

	if err := cmd.Wait(); err != nil {
		fmt.Println() // newline after bar
		fmt.Println(ui.Red("Mail transfer failed."))
		fmt.Println(ui.Red("Error:"), err)
		transferSuccess = false
	} else {
		transferSuccess = true
	}

	bar.Finish()
	fmt.Println() // newline after bar

	// Calculate transfer statistics
	duration := time.Since(startTime)
	if transferSuccess {
		// Estimate bytes transferred (this would be more accurate with actual parsing)
		bytesTransferred = 1024 * 1024 // 1MB estimate
		perfManager.UpdateStats(true, bytesTransferred)

		// Cache successful transfer
		transferInfo := map[string]interface{}{
			"timestamp": time.Now(),
			"duration":  duration,
			"bytes":     bytesTransferred,
		}
		perfManager.SetCachedData(cacheKey, transferInfo)

		fmt.Println(ui.Green("Mail transfer completed successfully!"))
		fmt.Printf("Transfer completed in %s\n", duration.Round(time.Second))
	} else {
		perfManager.UpdateStats(false, 0)
		fmt.Println(ui.Red("Mail transfer failed."))
	}
}
