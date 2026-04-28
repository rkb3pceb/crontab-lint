// Package suggester provides fix suggestions for common crontab mistakes.
package suggester

import (
	"fmt"
	"strings"
)

// Suggestion holds a human-readable fix recommendation.
type Suggestion struct {
	Field   string
	Message string
	Fix     string
}

// Suggest returns a list of fix suggestions based on a lint warning or error message.
func Suggest(field, value, lintMsg string) []Suggestion {
	var suggestions []Suggestion

	switch {
	case strings.Contains(lintMsg, "every minute"):
		suggestions = append(suggestions, Suggestion{
			Field:   field,
			Message: "Running every minute may cause high system load.",
			Fix:     "Consider '*/5 * * * *' to run every 5 minutes instead.",
		})

	case strings.Contains(lintMsg, "high frequency"):
		suggestions = append(suggestions, Suggestion{
			Field:   field,
			Message: "Interval is very short (less than 5 minutes).",
			Fix:     fmt.Sprintf("Replace '%s' with '*/5' or a larger step to reduce frequency.", value),
		})

	case strings.Contains(lintMsg, "both DOM and DOW"):
		suggestions = append(suggestions, Suggestion{
			Field:   "dom+dow",
			Message: "Setting both day-of-month and day-of-week uses OR logic, which may be unintended.",
			Fix:     "Set one of them to '*' to avoid ambiguous scheduling.",
		})

	case strings.Contains(lintMsg, "unreachable"):
		suggestions = append(suggestions, Suggestion{
			Field:   field,
			Message: fmt.Sprintf("Day-of-month value '%s' may never occur in some months.", value),
			Fix:     "Use a value between 1 and 28 to ensure the job runs every month.",
		})

	case strings.Contains(lintMsg, "out of range"), strings.Contains(lintMsg, "invalid"):
		suggestions = append(suggestions, Suggestion{
			Field:   field,
			Message: fmt.Sprintf("Value '%s' is not valid for field '%s'.", value, field),
			Fix:     suggestValidRange(field),
		})
	}

	return suggestions
}

// suggestValidRange returns the valid range hint for a given cron field name.
func suggestValidRange(field string) string {
	ranges := map[string]string{
		"minute":      "Use a value between 0 and 59.",
		"hour":        "Use a value between 0 and 23.",
		"dom":         "Use a value between 1 and 31.",
		"month":       "Use a value between 1 and 12 (or JAN-DEC).",
		"dow":         "Use a value between 0 and 7 (or SUN-SAT; 0 and 7 both mean Sunday).",
	}
	if hint, ok := ranges[strings.ToLower(field)]; ok {
		return hint
	}
	return "Check the cron field documentation for valid values."
}
