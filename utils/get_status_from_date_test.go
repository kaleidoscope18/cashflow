package utils

import (
	"cashflow/models"
	"testing"
)

type Expected struct {
	Date           string
	ExpectedResult models.Status
}

var expected = []Expected{
	{"2000/01/01", models.StatusDone},
	{"3000/12/30", models.StatusTodo},
	{GetTodayDate(), models.StatusDone},
}

func TestGetStatusFromDate(t *testing.T) {
	for _, d := range expected {
		trueResult := GetStatusFromDate(GetTodayDate(), d.Date)
		if d.ExpectedResult != trueResult {
			t.Fatalf(`GetStatusFromDate(today, "%s") should have given %s but resulted in %s instead`, d.Date, d.ExpectedResult, trueResult)
		}
	}
}
