// Package window computes cron firing times within a bounded time range.
package window

import (
	"fmt"
	"time"

	"github.com/your-org/crontab-lint/internal/scheduler"
)

// Result holds the output of a window query.
type Result struct {
	Expression string
	From       time.Time
	To         time.Time
	Times      []time.Time
	Count      int
}

// MaxResults caps the number of firing times returned to prevent runaway output.
const MaxResults = 1000

// Build returns all times the given cron expression fires between from (inclusive)
// and to (exclusive). An error is returned if the expression is invalid or if
// from is not before to.
func Build(expr string, from, to time.Time) (Result, error) {
	if !from.Before(to) {
		return Result{}, fmt.Errorf("window: 'from' must be before 'to'")
	}

	var times []time.Time
	current := from

	for len(times) < MaxResults {
		next, err := scheduler.Next(expr, current)
		if err != nil {
			return Result{}, fmt.Errorf("window: %w", err)
		}
		if !next.Before(to) {
			break
		}
		times = append(times, next)
		current = next.Add(time.Minute)
	}

	return Result{
		Expression: expr,
		From:       from,
		To:         to,
		Times:      times,
		Count:      len(times),
	}, nil
}
