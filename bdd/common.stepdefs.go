package bdd

import (
	"cashflow/models"
	"context"

	"github.com/cucumber/godog"
)

func thereIsAnAccount(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func thereIsAnAccountWithTransactions(ctx context.Context) (context.Context, error) {
	query := `{"query": "mutation { createTransactions(input: [{amount: -10, date: \"october 27, 2022\"}, {amount: 100.10, date: \"november 15, 2022\", description: \"Paie\"}]) }"}`

	return ctx, PostGraphQL(ctx.Value(url).(string), query, "createTransactions", nil)
}

func iListTheTransactions(ctx context.Context) (context.Context, error) {
	query := `{ "query": "query { listTransactions(from:\"1999-01-01T00:00:00.000Z\", to:\"3000-01-01T00:00:00.000Z\") { id date amount balance } }" }`

	var result []models.ComputedTransaction
	err := PostGraphQL(ctx.Value(url).(string), query, "listTransactions", &result)

	return context.WithValue(ctx, transactions, &result), err
}

func iListTheBalances(ctx context.Context) (context.Context, error) {
	query := `{ "query": "query { listBalances(from:\"1999-01-01T00:00:00.000Z\", to:\"3000-01-01T00:00:00.000Z\") { date amount } }" }`

	var result []models.Balance
	err := PostGraphQL(ctx.Value(url).(string), query, "listBalances", &result)

	return context.WithValue(ctx, balances, &result), err
}

func InitializeCommonStepDefs(ctx *godog.ScenarioContext) {
	ctx.Step(`^there is an account$`, thereIsAnAccount)
	ctx.Step(`^there is an account with transactions$`, thereIsAnAccountWithTransactions)

	ctx.Step(`^I list the balances$`, iListTheBalances)
	ctx.Step(`^I list the transactions$`, iListTheTransactions)
}
