package matcher_test

import (
	"testing"
	"time"

	"github.com/user/crontab-lint/internal/matcher"
)

func ts(year, month, day, hour, minute int) time.Time {
	return time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC)
}

func TestMatch_EveryMinute(t *testing.T) {
	result, err := matcher.Match("* * * * * echo hi", ts(2024, 6, 15, 10, 30))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Matches {
		t.Error("expected match for wildcard expression")
	}
}

func TestMatch_SpecificTime(t *testing.T) {
	result, err := matcher.Match("30 10 15 6 * echo hi", ts(2024, 6, 15, 10, 30))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Matches {
		t.Error("expected match for exact time expression")
	}
}

func TestMatch_SpecificTime_NoMatch(t *testing.T) {
	result, err := matcher.Match("0 9 * * * echo hi", ts(2024, 6, 15, 10, 30))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Matches {
		t.Error("expected no match")
	}
	if result.Fields["hour"] {
		t.Error("hour field should not match")
	}
}

func TestMatch_StepExpression(t *testing.T) {
	result, err := matcher.Match("*/15 * * * * echo hi", ts(2024, 6, 15, 10, 30))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Matches {
		t.Error("expected match for */15 at minute 30")
	}
}

func TestMatch_RangeExpression(t *testing.T) {
	result, err := matcher.Match("0 9-17 * * 1-5 echo hi", ts(2024, 6, 17, 12, 0))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Matches {
		t.Error("expected match for business hours on weekday")
	}
}

func TestMatch_Alias(t *testing.T) {
	result, err := matcher.Match("@hourly echo hi", ts(2024, 6, 15, 10, 0))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !result.Matches {
		t.Error("expected @hourly to match at minute 0")
	}
}

func TestMatch_InvalidExpression(t *testing.T) {
	_, err := matcher.Match("invalid", ts(2024, 6, 15, 10, 0))
	if err == nil {
		t.Error("expected error for invalid expression")
	}
}

func TestMatch_FieldsPopulated(t *testing.T) {
	result, err := matcher.Match("* * * * * echo hi", ts(2024, 6, 15, 10, 30))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, key := range []string{"minute", "hour", "dayOfMonth", "month", "dayOfWeek"} {
		if !result.Fields[key] {
			t.Errorf("expected field %q to match", key)
		}
	}
}
