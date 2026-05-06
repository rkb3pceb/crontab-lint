// Package audit provides structured audit trail generation for crontab
// expressions analyzed by crontab-lint.
//
// An audit Record captures the full analysis lifecycle of a single crontab
// expression: whether it parsed successfully, any linter warnings or errors
// that were raised, and the timestamp at which the analysis occurred.
//
// Each finding is represented as an Event with an associated Severity:
//
//   - SeverityInfo    — informational, e.g. successful parse
//   - SeverityWarning — non-fatal linter findings (e.g. every-minute jobs)
//   - SeverityError   — fatal issues that make the expression invalid
//
// Usage:
//
//	rec := audit.Build("*/5 * * * * /usr/bin/poll")
//	if !rec.Valid {
//	    log.Println("expression has errors")
//	}
//	for _, e := range rec.Events {
//	    fmt.Printf("[%s] %s\n", e.Severity, e.Message)
//	}
package audit
