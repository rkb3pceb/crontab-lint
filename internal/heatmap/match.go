package heatmap

import (
	"strconv"
	"strings"
)

// matchesMinute returns true when the given time values satisfy all five cron
// fields contained in fields (minute, hour, dom, month, dow).
func matchesMinute(fields []string, minute, hour, dom, month, dow int) bool {
	if len(fields) < 5 {
		return false
	}
	return fieldMatch(fields[0], minute, 0, 59) &&
		fieldMatch(fields[1], hour, 0, 23) &&
		fieldMatch(fields[2], dom, 1, 31) &&
		fieldMatch(fields[3], month, 1, 12) &&
		fieldMatch(fields[4], dow, 0, 6)
}

func fieldMatch(field string, value, min, max int) bool {
	if field == "*" {
		return true
	}
	for _, part := range strings.Split(field, ",") {
		if matchPart(part, value, min, max) {
			return true
		}
	}
	return false
}

func matchPart(part string, value, min, max int) bool {
	// Step: */n or a-b/n
	if idx := strings.Index(part, "/"); idx != -1 {
		step, err := strconv.Atoi(part[idx+1:])
		if err != nil || step <= 0 {
			return false
		}
		base := min
		rangeEnd := max
		prefix := part[:idx]
		if prefix != "*" {
			if dashIdx := strings.Index(prefix, "-"); dashIdx != -1 {
				base, _ = strconv.Atoi(prefix[:dashIdx])
				rangeEnd, _ = strconv.Atoi(prefix[dashIdx+1:])
			} else {
				base, _ = strconv.Atoi(prefix)
			}
		}
		for v := base; v <= rangeEnd; v += step {
			if v == value {
				return true
			}
		}
		return false
	}
	// Range: a-b
	if idx := strings.Index(part, "-"); idx != -1 {
		lo, _ := strconv.Atoi(part[:idx])
		hi, _ := strconv.Atoi(part[idx+1:])
		return value >= lo && value <= hi
	}
	// Literal
	n, err := strconv.Atoi(part)
	return err == nil && n == value
}
