// Package calendar provides a day-level calendar view of cron job execution
// times, showing which hours and minutes a given cron expression will fire
// within a 24-hour period.
//
// # Usage
//
// Build a view for a specific date:
//
//	view, err := calendar.BuildDay("30 9,17 * * * /usr/bin/backup", time.Now())
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Print(calendar.FormatDay(view))
//
// # Output Example
//
//	Calendar for 2024-06-05 (Wednesday) — 2 execution(s)
//	--------------------------------------------------
//	  09:00  →  [30]
//	  17:00  →  [30]
//
// The DayView struct exposes raw block data for programmatic use, while
// FormatDay renders a human-readable ASCII summary.
package calendar
