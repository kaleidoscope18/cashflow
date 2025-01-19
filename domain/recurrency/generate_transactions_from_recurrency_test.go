package recurrency

import (
	"cashflow/models"
	"cashflow/utils"
	"testing"
)

func TestGenerateTransactionsFromRecurrency(t *testing.T) {
	var from = utils.ParseDateToTime("2000-01-01")
	var to = utils.ParseDateToTime("2000-03-01")

	data := []models.Transaction{
		{
			Date:        "2000-01-01",
			Amount:      2000,
			Id:          "1",
			Description: "Payday",
			Recurrency:  "FREQ=MONTHLY;BYMONTHDAY=15",
		},
	}
	results, _ := GenerateTransactionsFromRecurrency(data, from, to)
	expected := []models.Transaction{
		{
			Date:        "2000-01-15",
			Amount:      2000,
			Id:          "1-0",
			Description: "Payday",
			Recurrency:  "",
		},
		{
			Date:        "2000-02-15",
			Amount:      2000,
			Id:          "1-1",
			Description: "Payday",
			Recurrency:  "",
		},
	}

	for i := range results {
		if expected[i].Id != results[i].Id {
			t.Errorf("Expected ID %s, got %s", expected[i].Id, results[i].Id)
		}
		if expected[i].Date != results[i].Date {
			t.Errorf("Expected date %s, got %s", expected[i].Date, results[i].Date)
		}
		if expected[i].Description != results[i].Description {
			t.Errorf("Expected description %s, got %s", expected[i].Description, results[i].Description)
		}
		if expected[i].Amount != results[i].Amount {
			t.Errorf("Expected amount %2.f, got %2.f", expected[i].Amount, results[i].Amount)
		}
		if expected[i].Recurrency != results[i].Recurrency {
			t.Errorf("Expected recurrency %s, got %s", expected[i].Recurrency, results[i].Recurrency)
		}
	}
}

func TestNoRecurrencyInTransaction(t *testing.T) {
	var from = utils.ParseDateToTime("2000-01-01")
	var to = utils.ParseDateToTime("2000-03-01")
	data := []models.Transaction{
		{
			Date:        "2000-01-01",
			Amount:      2000,
			Id:          "1",
			Description: "Payday",
			Recurrency:  "",
		},
	}
	_, err := GenerateTransactionsFromRecurrency(data, from, to)
	if err == nil {
		t.Errorf("Expected error, got nil")
	}
}

func TestShouldHandleRecurrencyOnWeekdaysOnceAMonth(t *testing.T) {
	var from = utils.ParseDateToTime("2025-01-01")
	var to = utils.ParseDateToTime("2025-07-01")
	data := []models.Transaction{
		{
			Id:          "1",
			Date:        "2025-01-15",
			Amount:      2000,
			Description: "Mortgage",
			Recurrency: "RRULE:FREQ=MONTHLY;BYMONTHDAY=15;BYDAY=MO,TU,WE,TH,FR\n" +
				"RRULE:FREQ=MONTHLY;BYMONTHDAY=16;BYDAY=MO\n" +
				"RRULE:FREQ=MONTHLY;BYMONTHDAY=17;BYDAY=MO",
		},
	}
	results, _ := GenerateTransactionsFromRecurrency(data, from, to)
	expected := []models.Transaction{
		{
			Date:        "2025-01-15",
			Amount:      2000,
			Id:          "1-0",
			Description: "Mortgage",
			Recurrency:  "",
		},
		{
			Date:        "2025-02-17",
			Amount:      2000,
			Id:          "1-1",
			Description: "Mortgage",
			Recurrency:  "",
		},
		{
			Date:        "2025-03-17",
			Amount:      2000,
			Id:          "1-2",
			Description: "Mortgage",
			Recurrency:  "",
		},
		{
			Date:        "2025-04-15",
			Amount:      2000,
			Id:          "1-3",
			Description: "Mortgage",
			Recurrency:  "",
		},
		{
			Date:        "2025-05-15",
			Amount:      2000,
			Id:          "1-4",
			Description: "Mortgage",
			Recurrency:  "",
		},
		{
			Date:        "2025-06-16",
			Amount:      2000,
			Id:          "1-5",
			Description: "Mortgage",
			Recurrency:  "",
		},
	}

	for i := range results {
		if expected[i].Id != results[i].Id {
			t.Errorf("Expected ID %s, got %s", expected[i].Id, results[i].Id)
		}
		if expected[i].Date != results[i].Date {
			t.Errorf("Expected date %s, got %s", expected[i].Date, results[i].Date)
		}
		if expected[i].Description != results[i].Description {
			t.Errorf("Expected description %s, got %s", expected[i].Description, results[i].Description)
		}
		if expected[i].Amount != results[i].Amount {
			t.Errorf("Expected amount %2.f, got %2.f", expected[i].Amount, results[i].Amount)
		}
		if expected[i].Recurrency != results[i].Recurrency {
			t.Errorf("Expected recurrency %s, got %s", expected[i].Recurrency, results[i].Recurrency)
		}
	}
}
