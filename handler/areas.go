package handler

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-population-types-api/config"
	"github.com/ONSdigital/dp-population-types-api/contract"
	"github.com/go-chi/chi/v5"

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
func NewAreas(cfg *config.Config, d datasetAPIClient, r responder, c cantabularClient) *Areas {
	return &Areas{
		cfg:      cfg,
		datasets: d,
		respond:  r,
		ctblr:    c,
	}
}

// Get is the handler for GET /population-types/{population-type}/area-types/{area-type}/areas
func (h *Areas) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cReq := cantabular.GetAreasRequest{
		Dataset:  chi.URLParam(r, "population-type"),
		Variable: chi.URLParam(r, "area-type"),
		Category: r.URL.Query().Get("q"),
	}

	// only return results for published population-types on web
	if !h.cfg.EnablePrivateEndpoints {
		if err := h.published(ctx, cReq.Dataset); err != nil {
			h.respond.Error(ctx, w, http.StatusNotFound, errors.New("population type not found"))
			return
		}
	}

	res, err := h.ctblr.GetAreas(ctx, cReq)
	if err != nil {
		msg := "failed to get areas"
		h.respond.Error(
			ctx,
			w,
			h.ctblr.StatusCode(err),
			&Error{
				err:     errors.Wrap(err, msg),
				message: msg,
			},
		)
		return
	}

	var resp contract.GetAreasResponse

	for _, variable := range res.Dataset.RuleBase.IsSourceOf.Search.Edges {
		for _, category := range variable.Node.Categories.Search.Edges {
			resp.Areas = append(resp.Areas, contract.Area{
				ID:       category.Node.Code,
				Label:    category.Node.Label,
				AreaType: variable.Node.Name,
			})
		}
	}

	h.respond.JSON(ctx, w, http.StatusOK, resp)
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
