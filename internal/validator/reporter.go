package validator

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// FormatReport formats the validation report in the specified format
func FormatReport(report *ValidationReport, format string) (string, error) {
	switch format {
	case "text":
		return formatText(report), nil
	case "json":
		return formatJSON(report)
	case "csv":
		return formatCSV(report)
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

// formatText formats the report as human-readable text
func formatText(report *ValidationReport) string {
	var sb strings.Builder

	// Header
	sb.WriteString("Link Validation Report\n")
	sb.WriteString("======================\n\n")

	// Summary
	sb.WriteString(fmt.Sprintf("Pages Processed:   %d\n", report.PagesProcessed))
	sb.WriteString(fmt.Sprintf("Total Links:       %d\n", report.TotalLinks))
	sb.WriteString(fmt.Sprintf("Unique URLs:       %d\n", report.UniqueURLs))
	sb.WriteString(fmt.Sprintf("✓ Success:         %d (%.1f%%)\n", report.SuccessLinks, percentage(report.SuccessLinks, report.TotalLinks)))
	sb.WriteString(fmt.Sprintf("✗ Broken:          %d (%.1f%%)\n", report.BrokenLinks, percentage(report.BrokenLinks, report.TotalLinks)))
	sb.WriteString(fmt.Sprintf("⚠ Warnings:        %d (%.1f%%)\n", report.WarningLinks, percentage(report.WarningLinks, report.TotalLinks)))
	sb.WriteString(fmt.Sprintf("Internal Links:    %d\n", report.InternalLinks))
	sb.WriteString(fmt.Sprintf("External Links:    %d\n", report.ExternalLinks))
	sb.WriteString(fmt.Sprintf("Cached Results:    %d\n", report.CachedLinks))
	sb.WriteString(fmt.Sprintf("Duration:          %s\n\n", report.Duration.Round(time.Millisecond)))

	// Links by tag type
	if len(report.LinksByTag) > 0 {
		sb.WriteString("Links by Type:\n")
		for tag, count := range report.LinksByTag {
			sb.WriteString(fmt.Sprintf("  <%s>: %d\n", tag, count))
		}
		sb.WriteString("\n")
	}

	// Broken links details
	if report.BrokenLinks > 0 {
		sb.WriteString("Broken Links:\n")
		sb.WriteString("-------------\n")
		for _, result := range report.Results {
			if result.IsBroken {
				sb.WriteString(fmt.Sprintf("\n✗ %s\n", result.TargetURL))
				sb.WriteString(fmt.Sprintf("  Source: %s\n", result.SourceURL))
				sb.WriteString(fmt.Sprintf("  Tag:    <%s>\n", result.Tag))
				if result.LinkText != "" {
					sb.WriteString(fmt.Sprintf("  Text:   %s\n", truncate(result.LinkText, 60)))
				}
				if result.Error != nil {
					sb.WriteString(fmt.Sprintf("  Error:  %v\n", result.Error))
				} else {
					sb.WriteString(fmt.Sprintf("  Status: %d %s\n", result.StatusCode, result.Status))
				}
			}
		}
		sb.WriteString("\n")
	}

	// Warning links (redirects)
	if report.WarningLinks > 0 {
		sb.WriteString("Warnings (Redirects):\n")
		sb.WriteString("--------------------\n")
		for _, result := range report.Results {
			if !result.IsBroken && result.StatusCode >= 300 && result.StatusCode < 400 {
				sb.WriteString(fmt.Sprintf("\n⚠ %s\n", result.TargetURL))
				sb.WriteString(fmt.Sprintf("  Source: %s\n", result.SourceURL))
				sb.WriteString(fmt.Sprintf("  Status: %d %s\n", result.StatusCode, result.Status))
			}
		}
		sb.WriteString("\n")
	}

	// Summary footer
	if report.BrokenLinks == 0 {
		sb.WriteString("✓ All links are valid!\n")
	} else {
		sb.WriteString(fmt.Sprintf("✗ Found %d broken link(s) that need attention.\n", report.BrokenLinks))
	}

	return sb.String()
}

// formatJSON formats the report as JSON
func formatJSON(report *ValidationReport) (string, error) {
	data, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// formatCSV formats the report as CSV
func formatCSV(report *ValidationReport) (string, error) {
	var sb strings.Builder
	writer := csv.NewWriter(&sb)

	// Header
	header := []string{"Source URL", "Target URL", "Status Code", "Status", "Is Broken", "Is External", "Tag", "Link Text", "Error", "Duration (ms)"}
	if err := writer.Write(header); err != nil {
		return "", err
	}

	// Data rows
	for _, result := range report.Results {
		errorStr := ""
		if result.Error != nil {
			errorStr = result.Error.Error()
		}

		row := []string{
			result.SourceURL,
			result.TargetURL,
			fmt.Sprintf("%d", result.StatusCode),
			result.Status,
			fmt.Sprintf("%t", result.IsBroken),
			fmt.Sprintf("%t", result.IsExternal),
			result.Tag,
			result.LinkText,
			errorStr,
			fmt.Sprintf("%d", result.Duration.Milliseconds()),
		}
		if err := writer.Write(row); err != nil {
			return "", err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}

	return sb.String(), nil
}

// WriteReportToFile writes the report to a file
func WriteReportToFile(report *ValidationReport, format, filename string) error {
	content, err := FormatReport(report, format)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, []byte(content), 0644)
}

// Helper functions

func percentage(part, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(part) / float64(total) * 100
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}
