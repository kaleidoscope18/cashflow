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
)

func TestBDD(t *testing.T) {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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
		graphURL := fmt.Sprintf("http://%s:%s/graphql", os.Getenv("APP_HOST"), os.Getenv("APP_PORT"))
		return context.WithValue(ctx, url, graphURL), nil
	})

	InitializeTransactionsScenarioStepDefs(ctx)
	InitializeBalancesScenarioStepDefs(ctx)
}
