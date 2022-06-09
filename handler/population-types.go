package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	dperrors "github.com/ONSdigital/dp-net/v2/errors"
	"github.com/ONSdigital/log.go/v2/log"

	"github.com/ONSdigital/dp-population-types-api/config"
	"github.com/ONSdigital/dp-population-types-api/contract"
)

type PopulationTypes struct {
	cfg        *config.Config
	respond    responder
	cantabular cantabularClient
	datasets   datasetAPIClient
}

func NewPopulationTypes(cfg *config.Config, r responder, c cantabularClient, d datasetAPIClient) *PopulationTypes {
	return &PopulationTypes{
		cfg:        cfg,
		respond:    r,
		cantabular: c,
		datasets:   d,
	}
}

func (h *PopulationTypes) Get(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	ptypes, err := h.cantabular.ListDatasets(ctx)
	if err != nil {
		h.respond.Error(ctx, w, dperrors.StatusCode(err), errors.Wrap(
			err,
			"failed to fetch population types",
		))
		return
	}

	log.Info(ctx, "population types found", log.Data{"datasets": ptypes})

	var published []string
	for _, p := range ptypes {
		q := dataset.QueryParams{
			IsBasedOn: p,
			Limit:     1,
		}

		datasets, err := h.datasets.GetDatasets(ctx, "", h.cfg.ServiceAuthToken, "", &q)
		if err != nil {
			if dperrors.StatusCode(err) == http.StatusNotFound {
				continue
			}
			h.respond.Error(ctx, w, dperrors.StatusCode(err), errors.Wrap(
				err,
				"failed to get datasets",
			))
			return
		}

		log.Info(ctx, "datasets found", log.Data{"datasets": datasets})

		var isPublished bool
		for _, d := range datasets.Items {
			if d.Current != nil {
				isPublished = true
				break
			}
		}
		if isPublished {
			published = append(published, p)
		}
	}

	if len(published) == 0 {
		h.respond.Error(ctx, w, http.StatusNotFound, errors.New("no population types found"))
		return
	}

	resp := contract.GetPopulationTypesResponse{
		PopulationTypes: contract.NewPopulationTypes(published),
	}

	h.respond.JSON(ctx, w, http.StatusOK, resp)
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
