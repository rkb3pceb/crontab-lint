package normalizer_test

import (
	"testing"

	"github.com/user/crontab-lint/internal/normalizer"
)

func TestNormalize_Alias(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"@yearly", "0 0 1 1 *"},
		{"@annually", "0 0 1 1 *"},
		{"@monthly", "0 0 1 * *"},
		{"@weekly", "0 0 * * 0"},
		{"@daily", "0 0 * * *"},
		{"@midnight", "0 0 * * *"},
		{"@hourly", "0 * * * *"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := normalizer.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q; want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNormalize_AliasWithCommand(t *testing.T) {
	input := "@daily /usr/bin/backup"
	want := "0 0 * * * /usr/bin/backup"
	got := normalizer.Normalize(input)
	if got != want {
		t.Errorf("Normalize(%q) = %q; want %q", input, got, want)
	}
}

func TestNormalize_MonthNames(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"0 0 1 jan *", "0 0 1 1 *"},
		{"0 0 1 dec *", "0 0 1 12 *"},
		{"0 0 1 jan-mar *", "0 0 1 1-3 *"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := normalizer.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q; want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNormalize_DowNames(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"0 0 * * sun", "0 0 * * 0"},
		{"0 0 * * mon-fri", "0 0 * * 1-5"},
		{"0 0 * * sat", "0 0 * * 6"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := normalizer.Normalize(tt.input)
			if got != tt.want {
				t.Errorf("Normalize(%q) = %q; want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestNormalize_TooFewFields(t *testing.T) {
	input := "* * *"
	got := normalizer.Normalize(input)
	if got != input {
		t.Errorf("Normalize(%q) = %q; want original %q", input, got, input)
	}
}

func TestNormalize_PlainExpression(t *testing.T) {
	input := "30 4 * * 1"
	got := normalizer.Normalize(input)
	if got != input {
		t.Errorf("Normalize(%q) = %q; want %q", input, got, input)
	}
}
