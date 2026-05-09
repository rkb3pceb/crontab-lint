// Package timeline generates a chronological view of cron job executions
// over a specified time window, useful for visualising scheduling density.
package timeline

import (
	"fmt"
	"time"

	"github.com/dkgv/crontab-lint/internal/scheduler"
)

// Entry represents a single scheduled execution moment within the window.
type Entry struct {
	At    time.Time
	Label string
}

// Result holds the timeline entries for one cron expression.
type Result struct {
	Expression string
	Window     time.Duration
	Entries    []Entry
	Truncated  bool
}

const maxEntries = 200

// Build generates timeline entries for expr within [from, from+window).
// At most maxEntries entries are returned; Truncated is set if more exist.
func Build(expr string, from time.Time, window time.Duration) (Result, error) {
	if window <= 0 {
		return Result{}, fmt.Errorf("timeline: window must be positive")
	}
	until := from.Add(window)

	var entries []Entry
	truncated := false
	current := from

	for {
		next, err := scheduler.Next(expr, current)
		if err != nil {
			return Result{}, fmt.Errorf("timeline: %w", err)
		}
		if !next.Before(until) {
			break
		}
		if len(entries) >= maxEntries {
			truncated = true
			break
		}
		entries = append(entries, Entry{
			At:    next,
			Label: next.Format("Mon 02 Jan 2006 15:04"),
		})
		current = next.Add(time.Minute)
	}

	return Result{
		Expression: expr,
		Window:     window,
		Entries:    entries,
		Truncated:  truncated,
	}, nil
}
