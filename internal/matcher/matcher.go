// Package matcher provides functionality to check whether a given crontab
// expression matches a specific time, useful for testing and validation.
package matcher

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/user/crontab-lint/internal/normalizer"
	"github.com/user/crontab-lint/internal/parser"
)

// MatchResult holds the result of matching a cron expression against a time.
type MatchResult struct {
	Matches  bool
	Expression string
	Time     time.Time
	Fields   map[string]bool
}

// Match checks whether the given cron expression matches the provided time.
// It returns a MatchResult describing which fields matched.
func Match(expression string, t time.Time) (MatchResult, error) {
	norm, err := normalizer.Normalize(expression)
	if err != nil {
		return MatchResult{}, fmt.Errorf("normalize: %w", err)
	}

	entry, err := parser.Parse(norm)
	if err != nil {
		return MatchResult{}, fmt.Errorf("parse: %w", err)
	}

	fields := map[string]bool{
		"minute":     fieldMatches(entry.Minute, t.Minute(), 0, 59),
		"hour":       fieldMatches(entry.Hour, t.Hour(), 0, 23),
		"dayOfMonth": fieldMatches(entry.DayOfMonth, t.Day(), 1, 31),
		"month":      fieldMatches(entry.Month, int(t.Month()), 1, 12),
		"dayOfWeek":  fieldMatches(entry.DayOfWeek, int(t.Weekday()), 0, 6),
	}

	all := fields["minute"] && fields["hour"] && fields["dayOfMonth"] &&
		fields["month"] && fields["dayOfWeek"]

	return MatchResult{
		Matches:    all,
		Expression: expression,
		Time:       t,
		Fields:     fields,
	}, nil
}

// fieldMatches returns true if the cron field expression matches the given value.
func fieldMatches(field string, value, min, max int) bool {
	if field == "*" {
		return true
	}

	// Step expression: */n or start/n
	if strings.Contains(field, "/") {
		parts := strings.SplitN(field, "/", 2)
		step, err := strconv.Atoi(parts[1])
		if err != nil || step <= 0 {
			return false
		}
		start := min
		if parts[0] != "*" {
			start, err = strconv.Atoi(parts[0])
			if err != nil {
				return false
			}
		}
		for v := start; v <= max; v += step {
			if v == value {
				return true
			}
		}
		return false
	}

	// List: a,b,c
	for _, part := range strings.Split(field, ",") {
		if strings.Contains(part, "-") {
			bounds := strings.SplitN(part, "-", 2)
			lo, e1 := strconv.Atoi(bounds[0])
			hi, e2 := strconv.Atoi(bounds[1])
			if e1 == nil && e2 == nil && value >= lo && value <= hi {
				return true
			}
		} else {
			v, err := strconv.Atoi(part)
			if err == nil && v == value {
				return true
			}
		}
	}
	return false
}
