package scanner

import (
        "fmt"
        "io/fs"
        "os"
        "path/filepath"
        "strings"

        "metadata-remover/logger"
        "metadata-remover/processor"
        "metadata-remover/stats"
        "metadata-remover/utils"
)

// Scanner handles the file scanning and processing
type Scanner struct {
        logger      *logger.Logger
        previewMode bool
        verbose     bool
        processor   *processor.Processor
}

// NewScanner creates a new scanner instance
func NewScanner(logger *logger.Logger, previewMode, verbose bool) *Scanner {
        return &Scanner{
                logger:      logger,
                previewMode: previewMode,
                verbose:     verbose,
                processor:   processor.NewProcessor(logger, previewMode),
        }
}

// ScanDirectory recursively scans a directory and processes files
func (s *Scanner) ScanDirectory(dirPath string, recursive bool) (int, int, error) {
        fileCount := 0
        processedCount := 0

        walkFunc := func(path string, info fs.FileInfo, err error) error {
                if err != nil {
                        s.logger.Error("Error accessing path %s: %v", path, err)
                        utils.PrintError(fmt.Sprintf("Error accessing path %s: %v", path, err))
                        return nil // Continue walking even if there's an error accessing a path
                }

                // Skip directories in non-recursive mode
                if !recursive && info.IsDir() && path != dirPath {
                        return filepath.SkipDir
                }

                // Skip directories and process files
                if !info.IsDir() {
                        fileCount++
                        err := s.ProcessFile(path)
                        if err == nil {
                                processedCount++
                        }
                }
                return nil
        }

        // Start walking
        err := filepath.Walk(dirPath, walkFunc)
        if err != nil {
                return fileCount, processedCount, err
        }

        return fileCount, processedCount, nil
}

// ProcessFile processes a single file
func (s *Scanner) ProcessFile(filePath string) error {
        // Get file info
        _, err := os.Stat(filePath)
        if err != nil {
                s.logger.Error("Error accessing file %s: %v", filePath, err)
                if s.verbose {
                        utils.PrintError(fmt.Sprintf("Error accessing file %s: %v", filePath, err))
                }
                return err
        }

        // Get file extension and check if supported
        ext := strings.ToLower(filepath.Ext(filePath))
        
        if s.verbose {
                utils.PrintInfo(fmt.Sprintf("Processing file: %s", filePath))
        }

        // Process file based on extension
        err = s.processor.ProcessFile(filePath, ext)
        if err != nil {
                s.logger.Error("Error processing file %s: %v", filePath, err)
                if s.verbose {
                        utils.PrintError(fmt.Sprintf("Error processing file %s: %v", filePath, err))
                }
                return err
        }

        return nil
}

// GetStats returns the metadata statistics collected during processing
func (s *Scanner) GetStats() *stats.MetadataStats {
        return s.processor.Stats
}
