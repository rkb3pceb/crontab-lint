// Package complexity provides a scoring mechanism for crontab expression complexity.
// A higher score indicates a more complex or harder-to-read expression.
package complexity

import (
	"strings"

	"github.com/example/crontab-lint/internal/parser"
)

// Result holds the complexity score and a breakdown of contributing factors.
type Result struct {
	Score    int
	Factors  []string
	Grade    string
}

// Score computes a complexity result for the given crontab expression.
// Returns an error if the expression cannot be parsed.
func Score(expr string) (Result, error) {
	fields, err := parser.Parse(expr)
	if err != nil {
		return Result{}, err
	}

	// fields[0..4] are minute, hour, dom, month, dow
	scheduleFields := fields[:5]

	var factors []string
	total := 0

	for _, f := range scheduleFields {
		if strings.Contains(f, ",") {
			parts := strings.Split(f, ",")
			points := len(parts)
			total += points
			factors = append(factors, "list expression adds "+itoa(points)+" point(s)")
		}
		if strings.Contains(f, "/") {
			total += 2
			factors = append(factors, "step expression adds 2 points")
		}
		if strings.Contains(f, "-") {
			total += 1
			factors = append(factors, "range expression adds 1 point")
		}
		if f != "*" {
			total += 1
		}
	}

	return Result{
		Score:   total,
		Factors: factors,
		Grade:   grade(total),
	}, nil
}

func grade(score int) string {
	switch {
	case score <= 2:
		return "A"
	case score <= 5:
		return "B"
	case score <= 9:
		return "C"
	case score <= 14:
		return "D"
	default:
		return "F"
	}
}

func itoa(n int) string {
	const digits = "0123456789"
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = digits[n%10]
		n /= 10
	}
	return string(buf[pos:])
}
