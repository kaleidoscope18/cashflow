package bdd

import (
	"cashflow/models"
	"context"
	"fmt"

	"github.com/cucumber/godog"
)

func thereIsAccount(ctx context.Context, accountName string) (context.Context, error) {
	return context.WithValue(ctx, account, accountName), nil
}

func thereIsNoAccount(ctx context.Context) (context.Context, error) {
	return context.WithValue(ctx, account, ""), nil
}

func thereIsAccountWithTransactions(ctx context.Context, accountName string) (context.Context, error) {
	query := `{"query": "mutation { createTransactions(input: [{amount: -10, date: \"october 27, 2022\"}, {amount: 100.10, date: \"november 15, 2022\", description: \"Paie\"}]) }"}`

	return ctx, PostGraphQL(ctx, query, "createTransactions", nil)
}

func iListTheTransactions(ctx context.Context) (context.Context, error) {
	query := `{"query": "query { listTransactions(from:\"1999-01-01T00:00:00.000Z\", to:\"3000-01-01T00:00:00.000Z\") { id date amount balance status description } }"}`

	var result []models.ComputedTransaction
	err := PostGraphQL(ctx, query, "listTransactions", &result)

	return context.WithValue(ctx, transactions, &result), err
}

func iListTheBalances(ctx context.Context) (context.Context, error) {
	query := `{ "query": "query { listBalances(from:\"1999-01-01T00:00:00.000Z\", to:\"3000-01-01T00:00:00.000Z\") { date amount } }" }`

	var result []models.Balance
	err := PostGraphQL(ctx, query, "listBalances", &result)

	return context.WithValue(ctx, balances, &result), err
}

func loginAs(ctx context.Context, role string) (context.Context, error) {
	if role == "authenticated" {
		userToken := ctx.Value(userToken)
		if userToken == nil {
			return ctx, fmt.Errorf("user token not found")
		}
		ctx = context.WithValue(ctx, token, userToken.(string))
		return ctx, nil
	}

	return ctx, nil
}

func InitializeCommonStepDefs(ctx *godog.ScenarioContext) {
	ctx.Step(`^login as "([^"]*)"$`, loginAs)
	ctx.Step(`^there is a (\w+) account$`, thereIsAccount)
	ctx.Step(`^there is a (\w+) account with transactions$`, thereIsAccountWithTransactions)
	ctx.Step(`^there is no account specified$`, thereIsNoAccount)
	ctx.Step(`^I list the balances$`, iListTheBalances)
	ctx.Step(`^I list the transactions$`, iListTheTransactions)
}
