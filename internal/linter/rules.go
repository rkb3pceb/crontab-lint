package linter

// rules.go defines the individual lint rules applied to parsed crontab expressions.
// Each rule is a function that inspects a CronExpression and returns zero or more LintIssues.

import (
	"fmt"
	"strings"

	"github.com/yourorg/crontab-lint/internal/parser"
)

// Rule is a function that evaluates a single lint concern against a parsed expression.
type Rule func(expr parser.CronExpression) []LintIssue

// DefaultRules returns the ordered set of rules applied during linting.
func DefaultRules() []Rule {
	return []Rule{
		RuleEveryMinute,
		RuleHighFrequency,
		RuleDomAndDowBothSet,
		RuleUnreachableDayOfMonth,
		RuleMidnightAmbiguity,
		RuleLeapDaySchedule,
	}
}

// RuleEveryMinute warns when the schedule runs every minute (* * * * *).
func RuleEveryMinute(expr parser.CronExpression) []LintIssue {
	if expr.Minute == "*" && expr.Hour == "*" && expr.DayOfMonth == "*" &&
		expr.Month == "*" && expr.DayOfWeek == "*" {
		return []LintIssue{{
			Severity: SeverityWarning,
			Message:  "schedule runs every minute, which may cause high load",
			Field:    "minute",
		}}
	}
	return nil
}

// RuleHighFrequency warns when the minute field uses a step that fires more than
// 12 times per hour (e.g. */4 or smaller).
func RuleHighFrequency(expr parser.CronExpression) []LintIssue {
	if strings.HasPrefix(expr.Minute, "*/") {
		var step int
		_, err := fmt.Sscanf(expr.Minute, "*/%d", &step)
		if err == nil && step > 0 && (60/step) > 12 {
			return []LintIssue{{
				Severity: SeverityWarning,
				Message:  fmt.Sprintf("schedule fires %d times per hour; consider a less frequent interval", 60/step),
				Field:    "minute",
			}}
		}
	}
	return nil
}

// RuleDomAndDowBothSet warns when both day-of-month and day-of-week are
// explicitly set, because cron treats this as a union (OR), which is often
// surprising to users.
func RuleDomAndDowBothSet(expr parser.CronExpression) []LintIssue {
	domSet := expr.DayOfMonth != "*" && expr.DayOfMonth != "?"
	dowSet := expr.DayOfWeek != "*" && expr.DayOfWeek != "?"
	if domSet && dowSet {
		return []LintIssue{{
			Severity: SeverityWarning,
			Message:  "both day-of-month and day-of-week are set; cron will trigger on either condition (OR logic), which may be unintentional",
			Field:    "day-of-month",
		}}
	}
	return nil
}

// RuleUnreachableDayOfMonth warns when day-of-month is set to 31 but month is
// one of the months that never has 31 days.
func RuleUnreachableDayOfMonth(expr parser.CronExpression) []LintIssue {
	// Months with fewer than 31 days: 2 (Feb), 4 (Apr), 6 (Jun), 9 (Sep), 11 (Nov).
	shortMonths := map[string]bool{
		"2": true, "4": true, "6": true, "9": true, "11": true,
	}
	if expr.DayOfMonth == "31" && shortMonths[expr.Month] {
		return []LintIssue{{
			Severity: SeverityError,
			Message:  fmt.Sprintf("day 31 never occurs in month %s; the job will never run", expr.Month),
			Field:    "day-of-month",
		}}
	}
	return nil
}

// RuleMidnightAmbiguity warns when hour is 0 and minute is 0 but the intent
// may be ambiguous due to daylight-saving-time transitions.
func RuleMidnightAmbiguity(expr parser.CronExpression) []LintIssue {
	if expr.Minute == "0" && expr.Hour == "0" {
		return []LintIssue{{
			Severity: SeverityInfo,
			Message:  "schedule runs at midnight; be aware that DST transitions can cause this job to run twice or be skipped depending on your cron daemon",
			Field:    "hour",
		}}
	}
	return nil
}

// RuleLeapDaySchedule warns when the schedule is set to February 29, which
// only exists in leap years.
func RuleLeapDaySchedule(expr parser.CronExpression) []LintIssue {
	if expr.DayOfMonth == "29" && expr.Month == "2" {
		return []LintIssue{{
			Severity: SeverityWarning,
			Message:  "schedule is set to February 29 (leap day); this job will only run in leap years",
			Field:    "day-of-month",
		}}
	}
	return nil
}
