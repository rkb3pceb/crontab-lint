// Package digest produces a compact fingerprint of a crontab expression,
// combining a stable hash with a human-readable summary of its key properties.
package digest

import (
	"crypto/sha256"
	"fmt"
	"strings"

	"github.com/nicholasgasior/crontab-lint/internal/normalizer"
	"github.com/nicholasgasior/crontab-lint/internal/parser"
)

// Result holds the computed digest for a crontab expression.
type Result struct {
	// Expression is the normalized form of the input.
	Expression string
	// Hash is a short hex fingerprint (first 8 bytes of SHA-256).
	Hash string
	// Fields contains the five cron fields in order.
	Fields []string
	// Signature is a compact human-readable key, e.g. "m=*/5 h=* dom=* mon=* dow=*".
	Signature string
}

// Compute derives a Digest Result from the given raw crontab expression.
// It returns an error if the expression cannot be parsed.
func Compute(expr string) (Result, error) {
	norm, err := normalizer.Normalize(expr)
	if err != nil {
		return Result{}, fmt.Errorf("digest: normalize: %w", err)
	}

	entry, err := parser.Parse(norm)
	if err != nil {
		return Result{}, fmt.Errorf("digest: parse: %w", err)
	}

	fields := []string{
		entry.Minute,
		entry.Hour,
		entry.DayOfMonth,
		entry.Month,
		entry.DayOfWeek,
	}

	schedPart := strings.Join(fields, " ")
	hash := sha256.Sum256([]byte(schedPart))
	hashHex := fmt.Sprintf("%x", hash[:8])

	sig := fmt.Sprintf("m=%s h=%s dom=%s mon=%s dow=%s",
		entry.Minute, entry.Hour, entry.DayOfMonth, entry.Month, entry.DayOfWeek)

	return Result{
		Expression: norm,
		Hash:       hashHex,
		Fields:     fields,
		Signature:  sig,
	}, nil
}
