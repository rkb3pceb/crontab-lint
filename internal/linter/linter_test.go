package linter_test

import (
	"testing"

	"github.com/yourorg/crontab-lint/internal/linter"
)

func TestLint_ValidExpression(t *testing.T) {
	r := linter.Lint("0 9 * * 1 /usr/bin/backup.sh")
	if !r.Valid {
		t.Errorf("expected valid, got diagnostics: %v", r.Diagnostics)
	}
	if r.Schedule == "" {
		t.Error("expected non-empty schedule description")
	}
}

func TestLint_InvalidExpression(t *testing.T) {
	r := linter.Lint("99 * * * * /bin/true")
	if r.Valid {
		t.Error("expected invalid result for out-of-range minute")
	}
	if len(r.Diagnostics) == 0 {
		t.Error("expected at least one diagnostic")
	}
	if r.Diagnostics[0].Severity != linter.SeverityError {
		t.Errorf("expected error severity, got %s", r.Diagnostics[0].Severity)
	}
}

func TestLint_EveryMinuteWarning(t *testing.T) {
	r := linter.Lint("* * * * * /bin/heartbeat")
	if !r.Valid {
		t.Errorf("expected valid, got: %v", r.Diagnostics)
	}
	found := false
	for _, d := range r.Diagnostics {
		if d.Severity == linter.SeverityWarning && d.Field == "minute/hour" {
			found = true
		}
	}
	if !found {
		t.Error("expected warning about every-minute schedule")
	}
}

func TestLint_DomAndDowBothSet(t *testing.T) {
	r := linter.Lint("0 12 15 * 5 /usr/bin/report")
	if !r.Valid {
		t.Errorf("expected valid, got: %v", r.Diagnostics)
	}
	found := false
	for _, d := range r.Diagnostics {
		if d.Severity == linter.SeverityInfo && d.Field == "dom/dow" {
			found = true
		}
	}
	if !found {
		t.Error("expected info diagnostic about dom/dow OR semantics")
	}
}

func TestLint_TooFewFields(t *testing.T) {
	r := linter.Lint("0 9 * *")
	if r.Valid {
		t.Error("expected invalid result for too few fields")
	}
}
