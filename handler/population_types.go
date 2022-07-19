package handler

import (
	"context"
	"fmt"
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

	// return all population types on publishing
	if h.cfg.EnablePrivateEndpoints {
		h.respond.JSON(ctx, w, http.StatusOK, contract.GetPopulationTypesResponse{
			PopulationTypes: contract.NewPopulationTypes(ptypes),
		})
		return
	}

	var published []string
	for _, p := range ptypes {
		if err := h.published(ctx, p); err != nil {
			if dperrors.StatusCode(err) == http.StatusNotFound {
				continue
			}
			h.respond.Error(ctx, w, dperrors.StatusCode(err), errors.Wrap(
				err,
				"failed to get datasets",
			))
			return
		}

		published = append(published, p)
	}

	if len(published) == 0 {
		h.respond.Error(
			ctx,
			w,
			http.StatusNotFound,
			errors.New("no population types found"),
		)
		return
	}

	resp := contract.GetPopulationTypesResponse{
		PopulationTypes: contract.NewPopulationTypes(published),
	}

	h.respond.JSON(ctx, w, http.StatusOK, resp)
}

func (h *PopulationTypes) GetAreaTypes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	isBasedOn := chi.URLParam(r, "population-type")

	cReq := cantabular.GetGeographyDimensionsRequest{
		Dataset: isBasedOn,
	}

	// only return results for published population-types on web
	if !h.cfg.EnablePrivateEndpoints {
		if err := h.published(ctx, cReq.Dataset); err != nil {
			h.respond.Error(ctx, w, http.StatusNotFound, errors.New("population type not found"))
			return
		}
	}

	res, err := h.cantabular.GetGeographyDimensions(ctx, cReq)
	if err != nil {
		h.respond.Error(
			ctx,
			w,
			h.cantabular.StatusCode(err),
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

	h.respond.JSON(ctx, w, http.StatusOK, resp)
}

func (h *PopulationTypes) GetAreaTypeParents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cReq := cantabular.GetParentsRequest{
		Dataset:  chi.URLParam(r, "population-type"),
		Variable: chi.URLParam(r, "area-type"),
	}

	// only return results for published population-types on web
	if !h.cfg.EnablePrivateEndpoints {
		if err := h.published(ctx, cReq.Dataset); err != nil {
			h.respond.Error(ctx, w, http.StatusNotFound, errors.New("population type not found"))
			return
		}
	}

	res, err := h.cantabular.GetParents(ctx, cReq)
	if err != nil {
		h.respond.Error(ctx, w, h.cantabular.StatusCode(err), errors.Wrap(err, "failed to get parents"))
		return
	}

	var resp contract.GetAreaTypeParentsResponse

	if len(res.Dataset.Variables.Edges) != 1 {
		h.respond.Error(ctx, w, http.StatusInternalServerError, fmt.Errorf("failed to get parents"))
		return
	}

	for _, e := range res.Dataset.Variables.Edges[0].Node.IsSourceOf.Edges {
		resp.AreaTypes = append(resp.AreaTypes, contract.AreaType{
			ID:         e.Node.Name,
			Label:      e.Node.Label,
			TotalCount: e.Node.Categories.TotalCount,
		})
	}

	h.respond.JSON(ctx, w, http.StatusOK, resp)
}

func (h *PopulationTypes) published(ctx context.Context, populationType string) error {
	datasets, err := h.datasets.GetDatasets(
		ctx,
		"",
		"",
		"",
		&dataset.QueryParams{
			IsBasedOn: populationType,
			Limit:     1000,
		},
	)
	if err != nil {
		return errors.Wrap(err, "failed to get datasets for population type")
	}

	if datasets.TotalCount == 0 {
		return errors.New("no published datasets found for population type")
	}

	return nil
}
