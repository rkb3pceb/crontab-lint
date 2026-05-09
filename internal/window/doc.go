// Package window provides bounded time-range querying for cron expressions.
//
// Given a cron expression and a [from, to) time range, Build returns every
// moment within that range at which the expression would fire. Results are
// capped at MaxResults entries to guard against expressions that fire very
// frequently over long windows.
//
// Example:
//
//	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
//	to   := time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC)
//
//	res, err := window.Build("0 9 * * * /deploy.sh", from, to)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Printf("%d firings on %s\n", res.Count, from.Format("2006-01-02"))
package window
