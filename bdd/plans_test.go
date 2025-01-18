package main

import (
	"cashflow/bdd/utils"
	"cashflow/models"
	cutils "cashflow/utils"
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
				recurrency: "FREQ=DAILY"
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

		if plans[0].StartDate.Format(cutils.DateLayout) != "2022-02-21" {
			t.Fatalf("Expected plan startDate to be '2022-02-21', got %s", plans[0].StartDate)
		}

		if plans[0].TotalAmount != 100 {
			t.Fatalf("Expected plan totalAmount to be '100', got %f", plans[0].TotalAmount)
		}

		// cleanup
	})

	t.Run("Create a plan between accounts", func(t *testing.T) {
		t.Parallel()

		// Given I named my plan "ping-pong table".

		// And I need a total amount of 934.34 dollars.

		// And I start this plan on 2025-01-01.

		// And the money comes from my chequing account.

		// And the money goes to my savings account.

		// And I want to transfer money every 3 days.

		// When I create a plan.

		// Then I should be able to list the transactions for this plan.

	})

	t.Run("Create a plan with external account", func(t *testing.T) {
		t.Parallel()

		// Given I named my plan "ping-pong table".

		// And I need a total amount of 934.34 dollars.

		// And I start this plan on 2025-01-01.

		// And the money comes from my chequing account.

		// And the money goes to an external account named "piggy bank".

		// And I want to transfer money every 3 days.

		// When I create a plan.

		// Then I should be able to list the transactions for this plan.

	})

	t.Run("Create a one-time plan", func(t *testing.T) {
		t.Parallel()

		// Given I named my plan "jelly beans".

		// And I need a total amount of 1 dollars.

		// And I start this plan on 2025-01-01.

		// And the money comes from my savings account.

		// And the money goes to my chequing account.

		// And I want to transfer money only once.

		// When I create a plan.

		// Then I should be able to list the single transaction for this plan.

	})
}
