package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
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

// Get is the handler for GET /area-types
func (h *PopulationTypes) GetAreaTypes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := cantabular.GetGeographyDimensionsRequest{
		Dataset: chi.URLParam(r, "population-type"),
	}

	res, err := h.cantabularClient.GetGeographyDimensions(ctx, req)
	if err != nil {
		h.responder.Error(
			ctx,
			w,
			h.cantabularClient.StatusCode(err), // Can be changed to ctblr.StatusCode(err) once added to Client
			errors.Wrap(err, "failed to get area-types"),
		)
		return
	}

	var resp contract.GetAreaTypesResponse

	if res != nil {
		for _, edge := range res.Dataset.RuleBase.IsSourceOf.Edges {
			resp.AreaTypes = append(resp.AreaTypes, contract.AreaType{
				ID:         edge.Node.Name,
				Label:      edge.Node.Label,
				TotalCount: edge.Node.Categories.TotalCount,
			})
		}
	}

	h.responder.JSON(ctx, w, http.StatusOK, resp)
}
