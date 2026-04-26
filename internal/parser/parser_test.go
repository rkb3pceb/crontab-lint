package parser

import (
	"testing"
)

func TestParse_Valid(t *testing.T) {
	expr, err := Parse("5 4 * * 0 /usr/bin/backup")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if expr.Minute != "5" {
		t.Errorf("expected minute=5, got %s", expr.Minute)
	}
	if expr.Command != "/usr/bin/backup" {
		t.Errorf("expected command=/usr/bin/backup, got %s", expr.Command)
	}
}

func TestParse_TooFewFields(t *testing.T) {
	_, err := Parse("5 4 * *")
	if err == nil {
		t.Fatal("expected error for too few fields")
	}
}

func TestParse_CommandWithSpaces(t *testing.T) {
	expr, err := Parse("0 0 * * * /bin/sh -c 'echo hello'")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if expr.Command != "/bin/sh -c 'echo hello'" {
		t.Errorf("unexpected command: %s", expr.Command)
	}
}

func TestValidateField_Wildcard(t *testing.T) {
	f := Field{Name: "minute", Min: 0, Max: 59}
	if err := ValidateField("*", f); err != nil {
		t.Errorf("unexpected error for wildcard: %v", err)
	}
}

func TestValidateField_ValidNumber(t *testing.T) {
	f := Field{Name: "hour", Min: 0, Max: 23}
	if err := ValidateField("12", f); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateField_OutOfRange(t *testing.T) {
	f := Field{Name: "hour", Min: 0, Max: 23}
	if err := ValidateField("25", f); err == nil {
		t.Error("expected error for out-of-range value")
	}
}

func TestValidateField_ValidRange(t *testing.T) {
	f := Field{Name: "minute", Min: 0, Max: 59}
	if err := ValidateField("10-30", f); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateField_InvalidRange(t *testing.T) {
	f := Field{Name: "minute", Min: 0, Max: 59}
	if err := ValidateField("50-70", f); err == nil {
		t.Error("expected error for out-of-range bound")
	}
}

func TestValidateField_ValidList(t *testing.T) {
	f := Field{Name: "month", Min: 1, Max: 12}
	if err := ValidateField("1,6,12", f); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateField_ValidStep(t *testing.T) {
	f := Field{Name: "minute", Min: 0, Max: 59}
	if err := ValidateField("*/15", f); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestValidateField_InvalidStep(t *testing.T) {
	f := Field{Name: "minute", Min: 0, Max: 59}
	if err := ValidateField("*/0", f); err == nil {
		t.Error("expected error for step=0")
	}
}
