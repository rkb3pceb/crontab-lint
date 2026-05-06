package snapshot_test

import (
	"testing"
	"time"

	"github.com/user/crontab-lint/internal/snapshot"
)

func TestTake_ValidExpression(t *testing.T) {
	snap := snapshot.Take("0 * * * * echo hello")
	if !snap.Valid {
		t.Fatalf("expected valid snapshot, got error: %s", snap.ParseError)
	}
	if snap.Expression != "0 * * * * echo hello" {
		t.Errorf("unexpected expression: %s", snap.Expression)
	}
	if snap.Fields["minute"] != "0" {
		t.Errorf("expected minute=0, got %s", snap.Fields["minute"])
	}
	if snap.Fields["hour"] != "*" {
		t.Errorf("expected hour=*, got %s", snap.Fields["hour"])
	}
	if snap.CapturedAt.IsZero() {
		t.Error("expected CapturedAt to be set")
	}
	if snap.CapturedAt.Location() != time.UTC {
		t.Error("expected CapturedAt to be UTC")
	}
}

func TestTake_InvalidExpression(t *testing.T) {
	snap := snapshot.Take("not-a-cron")
	if snap.Valid {
		t.Fatal("expected invalid snapshot")
	}
	if snap.ParseError == "" {
		t.Error("expected ParseError to be set")
	}
}

func TestTake_Alias(t *testing.T) {
	snap := snapshot.Take("@hourly echo hi")
	if !snap.Valid {
		t.Fatalf("expected valid snapshot, got: %s", snap.ParseError)
	}
	if snap.Normalized == "@hourly echo hi" {
		t.Error("expected normalized form to differ from alias")
	}
}

func TestTake_StatsPopulated(t *testing.T) {
	snap := snapshot.Take("* * * * * echo hi")
	if !snap.Valid {
		t.Fatalf("unexpected error: %s", snap.ParseError)
	}
	if snap.Stats.RunsPerDay == 0 {
		t.Error("expected Stats.RunsPerDay to be populated")
	}
}

func TestCompare_NoChange(t *testing.T) {
	a := snapshot.Take("0 9 * * 1 cmd")
	b := snapshot.Take("0 9 * * 1 cmd")
	diffs, err := snapshot.Compare(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, d := range diffs {
		if d.Changed {
			t.Errorf("expected no change in field %s", d.Field)
		}
	}
}

func TestCompare_FieldChanged(t *testing.T) {
	a := snapshot.Take("0 9 * * 1 cmd")
	b := snapshot.Take("30 10 * * 1 cmd")
	diffs, err := snapshot.Compare(a, b)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	changed := map[string]bool{}
	for _, d := range diffs {
		if d.Changed {
			changed[d.Field] = true
		}
	}
	if !changed["minute"] {
		t.Error("expected minute to be changed")
	}
	if !changed["hour"] {
		t.Error("expected hour to be changed")
	}
}

func TestCompare_InvalidBefore(t *testing.T) {
	a := snapshot.Take("bad")
	b := snapshot.Take("0 * * * * cmd")
	_, err := snapshot.Compare(a, b)
	if err == nil {
		t.Error("expected error for invalid before snapshot")
	}
}
