package bdd

import (
	"cashflow/models"
	"cashflow/utils"
	"context"
	"errors"
	"fmt"

	"github.com/cucumber/godog"
)

func iAddABalanceToIt(ctx context.Context) (context.Context, error) {
	query := `{"query": "mutation { createBalance(input: {amount: 1000, date: \"October 30, 2022\"}) { date amount } }"}`
	return ctx, PostGraphQL(ctx.Value(url).(string), query, "createBalance", nil)
}

func iAddABalanceWithoutADateToIt(ctx context.Context) (context.Context, error) {
	query := `{"query": "mutation { createBalance(input: {amount: 1000}) { date amount } }"}`
	return ctx, PostGraphQL(ctx.Value(url).(string), query, "createBalance", nil)
}

func itShouldBeInBalancesList(ctx context.Context) (context.Context, error) {
	b := *ctx.Value(balances).(*[]models.Balance)
	for _, balance := range b {
		if balance.Date == "2022/10/30" {
			return ctx, nil
		}
	}

	return ctx, errors.New("balance with date \"2022/10/30\" was not found")
}

func theNewBalanceShouldHaveTodaysDate(ctx context.Context) (context.Context, error) {
	today := utils.GetTodayDate()
	b := *ctx.Value(balances).(*[]models.Balance)

	for _, balance := range b {
		if balance.Date == today {
			return ctx, nil
		}
	}

	return ctx, errors.New("balance for today was not found")
}

func iShouldBeAbleToSeeTheTransactionsWithTheRightBalances(ctx context.Context) (context.Context, error) {
	transactions := ctx.Value(transactions).(*[]models.ComputedTransaction)

	if len(*transactions) != 2 {
		return ctx, fmt.Errorf("there should only have been 2 transactions, got %d", len(*transactions))
	}

	if (*transactions)[0].Balance != -10.00 || (*transactions)[1].Balance != 1100.10 {
		return ctx, fmt.Errorf("transactions should have balances of -10 and 1100.10, but were %.2f and %.2f", (*transactions)[0].Balance, (*transactions)[1].Balance)
	}

	return ctx, nil
}

func InitializeBalancesScenarioStepDefs(ctx *godog.ScenarioContext) {
	ctx.Step(`^I add a balance to it$`, iAddABalanceToIt)
	ctx.Step(`^it should be in balances list$`, itShouldBeInBalancesList)
	ctx.Step(`^I add a balance without a date to it$`, iAddABalanceWithoutADateToIt)
	ctx.Step(`^the new balance should have today\'s date$`, theNewBalanceShouldHaveTodaysDate)
	ctx.Step(`^I should be able to see the transactions with the right balances$`, iShouldBeAbleToSeeTheTransactionsWithTheRightBalances)
}
