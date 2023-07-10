package lib

import (
	"context"
	"errors"
	"fmt"
	"github.com/chxx-v/test/tools/cucumber"
	"github.com/cucumber/godog"
	"testing"
)

// godogsCtxKey is the key used to store the available godogs in the context.Context.
type godogsCtxKey struct{}

func thereAreGodogs(ctx context.Context, available int) (context.Context, error) {
	return context.WithValue(ctx, godogsCtxKey{}, available), nil
}

func iEat(ctx context.Context, num int) (context.Context, error) {
	available, ok := ctx.Value(godogsCtxKey{}).(int)
	if !ok {
		return ctx, errors.New("there are no godogs available")
	}

	if available < num {
		return ctx, fmt.Errorf("you cannot eat %d godogs, there are %d available", num, available)
	}

	available -= num

	return context.WithValue(ctx, godogsCtxKey{}, available), nil
}

func iEarn(ctx context.Context, num int) (context.Context, error) {
	available, _ := ctx.Value(godogsCtxKey{}).(int)

	available += num

	return context.WithValue(ctx, godogsCtxKey{}, available), nil
}

func thereShouldBeRemaining(ctx context.Context, remaining int) error {
	available, ok := ctx.Value(godogsCtxKey{}).(int)
	if !ok {
		return errors.New("there are no godogs available")
	}

	if available != remaining {
		return fmt.Errorf("expected %d godogs to be remaining, but there is %d", remaining, available)
	}

	return nil
}

func thereShouldBeTotal(ctx context.Context, total int) error {
	available, ok := ctx.Value(godogsCtxKey{}).(int)
	if !ok {
		return errors.New("there are no godogs available")
	}

	if available != total {
		return fmt.Errorf("expected %d godogs to be total, but there is %d", total, available)
	}

	return nil
}

func TestFeatures(t *testing.T) {
	opts := &godog.Options{
		Format: "pretty",
		Paths:  []string{"features"},
	}
	mdFiles, err := cucumber.MDFinder("features")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(mdFiles)
	if len(mdFiles) > 0 {
		f, err := cucumber.ExtractGherkinSnippetFromMarkdown(mdFiles)
		if err != nil {
			t.Fatal(err)
		}
		// Setting FeatureContents lets godog to ignore Path and loading the feature files.
		opts.FeatureContents = f
	}
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options:             opts,
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}

func InitializeScenario(sc *godog.ScenarioContext) {
	sc.Step(`^there are (\d+) godogs$`, thereAreGodogs)
	sc.Step(`^I eat (\d+)$`, iEat)
	sc.Step(`^there should be (\d+) remaining$`, thereShouldBeRemaining)
	sc.Step(`^I earn (\d+)$`, iEarn)
	sc.Step(`^there should be (\d+) total$`, thereShouldBeTotal)
}
