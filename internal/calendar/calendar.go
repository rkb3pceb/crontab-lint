// Package calendar generates a visual calendar view of cron job execution
// times for a given day or week, showing when a cron expression will fire.
package calendar

import (
	"fmt"
	"strings"
	"time"

	"github.com/example/crontab-lint/internal/matcher"
	"github.com/example/crontab-lint/internal/parser"
)

// HourBlock represents a single hour slot in the calendar view.
type HourBlock struct {
	Hour    int
	Minutes []int
	Fired   bool
}

// DayView holds the calendar view for a single day.
type DayView struct {
	Date   time.Time
	Blocks []HourBlock
	Total  int
}

// BuildDay returns a DayView showing which hours/minutes a cron expression
// fires on the given date.
func BuildDay(expr string, date time.Time) (*DayView, error) {
	if _, err := parser.Parse(expr); err != nil {
		return nil, fmt.Errorf("invalid expression: %w", err)
	}

	day := date.Truncate(24 * time.Hour)
	view := &DayView{
		Date:   day,
		Blocks: make([]HourBlock, 24),
	}

	for h := 0; h < 24; h++ {
		block := HourBlock{Hour: h}
		for m := 0; m < 60; m++ {
			t := time.Date(day.Year(), day.Month(), day.Day(), h, m, 0, 0, day.Location())
			if matcher.Match(expr, t) {
				block.Minutes = append(block.Minutes, m)
				block.Fired = true
				view.Total++
			}
		}
		view.Blocks[h] = block
	}

	return view, nil
}

// FormatDay renders the DayView as a compact ASCII calendar string.
func FormatDay(view *DayView) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Calendar for %s — %d execution(s)\n",
		view.Date.Format("2006-01-02 (Monday)"), view.Total))
	sb.WriteString(strings.Repeat("-", 50) + "\n")

	for _, block := range view.Blocks {
		if !block.Fired {
			continue
		}
		minStrs := make([]string, len(block.Minutes))
		for i, m := range block.Minutes {
			minStrs[i] = fmt.Sprintf("%02d", m)
		}
		sb.WriteString(fmt.Sprintf("  %02d:00  →  [%s]\n", block.Hour, strings.Join(minStrs, ", ")))
	}

	return sb.String()
}
