package logger

import (
	"os"
	"strings"
	"testing"
)

func TestLoggerCreation(t *testing.T) {
	tempFile, err := os.CreateTemp("", "logger_test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempPath := tempFile.Name()
	tempFile.Close()
	os.Remove(tempPath) // Remove so the logger can create it

	// Create logger
	logger, err := NewLogger(tempPath)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer func() {
		logger.Close()
		os.Remove(tempPath)
	}()

	// Verify log file exists
	if _, err := os.Stat(tempPath); os.IsNotExist(err) {
		t.Errorf("Log file was not created")
	}

	// Verify header was written
	content, err := os.ReadFile(tempPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	if !strings.Contains(string(content), "# Metadata Removal Utility Log") {
		t.Errorf("Log file doesn't contain expected header")
	}
	if !strings.Contains(string(content), "# Started:") {
		t.Errorf("Log file doesn't contain timestamp")
	}
}

func TestLogLevels(t *testing.T) {
	tempFile, err := os.CreateTemp("", "logger_test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempPath := tempFile.Name()
	tempFile.Close()
	os.Remove(tempPath) // Remove so the logger can create it

	// Create logger
	logger, err := NewLogger(tempPath)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}
	defer func() {
		logger.Close()
		os.Remove(tempPath)
	}()

	// Test each log level
	testCases := []struct {
		name     string
		logFunc  func(string, ...interface{}) error
		message  string
		expected string
	}{
		{
			name:     "Info level",
			logFunc:  logger.Info,
			message:  "Test info message",
			expected: "[INFO] Test info message",
		},
		{
			name:     "Success level",
			logFunc:  logger.Success,
			message:  "Test success message",
			expected: "[SUCCESS] Test success message",
		},
		{
			name:     "Warning level",
			logFunc:  logger.Warning,
			message:  "Test warning message",
			expected: "[WARNING] Test warning message",
		},
		{
			name:     "Error level",
			logFunc:  logger.Error,
			message:  "Test error message",
			expected: "[ERROR] Test error message",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.logFunc(tc.message)
			if err != nil {
				t.Fatalf("Failed to log message: %v", err)
			}
		})
	}

	// Read log file and verify all log levels were recorded
	content, err := os.ReadFile(tempPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}

	for _, tc := range testCases {
		if !strings.Contains(string(content), tc.expected) {
			t.Errorf("Log file doesn't contain %s message", strings.Fields(tc.name)[0])
		}
	}
}

func TestLoggerClose(t *testing.T) {
	tempFile, err := os.CreateTemp("", "logger_test")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	tempPath := tempFile.Name()
	tempFile.Close()
	os.Remove(tempPath) // Remove so the logger can create it

	// Create logger
	logger, err := NewLogger(tempPath)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// Write a log message
	err = logger.Info("Test before close")
	if err != nil {
		t.Fatalf("Failed to log message: %v", err)
	}

	// Close logger
	err = logger.Close()
	if err != nil {
		t.Fatalf("Failed to close logger: %v", err)
	}
	defer os.Remove(tempPath)

	// Verify footer was written
	content, err := os.ReadFile(tempPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
	}
	if !strings.Contains(string(content), "# Finished:") {
		t.Errorf("Log file doesn't contain finished timestamp")
	}

	// Verify logging after closing fails
	err = logger.Info("Test after close")
	if err == nil {
		t.Error("Expected error when logging after close")
	}
}
