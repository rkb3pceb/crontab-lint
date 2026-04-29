// Package scheduler computes future execution times for crontab expressions.
//
// It builds on top of the parser package to validate and interpret cron
// fields, then walks forward in time minute-by-minute to find matching
// timestamps.
//
// # Usage
//
//	next, err := scheduler.Next("0 9 * * 1 /bin/backup", time.Now())
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("Next run:", next)
//
//	times, err := scheduler.NextN("*/5 * * * * /bin/check", time.Now(), 5)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, t := range times {
//		fmt.Println(t)
//	}
//
// The scheduler respects all five standard cron fields (minute, hour,
// day-of-month, month, day-of-week) and supports wildcards, ranges,
// lists, and step values as parsed by the parser package.
package scheduler
