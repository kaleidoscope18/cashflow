package bdd

import (
	"cashflow/bdd/utils"
	"cashflow/models"
	"testing"
	"time"
)

func TestTransactionsStatuses(t *testing.T) {
	t.Parallel()

	/*
		As a customer
		        I need to be able to see the statuses of my transactions.
	*/

	utils.LoadEnv()
	token, err := utils.GetLoginToken()
	if err != nil {
		t.Fatalf(err.Error())
	}

	t.Run("Getting transactions statuses", func(t *testing.T) {
		t.Parallel()

		// Given I add a transaction today and another later.
		query := `
			mutation {
				createTransaction(input: {
					amount: 1000,
					date: "` + time.Now().UTC().Format("2006-01-02") + `",
					account: "chequing"
				})
		}`

		var t1 string
		err = utils.PostGraphQL(query, "createTransaction", token, &t1)
		if err != nil {
			t.Fatalf(err.Error())
		}

		query = `
		mutation {
			createTransaction(input: {
				amount: 1000,
				date: "` + time.Now().UTC().AddDate(0, 0, 2).Format("2006-01-02") + `",
				account: "chequing"
			})
		}`
		var t2 string
		err = utils.PostGraphQL(query, "createTransaction", token, &t2)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// Then I should be able to see the correct statuses for each transaction.
		query = `
			query {
				listTransactions(input: {
					account: "chequing"
					from: "` + time.Now().UTC().AddDate(0, -1, 0).Format("2006-01-02") + `"
					to: "` + time.Now().UTC().AddDate(0, 1, 0).Format("2006-01-02") + `"
				}) {
					id
					status
					date
					description
					amount
					balance
				}
			}
		`
		var result []models.ComputedTransaction
		err = utils.PostGraphQL(query, "listTransactions", token, &result)
		if err != nil {
			t.Fatalf(err.Error())
		}

		if len(result) != 2 {
			t.Fatalf("expected 2 transactions, got %d", len(result))
		}

		if result[0].Status != models.StatusDone {
			t.Fatalf("expected transaction %s to be done, got %s", result[0].Id, result[0].Status)
		}

		if result[1].Status != models.StatusTodo {
			t.Fatalf("expected transaction %s to be todo, got %s", result[1].Id, result[1].Status)
		}

		err = utils.CleanupTransactions([]string{t1, t2}, token)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})
}
