package bdd

import (
	"github.com/cucumber/godog"
)

func iAddABalanceToIt() error {
	return godog.ErrPending
}

func iAddARecurringTransactionToIt() error {
	return godog.ErrPending
}

func iListTheTransactionsBetweenTwoDates() error {
	return godog.ErrPending
}

func iShouldBeAbleToListTheBalances() error {
	return godog.ErrPending
}

func iShouldBeAbleToSeeAllRecurringTransactions() error {
	return godog.ErrPending
}

func iShouldBeAbleToSeeTheTransactionsWithTheRightBalances() error {
	return godog.ErrPending
}

func thereIsAnAccount() error {
	return godog.ErrPending
}

func InitializeTransactionsScenarioStepDefs(ctx *godog.ScenarioContext) {
	ctx.Step(`^I add a balance to it$`, iAddABalanceToIt)
	ctx.Step(`^I add a recurring transaction to it$`, iAddARecurringTransactionToIt)
	ctx.Step(`^I list the transactions between two dates$`, iListTheTransactionsBetweenTwoDates)
	ctx.Step(`^I should be able to list the balances$`, iShouldBeAbleToListTheBalances)
	ctx.Step(`^I should be able to see all recurring transactions$`, iShouldBeAbleToSeeAllRecurringTransactions)
	ctx.Step(`^I should be able to see the transactions with the right balances$`, iShouldBeAbleToSeeTheTransactionsWithTheRightBalances)
	ctx.Step(`^there is an account$`, thereIsAnAccount)
}
