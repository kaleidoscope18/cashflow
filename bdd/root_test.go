package bdd

import (
	"context"
	"os"
	"testing"
	"time"

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

func TestBDD(t *testing.T) {
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
		return ctx, nil
	})

	InitializeTransactionsScenarioStepDefs(ctx)
	InitializeBalancesScenarioStepDefs(ctx)
}
