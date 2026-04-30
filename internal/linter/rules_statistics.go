package linter

import (
	"fmt"

	"github.com/example/crontab-lint/internal/statistics"
)

// RuleHighDailyFrequency warns when a cron job runs more than 96 times per day
// (i.e., more than once every 15 minutes), which may indicate a misconfiguration
// or an unexpectedly resource-intensive schedule.
func RuleHighDailyFrequency(expr string) *Warning {
	stats, err := statistics.Compute(expr)
	if err != nil {
		return nil
	}
	if stats.RunsPerDay > 96 && stats.RunsPerDay < 1440 {
		return &Warning{
			Code:    "HIGH_DAILY_FREQUENCY",
			Message: fmt.Sprintf("job runs %.0f times per day; consider whether this frequency is intentional", stats.RunsPerDay),
		}
	}
	return nil
}

// RuleExcessiveMonthlyRuns warns when a job is estimated to run more than
// 10,000 times per month, flagging potential infrastructure cost concerns.
func RuleExcessiveMonthlyRuns(expr string) *Warning {
	stats, err := statistics.Compute(expr)
	if err != nil {
		return nil
	}
	if stats.RunsPerMonth > 10000 {
		return &Warning{
			Code:    "EXCESSIVE_MONTHLY_RUNS",
			Message: fmt.Sprintf("job is estimated to run %.0f times per month; review if this is expected", stats.RunsPerMonth),
		}
	}
	return nil
}

func init() {
	// Register statistics-based rules into DefaultRules at package init.
	DefaultRules = append(DefaultRules, RuleHighDailyFrequency, RuleExcessiveMonthlyRuns)
}
