package handler

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
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
