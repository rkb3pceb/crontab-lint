package explainer_test

import (
	"strings"
	"testing"

	"github.com/user/crontab-lint/internal/explainer"
)

func TestExplain_ValidExpression(t *testing.T) {
	expr := "0 9 * * 1 /usr/bin/backup"
	result, err := explainer.Explain(expr)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Expression != expr {
		t.Errorf("expected expression %q, got %q", expr, result.Expression)
	}
	if len(result.Fields) != 5 {
		t.Errorf("expected 5 field explanations, got %d", len(result.Fields))
	}
	if result.Summary == "" {
		t.Error("expected non-empty summary")
	}
}

func TestExplain_FieldLabels(t *testing.T) {
	result, err := explainer.Explain("30 6 * * * /bin/run")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expectedLabels := []string{"Minute", "Hour", "Day of Month", "Month", "Day of Week"}
	for i, fe := range result.Fields {
		if fe.Field != expectedLabels[i] {
			t.Errorf("field %d: expected label %q, got %q", i, expectedLabels[i], fe.Field)
		}
	}
}

func TestExplain_FieldValues(t *testing.T) {
	result, err := explainer.Explain("*/15 * * * * /bin/check")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Fields[0].Value != "*/15" {
		t.Errorf("expected minute value '*/15', got %q", result.Fields[0].Value)
	}
	if !strings.Contains(result.Fields[0].Description, "15") {
		t.Errorf("expected description to mention '15', got %q", result.Fields[0].Description)
	}
}

func TestExplain_InvalidExpression(t *testing.T) {
	_, err := explainer.Explain("not a cron")
	if err == nil {
		t.Error("expected error for invalid expression, got nil")
	}
}

func TestExplain_WildcardDescription(t *testing.T) {
	result, err := explainer.Explain("* * * * * /bin/run")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for _, fe := range result.Fields {
		if fe.Description == "" {
			t.Errorf("field %q has empty description", fe.Field)
		}
	}
}
