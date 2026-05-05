package complexity_test

import (
	"testing"

	"github.com/example/crontab-lint/internal/complexity"
)

func TestScore_SimpleWildcard(t *testing.T) {
	r, err := complexity.Score("* * * * * /bin/run")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Score != 0 {
		t.Errorf("expected score 0, got %d", r.Score)
	}
	if r.Grade != "A" {
		t.Errorf("expected grade A, got %s", r.Grade)
	}
}

func TestScore_StepExpression(t *testing.T) {
	r, err := complexity.Score("*/5 * * * * /bin/run")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// */5 counts as non-wildcard (+1) and step (+2) = 3, rest are wildcards
	if r.Score != 3 {
		t.Errorf("expected score 3, got %d", r.Score)
	}
	if r.Grade != "B" {
		t.Errorf("expected grade B, got %s", r.Grade)
	}
}

func TestScore_RangeExpression(t *testing.T) {
	r, err := complexity.Score("0 9-17 * * * /bin/run")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// minute=0: +1; hour=9-17: +1(non-wildcard)+1(range)=+2 => total 3
	if r.Score != 3 {
		t.Errorf("expected score 3, got %d", r.Score)
	}
}

func TestScore_ListExpression(t *testing.T) {
	r, err := complexity.Score("0 8,12,18 * * * /bin/run")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// minute=0: +1; hour=8,12,18: +1(non-wildcard)+3(list) = 4 => total 5
	if r.Score != 5 {
		t.Errorf("expected score 5, got %d", r.Score)
	}
	if r.Grade != "B" {
		t.Errorf("expected grade B, got %s", r.Grade)
	}
}

func TestScore_ComplexExpression(t *testing.T) {
	r, err := complexity.Score("*/10 8,12,18 1-15 1,6,12 1-5 /bin/run")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Score <= 9 {
		t.Errorf("expected high score (>9), got %d", r.Score)
	}
	if r.Grade != "D" && r.Grade != "F" {
		t.Errorf("expected grade D or F for complex expression, got %s", r.Grade)
	}
}

func TestScore_InvalidExpression(t *testing.T) {
	_, err := complexity.Score("bad expression")
	if err == nil {
		t.Fatal("expected error for invalid expression, got nil")
	}
}

func TestScore_FactorsPopulated(t *testing.T) {
	r, err := complexity.Score("*/5 8,12 * * * /bin/run")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Factors) == 0 {
		t.Error("expected factors to be populated for non-trivial expression")
	}
}
