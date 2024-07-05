package domain

import (
	"cashflow/models"
	"cashflow/utils"
	"fmt"
	"testing"
)

var startDate = utils.ParseDateToTime("1999/12/30")
var endDate = utils.ParseDateToTime("2000/01/11")

var balances = []models.Balance{
	{Date: "2000/01/01", Amount: 20},
	{Date: "2000/01/10", Amount: 100},
}
var transactions = []models.Transaction{
	{Id: "000", Amount: -2.99, Date: "1999/12/30", Description: "No balance yet"},
	{Id: "001", Amount: 2.00, Date: "2000/01/01", Description: "Balance will override this"},
	{Id: "002", Amount: -3.00, Date: "2000/01/02", Description: "Balance will decrease"},
	{Id: "003", Amount: 4.00, Date: "2000/01/03", Description: "Balance will increase"},
	{Id: "004", Amount: -5.00, Date: "2000/01/11", Description: "Balance will change and decrease"},
}
var expected = []models.ComputedTransaction{
	{Transaction: &transactions[0], Balance: -2.99},
	{Transaction: &transactions[1], Balance: 20},
	{Transaction: &transactions[2], Balance: 17},
	{Transaction: &transactions[3], Balance: 21},
	{Transaction: &transactions[4], Balance: 95},
}

func TestGetBalanceOnSameDayFound(t *testing.T) {
	result := getBalanceOnSameDay("2000/01/01", balances)

	if result == nil {
		t.Error("It should have found the balance on same day")
	}
}

func TestGetBalanceOnSameDayNotFound(t *testing.T) {
	result := getBalanceOnSameDay("2000/01/02", balances)

	if result != nil {
		t.Error("It not should have found the balance on same day")
	}
}

func TestGetLatestBalanceBefore(t *testing.T) {
	var testData = []struct {
		TransactionDate           string
		ExpectedLatestBalanceDate string
	}{
		{transactions[2].Date, balances[0].Date},
		{transactions[3].Date, balances[0].Date},
		{transactions[4].Date, balances[1].Date},
	}

	for _, d := range testData {
		result, _ := getLatestBalanceBefore(d.TransactionDate, balances)

		if d.ExpectedLatestBalanceDate != result.Date {
			t.Errorf(`getLatestBalanceBefore("%s", "%s") should have given balance on date %s but resulted in %s instead`,
				d.TransactionDate, fmt.Sprint(balances), d.ExpectedLatestBalanceDate, result.Date)
		}
	}
}

func TestGetLatestBalanceBeforeWithNoBalances(t *testing.T) {
	_, err := getLatestBalanceBefore(transactions[0].Date, []models.Balance{})

	if err == nil {
		t.Errorf(`getLatestBalanceBefore("%s", "%s") should have given an error`,
			transactions[0].Date, fmt.Sprint(balances))
	}
}

func TestGetLatestBalanceBeforeWithNoPreviousBalance(t *testing.T) {
	_, err := getLatestBalanceBefore(transactions[0].Date, balances)

	if err == nil {
		t.Errorf(`getLatestBalanceBefore("%s", "%s") should have given an error`,
			transactions[0].Date, fmt.Sprint(balances))
	}
}

func TestGetPreviousTransactionWithBalance(t *testing.T) {
	var testData = []struct {
		TransactionIndex int
		Expected         *models.ComputedTransaction
	}{
		{1, &expected[0]},
		{2, &expected[1]},
		{3, &expected[2]},
	}

	for _, d := range testData {
		previousTransactionWithBalance, err := getPreviousTransaction(d.TransactionIndex, expected)

		if err != nil {
			t.Fatalf(`getPreviousTransactionWithBalanceWithError("%d", "%s") should not have given an error`,
				d.TransactionIndex, fmt.Sprint(expected))
		}

		if d.Expected.Date != previousTransactionWithBalance.Date {
			t.Errorf(`getPreviousTransactionWithBalanceWithError("%d", "%s") should have given transaction on date %s`,
				d.TransactionIndex, fmt.Sprint(expected), d.Expected.Date)
		}
	}
}

func TestGetPreviousTransactionWithBalanceWithError(t *testing.T) {
	index := 0
	_, err := getPreviousTransaction(index, expected)

	if err == nil {
		t.Errorf(`getPreviousTransactionWithBalanceWithError("%d", "%s") should have given an error`, index, fmt.Sprint(expected))
	}
}

func TestGetBalanceForTransaction(t *testing.T) {
	var cases = []struct {
		*models.Transaction
		PreviousTransaction *models.ComputedTransaction
		LatestBalance       models.Balance
	}{
		{&transactions[1], &expected[0], balances[0]},
		{&transactions[2], &expected[1], balances[0]},
		{&transactions[3], &expected[2], balances[0]},
		{&transactions[4], &expected[3], balances[1]},
	}

	output := fmt.Sprintln("Expected | Actual")
	var fails = false
	for i, c := range cases {
		result := getBalanceForTransaction(*c.Transaction, *c.PreviousTransaction, c.LatestBalance)
		if result != expected[i+1].Balance {
			output += fmt.Sprintf("%.2f | %.2f\n", expected[i+1].Balance, result)
		} else {
			output += fmt.Sprintf("%.2f | %.2f [MATCH]\n", expected[i+1].Balance, result)
		}
	}

	if fails {
		fmt.Print(output)
		t.Errorf("Failed")
	}
}

func TestListTransactions(t *testing.T) {
	result, _ := listTransactions(transactions, balances, startDate, endDate)

	output := fmt.Sprintln("Expected | Actual")
	var fails = false
	for i, r := range result {
		if expected[i].Balance != r.Balance || expected[i].Date != r.Date {
			output += fmt.Sprintf("%.2f (%s) | %.2f (%s)\n", expected[i].Balance, expected[i].Date, r.Balance, r.Date)
			fails = true
		} else {
			output += fmt.Sprintf("%.2f (%s) | %.2f (%s) [MATCH]\n", expected[i].Balance, expected[i].Date, r.Balance, r.Date)
		}
	}

	if fails {
		fmt.Print(output)
		t.Errorf("Failed")
	}
}
