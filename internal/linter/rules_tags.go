package linter

import (
	"fmt"

	"github.com/example/crontab-lint/internal/tags"
)

// RuleTagHighFrequency warns when a crontab expression is tagged as high-frequency.
func RuleTagHighFrequency(expr string) *Result {
	t, err := tags.Extract(expr)
	if err != nil {
		return nil
	}
	for _, tag := range t {
		if tag.Name == "high-frequency" {
			return &Result{
				Level:   Warn,
				Message: "expression is tagged high-frequency: runs every minute",
				Code:    "TAG001",
			}
		}
	}
	return nil
}

// RuleTagDomSpecific warns when both dom and dow are constrained simultaneously.
func RuleTagDomAndDowSpecific(expr string) *Result {
	t, err := tags.Extract(expr)
	if err != nil {
		return nil
	}
	names := tags.Names(t)
	hasDom := containsTag(names, "dom-specific")
	hasDow := containsTag(names, "weekday-specific")
	if hasDom && hasDow {
		return &Result{
			Level:   Warn,
			Message: fmt.Sprintf("both day-of-month and day-of-week are constrained; cron uses OR semantics"),
			Code:    "TAG002",
		}
	}
	return nil
}

func containsTag(names []string, name string) bool {
	for _, n := range names {
		if n == name {
			return true
		}
	}
	return false
}

func init() {
	DefaultRules = append(DefaultRules, RuleTagHighFrequency, RuleTagDomAndDowSpecific)
}
