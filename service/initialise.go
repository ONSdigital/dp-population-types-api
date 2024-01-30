package service

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-population-types-api/config"
	"github.com/ONSdigital/dp-population-types-api/datastore"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dphttp "github.com/ONSdigital/dp-net/http"
	"github.com/ONSdigital/dp-net/v2/responder"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/pkg/errors"
)

type Init struct {
	CantabularClientFactory func(cfg cantabular.Config, ua dphttp.Clienter) *cantabular.Client
}

func NewInit() *Init {
	return &Init{
		CantabularClientFactory: cantabularNewClient,
	}
}

// GetHTTPServer creates an http server
func (i *Init) GetHTTPServer(bindAddr string, router http.Handler) HTTPServer {
	s := dphttp.NewServer(bindAddr, router)
	s.HandleOSSignals = false
	return s
}

// GetHealthCheck creates a healthcheck with versionInfo and sets teh HealthCheck flag to true
func (i *Init) GetHealthCheck(cfg *config.Config, buildTime, gitCommit, version string) (HealthChecker, error) {
	versionInfo, err := healthcheck.NewVersionInfo(buildTime, gitCommit, version)
	if err != nil {
		return nil, errors.Wrap(err, "Healthcheck version info failed")
	}
	hc := healthcheck.New(versionInfo, cfg.HealthCheckCriticalTimeout, cfg.HealthCheckInterval)
	return &hc, nil
}

// GetCantabularClient creates a cantabular client and sets the CantabularClient flag to true
func (i *Init) GetCantabularClient(cfg config.CantabularConfig) CantabularClient {
	return i.CantabularClientFactory(
		cantabular.Config{
			Host:           cfg.CantabularURL,
			ExtApiHost:     cfg.CantabularExtURL,
			GraphQLTimeout: cfg.DefaultRequestTimeout,
		},
		dphttp.ClientWithTimeout(nil, cfg.DefaultRequestTimeout),
	)
}

func (i *Init) GetDatasetAPIClient(cfg *config.Config) DatasetAPIClient {
	return dataset.NewAPIClient(cfg.DatasetAPIURL)
}

func (i *Init) GetMongoClient(ctx context.Context, cfg *config.Config) (MongoClient, error) {
	return datastore.NewClient(ctx, datastore.Config{
		MongoDriverConfig:  cfg.Mongo,
		MetadataCollection: cfg.MetadataCollection,
	})
}
func cantabularNewClient(cfg cantabular.Config, ua dphttp.Clienter) *cantabular.Client {
	return cantabular.NewClient(cfg, ua, nil)
}

func (i *Init) GetResponder() Responder {
	return responder.New()
}

func (i *Init) GetHTTPServerWithOtel(bindAddr string, router http.Handler) HTTPServer {
	otelHandler := otelhttp.NewHandler(router, "/")
	return dphttp.NewServer(bindAddr, otelHandler)
}
