package overlap_test

import (
	"testing"
	"time"

	"github.com/user/crontab-lint/internal/overlap"
)

var base = time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)

func TestDetect_NoCollisions(t *testing.T) {
	// exprA fires at :00, exprB fires at :30 — never overlap within one hour
	res, err := overlap.Detect("0 * * * * echo a", "30 * * * * echo b", base, time.Hour, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Collisions) != 0 {
		t.Errorf("expected 0 collisions, got %d", len(res.Collisions))
	}
}

func TestDetect_AllCollisions(t *testing.T) {
	// Both fire every minute — every tick is a collision
	res, err := overlap.Detect("* * * * * echo a", "* * * * * echo b", base, 5*time.Minute, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Collisions) == 0 {
		t.Error("expected collisions, got none")
	}
}

func TestDetect_PartialCollision(t *testing.T) {
	// exprA every minute, exprB every 15 minutes — collide at :00 and :15
	res, err := overlap.Detect("* * * * * echo a", "*/15 * * * * echo b", base, time.Hour, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Collisions) != 4 {
		t.Errorf("expected 4 collisions, got %d", len(res.Collisions))
	}
}

func TestDetect_MaxHitsRespected(t *testing.T) {
	res, err := overlap.Detect("* * * * * echo a", "* * * * * echo b", base, time.Hour, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.Collisions) > 3 {
		t.Errorf("expected at most 3 collisions, got %d", len(res.Collisions))
	}
}

func TestDetect_InvalidExprA(t *testing.T) {
	_, err := overlap.Detect("invalid", "* * * * * echo b", base, time.Hour, 5)
	if err == nil {
		t.Error("expected error for invalid expression A")
	}
}

func TestDetect_InvalidExprB(t *testing.T) {
	_, err := overlap.Detect("* * * * * echo a", "invalid", base, time.Hour, 5)
	if err == nil {
		t.Error("expected error for invalid expression B")
	}
}

func TestDetect_ZeroMaxHits(t *testing.T) {
	_, err := overlap.Detect("* * * * * echo a", "* * * * * echo b", base, time.Hour, 0)
	if err == nil {
		t.Error("expected error for maxHits=0")
	}
}

func TestDetect_NegativeWindow(t *testing.T) {
	_, err := overlap.Detect("* * * * * echo a", "* * * * * echo b", base, -time.Hour, 5)
	if err == nil {
		t.Error("expected error for negative window")
	}
}
