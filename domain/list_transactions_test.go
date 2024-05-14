package domain

import (
	"cashflow/models"
	"cashflow/repository"
	"fmt"
	"testing"
)

var tr, br = repository.Init(models.Mocked)
var balances = (*br).ListBalances()
var transactions = (*tr).ListTransactions()
var transactionsWithBalances = []*models.ComputedTransaction{
	{Transaction: &transactions[0], Balance: -10},
	{Transaction: &transactions[1], Balance: 30},
	{Transaction: &transactions[2], Balance: 15},
	{Transaction: &transactions[3], Balance: -85},
	{Transaction: &transactions[4], Balance: 78},
}

func TestGetLatestBalanceBefore(t *testing.T) {
	var testData = []struct {
		models.Transaction
		ExpectedLatestBalanceDate string
	}{
		{transactions[1], balances[0].Date},
		{transactions[2], balances[0].Date},
		{transactions[3], balances[0].Date},
		{transactions[4], balances[1].Date},
	}

	for _, d := range testData {
		trueResult, _ := getLatestBalanceBefore(d.Date, balances)

		if d.ExpectedLatestBalanceDate != trueResult.Date {
			t.Errorf(`getLatestBalanceBefore("%s", "%s") should have given balance on date %s but resulted in %s instead`,
				d.Date, fmt.Sprint(balances), d.ExpectedLatestBalanceDate, trueResult.Date)
		}
	}
}

func TestGetLatestBalanceBeforeWithNoPreviousBalance(t *testing.T) {
	_, err := getLatestBalanceBefore(transactions[0].Date, balances)

	if err == nil {
		t.Errorf(`getLatestBalanceBefore("%s", "%s") should have given an error`,
			transactions[0].Date, fmt.Sprint(balances))
	}
}

func TestGetLatestBalanceBeforeWithNoBalances(t *testing.T) {
	_, err := getLatestBalanceBefore(transactions[0].Date, []models.Balance{})

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
		{1, transactionsWithBalances[0]},
		{2, transactionsWithBalances[1]},
		{3, transactionsWithBalances[2]},
	}

	for _, d := range testData {
		previousTransactionWithBalance, err := getPreviousTransaction(d.TransactionIndex, transactionsWithBalances)

		if err != nil {
			t.Fatalf(`getPreviousTransactionWithBalanceWithError("%d", "%s") should not have given an error`,
				d.TransactionIndex, fmt.Sprint(transactionsWithBalances))
		}

		if d.Expected.Date != previousTransactionWithBalance.Date {
			t.Errorf(`getPreviousTransactionWithBalanceWithError("%d", "%s") should have given transaction on date %s`,
				d.TransactionIndex, fmt.Sprint(transactionsWithBalances), d.Expected.Date)
		}
	}
}

func TestGetPreviousTransactionWithBalanceWithError(t *testing.T) {
	index := 0
	_, err := getPreviousTransaction(index, transactionsWithBalances)

	if err == nil {
		t.Errorf(`getPreviousTransactionWithBalanceWithError("%d", "%s") should have given an error`, index, fmt.Sprint(transactionsWithBalances))
	}
}

func TestGetBalanceForTransaction(t *testing.T) {
	var cases = []struct {
		*models.Transaction
		PreviousTransaction *models.ComputedTransaction
		LatestBalance       models.Balance
	}{
		{&transactions[1], transactionsWithBalances[0], balances[0]},
		{&transactions[2], transactionsWithBalances[1], balances[0]},
		{&transactions[3], transactionsWithBalances[2], balances[0]},
		{&transactions[4], transactionsWithBalances[3], balances[1]},
	}

	for i, c := range cases {
		result := getBalanceForTransaction(*c.Transaction, *c.PreviousTransaction, c.LatestBalance)
		expectedBalance := transactionsWithBalances[i+1].Balance
		if result != transactionsWithBalances[i+1].Balance {
			t.Errorf(`getBalanceForTransaction(...) for transaction on %s should have given %f but resulted in %f`,
				c.Transaction.Date, expectedBalance, result)
		}
	}
}

func TestListTransactionsWithBalances(t *testing.T) {
	tr, br := repository.Init(models.Mocked)
	bs := NewBalanceService(br)
	ts := NewTransactionService(tr, &bs)
	result, _ := ts.ListTransactions(nil)

	for i, r := range result {
		if transactionsWithBalances[i].Balance != r.Balance || transactionsWithBalances[i].Date != r.Date {
			t.Fatalf(`FATAL ERROR!`)
		}
	}
}
