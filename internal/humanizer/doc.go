// Package humanizer provides utilities for converting cron field expressions
// into human-readable descriptions.
//
// It supports standard cron syntax including wildcards (*), step values (*/n),
// ranges (a-b), and lists (a,b,c). Month and day-of-week fields are rendered
// using their natural language names where applicable.
//
// # Field Types
//
// The package defines the following field types for use with [Describe]:
//
//   - [Minute]     – the minute field (0–59)
//   - [Hour]       – the hour field (0–23)
//   - [DayOfMonth] – the day-of-month field (1–31)
//   - [Month]      – the month field (1–12 or JAN–DEC)
//   - [DayOfWeek]  – the day-of-week field (0–7 or SUN–SAT)
//
// # Example usage
//
//	import "crontab-lint/internal/humanizer"
//
//	// Describe a single field
//	fmt.Println(humanizer.Describe("*/15", humanizer.Minute))
//	// Output: every 15 minute(s)
//
//	// Describe a full schedule
//	fields := []string{"0", "8", "*", "*", "1-5"}
//	fmt.Println(humanizer.DescribeSchedule(fields))
//	// Output: at 0 minute; at 8 hour; every day of month; every month; from Monday to Friday
package humanizer
