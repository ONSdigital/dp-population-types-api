package service

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular/gql"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/stream"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-population-types-api/config"
	"github.com/ONSdigital/dp-population-types-api/datastore"
)

//go:generate moq -out mock/initialiser.go -pkg mock . Initialiser
//go:generate moq -out mock/server.go -pkg mock . HTTPServer
//go:generate moq -out mock/health_check.go -pkg mock . HealthChecker
//go:generate moq -out mock/cantabular_client.go -pkg mock . CantabularClient
//go:generate moq -out mock/dataset_api_client.go -pkg mock . DatasetAPIClient
//go:generate moq -out mock/health_check.go -pkg mock . HealthChecker

// Initialiser defines the methods to initialise external services
type Initialiser interface {
	GetHTTPServer(bindAddr string, router http.Handler) HTTPServer
	GetCantabularClient(cfg config.CantabularConfig) CantabularClient
	GetHealthCheck(cfg *config.Config, time, commit, version string) (HealthChecker, error)
	GetDatasetAPIClient(cfg *config.Config) DatasetAPIClient
	GetMongoClient(ctx context.Context, cfg *config.Config) (MongoClient, error)
	GetResponder() Responder
	GetHTTPServerWithOtel(bindAddr string, router http.Handler) HTTPServer
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

type CantabularClient interface {
	ListDatasets(context.Context) (*cantabular.ListDatasetsResponse, error)
	GetDimensions(context.Context, cantabular.GetDimensionsRequest) (*cantabular.GetDimensionsResponse, error)
	GetDimensionsDescription(context.Context, cantabular.GetDimensionsDescriptionRequest) (*cantabular.GetDimensionsResponse, error)
	GetDimensionCategories(context.Context, cantabular.GetDimensionCategoriesRequest) (*cantabular.GetDimensionCategoriesResponse, error)
	GetGeographyDimensions(ctx context.Context, req cantabular.GetGeographyDimensionsRequest) (*cantabular.GetGeographyDimensionsResponse, error)
	GetAreas(context.Context, cantabular.GetAreasRequest) (*cantabular.GetAreasResponse, error)
	GetAreasTotalCount(context.Context, cantabular.GetAreasRequest) (int, error)
	GetArea(context.Context, cantabular.GetAreaRequest) (*cantabular.GetAreaResponse, error)
	GetParents(context.Context, cantabular.GetParentsRequest) (*cantabular.GetParentsResponse, error)
	GetParentAreaCount(ctx context.Context, req cantabular.GetParentAreaCountRequest) (*cantabular.GetParentAreaCountResult, error)
	GetBlockedAreaCount(ctx context.Context, req cantabular.GetBlockedAreaCountRequest) (*cantabular.GetBlockedAreaCountResult, error)
	GetCategorisations(context.Context, cantabular.GetCategorisationsRequest) (*cantabular.GetCategorisationsResponse, error)
	GetBaseVariable(context.Context, cantabular.GetBaseVariableRequest) (*cantabular.GetBaseVariableResponse, error)
	Checker(ctx context.Context, state *healthcheck.CheckState) error
	CheckerAPIExt(ctx context.Context, state *healthcheck.CheckState) error
	StatusCode(error) int
	StaticDatasetQuery(context.Context, cantabular.StaticDatasetQueryRequest) (*cantabular.StaticDatasetQuery, error)
	StaticDatasetType(ctx context.Context, datasetName string) (*gql.Dataset, error)
	StaticDatasetQueryStreamJSON(context.Context, cantabular.StaticDatasetQueryRequest, stream.Consumer) (cantabular.GetObservationsResponse, error)
	CheckQueryCount(context.Context, cantabular.StaticDatasetQueryRequest) (int, error)
}

// Responder handles responding to http requests
type Responder interface {
	JSON(context.Context, http.ResponseWriter, int, interface{})
	Error(context.Context, http.ResponseWriter, int, error)
	StatusCode(http.ResponseWriter, int)
	Bytes(context.Context, http.ResponseWriter, int, []byte)
}
type DatasetAPIClient interface {
	GetDatasets(ctx context.Context, uToken, svcToken, collectionID string, params *dataset.QueryParams) (dataset.List, error)
	Checker(ctx context.Context, state *healthcheck.CheckState) error
}

type MongoClient interface {
	GetDefaultDatasetMetadata(ctx context.Context, populationType string) (*datastore.DefaultDatasetMetadata, error)
	PutDefaultDatasetMetadata(ctx context.Context, metadata datastore.DefaultDatasetMetadata) error
	GetDefaultDatasetPopulationTypes(ctx context.Context) ([]string, error)
}
