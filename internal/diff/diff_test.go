package diff_test

import (
	"testing"

	"github.com/user/crontab-lint/internal/diff"
)

func TestCompare_NoChange(t *testing.T) {
	r, err := diff.Compare("0 9 * * 1 echo hi", "0 9 * * 1 echo hi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Changed {
		t.Error("expected no change")
	}
	if len(r.Diffs) != 0 {
		t.Errorf("expected 0 diffs, got %d", len(r.Diffs))
	}
}

func TestCompare_SingleFieldChanged(t *testing.T) {
	r, err := diff.Compare("0 9 * * 1 echo hi", "30 9 * * 1 echo hi")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !r.Changed {
		t.Fatal("expected change")
	}
	if len(r.Diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(r.Diffs))
	}
	d := r.Diffs[0]
	if d.Field != "minute" {
		t.Errorf("expected field 'minute', got %q", d.Field)
	}
	if d.From != "0" || d.To != "30" {
		t.Errorf("unexpected from/to: %q -> %q", d.From, d.To)
	}
}

func TestCompare_MultipleFieldsChanged(t *testing.T) {
	r, err := diff.Compare("0 9 * * 1 cmd", "0 18 * * 5 cmd")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Diffs) != 2 {
		t.Fatalf("expected 2 diffs, got %d", len(r.Diffs))
	}
}

func TestCompare_InvalidFrom(t *testing.T) {
	_, err := diff.Compare("bad", "0 9 * * 1 cmd")
	if err == nil {
		t.Error("expected error for invalid 'from' expression")
	}
}

func TestCompare_InvalidTo(t *testing.T) {
	_, err := diff.Compare("0 9 * * 1 cmd", "bad")
	if err == nil {
		t.Error("expected error for invalid 'to' expression")
	}
}

func TestCompare_CommandChange_NoScheduleDiff(t *testing.T) {
	r, err := diff.Compare("0 9 * * 1 old-cmd", "0 9 * * 1 new-cmd")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Changed {
		t.Error("expected no schedule change when only command differs")
	}
}
