// Package history computes the most recent past execution times
// for a given cron expression relative to a reference time.
//
// It is useful for auditing when a job last ran, debugging missed
// executions, and providing human-readable summaries of recent
// schedule activity.
//
// Usage:
//
//	res, err := history.Last("0 9 * * 1-5 /usr/bin/backup", time.Now(), 5)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, e := range res.Entries {
//		fmt.Println(e.Formatted)
//	}
//
// The n parameter must be between 1 and 100 inclusive.
package history
