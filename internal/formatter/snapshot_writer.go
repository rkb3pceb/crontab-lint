package formatter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/user/crontab-lint/internal/snapshot"
)

// WriteSnapshot writes a snapshot and its diffs (if any) to w in the
// configured output format.
func (f *Formatter) WriteSnapshot(w io.Writer, snap snapshot.Snapshot, diffs []snapshot.Diff) error {
	if f.JSON {
		return writeSnapshotJSON(w, snap, diffs)
	}
	return writeSnapshotText(w, snap, diffs)
}

func writeSnapshotText(w io.Writer, snap snapshot.Snapshot, diffs []snapshot.Diff) error {
	if !snap.Valid {
		_, err := fmt.Fprintf(w, "Snapshot invalid: %s\n", snap.ParseError)
		return err
	}

	fmt.Fprintf(w, "Snapshot captured at: %s\n", snap.CapturedAt.Format("2006-01-02 15:04:05 UTC"))
	fmt.Fprintf(w, "Expression : %s\n", snap.Expression)
	fmt.Fprintf(w, "Normalized : %s\n", snap.Normalized)

	fields := []string{"minute", "hour", "dom", "month", "dow"}
	fmt.Fprintln(w, "Fields:")
	for _, f := range fields {
		fmt.Fprintf(w, "  %-8s %s\n", f+":", snap.Fields[f])
	}

	fmt.Fprintf(w, "Runs/day   : %d\n", snap.Stats.RunsPerDay)

	if len(diffs) > 0 {
		fmt.Fprintln(w, "Changes:")
		for _, d := range diffs {
			if d.Changed {
				fmt.Fprintf(w, "  %-8s %s -> %s\n", d.Field+":", d.Before, d.After)
			}
		}
	}
	return nil
}

func writeSnapshotJSON(w io.Writer, snap snapshot.Snapshot, diffs []snapshot.Diff) error {
	payload := struct {
		Snapshot snapshot.Snapshot `json:"snapshot"`
		Diffs    []snapshot.Diff   `json:"diffs,omitempty"`
	}{
		Snapshot: snap,
		Diffs:    diffs,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}
