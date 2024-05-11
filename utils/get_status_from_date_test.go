package utils

import (
	"cashflow/models"
	"testing"
)

type StatusFromDatesComparisonTestData struct {
	Date1          string
	Date2          string
	ExpectedResult models.Status
}

var statusFromDatesComparisonTestData = []StatusFromDatesComparisonTestData{
	{"2000/02/01", "2000/01/01", models.StatusDone},
	{"2000/01/01", "2000/02/01", models.StatusTodo},
	{"2000/01/01", "2000/01/11", models.StatusTodo},
	{"2022/10/27", "2022/11/16", models.StatusTodo},
	{"2022/12/27", "2022/11/16", models.StatusDone},
}

func TestGetStatusFromDate(t *testing.T) {
	for _, d := range statusFromDatesComparisonTestData {
		trueResult := GetStatusFromDate(&d.Date1, &d.Date2)
		if d.ExpectedResult != trueResult {
			t.Fatalf(`GetStatusFromDate("%s", "%s") should have given %s but resulted in %s instead`, d.Date1, d.Date2, d.ExpectedResult, trueResult)
		}
	}
}
