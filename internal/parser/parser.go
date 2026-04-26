package parser

import (
	"fmt"
	"strconv"
	"strings"
)

// Field represents a single crontab field with its constraints.
type Field struct {
	Name string
	Min  int
	Max  int
}

// CronExpression holds the parsed fields of a crontab expression.
type CronExpression struct {
	Minute     string
	Hour       string
	DayOfMonth string
	Month      string
	DayOfWeek  string
	Command    string
}

var fields = []Field{
	{Name: "minute", Min: 0, Max: 59},
	{Name: "hour", Min: 0, Max: 23},
	{Name: "day of month", Min: 1, Max: 31},
	{Name: "month", Min: 1, Max: 12},
	{Name: "day of week", Min: 0, Max: 6},
}

// Parse splits a raw crontab line into a CronExpression.
func Parse(line string) (*CronExpression, error) {
	line = strings.TrimSpace(line)
	parts := strings.Fields(line)
	if len(parts) < 6 {
		return nil, fmt.Errorf("expected at least 6 fields, got %d", len(parts))
	}
	return &CronExpression{
		Minute:     parts[0],
		Hour:       parts[1],
		DayOfMonth: parts[2],
		Month:      parts[3],
		DayOfWeek:  parts[4],
		Command:    strings.Join(parts[5:], " "),
	}, nil
}

// ValidateField checks whether a cron field value is valid for the given constraints.
func ValidateField(value string, f Field) error {
	if value == "*" {
		return nil
	}
	// Handle step values: */2 or 1-5/2
	if strings.Contains(value, "/") {
		return validateStep(value, f)
	}
	// Handle ranges: 1-5
	if strings.Contains(value, "-") {
		return validateRange(value, f)
	}
	// Handle lists: 1,2,3
	if strings.Contains(value, ",") {
		return validateList(value, f)
	}
	return validateNumber(value, f)
}

func validateNumber(value string, f Field) error {
	n, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("%s: %q is not a valid number", f.Name, value)
	}
	if n < f.Min || n > f.Max {
		return fmt.Errorf("%s: %d out of range [%d-%d]", f.Name, n, f.Min, f.Max)
	}
	return nil
}

func validateRange(value string, f Field) error {
	parts := strings.SplitN(value, "-", 2)
	if err := validateNumber(parts[0], f); err != nil {
		return err
	}
	return validateNumber(parts[1], f)
}

func validateList(value string, f Field) error {
	for _, part := range strings.Split(value, ",") {
		if err := validateNumber(part, f); err != nil {
			return err
		}
	}
	return nil
}

func validateStep(value string, f Field) error {
	parts := strings.SplitN(value, "/", 2)
	if parts[0] != "*" {
		if err := validateRange(parts[0], f); err != nil {
			return err
		}
	}
	step, err := strconv.Atoi(parts[1])
	if err != nil || step < 1 {
		return fmt.Errorf("%s: step value %q must be a positive integer", f.Name, parts[1])
	}
	return nil
}

// Fields returns the ordered list of cron fields with their constraints.
func Fields() []Field {
	return fields
}
