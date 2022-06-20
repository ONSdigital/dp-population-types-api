package handler

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	dperrors "github.com/ONSdigital/dp-net/v2/errors"
	dprequest "github.com/ONSdigital/dp-net/v2/request"
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

func (h *PopulationTypes) GetPublic(w http.ResponseWriter, req *http.Request) {
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
			Limit:     1000,
		}

		datasets, err := h.datasets.GetDatasets(ctx, "", "", "", &q)
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

		log.Info(ctx, "datasets found", log.Data{"num_datasets": datasets.Count})

		if len(datasets.Items) > 0 {
			published = append(published, p)
		}
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

func (h *PopulationTypes) GetPrivate(w http.ResponseWriter, req *http.Request) {
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
			Limit:     1000,
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

		log.Info(ctx, "datasets found", log.Data{"num_datasets": datasets.Count})

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

// GetAreaTypes is the handler for GET /population-type/{id}/area-types
// Gets area types from cantabular. If not authenticated, then only return
// area-types if datasets using it published. Whether published determined by making
// unauthenticated request to dataset api.
func (h *PopulationTypes) GetAreaTypes(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authenticated := h.authenticate(r)
	isBasedOn := chi.URLParam(r, "population-type")

	req := cantabular.GetGeographyDimensionsRequest{
		Dataset: isBasedOn,
	}

	res, err := h.cantabular.GetGeographyDimensions(ctx, req)
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

	if !authenticated {
		datasets, err := h.datasets.GetDatasets(
			ctx,
			"",
			h.cfg.ServiceAuthToken,
			"",
			&dataset.QueryParams{
				IsBasedOn: isBasedOn,
				Limit:     1000,
			},
		)
		if err != nil {
			h.respond.Error(
				ctx,
				w,
				dperrors.StatusCode(err),
				errors.New("failed to get area types: failed to get population type"),
			)
			return

		}

		var isPublished bool
		for _, d := range datasets.Items {
			if d.Current != nil {
				isPublished = true
				break
			}
		}

		if !isPublished {
			h.respond.Error(
				ctx,
				w,
				http.StatusNotFound,
				errors.New("population type not found"),
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

	h.respond.JSON(ctx, w, http.StatusOK, resp)
}

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
