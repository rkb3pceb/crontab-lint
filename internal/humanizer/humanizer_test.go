package humanizer

import (
	"testing"
)

func TestDescribe_Wildcard(t *testing.T) {
	result := Describe("*", Minute)
	if result != "every minute" {
		t.Errorf("expected 'every minute', got %q", result)
	}
}

func TestDescribe_Step(t *testing.T) {
	result := Describe("*/5", Minute)
	if result != "every 5 minute(s)" {
		t.Errorf("expected 'every 5 minute(s)', got %q", result)
	}
}

func TestDescribe_Range(t *testing.T) {
	result := Describe("9-17", Hour)
	if result != "from 9 to 17" {
		t.Errorf("expected 'from 9 to 17', got %q", result)
	}
}

func TestDescribe_List(t *testing.T) {
	result := Describe("1,3,5", DayOfWeek)
	if result != "Monday, Wednesday, Friday" {
		t.Errorf("expected 'Monday, Wednesday, Friday', got %q", result)
	}
}

func TestDescribe_MonthName(t *testing.T) {
	result := Describe("12", Month)
	if result != "at December month" {
		t.Errorf("expected 'at December month', got %q", result)
	}
}

func TestDescribe_DayOfWeekName(t *testing.T) {
	result := Describe("0", DayOfWeek)
	if result != "at Sunday day of week" {
		t.Errorf("expected 'at Sunday day of week', got %q", result)
	}
}

func TestDescribeSchedule_Valid(t *testing.T) {
	fields := []string{"0", "9", "*", "*", "1-5"}
	result := DescribeSchedule(fields)
	expected := "at 0 minute; at 9 hour; every day of month; every month; from Monday to Friday"
	if result != expected {
		t.Errorf("expected %q, got %q", expected, result)
	}
}

func TestDescribeSchedule_InvalidFieldCount(t *testing.T) {
	fields := []string{"*", "*"}
	result := DescribeSchedule(fields)
	if result != "invalid cron expression" {
		t.Errorf("expected 'invalid cron expression', got %q", result)
	}
}
