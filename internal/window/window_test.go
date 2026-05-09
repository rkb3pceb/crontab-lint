package window_test

import (
	"testing"
	"time"

	"github.com/your-org/crontab-lint/internal/window"
)

func ref(s string) time.Time {
	t, _ := time.Parse("2006-01-02 15:04", s)
	return t
}

func TestBuild_HourlyJob(t *testing.T) {
	from := ref("2024-01-01 00:00")
	to := ref("2024-01-01 06:00")
	res, err := window.Build("0 * * * * /cmd", from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Count != 6 {
		t.Errorf("expected 6 firings, got %d", res.Count)
	}
}

func TestBuild_EveryMinute_Capped(t *testing.T) {
	from := ref("2024-01-01 00:00")
	to := from.Add(48 * time.Hour) // would produce >1000 results
	res, err := window.Build("* * * * * /cmd", from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Count != window.MaxResults {
		t.Errorf("expected results capped at %d, got %d", window.MaxResults, res.Count)
	}
}

func TestBuild_EmptyRange(t *testing.T) {
	from := ref("2024-01-01 12:00")
	to := ref("2024-01-01 12:00") // same time — no window
	_, err := window.Build("* * * * * /cmd", from, to)
	if err == nil {
		t.Fatal("expected error for zero-width window")
	}
}

func TestBuild_InvalidExpression(t *testing.T) {
	from := ref("2024-01-01 00:00")
	to := ref("2024-01-01 01:00")
	_, err := window.Build("invalid", from, to)
	if err == nil {
		t.Fatal("expected error for invalid expression")
	}
}

func TestBuild_SpecificTime_SingleFiring(t *testing.T) {
	from := ref("2024-03-15 08:00")
	to := ref("2024-03-15 09:00")
	res, err := window.Build("30 8 * * * /cmd", from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Count != 1 {
		t.Errorf("expected 1 firing, got %d", res.Count)
	}
	expected := ref("2024-03-15 08:30")
	if !res.Times[0].Equal(expected) {
		t.Errorf("expected firing at %v, got %v", expected, res.Times[0])
	}
}

func TestBuild_ResultFields(t *testing.T) {
	from := ref("2024-06-01 00:00")
	to := ref("2024-06-01 01:00")
	res, err := window.Build("0 * * * * /cmd", from, to)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Expression != "0 * * * * /cmd" {
		t.Errorf("expression not preserved")
	}
	if !res.From.Equal(from) || !res.To.Equal(to) {
		t.Errorf("from/to not preserved")
	}
}
