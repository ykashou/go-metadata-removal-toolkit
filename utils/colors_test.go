package utils

import (
        "bytes"
        "io"
        "os"
        "strings"
        "testing"
)

func TestColorFunctions(t *testing.T) {
        // Define color codes for verification
        colorCodes := map[string]string{
                "PrintSuccess": "\033[32m", // Green
                "PrintError":   "\033[31m", // Red
                "PrintWarning": "\033[33m", // Yellow
                "PrintInfo":    "\033[37m", // White
        }
        colorReset := "\033[0m"

        testCases := []struct {
                name     string
                function func(string)
                input    string
                colorCode string
        }{
                {
                        name:     "PrintSuccess",
                        function: PrintSuccess,
                        input:    "Test success message",
                        colorCode: colorCodes["PrintSuccess"],
                },
                {
                        name:     "PrintError",
                        function: PrintError,
                        input:    "Test error message",
                        colorCode: colorCodes["PrintError"],
                },
                {
                        name:     "PrintWarning",
                        function: PrintWarning,
                        input:    "Test warning message",
                        colorCode: colorCodes["PrintWarning"],
                },
                {
                        name:     "PrintInfo",
                        function: PrintInfo,
                        input:    "Test info message",
                        colorCode: colorCodes["PrintInfo"],
                },
        }

        for _, tc := range testCases {
                t.Run(tc.name, func(t *testing.T) {
                        // Redirect stdout to capture output
                        old := os.Stdout
                        r, w, _ := os.Pipe()
                        os.Stdout = w

                        // Call the function
                        tc.function(tc.input)

                        // Restore stdout
                        w.Close()
                        os.Stdout = old

                        // Read captured output
                        var buf bytes.Buffer
                        io.Copy(&buf, r)
                        output := buf.String()

                        // Check if the output contains the expected elements
                        if !strings.Contains(output, tc.colorCode) {
                                t.Errorf("Expected output to contain color code %q, but got %q", tc.colorCode, output)
                        }
                        if !strings.Contains(output, colorReset) {
                                t.Errorf("Expected output to contain color reset %q, but got %q", colorReset, output)
                        }
                        if !strings.Contains(output, tc.input) {
                                t.Errorf("Expected output to contain input %q, but got %q", tc.input, output)
                        }
                })
        }
}