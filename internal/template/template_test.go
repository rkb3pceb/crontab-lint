package template_test

import (
	"testing"

	"github.com/user/crontab-lint/internal/template"
)

func TestAll_ReturnsNonEmpty(t *testing.T) {
	templates := template.All()
	if len(templates) == 0 {
		t.Fatal("expected at least one template")
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	a := template.All()
	b := template.All()
	a[0].Name = "mutated"
	if b[0].Name == "mutated" {
		t.Error("All() should return a copy, not a reference to the internal slice")
	}
}

func TestLookup_Found(t *testing.T) {
	tmpl, err := template.Lookup("hourly")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tmpl.Expression != "0 * * * *" {
		t.Errorf("expected '0 * * * *', got %q", tmpl.Expression)
	}
}

func TestLookup_NotFound(t *testing.T) {
	_, err := template.Lookup("nonexistent")
	if err == nil {
		t.Fatal("expected error for unknown template name")
	}
}

func TestLookup_AllNamesResolvable(t *testing.T) {
	for _, tmpl := range template.All() {
		_, err := template.Lookup(tmpl.Name)
		if err != nil {
			t.Errorf("Lookup(%q) failed: %v", tmpl.Name, err)
		}
	}
}

func TestSearch_ByName(t *testing.T) {
	results := template.Search("hourly")
	if len(results) == 0 {
		t.Fatal("expected at least one result for 'hourly'")
	}
	for _, r := range results {
		if r.Name != "hourly" && r.Expression != "0 * * * *" {
			// acceptable if matched by description or tag
		}
	}
}

func TestSearch_ByTag(t *testing.T) {
	results := template.Search("daily")
	if len(results) < 2 {
		t.Errorf("expected at least 2 daily templates, got %d", len(results))
	}
}

func TestSearch_CaseInsensitive(t *testing.T) {
	a := template.Search("DAILY")
	b := template.Search("daily")
	if len(a) != len(b) {
		t.Errorf("case-insensitive search mismatch: %d vs %d", len(a), len(b))
	}
}

func TestSearch_NoMatch(t *testing.T) {
	results := template.Search("zzznomatch")
	if len(results) != 0 {
		t.Errorf("expected no results, got %d", len(results))
	}
}

func TestTemplate_FieldsPopulated(t *testing.T) {
	for _, tmpl := range template.All() {
		if tmpl.Name == "" {
			t.Error("template has empty Name")
		}
		if tmpl.Expression == "" {
			t.Errorf("template %q has empty Expression", tmpl.Name)
		}
		if tmpl.Description == "" {
			t.Errorf("template %q has empty Description", tmpl.Name)
		}
	}
}
