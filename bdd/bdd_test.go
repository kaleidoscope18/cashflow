package bdd

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/joho/godotenv"
	"github.com/spf13/pflag"
)

var GraphURL string
var UserToken string

var opts = godog.Options{
	Output:    colors.Colored(os.Stdout),
	Randomize: time.Now().UTC().UnixNano(),
	Format:    "pretty",
	Paths:     []string{"."}, // instead of []string{"features"} to avoid nested folders
	Strict:    true,
}

func TestBDD(t *testing.T) {
	pflag.Parse()
	args := pflag.Args()

	// run command with i.e. -- --tags=only,wip
	for i := 0; i < len(args); i++ {
		if args[i] == "--tags" && i+1 < len(args) {
			opts.Tags = args[i+1]
			break
		}
	}

	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	opts.TestingT = t

	suite := godog.TestSuite{
		Name:                 "Cashflow BDD",
		TestSuiteInitializer: InitializeTestSuite,
		ScenarioInitializer:  InitializeScenarios,
		Options:              &opts,
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeTestSuite(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		log.Println("Initializing test suite")
		GraphURL = fmt.Sprintf("http://%s:%s/graphql", os.Getenv("APP_HOST"), os.Getenv("APP_PORT"))
		log.Println("GraphQL URL:", GraphURL)
		UserToken = LoginBDDUser(GraphURL)
		log.Println("Test suite initialized")
	})
}

func InitializeScenarios(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		ctx = context.WithValue(ctx, url, GraphURL)
		ctx = context.WithValue(ctx, userToken, UserToken)
		return ctx, nil
	})

	InitializeCommonStepDefs(ctx)
	InitializeTransactionsScenarioStepDefs(ctx)
	InitializeBalancesScenarioStepDefs(ctx)

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		ctx = cleanupTransactions(ctx)
		ctx = cleanupBalances(ctx)
		return ctx, nil
	})
}
