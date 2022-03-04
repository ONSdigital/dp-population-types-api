package handler

import (
	"net/http"

	"github.com/pkg/errors"

	"github.com/ONSdigital/dp-population-types-api/contract"
)

type PopulationTypes struct {
	responder        responder
	cantabularClient cantabularClient
}

func NewPopulationTypes(responder responder, cantabularClient cantabularClient) *PopulationTypes {
	return &PopulationTypes{
		responder:        responder,
		cantabularClient: cantabularClient,
	}
}

func (h *PopulationTypes) Get(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	data, err := h.cantabularClient.ListDatasets(ctx)
	if err != nil {
		wrappedErr := errors.Wrap(err, "failed to fetch population types")
		h.responder.Error(ctx, w, http.StatusInternalServerError, wrappedErr)
	} else {
		body := contract.NewPopulationTypes(data)
		h.responder.JSON(ctx, w, http.StatusOK, body)
	}
}
