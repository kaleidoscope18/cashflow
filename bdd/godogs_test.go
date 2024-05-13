package bdd

import (
	"context"
	"os"
	"testing"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/assert"
)

var opts = godog.Options{Output: colors.Colored(os.Stdout)}

func init() {
	godog.BindCommandLineFlags("godog.", &opts)
}

func TestMain(m *testing.M) {
	pflag.Parse()
	opts.Paths = pflag.Args()

	status := godog.TestSuite{
		Name:                "godogs",
		ScenarioInitializer: InitializeScenario,
		Options:             &opts,
	}.Run()

	os.Exit(status)
}

func thereAreGodogs(available int) error {
	NumberOfGodogs = available
	return nil
}

func iEat(ctx context.Context, num int) error {
	if !assert.GreaterOrEqual(godog.T(ctx), NumberOfGodogs, num, "You cannot eat %d godogs, there are %d available", num, NumberOfGodogs) {
		return nil
	}
	NumberOfGodogs -= num
	return nil
}

func thereShouldBeRemaining(ctx context.Context, remaining int) error {
	assert.Equal(godog.T(ctx), NumberOfGodogs, remaining, "Expected %d godogs to be remaining, but there is %d", remaining, NumberOfGodogs)
	return nil
}

func thereShouldBeNoneRemaining(ctx context.Context) error {
	assert.Empty(godog.T(ctx), NumberOfGodogs, "Expected none godogs to be remaining, but there is %d", NumberOfGodogs)
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		NumberOfGodogs = 0 // clean the state before every scenario
		return ctx, nil
	})

	ctx.Step(`^there are (\d+) godogs$`, thereAreGodogs)
	ctx.Step(`^I eat (\d+)$`, iEat)
	ctx.Step(`^there should be (\d+) remaining$`, thereShouldBeRemaining)
	ctx.Step(`^there should be none remaining$`, thereShouldBeNoneRemaining)
}
