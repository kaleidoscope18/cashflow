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
					date: "1111-01-13"
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

	t.Run("Transactions with recurrency on weekdays only", func(t *testing.T) {
		t.Parallel()

		// Given I created a recurring transaction on the 11th once a month on weekdays only.
		// And I want transactions to be postponed to the next weekday if it falls on a weekend.
		query := `
			mutation {
				createTransaction(input: {
					amount: 100
					date: "1645-01-01"
					account: "chequing"
					recurrency: "RRULE:FREQ=MONTHLY;BYMONTHDAY=15;BYDAY=MO,TU,WE,TH,FR\nRRULE:FREQ=MONTHLY;BYMONTHDAY=16;BYDAY=MO\nRRULE:FREQ=MONTHLY;BYMONTHDAY=17;BYDAY=MO"
				})
			}
		`
		var id string
		err = utils.PostGraphQL(query, "createTransaction", token, &id)
		if err != nil {
			t.Fatalf(err.Error())
		}
		defer utils.CleanupTransactions([]string{id}, token)

		// Then I should see the transactions correctly.
		query = `
			query {
				listTransactions(input: {
					account: "chequing"
					from: "1645-01-01"
					to: "1646-01-01"
				}) {
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
					Date: "1645-01-16",
				},
				Balance: 100,
			},
			{
				Transaction: &models.Transaction{
					Date: "1645-02-15",
				},
				Balance: 200,
			},
			{
				Transaction: &models.Transaction{
					Date: "1645-03-15",
				},
				Balance: 300,
			},
			{
				Transaction: &models.Transaction{
					Date: "1645-04-17",
				},
				Balance: 400,
			},
			{
				Transaction: &models.Transaction{
					Date: "1645-05-15",
				},
				Balance: 500,
			},
			{
				Transaction: &models.Transaction{
					Date: "1645-06-15",
				},
				Balance: 600,
			},
			{
				Transaction: &models.Transaction{
					Date: "1645-07-17",
				},
				Balance: 700,
			},
			{
				Transaction: &models.Transaction{
					Date: "1645-08-15",
				},
				Balance: 800,
			},
			{
				Transaction: &models.Transaction{
					Date: "1645-09-15",
				},
				Balance: 900,
			},
			{
				Transaction: &models.Transaction{
					Date: "1645-10-16",
				},
				Balance: 1000,
			},
			{
				Transaction: &models.Transaction{
					Date: "1645-11-15",
				},
				Balance: 1100,
			},
			{
				Transaction: &models.Transaction{
					Date: "1645-12-15",
				},
				Balance: 1200,
			},
		}

		if len(result) != len(expected) {
			t.Errorf("expected %d transactions, got %d", len(expected), len(result))
		}

		for i, exp := range expected {
			got := result[i]
			if exp.Date != got.Date ||
				exp.Balance != got.Balance {
				t.Errorf("transaction %d mismatch:\nexpected date: %s, balance: %.2f\ngot date: %s, balance: %.2f",
					i, exp.Date, exp.Balance, got.Date, got.Balance)
			}
		}
	})

	t.Run("Transactions and balances together", func(t *testing.T) {
		t.Parallel()

		// Given I create a realistic budget starting on 2025-01-01.
		query := `
			mutation {
				createTransaction(
					input: {
						date: "2025-02-15"
						amount: -1000
						description: "payment-test-1"
						recurrency: "RRULE:FREQ=MONTHLY;BYMONTHDAY=15;BYDAY=MO,TU,WE,TH,FR\nRRULE:FREQ=MONTHLY;BYMONTHDAY=16;BYDAY=MO\nRRULE:FREQ=MONTHLY;BYMONTHDAY=17;BYDAY=MO"
						account: "chequing"
					}
				)
			}
		`
		var id1 string
		err = utils.PostGraphQL(query, "createTransaction", token, &id1)
		if err != nil {
			t.Fatalf(err.Error())
		}

		query = `
			mutation {
				createTransaction(
					input: {
						date: "2025-02-16"
						amount: -2000
						description: "payment-test-2"
						recurrency: "RRULE:FREQ=MONTHLY;BYMONTHDAY=16;BYDAY=MO,TU,WE,TH,FR\nRRULE:FREQ=MONTHLY;BYMONTHDAY=17;BYDAY=MO\nRRULE:FREQ=MONTHLY;BYMONTHDAY=18;BYDAY=MO"
						account: "chequing"
					}
				)
			}
		`
		var id2 string
		err = utils.PostGraphQL(query, "createTransaction", token, &id2)
		if err != nil {
			t.Fatalf(err.Error())
		}

		query = `
			mutation {
				createTransaction(
					input: {
						date: "2025-01-08"
						amount: 2000
						description: "payday"
						recurrency: "FREQ=WEEKLY;INTERVAL=2;BYDAY=TH"
						account: "chequing"
					}
				)
			}
		`
		var id3 string
		err = utils.PostGraphQL(query, "createTransaction", token, &id3)
		if err != nil {
			t.Fatalf(err.Error())
		}
		defer utils.CleanupTransactions([]string{id1, id2, id3}, token)

		// And I set the balance on 2025-01-19 to 318$.
		query = `
			mutation {
				createBalance(input: {
					date: "2025-01-19"
					amount: 318
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
		defer utils.CleanupBalance(balance.Date, balance.Account, token)

		// And I set the balance on 2025-04-16 to 1$.
		query = `
			mutation {
				createBalance(input: {
					date: "2025-04-16"
					amount: 1
					account: "chequing"
				}) {
					date
					amount
					account
				}
			}
		`
		var balance2 models.Balance
		err = utils.PostGraphQL(query, "createBalance", token, &balance2)
		if err != nil {
			t.Fatalf(err.Error())
		}
		defer utils.CleanupBalance(balance2.Date, balance2.Account, token)

		// Then I should see the transactions correctly.
		query = `
			query {
				listTransactions(input: {
					account: "chequing"
					from: "2025-01-01"
					to: "2025-06-01"
				}) {
					id
					date
					amount
					account
					balance
					description
				}
			}
		`
		var result []models.ComputedTransaction
		err = utils.PostGraphQL(query, "listTransactions", token, &result)
		if err != nil {
			t.Fatalf(err.Error())
		}

		expected := []models.ComputedTransaction{
			{Transaction: &models.Transaction{Amount: 2000, Date: "2025-01-09", Description: "payday", Account: "chequing"}, Balance: 2000},
			{Transaction: &models.Transaction{Amount: 2000, Date: "2025-01-23", Description: "payday", Account: "chequing"}, Balance: 2318},
			{Transaction: &models.Transaction{Amount: 2000, Date: "2025-02-06", Description: "payday", Account: "chequing"}, Balance: 4318},
			{Transaction: &models.Transaction{Amount: -1000, Date: "2025-02-17", Description: "payment-test-1", Account: "chequing"}, Balance: 3318},
			{Transaction: &models.Transaction{Amount: -2000, Date: "2025-02-17", Description: "payment-test-2", Account: "chequing"}, Balance: 1318},
			{Transaction: &models.Transaction{Amount: 2000, Date: "2025-02-20", Description: "payday", Account: "chequing"}, Balance: 3318},
			{Transaction: &models.Transaction{Amount: 2000, Date: "2025-03-06", Description: "payday", Account: "chequing"}, Balance: 5318},
			{Transaction: &models.Transaction{Amount: -1000, Date: "2025-03-17", Description: "payment-test-1", Account: "chequing"}, Balance: 4318},
			{Transaction: &models.Transaction{Amount: -2000, Date: "2025-03-17", Description: "payment-test-2", Account: "chequing"}, Balance: 2318},
			{Transaction: &models.Transaction{Amount: 2000, Date: "2025-03-20", Description: "payday", Account: "chequing"}, Balance: 4318},
			{Transaction: &models.Transaction{Amount: 2000, Date: "2025-04-03", Description: "payday", Account: "chequing"}, Balance: 6318},
			{Transaction: &models.Transaction{Amount: -1000, Date: "2025-04-15", Description: "payment-test-1", Account: "chequing"}, Balance: 5318},
			{Transaction: &models.Transaction{Amount: -2000, Date: "2025-04-16", Description: "payment-test-2", Account: "chequing"}, Balance: 3318},
			{Transaction: &models.Transaction{Amount: 2000, Date: "2025-04-17", Description: "payday", Account: "chequing"}, Balance: 2001},
			{Transaction: &models.Transaction{Amount: 2000, Date: "2025-05-01", Description: "payday", Account: "chequing"}, Balance: 4001},
			{Transaction: &models.Transaction{Amount: 2000, Date: "2025-05-15", Description: "payday", Account: "chequing"}, Balance: 6001},
			{Transaction: &models.Transaction{Amount: -1000, Date: "2025-05-15", Description: "payment-test-1", Account: "chequing"}, Balance: 5001},
			{Transaction: &models.Transaction{Amount: -2000, Date: "2025-05-16", Description: "payment-test-2", Account: "chequing"}, Balance: 3001},
			{Transaction: &models.Transaction{Amount: 2000, Date: "2025-05-29", Description: "payday", Account: "chequing"}, Balance: 5001},
		}

		if len(result) != len(expected) {
			t.Fatalf("expected %d transactions, got %d", len(expected), len(result))
		}

		for i, exp := range expected {
			got := result[i]
			if got.Amount != exp.Amount {
				t.Errorf("transaction %d: expected amount %f, got %f", i, exp.Amount, got.Amount)
			}
			if got.Date != exp.Date {
				t.Errorf("transaction %d: expected date %s, got %s", i, exp.Date, got.Date)
			}
			if got.Description != exp.Description {
				t.Errorf("transaction %d: expected description %s, got %s", i, exp.Description, got.Description)
			}
			if got.Account != exp.Account {
				t.Errorf("transaction %d: expected account %s, got %s", i, exp.Account, got.Account)
			}
			if got.Balance != exp.Balance {
				t.Errorf("transaction %d: expected balance %f, got %f", i, exp.Balance, got.Balance)
			}
			if got.Status != exp.Status {
				t.Errorf("transaction %d: expected status %s, got %s", i, exp.Status, got.Status)
			}
		}
	})
}
