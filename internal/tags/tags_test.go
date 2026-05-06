package tags_test

import (
	"testing"

	"github.com/example/crontab-lint/internal/tags"
)

func TestExtract_EveryMinute(t *testing.T) {
	result, err := tags.Extract("* * * * * /bin/job")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	names := tags.Names(result)
	if !contains(names, "high-frequency") {
		t.Errorf("expected high-frequency tag, got %v", names)
	}
}

func TestExtract_Hourly(t *testing.T) {
	result, err := tags.Extract("0 * * * * /bin/job")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contains(tags.Names(result), "hourly") {
		t.Errorf("expected hourly tag, got %v", tags.Names(result))
	}
}

func TestExtract_Daily(t *testing.T) {
	result, err := tags.Extract("0 0 * * * /bin/job")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contains(tags.Names(result), "daily") {
		t.Errorf("expected daily tag, got %v", tags.Names(result))
	}
}

func TestExtract_WeekdaySpecific(t *testing.T) {
	result, err := tags.Extract("0 9 * * 1-5 /bin/job")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contains(tags.Names(result), "weekday-specific") {
		t.Errorf("expected weekday-specific tag, got %v", tags.Names(result))
	}
}

func TestExtract_MonthRestricted(t *testing.T) {
	result, err := tags.Extract("0 0 1 6 * /bin/job")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	names := tags.Names(result)
	if !contains(names, "month-restricted") {
		t.Errorf("expected month-restricted tag, got %v", names)
	}
}

func TestExtract_Interval(t *testing.T) {
	result, err := tags.Extract("*/15 * * * * /bin/job")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !contains(tags.Names(result), "interval") {
		t.Errorf("expected interval tag, got %v", tags.Names(result))
	}
}

func TestExtract_InvalidExpression(t *testing.T) {
	_, err := tags.Extract("not a crontab")
	if err == nil {
		t.Error("expected error for invalid expression")
	}
}

func TestNames_ReturnsStrings(t *testing.T) {
	input := []tags.Tag{
		{Name: "daily", Description: "Runs once per day"},
		{Name: "custom", Description: "Custom schedule"},
	}
	got := tags.Names(input)
	if len(got) != 2 || got[0] != "daily" || got[1] != "custom" {
		t.Errorf("unexpected names: %v", got)
	}
}

func contains(slice []string, val string) bool {
	for _, s := range slice {
		if s == val {
			return true
		}
	}
	return false
}
