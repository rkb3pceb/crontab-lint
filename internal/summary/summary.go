// Package summary provides a high-level textual summary of a crontab expression,
// combining humanized schedule description, complexity grade, and next run times.
package summary

import (
	"fmt"
	"strings"
	"time"

	"github.com/yourorg/crontab-lint/internal/complexity"
	"github.com/yourorg/crontab-lint/internal/humanizer"
	"github.com/yourorg/crontab-lint/internal/scheduler"
)

// Result holds all summary information for a crontab expression.
type Result struct {
	Expression  string
	Description string
	Complexity  complexity.Result
	NextRuns    []time.Time
	Error       string
}

// Summarize produces a Result for the given crontab expression.
// nextCount controls how many upcoming run times are included (capped at 5).
func Summarize(expr string, from time.Time, nextCount int) Result {
	if nextCount < 1 {
		nextCount = 1
	}
	if nextCount > 5 {
		nextCount = 5
	}

	desc, err := humanizer.Describe(expr)
	if err != nil {
		return Result{
			Expression: expr,
			Error:      err.Error(),
		}
	}

	cx := complexity.Score(expr)

	next, err := scheduler.NextN(expr, from, nextCount)
	if err != nil {
		return Result{
			Expression:  expr,
			Description: desc,
			Complexity:  cx,
			Error:       err.Error(),
		}
	}

	return Result{
		Expression:  expr,
		Description: desc,
		Complexity:  cx,
		NextRuns:    next,
	}
}

// Format returns a human-readable multi-line string for the summary result.
func Format(r Result) string {
	if r.Error != "" {
		return fmt.Sprintf("Expression : %s\nError      : %s\n", r.Expression, r.Error)
	}

	var sb strings.Builder
	fmt.Fprintf(&sb, "Expression : %s\n", r.Expression)
	fmt.Fprintf(&sb, "Description: %s\n", r.Description)
	fmt.Fprintf(&sb, "Complexity : %s (score %d)\n", r.Complexity.Grade, r.Complexity.Score)
	if len(r.NextRuns) > 0 {
		sb.WriteString("Next runs  :\n")
		for _, t := range r.NextRuns {
			fmt.Fprintf(&sb, "  - %s\n", t.UTC().Format(time.RFC3339))
		}
	}
	return sb.String()
}
