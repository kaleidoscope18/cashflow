package main

import (
	"cashflow/bdd/utils"
	"cashflow/models"
	"testing"
	"time"
)

func TestBalances(t *testing.T) {
	t.Parallel()

	/*
		As a customer
		        I need to be able to set an end-of-day balance for any given day
		        And it will adjust the account's balances for each transaction.
	*/

	utils.LoadEnv()
	token, err := utils.GetLoginToken()
	if err != nil {
		t.Fatalf(err.Error())
	}

	t.Run("Adding a balance", func(t *testing.T) {
		t.Parallel()

		// When I add a balance to chequing account.
		query := `
		mutation { 
			createBalance(input: {
				amount: 1000,
				date: "2022-10-27",
				account: "chequing"
			}) {
				date
				amount
				account
			}
		}
		`
		var balance models.Balance
		err = utils.PostGraphQL(query, "createBalance", token, &balance)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// Then it should be in balances list.
		query = `
		query {
			listBalances(input: {
				account: "chequing"
				from: "2022-10-27"
				to: "2022-10-27"
			}) {
				date
				amount
				account
			}
		}
		`
		var balances []models.Balance
		err = utils.PostGraphQL(query, "listBalances", token, &balances)
		if err != nil {
			t.Fatalf(err.Error())
		}

		if len(balances) != 1 {
			t.Fatalf("expected 1 balance, got %d", len(balances))
		}

		if balances[0].Amount != 1000 {
			t.Fatalf("expected amount 1000, got %f", balances[0].Amount)
		}

		// and delete the balance
		query = `
		mutation { 
			deleteBalance(input: {
				date: "` + balances[0].Date + `"
				account: "chequing"
			})
		}
		`
		err = utils.PostGraphQL(query, "deleteBalance", token, nil)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("Adding a balance without a date", func(t *testing.T) {
		t.Parallel()

		// When I add a balance without a date to chequing account.
		query := `
		mutation { 
			createBalance(input: {
				amount: 1000,
				account: "chequing"
			}) {
				date
				amount
				account
			}
		}
		`
		var balance models.Balance
		err = utils.PostGraphQL(query, "createBalance", token, &balance)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// Then the new balance should have today's date.
		today := time.Now().UTC().Format("2006-01-02")
		if balance.Date != today {
			t.Fatalf("expected today's date (%s), got %s", today, balance.Date)
		}
	})

	t.Run("Listing transactions with balances", func(t *testing.T) {
		t.Parallel()

		// Given there is a chequing account with a transaction.
		query := `
		mutation { 
			createTransaction(input: {
				amount: 1000,
				date: "2031-10-30",
				account: "chequing"
				description: "test-list-transactions-with-balances"
			})
		}
		`
		var id string
		err = utils.PostGraphQL(query, "createTransaction", token, &id)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// When I add a balance to it.
		query = `
		mutation { 
			createBalance(input: {
				amount: 1000,
				date: "2031-10-29",
				account: "chequing"
			}) {
				date
				amount
				account
			}
		}
		`
		var balance models.Balance
		err = utils.PostGraphQL(query, "createBalance", token, &balance)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// Then I should be able to see the transactions with the right balance.
		query = `
		query {
			listTransactions(input: {
				account: "chequing"
				from: "2031-10-29"
				to: "2031-10-30"
			}) {
				id
				amount
				date
				balance
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

		if transactions[0].Balance != 2000.00 {
			t.Fatalf("expected balance of 2000, got %f", transactions[0].Balance)
		}

		// and delete the balance
		query = `
		mutation { 
			deleteBalance(input: {
				date: "` + balance.Date + `"
				account: "chequing"
			})
		}
		`
		err = utils.PostGraphQL(query, "deleteBalance", token, nil)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// and delete the transaction
		query = `
		mutation { 
			deleteTransaction(input: {
				id: "` + id + `"
				account: "chequing"
			})
		}
		`
		err = utils.PostGraphQL(query, "deleteTransaction", token, nil)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})
}
