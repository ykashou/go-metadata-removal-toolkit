package processor

// MetadataField represents a type of metadata found in files
type MetadataField struct {
	Name     string
	Count    int
	Examples []string // Store a few examples of metadata values (limited to prevent excessive memory usage)
}

// MetadataStats tracks statistics about metadata found in files
type MetadataStats struct {
	TotalFiles         int
	TotalMetadataFound int
	ByFileType         map[string]int            // Count of files by type
	ByMetadataType     map[string]*MetadataField // Statistics by metadata field type
	FileTypeMetadata   map[string]map[string]int // Count of metadata fields by file type
}

// NewMetadataStats creates a new stats tracker
func NewMetadataStats() *MetadataStats {
	return &MetadataStats{
		ByFileType:       make(map[string]int),
		ByMetadataType:   make(map[string]*MetadataField),
		FileTypeMetadata: make(map[string]map[string]int),
	}
}

// AddFile increments the counter for a specific file type
func (ms *MetadataStats) AddFile(fileType string) {
	ms.TotalFiles++
	ms.ByFileType[fileType]++

	// Initialize the inner map if it doesn't exist
	if _, ok := ms.FileTypeMetadata[fileType]; !ok {
		ms.FileTypeMetadata[fileType] = make(map[string]int)
	}
}

// AddMetadata tracks a metadata field found in a file
func (ms *MetadataStats) AddMetadata(fileType, fieldName, example string) {
	// Add to total count
	ms.TotalMetadataFound++

	// Track by metadata type
	if _, ok := ms.ByMetadataType[fieldName]; !ok {
		ms.ByMetadataType[fieldName] = &MetadataField{
			Name:     fieldName,
			Examples: make([]string, 0, 3), // Cap examples at 3
		}
	}

	ms.ByMetadataType[fieldName].Count++

	// Add example if we have space and it's not already in the list
	field := ms.ByMetadataType[fieldName]
	if len(example) > 0 && len(field.Examples) < 3 {
		// Check if example already exists
		exists := false
		for _, ex := range field.Examples {
			if ex == example {
				exists = true
				break
			}
		}

		if !exists {
			field.Examples = append(field.Examples, example)
		}
	}

	// Track by file type
	if _, ok := ms.FileTypeMetadata[fileType]; ok {
		ms.FileTypeMetadata[fileType][fieldName]++
	}
}

// MergeStats combines two MetadataStats objects
func (ms *MetadataStats) MergeStats(other *MetadataStats) {
	ms.TotalFiles += other.TotalFiles
	ms.TotalMetadataFound += other.TotalMetadataFound

	// Merge by file type
	for fileType, count := range other.ByFileType {
		ms.ByFileType[fileType] += count
	}

	// Merge by metadata type
	for fieldName, field := range other.ByMetadataType {
		if existing, ok := ms.ByMetadataType[fieldName]; ok {
			existing.Count += field.Count

			// Merge examples (keeping up to 3)
			for _, example := range field.Examples {
				if len(existing.Examples) < 3 {
					// Check if example already exists
					exists := false
					for _, ex := range existing.Examples {
						if ex == example {
							exists = true
							break
						}
					}

					if !exists {
						existing.Examples = append(existing.Examples, example)
					}
				}
			}
		} else {
			// Copy field data
			newField := &MetadataField{
				Name:     field.Name,
				Count:    field.Count,
				Examples: make([]string, 0, 3),
			}

			// Copy up to 3 examples
			for i, example := range field.Examples {
				if i < 3 {
					newField.Examples = append(newField.Examples, example)
				} else {
					break
				}
			}

			ms.ByMetadataType[fieldName] = newField
		}
	}

	// Merge by file type metadata
	for fileType, fields := range other.FileTypeMetadata {
		if _, ok := ms.FileTypeMetadata[fileType]; !ok {
			ms.FileTypeMetadata[fileType] = make(map[string]int)
		}

		for fieldName, count := range fields {
			ms.FileTypeMetadata[fileType][fieldName] += count
		}
	}
}
