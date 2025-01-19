package main

import (
	"cashflow/bdd/utils"
	"cashflow/models"
	"testing"
)

func TestPlans(t *testing.T) {
	t.Parallel()

	/*
		As an authenticated customer
		        I need to be able to create financial plans.
	*/

	utils.LoadEnv()
	token, err := utils.GetLoginToken()
	if err != nil {
		t.Fatalf(err.Error())
	}

	t.Run("List plans", func(t *testing.T) {
		t.Parallel()

		// Given I created plans.
		query := `
		mutation {
			createPlan(input: {
				name: "plan-1"
				fromAccount: "chequing"
				toAccount: "savings"
				startDate: "2022-02-21"
				recurrency: "FREQ=DAILY;COUNT=2"
				totalAmount: 100
			})
		}
		`
		var planId string
		err = utils.PostGraphQL(query, "createPlan", token, &planId)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// When I request the list of plans.
		query = `
		query {
			listPlans {
				id
				name
				toAccount
				fromAccount
				startDate
				recurrency
				totalAmount
			}
		}
		`
		var plans []models.Plan
		err = utils.PostGraphQL(query, "listPlans", token, &plans)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// Then I should receive a list of plans.
		if len(plans) != 1 {
			t.Fatalf("Expected 1 plan, got %d", len(plans))
		}

		if plans[0].Name != "plan-1" {
			t.Fatalf("Expected plan name to be 'plan-1', got %s", plans[0].Name)
		}

		if plans[0].FromAccount != "chequing" {
			t.Fatalf("Expected plan fromAccount to be 'chequing', got %s", plans[0].FromAccount)
		}

		if plans[0].ToAccount != "savings" {
			t.Fatalf("Expected plan toAccount to be 'savings', got %s", plans[0].ToAccount)
		}

		if plans[0].StartDate != "2022-02-21" {
			t.Fatalf("Expected plan startDate to be '2022-02-21', got %s", plans[0].StartDate)
		}

		if plans[0].TotalAmount != 100 {
			t.Fatalf("Expected plan totalAmount to be '100', got %f", plans[0].TotalAmount)
		}

		// cleanup
	})

	t.Run("Create a plan between existing accounts", func(t *testing.T) {
		t.Parallel()

		// Given I named my plan "ping-pong table".
		// And I need a total amount of 934.34 dollars.
		// And I start this plan on 2076-01-01.
		// And the money comes from my chequing account.
		// And the money goes to my savings account.
		// And I want to transfer money every 3 days, two times.
		// When I create a plan.
		query := `
		mutation {
			createPlan(input: {
				name: "ping-pong table"
				fromAccount: "chequing"
				toAccount: "savings"
				startDate: "2076-01-01"
				recurrency: "FREQ=DAILY;INTERVAL=3;COUNT=2"
				totalAmount: 934.34
			})
		}
		`
		var planId string
		err = utils.PostGraphQL(query, "createPlan", token, &planId)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// Then I should be able to list the transactions for this plan.
		query = `
		query {
			listTransactions(input: {
				from: "2076-01-01"
				to: "2076-01-10"
				account: "savings"
			}) {
				id
				amount
				date
				description
				account
			}
		}
		`
		var transactionsSavings []models.Transaction
		err = utils.PostGraphQL(query, "listTransactions", token, &transactionsSavings)
		if err != nil {
			t.Fatalf(err.Error())
		}

		query = `
		query {
			listTransactions(input: {
				from: "2076-01-01"
				to: "2076-01-10"
				account: "chequing"
			}) {
				id
				amount
				date
				description
				account
			}
		}
		`
		var transactionsChequing []models.Transaction
		err = utils.PostGraphQL(query, "listTransactions", token, &transactionsChequing)
		if err != nil {
			t.Fatalf(err.Error())
		}

		transactions := append(transactionsSavings, transactionsChequing...)

		expected := []models.Transaction{
			{
				Amount:      467.17,
				Date:        "2076-01-01",
				Description: "from plan: ping-pong table",
				Account:     "savings",
			},
			{
				Amount:      467.17,
				Date:        "2076-01-04",
				Description: "from plan: ping-pong table",
				Account:     "savings",
			},
			{
				Amount:      -467.17,
				Date:        "2076-01-01",
				Description: "from plan: ping-pong table",
				Account:     "chequing",
			},
			{
				Amount:      -467.17,
				Date:        "2076-01-04",
				Description: "from plan: ping-pong table",
				Account:     "chequing",
			},
		}

		for i, transaction := range transactions {
			if transaction.Amount != expected[i].Amount ||
				transaction.Date != expected[i].Date ||
				transaction.Description != expected[i].Description ||
				transaction.Account != expected[i].Account {
				t.Fatalf("Transaction %d does not match expected values.\nGot: %+v\nWant: %+v",
					i,
					transaction,
					expected[i])
			}
		}

		// cleanup
	})
}
