package handler

import (
	"net/http"

	"github.com/ONSdigital/dp-population-types-api/config"
	"github.com/ONSdigital/dp-population-types-api/contract"
	"github.com/ONSdigital/dp-population-types-api/datastore"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

type PopulationTypeMetadata struct {
	cfg         *config.Config
	respond     responder
	MongoClient Datastore
}

// NewMetadata returns a new Metadata Handler
func NewMetadata(cfg *config.Config, r responder, d Datastore) *PopulationTypeMetadata {
	return &PopulationTypeMetadata{
		cfg:         cfg,
		respond:     r,
		MongoClient: d,
	}
}

func (m *PopulationTypeMetadata) Put(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	populationType := chi.URLParam(r, "population-type")
	var req contract.PutMetadataRequest

	if err := parseRequest(r, &req); err != nil {
		m.respond.Error(
			ctx,
			w,
			http.StatusBadRequest,
			errors.Wrap(err, "query parameter error"),
		)
		return
	}

	metadata := datastore.DefaultDatasetMetadata{
		ID:               populationType,
		DefaultDatasetID: req.DefaultDatasetID,
		Edition:          req.Edition,
		Version:          req.Version,
	}

	if err := m.MongoClient.PutDefaultDatasetMetadata(ctx, metadata); err != nil {
		m.respond.Error(
			ctx,
			w,
			statusCode(err),
			errors.Wrap(err, "failed to get metadata"),
		)
		return
	}

	m.respond.JSON(ctx, w, http.StatusCreated, nil)
}

func (m *PopulationTypeMetadata) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	populationType := chi.URLParam(r, "population-type")

	metadata, err := m.MongoClient.GetDefaultDatasetMetadata(ctx, populationType)
	if err != nil {
		m.respond.Error(
			ctx,
			w,
			statusCode(err),
			errors.Wrap(err, "failed to get metadata"),
		)
		return
	}

	m.respond.JSON(ctx, w, http.StatusOK, contract.GetMetadataResponse{
		PopulationType:   populationType,
		DefaultDatasetID: metadata.DefaultDatasetID,
		Edition:          metadata.Edition,
		Version:          metadata.Version,
	})
}
