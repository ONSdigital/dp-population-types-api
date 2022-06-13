package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-api-clients-go/v2/identity"
	dprequest "github.com/ONSdigital/dp-net/v2/request"
	"github.com/ONSdigital/dp-population-types-api/contract"
)

type PopulationTypes struct {
	responder        responder
	cantabularClient cantabularClient
	datasetClient    DatasetAPIClient
	identityClient   *identity.Client
}

func NewPopulationTypes(responder responder, cantabularClient cantabularClient, identityClient *identity.Client, datasetClient DatasetAPIClient) *PopulationTypes {
	return &PopulationTypes{
		responder:        responder,
		cantabularClient: cantabularClient,
		identityClient:   identityClient,
		datasetClient:    datasetClient,
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

/*
   GET population-type/<id>/area-types
   Gets area types from cantabular.
   If not authenticated, then only return
   area-types if datasets using it published.
   Whether published determined by making
   unauthenticated request to dataset api.
*/
func (h *PopulationTypes) GetAreaTypes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	authenticated := h.authenticate(r)

	isBasedOn := chi.URLParam(r, "population-type")

	req := cantabular.GetGeographyDimensionsRequest{
		Dataset: isBasedOn,
	}

	res, err := h.cantabularClient.GetGeographyDimensions(ctx, req)
	if err != nil {
		h.responder.Error(
			ctx,
			w,
			h.cantabularClient.StatusCode(err),
			errors.Wrap(err, "failed to get area-types"),
		)
		return
	}

	var resp contract.GetAreaTypesResponse

	if !authenticated {
		datasets, err := h.datasetClient.GetDatasets(
			ctx,
			"",
			"",
			"",
			&dataset.QueryParams{IsBasedOn: isBasedOn, Limit: 100},
		)
		if err != nil {
			h.responder.Error(
				ctx,
				w,
				http.StatusInternalServerError,
				errors.New("failed to get area-types: internal server error"),
			)
			return

		}
		if datasets.TotalCount == 0 {
			h.responder.JSON(
				ctx,
				w,
				http.StatusNotFound,
				errors.New("dataset not found"),
			)
			return
		}
	}

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

// Required for GET area types
func (h *PopulationTypes) authenticate(r *http.Request) bool {
	var authorised bool

	var hasCallerIdentity, hasUserIdentity bool
	callerIdentity := dprequest.Caller(r.Context())

	if callerIdentity != "" {
		hasCallerIdentity = true
	}

	userIdentity := dprequest.User(r.Context())
	if userIdentity != "" {
		hasUserIdentity = true
	}

	if hasCallerIdentity || hasUserIdentity {
		authorised = true
	}

	return authorised
}
