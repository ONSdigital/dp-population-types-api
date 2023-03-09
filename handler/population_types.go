package handler

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"

	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	dperrors "github.com/ONSdigital/dp-net/v2/errors"
	"github.com/ONSdigital/log.go/v2/log"

	"github.com/ONSdigital/dp-population-types-api/config"
	"github.com/ONSdigital/dp-population-types-api/contract"
)

type PopulationTypes struct {
	cfg         *config.Config
	respond     responder
	cantabular  cantabularClient
	datasets    datasetAPIClient
	mongoClient Datastore
}

func NewPopulationTypes(cfg *config.Config, r responder, c cantabularClient, d datasetAPIClient, m Datastore) *PopulationTypes {
	return &PopulationTypes{
		cfg:         cfg,
		respond:     r,
		cantabular:  c,
		datasets:    d,
		mongoClient: m,
	}
}

func (h *PopulationTypes) GetAll(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	r := contract.GetPopulationTypesRequest{}
	if err := parseRequest(req, &r); err != nil {
		h.respond.Error(ctx, w, http.StatusBadRequest, &Error{
			err: errors.Wrap(err, "invalid request"),
		})
		return
	}

	ptypes, err := h.cantabular.ListDatasets(ctx)
	if err != nil {
		h.respond.Error(ctx, w, dperrors.StatusCode(err), &Error{
			err:     errors.Wrap(err, "failed to fetch population types"),
			message: "failed to get population types",
		})
		return
	}

	logData := log.Data{"datasets": ptypes}
	log.Info(ctx, "population types found", logData)
	lpt := len(ptypes.Datasets)

	var defaultDatasets []string

	// return just default datasets only for BYO
	if r.DefaultDatasets {
		response, err := h.mongoClient.GetDefaultDatasetPopulationTypes(ctx)
		if err != nil {
			h.respond.Error(
				ctx,
				w,
				http.StatusInternalServerError,
				errors.Wrap(err, "Failed to get metadata"),
			)
			return
		}
		defaultDatasets = response

		resp := contract.GetPopulationTypesResponse{
			PaginationResponse: contract.PaginationResponse{
				Limit:  r.Limit,
				Offset: r.Offset,
			},
		}

		for _, pt := range ptypes.Datasets {
			resp.Items = append(resp.Items, contract.PopulationType{
				Name:        pt.Name,
				Label:       pt.Label,
				Description: pt.Description,
				Type:        pt.Type,
			})
		}

		resp.Items = filterPopulationTypes(defaultDatasets, resp.Items)
		defaultDatasets = response
		resp.Paginate()
		resp.TotalCount = len(resp.Items)
		h.respond.JSON(ctx, w, http.StatusOK, resp)
		return
	}

	// return all population types on publishing
	if h.cfg.EnablePrivateEndpoints {
		resp := contract.GetPopulationTypesResponse{
			PaginationResponse: contract.PaginationResponse{
				Limit:      r.Limit,
				Offset:     r.Offset,
				TotalCount: lpt,
			},
		}
		for _, pt := range ptypes.Datasets {
			resp.Items = append(resp.Items, contract.PopulationType{
				Name:        pt.Name,
				Label:       pt.Label,
				Description: pt.Description,
				Type:        pt.Type,
			})
		}

		resp.Paginate()
		h.respond.JSON(ctx, w, http.StatusOK, resp)
		return
	}

	var published []contract.PopulationType
	for _, pt := range ptypes.Datasets {
		if err := h.published(ctx, pt.Name); err != nil {
			if dperrors.StatusCode(err) == http.StatusNotFound {
				continue
			}
			h.respond.Error(ctx, w, dperrors.StatusCode(err), &Error{
				err:     errors.Wrap(err, "failed to check published state"),
				message: "failed to get population types",
				logData: logData,
			})
			return
		}

		published = append(published, contract.PopulationType{
			Name:        pt.Name,
			Label:       pt.Label,
			Description: pt.Description,
			Type:        pt.Type,
		})
	}

	l := len(published)
	if l == 0 {
		h.respond.Error(
			ctx,
			w,
			http.StatusNotFound,
			errors.New("no population types found"),
		)
		return
	}

	resp := contract.GetPopulationTypesResponse{
		PaginationResponse: contract.PaginationResponse{
			Limit:      r.Limit,
			TotalCount: l,
			Offset:     r.Offset,
		},
		Items: published,
	}

	resp.Paginate()

	h.respond.JSON(ctx, w, http.StatusOK, resp)
}

func (h *PopulationTypes) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ptype := chi.URLParam(r, "population-type")
	logData := log.Data{
		"population_type": ptype,
	}

	if !h.cfg.EnablePrivateEndpoints {
		if err := h.published(ctx, ptype); err != nil {
			h.respond.Error(ctx, w, http.StatusNotFound, &Error{
				err:     errors.Wrap(err, "failed to check published state"),
				message: "population type not found",
				logData: logData,
			})
			return
		}
	}

	ptypes, err := h.cantabular.ListDatasets(ctx)
	if err != nil {
		h.respond.Error(ctx, w, dperrors.StatusCode(err), &Error{
			err:     errors.Wrap(err, "failed to fetch population types"),
			message: "failed to get population types",
		})
		return
	}

	var resp contract.GetPopulationTypeResponse

	for _, p := range ptypes.Datasets {
		if p.Name == ptype {
			resp.PopulationType = contract.PopulationType{
				Name:        p.Name,
				Label:       p.Label,
				Description: p.Description,
				Type:        p.Type,
			}
			break
		}
	}

	if resp.PopulationType.Name == "" {
		h.respond.Error(ctx, w, http.StatusNotFound, &Error{
			err:     errors.New("population type not found"),
			logData: logData,
		})
		return
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
