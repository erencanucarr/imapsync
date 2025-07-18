package ui

const (
	colorReset  = "\033[0m"
	colorCyan   = "\033[36m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorRed    = "\033[31m"
	colorBlue   = "\033[34m"
	colorPurple = "\033[35m"
	colorWhite  = "\033[37m"
	colorBold   = "\033[1m"
	colorDim    = "\033[2m"
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

func Blue(text string) string {
	return colorBlue + text + colorReset
}

func Purple(text string) string {
	return colorPurple + text + colorReset
}

func White(text string) string {
	return colorWhite + text + colorReset
}

func Bold(text string) string {
	return colorBold + text + colorReset
}

func Dim(text string) string {
	return colorDim + text + colorReset
}
