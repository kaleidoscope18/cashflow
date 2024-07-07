package bdd

import (
	"context"

	"github.com/cucumber/godog"
)

func iAddARecurringTransactionToIt(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func iShouldBeAbleToSeeAllRecurringTransactions(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func iAddATransactionToIt(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func thereIsAnExistingTransactionInChequingAccount(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func iRemoveIt(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func iShouldNotSeeTheNewTransaction(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func iShouldSeeTheNewTransaction(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func InitializeTransactionsScenarioStepDefs(ctx *godog.ScenarioContext) {
	ctx.Step(`^I add a recurring transaction to it$`, iAddARecurringTransactionToIt)
	ctx.Step(`^I add a transaction to it$`, iAddATransactionToIt)
	ctx.Step(`^I remove it$`, iRemoveIt)
	ctx.Step(`^I should be able to see all recurring transactions$`, iShouldBeAbleToSeeAllRecurringTransactions)
	ctx.Step(`^I should not see the new transaction$`, iShouldNotSeeTheNewTransaction)
	ctx.Step(`^I should see the new transaction$`, iShouldSeeTheNewTransaction)
	ctx.Step(`^there is an existing transaction in chequing account$`, thereIsAnExistingTransactionInChequingAccount)
}
