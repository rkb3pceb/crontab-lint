// Package grouper clusters multiple crontab expressions by their schedule
// pattern, helping identify duplicate or equivalent schedules across a set.
package grouper

import (
	"sort"

	"github.com/user/crontab-lint/internal/normalizer"
)

// Group represents a cluster of crontab expressions that share the same
// normalised schedule (i.e. they fire at exactly the same times).
type Group struct {
	// Key is the canonical normalised schedule used as the cluster key.
	Key string
	// Expressions contains every original expression belonging to this group.
	Expressions []string
}

// Cluster takes a slice of raw crontab expressions and returns groups of
// expressions that share an identical normalised schedule. Expressions that
// cannot be normalised are placed individually in a group whose Key equals
// the original expression.
//
// The returned slice is sorted by group size (largest first), then
// alphabetically by key so that output is deterministic.
func Cluster(expressions []string) []Group {
	index := make(map[string]*Group)
	order := []string{}

	for _, raw := range expressions {
		key := normaliseKey(raw)
		if _, exists := index[key]; !exists {
			index[key] = &Group{Key: key}
			order = append(order, key)
		}
		index[key].Expressions = append(index[key].Expressions, raw)
	}

	result := make([]Group, 0, len(index))
	for _, key := range order {
		result = append(result, *index[key])
	}

	sort.SliceStable(result, func(i, j int) bool {
		if len(result[i].Expressions) != len(result[j].Expressions) {
			return len(result[i].Expressions) > len(result[j].Expressions)
		}
		return result[i].Key < result[j].Key
	})

	return result
}

// normaliseKey returns the normalised schedule portion of expr, or the
// original string when normalisation fails.
func normaliseKey(expr string) string {
	norm, err := normalizer.Normalize(expr + " /bin/true")
	if err != nil {
		return expr
	}
	// Drop the command portion – keep only the five schedule fields.
	fields := splitFields(norm)
	if len(fields) < 5 {
		return expr
	}
	return fields[0] + " " + fields[1] + " " + fields[2] + " " + fields[3] + " " + fields[4]
}

func splitFields(s string) []string {
	var fields []string
	start := 0
	inField := false
	for i := 0; i < len(s); i++ {
		if s[i] == ' ' || s[i] == '\t' {
			if inField {
				fields = append(fields, s[start:i])
				inField = false
			}
		} else {
			if !inField {
				start = i
				inField = true
			}
		}
	}
	if inField {
		fields = append(fields, s[start:])
	}
	return fields
}
