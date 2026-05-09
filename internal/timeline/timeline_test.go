package timeline_test

import (
	"testing"
	"time"

	"github.com/dkgv/crontab-lint/internal/timeline"
)

func ref() time.Time {
	return time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
}

func TestBuild_EveryMinute(t *testing.T) {
	r, err := timeline.Build("* * * * * echo hi", ref(), time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Entries) != 60 {
		t.Errorf("expected 60 entries for hourly window, got %d", len(r.Entries))
	}
	if r.Truncated {
		t.Error("expected Truncated=false for 60 entries")
	}
}

func TestBuild_HourlyJob(t *testing.T) {
	r, err := timeline.Build("0 * * * * echo hi", ref(), 24*time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Entries) != 24 {
		t.Errorf("expected 24 entries for daily window, got %d", len(r.Entries))
	}
}

func TestBuild_TruncatesAtMax(t *testing.T) {
	// every minute for 4 hours = 240 > 200
	r, err := timeline.Build("* * * * * echo hi", ref(), 4*time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Entries) != 200 {
		t.Errorf("expected 200 entries (max), got %d", len(r.Entries))
	}
	if !r.Truncated {
		t.Error("expected Truncated=true")
	}
}

func TestBuild_InvalidExpression(t *testing.T) {
	_, err := timeline.Build("invalid", ref(), time.Hour)
	if err == nil {
		t.Error("expected error for invalid expression")
	}
}

func TestBuild_NegativeWindow(t *testing.T) {
	_, err := timeline.Build("* * * * * echo hi", ref(), -time.Hour)
	if err == nil {
		t.Error("expected error for negative window")
	}
}

func TestBuild_EntryLabel(t *testing.T) {
	r, err := timeline.Build("0 9 * * * echo hi", ref(), 24*time.Hour)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Entries) == 0 {
		t.Fatal("expected at least one entry")
	}
	if r.Entries[0].Label == "" {
		t.Error("expected non-empty label")
	}
}

func TestBuild_WindowAndExprStored(t *testing.T) {
	expr := "30 6 * * 1 echo hi"
	win := 7 * 24 * time.Hour
	r, err := timeline.Build(expr, ref(), win)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Expression != expr {
		t.Errorf("expression mismatch: got %q", r.Expression)
	}
	if r.Window != win {
		t.Errorf("window mismatch: got %v", r.Window)
	}
}
