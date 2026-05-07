// Package fingerprint produces a structural fingerprint for crontab expressions.
//
// A fingerprint describes the shape of each cron field rather than its exact
// value. This makes it possible to group expressions by their structural
// pattern — for example, identifying all jobs that run on a fixed minute with
// a step-based hour — without regard to the specific numbers used.
//
// Shape tokens:
//
//	"wildcard" — field is "*"
//	"step"     — field uses "/" (e.g. "*/5", "0/15")
//	"range"    — field uses "-" (e.g. "9-17")
//	"list"     — field uses "," (e.g. "1,3,5")
//	"literal"  — field is a single numeric value (e.g. "0", "12")
//
// Expressions are normalized before fingerprinting so that aliases such as
// @daily and name tokens such as MON are expanded to their canonical numeric
// form, ensuring consistent shape output.
//
// Example:
//
//	res, err := fingerprint.Compute("*/5 0 * * MON /bin/job")
//	// res.Shape == "step/literal/wildcard/wildcard/literal"
package fingerprint
