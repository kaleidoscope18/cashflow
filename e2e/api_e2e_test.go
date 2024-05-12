package e2e

import (
	"cashflow/api/graph"
	"cashflow/api/graph/generated"
	"cashflow/domain/transactions"
	"cashflow/models"
	"cashflow/repository"
	"testing"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/introspection"
	"github.com/stretchr/testify/require"
)

func initialize() *client.Client {
	storage := repository.New("InMemory")
	ts := transactions.New(storage)
	app := &models.App{
		TransactionService: ts,
	}
	client := client.New(handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{App: app}})))
	return client
}

func TestEnd2End(t *testing.T) {
	client := initialize()

	t.Run("introspection", func(t *testing.T) {
		// Make sure we can run the graphiql introspection query without errors
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
				Amount  float64
				Balance float64
				Date    string
			}
		}

		client.MustPost(`query {
			listTransactions {
				amount
				balance
				date
			}
		}`, &results)

		require.Equal(t, "2022/10/27", results.ListTransactions[0].Date)
		require.Equal(t, 885.00, results.ListTransactions[0].Balance)

		require.Equal(t, "2022/11/01", results.ListTransactions[1].Date)
		require.True(t, results.ListTransactions[1].Balance == 768.00 || results.ListTransactions[1].Balance == -448.00) // it depends on the order on same day (2 transactions on 1 date)

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
			ListBalances []struct {
				Amount float64
				Date   string
			}
		}

		client.MustPost(`query {
			listBalances {
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
