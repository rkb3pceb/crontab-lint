package formatter_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/example/crontab-lint/internal/formatter"
	"github.com/example/crontab-lint/internal/linter"
)

func makeResult(expr string, valid bool, human string, errs, warns []string) linter.Result {
	return linter.Result{
		Expression: expr,
		Valid:      valid,
		Human:      human,
		Errors:     errs,
		Warnings:   warns,
	}
}

func TestWriteText_ValidResult(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.New(formatter.FormatText, &buf)
	results := []linter.Result{
		makeResult("*/5 * * * * cmd", true, "every 5 minutes", nil, nil),
	}
	if err := f.Write(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[OK]") {
		t.Errorf("expected [OK] in output, got: %s", out)
	}
	if !strings.Contains(out, "every 5 minutes") {
		t.Errorf("expected human description in output, got: %s", out)
	}
}

func TestWriteText_ErrorResult(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.New(formatter.FormatText, &buf)
	results := []linter.Result{
		makeResult("99 * * * * cmd", false, "", []string{"minute out of range"}, nil),
	}
	if err := f.Write(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[ERROR]") {
		t.Errorf("expected [ERROR] in output, got: %s", out)
	}
	if !strings.Contains(out, "minute out of range") {
		t.Errorf("expected error message in output, got: %s", out)
	}
}

func TestWriteText_WarnResult(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.New(formatter.FormatText, &buf)
	results := []linter.Result{
		makeResult("* * * * * cmd", true, "every minute", nil, []string{"runs every minute"}),
	}
	if err := f.Write(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "[WARN]") {
		t.Errorf("expected [WARN] in output, got: %s", out)
	}
}

func TestWriteJSON_ValidResult(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.New(formatter.FormatJSON, &buf)
	results := []linter.Result{
		makeResult("0 * * * * cmd", true, "at minute 0", nil, nil),
	}
	if err := f.Write(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, `"valid":true`) {
		t.Errorf("expected valid:true in JSON, got: %s", out)
	}
	if !strings.Contains(out, `"errors":[]`) {
		t.Errorf("expected empty errors array, got: %s", out)
	}
}

func TestWriteJSON_MultipleResults(t *testing.T) {
	var buf bytes.Buffer
	f := formatter.New(formatter.FormatJSON, &buf)
	results := []linter.Result{
		makeResult("0 * * * * cmd", true, "at minute 0", nil, nil),
		makeResult("bad", false, "", []string{"too few fields"}, nil),
	}
	if err := f.Write(results); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if strings.Count(out, "expression") != 2 {
		t.Errorf("expected 2 expression entries in JSON, got: %s", out)
	}
}
