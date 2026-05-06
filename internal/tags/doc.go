// Package tags provides semantic tagging for crontab expressions.
//
// Tags are short, human-readable labels that classify the scheduling
// behavior of a crontab entry. They are derived purely from the
// structure of the five schedule fields and require no runtime execution.
//
// Available tags:
//
//   - high-frequency  — runs every minute ("* * * * *")
//   - hourly          — runs at the top of every hour ("0 * * * *")
//   - daily           — runs at midnight daily ("0 0 * * *")
//   - interval        — uses step syntax ("*/N")
//   - month-restricted — month field is not a wildcard
//   - weekday-specific — day-of-week field is constrained
//   - dom-specific    — day-of-month field is constrained
//   - multi-time      — uses list syntax in minute or hour field
//   - custom          — none of the above patterns matched
//
// Example:
//
//	tags, err := tags.Extract("0 9 * * 1-5 /usr/bin/report")
//	// tags: [{weekday-specific Restricted to specific weekdays}]
package tags
