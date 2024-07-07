package status

import (
	"cashflow/models"
	"cashflow/utils"
	"testing"
)

type Expected struct {
	Date           string
	ExpectedResult models.Status
}

var expected = []Expected{
	{"2000/01/01", models.StatusDone},
	{"3000/12/30", models.StatusTodo},
	{utils.GetTodayDate(), models.StatusDone},
}

func TestGetStatusFromDate(t *testing.T) {
	for _, d := range expected {
		trueResult := GetStatusFromDate(utils.GetTodayDate(), d.Date)
		if d.ExpectedResult != trueResult {
			t.Fatalf(`GetStatusFromDate(today, "%s") should have given %s but resulted in %s instead`, d.Date, d.ExpectedResult, trueResult)
		}
	}
}
