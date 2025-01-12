package main

import (
	"cashflow/bdd/utils"
	"cashflow/models"
	"testing"
)

func TestEditTransactions(t *testing.T) {
	t.Parallel()

	/*
		As a customer
		I need to be able to edit my transactions.
	*/

	utils.LoadEnv()
	token, err := utils.GetLoginToken()
	if err != nil {
		t.Fatalf(err.Error())
	}

	t.Run("Editing a transaction", func(t *testing.T) {
		t.Parallel()

		// When I add a transaction.
		query := `
		mutation { 
			createTransaction(input: {
				amount: -10, 
				date: "october 27, 2022", 
				account: "chequing"
			})
		}
		`

		var id string
		err = utils.PostGraphQL(query, "createTransaction", token, &id)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// And I edit the transaction.
		query = `
		mutation { 
			editTransaction(input: {
				type: ALL
				id: "` + id + `"
				amount: 999
				account: "chequing"
			})
		}
		`
		err = utils.PostGraphQL(query, "editTransaction", token, nil)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// Then I should see edited transaction.
		query = `
		query {
			listTransactions(input: {
				from: "2022-10-27"
				to: "2022-10-27"
				account: "chequing"
			}) {
				id
				amount
			}
		}
		`
		var transactions []models.ComputedTransaction
		err = utils.PostGraphQL(query, "listTransactions", token, &transactions)
		if err != nil {
			t.Fatalf(err.Error())
		}

		if len(transactions) != 1 {
			t.Fatalf("expected 1 transaction, got %d", len(transactions))
		}

		if transactions[0].Amount != 999 {
			t.Fatalf("expected amount 999, got %f", transactions[0].Amount)
		}

		err = utils.CleanupTransactions([]string{id}, token)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})
}
