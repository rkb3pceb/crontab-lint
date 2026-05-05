package calendar_test

import (
	"strings"
	"testing"
	"time"

	"github.com/example/crontab-lint/internal/calendar"
)

func day(year int, month time.Month, d int) time.Time {
	return time.Date(year, month, d, 0, 0, 0, 0, time.UTC)
}

func TestBuildDay_EveryHour(t *testing.T) {
	view, err := calendar.BuildDay("0 * * * * echo hi", day(2024, time.January, 15))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if view.Total != 24 {
		t.Errorf("expected 24 executions, got %d", view.Total)
	}
	for h := 0; h < 24; h++ {
		if !view.Blocks[h].Fired {
			t.Errorf("expected hour %d to fire", h)
		}
	}
}

func TestBuildDay_EveryMinute(t *testing.T) {
	view, err := calendar.BuildDay("* * * * * echo hi", day(2024, time.January, 15))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if view.Total != 1440 {
		t.Errorf("expected 1440 executions, got %d", view.Total)
	}
}

func TestBuildDay_SpecificTime(t *testing.T) {
	view, err := calendar.BuildDay("30 9 * * * echo hi", day(2024, time.March, 1))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if view.Total != 1 {
		t.Errorf("expected 1 execution, got %d", view.Total)
	}
	if !view.Blocks[9].Fired {
		t.Error("expected hour 9 to fire")
	}
	if len(view.Blocks[9].Minutes) != 1 || view.Blocks[9].Minutes[0] != 30 {
		t.Errorf("expected minute 30, got %v", view.Blocks[9].Minutes)
	}
}

func TestBuildDay_InvalidExpression(t *testing.T) {
	_, err := calendar.BuildDay("invalid", day(2024, time.January, 1))
	if err == nil {
		t.Error("expected error for invalid expression")
	}
}

func TestFormatDay_ContainsDate(t *testing.T) {
	view, err := calendar.BuildDay("0 12 * * * echo hi", day(2024, time.June, 5))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := calendar.FormatDay(view)
	if !strings.Contains(out, "2024-06-05") {
		t.Errorf("expected date in output, got:\n%s", out)
	}
	if !strings.Contains(out, "12:00") {
		t.Errorf("expected 12:00 in output, got:\n%s", out)
	}
}

func TestFormatDay_NoFirings(t *testing.T) {
	// Feb 30 doesn't exist — use a DOM that won't match the weekday
	// Use a specific DOM+DOW combo that won't fire on this date
	view, _ := calendar.BuildDay("0 0 31 2 * echo hi", day(2024, time.January, 15))
	out := calendar.FormatDay(view)
	if !strings.Contains(out, "0 execution") {
		t.Errorf("expected 0 executions in output, got:\n%s", out)
	}
}
