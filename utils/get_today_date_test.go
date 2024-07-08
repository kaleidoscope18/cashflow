package utils

import (
	"cashflow/dev"
	"testing"
)

func TestGetTodayDate(t *testing.T) {
	today := GetTodayDate()
	dev.PrintJson(today)
	location := ParseDateToTime(today).Location().String()
	if location != "UTC" {
		t.Fatalf("should have been %s but was %s", "UTC", location)
	}
}
