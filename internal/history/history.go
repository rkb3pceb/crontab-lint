// Package history provides functionality to track and display
// the last N execution times for a given cron expression.
package history

import (
	"fmt"
	"time"

	"github.com/user/crontab-lint/internal/scheduler"
)

// Entry represents a single past execution time.
type Entry struct {
	Time      time.Time
	Formatted string
}

// Result holds the computed history for a cron expression.
type Result struct {
	Expression string
	Entries    []Entry
}

// Last returns the last n execution times before the given reference time.
// It works by computing future times from a shifted origin and reversing.
func Last(expression string, ref time.Time, n int) (*Result, error) {
	if n <= 0 {
		return nil, fmt.Errorf("history: n must be greater than 0, got %d", n)
	}
	if n > 100 {
		return nil, fmt.Errorf("history: n must be <= 100, got %d", n)
	}

	// Step back far enough to find n matches.
	// Use a window of n*2 minutes plus a buffer.
	windowMinutes := n * 60 * 24 * 7 // up to a week back
	start := ref.Add(-time.Duration(windowMinutes) * time.Minute)

	times, err := scheduler.NextN(expression, start, windowMinutes)
	if err != nil {
		return nil, fmt.Errorf("history: %w", err)
	}

	// Filter to only times before ref.
	var before []time.Time
	for _, t := range times {
		if t.Before(ref) {
			before = append(before, t)
		}
	}

	// Take the last n.
	if len(before) > n {
		before = before[len(before)-n:]
	}

	entries := make([]Entry, len(before))
	for i, t := range before {
		entries[i] = Entry{
			Time:      t,
			Formatted: t.Format("2006-01-02 15:04:05 MST"),
		}
	}

	return &Result{
		Expression: expression,
		Entries:    entries,
	}, nil
}
