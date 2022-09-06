package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-population-types-api/config"
	"github.com/ONSdigital/dp-population-types-api/contract"
	"github.com/ONSdigital/log.go/v2/log"

	"github.com/pkg/errors"
)

// Areas handles requests to /area-types
type Areas struct {
	cfg      *config.Config
	datasets datasetAPIClient
	respond  responder
	ctblr    cantabularClient
}

// NewAreas returns a new Areas handler
func NewAreas(cfg *config.Config, r responder, c cantabularClient, d datasetAPIClient) *Areas {
	return &Areas{
		cfg:      cfg,
		datasets: d,
		respond:  r,
		ctblr:    c,
	}
}

// Get is the handler for GET /population-types/{population-type}/area-types/{area-type}/areas
func (h *Areas) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req contract.GetAreasRequest
	if err := parseRequest(r, &req); err != nil {
		h.respond.Error(ctx, w, http.StatusNotFound, &Error{
			err: errors.Wrap(err, "query parameter error"),
		})
		return
	}

	cReq := cantabular.GetAreasRequest{
		PaginationParams: cantabular.PaginationParams{
			Limit:  req.Limit,
			Offset: req.Offset,
		},
		Dataset:  chi.URLParam(r, "population-type"),
		Variable: chi.URLParam(r, "area-type"),
		Category: req.Category,
	}

	logData := log.Data{
		"population_type": cReq.Dataset,
		"area_type":       cReq.Variable,
		"query":           cReq.Category,
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

	res, err := h.ctblr.GetAreas(ctx, cReq)
	if err != nil {
		h.respond.Error(
			ctx,
			w,
			h.ctblr.StatusCode(err),
			&Error{
				err:     errors.Wrap(err, "failed to get areas"),
				message: "failed to get areas",
				logData: logData,
			},
		)
		return
	}

	resp := contract.GetAreasResponse{
		PaginationResponse: contract.PaginationResponse{
			Limit:  req.Limit,
			Offset: req.Offset,
		},
	}

	if res != nil {
		resp.Count = res.Count
		resp.TotalCount = res.TotalCount
		for _, variable := range res.Dataset.Variables.Edges {
			for _, category := range variable.Node.Categories.Search.Edges {
				resp.Areas = append(resp.Areas, contract.Area{
					ID:       category.Node.Code,
					Label:    category.Node.Label,
					AreaType: variable.Node.Name,
				})
			}
		}
	}

	h.respond.JSON(ctx, w, http.StatusOK, resp)
}

// GetID returns the information for a particular area
// for GET /population-types/{population-type}/
func (h *Areas) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cReq := cantabular.GetAreaRequest{
		Dataset:  chi.URLParam(r, "population-type"),
		Variable: chi.URLParam(r, "area-type"),
		Category: chi.URLParam(r, "area-id"),
	}

	logData := log.Data{
		"population_type": cReq.Dataset,
		"area_type":       cReq.Variable,
		"query":           cReq.Category,
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

	res, err := h.ctblr.GetArea(ctx, cReq)
	if err != nil {
		h.respond.Error(
			ctx,
			w,
			h.ctblr.StatusCode(err),
			&Error{
				err:     errors.Wrap(err, "failed to get area"),
				message: "failed to get area",
				logData: logData,
			},
		)
		return
	}

	var area *contract.Area

	for _, variable := range res.Dataset.Variables.Edges {
		for _, category := range variable.Node.Categories.Edges {
			area = &contract.Area{
				ID:       category.Node.Code,
				Label:    category.Node.Label,
				AreaType: variable.Node.Name,
			}
		}
	}

	if area == nil {
		h.respond.Error(
			ctx,
			w,
			http.StatusNotFound,
			&Error{
				err:     errors.Wrap(err, "failed to get area"),
				message: "failed to get area",
				logData: logData,
			},
		)
		return
	}
	// Stop gap until cantabular returns a better solution
	// this temporarily stops partial matches
	if area.ID != cReq.Category {
		h.respond.Error(
			ctx,
			w,
			http.StatusNotFound,
			&Error{
				err:     errors.Wrap(err, "failed to get area"),
				message: "failed to get area",
				logData: logData,
			},
		)
		return
	}
	h.respond.JSON(ctx, w, http.StatusOK, contract.GetAreaResponse{
		Area: *area,
	})
}

func (h *Areas) published(ctx context.Context, populationType string) error {
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
