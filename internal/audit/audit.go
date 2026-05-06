// Package audit provides a structured audit trail for crontab expressions,
// recording validation events, warnings, and metadata for each analyzed expression.
package audit

import (
	"fmt"
	"time"

	"github.com/dkarter/crontab-lint/internal/linter"
	"github.com/dkarter/crontab-lint/internal/parser"
)

// Severity represents the level of an audit event.
type Severity string

const (
	SeverityInfo    Severity = "info"
	SeverityWarning Severity = "warning"
	SeverityError   Severity = "error"
)

// Event represents a single audit entry for a crontab expression.
type Event struct {
	Timestamp  time.Time `json:"timestamp"`
	Expression string    `json:"expression"`
	Severity   Severity  `json:"severity"`
	Message    string    `json:"message"`
	Code       string    `json:"code,omitempty"`
}

// Record holds the full audit result for a single crontab expression.
type Record struct {
	Expression string    `json:"expression"`
	AnalyzedAt time.Time `json:"analyzed_at"`
	Valid      bool      `json:"valid"`
	Events     []Event   `json:"events"`
}

// Build produces an audit Record for the given crontab expression by
// running the parser and linter and collecting all findings.
func Build(expr string) Record {
	now := time.Now().UTC()
	rec := Record{
		Expression: expr,
		AnalyzedAt: now,
		Valid:      true,
	}

	_, err := parser.Parse(expr)
	if err != nil {
		rec.Valid = false
		rec.Events = append(rec.Events, Event{
			Timestamp:  now,
			Expression: expr,
			Severity:   SeverityError,
			Message:    fmt.Sprintf("parse error: %s", err.Error()),
			Code:       "PARSE_ERROR",
		})
		return rec
	}

	rec.Events = append(rec.Events, Event{
		Timestamp:  now,
		Expression: expr,
		Severity:   SeverityInfo,
		Message:    "expression parsed successfully",
		Code:       "PARSE_OK",
	})

	results := linter.Lint(expr)
	for _, r := range results {
		sev := SeverityWarning
		if r.Level == "error" {
			sev = SeverityError
			rec.Valid = false
		}
		rec.Events = append(rec.Events, Event{
			Timestamp:  now,
			Expression: expr,
			Severity:   sev,
			Message:    r.Message,
			Code:       r.Code,
		})
	}

	return rec
}
