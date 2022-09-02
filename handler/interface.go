package handler

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
)

type responder interface {
	JSON(context.Context, http.ResponseWriter, int, interface{})
	Error(context.Context, http.ResponseWriter, int, error)
}
type cantabularClient interface {
	ListDatasets(context.Context) ([]string, error)
	GetGeographyDimensions(context.Context, cantabular.GetGeographyDimensionsRequest) (*cantabular.GetGeographyDimensionsResponse, error)
	GetDimensions(context.Context, cantabular.GetDimensionsRequest) (*cantabular.GetDimensionsResponse, error)
	GetAreas(context.Context, cantabular.GetAreasRequest) (*cantabular.GetAreasResponse, error)
	GetParents(context.Context, cantabular.GetParentsRequest) (*cantabular.GetParentsResponse, error)
	StatusCode(error) int
}

type datasetAPIClient interface {
	GetDatasets(ctx context.Context, uToken, svcToken, collectionID string, params *dataset.QueryParams) (dataset.List, error)
	Checker(context.Context, *healthcheck.CheckState) error
}

type validator interface {
	Valid() error
}
