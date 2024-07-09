package utils

import (
	"testing"
)

func TestGetTodayDate(t *testing.T) {
	today := GetTodayDate()
	location := ParseDateToTime(today).Location().String()
	if location != "UTC" {
		t.Fatalf("should have been %s but was %s", "UTC", location)
	}
}
