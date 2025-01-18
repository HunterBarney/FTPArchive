package datetime

import (
	"testing"
	"time"
)

func TestAMPM13(t *testing.T) {
	hour := 13
	answer := getAMPM(hour)

	if answer != "PM" {
		t.Errorf("getAMPM(13) got %s; want PM", answer)
	}
}

func TestAMPM12(t *testing.T) {
	hour := 12
	answer := getAMPM(hour)
	if answer != "PM" {
		t.Errorf("getAMPM(12) got %s; want PM", answer)
	}
}

func TestAMPM10(t *testing.T) {
	hour := 10
	answer := getAMPM(hour)
	if answer != "AM" {
		t.Errorf("getAMPM(12) got %s; want AM", answer)
	}
}

func TestAMPM1(t *testing.T) {
	hour := 1
	answer := getAMPM(hour)
	if answer != "AM" {
		t.Errorf("getAMPM(1) got %s; want AM", answer)
	}
}

func TestTo12Hour5(t *testing.T) {
	hour := 5
	answer := toTwelveHour(hour)
	if answer != 5 {
		t.Errorf("toTwelveHour(5) got %d; want 5", answer)
	}
}

func TestTo12Hour12(t *testing.T) {
	hour := 12
	answer := toTwelveHour(hour)
	if answer != 12 {
		t.Errorf("toTwelveHour(12) got %d; want 12", answer)
	}
}

func TestTo12Hour13(t *testing.T) {
	hour := 13
	answer := toTwelveHour(hour)
	if answer != 1 {
		t.Errorf("toTwelveHour(13) got %d; want 1", answer)
	}
}

func TestTo12Hour23(t *testing.T) {
	hour := 23
	answer := toTwelveHour(hour)
	if answer != 11 {
		t.Errorf("toTwelveHour(23) got %d; want 11", answer)
	}
}

func TestFormatDateTimedd(t *testing.T) {
	testTime := time.Date(2024, 12, 5, 11, 20, 40, 0, time.UTC)
	format := "%dd"
	answer := FormatDateTime(format, testTime)

	if answer != "05" {
		t.Errorf("FormatDateTime(%s) got %s; want 05", format, answer)
	}
}

func TestFormatDateTimed(t *testing.T) {
	testTime := time.Date(2024, 12, 5, 11, 20, 40, 0, time.UTC)
	format := "%d"
	answer := FormatDateTime(format, testTime)

	if answer != "5" {
		t.Errorf("FormatDateTime(%s) got %s; want 5", format, answer)
	}
}

func TestFormatDateMM(t *testing.T) {
	testTime := time.Date(2024, 03, 5, 11, 20, 40, 0, time.UTC)
	format := "%MM"
	answer := FormatDateTime(format, testTime)

	if answer != "03" {
		t.Errorf("FormatDateTime(%s) got %s; want 03", format, answer)
	}
}

func TestFormatDateM(t *testing.T) {
	testTime := time.Date(2024, 03, 5, 11, 20, 40, 0, time.UTC)
	format := "%M"
	answer := FormatDateTime(format, testTime)

	if answer != "3" {
		t.Errorf("FormatDateTime(%s) got %s; want 3", format, answer)
	}
}

func TestFormatDateyyyy(t *testing.T) {
	testTime := time.Date(2024, 03, 5, 11, 20, 40, 0, time.UTC)
	format := "%yyyy"
	answer := FormatDateTime(format, testTime)

	if answer != "2024" {
		t.Errorf("FormatDateTime(%s) got %s; want 2024", format, answer)
	}
}

func TestFormatDateyy(t *testing.T) {
	testTime := time.Date(2024, 03, 5, 11, 20, 40, 0, time.UTC)
	format := "%yy"
	answer := FormatDateTime(format, testTime)

	if answer != "24" {
		t.Errorf("FormatDateTime(%s) got %s; want 24", format, answer)
	}
}

func TestFormatDatehhhh(t *testing.T) {
	testTime := time.Date(2024, 03, 5, 13, 20, 40, 0, time.UTC)
	format := "%hhhh"
	answer := FormatDateTime(format, testTime)

	if answer != "01" {
		t.Errorf("FormatDateTime(%s) got %s; want 01", format, answer)
	}
}

func TestFormatDatehhh(t *testing.T) {
	testTime := time.Date(2024, 03, 5, 13, 20, 40, 0, time.UTC)
	format := "%hhh"
	answer := FormatDateTime(format, testTime)

	if answer != "1" {
		t.Errorf("FormatDateTime(%s) got %s; want 1", format, answer)
	}
}

func TestFormatDatehh(t *testing.T) {
	testTime := time.Date(2024, 03, 5, 01, 20, 40, 0, time.UTC)
	format := "%hh"
	answer := FormatDateTime(format, testTime)

	if answer != "01" {
		t.Errorf("FormatDateTime(%s) got %s; want 01", format, answer)
	}
}

func TestFormatDateh(t *testing.T) {
	testTime := time.Date(2024, 03, 5, 01, 20, 40, 0, time.UTC)
	format := "%h"
	answer := FormatDateTime(format, testTime)

	if answer != "1" {
		t.Errorf("FormatDateTime(%s) got %s; want 01", format, answer)
	}
}

func TestFormatDatemm(t *testing.T) {
	testTime := time.Date(2024, 03, 5, 01, 20, 40, 0, time.UTC)
	format := "%mm"
	answer := FormatDateTime(format, testTime)

	if answer != "20" {
		t.Errorf("FormatDateTime(%s) got %s; want 20", format, answer)
	}
}

func TestFormatDateAM(t *testing.T) {
	testTime := time.Date(2024, 03, 5, 11, 20, 40, 0, time.UTC)
	format := "%tt"
	answer := FormatDateTime(format, testTime)

	if answer != "AM" {
		t.Errorf("FormatDateTime(%s) got %s; want AM", format, answer)
	}
}

func TestFormatDatePM(t *testing.T) {
	testTime := time.Date(2024, 03, 5, 14, 20, 40, 0, time.UTC)
	format := "%tt"
	answer := FormatDateTime(format, testTime)

	if answer != "PM" {
		t.Errorf("FormatDateTime(%s) got %s; want PM", format, answer)
	}
}
