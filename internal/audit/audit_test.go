package audit_test

import (
	"testing"
	"time"

	"github.com/dkarter/crontab-lint/internal/audit"
)

func TestBuild_ValidExpression(t *testing.T) {
	rec := audit.Build("0 9 * * 1-5 /usr/bin/backup")

	if !rec.Valid {
		t.Errorf("expected valid=true, got false")
	}
	if rec.Expression != "0 9 * * 1-5 /usr/bin/backup" {
		t.Errorf("unexpected expression: %s", rec.Expression)
	}
	if rec.AnalyzedAt.IsZero() {
		t.Error("expected AnalyzedAt to be set")
	}
	if time.Since(rec.AnalyzedAt) > 5*time.Second {
		t.Error("AnalyzedAt is too far in the past")
	}
	if len(rec.Events) == 0 {
		t.Error("expected at least one event")
	}
}

func TestBuild_ParseError(t *testing.T) {
	rec := audit.Build("invalid")

	if rec.Valid {
		t.Error("expected valid=false for unparseable expression")
	}
	if len(rec.Events) == 0 {
		t.Fatal("expected at least one event")
	}

	found := false
	for _, e := range rec.Events {
		if e.Severity == audit.SeverityError && e.Code == "PARSE_ERROR" {
			found = true
		}
	}
	if !found {
		t.Error("expected a PARSE_ERROR event")
	}
}

func TestBuild_EveryMinuteWarning(t *testing.T) {
	rec := audit.Build("* * * * * /bin/poll")

	hasWarning := false
	for _, e := range rec.Events {
		if e.Severity == audit.SeverityWarning {
			hasWarning = true
			break
		}
	}
	if !hasWarning {
		t.Error("expected at least one warning event for every-minute expression")
	}
}

func TestBuild_EventTimestamps(t *testing.T) {
	before := time.Now().UTC()
	rec := audit.Build("0 * * * * /bin/task")
	after := time.Now().UTC()

	for _, e := range rec.Events {
		if e.Timestamp.Before(before) || e.Timestamp.After(after) {
			t.Errorf("event timestamp %v out of expected range [%v, %v]", e.Timestamp, before, after)
		}
	}
}

func TestBuild_SeverityInfo_OnSuccess(t *testing.T) {
	rec := audit.Build("30 6 * * 1 /usr/bin/report")

	hasInfo := false
	for _, e := range rec.Events {
		if e.Severity == audit.SeverityInfo && e.Code == "PARSE_OK" {
			hasInfo = true
			break
		}
	}
	if !hasInfo {
		t.Error("expected a PARSE_OK info event for valid expression")
	}
}
