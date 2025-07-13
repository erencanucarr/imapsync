package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"

	"imapsync/internal/app"
	"imapsync/internal/i18n"
	"imapsync/internal/ui"
)

func main() {
	lang := flag.String("lang", "tr", "Language: tr or en")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println(ui.Cyan(i18n.T(*lang, "menu")))
		fmt.Println(i18n.T(*lang, "menu_setup"))
		fmt.Println(i18n.T(*lang, "menu_transfer"))
		fmt.Println(i18n.T(*lang, "menu_parallel"))
		fmt.Println(i18n.T(*lang, "menu_stats"))
		fmt.Println(i18n.T(*lang, "menu_developer"))
		fmt.Println(i18n.T(*lang, "menu_exit"))

		fmt.Print(ui.Green(i18n.T(*lang, "choice")))
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			app.SetupSystem(*lang)
		case "2":
			app.TransferMail(*lang)
		case "3":
			app.ParallelTransfer(*lang)
		case "4":
			app.ShowPerformanceStats(*lang)
		case "5":
			app.ShowDeveloper(*lang)
		case "6":
			fmt.Println(ui.Yellow(i18n.T(*lang, "exit")))
			return
		default:
			fmt.Println(ui.Red(i18n.T(*lang, "invalid")))
		}
	}
}
