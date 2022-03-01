package main

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"os/signal"
	"syscall"

	"github.com/ONSdigital/dp-population-types-api/config"
	"github.com/ONSdigital/dp-population-types-api/service"
	"github.com/ONSdigital/log.go/v2/log"
)

const serviceName = "dp-population-types-api"

var (
	// BuildTime represents the time in which the service was built
	BuildTime string
	// GitCommit represents the commit (SHA-1) hash of the service that is running
	GitCommit string
	// Version represents the version of the service that is running
	Version string

	// TODO: remove below explainer before commiting
	/* NOTE: replace the above with the below to run code with for example vscode debugger.
	   BuildTime string = "1601119818"
	   GitCommit string = "6584b786caac36b6214ffe04bf62f058d4021538"
	   Version   string = "v0.1.0"
	*/
)

func main() {
	log.Namespace = serviceName
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Fatal(ctx, "fatal runtime error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	svcErrors := make(chan error, 1)

	// Read config
	cfg, err := config.Get()
	if err != nil {
		return fmt.Errorf("unable to retrieve configuration, error: %w", err)
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Start service
	svc := service.New()
	init := service.NewInit()
	if err := svc.Init(ctx, init, cfg, BuildTime, GitCommit, Version); err != nil {
		return fmt.Errorf("failed to initialise service: %w", err)
	}
	svc.Start(ctx, svcErrors)

	// blocks until an os interrupt or a fatal error occurs
	select {
	case err := <-svcErrors:
		return errors.Wrap(err, "service error received")
	case sig := <-signals:
		log.Info(ctx, "os signal received", log.Data{"signal": sig})
	}

	if err := svc.Close(ctx); err != nil {
		return fmt.Errorf("failed to close service: %w", err)
	}

	return nil
}
