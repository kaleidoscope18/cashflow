package bdd

import (
	"cashflow/bdd/utils"
	"cashflow/models"
	"testing"
)

func TestPlans(t *testing.T) {
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
		var found models.Plan
		for _, plan := range plans {
			if plan.Id == planId {
				found = plan
				break
			}
		}
		if found.Id == "" {
			t.Fatalf("Could not find plan with id %s", planId)
		}

		if found.Name != "plan-1" {
			t.Fatalf("Expected plan name to be 'plan-1', got %s", found.Name)
		}

		if found.FromAccount != "chequing" {
			t.Fatalf("Expected plan fromAccount to be 'chequing', got %s", found.FromAccount)
		}

		if found.ToAccount != "savings" {
			t.Fatalf("Expected plan toAccount to be 'savings', got %s", found.ToAccount)
		}

		if found.StartDate != "2022-02-21" {
			t.Fatalf("Expected plan startDate to be '2022-02-21', got %s", found.StartDate)
		}

		if found.TotalAmount != 100 {
			t.Fatalf("Expected plan totalAmount to be '100', got %f", found.TotalAmount)
		}

		utils.CleanupPlan(planId, token)
	})

	t.Run("Create a plan between existing accounts", func(t *testing.T) {
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

		utils.CleanupPlan(planId, token)
	})

	t.Run("Create a plan to an external account", func(t *testing.T) {
		// Given I named my plan "popcorn maker".
		// And I need a total amount of 300 dollars.
		// And I start this plan on 2084-01-01.
		// And the money comes from my chequing account.
		// And the money goes to an external account named "piggy bank".
		// And I want to transfer money every 2 days, ten times.
		// When I create a plan.
		query := `
		mutation {
			createPlan(input: {
				name: "popcorn maker"
				fromAccount: "chequing"
				toAccount: "piggy bank"
				startDate: "2084-01-01"
				recurrency: "FREQ=DAILY;INTERVAL=2;COUNT=10"
				totalAmount: 300
			})
		}
		`
		var planId string
		err = utils.PostGraphQL(query, "createPlan", token, &planId)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// Then I should be able to list the plans.
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

		// Then I should be able to list the transactions for this plan.
		query = `
		query {
			listTransactions(input: {
				from: "2084-01-01"
				to: "2084-01-30"
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
		var transactions []models.Transaction
		err = utils.PostGraphQL(query, "listTransactions", token, &transactions)
		if err != nil {
			t.Fatalf(err.Error())
		}

		expected := []models.Transaction{
			{
				Amount:      -30,
				Date:        "2084-01-01",
				Description: "from plan: popcorn maker",
				Account:     "chequing",
			},
			{
				Amount:      -30,
				Date:        "2084-01-03",
				Description: "from plan: popcorn maker",
				Account:     "chequing",
			},
			{
				Amount:      -30,
				Date:        "2084-01-05",
				Description: "from plan: popcorn maker",
				Account:     "chequing",
			},
			{
				Amount:      -30,
				Date:        "2084-01-07",
				Description: "from plan: popcorn maker",
				Account:     "chequing",
			},
			{
				Amount:      -30,
				Date:        "2084-01-09",
				Description: "from plan: popcorn maker",
				Account:     "chequing",
			},
			{
				Amount:      -30,
				Date:        "2084-01-11",
				Description: "from plan: popcorn maker",
				Account:     "chequing",
			},
			{
				Amount:      -30,
				Date:        "2084-01-13",
				Description: "from plan: popcorn maker",
				Account:     "chequing",
			},
			{
				Amount:      -30,
				Date:        "2084-01-15",
				Description: "from plan: popcorn maker",
				Account:     "chequing",
			},
			{
				Amount:      -30,
				Date:        "2084-01-17",
				Description: "from plan: popcorn maker",
				Account:     "chequing",
			},
			{
				Amount:      -30,
				Date:        "2084-01-19",
				Description: "from plan: popcorn maker",
				Account:     "chequing",
			},
		}
		if len(transactions) != 10 {
			t.Fatalf("Expected 10 transactions, got %d", len(transactions))
		}
		for i, e := range expected {
			if transactions[i].Amount != e.Amount ||
				transactions[i].Date != e.Date ||
				transactions[i].Description != e.Description ||
				transactions[i].Account != e.Account {
				t.Fatalf("Transaction %d does not match expected values.\nGot: %+v\nWant: %+v",
					i,
					transactions[i],
					e)
			}
		}

		utils.CleanupPlan(planId, token)
	})

	t.Run("Create a plan from an external account", func(t *testing.T) {
		// Given I named my plan "smoothie machine".
		// And I need a total amount of 2 dollars.
		// And I start this plan on 2081-01-01.
		// And the money comes from an external account named "piggy bank".
		// And the money goes to my savings account.
		// And I want to transfer money every 2 days, two times.
		// When I create a plan.
		query := `
		mutation {
			createPlan(input: {
				name: "smoothie machine"
				fromAccount: "piggy bank"
				toAccount: "savings"
				startDate: "2081-01-01"
				recurrency: "FREQ=DAILY;INTERVAL=2;COUNT=2"
				totalAmount: 2
			})
		}
		`
		var planId string
		err = utils.PostGraphQL(query, "createPlan", token, &planId)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// Then I should be able to list the plans.
		query = `
		query {
			listPlans {
				id
			}
		}
		`
		var plans []models.Plan
		err = utils.PostGraphQL(query, "listPlans", token, &plans)
		if err != nil {
			t.Fatalf(err.Error())
		}

		found := false
		for _, plan := range plans {
			if plan.Id == planId {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Could not find plan with id %s", planId)
		}

		// Then I should be able to list the transactions for this plan only on savings.
		query = `
		query {
			listTransactions(input: {
				from: "2081-01-01"
				to: "2081-01-30"
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
		var transactions []models.Transaction
		err = utils.PostGraphQL(query, "listTransactions", token, &transactions)
		if err != nil {
			t.Fatalf(err.Error())
		}

		if len(transactions) != 2 {
			t.Fatalf("Expected 2 transactions, got %d", len(transactions))
		}

		utils.CleanupPlan(planId, token)
	})

	t.Run("Delete a plan", func(t *testing.T) {
		// Given I created a plan.
		query := `
		mutation {
			createPlan(input: {
				name: "popcorn maker"
				fromAccount: "chequing"
				toAccount: "piggy bank"
				startDate: "2099-01-01"
				recurrency: "FREQ=DAILY;COUNT=1"
				totalAmount: 300
			})
		}
		`
		var planId string
		err = utils.PostGraphQL(query, "createPlan", token, &planId)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// When I delete the plan.
		query = `
		mutation {
			deletePlan(id: "` + planId + `")
		}
		`
		err = utils.PostGraphQL(query, "deletePlan", token, nil)
		if err != nil {
			t.Fatalf(err.Error())
		}

		// Then I should not see the plan or its transactions anymore.
		query = `
		query {
			listPlans {
				id
			}
		}
		`
		var plans []models.Plan
		err = utils.PostGraphQL(query, "listPlans", token, &plans)
		if err != nil {
			t.Fatalf(err.Error())
		}
		for _, plan := range plans {
			if plan.Id == planId {
				t.Fatalf("Expected plan to be deleted, but found plan with id %s", planId)
			}
		}

		query = `
		query {
			listTransactions(input: {
				from: "2099-01-01"
				to: "2099-01-30"
				account: "chequing"
			}) {
				id
			}
		}
		`
		var transactions []models.Transaction
		err = utils.PostGraphQL(query, "listTransactions", token, &transactions)
		if err != nil {
			t.Fatalf(err.Error())
		}

		if len(transactions) != 0 {
			t.Fatalf("Expected 0 transactions, got %d", len(transactions))
		}

		utils.CleanupPlan(planId, token)
	})

	t.Run("Create a plan with both external accounts", func(t *testing.T) {
		// Given I create a plan with both external accounts.
		query := `
		mutation {
			createPlan(input: {
				name: "smoothie machine"
				fromAccount: "piggy bank"
				toAccount: "dream land"
				startDate: "2081-01-01"
				recurrency: "FREQ=DAILY;COUNT=1"
				totalAmount: 2
			})
		}
		`
		err = utils.PostGraphQL(query, "createPlan", token, nil)

		// Then I should not be able to create the plan.
		if err == nil {
			t.Fatalf("Expected error, got nil")
		}
	})
}
