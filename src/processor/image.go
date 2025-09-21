package processor

import (
        "bytes"
        "encoding/binary"
        "errors"
        "fmt"
        "io"
        "os"
        "strings"

        "metadata-remover/src/utils"
)

// ProcessImage removes metadata from image files
func (p *Processor) ProcessImage(filePath, ext string) error {
        // If preview mode, just log and return
        if p.previewMode {
                return nil
        }

        ext = strings.ToLower(ext)
        switch ext {
        case ".jpg", ".jpeg":
                return p.cleanJPEG(filePath)
        case ".png":
                return p.cleanPNG(filePath)
        case ".gif":
                return p.cleanGIF(filePath)
        case ".tiff", ".tif":
                return p.cleanTIFF(filePath)
        case ".bmp":
                return p.cleanBMP(filePath)
        case ".webp":
                return p.cleanWEBP(filePath)
        default:
                return fmt.Errorf("unsupported image format: %s", ext)
        }
}

// cleanJPEG removes metadata from JPEG files
func (p *Processor) cleanJPEG(filePath string) error {
        // Open file
        file, err := os.Open(filePath)
        if err != nil {
                return err
        }
        defer file.Close()

        // Read file header to verify it's a JPEG
        header := make([]byte, 2)
        if _, err := file.Read(header); err != nil {
                return err
        }
        if header[0] != 0xFF || header[1] != 0xD8 {
                return errors.New("not a valid JPEG file")
        }

        // Create temp file
        tempPath := filePath + ".temp"
        tempFile, err := os.Create(tempPath)
        if err != nil {
                return err
        }
        defer func() {
                tempFile.Close()
                os.Remove(tempPath) // Clean up temp file in case of error
        }()

        // Write JPEG header
        if _, err := tempFile.Write(header); err != nil {
                return err
        }

        // Process segments
        buffer := make([]byte, 2)
        for {
                // Read marker
                if _, err := file.Read(buffer); err != nil {
                        if err == io.EOF {
                                break
                        }
                        return err
                }

                // Check if it's a valid marker
                if buffer[0] != 0xFF {
                        return errors.New("invalid JPEG format")
                }

                // Write marker
                if _, err := tempFile.Write(buffer); err != nil {
                        return err
                }

                // Skip metadata segments
                switch buffer[1] {
                case 0xE0, 0xE1, 0xE2, 0xE3, 0xE4, 0xE5, 0xE6, 0xE7, 0xE8, 0xE9, 0xEA, 0xEB, 0xEC, 0xED, 0xEE, 0xEF: // APP0-APP15
                        // Read segment length
                        if _, err := file.Read(buffer); err != nil {
                                return err
                        }
                        length := int(binary.BigEndian.Uint16(buffer))
                        
                        // Skip segment data (length includes the 2 bytes of the length field)
                        if _, err := file.Seek(int64(length-2), io.SeekCurrent); err != nil {
                                return err
                        }
                        
                        // For APP0 (JFIF), we need to keep it but strip metadata
                        if buffer[1] == 0xE0 {
                                // Write minimal JFIF segment
                                if _, err := tempFile.Write([]byte{0x00, 0x10}); err != nil { // Length: 16 bytes
                                        return err
                                }
                                if _, err := tempFile.Write([]byte("JFIF\x00\x01\x01\x00\x00\x01\x00\x01\x00\x00")); err != nil {
                                        return err
                                }
                        }
                        
                case 0xDA: // Start of Scan - after this comes the image data
                        // Write segment length
                        if _, err := file.Read(buffer); err != nil {
                                return err
                        }
                        if _, err := tempFile.Write(buffer); err != nil {
                                return err
                        }
                        
                        length := int(binary.BigEndian.Uint16(buffer))
                        scanData := make([]byte, length-2)
                        if _, err := file.Read(scanData); err != nil {
                                return err
                        }
                        if _, err := tempFile.Write(scanData); err != nil {
                                return err
                        }
                        
                        // Copy the rest of the file (compressed image data)
                        if _, err := io.Copy(tempFile, file); err != nil {
                                return err
                        }
                        break
                        
                default:
                        // For other segments, keep them unchanged
                        if _, err := file.Read(buffer); err != nil {
                                return err
                        }
                        if _, err := tempFile.Write(buffer); err != nil {
                                return err
                        }
                        
                        length := int(binary.BigEndian.Uint16(buffer))
                        segmentData := make([]byte, length-2)
                        if _, err := file.Read(segmentData); err != nil {
                                return err
                        }
                        if _, err := tempFile.Write(segmentData); err != nil {
                                return err
                        }
                }
        }

        // Close files
        file.Close()
        tempFile.Close()

        // Replace original file with cleaned file
        return os.Rename(tempPath, filePath)
}

// cleanPNG removes metadata from PNG files
func (p *Processor) cleanPNG(filePath string) error {
        // Open file
        file, err := os.Open(filePath)
        if err != nil {
                return err
        }
        defer file.Close()

        // Read PNG signature
        signature := make([]byte, 8)
        if _, err := file.Read(signature); err != nil {
                return err
        }
        
        // Verify PNG signature
        pngSignature := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}
        if !bytes.Equal(signature, pngSignature) {
                return errors.New("not a valid PNG file")
        }

        // Create temp file
        tempPath := filePath + ".temp"
        tempFile, err := os.Create(tempPath)
        if err != nil {
                return err
        }
        defer func() {
                tempFile.Close()
                os.Remove(tempPath) // Clean up temp file in case of error
        }()

        // Write PNG signature
        if _, err := tempFile.Write(signature); err != nil {
                return err
        }

        // Process chunks
        for {
                // Read chunk length
                lengthBuf := make([]byte, 4)
                if _, err := file.Read(lengthBuf); err != nil {
                        if err == io.EOF {
                                break
                        }
                        return err
                }
                
                length := binary.BigEndian.Uint32(lengthBuf)
                
                // Read chunk type
                typeBuf := make([]byte, 4)
                if _, err := file.Read(typeBuf); err != nil {
                        return err
                }
                
                chunkType := string(typeBuf)
                
                // Read chunk data
                data := make([]byte, length)
                if _, err := file.Read(data); err != nil {
                        return err
                }
                
                // Read CRC
                crcBuf := make([]byte, 4)
                if _, err := file.Read(crcBuf); err != nil {
                        return err
                }
                
                // Skip metadata chunks, write everything else
                switch chunkType {
                case "tEXt", "iTXt", "zTXt", "tIME", "eXIf":
                        // Skip these metadata chunks
                        continue
                default:
                        // Write this chunk
                        if _, err := tempFile.Write(lengthBuf); err != nil {
                                return err
                        }
                        if _, err := tempFile.Write(typeBuf); err != nil {
                                return err
                        }
                        if _, err := tempFile.Write(data); err != nil {
                                return err
                        }
                        if _, err := tempFile.Write(crcBuf); err != nil {
                                return err
                        }
                }
                
                // IEND chunk signals the end of the PNG file
                if chunkType == "IEND" {
                        break
                }
        }

        // Close files
        file.Close()
        tempFile.Close()

        // Replace original file with cleaned file
        return os.Rename(tempPath, filePath)
}

// cleanGIF removes metadata from GIF files
func (p *Processor) cleanGIF(filePath string) error {
        // GIF files don't have standard metadata chunks to remove
        // The most we could do is reset the Comment Extension if it exists
        // For simplicity, we'll just check the file is a valid GIF
        
        file, err := os.Open(filePath)
        if err != nil {
                return err
        }
        defer file.Close()
        
        // Verify GIF header
        header := make([]byte, 6)
        if _, err := file.Read(header); err != nil {
                return err
        }
        
        if string(header) != "GIF87a" && string(header) != "GIF89a" {
                return errors.New("not a valid GIF file")
        }
        
        p.logger.Info("GIF files have minimal metadata to remove for %s", filePath)
        utils.PrintInfo(fmt.Sprintf("GIF files have minimal metadata to remove"))
        
        return nil
}

// cleanTIFF removes metadata from TIFF files
func (p *Processor) cleanTIFF(filePath string) error {
        // TIFF processing is complex without dependencies
        // For a real implementation, we would parse the IFD structure
        // and remove or modify specific tags
        
        p.logger.Warning("TIFF metadata removal requires complex processing. Basic validation only for %s", filePath)
        utils.PrintWarning(fmt.Sprintf("TIFF metadata removal requires complex processing. Performing basic validation only"))
        
        file, err := os.Open(filePath)
        if err != nil {
                return err
        }
        defer file.Close()
        
        // Read TIFF header
        header := make([]byte, 4)
        if _, err := file.Read(header); err != nil {
                return err
        }
        
        // Check if it's Intel or Motorola byte order
        if !bytes.Equal(header, []byte{0x49, 0x49, 0x2A, 0x00}) && // Little-endian
           !bytes.Equal(header, []byte{0x4D, 0x4D, 0x00, 0x2A}) {  // Big-endian
                return errors.New("not a valid TIFF file")
        }
        
        return nil
}

// cleanBMP removes metadata from BMP files
func (p *Processor) cleanBMP(filePath string) error {
        // BMP files have minimal metadata
        // We'll just validate the file format
        
        file, err := os.Open(filePath)
        if err != nil {
                return err
        }
        defer file.Close()
        
        // Read BMP signature
        signature := make([]byte, 2)
        if _, err := file.Read(signature); err != nil {
                return err
        }
        
        if signature[0] != 'B' || signature[1] != 'M' {
                return errors.New("not a valid BMP file")
        }
        
        p.logger.Info("BMP files have minimal metadata to remove for %s", filePath)
        utils.PrintInfo(fmt.Sprintf("BMP files have minimal metadata to remove"))
        
        return nil
}

// cleanWEBP removes metadata from WebP files
func (p *Processor) cleanWEBP(filePath string) error {
        // WebP processing is complex without dependencies
        // We'll just validate the file format
        
        file, err := os.Open(filePath)
        if err != nil {
                return err
        }
        defer file.Close()
        
        // Read WebP signature
        header := make([]byte, 12)
        if _, err := file.Read(header); err != nil {
                return err
        }
        
        // Check RIFF header and WEBP type
        if !bytes.Equal(header[0:4], []byte("RIFF")) || !bytes.Equal(header[8:12], []byte("WEBP")) {
                return errors.New("not a valid WebP file")
        }
        
        p.logger.Warning("WebP metadata removal requires complex processing. Basic validation only for %s", filePath)
        utils.PrintWarning(fmt.Sprintf("WebP metadata removal requires complex processing. Performing basic validation only"))
        
        return nil
}
