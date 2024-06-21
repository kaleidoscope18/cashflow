package e2e

import (
	"cashflow/models"
	"testing"

	"github.com/99designs/gqlgen/graphql/introspection"
	"github.com/stretchr/testify/require"
)

func TestGraphQLApiE2E(t *testing.T) {
	client := Initialize()

	t.Run("graphql introspection query", func(t *testing.T) {
		var resp interface{}
		client.MustPost(introspection.Query, &resp)
	})

	t.Run("mutations", func(t *testing.T) {
		var results struct {
			Balance1     models.Balance
			Balance2     models.Balance
			Transactions []models.Transaction
		}
		client.MustPost(`mutation {
			balance1:createBalance(input:{amount:1000,date:"october 15, 2022"}) {
				date
				amount
			}
			balance2:createBalance(input:{amount:2000,date:"november 03, 2022"}) {
				date
				amount
			}
			transactions:createTransactions(
				input: [
				  {amount: -115, date: "october 27, 2022"},
				  {amount: -117, date: "november 1, 2022"},
				  {amount: -1333, date: "november 1, 2022"},
				  {amount: -91, date: "november 4, 2022"},
				  {amount: 2800.69, date: "november 15, 2022"},
				  {amount: -1374, date: "november 16, 2022"},
				]
			  ) {
				date
				amount
			  }
		}`, &results)

		require.Equal(t, "2022/10/15", results.Balance1.Date)
		require.Equal(t, 1000.00, results.Balance1.Amount)

		require.Equal(t, "2022/11/03", results.Balance2.Date)
		require.Equal(t, 2000.00, results.Balance2.Amount)

		require.Equal(t, "2022/10/27", results.Transactions[0].Date)
		require.Equal(t, -115.00, results.Transactions[0].Amount)

		require.Equal(t, "2022/11/01", results.Transactions[1].Date)
		require.Equal(t, -117.00, results.Transactions[1].Amount)

		require.Equal(t, "2022/11/01", results.Transactions[2].Date)
		require.Equal(t, -1333.00, results.Transactions[2].Amount)

		require.Equal(t, "2022/11/04", results.Transactions[3].Date)
		require.Equal(t, -91.00, results.Transactions[3].Amount)

		require.Equal(t, "2022/11/15", results.Transactions[4].Date)
		require.Equal(t, 2800.69, results.Transactions[4].Amount)

		require.Equal(t, "2022/11/16", results.Transactions[5].Date)
		require.Equal(t, -1374.00, results.Transactions[5].Amount)
	})

	t.Run("query transactions with balance", func(t *testing.T) {
		var results struct {
			ListTransactions []struct {
				Id      string
				Amount  float64
				Balance float64
				Date    string
			}
		}

		client.MustPost(`query {
			listTransactions(from: "2000-01-01T00:00:00.000Z", to: "2020-01-01T00:00:00.000Z") {
				id
				amount
				balance
				date
			}
		}`, &results)

		require.Equal(t, "2022/10/27", results.ListTransactions[0].Date)
		require.Equal(t, 885.00, results.ListTransactions[0].Balance)

		require.Equal(t, "2022/11/01", results.ListTransactions[1].Date)
		require.Equal(t, 768.00, results.ListTransactions[1].Balance)

		require.Equal(t, "2022/11/01", results.ListTransactions[2].Date)
		require.Equal(t, -565.00, results.ListTransactions[2].Balance)

		require.Equal(t, "2022/11/04", results.ListTransactions[3].Date)
		require.Equal(t, 1909.00, results.ListTransactions[3].Balance)

		require.Equal(t, "2022/11/15", results.ListTransactions[4].Date)
		require.Equal(t, 4709.69, results.ListTransactions[4].Balance)

		require.Equal(t, "2022/11/16", results.ListTransactions[5].Date)
		require.Equal(t, 3335.69, results.ListTransactions[5].Balance)
	})

	t.Run("query balances", func(t *testing.T) {
		var results struct {
			ListBalances []models.Balance
		}

		client.MustPost(`query {
			listBalances(from: "2000-01-01T00:00:00.000Z", to: "2020-01-01T00:00:00.000Z") {
				amount
				date
			}
		}`, &results)

		require.Equal(t, "2022/10/15", results.ListBalances[0].Date)
		require.Equal(t, 1000.00, results.ListBalances[0].Amount)

		require.Equal(t, "2022/11/03", results.ListBalances[1].Date)
		require.Equal(t, 2000.00, results.ListBalances[1].Amount)
	})
}
