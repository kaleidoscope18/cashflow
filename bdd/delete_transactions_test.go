package bdd

import (
	"cashflow/bdd/utils"
	"cashflow/models"
	"slices"
	"testing"
)

func TestDeleteTransactions(t *testing.T) {
	/*
		As a customer
		        I need to be able to delete my transactions.
	*/

	utils.LoadEnv()
	token, err := utils.GetLoginToken()
	if err != nil {
		t.Fatalf(err.Error())
	}

	t.Run("Deleting a non-recurring transaction", func(t *testing.T) {
		// Given there is a chequing account with transactions.
		query := `
			mutation {
				createTransaction(input: {
					amount: 100,
					date: "2029-09-20",
					description: "test-delete-non-recurring",
					account: "chequing"
				})
			}
		`
		var id string
		err = utils.PostGraphQL(query, "createTransaction", token, &id)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// When I remove a transaction.
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

		// Then I should not see the transaction.
		query = `
			query {
				listTransactions(input: {
					account: "chequing"
					from: "2029-09-19"
					to: "2029-09-21"
				}) {
					id
					amount
					date
					description
					balance
					account
				}
			}
		`
		var transactions []models.ComputedTransaction
		err = utils.PostGraphQL(query, "listTransactions", token, &transactions)
		if err != nil {
			t.Fatalf(err.Error())
		}

		if len(transactions) != 0 {
			t.Fatalf("expected 0 transactions, got %d", len(transactions))
		}
	})

	t.Run("Deleting a recurring transaction", func(t *testing.T) {
		// When I add a recurring transaction to it on "2030-09-20" until "2031-09-31".
		query := `
			mutation {
				createTransaction(input: {
					amount: 100,
					date: "2030-09-20",
					description: "test-delete-recurring",
					recurrency: "FREQ=MONTHLY;BYMONTHDAY=1,15;UNTIL=20310930"
					account: "chequing"
				})
			}
		`
		var id string
		err = utils.PostGraphQL(query, "createTransaction", token, &id)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// And I delete the recurring transaction on "2030-12-15".
		query = `
			mutation { 
				deleteRecurringTransaction(input: {
					id: "` + id + `"
					date: "2030-12-15"
					account: "chequing"
				})
			}
		`
		err = utils.PostGraphQL(query, "deleteRecurringTransaction", token, nil)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// Then I should not see the recurring transaction on "2030-12-15" but see the rest.
		query = `
			query {
				listTransactions(input: {
					account: "chequing"
					from: "2030-09-20"
					to: "2031-09-30"
				}) {
					id
					amount
					date
					balance
					description
					account
				}
			}
		`
		var result []models.ComputedTransaction
		err = utils.PostGraphQL(query, "listTransactions", token, &result)
		if err != nil {
			t.Fatalf(err.Error())
		}

		var ids []string
		for _, r := range result {
			if r.Date == "2030-12-15" {
				t.Fatalf("found deleted transaction on %s with ID %s", r.Date, r.Id)
			}
			ids = append(ids, utils.ExtractIdFromSubid(r.Id))
		}

		// The test expects 23 transactions because:
		// - The recurring transaction occurs on the 1st and 15th of each month
		// - From Sept 20, 2030 to Sept 30, 2031 = ~12 months
		// - 12 months * 2 occurrences per month = 24 transactions
		// - Minus 1 deleted transaction on Dec 15, 2030
		// = 23 transactions total
		if len(result) != 23 {
			t.Fatalf("expected 23 transactions, got %d", len(result))
		}

		uniqueIds := slices.Compact(ids)
		if len(uniqueIds) != 2 {
			t.Fatalf("expected 2 unique ids, got %d", len(uniqueIds))
		}

		err = utils.CleanupTransactions(uniqueIds, token)
		if err != nil {
			t.Fatalf("failed to cleanup transactions with ids %v: %v", uniqueIds, err)
		}
	})
}
