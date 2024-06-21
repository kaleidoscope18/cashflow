package e2e

import (
	"cashflow/api/graph"
	"cashflow/api/graph/generated"
	"cashflow/domain"
	"cashflow/models"
	"cashflow/repository"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
)

func Initialize() *client.Client {
	tr, br, _ := repository.Init(models.InMemory)
	defer repository.Close()

	bs := domain.NewBalanceService(br)
	ts := domain.NewTransactionService(tr, &bs)

	app := &models.App{
		TransactionService: &ts,
		BalanceService:     &bs,
	}

	server := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{App: app}}))
	client := client.New(server)
	return client
}
