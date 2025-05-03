package processor

import (
        "errors"
        "fmt"
        "path/filepath"
        "strings"

        "metadata-remover/logger"
        "metadata-remover/utils"
)

// Processor handles metadata removal from various file types
type Processor struct {
        logger      *logger.Logger
        previewMode bool
}

// File type constants
const (
        TypeImage    = "image"
        TypePDF      = "pdf"
        TypeDocument = "document"
        TypeUnknown  = "unknown"
)

// NewProcessor creates a new processor
func NewProcessor(logger *logger.Logger, previewMode bool) *Processor {
        return &Processor{
                logger:      logger,
                previewMode: previewMode,
        }
}

// ProcessFile processes a file based on its extension
func (p *Processor) ProcessFile(filePath, ext string) error {
        // Determine file type based on extension
        fileType := p.getFileType(ext)
        
        if fileType == TypeUnknown {
                p.logger.Warning("Unsupported file type: %s", ext)
                utils.PrintWarning(fmt.Sprintf("Unsupported file type: %s (skipping %s)", ext, filepath.Base(filePath)))
                return errors.New("unsupported file type")
        }

        // Process file based on type
        var err error
        switch fileType {
        case TypeImage:
                err = p.ProcessImage(filePath, ext)
        case TypePDF:
                err = p.ProcessPDF(filePath)
        case TypeDocument:
                err = p.ProcessDocument(filePath, ext)
        }

        if err != nil {
                return err
        }

        if p.previewMode {
                p.logger.Info("Preview mode: Would process %s", filePath)
                utils.PrintInfo(fmt.Sprintf("Preview mode: Would remove metadata from %s", filepath.Base(filePath)))
        } else {
                p.logger.Success("Successfully removed metadata from %s", filePath)
                utils.PrintSuccess(fmt.Sprintf("Successfully removed metadata from %s", filepath.Base(filePath)))
        }

        return nil
}

// getFileType determines the file type based on extension
func (p *Processor) getFileType(ext string) string {
        ext = strings.ToLower(ext)
        
        // Image file extensions
        imageExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".tif", ".webp"}
        for _, imgExt := range imageExtensions {
                if ext == imgExt {
                        return TypeImage
                }
        }

        // PDF file extension
        if ext == ".pdf" {
                return TypePDF
        }

        // Document file extensions
        docExtensions := []string{".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx", ".odt", ".ods", ".odp", ".rtf", ".txt"}
        for _, docExt := range docExtensions {
                if ext == docExt {
                        return TypeDocument
                }
        }

        return TypeUnknown
}
