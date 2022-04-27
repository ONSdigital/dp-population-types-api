package service

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	dphttp "github.com/ONSdigital/dp-net/http"
	"github.com/ONSdigital/dp-net/v2/responder"
	"github.com/ONSdigital/dp-population-types-api/config"
)

// Init implements the Initialiser interface to initialise dependencies
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
		dphttp.NewClient(),
	)
}

func cantabularNewClient(cfg cantabular.Config, ua dphttp.Clienter) *cantabular.Client {
	return cantabular.NewClient(cfg, ua, nil)
}

func (i *Init) GetResponder() Responder {
	return responder.New()
}
