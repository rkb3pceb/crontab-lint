// Package normalizer provides utilities for normalizing crontab expressions
// into a canonical form, expanding aliases and shorthand notations.
package normalizer

import (
	"strings"
)

// aliases maps common crontab shorthand expressions to their canonical form.
var aliases = map[string]string{
	"@yearly":   "0 0 1 1 *",
	"@annually": "0 0 1 1 *",
	"@monthly":  "0 0 1 * *",
	"@weekly":   "0 0 * * 0",
	"@daily":    "0 0 * * *",
	"@midnight": "0 0 * * *",
	"@hourly":   "0 * * * *",
}

// monthNames maps month name abbreviations to their numeric equivalents.
var monthNames = map[string]string{
	"jan": "1", "feb": "2", "mar": "3", "apr": "4",
	"may": "5", "jun": "6", "jul": "7", "aug": "8",
	"sep": "9", "oct": "10", "nov": "11", "dec": "12",
}

// dowNames maps day-of-week name abbreviations to their numeric equivalents.
var dowNames = map[string]string{
	"sun": "0", "mon": "1", "tue": "2", "wed": "3",
	"thu": "4", "fri": "5", "sat": "6",
}

// Normalize takes a raw crontab expression and returns a canonical 5-field
// expression. Aliases like @daily are expanded, and named months/days are
// replaced with their numeric equivalents. The original command portion
// (if present) is preserved unchanged.
func Normalize(expr string) string {
	expr = strings.TrimSpace(expr)

	// Check for alias shorthand (no command)
	if canonical, ok := aliases[strings.ToLower(expr)]; ok {
		return canonical
	}

	// Check for alias shorthand with a command: "@daily /usr/bin/backup"
	for alias, canonical := range aliases {
		if strings.HasPrefix(strings.ToLower(expr), alias+" ") {
			command := strings.TrimSpace(expr[len(alias):])
			return canonical + " " + command
		}
	}

	parts := strings.Fields(expr)
	if len(parts) < 5 {
		return expr
	}

	fields := make([]string, 5)
	for i := 0; i < 5; i++ {
		fields[i] = normalizeField(parts[i], i)
	}

	result := strings.Join(fields, " ")
	if len(parts) > 5 {
		result += " " + strings.Join(parts[5:], " ")
	}
	return result
}

// normalizeField replaces named values in a single cron field with numbers.
func normalizeField(field string, index int) string {
	field = strings.ToLower(field)
	switch index {
	case 3: // month
		return replaceNames(field, monthNames)
	case 4: // day of week
		return replaceNames(field, dowNames)
	}
	return field
}

// replaceNames substitutes named tokens within a field expression.
func replaceNames(field string, names map[string]string) string {
	for name, num := range names {
		field = strings.ReplaceAll(field, name, num)
	}
	return field
}
