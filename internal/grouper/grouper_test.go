package grouper_test

import (
	"testing"

	"github.com/user/crontab-lint/internal/grouper"
)

func TestCluster_SingleExpression(t *testing.T) {
	groups := grouper.Cluster([]string{"* * * * *"})
	if len(groups) != 1 {
		t.Fatalf("expected 1 group, got %d", len(groups))
	}
	if len(groups[0].Expressions) != 1 {
		t.Errorf("expected 1 expression in group, got %d", len(groups[0].Expressions))
	}
}

func TestCluster_DuplicatesGrouped(t *testing.T) {
	exprs := []string{"0 * * * *", "0 * * * *", "0 * * * *"}
	groups := grouper.Cluster(exprs)
	if len(groups) != 1 {
		t.Fatalf("expected 1 group for identical expressions, got %d", len(groups))
	}
	if len(groups[0].Expressions) != 3 {
		t.Errorf("expected 3 expressions, got %d", len(groups[0].Expressions))
	}
}

func TestCluster_DistinctSchedulesSeparated(t *testing.T) {
	exprs := []string{"0 * * * *", "0 0 * * *", "*/5 * * * *"}
	groups := grouper.Cluster(exprs)
	if len(groups) != 3 {
		t.Fatalf("expected 3 groups, got %d", len(groups))
	}
}

func TestCluster_AliasNormalisedToSameGroup(t *testing.T) {
	// @hourly and "0 * * * *" should resolve to the same normalised key.
	exprs := []string{"@hourly", "0 * * * *"}
	groups := grouper.Cluster(exprs)
	if len(groups) != 1 {
		t.Fatalf("expected aliases to collapse into 1 group, got %d groups", len(groups))
	}
	if len(groups[0].Expressions) != 2 {
		t.Errorf("expected 2 expressions, got %d", len(groups[0].Expressions))
	}
}

func TestCluster_LargestGroupFirst(t *testing.T) {
	exprs := []string{
		"0 0 * * *",
		"*/5 * * * *", "*/5 * * * *", "*/5 * * * *",
		"0 12 * * 1",
	}
	groups := grouper.Cluster(exprs)
	if len(groups[0].Expressions) < len(groups[len(groups)-1].Expressions) {
		t.Errorf("groups not sorted by size descending")
	}
}

func TestCluster_EmptyInput(t *testing.T) {
	groups := grouper.Cluster([]string{})
	if len(groups) != 0 {
		t.Errorf("expected 0 groups for empty input, got %d", len(groups))
	}
}

func TestCluster_InvalidExpressionIsolated(t *testing.T) {
	exprs := []string{"not-a-cron", "* * * * *"}
	groups := grouper.Cluster(exprs)
	// Each should form its own group since the invalid one cannot be normalised.
	if len(groups) != 2 {
		t.Fatalf("expected 2 groups, got %d", len(groups))
	}
}
