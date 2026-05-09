package formatter_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/user/crontab-lint/internal/formatter"
	"github.com/user/crontab-lint/internal/grouper"
)

func makeGroups() []grouper.Group {
	return []grouper.Group{
		{Key: "*/5 * * * *", Expressions: []string{"*/5 * * * *", "*/5 * * * *"}},
		{Key: "0 0 * * *", Expressions: []string{"0 0 * * *"}},
	}
}

func TestWriteGroups_Text(t *testing.T) {
	f := formatter.New("text")
	var buf bytes.Buffer
	if err := f.WriteGroups(&buf, makeGroups()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "*/5 * * * *") {
		t.Errorf("expected schedule key in output")
	}
	if !strings.Contains(out, "[2 duplicates]") {
		t.Errorf("expected duplicate label for group of 2")
	}
	if !strings.Contains(out, "0 0 * * *") {
		t.Errorf("expected second schedule in output")
	}
}

func TestWriteGroups_Text_Empty(t *testing.T) {
	f := formatter.New("text")
	var buf bytes.Buffer
	if err := f.WriteGroups(&buf, []grouper.Group{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(buf.String(), "No expressions") {
		t.Errorf("expected empty-state message")
	}
}

func TestWriteGroups_JSON(t *testing.T) {
	f := formatter.New("json")
	var buf bytes.Buffer
	if err := f.WriteGroups(&buf, makeGroups()); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	var out []struct {
		Key         string   `json:"key"`
		Count       int      `json:"count"`
		Expressions []string `json:"expressions"`
	}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if len(out) != 2 {
		t.Fatalf("expected 2 groups in JSON, got %d", len(out))
	}
	if out[0].Count != 2 {
		t.Errorf("expected count 2 for first group, got %d", out[0].Count)
	}
	if out[1].Count != 1 {
		t.Errorf("expected count 1 for second group, got %d", out[1].Count)
	}
}
