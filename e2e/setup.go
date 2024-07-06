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
	err := repository.Init(models.InMemory)
	if err != nil {
		panic(err.Error())
	}
	defer repository.Close()

	r := repository.Get()
	bs := domain.NewBalanceService(r)
	ts := domain.NewTransactionService(r, &bs)

	app := &models.App{
		TransactionService: &ts,
		BalanceService:     &bs,
	}

	server := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{App: app}}))
	client := client.New(server)
	return client
}
