package bdd

import (
	"cashflow/models"
	"cashflow/utils"
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cucumber/godog"
)

var (
	balanceDate        = contextKey("balance-date")
	newTransactionsIds = contextKey("new-transactions-ids")
	newTransactions    = contextKey("new-transactions")
)

func thereIsAnAccount(ctx context.Context) (context.Context, error) {
	return ctx, nil
}

func thereIsAnAccountWithTransactions(ctx context.Context) (context.Context, error) {
	query := `{"query": "mutation { createTransactions(input: [{amount: -10, date: \"october 27, 2022\"}, {amount: 100.10, date: \"november 15, 2022\", description: \"Paie\"}]) }"}`

	var ids []string
	err := PostGraphQL(ctx.Value(url).(string), query, "createTransactions", &ids)
	return context.WithValue(ctx, newTransactionsIds, &ids), err
}

func iAddABalanceToIt(ctx context.Context) (context.Context, error) {
	query := `{"query": "mutation { createBalance(input: {amount: 1000, date: \"October 30, 2022\"}) { date amount } }"}`

	var result models.Balance
	return context.WithValue(ctx, balanceDate, "2022/10/30"),
		PostGraphQL(ctx.Value(url).(string), query, "createBalance", &result)
}

func iListTheTransactions(ctx context.Context) (context.Context, error) {
	query := `{ "query": "query { listTransactions(from:\"1999-01-01T00:00:00.000Z\", to:\"3000-01-01T00:00:00.000Z\") { id date amount balance } }" }`

	var result []models.ComputedTransaction
	err := PostGraphQL(ctx.Value(url).(string), query, "listTransactions", &result)

	return context.WithValue(ctx, newTransactions, &result), err
}

func itShouldBeInBalancesList(ctx context.Context) (context.Context, error) {
	query := `{ "query": "query { listBalances(from:\"1999-01-01T00:00:00.000Z\", to:\"3000-01-01T00:00:00.000Z\") { date amount } }" }`

	var balances []models.Balance
	PostGraphQL(ctx.Value(url).(string), query, "listBalances", &balances)

	for _, balance := range balances {
		if balance.Date == "2022/10/30" {
			return ctx, nil
		}
	}

	return ctx, errors.New("balance with date \"2022/10/30\" was not found")
}

func iShouldBeAbleToSeeTheTransactionsWithTheRightBalances(ctx context.Context) (context.Context, error) {
	transactions := ctx.Value(newTransactions).(*[]models.ComputedTransaction)

	if len(*transactions) != 2 {
		return ctx, fmt.Errorf("There should only have been 2 transactions, got %d", len(*transactions))
	}

	if (*transactions)[0].Balance != -10 || (*transactions)[1].Balance != 1090.10 {
		return ctx, fmt.Errorf("The transactions should have balances of -10 and 1090.10, but were %.2f and %.2f", (*transactions)[0].Balance, (*transactions)[1].Balance)
	}

	return ctx, nil
}

func iAddABalanceWithoutADateToIt(ctx context.Context) (context.Context, error) {
	query := `{"query": "mutation { createBalance(input: {amount: 1000}) { date amount } }"}`

	var result models.Balance

	err := PostGraphQL(ctx.Value(url).(string), query, "createBalance", &result)
	return context.WithValue(ctx, balanceDate, result.Date), err
}

func theNewBalanceShouldHaveTodaysDate(ctx context.Context) (context.Context, error) {
	balanceDate := ctx.Value(balanceDate).(string)

	if balanceDate != utils.GetTodayDate() {
		return ctx, fmt.Errorf("date was supposed to be %s but was %s", utils.GetTodayDate(), balanceDate)
	}

	return ctx, nil
}

func InitializeBalancesScenarioStepDefs(ctx *godog.ScenarioContext) {
	ctx.Step(`^there is an account$`, thereIsAnAccount)
	ctx.Step(`^I add a balance to it$`, iAddABalanceToIt)
	ctx.Step(`^it should be in balances list$`, itShouldBeInBalancesList)
	ctx.Step(`^I add a balance without a date to it$`, iAddABalanceWithoutADateToIt)
	ctx.Step(`^the new balance should have today\'s date$`, theNewBalanceShouldHaveTodaysDate)
	ctx.Step(`^there is an account with transactions$`, thereIsAnAccountWithTransactions)
	ctx.Step(`^there is an account$`, thereIsAnAccount)
	ctx.Step(`^I add a balance to it$`, iAddABalanceToIt)
	ctx.Step(`^I list the transactions$`, iListTheTransactions)
	ctx.Step(`^I should be able to see the transactions with the right balances$`, iShouldBeAbleToSeeTheTransactionsWithTheRightBalances)

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		if ctx.Value(newTransactions) != nil {
			transactionsToDelete := *ctx.Value(newTransactions).(*[]models.ComputedTransaction)
			ids := make([]string, len(transactionsToDelete))
			for i, transaction := range transactionsToDelete {
				ids[i] = transaction.Id
			}
			jsonBytes, _ := json.Marshal(ids)

			query := fmt.Sprintf(`"query": "mutation { deleteTransactions(ids: %s) }"`, string(jsonBytes))
			PostGraphQL(ctx.Value(url).(string), query, "deleteTransactions", nil)
			ctx = context.WithValue(ctx, newTransactions, nil)
		}

		return ctx, nil
	})
}
