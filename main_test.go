package main

import (
	"context"
	"flag"
	"io"
	glog "log"
	"os"
	"testing"

	"github.com/ONSdigital/log.go/v2/log"

	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"

	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-population-types-api/features/steps"
)

var componentFlag = flag.Bool("component", false, "perform component tests")
var loggingFlag = flag.Bool("logging", false, "print logging")

type ComponentTest struct {
	MongoFeature *componenttest.MongoFeature
}

func (f *ComponentTest) InitializeScenario(ctx *godog.ScenarioContext) {
	authFeature := componenttest.NewAuthorizationFeature()
	zebedeeURL := authFeature.FakeAuthService.ResolveURL("")

	component, err := steps.NewComponent(zebedeeURL)
	if err != nil {
		panic(err)
	}
	apiFeature := componenttest.NewAPIFeature(component.InitialiseService)
	component.APIFeature = apiFeature

	ctx.Before(func(ctx context.Context, sc *godog.Scenario) (context.Context, error) {
		apiFeature.Reset()
		if err := component.Reset(); err != nil {
			glog.Panicf("unable to initialise scenario: %s", err)
		}
		authFeature.Reset()

		return ctx, nil
	})

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		afterErr := component.Close()
		return ctx, afterErr
	})

	component.RegisterSteps(ctx)
	component.APIFeature.RegisterSteps(ctx)
	authFeature.RegisterSteps(ctx)
}

func (f *ComponentTest) InitializeTestSuite(ctx *godog.TestSuiteContext) {

}

func TestComponent(t *testing.T) {
	os.Clearenv()
	if *componentFlag {

		if !*loggingFlag {
			// discarding production logging only during the test run, in order to make the BDD output readable
			log.SetDestination(io.Discard, io.Discard)
			defer func() { log.SetDestination(os.Stdout, os.Stderr) }()
		}

		status := 0

		var opts = godog.Options{
			Output: colors.Colored(os.Stdout),
			Format: "pretty",
			Paths:  flag.Args(),
		}

		f := &ComponentTest{}

		status = godog.TestSuite{
			Name:                 "feature_tests",
			ScenarioInitializer:  f.InitializeScenario,
			TestSuiteInitializer: f.InitializeTestSuite,
			Options:              &opts,
		}.Run()

		if status > 0 {
			t.Fail()
		}
	} else {
		t.Skip("component flag required to run component tests")
	}
}
