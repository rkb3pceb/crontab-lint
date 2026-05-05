// Package heatmap builds a 7-day × 24-hour execution frequency grid for a
// cron expression.
//
// For each cell (day-of-week, hour) the grid records how many minutes within
// that hour the job would fire.  This lets callers render ASCII heat-maps,
// JSON payloads, or SVG visualisations showing the busiest windows at a glance.
//
// # Usage
//
//	hm, err := heatmap.Build("*/5 * * * * /usr/bin/backup")
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, cell := range hm.Cells() {
//		fmt.Printf("%s %02d:xx → %d runs\n",
//			heatmap.DayLabels[cell.Day], cell.Hour, cell.Hits)
//	}
//
// The returned [Map] also exposes the raw [Map.Grid] array and [Map.MaxHits]
// for normalisation when drawing colour scales.
package heatmap
