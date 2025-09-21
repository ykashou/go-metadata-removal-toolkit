package utils

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// GetFileHash returns the SHA-256 hash of a file
func GetFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hasher.Sum(nil)), nil
}

// GetFileSize returns the size of a file in bytes
func GetFileSize(filePath string) (int64, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return fileInfo.Size(), nil
}

// GetFileModTime returns the modification time of a file
func GetFileModTime(filePath string) (time.Time, error) {
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return time.Time{}, err
	}
	return fileInfo.ModTime(), nil
}

// CreateBackup creates a backup of a file
func CreateBackup(filePath string) (string, error) {
	// Read original file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	// Create backup file name
	backupPath := fmt.Sprintf("%s.bak.%s", filePath, time.Now().Format("20060102_150405"))

	// Write backup file
	err = os.WriteFile(backupPath, data, 0644)
	if err != nil {
		return "", err
	}

	return backupPath, nil
}

// IsDirectory checks if a path is a directory
func IsDirectory(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return fileInfo.IsDir(), nil
}

// GetFileExtension returns the extension of a file
func GetFileExtension(filePath string) string {
	return filepath.Ext(filePath)
}
