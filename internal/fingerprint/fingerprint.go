// Package fingerprint produces a structural fingerprint of a crontab expression
// that is stable across semantically equivalent representations.
//
// Unlike digest, which hashes the normalized expression string, fingerprint
// encodes the structural shape of each field (wildcard, step, range, list, literal)
// as a compact descriptor useful for grouping or deduplication by shape.
package fingerprint

import (
	"fmt"
	"strings"

	"github.com/user/crontab-lint/internal/normalizer"
	"github.com/user/crontab-lint/internal/parser"
)

// Result holds the structural fingerprint of a crontab expression.
type Result struct {
	// Shape is a slash-separated descriptor of each field's structural type.
	// e.g. "wildcard/step/range/list/literal"
	Shape string

	// FieldShapes holds the individual shape token for each cron field.
	FieldShapes []string

	// Expression is the normalized expression that was fingerprinted.
	Expression string
}

// Compute returns a structural fingerprint for the given crontab expression.
// The expression may use aliases (e.g. @daily) or name tokens (e.g. MON).
func Compute(expr string) (Result, error) {
	norm, err := normalizer.Normalize(expr)
	if err != nil {
		return Result{}, fmt.Errorf("fingerprint: normalize: %w", err)
	}

	entry, err := parser.Parse(norm)
	if err != nil {
		return Result{}, fmt.Errorf("fingerprint: parse: %w", err)
	}

	fields := []string{
		entry.Minute,
		entry.Hour,
		entry.DayOfMonth,
		entry.Month,
		entry.DayOfWeek,
	}

	shapes := make([]string, len(fields))
	for i, f := range fields {
		shapes[i] = shapeOf(f)
	}

	return Result{
		Shape:       strings.Join(shapes, "/"),
		FieldShapes: shapes,
		Expression:  norm,
	}, nil
}

// shapeOf returns a single token describing the structural type of a cron field.
func shapeOf(field string) string {
	switch {
	case field == "*":
		return "wildcard"
	case strings.Contains(field, ","):
		return "list"
	case strings.Contains(field, "/"):
		return "step"
	case strings.Contains(field, "-"):
		return "range"
	default:
		return "literal"
	}
}
