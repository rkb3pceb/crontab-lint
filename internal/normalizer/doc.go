// Package normalizer provides canonicalization of crontab expressions.
//
// It handles two categories of normalization:
//
//  1. Alias expansion — shorthand strings such as @daily, @weekly, and
//     @hourly are expanded to their equivalent 5-field cron expressions.
//     The command portion of the entry, if present, is preserved.
//
//  2. Named-value substitution — month names (jan–dec) and day-of-week
//     names (sun–sat) are replaced with their numeric counterparts so
//     that downstream components (parser, linter, scheduler) only need
//     to handle numeric fields.
//
// Example usage:
//
//	expr := normalizer.Normalize("@daily /usr/bin/backup")
//	// expr == "0 0 * * * /usr/bin/backup"
//
//	expr = normalizer.Normalize("0 0 * jan mon")
//	// expr == "0 0 * 1 1"
//
// Normalize is idempotent: calling it on an already-normalized expression
// returns the same expression unchanged.
package normalizer
