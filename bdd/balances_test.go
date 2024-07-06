package bdd

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/cucumber/godog"
)

func thereIsAnAccount(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func iAddABalanceToIt(ctx context.Context) (context.Context, error) {
	postBody := strings.NewReader(`{"key":"value"}`)
	resp, err := http.Post("https://postman-echo.com/post", "application/json", postBody)
	if err != nil {
		return ctx, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ctx, errors.New("problem status")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ctx, errors.New("problem body deserialization")
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return ctx, errors.New("problem body validation")
	}

	if result["json"].(map[string]interface{})["key"] != "value" {
		return ctx, errors.New("problem body value")
	}

	return ctx, godog.ErrUndefined
}

func iListTheTransactions(ctx context.Context) (context.Context, error) {
	return ctx, godog.ErrPending
}

func iShouldBeAbleToListTheBalances(ctx context.Context) (context.Context, error) {
	return ctx, godog.ErrPending
}

func iShouldBeAbleToSeeTheTransactionsWithTheRightBalances(ctx context.Context) (context.Context, error) {
	return ctx, godog.ErrPending
}

func InitializeBalancesScenarioStepDefs(ctx *godog.ScenarioContext) {
	// Scenario: Setting the balance
	ctx.Step(`^there is an account$`, thereIsAnAccount)
	ctx.Step(`^I add a balance to it$`, iAddABalanceToIt)
	ctx.Step(`^I should be able to list the balances$`, iShouldBeAbleToListTheBalances)

	// Scenario: Listing transactions with balances
	ctx.Step(`^there is an account$`, thereIsAnAccount)
	ctx.Step(`^I add a balance to it$`, iAddABalanceToIt)
	ctx.Step(`^I list the transactions$`, iListTheTransactions)
	ctx.Step(`^I should be able to see the transactions with the right balances$`, iShouldBeAbleToSeeTheTransactionsWithTheRightBalances)
}
