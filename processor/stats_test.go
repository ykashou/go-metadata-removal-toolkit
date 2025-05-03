package processor

import (
	"testing"
)

func TestMetadataStats(t *testing.T) {
	t.Run("Basic stat tracking", func(t *testing.T) {
		stats := NewMetadataStats()
		
		// Add files of different types
		stats.AddFile(TypeImage)
		stats.AddFile(TypeImage)
		stats.AddFile(TypePDF)
		
		// Verify file counts
		if stats.TotalFiles != 3 {
			t.Errorf("Expected 3 total files, got %d", stats.TotalFiles)
		}
		
		if stats.ByFileType[TypeImage] != 2 {
			t.Errorf("Expected 2 image files, got %d", stats.ByFileType[TypeImage])
		}
		
		if stats.ByFileType[TypePDF] != 1 {
			t.Errorf("Expected 1 PDF file, got %d", stats.ByFileType[TypePDF])
		}
	})
	
	t.Run("Metadata tracking", func(t *testing.T) {
		stats := NewMetadataStats()
		
		// Add a file
		stats.AddFile(TypeImage)
		
		// Add metadata for that file
		stats.AddMetadata(TypeImage, "Author", "John Smith")
		stats.AddMetadata(TypeImage, "Creation Date", "2023-01-01")
		stats.AddMetadata(TypeImage, "Author", "Jane Doe")
		
		// Verify metadata counts
		if stats.TotalMetadataFound != 3 {
			t.Errorf("Expected 3 metadata fields found, got %d", stats.TotalMetadataFound)
		}
		
		if stats.ByMetadataType["Author"].Count != 2 {
			t.Errorf("Expected 2 Author fields, got %d", stats.ByMetadataType["Author"].Count)
		}
		
		// Verify examples were captured
		authorExamples := stats.ByMetadataType["Author"].Examples
		if len(authorExamples) != 2 {
			t.Errorf("Expected 2 examples for Author, got %d", len(authorExamples))
		}
		
		// Verify examples
		foundJohn := false
		foundJane := false
		for _, ex := range authorExamples {
			if ex == "John Smith" {
				foundJohn = true
			}
			if ex == "Jane Doe" {
				foundJane = true
			}
		}
		
		if !foundJohn || !foundJane {
			t.Errorf("Missing expected examples for Author field")
		}
	})
	
	t.Run("Example limiting", func(t *testing.T) {
		stats := NewMetadataStats()
		
		// Add a file
		stats.AddFile(TypeDocument)
		
		// Add metadata with many examples
		stats.AddMetadata(TypeDocument, "Author", "Person 1")
		stats.AddMetadata(TypeDocument, "Author", "Person 2")
		stats.AddMetadata(TypeDocument, "Author", "Person 3")
		stats.AddMetadata(TypeDocument, "Author", "Person 4")
		stats.AddMetadata(TypeDocument, "Author", "Person 5")
		
		// Verify example limiting
		authorExamples := stats.ByMetadataType["Author"].Examples
		if len(authorExamples) > 3 {
			t.Errorf("Expected maximum 3 examples, got %d", len(authorExamples))
		}
	})
	
	t.Run("Stats merging", func(t *testing.T) {
		stats1 := NewMetadataStats()
		stats1.AddFile(TypeImage)
		stats1.AddMetadata(TypeImage, "GPS", "40.7128° N, 74.0060° W")
		
		stats2 := NewMetadataStats()
		stats2.AddFile(TypePDF)
		stats2.AddMetadata(TypePDF, "Author", "Test User")
		
		// Merge stats2 into stats1
		stats1.MergeStats(stats2)
		
		// Verify merged results
		if stats1.TotalFiles != 2 {
			t.Errorf("Expected 2 total files after merge, got %d", stats1.TotalFiles)
		}
		
		if stats1.TotalMetadataFound != 2 {
			t.Errorf("Expected 2 total metadata fields after merge, got %d", stats1.TotalMetadataFound)
		}
		
		if stats1.ByFileType[TypePDF] != 1 {
			t.Errorf("Expected 1 PDF file after merge, got %d", stats1.ByFileType[TypePDF])
		}
		
		if stats1.ByMetadataType["Author"].Count != 1 {
			t.Errorf("Expected 1 Author field after merge, got %d", stats1.ByMetadataType["Author"].Count)
		}
		
		if stats1.ByMetadataType["GPS"].Count != 1 {
			t.Errorf("Expected 1 GPS field after merge, got %d", stats1.ByMetadataType["GPS"].Count)
		}
		
		// Verify examples were merged
		if len(stats1.ByMetadataType["Author"].Examples) != 1 || 
		   stats1.ByMetadataType["Author"].Examples[0] != "Test User" {
			t.Errorf("Author example not correctly merged")
		}
		
		if len(stats1.ByMetadataType["GPS"].Examples) != 1 || 
		   stats1.ByMetadataType["GPS"].Examples[0] != "40.7128° N, 74.0060° W" {
			t.Errorf("GPS example not correctly preserved in merge")
		}
	})
}