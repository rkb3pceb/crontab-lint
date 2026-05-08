// Package ranking scores and ranks multiple crontab expressions by
// overall quality, combining complexity, frequency, and lint severity.
package ranking

import (
	"sort"

	"github.com/example/crontab-lint/internal/complexity"
	"github.com/example/crontab-lint/internal/linter"
	"github.com/example/crontab-lint/internal/statistics"
)

// Entry holds a crontab expression and its computed rank score.
type Entry struct {
	Expression string
	Score      int    // higher is better
	Grade      string // A–F derived from score
	Errors     int
	Warnings   int
}

// Rank evaluates each expression and returns entries sorted best-first.
// Expressions that fail to parse are included last with score 0.
func Rank(expressions []string) []Entry {
	entries := make([]Entry, 0, len(expressions))
	for _, expr := range expressions {
		entries = append(entries, evaluate(expr))
	}
	sort.SliceStable(entries, func(i, j int) bool {
		return entries[i].Score > entries[j].Score
	})
	return entries
}

func evaluate(expr string) Entry {
	e := Entry{Expression: expr}

	results := linter.Lint(expr)
	for _, r := range results {
		switch r.Severity {
		case "error":
			e.Errors++
		case "warning":
			e.Warnings++
		}
	}

	if e.Errors > 0 {
		e.Grade = "F"
		return e
	}

	cx := complexity.Score(expr)
	st, statsErr := statistics.Compute(expr)

	// Base score from complexity (lower complexity → higher base)
	base := 100 - cx.Total
	if base < 0 {
		base = 0
	}

	// Penalise by frequency if stats available
	freqPenalty := 0
	if statsErr == nil {
		freqPenalty = st.RunsPerDay / 10
	}

	warningPenalty := e.Warnings * 5

	e.Score = base - freqPenalty - warningPenalty
	if e.Score < 0 {
		e.Score = 0
	}
	e.Grade = gradeFrom(e.Score)
	return e
}

func gradeFrom(score int) string {
	switch {
	case score >= 80:
		return "A"
	case score >= 60:
		return "B"
	case score >= 40:
		return "C"
	case score >= 20:
		return "D"
	default:
		return "F"
	}
}
