package bdd

import (
	"cashflow/models"
	"cashflow/utils"
	"context"
	"fmt"
	"strings"

	"github.com/cucumber/godog"
)

func iAddARecurringTransactionToIt(ctx context.Context) (context.Context, error) {
	query := `{"query": "mutation {createTransaction(input: {amount: 1000, date: \"september 20, 2022\", recurrency: \"FREQ=MONTHLY;BYMONTHDAY=15,1;UNTIL=20231031T000000Z\", description: \"pay\"})}"}`
	return ctx, PostGraphQL(ctx.Value(url).(string), query, "createTransaction", nil)
}

func iShouldBeAbleToSeeAllRecurringTransactions(ctx context.Context) (context.Context, error) {
	transactions := *ctx.Value(transactions).(*[]models.ComputedTransaction)

	if len(transactions) != 26 {
		return ctx, fmt.Errorf("should have been 26 transactions (2 per month, for 13 months), but got %d", len(transactions))
	}

	expected := []models.ComputedTransaction{
		{
			Transaction: &models.Transaction{
				Date:   "2022/10/01",
				Amount: 1000,
			},
			Balance: 1000,
		},
		{
			Transaction: &models.Transaction{
				Date:   "2022/10/15",
				Amount: 1000,
			},
			Balance: 2000,
		},
		{
			Transaction: &models.Transaction{
				Date:   "2022/11/01",
				Amount: 1000,
			},
			Balance: 3000,
		},
		{
			Transaction: &models.Transaction{
				Date:   "2022/11/15",
				Amount: 1000,
			},
			Balance: 4000,
		},
		{
			Transaction: &models.Transaction{
				Date:   "2022/12/01",
				Amount: 1000,
			},
			Balance: 5000,
		},
		{
			Transaction: &models.Transaction{
				Date:   "2022/12/15",
				Amount: 1000,
			},
			Balance: 6000,
		},
		{
			Transaction: &models.Transaction{
				Date:   "2023/01/01",
				Amount: 1000,
			},
			Balance: 7000,
		},
		{
			Transaction: &models.Transaction{
				Date:   "2023/01/15",
				Amount: 1000,
			},
			Balance: 8000,
		},
		{
			Transaction: &models.Transaction{
				Date:   "2023/02/01",
				Amount: 1000,
			},
			Balance: 9000,
		},
		{
			Transaction: &models.Transaction{
				Date:   "2023/02/15",
				Amount: 1000,
			},
			Balance: 10000,
		},
		{
			Transaction: &models.Transaction{
				Date:   "2023/03/01",
				Amount: 1000,
			},
			Balance: 11000,
		},
		{
			Transaction: &models.Transaction{
				Date:   "2023/03/15",
				Amount: 1000,
			},
			Balance: 12000,
		},
		{
			Transaction: &models.Transaction{
				Date:   "2023/04/01",
				Amount: 1000,
			},
			Balance: 13000,
		},
		{
			Transaction: &models.Transaction{
				Date:   "2023/04/15",
				Amount: 1000,
			},
			Balance: 14000,
		},
		{
			Transaction: &models.Transaction{
				Date:   "2023/05/01",
				Amount: 1000,
			},
			Balance: 15000,
		},
		{
			Transaction: &models.Transaction{
				Date:   "2023/05/15",
				Amount: 1000,
			},
			Balance: 16000,
		},
		{
			Transaction: &models.Transaction{
				Date:   "2023/06/01",
				Amount: 1000,
			},
			Balance: 17000,
		},
		{
			Transaction: &models.Transaction{
				Date:   "2023/06/15",
				Amount: 1000,
			},
			Balance: 18000,
		},
		{
			Transaction: &models.Transaction{
				Date:   "2023/07/01",
				Amount: 1000,
			},
			Balance: 19000,
		},
		{
			Transaction: &models.Transaction{
				Date:   "2023/07/15",
				Amount: 1000,
			},
			Balance: 20000,
		},
		{
			Transaction: &models.Transaction{
				Date:   "2023/08/01",
				Amount: 1000,
			},
			Balance: 21000,
		},
		{
			Transaction: &models.Transaction{
				Date:   "2023/08/15",
				Amount: 1000,
			},
			Balance: 22000,
		},
		{
			Transaction: &models.Transaction{
				Date:   "2023/09/01",
				Amount: 1000,
			},
			Balance: 23000,
		},
		{
			Transaction: &models.Transaction{
				Date:   "2023/09/15",
				Amount: 1000,
			},
			Balance: 24000,
		},
		{
			Transaction: &models.Transaction{
				Date:   "2023/10/01",
				Amount: 1000,
			},
			Balance: 25000,
		},
		{
			Transaction: &models.Transaction{
				Date:   "2023/10/15",
				Amount: 1000,
			},
			Balance: 26000,
		},
	}

	for i, t := range expected {
		if t.Balance != transactions[i].Balance || t.Date != transactions[i].Date {
			utils.PrintJson(transactions)
			return ctx, fmt.Errorf("for transaction with id %s, expected balance of %.2f on %s but got %.2f on %s",
				transactions[i].Id, t.Balance, t.Date, transactions[i].Balance, transactions[i].Date)
		}

		if !strings.Contains(transactions[i].Id, fmt.Sprintf("-%d", i)) {
			return ctx, fmt.Errorf("for transaction with id %s, expected id of x-%d", transactions[i].Id, i)
		}
	}

	return ctx, nil
}

func iAddATransactionToIt(ctx context.Context) (context.Context, error) {
	query := `{"query": "mutation {createTransaction(input: {amount: 1000, date: \"october 20, 2022\"})}"}`
	return ctx, PostGraphQL(ctx.Value(url).(string), query, "createTransaction", nil)
}

func iRemoveATransaction(ctx context.Context) (context.Context, error) {
	transactions := *ctx.Value(transactions).(*[]models.ComputedTransaction)
	query := fmt.Sprintf(`{"query": "mutation {deleteTransaction(id: \"%s\")}"}`, transactions[0].Id)
	return ctx, PostGraphQL(ctx.Value(url).(string), query, "createTransaction", nil)
}

func iShouldSeeRemainingTransactions(ctx context.Context) (context.Context, error) {
	transactions := *ctx.Value(transactions).(*[]models.ComputedTransaction)

	if len(transactions) != 1 {
		return ctx, fmt.Errorf("expected 1 transaction, got %d", len(transactions))
	}
	return ctx, nil
}

func iShouldSeeTheNewTransaction(ctx context.Context) (context.Context, error) {
	transactions := *ctx.Value(transactions).(*[]models.ComputedTransaction)

	if transactions[0].Date != "2022/10/20" || transactions[0].Amount != 1000.00 || transactions[0].Balance != 1000.00 {
		return ctx, fmt.Errorf("expected amount and balance of 1000.00 on 2022/10/20 but got %.2f on %s",
			transactions[0].Balance, transactions[0].Date)
	}

	return ctx, nil
}

func InitializeTransactionsScenarioStepDefs(ctx *godog.ScenarioContext) {
	ctx.Step(`^I add a recurring transaction to it$`, iAddARecurringTransactionToIt)
	ctx.Step(`^I add a transaction to it$`, iAddATransactionToIt)
	ctx.Step(`^I remove a transaction$`, iRemoveATransaction)
	ctx.Step(`^I should be able to see all recurring transactions$`, iShouldBeAbleToSeeAllRecurringTransactions)
	ctx.Step(`^I should see remaining transactions$`, iShouldSeeRemainingTransactions)
	ctx.Step(`^I should see the new transaction$`, iShouldSeeTheNewTransaction)
}
