package scheduler_test

import (
	"testing"
	"time"

	"github.com/user/crontab-lint/internal/scheduler"
)

var base = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestNext_EveryMinute(t *testing.T) {
	next, err := scheduler.Next("* * * * * /bin/true", base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := base.Add(time.Minute)
	if !next.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, next)
	}
}

func TestNext_SpecificMinute(t *testing.T) {
	// "30 * * * *" — fires at minute 30 of every hour.
	next, err := scheduler.Next("30 * * * * /bin/true", base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := time.Date(2024, 1, 1, 0, 30, 0, 0, time.UTC)
	if !next.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, next)
	}
}

func TestNextN_ReturnsN(t *testing.T) {
	times, err := scheduler.NextN("0 * * * * /bin/true", base, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(times) != 3 {
		t.Fatalf("expected 3 times, got %d", len(times))
	}
	expected := []time.Time{
		time.Date(2024, 1, 1, 1, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 2, 0, 0, 0, time.UTC),
		time.Date(2024, 1, 1, 3, 0, 0, 0, time.UTC),
	}
	for i, e := range expected {
		if !times[i].Equal(e) {
			t.Errorf("times[%d]: expected %v, got %v", i, e, times[i])
		}
	}
}

func TestNextN_InvalidN(t *testing.T) {
	_, err := scheduler.NextN("* * * * * /bin/true", base, 0)
	if err == nil {
		t.Fatal("expected error for n=0, got nil")
	}
}

func TestNext_InvalidExpression(t *testing.T) {
	_, err := scheduler.Next("invalid", base)
	if err == nil {
		t.Fatal("expected error for invalid expression, got nil")
	}
}

func TestNext_StepExpression(t *testing.T) {
	// "*/15 * * * *" fires at 0, 15, 30, 45 minutes.
	next, err := scheduler.Next("*/15 * * * * /bin/true", base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := time.Date(2024, 1, 1, 0, 15, 0, 0, time.UTC)
	if !next.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, next)
	}
}

func TestNext_MonthConstrained(t *testing.T) {
	// "0 9 1 6 * /bin/true" — 09:00 on June 1st.
	next, err := scheduler.Next("0 9 1 6 * /bin/true", base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := time.Date(2024, 6, 1, 9, 0, 0, 0, time.UTC)
	if !next.Equal(expected) {
		t.Errorf("expected %v, got %v", expected, next)
	}
}
