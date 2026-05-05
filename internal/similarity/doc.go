// Package similarity measures how alike two crontab expressions are.
//
// It parses both expressions using the shared parser, then scores each of the
// five schedule fields (minute, hour, day-of-month, month, day-of-week)
// independently using a Jaccard-style token overlap metric:
//
//	- Identical fields score 1.0.
//	- A wildcard (*) compared against a specific value scores 0.5, reflecting
//	  that the wildcard partially overlaps every possible value.
//	- Comma-separated lists are split into token sets and scored by
//	  intersection / union (Jaccard index).
//
// The overall Score is the arithmetic mean of the five field scores, giving a
// value in [0.0, 1.0].
//
// Typical usage:
//
//	r, err := similarity.Compare("0 9 * * 1", "0 10 * * 1")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Similarity: %.0f%%\n", r.Score*100)
package similarity
