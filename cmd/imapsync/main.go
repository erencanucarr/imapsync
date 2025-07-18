package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"imapsync/internal/app"
	"imapsync/internal/ui"
)

func main() {
	// Default to TUI mode, but allow CLI mode with -cli flag
	cliMode := flag.Bool("cli", false, "Enable CLI mode (default is TUI)")
	flag.Parse()

	if !*cliMode {
		// Start TUI mode by default
		app.StartSimpleInterface()
		return
	}

	// Original CLI mode (only when -cli flag is used)
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println(ui.Cyan("Please select an option:"))
		fmt.Println("1 - Setup System")
		fmt.Println("2 - Transfer Mail")
		fmt.Println("3 - Parallel Transfer")
		fmt.Println("4 - Performance Stats")
		fmt.Println("5 - Developer")
		fmt.Println("6 - Modern TUI Interface")
		fmt.Println("7 - Exit")

		fmt.Print(ui.Green("Choice (1-7): "))
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			app.SetupSystem()
		case "2":
			app.TransferMail()
		case "3":
			app.ParallelTransfer()
		case "4":
			app.ShowPerformanceStats()
		case "5":
			app.ShowDeveloper()
		case "6":
			app.StartSimpleInterface()
		case "7":
			fmt.Println(ui.Yellow("Exiting program..."))
			return
		default:
			fmt.Println(ui.Red("Invalid choice. Please enter 1-7."))
		}
	}
}
