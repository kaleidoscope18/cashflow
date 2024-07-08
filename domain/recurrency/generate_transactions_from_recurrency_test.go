package recurrency

import (
	"cashflow/models"
	"cashflow/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

var from = utils.ParseDateToTime("2000/01/01")
var to = utils.ParseDateToTime("2000/03/01")

func TestGenerateTransactionsFromRecurrency(t *testing.T) {
	data := []models.Transaction{
		{
			Date:        "2000/01/01",
			Amount:      2000,
			Id:          "1",
			Description: "Payday",
			Recurrency:  "FREQ=MONTHLY;BYMONTHDAY=15",
		},
	}
	results, _ := GenerateTransactionsFromRecurrency(data, from, to)
	expected := []models.Transaction{
		{
			Date:        "2000/01/15",
			Amount:      2000,
			Id:          "1-0",
			Description: "Payday",
			Recurrency:  "",
		},
		{
			Date:        "2000/02/15",
			Amount:      2000,
			Id:          "1-1",
			Description: "Payday",
			Recurrency:  "",
		},
	}

	for i := range results {
		require.Equal(t, expected[i].Id, results[i].Id)
		require.Equal(t, expected[i].Date, results[i].Date)
		require.Equal(t, expected[i].Description, results[i].Description)
		require.Equal(t, expected[i].Amount, results[i].Amount)
		require.Equal(t, expected[i].Recurrency, results[i].Recurrency)
	}
}

func TestNoRecurrencyInTransaction(t *testing.T) {
	data := []models.Transaction{
		{
			Date:        "2000/01/01",
			Amount:      2000,
			Id:          "1",
			Description: "Payday",
			Recurrency:  "",
		},
	}
	_, err := GenerateTransactionsFromRecurrency(data, from, to)
	require.Error(t, err)
}
