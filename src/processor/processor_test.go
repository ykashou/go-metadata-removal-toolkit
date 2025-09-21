package processor

import (
	"os"
	"path/filepath"
	"testing"

	"metadata-remover/src/logger"
)

func setupProcessorTest(t *testing.T) (string, *logger.Logger, *Processor, func()) {
	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", "processor_test")
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

	// Create processor instances
	proc := NewProcessor(log, false)

	// Return cleanup function
	cleanup := func() {
		log.Close()
		os.RemoveAll(tempDir)
	}

	return tempDir, log, proc, cleanup
}

func TestNewProcessor(t *testing.T) {
	_, log, _, cleanup := setupProcessorTest(t)
	defer cleanup()

	// Test processor creation with different options
	testCases := []struct {
		name       string
		previewMode bool
	}{
		{
			name:       "Normal mode",
			previewMode: false,
		},
		{
			name:       "Preview mode",
			previewMode: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			proc := NewProcessor(log, tc.previewMode)
			if proc == nil {
				t.Fatal("Failed to create processor")
			}
			if proc.previewMode != tc.previewMode {
				t.Errorf("Expected previewMode to be %v, got %v", tc.previewMode, proc.previewMode)
			}
			if proc.logger != log {
				t.Error("Expected logger to be set correctly")
			}
		})
	}
}

func TestGetFileType(t *testing.T) {
	_, _, proc, cleanup := setupProcessorTest(t)
	defer cleanup()

	testCases := []struct {
		name     string
		ext      string
		expected string
	}{
		{
			name:     "JPEG Image",
			ext:      ".jpg",
			expected: TypeImage,
		},
		{
			name:     "PNG Image",
			ext:      ".png",
			expected: TypeImage,
		},
		{
			name:     "PDF Document",
			ext:      ".pdf",
			expected: TypePDF,
		},
		{
			name:     "Word Document",
			ext:      ".docx",
			expected: TypeDocument,
		},
		{
			name:     "Excel Document",
			ext:      ".xlsx",
			expected: TypeDocument,
		},
		{
			name:     "Text Document",
			ext:      ".txt",
			expected: TypeDocument,
		},
		{
			name:     "Unsupported File Type",
			ext:      ".xyz",
			expected: TypeUnknown,
		},
		{
			name:     "Case Insensitivity",
			ext:      ".JPG",
			expected: TypeImage,
		},
		{
			name:     "Empty Extension",
			ext:      "",
			expected: TypeUnknown,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := proc.getFileType(tc.ext)
			if result != tc.expected {
				t.Errorf("Expected file type %q, got %q", tc.expected, result)
			}
		})
	}
}

func TestProcessFile(t *testing.T) {
	tempDir, _, proc, cleanup := setupProcessorTest(t)
	defer cleanup()

	// Create test files
	testFiles := []struct {
		path     string
		content  []byte
		ext      string
		fileType string
	}{
		{
			path:     filepath.Join(tempDir, "test.jpg"),
			content:  []byte("JPEG image data"),
			ext:      ".jpg",
			fileType: TypeImage,
		},
		{
			path:     filepath.Join(tempDir, "test.pdf"),
			content:  []byte("%PDF-1.5\nTest content"),
			ext:      ".pdf",
			fileType: TypePDF,
		},
		{
			path:     filepath.Join(tempDir, "test.docx"),
			content:  []byte("DOCX document content"),
			ext:      ".docx",
			fileType: TypeDocument,
		},
		{
			path:     filepath.Join(tempDir, "test.xyz"),
			content:  []byte("Unknown file type"),
			ext:      ".xyz",
			fileType: TypeUnknown,
		},
	}

	// Create the test files
	for _, tf := range testFiles {
		err := os.WriteFile(tf.path, tf.content, 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", tf.path, err)
		}
	}

	// Test processing files
	for _, tf := range testFiles {
		t.Run("Process "+tf.fileType, func(t *testing.T) {
			err := proc.ProcessFile(tf.path, tf.ext)
			
			if tf.fileType == TypeUnknown {
				// Should error for unknown file types
				if err == nil {
					t.Error("Expected error for unknown file type")
				}
			} else {
				// Preview mode may not actually process the file
				if proc.previewMode && err != nil {
					t.Errorf("Expected no error in preview mode, got %v", err)
				}
			}
		})
	}

	// Test preview mode
	previewProc := NewProcessor(proc.logger, true)
	for _, tf := range testFiles {
		if tf.fileType == TypeUnknown {
			continue // Skip unknown file types
		}

		t.Run("Preview "+tf.fileType, func(t *testing.T) {
			err := previewProc.ProcessFile(tf.path, tf.ext)
			if err != nil {
				t.Errorf("Expected no error in preview mode, got %v", err)
			}
			
			// Verify file content was not changed in preview mode
			content, err := os.ReadFile(tf.path)
			if err != nil {
				t.Fatalf("Failed to read file: %v", err)
			}
			
			if string(content) != string(tf.content) {
				t.Error("File was modified in preview mode")
			}
		})
	}
}