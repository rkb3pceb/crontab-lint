package formatter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/user/crontab-lint/internal/grouper"
)

// WriteGroups writes a slice of grouper.Group values to w using the format
// configured on f ("text" or "json").
func (f *Formatter) WriteGroups(w io.Writer, groups []grouper.Group) error {
	if f.format == "json" {
		return writeGroupsJSON(w, groups)
	}
	return writeGroupsText(w, groups)
}

func writeGroupsText(w io.Writer, groups []grouper.Group) error {
	if len(groups) == 0 {
		_, err := fmt.Fprintln(w, "No expressions to group.")
		return err
	}
	for _, g := range groups {
		dupeLabel := ""
		if len(g.Expressions) > 1 {
			dupeLabel = fmt.Sprintf(" [%d duplicates]", len(g.Expressions))
		}
		if _, err := fmt.Fprintf(w, "Schedule: %s%s\n", g.Key, dupeLabel); err != nil {
			return err
		}
		for _, expr := range g.Expressions {
			if _, err := fmt.Fprintf(w, "  %s\n", expr); err != nil {
				return err
			}
		}
	}
	return nil
}

func writeGroupsJSON(w io.Writer, groups []grouper.Group) error {
	type jsonGroup struct {
		Key         string   `json:"key"`
		Count       int      `json:"count"`
		Expressions []string `json:"expressions"`
	}
	out := make([]jsonGroup, len(groups))
	for i, g := range groups {
		out[i] = jsonGroup{
			Key:         g.Key,
			Count:       len(g.Expressions),
			Expressions: g.Expressions,
		}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
