package main

import (
	"cashflow/bdd/utils"
	"cashflow/models"
	"testing"
)

func TestAccounts(t *testing.T) {
	t.Parallel()

	/*
		As an authenticated customer
		        I need to be able to manage my accounts.
	*/

	utils.LoadEnv()
	token, err := utils.GetLoginToken()
	if err != nil {
		t.Fatalf(err.Error())
	}

	t.Run("Filtering transactions from specific account", func(t *testing.T) {
		t.Parallel()

		// When I add a transaction on "2022-02-21" to savings account.
		query := `
		mutation { 
			createTransaction(input: {
				amount: 100,
				date: "2022-02-21",
				account: "savings"
				description: "savings-transaction-test"
			})
		}
		`
		var id1 string
		err = utils.PostGraphQL(query, "createTransaction", token, &id1)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// When I add a transaction on "2022-02-22" to chequing account.
		query = `
		mutation { 
			createTransaction(input: {
				amount: 100,
				date: "2022-02-22",
				account: "chequing"
				description: "chequing-transaction-test"
			})
		}
		`
		var id2 string
		err = utils.PostGraphQL(query, "createTransaction", token, &id2)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// Then I should see the new transaction only on specified account.
		query = `
		query {
			listTransactions(input: {
				account: "savings"
				from: "2022-02-20"
				to: "2022-02-23"
			}) {
				id
			}
		}
		`
		var transactions []models.Transaction
		err = utils.PostGraphQL(query, "listTransactions", token, &transactions)
		if err != nil {
			t.Fatalf(err.Error())
		}

		if len(transactions) != 1 {
			t.Fatalf("Expected 1 transaction, got %d", len(transactions))
		}

		err = utils.CleanupTransactions([]string{id2}, token)
		if err != nil {
			t.Fatalf(err.Error())
		}

		err = utils.CleanupSavingsTransactions([]string{id1}, token)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})
}
