package ui

const (
    colorReset  = "\033[0m"
    colorCyan   = "\033[36m"
    colorGreen  = "\033[32m"
    colorYellow = "\033[33m"
    colorRed    = "\033[31m"
)

func Cyan(text string) string {
    return colorCyan + text + colorReset
}

func Green(text string) string {
    return colorGreen + text + colorReset
}

func Yellow(text string) string {
    return colorYellow + text + colorReset
}

func Red(text string) string {
    return colorRed + text + colorReset
}
