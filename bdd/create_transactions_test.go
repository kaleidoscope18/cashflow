package bdd

import (
	"cashflow/bdd/utils"
	"cashflow/models"
	cfutils "cashflow/utils"
	"fmt"
	"strings"
	"testing"
)

func TestCreateTransactions(t *testing.T) {
	t.Parallel()

	/*
		As an authenticated customer
		        I need to be able to create my transactions.
	*/

	utils.LoadEnv()
	token, err := utils.GetLoginToken()
	if err != nil {
		t.Fatalf(err.Error())
	}

	t.Run("Creating a transaction", func(t *testing.T) {
		t.Parallel()

		// When I add a transaction.
		query := `
			mutation {
				createTransaction(input: {
					amount: 1000,
					date: "2025-11-30",
					account: "chequing"
					description: "test-create-transaction"
				})
			}
		`
		var id string
		err = utils.PostGraphQL(query, "createTransaction", token, &id)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// Then I should see the new transaction.
		query = `
			query {
				listTransactions(input: {
					account: "chequing"
					from: "2025-11-30"
					to: "2025-11-30"
				}) {
					id
					amount
					date
					account
				}
			}
		`
		var transactions []models.Transaction
		err = utils.PostGraphQL(query, "listTransactions", token, &transactions)
		if err != nil {
			t.Fatalf(err.Error())
		}
		if len(transactions) != 1 {
			t.Fatalf("expected 1 transaction, got %d", len(transactions))
		}
		if transactions[0].Amount != 1000 {
			t.Fatalf("expected amount of 1000, got %f", transactions[0].Amount)
		}

		err = utils.CleanupTransactions([]string{id}, token)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("Creating a recurring transaction", func(t *testing.T) {
		t.Parallel()

		// When I add a recurring transaction to it on "1940-09-20" until "1941-10-31".
		query := `
			mutation {
				createTransaction(input: {
					amount: 1000,
					date: "1940-09-20",
					account: "chequing"
					recurrency: "FREQ=MONTHLY;BYMONTHDAY=1,15;UNTIL=19411031"
					description: "test-create-recurring-transaction"
				})
			}
		`
		var id string
		err = utils.PostGraphQL(query, "createTransaction", token, &id)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// Then I should be able to see all recurring transactions.
		query = `
			query {
				listTransactions(input: {
					account: "chequing"
					from: "1940-09-20"
					to: "1941-10-31"
				}) {
					id
					amount
					date
					balance
				}
			}
		`
		var result []models.ComputedTransaction
		err = utils.PostGraphQL(query, "listTransactions", token, &result)
		if err != nil {
			t.Fatalf(err.Error())
		}

		expected := []models.ComputedTransaction{
			{
				Transaction: &models.Transaction{
					Date:   "1940-10-01",
					Amount: 1000,
				},
				Balance: 1000,
			},
			{
				Transaction: &models.Transaction{
					Date:   "1940-10-15",
					Amount: 1000,
				},
				Balance: 2000,
			},
			{
				Transaction: &models.Transaction{
					Date:   "1940-11-01",
					Amount: 1000,
				},
				Balance: 3000,
			},
			{
				Transaction: &models.Transaction{
					Date:   "1940-11-15",
					Amount: 1000,
				},
				Balance: 4000,
			},
			{
				Transaction: &models.Transaction{
					Date:   "1940-12-01",
					Amount: 1000,
				},
				Balance: 5000,
			},
			{
				Transaction: &models.Transaction{
					Date:   "1940-12-15",
					Amount: 1000,
				},
				Balance: 6000,
			},
			{
				Transaction: &models.Transaction{
					Date:   "1941-01-01",
					Amount: 1000,
				},
				Balance: 7000,
			},
			{
				Transaction: &models.Transaction{
					Date:   "1941-01-15",
					Amount: 1000,
				},
				Balance: 8000,
			},
			{
				Transaction: &models.Transaction{
					Date:   "1941-02-01",
					Amount: 1000,
				},
				Balance: 9000,
			},
			{
				Transaction: &models.Transaction{
					Date:   "1941-02-15",
					Amount: 1000,
				},
				Balance: 10000,
			},
			{
				Transaction: &models.Transaction{
					Date:   "1941-03-01",
					Amount: 1000,
				},
				Balance: 11000,
			},
			{
				Transaction: &models.Transaction{
					Date:   "1941-03-15",
					Amount: 1000,
				},
				Balance: 12000,
			},
			{
				Transaction: &models.Transaction{
					Date:   "1941-04-01",
					Amount: 1000,
				},
				Balance: 13000,
			},
			{
				Transaction: &models.Transaction{
					Date:   "1941-04-15",
					Amount: 1000,
				},
				Balance: 14000,
			},
			{
				Transaction: &models.Transaction{
					Date:   "1941-05-01",
					Amount: 1000,
				},
				Balance: 15000,
			},
			{
				Transaction: &models.Transaction{
					Date:   "1941-05-15",
					Amount: 1000,
				},
				Balance: 16000,
			},
			{
				Transaction: &models.Transaction{
					Date:   "1941-06-01",
					Amount: 1000,
				},
				Balance: 17000,
			},
			{
				Transaction: &models.Transaction{
					Date:   "1941-06-15",
					Amount: 1000,
				},
				Balance: 18000,
			},
			{
				Transaction: &models.Transaction{
					Date:   "1941-07-01",
					Amount: 1000,
				},
				Balance: 19000,
			},
			{
				Transaction: &models.Transaction{
					Date:   "1941-07-15",
					Amount: 1000,
				},
				Balance: 20000,
			},
			{
				Transaction: &models.Transaction{
					Date:   "1941-08-01",
					Amount: 1000,
				},
				Balance: 21000,
			},
			{
				Transaction: &models.Transaction{
					Date:   "1941-08-15",
					Amount: 1000,
				},
				Balance: 22000,
			},
			{
				Transaction: &models.Transaction{
					Date:   "1941-09-01",
					Amount: 1000,
				},
				Balance: 23000,
			},
			{
				Transaction: &models.Transaction{
					Date:   "1941-09-15",
					Amount: 1000,
				},
				Balance: 24000,
			},
			{
				Transaction: &models.Transaction{
					Date:   "1941-10-01",
					Amount: 1000,
				},
				Balance: 25000,
			},
			{
				Transaction: &models.Transaction{
					Date:   "1941-10-15",
					Amount: 1000,
				},
				Balance: 26000,
			},
		}

		for i, e := range expected {
			if e.Balance != result[i].Balance || e.Date != result[i].Date {
				t.Fatalf("for transaction with id %s, expected balance of %.2f on %s but got %.2f on %s",
					result[i].Id, e.Balance, e.Date, result[i].Balance, result[i].Date)
			}

			if !strings.Contains(result[i].Id, fmt.Sprintf("-%d", i)) {
				t.Fatalf("for transaction with id %s, expected id of x-%d", result[i].Id, i)
			}
		}

		err = utils.CleanupTransactions([]string{id}, token)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("Creating a mix of recurring and non recurring transactions", func(t *testing.T) {
		t.Parallel()

		// When I add a recurring transaction to it on "2020-09-20" until "2021-10-31".
		query := `
			mutation {
				createTransaction(input: {
					amount: 1000,
					date: "2020-09-20",
					account: "chequing"
					recurrency: "FREQ=MONTHLY;BYMONTHDAY=1,15;UNTIL=20211031"
					description: "test-create-mix-1"
				})
			}
		`
		var t1 string
		err = utils.PostGraphQL(query, "createTransaction", token, &t1)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// And I add a transaction on "2020-12-13".
		query = `
			mutation {
				createTransaction(input: {
					amount: 1000,
					date: "2020-12-13",
					account: "chequing"
					description: "test-create-mix-2"
				})
			}
		`
		var t2 string
		err = utils.PostGraphQL(query, "createTransaction", token, &t2)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// And I add a transaction on "2021-01-18".
		query = `
			mutation {
				createTransaction(input: {
					amount: 1000,
					date: "2021-01-18",
					account: "chequing"
					description: "test-create-mix-3"
				})
			}
		`
		var t3 string
		err = utils.PostGraphQL(query, "createTransaction", token, &t3)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// Then I should be able to see all transactions in order.
		query = `
			query {
				listTransactions(input: {
					account: "chequing"
					from: "2020-09-20"
					to: "2021-10-31"
				}) {
					id
					amount
					date
					account
					balance
				}
			}
		`

		var result []models.ComputedTransaction
		err = utils.PostGraphQL(query, "listTransactions", token, &result)
		if err != nil {
			t.Fatalf(err.Error())
		}

		for i, r := range result {
			if i == 0 {
				continue
			}

			a := cfutils.ParseDateToTime(result[i-1].Date)
			b := cfutils.ParseDateToTime(r.Date)
			if !b.After(a) && !b.Equal(a) {
				t.Fatalf("expected transaction on %s to be after %s", r.Date, result[i-1].Date)
			}
		}

		err = utils.CleanupTransactions([]string{t1, t2, t3}, token)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})
}
