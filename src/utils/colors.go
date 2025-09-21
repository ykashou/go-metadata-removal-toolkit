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
	colorBlue   = "\033[34m" // Heading
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

// Green colorizes a string in green
func Green(text string) string {
	return colorGreen + text + colorReset
}

// Red colorizes a string in red
func Red(text string) string {
	return colorRed + text + colorReset
}

// Yellow colorizes a string in yellow
func Yellow(text string) string {
	return colorYellow + text + colorReset
}

// Blue colorizes a string in blue
func Blue(text string) string {
	return colorBlue + text + colorReset
}

// White colorizes a string in white
func White(text string) string {
	return colorWhite + text + colorReset
}
