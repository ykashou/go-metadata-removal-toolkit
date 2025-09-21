package processor

import (
	"os"
	"path/filepath"
	"testing"

	"metadata-remover/src/logger"
)

func setupImageTest(t *testing.T) (string, *Processor, func()) {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "image_processor_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Create logger
	logPath := filepath.Join(tempDir, "image_test.log")
	log, err := logger.NewLogger(logPath)
	if err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Create processor
	proc := NewProcessor(log, false)

	// Return cleanup function
	cleanup := func() {
		log.Close()
		os.RemoveAll(tempDir)
	}

	return tempDir, proc, cleanup
}

func TestProcessImage(t *testing.T) {
	tempDir, proc, cleanup := setupImageTest(t)
	defer cleanup()

	// Test cases for different image formats
	testCases := []struct {
		name       string
		ext        string
		shouldFail bool
	}{
		{
			name:       "JPEG image",
			ext:        ".jpg",
			shouldFail: false,
		},
		{
			name:       "PNG image",
			ext:        ".png",
			shouldFail: false,
		},
		{
			name:       "GIF image",
			ext:        ".gif",
			shouldFail: false,
		},
		{
			name:       "BMP image",
			ext:        ".bmp",
			shouldFail: false,
		},
		{
			name:       "TIFF image",
			ext:        ".tiff",
			shouldFail: false,
		},
		{
			name:       "WebP image",
			ext:        ".webp",
			shouldFail: false,
		},
		{
			name:       "Unsupported format",
			ext:        ".xyz",
			shouldFail: true,
		},
	}

	// Test each format in preview mode (shouldn't fail)
	previewProc := NewProcessor(proc.logger, true)
	for _, tc := range testCases {
		t.Run("Preview "+tc.name, func(t *testing.T) {
			err := previewProc.ProcessImage("dummy_path"+tc.ext, tc.ext)
			if err != nil {
				t.Errorf("Expected no error in preview mode, got: %v", err)
			}
		})
	}

	// Create dummy files for testing
	for _, tc := range testCases {
		if tc.shouldFail {
			continue
		}

		// Create dummy file with proper header for testing
		filePath := filepath.Join(tempDir, "test"+tc.ext)
		var fileContent []byte

		switch tc.ext {
		case ".jpg", ".jpeg":
			fileContent = []byte{0xFF, 0xD8, 0xFF, 0xE0, 0x00, 0x10, 0x4A, 0x46, 0x49, 0x46, 0x00}
		case ".png":
			fileContent = []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
		case ".gif":
			fileContent = []byte("GIF89a")
		case ".bmp":
			fileContent = []byte{'B', 'M'}
		case ".tiff", ".tif":
			fileContent = []byte{0x49, 0x49, 0x2A, 0x00}
		case ".webp":
			fileContent = []byte{'R', 'I', 'F', 'F', 0x00, 0x00, 0x00, 0x00, 'W', 'E', 'B', 'P'}
		}

		err := os.WriteFile(filePath, fileContent, 0644)
		if err != nil {
			t.Fatalf("Failed to create test file %s: %v", filePath, err)
		}

		// Test actual processing
		t.Run("Process "+tc.name, func(t *testing.T) {
			err := proc.ProcessImage(filePath, tc.ext)
			if err != nil {
				t.Errorf("Expected no error, got: %v", err)
			}
		})
	}

	// Test with invalid files
	t.Run("Invalid JPEG", func(t *testing.T) {
		invalidPath := filepath.Join(tempDir, "invalid.jpg")
		err := os.WriteFile(invalidPath, []byte("Not a JPEG file"), 0644)
		if err != nil {
			t.Fatalf("Failed to create invalid test file: %v", err)
		}

		err = proc.ProcessImage(invalidPath, ".jpg")
		if err == nil {
			t.Error("Expected error for invalid JPEG file")
		}
	})

	t.Run("Invalid PNG", func(t *testing.T) {
		invalidPath := filepath.Join(tempDir, "invalid.png")
		err := os.WriteFile(invalidPath, []byte("Not a PNG file"), 0644)
		if err != nil {
			t.Fatalf("Failed to create invalid test file: %v", err)
		}

		err = proc.ProcessImage(invalidPath, ".png")
		if err == nil {
			t.Error("Expected error for invalid PNG file")
		}
	})

	t.Run("Nonexistent file", func(t *testing.T) {
		err := proc.ProcessImage(filepath.Join(tempDir, "nonexistent.jpg"), ".jpg")
		if err == nil {
			t.Error("Expected error for nonexistent file")
		}
	})
}