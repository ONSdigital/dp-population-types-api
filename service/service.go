package service

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"

	"github.com/ONSdigital/dp-api-clients-go/v2/identity"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-population-types-api/config"
	"github.com/ONSdigital/log.go/v2/log"
)

// Service contains all the configs, server and clients to run the API
type Service struct {
	Config           *config.Config
	Server           HTTPServer
	Router           *chi.Mux
	responder        Responder
	cantabularClient CantabularClient
	datasetAPIClient DatasetAPIClient
	mongoClient      MongoClient
	HealthCheck      HealthChecker
	identityClient   *identity.Client
}

func New() *Service {
	return &Service{}
}

func (svc *Service) Init(ctx context.Context, init Initialiser, cfg *config.Config, buildTime, gitCommit, version string) error {
	var err error

	if cfg == nil {
		return errors.New("nil config passed to service init")
	}
	log.Info(ctx, "initialising service with config", log.Data{"config": cfg})

	svc.Config = cfg
	svc.identityClient = identity.New(cfg.ZebedeeURL)
	svc.HealthCheck, err = init.GetHealthCheck(cfg, buildTime, gitCommit, version)
	if err != nil {
		return errors.Wrap(err, "failed to get healthcheck")
	}

	svc.responder = init.GetResponder()
	svc.cantabularClient = init.GetCantabularClient(cfg.CantabularConfig)
	svc.datasetAPIClient = init.GetDatasetAPIClient(cfg)
	svc.mongoClient, err = init.GetMongoClient(ctx, cfg)
	if err != nil {
		return errors.Wrap(err, "failed to get mongo client")
	}

	svc.buildRoutes(ctx)

	if cfg.OtelEnabled {
		svc.Server = init.GetHTTPServerWithOtel(cfg.BindAddr, svc.Router)
	} else {
		svc.Server = init.GetHTTPServer(cfg.BindAddr, svc.Router)
	}

	if err = svc.registerCheckers(ctx); err != nil {
		return errors.Wrap(err, "unable to register checkers")
	}

	return nil
}

// Start the service
func (svc *Service) Start(ctx context.Context, svcErrors chan error) {
	svc.HealthCheck.Start(ctx)

	go func() {
		if err := svc.Server.ListenAndServe(); err != nil {
			svcErrors <- errors.Wrap(err, "failed to start main http server")
		}
	}()
}

// Close gracefully shuts the service down in the required order, with timeout
func (svc *Service) Close(ctx context.Context) error {
	timeout := svc.Config.GracefulShutdownTimeout
	log.Info(ctx, "commencing graceful shutdown", log.Data{"graceful_shutdown_timeout": timeout})
	ctx, cancel := context.WithTimeout(ctx, timeout)

	// track shutdown gracefully closes up
	var hasShutdownError bool

	go func() {
		defer cancel()

		// stop healthcheck, as it depends on everything else
		if svc.HealthCheck != nil {
			svc.HealthCheck.Stop()
		}

		// stop any incoming requests before closing any outbound connections
		if err := svc.Server.Shutdown(ctx); err != nil {
			log.Error(ctx, "failed to shutdown http server", err)
			hasShutdownError = true
		}

	}()

	// wait for shutdown success (via cancel) or failure (timeout)
	<-ctx.Done()

	// timeout expired
	if ctx.Err() == context.DeadlineExceeded {
		return errors.Wrap(ctx.Err(), "shutdown timed out")
	}

	// other error
	if hasShutdownError {
		return errors.New("failed to shutdown gracefully")
	}

	log.Info(ctx, "graceful shutdown was successful")
	return nil
}

func (svc *Service) registerCheckers(ctx context.Context) error {
	cantabularSrvChecker := svc.cantabularClient.Checker
	cantabularAPIExtChecker := svc.cantabularClient.CheckerAPIExt

	if !svc.Config.CantabularHealthcheckEnabled {
		cantabularSrvChecker = func(ctx context.Context, state *healthcheck.CheckState) error {
			return state.Update(healthcheck.StatusOK, "Cantabular Server healthcheck placeholder", http.StatusOK)
		}
		cantabularAPIExtChecker = func(ctx context.Context, state *healthcheck.CheckState) error {
			return state.Update(healthcheck.StatusOK, "Cantabular API Ext healthcheck placeholder", http.StatusOK)
		}
	}

	if err := svc.HealthCheck.AddCheck("Cantabular Server", cantabularSrvChecker); err != nil {
		return errors.Wrap(err, "error adding check for cantabular server")
	}

	if err := svc.HealthCheck.AddCheck("Cantabular Api Ext", cantabularAPIExtChecker); err != nil {
		return errors.Wrap(err, "error adding check for cantabular api ext")
	}

	if err := svc.HealthCheck.AddCheck("dataset api", svc.datasetAPIClient.Checker); err != nil {
		return errors.Wrap(err, "error adding check for dataset api client")
	}

	return nil
}
