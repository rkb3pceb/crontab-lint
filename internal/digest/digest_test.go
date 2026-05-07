package digest_test

import (
	"strings"
	"testing"

	"github.com/nicholasgasior/crontab-lint/internal/digest"
)

func TestCompute_BasicExpression(t *testing.T) {
	r, err := digest.Compute("*/5 * * * * echo hello")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r.Hash == "" {
		t.Error("expected non-empty hash")
	}
	if len(r.Hash) != 16 {
		t.Errorf("expected 16-char hash, got %d: %s", len(r.Hash), r.Hash)
	}
	if len(r.Fields) != 5 {
		t.Errorf("expected 5 fields, got %d", len(r.Fields))
	}
}

func TestCompute_SignatureFormat(t *testing.T) {
	r, err := digest.Compute("0 9 * * 1 run")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(r.Signature, "m=") {
		t.Errorf("signature should start with m=, got: %s", r.Signature)
	}
	for _, key := range []string{"m=", "h=", "dom=", "mon=", "dow="} {
		if !strings.Contains(r.Signature, key) {
			t.Errorf("signature missing key %q: %s", key, r.Signature)
		}
	}
}

func TestCompute_DeterministicHash(t *testing.T) {
	expr := "0 12 * * * job"
	r1, _ := digest.Compute(expr)
	r2, _ := digest.Compute(expr)
	if r1.Hash != r2.Hash {
		t.Errorf("hash not deterministic: %s != %s", r1.Hash, r2.Hash)
	}
}

func TestCompute_DifferentExpressionsProduceDifferentHashes(t *testing.T) {
	r1, _ := digest.Compute("0 * * * * job")
	r2, _ := digest.Compute("30 * * * * job")
	if r1.Hash == r2.Hash {
		t.Error("expected different hashes for different expressions")
	}
}

func TestCompute_AliasNormalized(t *testing.T) {
	r1, err := digest.Compute("@hourly echo x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r2, err := digest.Compute("0 * * * * echo x")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if r1.Hash != r2.Hash {
		t.Errorf("alias and canonical form should produce same hash: %s != %s", r1.Hash, r2.Hash)
	}
}

func TestCompute_InvalidExpression(t *testing.T) {
	_, err := digest.Compute("not-a-cron")
	if err == nil {
		t.Error("expected error for invalid expression")
	}
}
