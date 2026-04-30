package statistics_test

import (
	"testing"

	"github.com/example/crontab-lint/internal/statistics"
)

func TestCompute_EveryMinute(t *testing.T) {
	stats, err := statistics.Compute("* * * * * /bin/task")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.RunsPerDay != 1440 {
		t.Errorf("expected 1440 runs/day, got %v", stats.RunsPerDay)
	}
	if stats.Frequency != "every minute" {
		t.Errorf("unexpected frequency: %s", stats.Frequency)
	}
}

func TestCompute_HourlyJob(t *testing.T) {
	stats, err := statistics.Compute("0 * * * * /bin/task")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.RunsPerDay != 24 {
		t.Errorf("expected 24 runs/day, got %v", stats.RunsPerDay)
	}
}

func TestCompute_DailyJob(t *testing.T) {
	stats, err := statistics.Compute("0 9 * * * /bin/task")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if stats.RunsPerDay != 1 {
		t.Errorf("expected 1 run/day, got %v", stats.RunsPerDay)
	}
	if stats.RunsPerWeek != 7 {
		t.Errorf("expected 7 runs/week, got %v", stats.RunsPerWeek)
	}
}

func TestCompute_WeekdayOnly(t *testing.T) {
	stats, err := statistics.Compute("0 9 * * 1-5 /bin/task")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 1 run/day * 2 weekdays matched by range "1-5" counted as 1 list item
	if stats.RunsPerDay != 1 {
		t.Errorf("expected 1 run/day, got %v", stats.RunsPerDay)
	}
}

func TestCompute_StepExpression(t *testing.T) {
	stats, err := statistics.Compute("*/15 * * * * /bin/task")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// 60/15 = 4 minutes per hour * 24 hours
	if stats.RunsPerDay != 96 {
		t.Errorf("expected 96 runs/day, got %v", stats.RunsPerDay)
	}
}

func TestCompute_InvalidExpression(t *testing.T) {
	_, err := statistics.Compute("not-a-cron")
	if err == nil {
		t.Error("expected error for invalid expression")
	}
}

func TestCompute_FrequencyLabel(t *testing.T) {
	tests := []struct {
		expr    string
		label   string
	}{
		{"0 9 * * * /t", "daily or more"},
		{"0 * * * * /t", "hourly or more"},
		{"*/5 * * * * /t", "high frequency (multiple times per hour)"},
	}
	for _, tt := range tests {
		stats, err := statistics.Compute(tt.expr)
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", tt.expr, err)
		}
		if stats.Frequency != tt.label {
			t.Errorf("%s: expected %q, got %q", tt.expr, tt.label, stats.Frequency)
		}
	}
}
