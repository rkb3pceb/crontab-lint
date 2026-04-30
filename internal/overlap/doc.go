// Package overlap provides collision detection for crontab expressions.
//
// Given two cron expressions and a time window, Detect identifies all moments
// within that window where both jobs are scheduled to fire simultaneously.
// This is useful for detecting resource contention or unintended concurrency
// between scheduled tasks.
//
// Example usage:
//
//	res, err := overlap.Detect(
//		"0 * * * * /usr/bin/backup",
//		"*/30 * * * * /usr/bin/report",
//		time.Now(),
//		24*time.Hour,
//		10,
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("Found %d collision(s)\n", len(res.Collisions))
package overlap
