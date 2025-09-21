package main

import (
	"flag"
	"os"
	"testing"
)

// TestMainCommandLineFlags tests the command line flag parsing
func TestMainCommandLineFlags(t *testing.T) {
	// Save original flag values and command-line arguments
	originalDirPath := dirPath
	originalRecursive := recursive
	originalPreviewMode := previewMode
	originalVerboseMode := verboseMode
	originalOutputFormat := outputFormat
	originalVersion := version
	originalArgs := os.Args

	// Restore the original values after the test
	defer func() {
		dirPath = originalDirPath
		recursive = originalRecursive
		previewMode = originalPreviewMode
		verboseMode = originalVerboseMode
		outputFormat = originalOutputFormat
		version = originalVersion
		os.Args = originalArgs
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	}()

	testCases := []struct {
		name           string
		args           []string
		expectedPath   string
		expectedRec    bool
		expectedPrev   bool
		expectedVerb   bool
		expectedOutput string
		expectedVer    bool
	}{
		{
			name:           "Default values",
			args:           []string{"cmd"},
			expectedPath:   ".", // Default path
			expectedRec:    false,
			expectedPrev:   false,
			expectedVerb:   false,
			expectedOutput: "terminal",
			expectedVer:    false,
		},
		{
			name:           "Simple path",
			args:           []string{"cmd", "-path", "/test/path"},
			expectedPath:   "/test/path",
			expectedRec:    false,
			expectedPrev:   false,
			expectedVerb:   false,
			expectedOutput: "terminal",
			expectedVer:    false,
		},
		{
			name:           "Recursive flag",
			args:           []string{"cmd", "-r"},
			expectedPath:   ".",
			expectedRec:    true,
			expectedPrev:   false,
			expectedVerb:   false,
			expectedOutput: "terminal",
			expectedVer:    false,
		},
		{
			name:           "Preview mode",
			args:           []string{"cmd", "-preview"},
			expectedPath:   ".",
			expectedRec:    false,
			expectedPrev:   true,
			expectedVerb:   false,
			expectedOutput: "terminal",
			expectedVer:    false,
		},
		{
			name:           "Verbose mode",
			args:           []string{"cmd", "-v"},
			expectedPath:   ".",
			expectedRec:    false,
			expectedPrev:   false,
			expectedVerb:   true,
			expectedOutput: "terminal",
			expectedVer:    false,
		},
		{
			name:           "JSON output",
			args:           []string{"cmd", "-output", "json"},
			expectedPath:   ".",
			expectedRec:    false,
			expectedPrev:   false,
			expectedVerb:   false,
			expectedOutput: "json",
			expectedVer:    false,
		},
		{
			name:           "Version flag",
			args:           []string{"cmd", "-version"},
			expectedPath:   ".",
			expectedRec:    false,
			expectedPrev:   false,
			expectedVerb:   false,
			expectedOutput: "terminal",
			expectedVer:    true,
		},
		{
			name:           "Multiple flags",
			args:           []string{"cmd", "-p", "/some/path", "-r", "-v", "-preview"},
			expectedPath:   "/some/path",
			expectedRec:    true,
			expectedPrev:   true,
			expectedVerb:   true,
			expectedOutput: "terminal",
			expectedVer:    false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Reset flags for each test case
			dirPath = originalDirPath
			recursive = originalRecursive
			previewMode = originalPreviewMode
			verboseMode = originalVerboseMode
			outputFormat = originalOutputFormat
			version = originalVersion
			flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
			// Re-declare flags for this test
			flag.StringVar(&dirPath, "path", ".", "Path to directory or file to process")
			flag.BoolVar(&recursive, "recursive", false, "Recursively process subdirectories")
			flag.BoolVar(&previewMode, "preview", false, "Preview mode (no actual changes)")
			flag.BoolVar(&verboseMode, "verbose", false, "Verbose output")
			flag.StringVar(&outputFormat, "output", "terminal", "Output format (terminal, json)")
			flag.BoolVar(&version, "version", false, "Show version information")
			flag.StringVar(&dirPath, "p", ".", "Path to directory or file to process (shorthand)")
			flag.BoolVar(&recursive, "r", false, "Recursively process subdirectories (shorthand)")
			flag.BoolVar(&previewMode, "preview-only", false, "Preview mode (no actual changes)")
			flag.BoolVar(&verboseMode, "v", false, "Verbose output (shorthand)")

			// Set command-line arguments for this test
			os.Args = tc.args

			// Parse flags
			flag.Parse()

			// Check if flags were set correctly
			if dirPath != tc.expectedPath {
				t.Errorf("Expected dirPath to be %q, got %q", tc.expectedPath, dirPath)
			}
			if recursive != tc.expectedRec {
				t.Errorf("Expected recursive to be %v, got %v", tc.expectedRec, recursive)
			}
			if previewMode != tc.expectedPrev {
				t.Errorf("Expected previewMode to be %v, got %v", tc.expectedPrev, previewMode)
			}
			if verboseMode != tc.expectedVerb {
				t.Errorf("Expected verboseMode to be %v, got %v", tc.expectedVerb, verboseMode)
			}
			if outputFormat != tc.expectedOutput {
				t.Errorf("Expected outputFormat to be %q, got %q", tc.expectedOutput, outputFormat)
			}
			if version != tc.expectedVer {
				t.Errorf("Expected version to be %v, got %v", tc.expectedVer, version)
			}
		})
	}
}

// Since the main function has a lot of side effects (file operations, logging, etc.),
// it's difficult to test it directly without mocking or refactoring.
// In a real-world scenario, we would ideally move most of the functionality into
// separate functions that could be tested independently.
func TestMainExportedFunctions(t *testing.T) {
	// This test primarily ensures that exported functions and variables exist
	// and are accessible

	t.Run("AppVersion", func(t *testing.T) {
		if appVersion == "" {
			t.Error("appVersion should not be empty")
		}
	})

	// Similarly, we could test other exported functions if they existed
}
