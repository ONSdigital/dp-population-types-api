package handler

import (
	"context"
	"net/http"

<<<<<<< HEAD
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
=======
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
>>>>>>> 63d8084 (Add published check ro GET /population-types)
)

type responder interface {
	JSON(context.Context, http.ResponseWriter, int, interface{})
	Error(context.Context, http.ResponseWriter, int, error)
}

// CantabularClient fetches lists of datasets
type cantabularClient interface {
	ListDatasets(ctx context.Context) ([]string, error)
	GetGeographyDimensions(ctx context.Context, req cantabular.GetGeographyDimensionsRequest) (*cantabular.GetGeographyDimensionsResponse, error)
	StatusCode(error) int
}

type datasetAPIClient interface {
	GetDatasets(ctx context.Context, uToken, svcToken, collectionID string, params *dataset.QueryParams) (dataset.List, error)
}
