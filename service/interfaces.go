package service

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-population-types-api/config"
)

//go:generate moq -out mock/initialiser.go -pkg mock . Initialiser
//go:generate moq -out mock/server.go -pkg mock . HTTPServer
//go:generate moq -out mock/healthCheck.go -pkg mock . HealthChecker
//go:generate moq -out mock/cantabularClient.go -pkg mock . CantabularClient

// Initialiser defines the methods to initialise external services
type Initialiser interface {
	GetHTTPServer(bindAddr string, router http.Handler) HTTPServer
	GetCantabularClient(cfg config.CantabularConfig) CantabularClient
	GetHealthCheck(cfg *config.Config, time, commit, version string) (HealthChecker, error)
	GetResponder() Responder
}

// HTTPServer defines the required methods from the HTTP server
type HTTPServer interface {
	ListenAndServe() error
	Shutdown(ctx context.Context) error
}

// HealthChecker defines the required methods from Healthcheck
type HealthChecker interface {
	Handler(w http.ResponseWriter, req *http.Request)
	Start(ctx context.Context)
	Stop()
	AddCheck(name string, checker healthcheck.Checker) (err error)
}

// CantabularClient fetches lists of datasets
type CantabularClient interface {
	ListDatasets(ctx context.Context) ([]string, error)
	Checker(ctx context.Context, state *healthcheck.CheckState) error
}

// Responder handles responding to http requests
type Responder interface {
	JSON(context.Context, http.ResponseWriter, int, interface{})
	Error(context.Context, http.ResponseWriter, int, error)
	StatusCode(http.ResponseWriter, int)
	Bytes(context.Context, http.ResponseWriter, int, []byte)
}
