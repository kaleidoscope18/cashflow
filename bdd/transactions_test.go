package bdd

import (
	"context"

	"github.com/cucumber/godog"
)

func thereIsAChequingAccount() error {
	return nil
}

func iAddATransactionToIt() error {
	return nil
}

func iShouldBeAbleToSeeTheTransactions() error {
	return nil
}

func thereIsAnExistingTransactionInChequingAccount() error {
	return nil
}

func iRemoveIt() error {
	return nil
}

func itShouldBeRemoved() error {
	return nil
}

func iListTheTransactions() error {
	return nil
}

func InitializeTransactionsScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		return ctx, nil
	})

	ctx.Step(`^there is a chequing account$`, thereIsAChequingAccount)
	ctx.Step(`^I add a transaction to it$`, iAddATransactionToIt)
	ctx.Step(`^I should be able to see the transactions$`, iShouldBeAbleToSeeTheTransactions)
	ctx.Step(`^there is an existing transaction in chequing account$`, thereIsAnExistingTransactionInChequingAccount)
	ctx.Step(`^I remove it$`, iRemoveIt)
	ctx.Step(`^it should be removed$`, itShouldBeRemoved)
	ctx.Step(`^I list the transactions$`, iListTheTransactions)
}
