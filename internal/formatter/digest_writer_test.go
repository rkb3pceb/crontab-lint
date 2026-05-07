package formatter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/nicholasgasior/crontab-lint/internal/digest"
	"github.com/nicholasgasior/crontab-lint/internal/formatter"
)

func makeDigest() digest.Result {
	return digest.Result{
		Expression: "0 * * * * job",
		Hash:       "abcdef1234567890",
		Fields:     []string{"0", "*", "*", "*", "*"},
		Signature:  "m=0 h=* dom=* mon=* dow=*",
	}
}

func TestWriteDigest_Text(t *testing.T) {
	f := formatter.New(false)
	var buf bytes.Buffer
	if err := f.WriteDigest(&buf, makeDigest()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	for _, want := range []string{"Hash", "abcdef1234567890", "Signature", "m=0"} {
		if !strings.Contains(out, want) {
			t.Errorf("text output missing %q:\n%s", want, out)
		}
	}
}

func TestWriteDigest_JSON(t *testing.T) {
	f := formatter.New(true)
	var buf bytes.Buffer
	if err := f.WriteDigest(&buf, makeDigest()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var got struct {
		Expression string   `json:"expression"`
		Hash       string   `json:"hash"`
		Fields     []string `json:"fields"`
		Signature  string   `json:"signature"`
	}
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if got.Hash != "abcdef1234567890" {
		t.Errorf("expected hash abcdef1234567890, got %s", got.Hash)
	}
	if len(got.Fields) != 5 {
		t.Errorf("expected 5 fields, got %d", len(got.Fields))
	}
	if got.Signature != "m=0 h=* dom=* mon=* dow=*" {
		t.Errorf("unexpected signature: %s", got.Signature)
	}
}
