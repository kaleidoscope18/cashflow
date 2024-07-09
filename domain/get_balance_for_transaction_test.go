package domain

import (
	"cashflow/models"
	"fmt"
	"testing"
)

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
