package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/ONSdigital/dp-population-types-api/config"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/go-chi/chi/v5"
)

// Service contains all the configs, server and clients to run the API
type Service struct {
	Config           *config.Config
	Server           HTTPServer
	router           *chi.Mux
	responder        Responder
	cantabularClient CantabularClient
	HealthCheck      HealthChecker
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

	svc.HealthCheck, err = init.GetHealthCheck(cfg, buildTime, gitCommit, version)
	if err != nil {
		return fmt.Errorf("failed to get healthcheck: %w", err)
	}

	svc.responder = init.GetResponder()
	svc.cantabularClient = init.GetCantabularClient(cfg.CantabularConfig)

	svc.buildRoutes(ctx)
	svc.Server = init.GetHTTPServer(cfg.BindAddr, svc.router)

	if err := svc.registerCheckers(); err != nil {
		return fmt.Errorf("unable to register checkers: %w", err)
	}

	return nil
}

// Start the service
func (svc *Service) Start(ctx context.Context, svcErrors chan error) {
	svc.HealthCheck.Start(ctx)

	go func() {
		if err := svc.Server.ListenAndServe(); err != nil {
			svcErrors <- fmt.Errorf("failed to start main http server: %w", err)
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
		return fmt.Errorf("shutdown timed out: %w", ctx.Err())
	}

	// other error
	if hasShutdownError {
		return errors.New("failed to shutdown gracefully")
	}

	log.Info(ctx, "graceful shutdown was successful")
	return nil
}

func (svc *Service) registerCheckers() (err error) {

	if svc.cantabularClient != nil {
		if err := svc.HealthCheck.AddCheck("Cantabular client", svc.cantabularClient.Checker); err != nil {
			return fmt.Errorf("error adding check for cantabular client")
		}
	}

	return nil
}

//
//// Run the service
//func Run(ctx context.Context, cfg *config.Config, serviceList *ExternalServiceList, buildTime, gitCommit, version string, svcErrors chan error) (*Service, error) {
//
//	log.Info(ctx, "running service")
//
//	log.Info(ctx, "using service configuration", log.Data{"config": cfg})
//
//	// Get HTTP Server and ... // TODO: Add any middleware that your service requires
//	r := mux.NewRouter()
//
//	s := serviceList.GetHTTPServer(cfg.BindAddr, r)
//
//	// TODO: Add other(s) to serviceList here
//
//	// Setup the API
//	a := api.Setup(ctx, r)
//
//	hc, err := serviceList.GetHealthCheck(cfg, buildTime, gitCommit, version)
//
//	if err != nil {
//		log.Fatal(ctx, "could not instantiate healthcheck", err)
//		return nil, err
//	}
//
//	cantabularClient := serviceList.GetCantabularClient(ctx, cfg.CantabularConfig)
//
//	if err := registerCheckers(ctx, hc, cantabularClient); err != nil {
//		return nil, errors.Wrap(err, "unable to register checkers")
//	}
//
//	r.StrictSlash(true).Path("/health").HandlerFunc(hc.Handler)
//	hc.Start(ctx)
//
//	// Run the http server in a new go-routine
//	go func() {
//		if err := s.ListenAndServe(); err != nil {
//			svcErrors <- errors.Wrap(err, "failure in http listen and serve")
//		}
//	}()
//
//	return &Service{
//		Config:           cfg,
//		Router:           r,
//		Api:              a,
//		HealthCheck:      hc,
//		ServiceList:      serviceList,
//		Server:           s,
//		CantabularClient: cantabularClient,
//	}, nil
//}
//
//// Close gracefully shuts the service down in the required order, with timeout
//func (svc *Service) Close(ctx context.Context) error {
//	timeout := svc.Config.GracefulShutdownTimeout
//	log.Info(ctx, "commencing graceful shutdown", log.Data{"graceful_shutdown_timeout": timeout})
//	ctx, cancel := context.WithTimeout(ctx, timeout)
//
//	// track shutown gracefully closes up
//	var hasShutdownError bool
//
//	go func() {
//		defer cancel()
//
//		// stop healthcheck, as it depends on everything else
//		if svc.ServiceList.HealthCheck {
//			svc.HealthCheck.Stop()
//		}
//
//		// stop any incoming requests before closing any outbound connections
//		if err := svc.Server.Shutdown(ctx); err != nil {
//			log.Error(ctx, "failed to shutdown http server", err)
//			hasShutdownError = true
//		}
//
//		// TODO: Close other dependencies, in the expected order
//	}()
//
//	// wait for shutdown success (via cancel) or failure (timeout)
//	<-ctx.Done()
//
//	// timeout expired
//	if ctx.Err() == context.DeadlineExceeded {
//		log.Error(ctx, "shutdown timed out", ctx.Err())
//		return ctx.Err()
//	}
//
//	// other error
//	if hasShutdownError {
//		err := errors.New("failed to shutdown gracefully")
//		log.Error(ctx, "failed to shutdown gracefully ", err)
//		return err
//	}
//
//	log.Info(ctx, "graceful shutdown was successful")
//	return nil
//}
