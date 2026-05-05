package similarity_test

import (
	"testing"

	"github.com/nicholasgasior/crontab-lint/internal/similarity"
)

func TestCompare_IdenticalExpressions(t *testing.T) {
	r, err := similarity.Compare("0 9 * * 1", "0 9 * * 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Score != 1.0 {
		t.Errorf("expected score 1.0, got %.2f", r.Score)
	}
}

func TestCompare_CompletelyDifferent(t *testing.T) {
	r, err := similarity.Compare("0 9 * * 1", "30 18 15 6 5")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Score >= 0.5 {
		t.Errorf("expected low score, got %.2f", r.Score)
	}
}

func TestCompare_WildcardVsSpecific(t *testing.T) {
	r, err := similarity.Compare("* * * * *", "0 9 * * 1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Minute and hour differ (wildcard vs specific → 0.5 each),
	// dom and month are both wildcard → 1.0 each, dow differs → 0.5.
	// Expected: (0.5+0.5+1.0+1.0+0.5)/5 = 0.7
	if r.Score < 0.6 || r.Score > 0.8 {
		t.Errorf("expected score near 0.7, got %.2f", r.Score)
	}
}

func TestCompare_PartialListOverlap(t *testing.T) {
	r, err := similarity.Compare("0,30 * * * *", "0,15,30 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// minute: intersection={0,30}=2, union={0,15,30}=3 → 2/3 ≈ 0.667
	minScore := r.FieldScores[0]
	if minScore < 0.6 || minScore > 0.75 {
		t.Errorf("expected minute field score ~0.667, got %.3f", minScore)
	}
}

func TestCompare_InvalidExprA(t *testing.T) {
	_, err := similarity.Compare("invalid", "0 9 * * 1")
	if err == nil {
		t.Error("expected error for invalid expression A")
	}
}

func TestCompare_InvalidExprB(t *testing.T) {
	_, err := similarity.Compare("0 9 * * 1", "not valid")
	if err == nil {
		t.Error("expected error for invalid expression B")
	}
}

func TestCompare_FieldScoresLength(t *testing.T) {
	r, err := similarity.Compare("*/5 * * * *", "*/10 * * * *")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.FieldScores) != 5 {
		t.Errorf("expected 5 field scores, got %d", len(r.FieldScores))
	}
}

func TestFieldName(t *testing.T) {
	expected := []string{"minute", "hour", "dom", "month", "dow"}
	for i, name := range expected {
		if got := similarity.FieldName(i); got != name {
			t.Errorf("FieldName(%d): expected %q, got %q", i, name, got)
		}
	}
	if got := similarity.FieldName(99); got != "unknown" {
		t.Errorf("expected 'unknown' for out-of-range index, got %q", got)
	}
}
