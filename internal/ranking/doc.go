// Package ranking provides quality-based ranking for crontab expressions.
//
// It combines three signals to produce a single integer score (0–100) and
// a letter grade (A–F) for each expression:
//
//   - Complexity: expressions with many list/range/step tokens are penalised.
//   - Frequency: jobs that run very often receive a frequency penalty.
//   - Lint severity: warnings reduce the score; errors result in grade F.
//
// Usage:
//
//	entries := ranking.Rank([]string{
//		"0 0 * * * /bin/daily",
//		"* * * * * /bin/every-minute",
//	})
//	for _, e := range entries {
//		fmt.Printf("%s  score=%d  grade=%s\n", e.Expression, e.Score, e.Grade)
//	}
//
// Entries are returned sorted best-first (highest score first).
package ranking
