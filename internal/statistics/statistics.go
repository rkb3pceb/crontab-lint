// Package statistics provides frequency and execution analysis
// for crontab expressions, including estimated runs per day/week/month.
package statistics

import (
	"fmt"
	"math"

	"github.com/example/crontab-lint/internal/parser"
)

// Stats holds computed frequency statistics for a cron expression.
type Stats struct {
	RunsPerDay   float64
	RunsPerWeek  float64
	RunsPerMonth float64
	Frequency    string
}

// Compute returns frequency statistics for the given cron expression.
// It estimates execution counts based on field cardinality.
func Compute(expr string) (*Stats, error) {
	entry, err := parser.Parse(expr)
	if err != nil {
		return nil, fmt.Errorf("statistics: %w", err)
	}

	minuteCount := countField(entry.Minute, 0, 59)
	hourCount := countField(entry.Hour, 0, 23)
	dayCount := countField(entry.DayOfMonth, 1, 31)
	weekdayCount := countField(entry.DayOfWeek, 0, 6)

	// Effective days per week: use DOW if restricted, else assume 7
	effectiveDaysPerWeek := float64(weekdayCount)
	if entry.DayOfWeek == "*" {
		effectiveDaysPerWeek = 7
	}

	// Effective days per month
	effectiveDaysPerMonth := float64(dayCount)
	if entry.DayOfMonth == "*" {
		effectiveDaysPerMonth = 30
	}

	runsPerDay := float64(minuteCount) * float64(hourCount)
	runsPerWeek := runsPerDay * effectiveDaysPerWeek
	runsPerMonth := runsPerDay * effectiveDaysPerMonth

	return &Stats{
		RunsPerDay:   math.Round(runsPerDay*100) / 100,
		RunsPerWeek:  math.Round(runsPerWeek*100) / 100,
		RunsPerMonth: math.Round(runsPerMonth*100) / 100,
		Frequency:    describeFrequency(runsPerDay),
	}, nil
}

// countField estimates how many distinct values a cron field matches
// within the given inclusive range [min, max].
func countField(field string, min, max int) int {
	if field == "*" {
		return max - min + 1
	}
	total := max - min + 1
	// step expression e.g. */5
	if len(field) > 2 && field[:2] == "*/" {
		step := 1
		fmt.Sscanf(field[2:], "%d", &step)
		if step < 1 {
			step = 1
		}
		return (total + step - 1) / step
	}
	// list: count commas + 1
	count := 1
	for _, c := range field {
		if c == ',' {
			count++
		}
	}
	return count
}

func describeFrequency(runsPerDay float64) string {
	switch {
	case runsPerDay >= 1440:
		return "every minute"
	case runsPerDay >= 60:
		return "high frequency (multiple times per hour)"
	case runsPerDay >= 24:
		return "hourly or more"
	case runsPerDay >= 1:
		return "daily or more"
	default:
		return "less than once per day"
	}
}
