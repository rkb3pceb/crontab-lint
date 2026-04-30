// Package overlap detects scheduling conflicts between two or more crontab
// expressions — for example, two jobs that would fire at the same time.
package overlap

import (
	"fmt"
	"time"

	"github.com/user/crontab-lint/internal/scheduler"
)

// Result describes a detected overlap between two crontab expressions.
type Result struct {
	// ExprA is the first crontab expression.
	ExprA string
	// ExprB is the second crontab expression.
	ExprB string
	// Collisions holds the times at which both expressions fire simultaneously.
	Collisions []time.Time
}

// Detect checks whether exprA and exprB share any execution times within the
// window [from, from+window). It returns up to maxHits collisions.
// An error is returned if either expression is invalid.
func Detect(exprA, exprB string, from time.Time, window time.Duration, maxHits int) (*Result, error) {
	if maxHits <= 0 {
		return nil, fmt.Errorf("overlap: maxHits must be greater than zero")
	}
	if window <= 0 {
		return nil, fmt.Errorf("overlap: window must be a positive duration")
	}

	end := from.Add(window)

	timesA, err := collectTimes(exprA, from, end)
	if err != nil {
		return nil, fmt.Errorf("overlap: expression A: %w", err)
	}

	timesB, err := collectTimes(exprB, from, end)
	if err != nil {
		return nil, fmt.Errorf("overlap: expression B: %w", err)
	}

	setB := make(map[time.Time]struct{}, len(timesB))
	for _, t := range timesB {
		setB[t] = struct{}{}
	}

	var collisions []time.Time
	for _, t := range timesA {
		if _, ok := setB[t]; ok {
			collisions = append(collisions, t)
			if len(collisions) >= maxHits {
				break
			}
		}
	}

	return &Result{
		ExprA:      exprA,
		ExprB:      exprB,
		Collisions: collisions,
	}, nil
}

// collectTimes gathers all firing times for expr in [from, end).
func collectTimes(expr string, from, end time.Time) ([]time.Time, error) {
	var times []time.Time
	current := from
	for {
		next, err := scheduler.Next(expr, current)
		if err != nil {
			return nil, err
		}
		if !next.Before(end) {
			break
		}
		times = append(times, next)
		current = next.Add(time.Minute)
	}
	return times, nil
}
