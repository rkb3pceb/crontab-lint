package ranking_test

import (
	"testing"

	"github.com/example/crontab-lint/internal/ranking"
)

func TestRank_ReturnsSameCount(t *testing.T) {
	exprs := []string{"0 * * * * /bin/hourly", "* * * * * /bin/every", "0 0 * * * /bin/daily"}
	got := ranking.Rank(exprs)
	if len(got) != len(exprs) {
		t.Fatalf("expected %d entries, got %d", len(exprs), len(got))
	}
}

func TestRank_SortedBestFirst(t *testing.T) {
	exprs := []string{"* * * * * /bin/every", "0 0 * * * /bin/daily"}
	got := ranking.Rank(exprs)
	if got[0].Score < got[1].Score {
		t.Errorf("expected entries sorted best-first; got scores %d, %d", got[0].Score, got[1].Score)
	}
}

func TestRank_InvalidExpressionGradedF(t *testing.T) {
	exprs := []string{"not-a-cron"}
	got := ranking.Rank(exprs)
	if got[0].Grade != "F" {
		t.Errorf("expected grade F for invalid expression, got %s", got[0].Grade)
	}
	if got[0].Errors == 0 {
		t.Error("expected at least one error for invalid expression")
	}
}

func TestRank_DailyJobHigherThanEveryMinute(t *testing.T) {
	exprs := []string{"* * * * * /bin/every", "0 0 * * * /bin/daily"}
	got := ranking.Rank(exprs)
	// daily should rank higher (lower frequency penalty)
	if got[0].Expression != "0 0 * * * /bin/daily" {
		t.Errorf("expected daily job ranked first, got %q", got[0].Expression)
	}
}

func TestRank_GradeA_ForSimpleDaily(t *testing.T) {
	got := ranking.Rank([]string{"0 0 * * * /bin/daily"})
	if got[0].Grade == "F" {
		t.Errorf("expected non-F grade for simple daily job, got F")
	}
}

func TestRank_EmptyInput(t *testing.T) {
	got := ranking.Rank([]string{})
	if len(got) != 0 {
		t.Errorf("expected empty result for empty input")
	}
}

func TestRank_ScoreNonNegative(t *testing.T) {
	exprs := []string{"* * * * * /bin/every", "0/5 * * * * /bin/five", "0 0 1 1 * /bin/yearly"}
	for _, e := range ranking.Rank(exprs) {
		if e.Score < 0 {
			t.Errorf("score should not be negative, got %d for %q", e.Score, e.Expression)
		}
	}
}
