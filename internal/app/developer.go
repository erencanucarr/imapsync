package app

import (
    "fmt"
    "imapsync/internal/ui"
)

// ShowDeveloper displays developer info (placeholder)
func ShowDeveloper(lang string) {
    fmt.Println(ui.Cyan("Developer: Erencan Uçar"))
    fmt.Println(ui.Green("GitHub: https://github.com/erencanucarr"))
    fmt.Println(ui.Green("LinkedIn: https://www.linkedin.com/in/erencanucarr/"))
    fmt.Println(ui.Yellow("Press Enter to continue..."))
    fmt.Scanln()
}
