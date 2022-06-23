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
	ListDatasets(ctx context.Context) ([]string, error)
	GetGeographyDimensions(ctx context.Context, req cantabular.GetGeographyDimensionsRequest) (*cantabular.GetGeographyDimensionsResponse, error)
	GetAreas(context.Context, cantabular.GetAreasRequest) (*cantabular.GetAreasResponse, error)
	StatusCode(error) int
}

type datasetAPIClient interface {
	GetDatasets(ctx context.Context, uToken, svcToken, collectionID string, params *dataset.QueryParams) (dataset.List, error)
	Checker(context.Context, *healthcheck.CheckState) error
}

type validator interface {
	Valid() error
}
