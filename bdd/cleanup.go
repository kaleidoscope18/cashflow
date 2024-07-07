package bdd

import (
	"cashflow/models"
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

func cleanupTransactions(ctx context.Context) context.Context {
	if ctx.Value(transactions) != nil {
		transactionsToDelete := *ctx.Value(transactions).(*[]models.ComputedTransaction)
		ids := make([]string, len(transactionsToDelete))
		for i, transaction := range transactionsToDelete {
			ids[i] = transaction.Id
		}
		jsonBytes, _ := json.Marshal(ids)

		query := fmt.Sprintf(`{"query": "mutation { deleteTransactions(ids: %s) }"}`, strings.Replace(string(jsonBytes), `"`, `\"`, -1))
		PostGraphQL(ctx.Value(url).(string), query, "deleteTransactions", nil)
	}

	return context.WithValue(ctx, transactions, nil)
}

func cleanupBalances(ctx context.Context) context.Context {
	if ctx.Value(balances) != nil {
		balancesToDelete := *ctx.Value(balances).(*[]models.Balance)
		for _, b := range balancesToDelete {
			query := fmt.Sprintf(`{"query": "mutation { deleteBalance(date: \"%s\") }"}`, b.Date)
			PostGraphQL(ctx.Value(url).(string), query, "deleteBalance", nil)
		}
	}

	return context.WithValue(ctx, balances, nil)
}
