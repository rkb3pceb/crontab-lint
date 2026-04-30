package formatter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/user/crontab-lint/internal/formatter"
	"github.com/user/crontab-lint/internal/history"
)

func makeHistory() *history.Result {
	t1 := time.Date(2024, 3, 15, 11, 0, 0, 0, time.UTC)
	t2 := time.Date(2024, 3, 15, 10, 0, 0, 0, time.UTC)
	return &history.Result{
		Expression: "0 * * * * echo hi",
		Entries: []history.Entry{
			{Time: t2, Formatted: t2.Format("2006-01-02 15:04:05 MST")},
			{Time: t1, Formatted: t1.Format("2006-01-02 15:04:05 MST")},
		},
	}
}

func TestWriteHistory_Text(t *testing.T) {
	f := formatter.New(formatter.FormatText)
	var buf bytes.Buffer
	if err := f.WriteHistory(&buf, makeHistory()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "0 * * * * echo hi") {
		t.Error("expected expression in output")
	}
	if !strings.Contains(out, "2024-03-15") {
		t.Error("expected date in output")
	}
	if !strings.Contains(out, "1.") {
		t.Error("expected numbered list in output")
	}
}

func TestWriteHistory_JSON(t *testing.T) {
	f := formatter.New(formatter.FormatJSON)
	var buf bytes.Buffer
	if err := f.WriteHistory(&buf, makeHistory()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if out["expression"] != "0 * * * * echo hi" {
		t.Errorf("unexpected expression: %v", out["expression"])
	}
	entries, ok := out["entries"].([]interface{})
	if !ok || len(entries) != 2 {
		t.Errorf("expected 2 entries, got %v", out["entries"])
	}
}

func TestWriteHistory_Text_Empty(t *testing.T) {
	f := formatter.New(formatter.FormatText)
	var buf bytes.Buffer
	res := &history.Result{Expression: "0 0 31 2 * echo hi", Entries: nil}
	if err := f.WriteHistory(&buf, res); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "no executions") {
		t.Error("expected 'no executions' message for empty result")
	}
}
