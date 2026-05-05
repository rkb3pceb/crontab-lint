// Package summary produces a consolidated, human-readable overview of a
// crontab expression by combining outputs from the humanizer, complexity,
// and scheduler packages into a single Result value.
//
// Typical usage:
//
//	import "github.com/yourorg/crontab-lint/internal/summary"
//
//	r := summary.Summarize("0 9 * * 1 /backup.sh", time.Now(), 3)
//	fmt.Print(summary.Format(r))
//
// Output example:
//
//	Expression : 0 9 * * 1 /backup.sh
//	Description: At 09:00, only on Monday
//	Complexity : A (score 2)
//	Next runs  :
//	  - 2024-01-22T09:00:00Z
//	  - 2024-01-29T09:00:00Z
//	  - 2024-02-05T09:00:00Z
package summary
