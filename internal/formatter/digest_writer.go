package formatter

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/nicholasgasior/crontab-lint/internal/digest"
)

// WriteDigest writes a digest.Result to w in either text or JSON format
// depending on whether the Formatter was created with JSON mode enabled.
func (f *Formatter) WriteDigest(w io.Writer, r digest.Result) error {
	if f.json {
		return writeDigestJSON(w, r)
	}
	return writeDigestText(w, r)
}

func writeDigestText(w io.Writer, r digest.Result) error {
	_, err := fmt.Fprintf(w,
		"Expression : %s\nHash       : %s\nSignature  : %s\n",
		r.Expression, r.Hash, r.Signature,
	)
	return err
}

func writeDigestJSON(w io.Writer, r digest.Result) error {
	payload := struct {
		Expression string   `json:"expression"`
		Hash       string   `json:"hash"`
		Fields     []string `json:"fields"`
		Signature  string   `json:"signature"`
	}{
		Expression: r.Expression,
		Hash:       r.Hash,
		Fields:     r.Fields,
		Signature:  r.Signature,
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}
