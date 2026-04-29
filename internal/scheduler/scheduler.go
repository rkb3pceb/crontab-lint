// Package scheduler provides utilities for computing the next scheduled
// run times for a given crontab expression.
package scheduler

import (
	"fmt"
	"time"

	"github.com/user/crontab-lint/internal/parser"
)

// NextN returns the next n scheduled times after the given start time
// for the provided crontab expression. Returns an error if the expression
// is invalid or if n is less than 1.
func NextN(expr string, start time.Time, n int) ([]time.Time, error) {
	if n < 1 {
		return nil, fmt.Errorf("n must be at least 1, got %d", n)
	}

	entry, err := parser.Parse(expr)
	if err != nil {
		return nil, fmt.Errorf("invalid expression: %w", err)
	}

	results := make([]time.Time, 0, n)
	current := start.Add(time.Minute).Truncate(time.Minute)

	for len(results) < n {
		if matches(entry, current) {
			results = append(results, current)
		}
		current = current.Add(time.Minute)
		// Safety limit: scan at most 4 years ahead.
		if current.After(start.Add(4 * 365 * 24 * time.Hour)) {
			break
		}
	}

	return results, nil
}

// Next returns the next scheduled time after start for the given expression.
func Next(expr string, start time.Time) (time.Time, error) {
	times, err := NextN(expr, start, 1)
	if err != nil {
		return time.Time{}, err
	}
	if len(times) == 0 {
		return time.Time{}, fmt.Errorf("no matching time found within 4 years")
	}
	return times[0], nil
}

// matches reports whether t satisfies all cron fields in entry.
func matches(entry parser.Entry, t time.Time) bool {
	if !fieldMatches(entry.Minute, t.Minute(), 0, 59) {
		return false
	}
	if !fieldMatches(entry.Hour, t.Hour(), 0, 23) {
		return false
	}
	if !fieldMatches(entry.DayOfMonth, t.Day(), 1, 31) {
		return false
	}
	if !fieldMatches(entry.Month, int(t.Month()), 1, 12) {
		return false
	}
	// time.Weekday: Sunday=0 ... Saturday=6
	if !fieldMatches(entry.DayOfWeek, int(t.Weekday()), 0, 6) {
		return false
	}
	return true
}

// fieldMatches checks whether value matches the cron field token.
func fieldMatches(field string, value, min, max int) bool {
	if field == "*" {
		return true
	}
	values, err := parser.ExpandField(field, min, max)
	if err != nil {
		return false
	}
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}
