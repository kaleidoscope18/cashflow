package main

import (
	"cashflow/bdd/utils"
	"cashflow/models"
	cutils "cashflow/utils"
	"slices"
	"strings"
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
				description: "test-edit-transaction"
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

	t.Run("Editing a recurring transaction", func(t *testing.T) {
		t.Parallel()

		// When I add a recurring transaction between "1999-09-20" until "2000-09-30".
		query := `
		mutation {
			createTransaction(input: {
				amount: 1,
				date: "1999-09-20",
				recurrency: "FREQ=MONTHLY;BYMONTHDAY=1;UNTIL=20000930"
				account: "chequing"
				description: "test-edit-recurring"
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

		// Then I should see edited transactions.
		query = `
		query {
			listTransactions(input: {
				from: "1999-09-20"
				to: "2000-09-30"
				account: "chequing"
			}) {
				id
				amount
				date
			}
		}
		`
		var transactions []models.ComputedTransaction
		err = utils.PostGraphQL(query, "listTransactions", token, &transactions)
		if err != nil {
			t.Fatalf(err.Error())
		}

		if len(transactions) != 12 {
			t.Fatalf("expected 12 transactions, got %d", len(transactions))
		}

		for _, transaction := range transactions {
			if transaction.Amount != 999 {
				t.Fatalf("expected amount 999, got %f", transaction.Amount)
			}
		}

		err = utils.CleanupTransactions([]string{id}, token)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	t.Run("Editing a recurring transaction from given date", func(t *testing.T) {
		t.Parallel()

		// When I add a recurring transaction to it on "2029-09-20" until "2030-10-31".
		query := `
		mutation {
			createTransaction(input: {
				amount: 1,
				date: "2029-09-20",
				recurrency: "FREQ=MONTHLY;BYMONTHDAY=1;UNTIL=20301031"
				account: "chequing"
				description: "test-edit-from-date-recurring"
			})
		}
		`
		var id string
		err = utils.PostGraphQL(query, "createTransaction", token, &id)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// And I edit the transaction from "2029-12-15".
		query = `
		mutation {
			editTransaction(input: {
				type: FROM_DATE
				id: "` + id + `"
				amount: 999
				account: "chequing"
				date: "2029-12-15"
			})
		}
		`
		err = utils.PostGraphQL(query, "editTransaction", token, nil)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// Then I should see unchanged transactions before "2029-12-15" and the rest should change.
		query = `
		query {
			listTransactions(input: {
				from: "2029-09-20"
				to: "2030-10-31"
				account: "chequing"
			}) {
				id
				date
				amount
				description
			}
		}
		`
		var transactions []models.ComputedTransaction
		err = utils.PostGraphQL(query, "listTransactions", token, &transactions)
		if err != nil {
			t.Fatalf(err.Error())
		}

		duplicateCount := 0
		var ids []string
		for _, r := range transactions {
			ids = append(ids, utils.ExtractIdFromSubid(r.Id))
			if cutils.ParseDateToTime(r.Date).After(cutils.ParseDateToTime("2029-12-15")) {
				if r.Amount != 999.00 {
					t.Fatalf("expected amount of 999.00 for transaction %s on %s but got %.2f", r.Id, r.Date, r.Amount)
				}
			}

			if cutils.ParseDateToTime(r.Date).Equal(cutils.ParseDateToTime("2029-12-15")) {
				duplicateCount++
				if duplicateCount > 1 {
					t.Fatalf("expected only one transaction on %s but got %d", "2029-12-15", duplicateCount)
				}
			}

			if r.Description == "" {
				t.Fatalf("expected description for transaction %s on %s but got empty strings", r.Id, r.Date)
			}
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

	t.Run("Editing a recurring transaction from given date but no date is provided", func(t *testing.T) {
		t.Parallel()

		// When I add a recurring transaction to it on "1980-09-20" until "1981-10-31".
		query := `
		mutation {
			createTransaction(input: {
				amount: 1,
				date: "1980-09-20",
				recurrency: "FREQ=MONTHLY;BYMONTHDAY=1;UNTIL=19811031"
				account: "chequing"
				description: "test-edit-from-date-recurring"
			})
		}
		`
		var id string
		err = utils.PostGraphQL(query, "createTransaction", token, &id)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// And I edit the transaction info from date without a date.
		query = `
		mutation {
			editTransaction(input: {
				type: FROM_DATE
				id: "` + id + `"
				amount: 999
				account: "chequing"
			})
		}
		`

		// Then I should have an expected date error.
		err = utils.PostGraphQL(query, "editTransaction", token, nil)
		if err == nil {
			t.Fatalf("expected error but got nil")
		}
		if !strings.Contains(err.Error(), "date is required for FROM_DATE type") {
			t.Fatalf("expected error message containing \"date is required for FROM_DATE type\" but got \"%s\"", err.Error())
		}

		err = utils.CleanupTransactions([]string{id}, token)
		if err != nil {
			t.Fatalf(err.Error())
		}
	})

	// TODO: Fix this test
	t.Run("Editing a recurring transaction only on date", func(t *testing.T) {
		t.Parallel()

		// When I add a recurring transaction to it on "1970-09-20" until "1971-10-31".
		query := `
		mutation {
			createTransaction(input: {
				amount: 1,
				date: "1970-09-20",
				recurrency: "FREQ=MONTHLY;BYMONTHDAY=1;UNTIL=19711031"
				account: "chequing"
				description: "test-edit-only-date-recurring"
			})
		}
		`
		var id string
		err = utils.PostGraphQL(query, "createTransaction", token, &id)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// And I edit the transaction on "1970-12-15".
		query = `
		mutation {
			editTransaction(input: {
				type: ON_DATE_ONLY
				id: "` + id + `"
				amount: 999
				account: "chequing"
				date: "1970-12-15"
			})
		}
		`
		err = utils.PostGraphQL(query, "editTransaction", token, nil)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// Then I should see unchanged transactions except transaction on "1970-12-15".
		query = `
		query {
			listTransactions(input: {
				from: "1970-09-20"
				to: "1971-10-31"
				account: "chequing"
			}) {
				id
				date
				amount
			}
		}
		`
		var transactions []models.ComputedTransaction
		err = utils.PostGraphQL(query, "listTransactions", token, &transactions)
		if err != nil {
			t.Fatalf(err.Error())
		}

		var ids []string
		for _, r := range transactions {
			ids = append(ids, utils.ExtractIdFromSubid(r.Id))
			if cutils.ParseDateToTime(r.Date).Equal(cutils.ParseDateToTime("1970-12-15")) && r.Amount != 999.00 {
				t.Fatalf("expected transaction on %s to be changed to 999.00 but got %.2f", "1970-12-15", r.Amount)
			}

			if !cutils.ParseDateToTime(r.Date).Equal(cutils.ParseDateToTime("1970-12-15")) && r.Amount != 1000 {
				t.Fatalf("expected transaction on %s to be unchanged but got %.2f", "1970-12-15", r.Amount)
			}
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

	// TODO: Fix this test
	t.Run("Editing a recurring transaction only on date but no date is provided", func(t *testing.T) {
		t.Parallel()

		// When I add a recurring transaction to it on "2022-09-20" until "2023-10-31".

		// And I edit the transaction info on date only without a date.

		// Then I should have an expected date error.

	})
}
