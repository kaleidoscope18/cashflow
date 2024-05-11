package utils

import (
	"testing"
)

type DatesComparisonTestData struct {
	Date1          string
	Date2          string
	ExpectedResult bool
}

var datesComparisonTestData = []DatesComparisonTestData{
	{"2000/02/01", "2000/01/01", false},
	{"2000/01/01", "2000/02/01", true},
	{"2000/01/01", "2000/01/11", true},
	{"2022/10/27", "2022/11/16", true},
}

func TestIsDateBefore(t *testing.T) {
	for _, d := range datesComparisonTestData {
		trueResult := IsDateBefore(d.Date1, d.Date2)
		if d.ExpectedResult != trueResult {
			t.Fatalf(`IsDateBefore("%s", "%s") should have given %t but resulted in %t instead`, d.Date1, d.Date2, d.ExpectedResult, trueResult)
		}
	}
}
