package fingerprint_test

import (
	"strings"
	"testing"

	"github.com/user/crontab-lint/internal/fingerprint"
)

func TestCompute_EveryMinute(t *testing.T) {
	res, err := fingerprint.Compute("* * * * * /bin/job")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.Shape != "wildcard/wildcard/wildcard/wildcard/wildcard" {
		t.Errorf("unexpected shape: %s", res.Shape)
	}
}

func TestCompute_StepExpression(t *testing.T) {
	res, err := fingerprint.Compute("*/5 * * * * /bin/job")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.FieldShapes[0] != "step" {
		t.Errorf("expected step for minute, got %s", res.FieldShapes[0])
	}
}

func TestCompute_RangeExpression(t *testing.T) {
	res, err := fingerprint.Compute("0 9-17 * * 1-5 /bin/job")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.FieldShapes[1] != "range" {
		t.Errorf("expected range for hour, got %s", res.FieldShapes[1])
	}
	if res.FieldShapes[4] != "range" {
		t.Errorf("expected range for dow, got %s", res.FieldShapes[4])
	}
}

func TestCompute_ListExpression(t *testing.T) {
	res, err := fingerprint.Compute("0 8,12,18 * * * /bin/job")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.FieldShapes[1] != "list" {
		t.Errorf("expected list for hour, got %s", res.FieldShapes[1])
	}
}

func TestCompute_LiteralExpression(t *testing.T) {
	res, err := fingerprint.Compute("0 6 15 3 * /bin/job")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := "literal/literal/literal/literal/wildcard"
	if res.Shape != expected {
		t.Errorf("expected %s, got %s", expected, res.Shape)
	}
}

func TestCompute_AliasNormalized(t *testing.T) {
	a, err := fingerprint.Compute("@daily /bin/job")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, err := fingerprint.Compute("0 0 * * * /bin/job")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a.Shape != b.Shape {
		t.Errorf("alias and expanded should share shape: %s vs %s", a.Shape, b.Shape)
	}
}

func TestCompute_InvalidExpression(t *testing.T) {
	_, err := fingerprint.Compute("not a cron")
	if err == nil {
		t.Fatal("expected error for invalid expression")
	}
	if !strings.Contains(err.Error(), "fingerprint") {
		t.Errorf("error should mention fingerprint package: %v", err)
	}
}

func TestCompute_FieldShapesLength(t *testing.T) {
	res, err := fingerprint.Compute("*/10 */2 * * * /bin/job")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.FieldShapes) != 5 {
		t.Errorf("expected 5 field shapes, got %d", len(res.FieldShapes))
	}
}
