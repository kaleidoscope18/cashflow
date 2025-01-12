package recurrency

import (
	"cashflow/models"
	"reflect"
	"testing"
)

func TestSplitTransactionsWithRecurrency(t *testing.T) {
	expectedWith := []models.Transaction{
		{
			Date:        "2000-01-01",
			Amount:      2000,
			Id:          "1",
			Description: "",
			Recurrency:  "FREQ=MONTHLY;BYMONTHDAY=15",
		},
	}

	expectedWithout := []models.Transaction{
		{
			Date:        "2000-01-01",
			Amount:      2000,
			Id:          "2",
			Description: "",
			Recurrency:  "",
		},
	}

	resultWith, resultWithout := SplitTransactionsWithRecurrency(append(expectedWith, expectedWithout...))

	if !reflect.DeepEqual(expectedWith, resultWith) || !reflect.DeepEqual(expectedWithout, resultWithout) {
		t.Errorf("Expected %v, got %v", expectedWith, resultWith)
		t.Errorf("Expected %v, got %v", expectedWithout, resultWithout)
	}
}
