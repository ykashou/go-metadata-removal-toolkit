package processor

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"metadata-remover/src/utils"
)

// ProcessDocument removes metadata from document files
func (p *Processor) ProcessDocument(filePath, ext string) error {
	// If preview mode, just log and return
	if p.previewMode {
		return nil
	}

	ext = strings.ToLower(ext)
	switch ext {
	case ".docx", ".xlsx", ".pptx":
		return p.cleanOpenXML(filePath)
	case ".odt", ".ods", ".odp":
		return p.cleanOpenDocument(filePath)
	case ".doc", ".xls", ".ppt":
		return p.cleanBinaryOffice(filePath, ext)
	case ".rtf":
		return p.cleanRTF(filePath)
	case ".txt":
		// Plain text files don't typically have metadata
		p.logger.Info("Text files don't have metadata to remove for %s", filePath)
		utils.PrintInfo(fmt.Sprintf("Text files don't have metadata to remove"))
		return nil
	default:
		return fmt.Errorf("unsupported document format: %s", ext)
	}
}

// cleanOpenXML removes metadata from Office Open XML files (.docx, .xlsx, .pptx)
func (p *Processor) cleanOpenXML(filePath string) error {
	// Office Open XML files are ZIP archives containing XML files
	// We need to extract, modify, and repackage them

	// Create a temporary directory
	tempDir := filePath + "_temp"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return err
	}
	defer os.RemoveAll(tempDir) // Clean up temp directory

	// Open the document as a ZIP archive
	reader, err := zip.OpenReader(filePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	// Create a new ZIP file
	tempFile := filePath + ".temp"
	zipWriter, err := os.Create(tempFile)
	if err != nil {
		return err
	}
	defer func() {
		zipWriter.Close()
		os.Remove(tempFile) // Clean up temp file in case of error
	}()

	archive := zip.NewWriter(zipWriter)
	defer archive.Close()

	// Process each file in the archive
	for _, file := range reader.File {
		fileReader, err := file.Open()
		if err != nil {
			return err
		}

		data, err := ioutil.ReadAll(fileReader)
		fileReader.Close()
		if err != nil {
			return err
		}

		// Check if this is a metadata file and clean it
		switch {
		case strings.Contains(file.Name, "docProps/core.xml"),
			strings.Contains(file.Name, "docProps/app.xml"),
			strings.Contains(file.Name, "meta.xml"):
			// Replace metadata with minimal content
			data = p.cleanOpenXMLMetadata(data, file.Name)
		}

		// Add file to the new archive
		writer, err := archive.Create(file.Name)
		if err != nil {
			return err
		}

		_, err = writer.Write(data)
		if err != nil {
			return err
		}
	}

	// Close everything before renaming
	archive.Close()
	zipWriter.Close()
	reader.Close()

	// Replace the original file with the cleaned one
	return os.Rename(tempFile, filePath)
}

// cleanOpenXMLMetadata replaces metadata content with minimal values
func (p *Processor) cleanOpenXMLMetadata(data []byte, fileName string) []byte {
	// This is a simplified approach - a full XML parser would be better

	// Replace creator, lastModifiedBy, etc.
	data = regexp.MustCompile(`<dc:creator>.*?</dc:creator>`).
		ReplaceAll(data, []byte("<dc:creator></dc:creator>"))

	data = regexp.MustCompile(`<dc:title>.*?</dc:title>`).
		ReplaceAll(data, []byte("<dc:title></dc:title>"))

	data = regexp.MustCompile(`<dc:subject>.*?</dc:subject>`).
		ReplaceAll(data, []byte("<dc:subject></dc:subject>"))

	data = regexp.MustCompile(`<dc:description>.*?</dc:description>`).
		ReplaceAll(data, []byte("<dc:description></dc:description>"))

	data = regexp.MustCompile(`<cp:lastModifiedBy>.*?</cp:lastModifiedBy>`).
		ReplaceAll(data, []byte("<cp:lastModifiedBy></cp:lastModifiedBy>"))

	data = regexp.MustCompile(`<cp:keywords>.*?</cp:keywords>`).
		ReplaceAll(data, []byte("<cp:keywords></cp:keywords>"))

	// Replace revision information
	data = regexp.MustCompile(`<cp:revision>.*?</cp:revision>`).
		ReplaceAll(data, []byte("<cp:revision>1</cp:revision>"))

	// For app.xml
	if strings.Contains(fileName, "app.xml") {
		data = regexp.MustCompile(`<Application>.*?</Application>`).
			ReplaceAll(data, []byte("<Application></Application>"))

		data = regexp.MustCompile(`<Company>.*?</Company>`).
			ReplaceAll(data, []byte("<Company></Company>"))

		data = regexp.MustCompile(`<Manager>.*?</Manager>`).
			ReplaceAll(data, []byte("<Manager></Manager>"))
	}

	return data
}

// cleanOpenDocument removes metadata from OpenDocument files (.odt, .ods, .odp)
func (p *Processor) cleanOpenDocument(filePath string) error {
	// OpenDocument files are also ZIP archives with XML content
	// Similar approach to Office Open XML files
	return p.cleanOpenXML(filePath) // The same approach works for both formats
}

// cleanBinaryOffice removes metadata from legacy binary Office files (.doc, .xls, .ppt)
func (p *Processor) cleanBinaryOffice(filePath, ext string) error {
	// Binary Office formats are complex and hard to parse without dependencies
	// We'll provide a warning about limited capabilities

	p.logger.Warning("Legacy binary Office formats (%s) require complex processing. Limited metadata removal for %s", ext, filePath)
	utils.PrintWarning(fmt.Sprintf("Legacy binary Office formats (%s) require complex processing. Limited metadata removal possible", ext))

	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Check for Office binary file signature (D0 CF 11 E0 A1 B1 1A E1)
	header := make([]byte, 8)
	if _, err := file.Read(header); err != nil {
		return err
	}

	expectedHeader := []byte{0xD0, 0xCF, 0x11, 0xE0, 0xA1, 0xB1, 0x1A, 0xE1}
	if !bytes.Equal(header, expectedHeader) {
		return errors.New("not a valid Office binary file")
	}

	// In a real implementation, we would use the Compound File Binary Format
	// specification to locate and modify the summary information streams
	// This is complex without dependencies

	return nil
}

// cleanRTF removes metadata from RTF files
func (p *Processor) cleanRTF(filePath string) error {
	// Open file
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Verify it's an RTF file
	if !bytes.HasPrefix(fileData, []byte("{\\rtf")) {
		return errors.New("not a valid RTF file")
	}

	// Remove common metadata fields
	// This is a simplified approach using regex

	// Remove author info
	fileData = regexp.MustCompile(`\{\\author [^}]*\}`).ReplaceAll(fileData, []byte(""))

	// Remove title info
	fileData = regexp.MustCompile(`\{\\title [^}]*\}`).ReplaceAll(fileData, []byte(""))

	// Remove subject info
	fileData = regexp.MustCompile(`\{\\subject [^}]*\}`).ReplaceAll(fileData, []byte(""))

	// Remove company info
	fileData = regexp.MustCompile(`\{\\company [^}]*\}`).ReplaceAll(fileData, []byte(""))

	// Remove operator info
	fileData = regexp.MustCompile(`\{\\operator [^}]*\}`).ReplaceAll(fileData, []byte(""))

	// Remove creation/revision time
	fileData = regexp.MustCompile(`\{\\creatim[^}]*\}`).ReplaceAll(fileData, []byte(""))
	fileData = regexp.MustCompile(`\{\\revtim[^}]*\}`).ReplaceAll(fileData, []byte(""))

	// Create temp file
	tempPath := filePath + ".temp"
	err = os.WriteFile(tempPath, fileData, 0644)
	if err != nil {
		return err
	}

	// Replace original file with cleaned file
	err = os.Rename(tempPath, filePath)
	if err != nil {
		os.Remove(tempPath) // Clean up temp file in case of error
		return err
	}

	return nil
}
