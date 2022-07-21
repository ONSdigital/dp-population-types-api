package handler

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-population-types-api/contract"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

func (h *PopulationTypes) GetAreaTypes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cReq := cantabular.GetGeographyDimensionsRequest{
		Dataset: chi.URLParam(r, "population-type"),
	}

	logData := log.Data{
		"population_type": cReq.Dataset,
	}

	// only return results for published population-types on web
	if !h.cfg.EnablePrivateEndpoints {
		if err := h.published(ctx, cReq.Dataset); err != nil {
			h.respond.Error(ctx, w, http.StatusNotFound, &Error{
				err:     errors.Wrap(err, "failed to check published state"),
				message: "population type not found",
				logData: logData,
			})
			return
		}
	}

	res, err := h.cantabular.GetGeographyDimensions(ctx, cReq)
	if err != nil {
		h.respond.Error(ctx, w, h.cantabular.StatusCode(err), &Error{
			err:     errors.Wrap(err, "failed to get geography dimensions"),
			message: "failed to get area-types",
			logData: logData,
		})
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

	logData := log.Data{
		"population_type": cReq.Dataset,
		"area_type":       cReq.Variable,
	}
	// only return results for published population-types on web
	if !h.cfg.EnablePrivateEndpoints {
		if err := h.published(ctx, cReq.Dataset); err != nil {
			h.respond.Error(ctx, w, http.StatusNotFound, &Error{
				err:     errors.Wrap(err, "failed to check published state"),
				message: "population type not found",
				logData: logData,
			})
			return
		}
	}

	res, err := h.cantabular.GetParents(ctx, cReq)
	if err != nil {
		h.respond.Error(ctx, w, h.cantabular.StatusCode(err), &Error{
			err:     errors.Wrap(err, "failed to get parents"),
			message: "failed to get parents",
			logData: logData,
		})
		return
	}

	var resp contract.GetAreaTypeParentsResponse

	if l := len(res.Dataset.Variables.Edges); l != 1 {
		logData["edges_expected_length"] = 1
		logData["edges_length"] = l
		h.respond.Error(ctx, w, http.StatusInternalServerError, &Error{
			err:     errors.New("invalid response from Cantabular"),
			message: "failed to get parents",
			logData: logData,
		})
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
