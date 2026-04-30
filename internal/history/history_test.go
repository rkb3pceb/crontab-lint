package history_test

import (
	"testing"
	"time"

	"github.com/user/crontab-lint/internal/history"
)

func ref() time.Time {
	// 2024-03-15 12:30:00 UTC, a Friday
	return time.Date(2024, 3, 15, 12, 30, 0, 0, time.UTC)
}

func TestLast_EveryMinute(t *testing.T) {
	res, err := history.Last("* * * * * echo hi", ref(), 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 5 {
		t.Errorf("expected 5 entries, got %d", len(res.Entries))
	}
	// All entries must be before ref.
	for _, e := range res.Entries {
		if !e.Time.Before(ref()) {
			t.Errorf("entry %v is not before ref %v", e.Time, ref())
		}
	}
}

func TestLast_HourlyJob(t *testing.T) {
	res, err := history.Last("0 * * * * echo hi", ref(), 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Entries) != 3 {
		t.Errorf("expected 3 entries, got %d", len(res.Entries))
	}
	// Each entry should be at minute 0.
	for _, e := range res.Entries {
		if e.Time.Minute() != 0 {
			t.Errorf("expected minute 0, got %d", e.Time.Minute())
		}
	}
}

func TestLast_InvalidN_Zero(t *testing.T) {
	_, err := history.Last("* * * * * echo hi", ref(), 0)
	if err == nil {
		t.Error("expected error for n=0")
	}
}

func TestLast_InvalidN_TooLarge(t *testing.T) {
	_, err := history.Last("* * * * * echo hi", ref(), 101)
	if err == nil {
		t.Error("expected error for n=101")
	}
}

func TestLast_InvalidExpression(t *testing.T) {
	_, err := history.Last("99 * * * * echo hi", ref(), 3)
	if err == nil {
		t.Error("expected error for invalid expression")
	}
}

func TestLast_FormattedNotEmpty(t *testing.T) {
	res, err := history.Last("0 12 * * * echo hi", ref(), 2)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, e := range res.Entries {
		if e.Formatted == "" {
			t.Error("expected non-empty Formatted string")
		}
	}
}
