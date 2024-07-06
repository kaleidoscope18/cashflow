package bdd

import (
	"github.com/cucumber/godog"
)

func thereIsAChequingAccount() error {
	return godog.ErrPending
}

func iAddARecurringTransactionToIt() error {
	return godog.ErrPending
}

func iListTheTransactionsBetweenTwoDates() error {
	return godog.ErrPending
}

func iShouldBeAbleToSeeAllRecurringTransactions() error {
	return godog.ErrPending
}

func iAddATransactionToIt() error {
	return godog.ErrPending
}

func iRemoveIt() error {
	return godog.ErrPending
}

func iShouldBeAbleToSeeTheTransactions() error {
	return godog.ErrPending
}

func itShouldBeRemoved() error {
	return godog.ErrPending
}

func thereIsAnExistingTransactionInChequingAccount() error {
	return godog.ErrPending
}

func InitializeTransactionsScenarioStepDefs(ctx *godog.ScenarioContext) {
	// Scenario: Adding a transaction
	ctx.Step(`^there is a chequing account$`, thereIsAChequingAccount)
	ctx.Step(`^I add a transaction to it$`, iAddATransactionToIt)
	ctx.Step(`^I should be able to see the transactions$`, iShouldBeAbleToSeeTheTransactions)

	// Scenario: Removing a transaction
	ctx.Step(`^there is an existing transaction in chequing account$`, thereIsAnExistingTransactionInChequingAccount)
	ctx.Step(`^I remove it$`, iRemoveIt)
	ctx.Step(`^it should be removed$`, itShouldBeRemoved)

	// Scenario: Listing transactions
	ctx.Step(`^there is a chequing account$`, thereIsAChequingAccount)
	ctx.Step(`^I list the transactions$`, iListTheTransactions)
	ctx.Step(`^I should be able to see the transactions$`, iShouldBeAbleToSeeTheTransactions)

	// Scenario: Adding a recurring transaction
	ctx.Step(`^there is a chequing account$`, thereIsAChequingAccount)
	ctx.Step(`^I add a recurring transaction to it$`, iAddARecurringTransactionToIt)
	ctx.Step(`^I list the transactions between two dates$`, iListTheTransactionsBetweenTwoDates)
	ctx.Step(`^I should be able to see all recurring transactions$`, iShouldBeAbleToSeeAllRecurringTransactions)
}
