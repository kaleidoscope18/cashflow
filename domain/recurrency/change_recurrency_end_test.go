package recurrency

import (
	"testing"
)

func TestChangeRecurrencyEnd(t *testing.T) {
	var recurrency string = "FREQ=MONTHLY;UNTIL=20000301"

	results, err := ChangeRecurrencyEnd(recurrency, "2000-02-01")

	expected := "FREQ=MONTHLY;UNTIL=20000131T000000Z"

	if err != nil || results != expected {
		t.Errorf("Expected %s, got %s", expected, results)
	}
}
