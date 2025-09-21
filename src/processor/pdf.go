package processor

import (
	"bytes"
	"errors"
	"os"
	"regexp"
)

// ProcessPDF removes metadata from PDF files
func (p *Processor) ProcessPDF(filePath string) error {
	// If preview mode, just log and return
	if p.previewMode {
		return nil
	}

	// Open file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Verify it's a PDF file
	header := make([]byte, 5)
	if _, err := file.Read(header); err != nil {
		return err
	}
	if !bytes.Equal(header, []byte("%PDF-")) {
		return errors.New("not a valid PDF file")
	}

	// Reopen the file to process it
	file.Close()
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Find and remove common metadata dictionaries
	// This is a simplified approach - a full PDF parser would be better
	cleanedContent := p.removeInfoDictionary(fileContent)
	cleanedContent = p.removeXMPMetadata(cleanedContent)
	cleanedContent = p.removeDocumentInfo(cleanedContent)

	// Create temp file
	tempPath := filePath + ".temp"
	err = os.WriteFile(tempPath, cleanedContent, 0644)
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

// removeInfoDictionary removes the Info dictionary from PDF content
func (p *Processor) removeInfoDictionary(content []byte) []byte {
	// Pattern to match the Info dictionary
	infoPattern := regexp.MustCompile(`/Info\s+\d+\s+\d+\s+R`)

	// Find all instances
	matches := infoPattern.FindAllIndex(content, -1)

	if len(matches) == 0 {
		return content
	}

	// Generate replacement (keep structure but remove content)
	replacement := []byte("/Info 0 0 R")

	// Create new content with replaced Info references
	newContent := make([]byte, 0, len(content))
	lastPos := 0

	for _, match := range matches {
		newContent = append(newContent, content[lastPos:match[0]]...)
		newContent = append(newContent, replacement...)
		lastPos = match[1]
	}

	newContent = append(newContent, content[lastPos:]...)

	return newContent
}

// removeXMPMetadata removes XMP metadata from PDF content
func (p *Processor) removeXMPMetadata(content []byte) []byte {
	// Pattern to match XMP metadata streams
	// This is a simplified approach - a full XML parser would be better
	startPattern := []byte("<x:xmpmeta")
	endPattern := []byte("</x:xmpmeta>")

	// Find start position
	startPos := bytes.Index(content, startPattern)
	if startPos == -1 {
		return content
	}

	// Find end position
	endPos := bytes.Index(content[startPos:], endPattern)
	if endPos == -1 {
		return content
	}
	endPos += startPos + len(endPattern)

	// Create new content without XMP metadata
	newContent := make([]byte, 0, len(content)-(endPos-startPos))
	newContent = append(newContent, content[:startPos]...)
	newContent = append(newContent, content[endPos:]...)

	return newContent
}

// removeDocumentInfo removes document information from PDF content
func (p *Processor) removeDocumentInfo(content []byte) []byte {
	// Pattern to match document information dictionaries
	pattern := regexp.MustCompile(`<<\s*(/Title|\s*/Author|\s*/Subject|\s*/Keywords|\s*/Creator|\s*/Producer|\s*/CreationDate|\s*/ModDate|\s*/Trapped)[^>]*>>`)

	// Find all matches
	matches := pattern.FindAllSubmatchIndex(content, -1)

	if len(matches) == 0 {
		return content
	}

	// Create new content with empty document information
	newContent := make([]byte, 0, len(content))
	lastPos := 0

	for _, match := range matches {
		newContent = append(newContent, content[lastPos:match[0]]...)
		newContent = append(newContent, []byte("<< >>")...)
		lastPos = match[1]
	}

	newContent = append(newContent, content[lastPos:]...)

	return newContent
}
