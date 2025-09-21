package logger

import (
	"fmt"
	"os"
	"time"
)

// LogLevel defines the severity of log messages
type LogLevel int

const (
	INFO LogLevel = iota
	SUCCESS
	WARNING
	ERROR
)

// Logger handles logging to a file
type Logger struct {
	file *os.File
}

// NewLogger creates a new logger that writes to a file
func NewLogger(filePath string) (*Logger, error) {
	file, err := os.Create(filePath)
	if err != nil {
		return nil, err
	}

	// Write header to log file
	timestamp := time.Now().Format(time.RFC3339)
	header := fmt.Sprintf("# Metadata Removal Utility Log\n# Started: %s\n\n", timestamp)
	_, err = file.WriteString(header)
	if err != nil {
		file.Close()
		return nil, err
	}

	return &Logger{
		file: file,
	}, nil
}

// Close closes the log file
func (l *Logger) Close() error {
	if l.file != nil {
		// Write footer to log file
		timestamp := time.Now().Format(time.RFC3339)
		footer := fmt.Sprintf("\n# Finished: %s\n", timestamp)
		_, err := l.file.WriteString(footer)
		if err != nil {
			return err
		}

		return l.file.Close()
	}
	return nil
}

// log writes a message to the log file with timestamp and level
func (l *Logger) log(level LogLevel, format string, args ...interface{}) error {
	if l.file == nil {
		return fmt.Errorf("logger is not initialized")
	}

	// Format log message
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	var levelStr string

	switch level {
	case INFO:
		levelStr = "INFO"
	case SUCCESS:
		levelStr = "SUCCESS"
	case WARNING:
		levelStr = "WARNING"
	case ERROR:
		levelStr = "ERROR"
	}

	message := fmt.Sprintf(format, args...)
	logLine := fmt.Sprintf("[%s] [%s] %s\n", timestamp, levelStr, message)

	// Write to log file
	_, err := l.file.WriteString(logLine)
	return err
}

// Info logs an informational message
func (l *Logger) Info(format string, args ...interface{}) error {
	return l.log(INFO, format, args...)
}

// Success logs a success message
func (l *Logger) Success(format string, args ...interface{}) error {
	return l.log(SUCCESS, format, args...)
}

// Warning logs a warning message
func (l *Logger) Warning(format string, args ...interface{}) error {
	return l.log(WARNING, format, args...)
}

// Error logs an error message
func (l *Logger) Error(format string, args ...interface{}) error {
	return l.log(ERROR, format, args...)
}
