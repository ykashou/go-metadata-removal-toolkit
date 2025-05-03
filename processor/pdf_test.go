package processor

import (
        "os"
        "path/filepath"
        "testing"

        "metadata-remover/logger"
)

func setupPDFTest(t *testing.T) (string, *Processor, func()) {
        // Create temporary directory
        tempDir, err := os.MkdirTemp("", "pdf_processor_test")
        if err != nil {
                t.Fatalf("Failed to create temp directory: %v", err)
        }

        // Create logger
        logPath := filepath.Join(tempDir, "pdf_test.log")
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

func TestProcessPDF(t *testing.T) {
        tempDir, proc, cleanup := setupPDFTest(t)
        defer cleanup()

        // Create dummy PDF file with minimal header
        validPDFPath := filepath.Join(tempDir, "valid.pdf")
        err := os.WriteFile(validPDFPath, []byte("%PDF-1.5\nSample PDF content with metadata\n/Info 1 0 R\n"), 0644)
        if err != nil {
                t.Fatalf("Failed to create test PDF file: %v", err)
        }

        // Create invalid PDF file
        invalidPDFPath := filepath.Join(tempDir, "invalid.pdf")
        err = os.WriteFile(invalidPDFPath, []byte("Not a PDF file"), 0644)
        if err != nil {
                t.Fatalf("Failed to create invalid test file: %v", err)
        }

        // Test processing valid PDF file
        t.Run("Process valid PDF", func(t *testing.T) {
                err := proc.ProcessPDF(validPDFPath)
                if err != nil {
                        t.Errorf("Expected no error for valid PDF, got: %v", err)
                }

                // Verify the file still exists after processing
                if _, err := os.Stat(validPDFPath); os.IsNotExist(err) {
                        t.Error("PDF file should still exist after processing")
                }
        })

        // Test processing in preview mode
        t.Run("Preview mode", func(t *testing.T) {
                previewProc := NewProcessor(proc.logger, true)
                err := previewProc.ProcessPDF(validPDFPath)
                if err != nil {
                        t.Errorf("Expected no error in preview mode, got: %v", err)
                }
        })

        // Test processing invalid PDF file
        t.Run("Process invalid PDF", func(t *testing.T) {
                err := proc.ProcessPDF(invalidPDFPath)
                if err == nil {
                        t.Error("Expected error for invalid PDF file")
                }
        })

        // Test with nonexistent file
        t.Run("Nonexistent file", func(t *testing.T) {
                err := proc.ProcessPDF(filepath.Join(tempDir, "nonexistent.pdf"))
                if err == nil {
                        t.Error("Expected error for nonexistent file")
                }
        })
}

func TestPDFMetadataRemoval(t *testing.T) {
        tempDir, proc, cleanup := setupPDFTest(t)
        defer cleanup()

        // Create a PDF file with various types of metadata
        pdfWithMetadata := filepath.Join(tempDir, "metadata.pdf")
        pdfContent := `%PDF-1.5
1 0 obj
<< /Title (Test Document) /Author (Test Author) /Subject (Test Subject) /Keywords (test, pdf, metadata) >>
endobj
2 0 obj
<< /Info 1 0 R >>
endobj
%% Document with XMP metadata
<x:xmpmeta xmlns:x="adobe:ns:meta/">
<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">
</rdf:RDF>
</x:xmpmeta>
%%EOF`

        err := os.WriteFile(pdfWithMetadata, []byte(pdfContent), 0644)
        if err != nil {
                t.Fatalf("Failed to create test PDF with metadata: %v", err)
        }

        // Process the PDF
        err = proc.ProcessPDF(pdfWithMetadata)
        if err != nil {
                t.Fatalf("Failed to process PDF with metadata: %v", err)
        }

        // Read the processed file
        processedContent, err := os.ReadFile(pdfWithMetadata)
        if err != nil {
                t.Fatalf("Failed to read processed PDF: %v", err)
        }

        contentStr := string(processedContent)

        // Check that metadata was modified or removed
        t.Run("Info dictionary modification", func(t *testing.T) {
                if proc.previewMode {
                        // In preview mode, metadata should remain unchanged
                        return
                }

                if pdfContains(contentStr, "/Info 1 0 R") && !pdfContains(contentStr, "/Info 0 0 R") {
                        t.Error("Info reference was not replaced")
                }
        })

        t.Run("XMP metadata removal", func(t *testing.T) {
                if proc.previewMode {
                        // In preview mode, metadata should remain unchanged
                        return
                }

                if pdfContains(contentStr, "<x:xmpmeta") {
                        t.Error("XMP metadata was not removed")
                }
        })

        t.Run("Document info removal", func(t *testing.T) {
                if proc.previewMode {
                        // In preview mode, metadata should remain unchanged
                        return
                }

                if pdfContains(contentStr, "<< /Title (Test Document)") && !pdfContains(contentStr, "<< >>") {
                        t.Error("Document info dictionary was not replaced")
                }
        })
}

// pdfContains checks if a string contains a substring (PDF context)
func pdfContains(s, substr string) bool {
        if len(s) < len(substr) {
                return false
        }
        for i := 0; i <= len(s)-len(substr); i++ {
                if s[i:i+len(substr)] == substr {
                        return true
                }
        }
        return false
}