package utils

import (
        "encoding/json"
        "fmt"
        "metadata-remover/src/stats"
        "sort"
        "strings"
)

// FormatStats generates a human-readable report of metadata statistics
func FormatStats(stats *stats.MetadataStats, format string) string {
        if format == "json" {
                return FormatStatsAsJSON(stats)
        }
        return FormatStatsAsText(stats)
}

// FormatStatsAsText generates a human-readable text report of metadata statistics
func FormatStatsAsText(stats *stats.MetadataStats) string {
        var sb strings.Builder

        // Title
        sb.WriteString("\n")
        sb.WriteString(Blue("=== METADATA REMOVAL STATISTICS ==="))
        sb.WriteString("\n\n")

        // Summary section
        sb.WriteString(Blue("SUMMARY:"))
        sb.WriteString("\n")
        sb.WriteString(fmt.Sprintf("Total files processed: %d\n", stats.TotalFiles))
        sb.WriteString(fmt.Sprintf("Total metadata fields found: %d\n", stats.TotalMetadataFound))
        sb.WriteString("\n")

        // Files by type section
        sb.WriteString(Blue("FILES BY TYPE:"))
        sb.WriteString("\n")
        if len(stats.ByFileType) == 0 {
                sb.WriteString("No files processed.\n")
        } else {
                // Convert map to slice for sorting
                type fileTypeCount struct {
                        Type  string
                        Count int
                }
                fileTypes := make([]fileTypeCount, 0, len(stats.ByFileType))
                for fileType, count := range stats.ByFileType {
                        fileTypes = append(fileTypes, fileTypeCount{fileType, count})
                }

                // Sort by count (descending)
                sort.Slice(fileTypes, func(i, j int) bool {
                        return fileTypes[i].Count > fileTypes[j].Count
                })

                // Print each file type
                for _, ft := range fileTypes {
                        sb.WriteString(fmt.Sprintf("  %s: %d files\n", formatFileType(ft.Type), ft.Count))
                }
        }
        sb.WriteString("\n")

        // Metadata by type section
        sb.WriteString(Blue("METADATA FIELDS FOUND:"))
        sb.WriteString("\n")
        if len(stats.ByMetadataType) == 0 {
                sb.WriteString("No metadata found in processed files.\n")
        } else {
                // Convert map to slice for sorting
                type metadataTypeCount struct {
                        Type     string
                        Count    int
                        Examples []string
                }
                metadataTypes := make([]metadataTypeCount, 0, len(stats.ByMetadataType))
                for metaType, field := range stats.ByMetadataType {
                        metadataTypes = append(metadataTypes, metadataTypeCount{
                                Type:     metaType,
                                Count:    field.Count,
                                Examples: field.Examples,
                        })
                }

                // Sort by count (descending)
                sort.Slice(metadataTypes, func(i, j int) bool {
                        return metadataTypes[i].Count > metadataTypes[j].Count
                })

                // Print each metadata type
                for _, mt := range metadataTypes {
                        sb.WriteString(fmt.Sprintf("  %s: %d occurrences\n", Yellow(mt.Type), mt.Count))
                        
                        // Print examples if available
                        if len(mt.Examples) > 0 {
                                sb.WriteString("    Examples: ")
                                for i, example := range mt.Examples {
                                        if i > 0 {
                                                sb.WriteString(", ")
                                        }
                                        // Truncate very long examples
                                        if len(example) > 50 {
                                                example = example[:47] + "..."
                                        }
                                        sb.WriteString(fmt.Sprintf("\"%s\"", example))
                                }
                                sb.WriteString("\n")
                        }
                }
        }
        sb.WriteString("\n")

        // Metadata by file type section
        sb.WriteString(Blue("METADATA BY FILE TYPE:"))
        sb.WriteString("\n")
        if len(stats.FileTypeMetadata) == 0 {
                sb.WriteString("No metadata found in processed files.\n")
        } else {
                // Convert map to slice for sorting
                type fileTypeData struct {
                        Type     string
                        Metadata map[string]int
                }
                fileTypes := make([]fileTypeData, 0, len(stats.FileTypeMetadata))
                for fileType, metadata := range stats.FileTypeMetadata {
                        fileTypes = append(fileTypes, fileTypeData{
                                Type:     fileType,
                                Metadata: metadata,
                        })
                }

                // Sort by file type
                sort.Slice(fileTypes, func(i, j int) bool {
                        return fileTypes[i].Type < fileTypes[j].Type
                })

                // Print each file type and its metadata
                for _, ft := range fileTypes {
                        sb.WriteString(fmt.Sprintf("  %s:\n", formatFileType(ft.Type)))
                        
                        // Convert metadata map to slice for sorting
                        type metaCount struct {
                                Name  string
                                Count int
                        }
                        metadata := make([]metaCount, 0, len(ft.Metadata))
                        for name, count := range ft.Metadata {
                                metadata = append(metadata, metaCount{name, count})
                        }
                        
                        // Sort by count (descending)
                        sort.Slice(metadata, func(i, j int) bool {
                                return metadata[i].Count > metadata[j].Count
                        })
                        
                        // Print each metadata field
                        for _, meta := range metadata {
                                sb.WriteString(fmt.Sprintf("    - %s: %d occurrences\n", meta.Name, meta.Count))
                        }
                }
        }

        return sb.String()
}

// FormatStatsAsJSON generates a JSON representation of metadata statistics
func FormatStatsAsJSON(stats *stats.MetadataStats) string {
        jsonBytes, err := json.MarshalIndent(stats, "", "  ")
        if err != nil {
                return fmt.Sprintf("Error generating JSON: %v", err)
        }
        return string(jsonBytes)
}

// formatFileType returns a formatted string for the file type
func formatFileType(fileType string) string {
        switch fileType {
        case stats.TypeImage:
                return "Images"
        case stats.TypePDF:
                return "PDFs"
        case stats.TypeDocument:
                return "Documents"
        default:
                return fmt.Sprintf("%s files", strings.ToUpper(fileType[:1]) + fileType[1:])
        }
}