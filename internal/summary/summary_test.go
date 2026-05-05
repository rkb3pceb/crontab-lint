package summary_test

import (
	"strings"
	"testing"
	"time"

	"github.com/yourorg/crontab-lint/internal/summary"
)

var ref = time.Date(2024, 1, 15, 12, 0, 0, 0, time.UTC)

func TestSummarize_ValidExpression(t *testing.T) {
	r := summary.Summarize("0 * * * * /bin/task", ref, 3)
	if r.Error != "" {
		t.Fatalf("unexpected error: %s", r.Error)
	}
	if r.Description == "" {
		t.Error("expected non-empty description")
	}
	if len(r.NextRuns) != 3 {
		t.Errorf("expected 3 next runs, got %d", len(r.NextRuns))
	}
	if r.Complexity.Grade == "" {
		t.Error("expected non-empty complexity grade")
	}
}

func TestSummarize_InvalidExpression(t *testing.T) {
	r := summary.Summarize("invalid", ref, 3)
	if r.Error == "" {
		t.Error("expected error for invalid expression")
	}
	if len(r.NextRuns) != 0 {
		t.Error("expected no next runs for invalid expression")
	}
}

func TestSummarize_NextCountCapped(t *testing.T) {
	r := summary.Summarize("* * * * * /bin/task", ref, 99)
	if r.Error != "" {
		t.Fatalf("unexpected error: %s", r.Error)
	}
	if len(r.NextRuns) > 5 {
		t.Errorf("next runs should be capped at 5, got %d", len(r.NextRuns))
	}
}

func TestSummarize_NextCountMinimum(t *testing.T) {
	r := summary.Summarize("0 9 * * 1 /bin/task", ref, 0)
	if r.Error != "" {
		t.Fatalf("unexpected error: %s", r.Error)
	}
	if len(r.NextRuns) < 1 {
		t.Error("expected at least 1 next run")
	}
}

func TestFormat_ValidResult(t *testing.T) {
	r := summary.Summarize("0 9 * * 1 /bin/task", ref, 2)
	out := summary.Format(r)
	for _, want := range []string{"Expression", "Description", "Complexity", "Next runs"} {
		if !strings.Contains(out, want) {
			t.Errorf("Format output missing %q", want)
		}
	}
}

func TestFormat_ErrorResult(t *testing.T) {
	r := summary.Summarize("bad expr", ref, 1)
	out := summary.Format(r)
	if !strings.Contains(out, "Error") {
		t.Error("Format output missing 'Error' for invalid expression")
	}
	if strings.Contains(out, "Description") {
		t.Error("Format should not include Description on error")
	}
}
