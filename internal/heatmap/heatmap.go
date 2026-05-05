// Package heatmap generates a 24x7 execution heatmap for a cron expression,
// showing how many times a job runs in each hour-of-day × day-of-week cell.
package heatmap

import (
	"fmt"

	"github.com/nicholasgasior/crontab-lint/internal/parser"
)

// DayLabels are the short names for days of the week (0 = Sunday).
var DayLabels = [7]string{"Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"}

// Cell holds the execution count for a single hour/day combination.
type Cell struct {
	Hour int
	Day  int // 0 = Sunday … 6 = Saturday
	Hits int
}

// Map is a 7-day × 24-hour grid of hit counts.
type Map struct {
	Expression string
	// Grid[day][hour] = number of minutes in that hour the job fires.
	Grid [7][24]int
	MaxHits int
}

// Build computes the heatmap for the given cron expression.
// It iterates over every minute of a representative week and counts matches.
func Build(expr string) (*Map, error) {
	fields, err := parser.Parse(expr)
	if err != nil {
		return nil, fmt.Errorf("heatmap: %w", err)
	}

	hm := &Map{Expression: expr}

	// Iterate all day-of-week (0-6), hours (0-23), minutes (0-59).
	for dow := 0; dow < 7; dow++ {
		for hour := 0; hour < 24; hour++ {
			for min := 0; min < 60; min++ {
				if matchesMinute(fields, min, hour, 1, 1, dow) {
					hm.Grid[dow][hour]++
				}
			}
			if hm.Grid[dow][hour] > hm.MaxHits {
				hm.MaxHits = hm.Grid[dow][hour]
			}
		}
	}

	return hm, nil
}

// Cells returns a flat slice of all non-zero cells.
func (m *Map) Cells() []Cell {
	var out []Cell
	for dow := 0; dow < 7; dow++ {
		for hour := 0; hour < 24; hour++ {
			if m.Grid[dow][hour] > 0 {
				out = append(out, Cell{Hour: hour, Day: dow, Hits: m.Grid[dow][hour]})
			}
		}
	}
	return out
}
