// Package timeline builds a chronological list of scheduled execution times
// for a cron expression within a caller-specified time window.
//
// # Usage
//
//	result, err := timeline.Build("0 9 * * 1-5 run-job", time.Now(), 7*24*time.Hour)
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, e := range result.Entries {
//		fmt.Println(e.Label)
//	}
//	if result.Truncated {
//		fmt.Println("... output truncated")
//	}
//
// # Limits
//
// At most 200 entries are returned per call. If more executions fall within
// the requested window the Result.Truncated flag is set to true so callers
// can inform the user that the listing is incomplete.
package timeline
