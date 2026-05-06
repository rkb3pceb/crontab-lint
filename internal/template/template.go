// Package template provides named crontab expression templates for common
// scheduling patterns, allowing users to look up or search well-known schedules.
package template

import (
	"fmt"
	"strings"
)

// Template represents a named crontab schedule with a description.
type Template struct {
	Name        string
	Expression  string
	Description string
	Tags        []string
}

// catalog holds all built-in templates.
var catalog = []Template{
	{Name: "every-minute", Expression: "* * * * *", Description: "Run every minute", Tags: []string{"frequent", "debug"}},
	{Name: "hourly", Expression: "0 * * * *", Description: "Run at the start of every hour", Tags: []string{"hourly"}},
	{Name: "daily-midnight", Expression: "0 0 * * *", Description: "Run once a day at midnight", Tags: []string{"daily"}},
	{Name: "daily-noon", Expression: "0 12 * * *", Description: "Run once a day at noon", Tags: []string{"daily"}},
	{Name: "weekly-sunday", Expression: "0 0 * * 0", Description: "Run once a week on Sunday at midnight", Tags: []string{"weekly"}},
	{Name: "weekly-monday", Expression: "0 0 * * 1", Description: "Run once a week on Monday at midnight", Tags: []string{"weekly"}},
	{Name: "monthly", Expression: "0 0 1 * *", Description: "Run once a month on the 1st at midnight", Tags: []string{"monthly"}},
	{Name: "yearly", Expression: "0 0 1 1 *", Description: "Run once a year on January 1st at midnight", Tags: []string{"yearly"}},
	{Name: "every-5-minutes", Expression: "*/5 * * * *", Description: "Run every 5 minutes", Tags: []string{"frequent"}},
	{Name: "every-15-minutes", Expression: "*/15 * * * *", Description: "Run every 15 minutes", Tags: []string{"frequent"}},
	{Name: "every-30-minutes", Expression: "*/30 * * * *", Description: "Run every 30 minutes", Tags: []string{"frequent"}},
	{Name: "weekdays", Expression: "0 9 * * 1-5", Description: "Run at 9am on weekdays (Mon–Fri)", Tags: []string{"weekday", "business"}},
	{Name: "weekends", Expression: "0 10 * * 6,0", Description: "Run at 10am on weekends (Sat & Sun)", Tags: []string{"weekend"}},
}

// All returns a copy of all available templates.
func All() []Template {
	out := make([]Template, len(catalog))
	copy(out, catalog)
	return out
}

// Lookup returns the template with the given name, or an error if not found.
func Lookup(name string) (Template, error) {
	for _, t := range catalog {
		if t.Name == name {
			return t, nil
		}
	}
	return Template{}, fmt.Errorf("template %q not found", name)
}

// Search returns all templates whose name, description, or tags contain the
// given query string (case-insensitive).
func Search(query string) []Template {
	q := strings.ToLower(query)
	var results []Template
	for _, t := range catalog {
		if strings.Contains(strings.ToLower(t.Name), q) ||
			strings.Contains(strings.ToLower(t.Description), q) ||
			containsTag(t.Tags, q) {
			results = append(results, t)
		}
	}
	return results
}

func containsTag(tags []string, q string) bool {
	for _, tag := range tags {
		if strings.Contains(tag, q) {
			return true
		}
	}
	return false
}
