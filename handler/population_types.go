package handler

import (
	"context"
	"net/http"

	"github.com/pkg/errors"

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
		h.respond.Error(ctx, w, dperrors.StatusCode(err), &Error{
			err:     errors.Wrap(err, "failed to fetch population types"),
			message: "failed to get population types",
		})
		return
	}

	logData := log.Data{"datasets": ptypes}
	log.Info(ctx, "population types found", logData)

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
			h.respond.Error(ctx, w, dperrors.StatusCode(err), &Error{
				err:     errors.Wrap(err, "failed to check published state"),
				message: "failed to get population types",
				logData: logData,
			})
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
