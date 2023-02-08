package steps

import (
	"context"
	"net/http"
	"testing"

	"github.com/maxcnunes/httpfake"
	"github.com/pkg/errors"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-net/v2/responder"
	"github.com/ONSdigital/dp-population-types-api/config"
	"github.com/ONSdigital/dp-population-types-api/datastore"
	"github.com/ONSdigital/dp-population-types-api/features/mock"
	"github.com/ONSdigital/dp-population-types-api/service"
	svcmock "github.com/ONSdigital/dp-population-types-api/service/mock"
)

const fakeCantabularFailedToRespondErrorMessage = "cantabular failed to respond"

type PopulationTypesComponent struct {
	componenttest.ErrorFeature
	svc                          *service.Service
	errorChan                    chan error
	Config                       *config.Config
	ServiceRunning               bool
	HttpServer                   *http.Server
	fakeCantabularIsUnresponsive bool
	service                      *service.Service
	InitialiserMock              service.Initialiser
	MongoFeature                 *MongoFeature
	datasetAPI                   *httpfake.HTTPFake
	CantabularApiExt             *httpfake.HTTPFake
	CantabularSrv                *httpfake.HTTPFake
	fakeCantabular               *mock.CantabularClient
}

func NewComponent(t testing.TB, zebedeeURL string) (*PopulationTypesComponent, error) {
	config, err := config.Get()
	if err != nil {
		return nil, err
	}

	config.EnablePermissionsAuth = false

	config.ZebedeeURL = zebedeeURL
	c := &PopulationTypesComponent{
		errorChan:        make(chan error),
		ServiceRunning:   false,
		Config:           config,
		HttpServer:       &http.Server{},
		datasetAPI:       httpfake.New(),
		CantabularSrv:    httpfake.New(),
		CantabularApiExt: httpfake.New(httpfake.WithTesting(t)),
	}

	c.datasetAPI = httpfake.New()
	c.Config.DatasetAPIURL = c.datasetAPI.ResolveURL("")
	c.fakeCantabular = &mock.CantabularClient{
		Healthy:    true,
		BadRequest: false,
		NotFound:   false,
	}

	c.MongoFeature = NewMongoFeature(c.ErrorFeature, config)

	return c, nil
}

func (c *PopulationTypesComponent) Reset() error {
	c.datasetAPI.Reset()
	if _, err := c.InitialiseService(); err != nil {
		return errors.Wrap(err, "failed to reset component.")
	}

	return nil
}

func (c *PopulationTypesComponent) Close() error {
	if c.svc != nil && c.ServiceRunning {
		c.svc.Close(context.Background())
		c.ServiceRunning = false
	}
	return nil
}

func (c *PopulationTypesComponent) InitialiseService() (http.Handler, error) {
	ctx := context.Background()

	c.InitialiserMock = &svcmock.InitialiserMock{
		GetHealthCheckFunc:      c.GetHealthcheck,
		GetHTTPServerFunc:       c.GetHttpServer,
		GetCantabularClientFunc: c.GetCantabularClient,
		GetResponderFunc:        c.GetResponder,
		GetDatasetAPIClientFunc: c.GetDatasetAPIClient,
		GetMongoClientFunc:      c.GetMongoClient,
	}

	c.service = service.New()
	c.Config.DatasetAPIURL = c.datasetAPI.ResolveURL("")

	if err := c.service.Init(ctx, c.InitialiserMock, c.Config, "", "1", ""); err != nil {
		return nil, errors.Wrap(err, "error initialising service")
	}

	c.service.Start(ctx, c.errorChan)
	c.ServiceRunning = true

	return c.HttpServer.Handler, nil
}

func (c *PopulationTypesComponent) GetResponder() service.Responder {
	return responder.New()
}

func (c *PopulationTypesComponent) GetHttpServer(addr string, router http.Handler) service.HTTPServer {
	c.HttpServer.Addr = addr
	c.HttpServer.Handler = router
	return c.HttpServer
}

func (c *PopulationTypesComponent) GetHealthcheck(cfg *config.Config, buildTime string, gitCommit string, version string) (service.HealthChecker, error) {
	return &svcmock.HealthCheckerMock{
		AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
		StartFunc:    func(ctx context.Context) {},
		StopFunc:     func() {},
	}, nil
}

func (c *PopulationTypesComponent) GetCantabularClient(_ config.CantabularConfig) service.CantabularClient {
	return c.fakeCantabular
}

func (c *PopulationTypesComponent) GetMongoClient(_ context.Context, _ *config.Config) (service.MongoClient, error) {
	return datastore.NewClient(context.Background(), datastore.Config{
		MongoDriverConfig:  c.Config.Mongo,
		MetadataCollection: "filterMetadata",
	})
}
func (c *PopulationTypesComponent) GetDatasetAPIClient(_ *config.Config) service.DatasetAPIClient {
	return dataset.NewAPIClient(c.Config.DatasetAPIURL)
}
