// Package explainer provides field-by-field breakdown of crontab expressions
// into structured, human-readable explanations.
package explainer

import (
	"fmt"
	"strings"

	"github.com/user/crontab-lint/internal/humanizer"
	"github.com/user/crontab-lint/internal/parser"
)

// FieldExplanation holds the label and description for a single crontab field.
type FieldExplanation struct {
	Field       string `json:"field"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

// Explanation is the full breakdown of a crontab expression.
type Explanation struct {
	Expression string             `json:"expression"`
	Summary    string             `json:"summary"`
	Fields     []FieldExplanation `json:"fields"`
}

var fieldLabels = []string{
	"Minute",
	"Hour",
	"Day of Month",
	"Month",
	"Day of Week",
}

// Explain parses a crontab expression and returns a structured field-by-field
// explanation. Returns an error if the expression cannot be parsed.
func Explain(expr string) (*Explanation, error) {
	entry, err := parser.Parse(expr)
	if err != nil {
		return nil, fmt.Errorf("cannot explain invalid expression: %w", err)
	}

	fields := strings.Fields(expr)
	scheduleFields := fields[:5]

	var explanations []FieldExplanation
	for i, val := range scheduleFields {
		desc := humanizer.Describe(val, i)
		explanations = append(explanations, FieldExplanation{
			Field:       fieldLabels[i],
			Value:       val,
			Description: desc,
		})
	}

	summary := humanizer.DescribeSchedule(entry)

	return &Explanation{
		Expression: expr,
		Summary:    summary,
		Fields:     explanations,
	}, nil
}
