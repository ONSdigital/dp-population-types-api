package steps

import (
	"context"
	"net/http"

	"github.com/pkg/errors"

	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-population-types-api/config"
	"github.com/ONSdigital/dp-population-types-api/service"
	"github.com/ONSdigital/dp-population-types-api/service/mock"
)

type PopulationTypesComponent struct {
	componenttest.ErrorFeature
	svc                          *service.Service
	errorChan                    chan error
	Config                       *config.Config
	ServiceRunning               bool
	HttpServer                   *http.Server
	apiFeature                   *componenttest.APIFeature
	fakeCantabularDatasets       []string
	fakeCantabularIsUnresponsive bool
	service                      *service.Service
	InitialiserMock              service.Initialiser
}

func NewComponent() (*PopulationTypesComponent, error) {

	config, err := config.Get()
	if err != nil {
		return nil, err
	}

	c := &PopulationTypesComponent{
		errorChan:      make(chan error),
		ServiceRunning: false,
		Config:         config,
		HttpServer:     &http.Server{},
	}

	c.apiFeature = componenttest.NewAPIFeature(c.InitialiseService)

	return c, nil
}

func (c *PopulationTypesComponent) Reset() *PopulationTypesComponent {
	c.apiFeature.Reset()
	return c
}

func (c *PopulationTypesComponent) Close() error {
	if c.svc != nil && c.ServiceRunning {
		c.svc.Close(context.Background())
		c.ServiceRunning = false
	}
	return nil
}

func (c *PopulationTypesComponent) InitialiseService() (http.Handler, error) {

	c.InitialiserMock = &mock.InitialiserMock{
		GetHealthCheckFunc:      c.GetHealthcheck,
		GetHTTPServerFunc:       c.GetHttpServer,
		GetCantabularClientFunc: c.GetCantabularClient,
		GetResponderFunc:        c.GetResponder,
	}

	var err error

	c.service = service.New()

	err = c.service.Init(context.Background(), c.InitialiserMock, c.Config, "", "1", "")
	if err != nil {
		return nil, errors.Wrap(err, "error initialising service")
	}

	c.service.Start(context.Background(), c.errorChan)

	c.ServiceRunning = true
	return c.HttpServer.Handler, nil
}

func (c *PopulationTypesComponent) GetResponder() service.Responder {
	return nil
}

func (c *PopulationTypesComponent) GetHttpServer(addr string, router http.Handler) service.HTTPServer {
	c.HttpServer.Addr = addr
	c.HttpServer.Handler = router
	return c.HttpServer
}

func (c *PopulationTypesComponent) GetHealthcheck(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
	return &mock.HealthCheckerMock{
		AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
		StartFunc:    func(ctx context.Context) {},
		StopFunc:     func() {},
	}, nil
}

func (c *PopulationTypesComponent) GetCantabularClient(cfg config.CantabularConfig) service.CantabularClient {
	return nil
}
