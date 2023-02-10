package handler

import (
	"net/http"

	"github.com/ONSdigital/dp-population-types-api/config"
	"github.com/ONSdigital/dp-population-types-api/contract"
	"github.com/ONSdigital/dp-population-types-api/datastore"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

type Metadata struct {
	cfg         *config.Config
	respond     responder
	MongoClient Datastore
}

// NewMetada returns a new Metadata Handler
func NewMetada(cfg *config.Config, r responder, databaseClient Datastore) *Metadata {
	return &Metadata{
		cfg:         cfg,
		respond:     r,
		MongoClient: databaseClient,
	}
}

func (m *Metadata) PutMetadata(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	populationType := chi.URLParam(r, "population-type")
	var req contract.PutMetadataRequest
	if err := parseRequest(r, &req); err != nil {
		m.respond.Error(ctx, w, http.StatusNotFound, &Error{
			err: errors.Wrap(err, "query parameter error"),
		})
		return
	}

	metadata := datastore.DefaultDatasetMetadata{
		ID:               populationType,
		DefaultDatasetID: req.DefaultDatasetID,
	}

	if err := m.MongoClient.PutDefaultDatasetMetadata(ctx, metadata); err != nil {
		m.respond.Error(ctx, w, http.StatusInternalServerError, errors.New("Failed to get metadata"))
		return
	}

	m.respond.JSON(ctx, w, http.StatusCreated, nil)

}

func (m *Metadata) GetMetadata(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	populationType := chi.URLParam(r, "population-type")

	metadata, err := m.MongoClient.GetDefaultDatasetMetadata(ctx, populationType)
	if err != nil {
		m.respond.Error(
			ctx,
			w,
			statusCode(err),
			errors.New("Failed to get metadata"),
		)
		return
	}

	m.respond.JSON(ctx, w, http.StatusOK, contract.GetMetadataResponse{
		PopulationType:   populationType,
		DefaultDatasetID: metadata.DefaultDatasetID,
	})
}
