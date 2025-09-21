package utils

import (
        "metadata-remover/src/stats"
        "strings"
        "testing"
)

func TestFormatStats(t *testing.T) {
        // Create test stats
        testStats := stats.NewMetadataStats()
        
        // Add some test data
        testStats.AddFile(stats.TypeImage)
        testStats.AddFile(stats.TypeImage)
        testStats.AddFile(stats.TypePDF)
        
        testStats.AddMetadata(stats.TypeImage, "Author", "John Doe")
        testStats.AddMetadata(stats.TypeImage, "Creation Date", "2023-01-01")
        testStats.AddMetadata(stats.TypeImage, "GPS", "40.7128° N, 74.0060° W")
        testStats.AddMetadata(stats.TypePDF, "Author", "Jane Smith")
        
        t.Run("Text format", func(t *testing.T) {
                result := FormatStats(testStats, "terminal")
                
                // Check if the output contains expected sections
                if !strings.Contains(result, "METADATA REMOVAL STATISTICS") {
                        t.Error("Missing title in output")
                }
                
                if !strings.Contains(result, "Total files processed: 3") {
                        t.Error("Missing or incorrect file count")
                }
                
                if !strings.Contains(result, "Total metadata fields found: 4") {
                        t.Error("Missing or incorrect metadata count")
                }
                
                if !strings.Contains(result, "Images: 2 files") {
                        t.Error("Missing or incorrect image file count")
                }
                
                if !strings.Contains(result, "PDFs: 1 files") {
                        t.Error("Missing or incorrect PDF file count")
                }
                
		// Check for Author count (accounting for color codes)
		if !strings.Contains(result, "Author") || !strings.Contains(result, "2 occurrences") {
			t.Error("Missing or incorrect Author metadata count")
		}
                
                // Check for examples
                if !strings.Contains(result, "\"John Doe\"") || !strings.Contains(result, "\"Jane Smith\"") {
                        t.Error("Missing examples in output")
                }
        })
        
        t.Run("JSON format", func(t *testing.T) {
                result := FormatStats(testStats, "json")
                
                // Check if the output is valid JSON
                if !strings.HasPrefix(result, "{") || !strings.HasSuffix(result, "}") {
                        t.Error("Output is not valid JSON")
                }
                
                // Check if the output contains expected fields
                if !strings.Contains(result, "\"TotalFiles\": 3") {
                        t.Error("Missing or incorrect file count in JSON")
                }
                
                if !strings.Contains(result, "\"TotalMetadataFound\": 4") {
                        t.Error("Missing or incorrect metadata count in JSON")
                }
                
                if !strings.Contains(result, "\"ByFileType\"") {
                        t.Error("Missing ByFileType section in JSON")
                }
                
                if !strings.Contains(result, "\"ByMetadataType\"") {
                        t.Error("Missing ByMetadataType section in JSON")
                }
        })
}