package frequency_test

import (
	"testing"

	"github.com/nicholasgasior/crontab-lint/internal/frequency"
)

func TestEstimate_EveryMinute(t *testing.T) {
	r, err := frequency.Estimate("* * * * * /bin/job")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.PerHour != 60 {
		t.Errorf("PerHour: want 60, got %.0f", r.PerHour)
	}
	if r.PerDay != 1440 {
		t.Errorf("PerDay: want 1440, got %.0f", r.PerDay)
	}
	if r.Label != "every minute" {
		t.Errorf("Label: want 'every minute', got %q", r.Label)
	}
}

func TestEstimate_Hourly(t *testing.T) {
	r, err := frequency.Estimate("0 * * * * /bin/job")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.PerHour != 1 {
		t.Errorf("PerHour: want 1, got %.0f", r.PerHour)
	}
	if r.PerDay != 24 {
		t.Errorf("PerDay: want 24, got %.0f", r.PerDay)
	}
	if r.Label != "hourly" {
		t.Errorf("Label: want 'hourly', got %q", r.Label)
	}
}

func TestEstimate_Daily(t *testing.T) {
	r, err := frequency.Estimate("0 9 * * * /bin/job")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.PerDay != 1 {
		t.Errorf("PerDay: want 1, got %.0f", r.PerDay)
	}
	if r.Label != "daily" {
		t.Errorf("Label: want 'daily', got %q", r.Label)
	}
}

func TestEstimate_StepExpression(t *testing.T) {
	// */15 * * * * fires 4 times per hour, 96 times per day
	r, err := frequency.Estimate("*/15 * * * * /bin/job")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.PerHour != 4 {
		t.Errorf("PerHour: want 4, got %.0f", r.PerHour)
	}
	if r.PerDay != 96 {
		t.Errorf("PerDay: want 96, got %.0f", r.PerDay)
	}
}

func TestEstimate_WeekdayOnly(t *testing.T) {
	// 0 9 * * 1-5  — weekdays only
	r, err := frequency.Estimate("0 9 * * 1-5 /bin/job")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.PerWeek != 5 {
		t.Errorf("PerWeek: want 5, got %.2f", r.PerWeek)
	}
	if r.Label != "~5x per week" {
		t.Errorf("Label: want '~5x per week', got %q", r.Label)
	}
}

func TestEstimate_InvalidExpression(t *testing.T) {
	_, err := frequency.Estimate("not-a-cron")
	if err == nil {
		t.Fatal("expected error for invalid expression, got nil")
	}
}

func TestEstimate_ListField(t *testing.T) {
	// 0 8,12,18 * * * — 3 times per day
	r, err := frequency.Estimate("0 8,12,18 * * * /bin/job")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.PerDay != 3 {
		t.Errorf("PerDay: want 3, got %.0f", r.PerDay)
	}
}
