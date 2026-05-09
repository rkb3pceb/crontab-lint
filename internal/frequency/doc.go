// Package frequency estimates how often a cron expression fires and exposes
// the result as a structured [Result] value.
//
// # Overview
//
// Given a standard five-field cron expression (minute, hour, day-of-month,
// month, day-of-week followed by a command), [Estimate] calculates:
//
//   - [Result.PerHour]  — average fires per hour
//   - [Result.PerDay]   — average fires per 24-hour day
//   - [Result.PerWeek]  — average fires per 7-day week
//   - [Result.Label]    — a concise human-readable summary
//
// # Usage
//
//	r, err := frequency.Estimate("*/15 * * * * /usr/bin/backup")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(r.Label)   // "multiple times per hour"
//	fmt.Println(r.PerDay)  // 96
//
// # Limitations
//
// The month field is not factored into PerDay / PerWeek estimates; the
// calculation assumes the job runs every month. This keeps the API simple
// while remaining accurate for the vast majority of crontab entries.
package frequency
