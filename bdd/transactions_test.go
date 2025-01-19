package bdd

import (
	"cashflow/bdd/utils"
	"cashflow/models"
	"testing"
)

func TestTransactions(t *testing.T) {
	t.Parallel()

	/*
		As an authenticated customer
		        I need to be able to list all my transactions.
	*/

	utils.LoadEnv()
	token, err := utils.GetLoginToken()
	if err != nil {
		t.Fatalf(err.Error())
	}

	t.Run("List transactions", func(t *testing.T) {
		t.Parallel()

		// Given I created a plan starting on 1111-01-11.
		query := `
		mutation {
			createPlan(input: {
				name: "test-transactions-1"
				fromAccount: "savings"
				toAccount: "chequing"
				startDate: "1111-01-11"
				recurrency: "FREQ=WEEKLY;COUNT=2"
				totalAmount: 100
			})
		}
		`
		var planId string
		err := utils.PostGraphQL(query, "createPlan", token, &planId)
		if err != nil {
			t.Fatalf(err.Error())
		}
		defer utils.CleanupPlan(planId, token)

		// And I created a transaction on 1111-01-12.
		query = `
			mutation {
				createTransaction(input: {
					amount: 22,
					date: "1111-01-12",
					account: "chequing"
					description: "test-transactions-2"
				})
			}
		`
		var id1 string
		err = utils.PostGraphQL(query, "createTransaction", token, &id1)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// And I created a recurring transaction on 1111-01-13.
		query = `
			mutation {
				createTransaction(input: {
					amount: 55,
					date: "1111-01-13",
					account: "chequing"
					recurrency: "FREQ=DAILY;COUNT=2"
					description: "test-transactions-3"
				})
			}
		`
		var id2 string
		err = utils.PostGraphQL(query, "createTransaction", token, &id2)
		if err != nil {
			t.Fatalf(err.Error())
		}
		defer utils.CleanupTransactions([]string{id1, id2}, token)

		// When I list the transactions for this plan.
		query = `
			query {
				listTransactions(input: {
					account: "chequing"
					from: "1111-01-11"
					to: "1111-02-11"
				}) {
					id
					description
					amount
					date
					balance
					status
					account
				}
			}
		`
		var result []models.ComputedTransaction
		err = utils.PostGraphQL(query, "listTransactions", token, &result)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// Then I should see the transactions correctly.
		expected := []models.ComputedTransaction{
			{
				Transaction: &models.Transaction{
					Amount:      50,
					Date:        "1111-01-11",
					Description: "from plan: test-transactions-1",
					Account:     "chequing",
				},
				Balance: 50,
				Status:  "DONE",
			},
			{
				Transaction: &models.Transaction{
					Amount:      22,
					Date:        "1111-01-12",
					Description: "test-transactions-2",
					Account:     "chequing",
				},
				Balance: 72,
				Status:  "DONE",
			},
			{
				Transaction: &models.Transaction{
					Amount:      55,
					Date:        "1111-01-13",
					Description: "test-transactions-3",
					Account:     "chequing",
				},
				Balance: 127,
				Status:  "DONE",
			},
			{
				Transaction: &models.Transaction{
					Amount:      55,
					Date:        "1111-01-14",
					Description: "test-transactions-3",
					Account:     "chequing",
				},
				Balance: 182,
				Status:  "DONE",
			},
			{
				Transaction: &models.Transaction{
					Amount:      50,
					Date:        "1111-01-18",
					Description: "from plan: test-transactions-1",
					Account:     "chequing",
				},
				Balance: 232,
				Status:  "DONE",
			},
		}

		if len(expected) != len(result) {
			t.Errorf("length mismatch: expected %d transactions, got %d", len(expected), len(result))
			return
		}

		for i := range expected {
			exp := expected[i]
			got := result[i]
			if exp.Transaction.Amount != got.Transaction.Amount ||
				exp.Transaction.Date != got.Transaction.Date ||
				exp.Transaction.Description != got.Transaction.Description ||
				exp.Transaction.Account != got.Transaction.Account ||
				exp.Balance != got.Balance ||
				exp.Status != got.Status {
				t.Errorf("transaction %d mismatch:\nexpected: %+v\ngot: %+v", i, exp, got)
			}
		}
	})
}
