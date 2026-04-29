// Package diff compares two crontab expressions and reports their differences.
package diff

import (
	"fmt"
	"strings"

	"github.com/user/crontab-lint/internal/parser"
)

// FieldDiff represents a difference in a single crontab field.
type FieldDiff struct {
	Field string
	From  string
	To    string
}

// Result holds the outcome of comparing two crontab expressions.
type Result struct {
	From    string
	To      string
	Diffs   []FieldDiff
	Changed bool
}

var fieldNames = []string{"minute", "hour", "day-of-month", "month", "day-of-week"}

// Compare takes two crontab expressions and returns a Result describing
// which fields changed between them. Returns an error if either expression
// is invalid.
func Compare(from, to string) (*Result, error) {
	entriFrom, err := parser.Parse(from)
	if err != nil {
		return nil, fmt.Errorf("invalid 'from' expression: %w", err)
	}
	entriTo, err := parser.Parse(to)
	if err != nil {
		return nil, fmt.Errorf("invalid 'to' expression: %w", err)
	}

	fromFields := scheduleFields(entriFrom.Schedule)
	toFields := scheduleFields(entriTo.Schedule)

	result := &Result{From: from, To: to}
	for i, name := range fieldNames {
		if fromFields[i] != toFields[i] {
			result.Diffs = append(result.Diffs, FieldDiff{
				Field: name,
				From:  fromFields[i],
				To:    toFields[i],
			})
			result.Changed = true
		}
	}
	return result, nil
}

// scheduleFields splits a 5-field schedule string into its component parts.
func scheduleFields(schedule string) []string {
	parts := strings.Fields(schedule)
	if len(parts) >= 5 {
		return parts[:5]
	}
	return parts
}
