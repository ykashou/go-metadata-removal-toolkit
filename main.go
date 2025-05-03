package main

import (
        "flag"
        "fmt"
        "os"
        "path/filepath"
        "time"

        "metadata-remover/logger"
        "metadata-remover/scanner"
        "metadata-remover/utils"
)

var (
        dirPath      string
        recursive    bool
        previewMode  bool
        verboseMode  bool
        outputFormat string
        version      bool
)

const (
        appVersion = "1.0.0"
)

func init() {
        flag.StringVar(&dirPath, "path", ".", "Path to directory or file to process")
        flag.BoolVar(&recursive, "recursive", false, "Recursively process subdirectories")
        flag.BoolVar(&previewMode, "preview", false, "Preview mode (no actual changes)")
        flag.BoolVar(&verboseMode, "verbose", false, "Verbose output")
        flag.StringVar(&outputFormat, "output", "terminal", "Output format (terminal, json)")
        flag.BoolVar(&version, "version", false, "Show version information")

        // Add aliases for flags
        flag.StringVar(&dirPath, "p", ".", "Path to directory or file to process (shorthand)")
        flag.BoolVar(&recursive, "r", false, "Recursively process subdirectories (shorthand)")
        flag.BoolVar(&previewMode, "preview-only", false, "Preview mode (no actual changes)")
        flag.BoolVar(&verboseMode, "v", false, "Verbose output (shorthand)")
}

func main() {
        flag.Parse()

        if version {
                fmt.Printf("go-metadata-removal-utility v%s\n", appVersion)
                os.Exit(0)
        }

        // Check if path exists
        _, err := os.Stat(dirPath)
        if err != nil {
                utils.PrintError(fmt.Sprintf("Error accessing path %s: %v", dirPath, err))
                os.Exit(1)
        }

        // Create logger
        logFileName := fmt.Sprintf("metadata_removal_%s.log", time.Now().Format("20060102_150405"))
        logFilePath := filepath.Join(".", logFileName)
        log, err := logger.NewLogger(logFilePath)
        if err != nil {
                utils.PrintError(fmt.Sprintf("Error creating log file: %v", err))
                os.Exit(1)
        }
        defer log.Close()

        // Initialize scanner
        s := scanner.NewScanner(log, previewMode, verboseMode)

        // Print initial information
        utils.PrintInfo(fmt.Sprintf("Starting metadata removal utility"))
        utils.PrintInfo(fmt.Sprintf("Path: %s", dirPath))
        utils.PrintInfo(fmt.Sprintf("Recursive mode: %v", recursive))
        utils.PrintInfo(fmt.Sprintf("Preview mode: %v", previewMode))
        utils.PrintInfo(fmt.Sprintf("Log file: %s", logFilePath))
        utils.PrintInfo("")

        // Process files
        startTime := time.Now()

        // Get file information
        fileInfo, err := os.Stat(dirPath)
        if err != nil {
                utils.PrintError(fmt.Sprintf("Error accessing path %s: %v", dirPath, err))
                os.Exit(1)
        }

        var fileCount, processedCount int
        // Check if path is a file or directory
        if fileInfo.IsDir() {
                fileCount, processedCount, err = s.ScanDirectory(dirPath, recursive)
        } else {
                // Single file mode
                err = s.ProcessFile(dirPath)
                if err == nil {
                        fileCount = 1
                        processedCount = 1
                }
        }

        if err != nil {
                utils.PrintError(fmt.Sprintf("Error during processing: %v", err))
                os.Exit(1)
        }

        // Print summary
        duration := time.Since(startTime)
        utils.PrintInfo("")
        utils.PrintSuccess(fmt.Sprintf("Processing complete!"))
        utils.PrintSuccess(fmt.Sprintf("Files scanned: %d", fileCount))
        utils.PrintSuccess(fmt.Sprintf("Files processed: %d", processedCount))
        utils.PrintSuccess(fmt.Sprintf("Time taken: %v", duration))
        utils.PrintSuccess(fmt.Sprintf("Log file: %s", logFilePath))

        // Print detailed metadata statistics
        if processedCount > 0 {
                // Get statistics from the scanner
                metadataStats := s.GetStats()
                
                // Format and display stats based on output format
                switch outputFormat {
                case "json":
                        fmt.Println(utils.FormatStatsAsJSON(metadataStats))
                default:
                        fmt.Println(utils.FormatStatsAsText(metadataStats))
                }
        }

        if previewMode {
                utils.PrintWarning("Preview mode: No files were modified")
        }
}
