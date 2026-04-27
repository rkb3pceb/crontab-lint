package humanizer

import (
	"fmt"
	"strings"
)

// FieldType represents the type of a cron field.
type FieldType int

const (
	Minute FieldType = iota
	Hour
	DayOfMonth
	Month
	DayOfWeek
)

var fieldNames = map[FieldType]string{
	Minute:     "minute",
	Hour:       "hour",
	DayOfMonth: "day of month",
	Month:      "month",
	DayOfWeek:  "day of week",
}

var monthNames = []string{
	"", "January", "February", "March", "April", "May", "June",
	"July", "August", "September", "October", "November", "December",
}

var dayNames = []string{
	"Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday",
}

// Describe returns a human-readable description of a single cron field value.
func Describe(field string, ft FieldType) string {
	if field == "*" {
		return fmt.Sprintf("every %s", fieldNames[ft])
	}

	if strings.Contains(field, "/") {
		parts := strings.SplitN(field, "/", 2)
		return fmt.Sprintf("every %s %s(s)", parts[1], fieldNames[ft])
	}

	if strings.Contains(field, "-") {
		parts := strings.SplitN(field, "-", 2)
		return fmt.Sprintf("from %s to %s", labelFor(parts[0], ft), labelFor(parts[1], ft))
	}

	if strings.Contains(field, ",") {
		parts := strings.Split(field, ",")
		labels := make([]string, len(parts))
		for i, p := range parts {
			labels[i] = labelFor(p, ft)
		}
		return strings.Join(labels, ", ")
	}

	return fmt.Sprintf("at %s %s", labelFor(field, ft), fieldNames[ft])
}

// DescribeSchedule returns a full human-readable description of a 5-field cron expression.
func DescribeSchedule(fields []string) string {
	if len(fields) != 5 {
		return "invalid cron expression"
	}
	parts := []string{
		Describe(fields[0], Minute),
		Describe(fields[1], Hour),
		Describe(fields[2], DayOfMonth),
		Describe(fields[3], Month),
		Describe(fields[4], DayOfWeek),
	}
	return strings.Join(parts, "; ")
}

func labelFor(value string, ft FieldType) string {
	switch ft {
	case Month:
		if n := parseIntSafe(value); n >= 1 && n <= 12 {
			return monthNames[n]
		}
	case DayOfWeek:
		if n := parseIntSafe(value); n >= 0 && n <= 6 {
			return dayNames[n]
		}
	}
	return value
}

func parseIntSafe(s string) int {
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	if err != nil {
		return -1
	}
	return n
}
