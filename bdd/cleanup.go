package bdd

import (
	"cashflow/models"
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

func cleanupTransactions(ctx context.Context) context.Context {
	query := `{ "query": "query { listTransactions(from:\"1999-01-01T00:00:00.000Z\", to:\"3000-01-01T00:00:00.000Z\") { id date amount balance } }" }`
	var result []models.ComputedTransaction
	PostGraphQL(ctx.Value(url).(string), query, "listTransactions", &result)

	if len(result) != 0 {
		ids := make([]string, len(result))
		for i, transaction := range result {
			if strings.Contains(transaction.Id, "-") {
				ids[i] = strings.Split(transaction.Id, "-")[0]
			} else {
				ids[i] = transaction.Id
			}
		}
		jsonBytes, _ := json.Marshal(ids)

		query := fmt.Sprintf(`{"query": "mutation { deleteTransactions(ids: %s) }"}`, strings.Replace(string(jsonBytes), `"`, `\"`, -1))
		PostGraphQL(ctx.Value(url).(string), query, "deleteTransactions", nil)
	}

	return context.WithValue(ctx, transactions, nil)
}

func cleanupBalances(ctx context.Context) context.Context {
	query := `{ "query": "query { listBalances(from:\"1999-01-01T00:00:00.000Z\", to:\"3000-01-01T00:00:00.000Z\") { date amount } }" }`
	var result []models.Balance
	PostGraphQL(ctx.Value(url).(string), query, "listBalances", &result)

	if len(result) != 0 {
		for _, b := range result {
			query := fmt.Sprintf(`{"query": "mutation { deleteBalance(date: \"%s\") }"}`, b.Date)
			PostGraphQL(ctx.Value(url).(string), query, "deleteBalance", nil)
		}
	}

	return context.WithValue(ctx, balances, nil)
}
