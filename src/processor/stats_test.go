package processor

import (
	"metadata-remover/src/stats"
	"testing"
)

func TestMetadataStats(t *testing.T) {
	t.Run("Basic stat tracking", func(t *testing.T) {
		statsTracker := stats.NewMetadataStats()

		// Add files of different types
		statsTracker.AddFile(stats.TypeImage)
		statsTracker.AddFile(stats.TypeImage)
		statsTracker.AddFile(stats.TypePDF)

		// Verify file counts
		if statsTracker.TotalFiles != 3 {
			t.Errorf("Expected 3 total files, got %d", statsTracker.TotalFiles)
		}

		if statsTracker.ByFileType[stats.TypeImage] != 2 {
			t.Errorf("Expected 2 image files, got %d", statsTracker.ByFileType[stats.TypeImage])
		}

		if statsTracker.ByFileType[stats.TypePDF] != 1 {
			t.Errorf("Expected 1 PDF file, got %d", statsTracker.ByFileType[stats.TypePDF])
		}
	})

	t.Run("Metadata tracking", func(t *testing.T) {
		statsTracker := stats.NewMetadataStats()

		// Add a file
		statsTracker.AddFile(stats.TypeImage)

		// Add metadata for that file
		statsTracker.AddMetadata(stats.TypeImage, "Author", "John Smith")
		statsTracker.AddMetadata(stats.TypeImage, "Creation Date", "2023-01-01")
		statsTracker.AddMetadata(stats.TypeImage, "Author", "Jane Doe")

		// Verify metadata counts
		if statsTracker.TotalMetadataFound != 3 {
			t.Errorf("Expected 3 metadata fields found, got %d", statsTracker.TotalMetadataFound)
		}

		if statsTracker.ByMetadataType["Author"].Count != 2 {
			t.Errorf("Expected 2 Author fields, got %d", statsTracker.ByMetadataType["Author"].Count)
		}

		// Verify examples were captured
		authorExamples := statsTracker.ByMetadataType["Author"].Examples
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
		statsTracker := stats.NewMetadataStats()

		// Add a file
		statsTracker.AddFile(stats.TypeDocument)

		// Add metadata with many examples
		statsTracker.AddMetadata(stats.TypeDocument, "Author", "Person 1")
		statsTracker.AddMetadata(stats.TypeDocument, "Author", "Person 2")
		statsTracker.AddMetadata(stats.TypeDocument, "Author", "Person 3")
		statsTracker.AddMetadata(stats.TypeDocument, "Author", "Person 4")
		statsTracker.AddMetadata(stats.TypeDocument, "Author", "Person 5")

		// Verify example limiting
		authorExamples := statsTracker.ByMetadataType["Author"].Examples
		if len(authorExamples) > 3 {
			t.Errorf("Expected maximum 3 examples, got %d", len(authorExamples))
		}
	})

	t.Run("Stats merging", func(t *testing.T) {
		stats1 := stats.NewMetadataStats()
		stats1.AddFile(stats.TypeImage)
		stats1.AddMetadata(stats.TypeImage, "GPS", "40.7128° N, 74.0060° W")

		stats2 := stats.NewMetadataStats()
		stats2.AddFile(stats.TypePDF)
		stats2.AddMetadata(stats.TypePDF, "Author", "Test User")

		// Merge stats2 into stats1
		stats1.MergeStats(stats2)

		// Verify merged results
		if stats1.TotalFiles != 2 {
			t.Errorf("Expected 2 total files after merge, got %d", stats1.TotalFiles)
		}

		if stats1.TotalMetadataFound != 2 {
			t.Errorf("Expected 2 total metadata fields after merge, got %d", stats1.TotalMetadataFound)
		}

		if stats1.ByFileType[stats.TypePDF] != 1 {
			t.Errorf("Expected 1 PDF file after merge, got %d", stats1.ByFileType[stats.TypePDF])
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
