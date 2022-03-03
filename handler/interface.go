package handler

import (
	"context"
	"net/http"
)

type responder interface {
	JSON(context.Context, http.ResponseWriter, int, interface{})
	Error(context.Context, http.ResponseWriter, int, error)
}

// CantabularClient fetches lists of datasets
type cantabularClient interface {
	ListDatasets(ctx context.Context) ([]string, error)
}
