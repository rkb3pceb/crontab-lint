// Package linter provides static analysis and diagnostics for crontab expressions.
package linter

import (
	"fmt"
	"strings"

	"github.com/yourorg/crontab-lint/internal/humanizer"
	"github.com/yourorg/crontab-lint/internal/parser"
)

// Severity represents the level of a diagnostic finding.
type Severity string

const (
	SeverityError   Severity = "error"
	SeverityWarning Severity = "warning"
	SeverityInfo    Severity = "info"
)

// Diagnostic holds a single lint finding for a crontab expression.
type Diagnostic struct {
	Severity Severity
	Field    string
	Message  string
}

// Result holds the full lint result for a crontab line.
type Result struct {
	Raw         string
	Valid        bool
	Schedule    string
	Diagnostics []Diagnostic
}

// Lint parses and analyzes a single crontab expression line.
// It returns a Result with diagnostics and a human-readable schedule description.
func Lint(line string) Result {
	result := Result{Raw: line}

	entry, err := parser.Parse(line)
	if err != nil {
		result.Valid = false
		result.Diagnostics = append(result.Diagnostics, Diagnostic{
			Severity: SeverityError,
			Field:    "expression",
			Message:  err.Error(),
		})
		return result
	}

	result.Valid = true

	// Generate human-readable description.
	desc, descErr := humanizer.Describe(fmt.Sprintf("%s %s %s %s %s",
		entry.Minute, entry.Hour, entry.DayOfMonth, entry.Month, entry.DayOfWeek))
	if descErr == nil {
		result.Schedule = desc
	}

	// Warn on broad wildcard schedules that run every minute.
	if entry.Minute == "*" && entry.Hour == "*" {
		result.Diagnostics = append(result.Diagnostics, Diagnostic{
			Severity: SeverityWarning,
			Field:    "minute/hour",
			Message:  "schedule runs every minute — consider restricting frequency",
		})
	}

	// Info: empty command check.
	if strings.TrimSpace(entry.Command) == "" {
		result.Diagnostics = append(result.Diagnostics, Diagnostic{
			Severity: SeverityError,
			Field:    "command",
			Message:  "command is empty",
		})
		result.Valid = false
	}

	// Info: day-of-week and day-of-month both set (non-wildcard) — OR semantics may surprise users.
	if entry.DayOfMonth != "*" && entry.DayOfWeek != "*" {
		result.Diagnostics = append(result.Diagnostics, Diagnostic{
			Severity: SeverityInfo,
			Field:    "dom/dow",
			Message:  "both day-of-month and day-of-week are set; cron uses OR semantics between them",
		})
	}

	return result
}
