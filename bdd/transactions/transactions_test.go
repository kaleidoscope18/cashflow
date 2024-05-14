package bdd

import (
	"context"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

var NumberOfGodogs = 0

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

func TestMain(m *testing.M) {
	var options = godog.Options{
		Output: colors.Colored(os.Stdout),
		Format: "pretty",
		Paths:  []string{"transactions.feature"},
	}

	godog.BindCommandLineFlags("godog.", &options)

	status := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options:             &options,
	}.Run()

	os.Exit(status)
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		NumberOfGodogs = 0 // clean the state before every scenario
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
