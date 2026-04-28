package suggester_test

import (
	"testing"

	"github.com/yourorg/crontab-lint/internal/suggester"
)

func TestSuggest_EveryMinute(t *testing.T) {
	suggestions := suggester.Suggest("minute", "*", "runs every minute which may cause high load")
	if len(suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(suggestions))
	}
	if suggestions[0].Fix == "" {
		t.Error("expected a non-empty Fix string")
	}
}

func TestSuggest_HighFrequency(t *testing.T) {
	suggestions := suggester.Suggest("minute", "*/2", "high frequency interval detected")
	if len(suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(suggestions))
	}
	if suggestions[0].Field != "minute" {
		t.Errorf("expected field 'minute', got '%s'", suggestions[0].Field)
	}
}

func TestSuggest_DomAndDow(t *testing.T) {
	suggestions := suggester.Suggest("dom", "15", "both DOM and DOW are set")
	if len(suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(suggestions))
	}
	if suggestions[0].Field != "dom+dow" {
		t.Errorf("expected field 'dom+dow', got '%s'", suggestions[0].Field)
	}
}

func TestSuggest_Unreachable(t *testing.T) {
	suggestions := suggester.Suggest("dom", "31", "unreachable day-of-month in short months")
	if len(suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(suggestions))
	}
}

func TestSuggest_OutOfRange(t *testing.T) {
	suggestions := suggester.Suggest("hour", "25", "out of range value")
	if len(suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(suggestions))
	}
	expected := "Use a value between 0 and 23."
	if suggestions[0].Fix != expected {
		t.Errorf("expected fix %q, got %q", expected, suggestions[0].Fix)
	}
}

func TestSuggest_UnknownField(t *testing.T) {
	suggestions := suggester.Suggest("unknown", "99", "invalid value")
	if len(suggestions) != 1 {
		t.Fatalf("expected 1 suggestion, got %d", len(suggestions))
	}
	if suggestions[0].Fix == "" {
		t.Error("expected a fallback fix hint")
	}
}

func TestSuggest_NoMatch(t *testing.T) {
	suggestions := suggester.Suggest("minute", "5", "some unrelated message")
	if len(suggestions) != 0 {
		t.Errorf("expected 0 suggestions, got %d", len(suggestions))
	}
}
