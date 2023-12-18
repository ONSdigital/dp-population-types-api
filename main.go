package main

import (
	"context"
	goerrors "errors"
	"os"
	"os/signal"
	"runtime/debug"
	"syscall"

	"github.com/pkg/errors"

	"github.com/ONSdigital/dp-otel-go"
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

	// NOTE: replace the above with the below to run code with for example vscode debugger.
	//BuildTime string = "1601119818"
	//GitCommit string = "6584b786caac36b6214ffe04bf62f058d4021538"
	//Version   string = "v0.1.0"
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
		return errors.Wrap(err, "unable to retrieve configuration, error: %w")
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Set up OpenTelemetry
	otelConfig := dpotelgo.Config{
		OtelServiceName:          cfg.OTServiceName,
		OtelExporterOtlpEndpoint: cfg.OTExporterOTLPEndpoint,
		OtelBatchTimeout:         cfg.OTBatchTimeout,
	}

	otelShutdown, err := dpotelgo.SetupOTelSDK(ctx, otelConfig)

	if err != nil {
		log.Error(ctx, "error setting up OpenTelemetry - hint: ensure OTEL_EXPORTER_OTLP_ENDPOINT is set", err)
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = goerrors.Join(err, otelShutdown(context.Background()))
	}()

	// Start service
	svc := service.New()
	init := service.NewInit()
	if err := svc.Init(ctx, init, cfg, BuildTime, GitCommit, Version); err != nil {
		return errors.Wrap(err, "failed to initialise service")
	}

	// The following value could be read from some a config setting ...
	debug.SetGCPercent(25)

	svc.Start(ctx, svcErrors)

	// blocks until an os interrupt or a fatal error occurs
	select {
	case err := <-svcErrors:
		return errors.Wrap(err, "service error received")
	case sig := <-signals:
		log.Info(ctx, "os signal received", log.Data{"signal": sig})
	}

	if err := svc.Close(ctx); err != nil {
		return errors.Wrap(err, "failed to close service")
	}

	return nil
}
