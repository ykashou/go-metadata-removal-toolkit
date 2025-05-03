package utils

import (
	"fmt"
)

// Color codes for terminal output
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m" // Error
	colorGreen  = "\033[32m" // Success
	colorYellow = "\033[33m" // Warning
	colorWhite  = "\033[37m" // Info
)

// PrintSuccess prints a success message in green
func PrintSuccess(message string) {
	fmt.Printf("%s%s%s\n", colorGreen, message, colorReset)
}

// PrintWarning prints a warning message in yellow
func PrintWarning(message string) {
	fmt.Printf("%s%s%s\n", colorYellow, message, colorReset)
}

// PrintError prints an error message in red
func PrintError(message string) {
	fmt.Printf("%s%s%s\n", colorRed, message, colorReset)
}

// PrintInfo prints an info message in white
func PrintInfo(message string) {
	fmt.Printf("%s%s%s\n", colorWhite, message, colorReset)
}
