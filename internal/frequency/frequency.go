// Package frequency estimates how often a cron expression fires per unit of
// time and exposes the result as a structured value with a human-readable label.
package frequency

import (
	"fmt"
	"strings"

	"github.com/nicholasgasior/crontab-lint/internal/parser"
)

// Result holds the estimated execution frequency for a cron expression.
type Result struct {
	// PerDay is the approximate number of times the job fires in a 24-hour period.
	PerDay float64
	// PerHour is the approximate number of times the job fires in one hour.
	PerHour float64
	// PerWeek is the approximate number of times the job fires in a 7-day week.
	PerWeek float64
	// Label is a short human-readable description such as "every minute" or
	// "~3x per day".
	Label string
}

// Estimate parses expr and returns a frequency Result.
// An error is returned when the expression cannot be parsed.
func Estimate(expr string) (Result, error) {
	entry, err := parser.Parse(expr)
	if err != nil {
		return Result{}, fmt.Errorf("frequency: %w", err)
	}

	minuteCount := countField(entry.Minute, 0, 59)
	hourCount := countField(entry.Hour, 0, 23)
	dayCount := countField(entry.DayOfMonth, 1, 31)
	weekdayCount := countField(entry.DayOfWeek, 0, 6)

	// Effective days per week: if both DOM and DOW are wildcards use 7,
	// otherwise take the more restrictive constraint.
	domWild := strings.TrimSpace(entry.DayOfMonth) == "*"
	dowWild := strings.TrimSpace(entry.DayOfWeek) == "*"

	var daysPerWeek float64
	switch {
	case domWild && dowWild:
		daysPerWeek = 7
	case domWild:
		daysPerWeek = float64(weekdayCount)
	case dowWild:
		daysPerWeek = float64(dayCount) * 7.0 / 31.0
	default:
		// Both set — cron fires when either matches; approximate union.
		domFrac := float64(dayCount) / 31.0
		dowFrac := float64(weekdayCount) / 7.0
		union := domFrac + dowFrac - domFrac*dowFrac
		daysPerWeek = union * 7.0
	}

	firesPerHour := float64(minuteCount)
	firesPerDay := firesPerHour * float64(hourCount)
	firesPerWeek := firesPerDay * daysPerWeek

	return Result{
		PerHour: firesPerHour,
		PerDay:  firesPerDay,
		PerWeek: firesPerWeek,
		Label:   label(firesPerHour, firesPerDay, firesPerWeek),
	}, nil
}

// countField returns the number of distinct values a cron field matches within
// [min, max] (inclusive). A wildcard "*" matches all values in the range.
func countField(field string, min, max int) int {
	field = strings.TrimSpace(field)
	if field == "*" {
		return max - min + 1
	}
	if strings.Contains(field, "/") {
		parts := strings.SplitN(field, "/", 2)
		step := 1
		fmt.Sscanf(parts[1], "%d", &step)
		if step < 1 {
			step = 1
		}
		return (max-min)/step + 1
	}
	if strings.Contains(field, ",") {
		return len(strings.Split(field, ","))
	}
	if strings.Contains(field, "-") {
		var lo, hi int
		fmt.Sscanf(field, "%d-%d", &lo, &hi)
		if hi >= lo {
			return hi - lo + 1
		}
	}
	return 1
}

func label(perHour, perDay, perWeek float64) string {
	switch {
	case perHour >= 60:
		return "every minute"
	case perHour >= 30:
		return "every 2 minutes"
	case perHour >= 12:
		return "multiple times per hour"
	case perHour >= 2:
		return fmt.Sprintf("~%.0fx per hour", perHour)
	case perDay >= 24:
		return "hourly"
	case perDay >= 2:
		return fmt.Sprintf("~%.0fx per day", perDay)
	case perDay >= 1:
		return "daily"
	case perWeek >= 1:
		return fmt.Sprintf("~%.0fx per week", perWeek)
	default:
		return "less than once per week"
	}
}
