package utils

import (
	"cashflow/models"
	"fmt"
	"testing"
)

func TestSortByDateBalances(t *testing.T) {
	unorderedData := []models.Balance{
		{Amount: 2.00, Date: "2022/11/14"},
		{Amount: 1.00, Date: "2022/10/01"},
		{Amount: 3.00, Date: "2022/11/15"},
	}

	result := SortByDate(unorderedData)

	for i, r := range result {
		if r.Amount != float64(i+1) {
			t.Errorf(`Wrong order, the balance's amount at %d should be %f but was %f`, i+1, float64(i+1), r.Amount)
		}
	}
	fmt.Println(fmt.Sprint(result))
}

func TestSortByDateTransactions(t *testing.T) {
	unorderedData := []models.Transaction{
		{Amount: -91, Date: "2022/11/04"},
		{Amount: -117, Date: "2022/11/01"},
		{Amount: -1333, Date: "2022/11/01"},
		{Amount: 2800.69, Date: "2022/11/01"},
		{Amount: -788, Date: "2022/11/01"},
		{Amount: 2763.69, Date: "2022/11/15"},
		{Amount: -1000, Date: "2022/11/16"},
		{Amount: -115, Date: "2022/10/27"},
	}

	orderedData := []models.Transaction{
		{Amount: -115, Date: "2022/10/27"},
		{Amount: -117, Date: "2022/11/01"},
		{Amount: -1333, Date: "2022/11/01"},
		{Amount: 2800.69, Date: "2022/11/01"},
		{Amount: -788, Date: "2022/11/01"},
		{Amount: -91, Date: "2022/11/04"},
		{Amount: 2763.69, Date: "2022/11/15"},
		{Amount: -1000, Date: "2022/11/16"},
	}

	result := SortByDate(unorderedData)

	for i, r := range result {
		if r.Amount != orderedData[i].Amount {
			t.Errorf(`Wrong order, the balance's amount at %d should be %f but was %f`, i, orderedData[i].Amount, r.Amount)
		}
	}
	fmt.Println(fmt.Sprint(result))
}
