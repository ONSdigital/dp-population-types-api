package handler

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular/gql"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/stream"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/dp-population-types-api/datastore"
)

type responder interface {
	JSON(context.Context, http.ResponseWriter, int, interface{})
	Error(context.Context, http.ResponseWriter, int, error)
}
type cantabularClient interface {
	ListDatasets(context.Context) (*cantabular.ListDatasetsResponse, error)
	GetGeographyDimensions(context.Context, cantabular.GetGeographyDimensionsRequest) (*cantabular.GetGeographyDimensionsResponse, error)
	GetDimensions(context.Context, cantabular.GetDimensionsRequest) (*cantabular.GetDimensionsResponse, error)
	GetDimensionsDescription(context.Context, cantabular.GetDimensionsDescriptionRequest) (*cantabular.GetDimensionsResponse, error)
	GetDimensionCategories(context.Context, cantabular.GetDimensionCategoriesRequest) (*cantabular.GetDimensionCategoriesResponse, error)
	GetAreas(context.Context, cantabular.GetAreasRequest) (*cantabular.GetAreasResponse, error)
	GetAreasTotalCount(context.Context, cantabular.GetAreasRequest) (int, error)
	GetArea(context.Context, cantabular.GetAreaRequest) (*cantabular.GetAreaResponse, error)
	GetParents(context.Context, cantabular.GetParentsRequest) (*cantabular.GetParentsResponse, error)
	GetParentAreaCount(context.Context, cantabular.GetParentAreaCountRequest) (*cantabular.GetParentAreaCountResult, error)
	GetBlockedAreaCount(context.Context, cantabular.GetBlockedAreaCountRequest) (*cantabular.GetBlockedAreaCountResult, error)
	GetCategorisations(context.Context, cantabular.GetCategorisationsRequest) (*cantabular.GetCategorisationsResponse, error)
	GetBaseVariable(context.Context, cantabular.GetBaseVariableRequest) (*cantabular.GetBaseVariableResponse, error)
	StatusCode(error) int
	StaticDatasetQuery(context.Context, cantabular.StaticDatasetQueryRequest) (*cantabular.StaticDatasetQuery, error)
	StaticDatasetType(ctx context.Context, datasetName string) (*gql.Dataset, error)
	StaticDatasetQueryStreamJSON(context.Context, cantabular.StaticDatasetQueryRequest, stream.Consumer) (cantabular.GetObservationsResponse, error)
	CheckQueryCount(context.Context, cantabular.StaticDatasetQueryRequest) (int, error)
}

type datasetAPIClient interface {
	GetDatasets(ctx context.Context, uToken, svcToken, collectionID string, params *dataset.QueryParams) (dataset.List, error)
	Checker(context.Context, *healthcheck.CheckState) error
}

type Datastore interface {
	GetDefaultDatasetMetadata(ctx context.Context, populationType string) (*datastore.DefaultDatasetMetadata, error)
	PutDefaultDatasetMetadata(ctx context.Context, metadata datastore.DefaultDatasetMetadata) error
	GetDefaultDatasetPopulationTypes(ctx context.Context) ([]string, error)
}

type validator interface {
	Valid() error
}

type coder interface {
	Code() int
}
