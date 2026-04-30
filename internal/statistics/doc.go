// Package statistics provides execution frequency analysis for crontab expressions.
//
// It estimates how many times a cron job will run per day, week, and month
// based on the cardinality of each schedule field.
//
// Usage:
//
//	stats, err := statistics.Compute("*/15 * * * * /usr/bin/backup")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Runs per day: %.0f\n", stats.RunsPerDay)
//	fmt.Printf("Frequency:    %s\n", stats.Frequency)
//
// The Frequency field returns a human-readable label such as:
//   - "every minute"
//   - "high frequency (multiple times per hour)"
//   - "hourly or more"
//   - "daily or more"
//   - "less than once per day"
package statistics
