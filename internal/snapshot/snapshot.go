// Package snapshot captures and compares crontab expression state over time,
// enabling change detection between two points in analysis.
package snapshot

import (
	"fmt"
	"time"

	"github.com/user/crontab-lint/internal/normalizer"
	"github.com/user/crontab-lint/internal/parser"
	"github.com/user/crontab-lint/internal/statistics"
)

// Snapshot holds a point-in-time capture of a crontab expression's metadata.
type Snapshot struct {
	Expression  string            `json:"expression"`
	Normalized  string            `json:"normalized"`
	CapturedAt  time.Time         `json:"captured_at"`
	Fields      map[string]string `json:"fields"`
	Stats       statistics.Result `json:"stats"`
	Valid       bool              `json:"valid"`
	ParseError  string            `json:"parse_error,omitempty"`
}

// Take captures a snapshot of the given crontab expression at the current time.
func Take(expr string) Snapshot {
	snap := Snapshot{
		Expression: expr,
		CapturedAt: time.Now().UTC(),
		Fields:     make(map[string]string),
	}

	norm, err := normalizer.Normalize(expr)
	if err != nil {
		snap.Valid = false
		snap.ParseError = err.Error()
		return snap
	}
	snap.Normalized = norm

	entry, err := parser.Parse(norm)
	if err != nil {
		snap.Valid = false
		snap.ParseError = err.Error()
		return snap
	}
	snap.Valid = true

	labels := []string{"minute", "hour", "dom", "month", "dow"}
	parts := []string{entry.Minute, entry.Hour, entry.Dom, entry.Month, entry.Dow}
	for i, label := range labels {
		snap.Fields[label] = parts[i]
	}

	result, err := statistics.Compute(norm)
	if err == nil {
		snap.Stats = result
	}

	return snap
}

// Diff describes a change between two snapshots.
type Diff struct {
	Field    string `json:"field"`
	Before   string `json:"before"`
	After    string `json:"after"`
	Changed  bool   `json:"changed"`
}

// Compare returns per-field diffs between two snapshots.
func Compare(before, after Snapshot) ([]Diff, error) {
	if !before.Valid {
		return nil, fmt.Errorf("before snapshot is invalid: %s", before.ParseError)
	}
	if !after.Valid {
		return nil, fmt.Errorf("after snapshot is invalid: %s", after.ParseError)
	}

	fields := []string{"minute", "hour", "dom", "month", "dow"}
	diffs := make([]Diff, 0, len(fields))
	for _, f := range fields {
		b := before.Fields[f]
		a := after.Fields[f]
		diffs = append(diffs, Diff{
			Field:   f,
			Before:  b,
			After:   a,
			Changed: b != a,
		})
	}
	return diffs, nil
}
