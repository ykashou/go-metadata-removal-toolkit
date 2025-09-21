package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestGetFileExtension(t *testing.T) {
	testCases := []struct {
		name     string
		filePath string
		expected string
	}{
		{
			name:     "Regular file with extension",
			filePath: "document.docx",
			expected: ".docx",
		},
		{
			name:     "File with multiple dots",
			filePath: "archive.tar.gz",
			expected: ".gz",
		},
		{
			name:     "File with no extension",
			filePath: "README",
			expected: "",
		},
		{
			name:     "File with path",
			filePath: "/path/to/image.jpg",
			expected: ".jpg",
		},
		{
			name:     "Hidden file with extension",
			filePath: ".config.json",
			expected: ".json",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GetFileExtension(tc.filePath)
			if result != tc.expected {
				t.Errorf("Expected extension %q, but got %q", tc.expected, result)
			}
		})
	}
}

func TestIsDirectory(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "fileutils_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a temporary file
	tempFile, err := os.CreateTemp(tempDir, "test_file")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempFile.Close()

	testCases := []struct {
		name     string
		path     string
		expected bool
		hasError bool
	}{
		{
			name:     "Directory path",
			path:     tempDir,
			expected: true,
			hasError: false,
		},
		{
			name:     "File path",
			path:     tempFile.Name(),
			expected: false,
			hasError: false,
		},
		{
			name:     "Non-existent path",
			path:     filepath.Join(tempDir, "non-existent"),
			expected: false,
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := IsDirectory(tc.path)
			if (err != nil) != tc.hasError {
				t.Errorf("Expected error: %v, but got: %v", tc.hasError, err)
			}
			if !tc.hasError && result != tc.expected {
				t.Errorf("Expected isDirectory = %v, but got %v", tc.expected, result)
			}
		})
	}
}

func TestCreateBackup(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "backup_test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Write some data to the file
	testData := []byte("This is test data for backup")
	if err := os.WriteFile(tempFile.Name(), testData, 0644); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// Create backup
	backupPath, err := CreateBackup(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to create backup: %v", err)
	}
	defer os.Remove(backupPath)

	// Verify backup file exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		t.Errorf("Backup file does not exist")
	}

	// Verify backup contains the same data
	backupData, err := os.ReadFile(backupPath)
	if err != nil {
		t.Fatalf("Failed to read backup file: %v", err)
	}
	if string(backupData) != string(testData) {
		t.Errorf("Backup data does not match original data")
	}

	// Verify backup name format
	if !strings.HasPrefix(backupPath, tempFile.Name()+".bak.") {
		t.Errorf("Backup file name doesn't match expected format, got: %s", backupPath)
	}
}

func TestGetFileHash(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "hash_test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Test with empty file
	hash1, err := GetFileHash(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to get hash for empty file: %v", err)
	}
	if hash1 == "" {
		t.Error("Expected non-empty hash for empty file")
	}

	// Test with content
	testData := []byte("This is test data for hashing")
	if err := os.WriteFile(tempFile.Name(), testData, 0644); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	hash2, err := GetFileHash(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to get hash for file with content: %v", err)
	}
	if hash2 == "" {
		t.Error("Expected non-empty hash for file with content")
	}
	if hash1 == hash2 {
		t.Error("Expected different hashes for different content")
	}

	// Test with non-existent file
	_, err = GetFileHash(tempFile.Name() + "_nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestGetFileSize(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "size_test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Test with empty file
	size1, err := GetFileSize(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to get size for empty file: %v", err)
	}
	if size1 != 0 {
		t.Errorf("Expected size 0 for empty file, got %d", size1)
	}

	// Test with content
	testData := []byte("This is test data for size calculation")
	if err := os.WriteFile(tempFile.Name(), testData, 0644); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	size2, err := GetFileSize(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to get size for file with content: %v", err)
	}
	if size2 != int64(len(testData)) {
		t.Errorf("Expected size %d, got %d", len(testData), size2)
	}

	// Test with non-existent file
	_, err = GetFileSize(tempFile.Name() + "_nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}

func TestGetFileModTime(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "modtime_test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	// Get current time before modification
	beforeMod := time.Now().Add(-1 * time.Second)

	// Write to file to update modification time
	testData := []byte("This is test data for modification time")
	if err := os.WriteFile(tempFile.Name(), testData, 0644); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// Get modification time
	modTime, err := GetFileModTime(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to get modification time: %v", err)
	}

	// Verify modification time is after our before time
	if modTime.Before(beforeMod) {
		t.Errorf("Expected modTime after %v, got %v", beforeMod, modTime)
	}

	// Test with non-existent file
	_, err = GetFileModTime(tempFile.Name() + "_nonexistent")
	if err == nil {
		t.Error("Expected error for non-existent file")
	}
}
