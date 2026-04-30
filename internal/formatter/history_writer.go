package formatter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/user/crontab-lint/internal/history"
)

// WriteHistory writes history.Result to w in the requested format.
func (f *Formatter) WriteHistory(w io.Writer, res *history.Result) error {
	switch f.format {
	case FormatJSON:
		return writeHistoryJSON(w, res)
	default:
		return writeHistoryText(w, res)
	}
}

func writeHistoryText(w io.Writer, res *history.Result) error {
	fmt.Fprintf(w, "Recent executions for: %s\n", res.Expression)
	if len(res.Entries) == 0 {
		fmt.Fprintln(w, "  (no executions found in the search window)")
		return nil
	}
	for i, e := range res.Entries {
		fmt.Fprintf(w, "  %2d. %s\n", i+1, e.Formatted)
	}
	return nil
}

func writeHistoryJSON(w io.Writer, res *history.Result) error {
	type jsonEntry struct {
		Formatted string `json:"formatted"`
		Unix      int64  `json:"unix"`
	}
	type jsonResult struct {
		Expression string       `json:"expression"`
		Entries    []jsonEntry  `json:"entries"`
	}
	out := jsonResult{
		Expression: res.Expression,
		Entries:    make([]jsonEntry, len(res.Entries)),
	}
	for i, e := range res.Entries {
		out.Entries[i] = jsonEntry{
			Formatted: e.Formatted,
			Unix:      e.Time.Unix(),
		}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
