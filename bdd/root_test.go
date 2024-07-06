package bdd

import (
	"cashflow/api/graph"
	"cashflow/api/graph/generated"
	"cashflow/domain"
	"cashflow/models"
	"cashflow/repository"
	"context"
	"os"
	"testing"
	"time"

	"github.com/99designs/gqlgen/client"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/introspection"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
)

var options = godog.Options{
	Output:    colors.Colored(os.Stdout),
	Format:    "pretty",
	Randomize: time.Now().UTC().UnixNano(),
}

func init() {
	godog.BindCommandLineFlags("godog.", &options)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenarios,
		Options: &godog.Options{
			Output:    colors.Colored(os.Stdout),
			Randomize: time.Now().UTC().UnixNano(),
			Format:    "pretty",
			Paths:     []string{"features"},
			TestingT:  t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeScenarios(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		InitializeClient()
		return ctx, nil
	})

	InitializeTransactionsScenarioStepDefs(ctx)
}

func InitializeClient() *client.Client {
	err := repository.Init(models.InMemory)
	if err != nil {
		panic(err.Error())
	}

	tr, br := repository.GetRepos()
	bs := domain.NewBalanceService(br)
	ts := domain.NewTransactionService(tr, &bs)

	app := &models.App{
		TransactionService: &ts,
		BalanceService:     &bs,
	}

	server := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{App: app}}))
	client := client.New(server)

	// verify if graphql introspection is good
	var resp interface{}
	client.MustPost(introspection.Query, &resp)

	return client
}

func CleanupClient() {
	repository.Close()
}
