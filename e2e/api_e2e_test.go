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

	t.Run("seed", func(t *testing.T) {
		var results struct {
			Balance1        models.Balance
			Balance2        models.Balance
			BalanceToDelete models.Balance
			Transactions    []string
			DeleteBalance   string
		}
		client.MustPost(`mutation {
			balance1: 		createBalance(input: {amount: 1000, date: "october 15, 2022"}) {
								date
								amount
							}
			balance2: 		createBalance(input: {amount: 2000, date: "november 03, 2022"}) {
								date
								amount
							}
			transactions: 	createTransactions(input: [
								{amount: -115, date: "october 27, 2022"}, 
								{amount: -117, date: "november 1, 2022"}, 
								{amount: -1333, date: "november 1, 2022"}, 
								{amount: -91, date: "november 4, 2022"}, 
								{amount: 2800.69, date: "november 15, 2022"}, 
								{amount: -1374, date: "november 16, 2022"}
							])
		  }`, &results)

		require.Equal(t, "2022/10/15", results.Balance1.Date)
		require.Equal(t, 1000.00, results.Balance1.Amount)

		require.Equal(t, "2022/11/03", results.Balance2.Date)
		require.Equal(t, 2000.00, results.Balance2.Amount)

		require.Equal(t, 6, len(results.Transactions))

		var listTransactionsResult struct {
			ListTransactions []struct {
				Id      string
				Amount  float64
				Balance float64
				Date    string
			}
		}

		client.MustPost(`query {
			listTransactions(from: "2020-01-01T00:00:00.000Z", to: "2023-01-01T00:00:00.000Z") {
				id
				amount
				balance
				date
			}
		}`, &listTransactionsResult)

		require.Equal(t, "2022/10/27", listTransactionsResult.ListTransactions[0].Date)
		require.Equal(t, 885.00, listTransactionsResult.ListTransactions[0].Balance)

		require.Equal(t, "2022/11/01", listTransactionsResult.ListTransactions[1].Date)
		require.Equal(t, 768.00, listTransactionsResult.ListTransactions[1].Balance)

		require.Equal(t, "2022/11/01", listTransactionsResult.ListTransactions[2].Date)
		require.Equal(t, -565.00, listTransactionsResult.ListTransactions[2].Balance)

		require.Equal(t, "2022/11/04", listTransactionsResult.ListTransactions[3].Date)
		require.Equal(t, 1909.00, listTransactionsResult.ListTransactions[3].Balance)

		require.Equal(t, "2022/11/15", listTransactionsResult.ListTransactions[4].Date)
		require.Equal(t, 4709.69, listTransactionsResult.ListTransactions[4].Balance)

		require.Equal(t, "2022/11/16", listTransactionsResult.ListTransactions[5].Date)
		require.Equal(t, 3335.69, listTransactionsResult.ListTransactions[5].Balance)

		var listBalancesResult struct {
			ListBalances []models.Balance
		}

		client.MustPost(`query {
			listBalances(from: "2020-01-01T00:00:00.000Z", to: "2023-01-01T00:00:00.000Z") {
				amount
				date
			}
		}`, &listBalancesResult)

		require.Equal(t, "2022/10/15", listBalancesResult.ListBalances[0].Date)
		require.Equal(t, 1000.00, listBalancesResult.ListBalances[0].Amount)

		require.Equal(t, "2022/11/03", listBalancesResult.ListBalances[1].Date)
		require.Equal(t, 2000.00, listBalancesResult.ListBalances[1].Amount)
	})

	t.Run("delete all", func(t *testing.T) {
		var response struct {
			DeleteBalanceOne   string
			DeleteBalanceTwo   string
			DeleteTransactions []string
		}
		client.MustPost(`mutation {
			deleteBalanceOne: 	deleteBalance(date: "2022/10/15")
			deleteBalanceTwo: 	deleteBalance(date: "2022/11/03")
								deleteTransactions(ids:["1","2","3","4","5","6"])
		}`, &response)

		var results struct {
			ListBalances []struct {
				Date string
			}
			ListTransactions []struct {
				Id   string
				Date string
			}
		}

		client.MustPost(`query {
			listBalances(from: "2021-01-01T00:00:00.000Z", to: "2023-01-01T00:00:00.000Z") {
				date
			}
			listTransactions(from: "2021-01-01T00:00:00.000Z", to: "2023-01-01T00:00:00.000Z") {
				id
				date
			}
		}`, &results)

		require.Empty(t, results.ListBalances)
		require.Empty(t, results.ListTransactions)
	})
}
