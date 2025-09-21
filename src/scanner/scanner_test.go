package scanner

import (
        "os"
        "path/filepath"
        "testing"

        "metadata-remover/src/logger"
)

func setupTestEnvironment(t *testing.T) (string, *logger.Logger, func()) {
        // Create a temporary directory for testing
        tempDir, err := os.MkdirTemp("", "scanner_test")
        if err != nil {
                t.Fatalf("Failed to create temp directory: %v", err)
        }

        // Create a log file
        logPath := filepath.Join(tempDir, "test.log")
        log, err := logger.NewLogger(logPath)
        if err != nil {
                os.RemoveAll(tempDir)
                t.Fatalf("Failed to create logger: %v", err)
        }

        // Create sample files for testing
        createTestFiles(t, tempDir)

        // Return cleanup function
        cleanup := func() {
                log.Close()
                os.RemoveAll(tempDir)
        }

        return tempDir, log, cleanup
}

func createTestFiles(t *testing.T, dir string) {
        // Create sub-directories
        subDir := filepath.Join(dir, "subdir")
        err := os.Mkdir(subDir, 0755)
        if err != nil {
                t.Fatalf("Failed to create subdirectory: %v", err)
        }

        // Create test files with different extensions
        testFiles := []struct {
                path    string
                content []byte
        }{
                {
                        path:    filepath.Join(dir, "test.jpg"),
                        content: []byte("JPEG test file content"),
                },
                {
                        path:    filepath.Join(dir, "test.pdf"),
                        content: []byte("%PDF-1.5\nSample PDF content"),
                },
                {
                        path:    filepath.Join(dir, "test.docx"),
                        content: []byte("DOCX test file content"),
                },
                {
                        path:    filepath.Join(subDir, "nested.png"),
                        content: []byte("PNG test file content"),
                },
        }

        for _, tf := range testFiles {
                err := os.WriteFile(tf.path, tf.content, 0644)
                if err != nil {
                        t.Fatalf("Failed to create test file %s: %v", tf.path, err)
                }
        }
}

func TestNewScanner(t *testing.T) {
        _, log, cleanup := setupTestEnvironment(t)
        defer cleanup()

        // Test scanner creation with different options
        testCases := []struct {
                name       string
                previewMode bool
                verbose    bool
        }{
                {
                        name:       "Default mode",
                        previewMode: false,
                        verbose:    false,
                },
                {
                        name:       "Preview mode",
                        previewMode: true,
                        verbose:    false,
                },
                {
                        name:       "Verbose mode",
                        previewMode: false,
                        verbose:    true,
                },
                {
                        name:       "Preview and verbose mode",
                        previewMode: true,
                        verbose:    true,
                },
        }

        for _, tc := range testCases {
                t.Run(tc.name, func(t *testing.T) {
                        scanner := NewScanner(log, tc.previewMode, tc.verbose)
                        if scanner == nil {
                                t.Error("Expected scanner to be created")
                        }
                        if scanner.previewMode != tc.previewMode {
                                t.Errorf("Expected previewMode to be %v, got %v", tc.previewMode, scanner.previewMode)
                        }
                        if scanner.verbose != tc.verbose {
                                t.Errorf("Expected verbose to be %v, got %v", tc.verbose, scanner.verbose)
                        }
                        if scanner.logger != log {
                                t.Error("Expected logger to be set correctly")
                        }
                        if scanner.processor == nil {
                                t.Error("Expected processor to be created")
                        }
                })
        }
}

func TestScanDirectory(t *testing.T) {
        tempDir, log, cleanup := setupTestEnvironment(t)
        defer cleanup()

        testCases := []struct {
                name       string
                recursive  bool
                previewMode bool
                minFileCount   int // Minimum number of files that should be found
        }{
                {
                        name:       "Non-recursive scan",
                        recursive:  false,
                        previewMode: false,
                        minFileCount:   3, // Files in the root directory
                },
                {
                        name:       "Recursive scan",
                        recursive:  true,
                        previewMode: false,
                        minFileCount:   4, // All files including subdirectories
                },
                {
                        name:       "Preview mode",
                        recursive:  true,
                        previewMode: true,
                        minFileCount:   4, // All files including subdirectories
                },
        }

        for _, tc := range testCases {
                t.Run(tc.name, func(t *testing.T) {
                        scanner := NewScanner(log, tc.previewMode, true)
                        fileCount, processedCount, err := scanner.ScanDirectory(tempDir, tc.recursive)
                        
                        if err != nil {
                                t.Errorf("Expected no error, got %v", err)
                        }
                        
                        if fileCount < tc.minFileCount {
                                t.Errorf("Expected at least %d files, got %d", tc.minFileCount, fileCount)
                        }
                        
                        // In preview mode, files are counted but not actually processed
                        if tc.previewMode {
                                // Test log files are unsupported even in preview mode, so it's normal to have fewer processed files
                                if processedCount < fileCount-1 {
                                        t.Errorf("In preview mode, most files should be marked as processed: expected at least %d, got %d", fileCount-1, processedCount)
                                }
                        } else {
                                // Some processing may fail in real mode, but we should process at least some files
                                if processedCount == 0 && fileCount > 0 {
                                        t.Errorf("Expected some files to be processed, got %d of %d", processedCount, fileCount)
                                }
                        }
                })
        }
}

func TestProcessFile(t *testing.T) {
        tempDir, log, cleanup := setupTestEnvironment(t)
        defer cleanup()

        // Test file scenarios
        testCases := []struct {
                name       string
                filePath   string
                previewMode bool
                shouldError bool
        }{
                {
                        name:       "Process JPEG file",
                        filePath:   filepath.Join(tempDir, "test.jpg"),
                        previewMode: false,
                        shouldError: true, // Test files don't have valid headers, so they'll fail validation
                },
                {
                        name:       "Process PDF file",
                        filePath:   filepath.Join(tempDir, "test.pdf"),
                        previewMode: false,
                        shouldError: false,
                },
                {
                        name:       "Process non-existent file",
                        filePath:   filepath.Join(tempDir, "nonexistent.txt"),
                        previewMode: false,
                        shouldError: true,
                },
                {
                        name:       "Preview mode",
                        filePath:   filepath.Join(tempDir, "test.jpg"),
                        previewMode: true,
                        shouldError: false,
                },
        }

        for _, tc := range testCases {
                t.Run(tc.name, func(t *testing.T) {
                        scanner := NewScanner(log, tc.previewMode, true)
                        err := scanner.ProcessFile(tc.filePath)
                        
                        if tc.shouldError && err == nil {
                                t.Error("Expected error, got none")
                        }
                        if !tc.shouldError && err != nil {
                                t.Errorf("Expected no error, got %v", err)
                        }
                })
        }
}