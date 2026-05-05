// Package similarity computes a likeness score between two crontab expressions.
// It compares each schedule field individually and returns a score from 0.0
// (completely different) to 1.0 (identical schedules).
package similarity

import (
	"strings"

	"github.com/nicholasgasior/crontab-lint/internal/parser"
)

// Result holds the outcome of comparing two crontab expressions.
type Result struct {
	// Score is a value in [0.0, 1.0] indicating how similar the two expressions are.
	Score float64
	// FieldScores contains the per-field similarity score (minute, hour, dom, month, dow).
	FieldScores [5]float64
	// ExprA is the first expression that was compared.
	ExprA string
	// ExprB is the second expression that was compared.
	ExprB string
}

var fieldNames = [5]string{"minute", "hour", "dom", "month", "dow"}

// Compare returns a Result describing how similar exprA and exprB are.
// Both expressions must be valid 5-or-6 field crontab strings; if either
// fails to parse the returned score is 0.0 and an error is returned.
func Compare(exprA, exprB string) (Result, error) {
	result := Result{ExprA: exprA, ExprB: exprB}

	schedA, err := parser.Parse(exprA)
	if err != nil {
		return result, err
	}
	schedB, err := parser.Parse(exprB)
	if err != nil {
		return result, err
	}

	fields := [5][2]string{
		{schedA.Minute, schedB.Minute},
		{schedA.Hour, schedB.Hour},
		{schedA.DayOfMonth, schedB.DayOfMonth},
		{schedA.Month, schedB.Month},
		{schedA.DayOfWeek, schedB.DayOfWeek},
	}

	var total float64
	for i, pair := range fields {
		s := fieldSimilarity(pair[0], pair[1])
		result.FieldScores[i] = s
		total += s
	}
	result.Score = total / 5.0
	return result, nil
}

// FieldName returns the human-readable label for field index i (0–4).
func FieldName(i int) string {
	if i < 0 || i >= len(fieldNames) {
		return "unknown"
	}
	return fieldNames[i]
}

// fieldSimilarity returns a score in [0.0, 1.0] for two individual field values.
func fieldSimilarity(a, b string) float64 {
	if a == b {
		return 1.0
	}
	// Normalise to lower-case for comparison.
	a = strings.ToLower(strings.TrimSpace(a))
	b = strings.ToLower(strings.TrimSpace(b))
	if a == b {
		return 1.0
	}
	// Both wildcards — already handled above, but guard anyway.
	if a == "*" || b == "*" {
		return 0.5
	}
	// Shared tokens (comma-separated lists).
	aTokens := tokenSet(a)
	bTokens := tokenSet(b)
	intersection := 0
	for t := range aTokens {
		if bTokens[t] {
			intersection++
		}
	}
	union := len(aTokens) + len(bTokens) - intersection
	if union == 0 {
		return 1.0
	}
	return float64(intersection) / float64(union)
}

func tokenSet(field string) map[string]bool {
	set := make(map[string]bool)
	for _, t := range strings.Split(field, ",") {
		set[strings.TrimSpace(t)] = true
	}
	return set
}
