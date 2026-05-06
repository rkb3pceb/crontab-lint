// Package tags provides functionality for extracting and classifying
// semantic tags from crontab expressions based on their schedule patterns.
package tags

import (
	"strings"

	"github.com/example/crontab-lint/internal/parser"
)

// Tag represents a semantic label applied to a crontab expression.
type Tag struct {
	Name        string
	Description string
}

// Extract analyzes a crontab expression and returns a list of semantic tags
// that describe its scheduling characteristics.
func Extract(expr string) ([]Tag, error) {
	entry, err := parser.Parse(expr)
	if err != nil {
		return nil, err
	}

	fields := strings.Fields(entry.Schedule)
	if len(fields) < 5 {
		return nil, fmt.Errorf("too few fields in schedule")
	}

	minute, hour, dom, month, dow := fields[0], fields[1], fields[2], fields[3], fields[4]

	var tags []Tag

	if minute == "*" && hour == "*" {
		tags = append(tags, Tag{Name: "high-frequency", Description: "Runs every minute"})
	} else if minute == "0" && hour == "*" {
		tags = append(tags, Tag{Name: "hourly", Description: "Runs once per hour"})
	} else if minute == "0" && hour == "0" {
		tags = append(tags, Tag{Name: "daily", Description: "Runs once per day"})
	} else if strings.HasPrefix(minute, "*/") || strings.HasPrefix(hour, "*/") {
		tags = append(tags, Tag{Name: "interval", Description: "Runs on a fixed interval"})
	}

	if month != "*" {
		tags = append(tags, Tag{Name: "month-restricted", Description: "Restricted to specific months"})
	}

	if dow != "*" && dom == "*" {
		tags = append(tags, Tag{Name: "weekday-specific", Description: "Restricted to specific weekdays"})
	}

	if dom != "*" && dow == "*" {
		tags = append(tags, Tag{Name: "dom-specific", Description: "Restricted to specific days of month"})
	}

	if strings.Contains(minute, ",") || strings.Contains(hour, ",") {
		tags = append(tags, Tag{Name: "multi-time", Description: "Runs at multiple times"})
	}

	if len(tags) == 0 {
		tags = append(tags, Tag{Name: "custom", Description: "Custom schedule"})
	}

	return tags, nil
}

// Names returns only the tag names from a slice of tags.
func Names(tags []Tag) []string {
	names := make([]string, len(tags))
	for i, t := range tags {
		names[i] = t.Name
	}
	return names
}
