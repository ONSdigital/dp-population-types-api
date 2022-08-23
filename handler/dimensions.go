package handler

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-population-types-api/config"
	"github.com/ONSdigital/dp-population-types-api/contract"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

type Dimensions struct {
	cfg        *config.Config
	respond    responder
	cantabular cantabularClient
	datasets   datasetAPIClient
}

// NewDimensions returns a new AreaTypes handler set
func NewDimensions(cfg *config.Config, r responder, c cantabularClient, d datasetAPIClient) *Dimensions {
	return &Dimensions{
		cfg:        cfg,
		respond:    r,
		cantabular: c,
		datasets:   d,
	}
}

// GetAll is the handler for GET /population-types/{population-type}/dimensions
func (h *Dimensions) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := contract.GetDimensionsRequest{
		PopulationType: chi.URLParam(r, "population-type"),
	}
	if err := parseRequest(r, &req); err != nil {
		h.respond.Error(ctx, w, http.StatusBadRequest, &Error{
			err: errors.Wrap(err, "invalid request"),
		})
		return
	}

	logData := log.Data{
		"request": req,
	}

	// only return results for published population-types on web
	if !h.cfg.EnablePrivateEndpoints {
		if err := h.published(ctx, req.PopulationType); err != nil {
			h.respond.Error(ctx, w, http.StatusNotFound, &Error{
				err:     errors.Wrap(err, "failed to check published state"),
				message: "population type not found",
				logData: logData,
			})
			return
		}
	}

	cReq := cantabular.GetDimensionsRequest{
		Dataset: req.PopulationType,
		Text:    req.SearchText,
		PaginationParams: cantabular.PaginationParams{
			Limit:  req.QueryParams.Limit,
			Offset: req.QueryParams.Offset,
		},
	}

	res, err := h.cantabular.GetDimensions(ctx, cReq)
	if err != nil {
		h.respond.Error(ctx, w, h.cantabular.StatusCode(err), &Error{
			err:     errors.Wrap(err, "failed to get dimensions"),
			message: "failed to get dimensions",
			logData: logData,
		})
		return
	}

	resp := contract.GetDimensionsResponse{
		PaginationResponse: contract.PaginationResponse{
			Limit:  req.Limit,
			Offset: req.Offset,
		},
	}

	if res != nil {
		resp.Count = res.Count
		resp.TotalCount = res.TotalCount
		for _, edge := range res.Dataset.Variables.Search.Edges {
			resp.Dimensions = append(resp.Dimensions, contract.Dimension{
				ID:         edge.Node.Name,
				Label:      edge.Node.Label,
				TotalCount: edge.Node.Categories.TotalCount,
			})
		}
	}

	h.respond.JSON(ctx, w, http.StatusOK, resp)
}

func (h *Dimensions) published(ctx context.Context, populationType string) error {
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
