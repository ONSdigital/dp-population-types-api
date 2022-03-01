package service

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ONSdigital/dp-population-types-api/handler"
	"github.com/ONSdigital/log.go/v2/log"
)

func (svc *Service) buildRoutes(ctx context.Context) {
	svc.Router = chi.NewRouter()
	if svc.HealthCheck != nil {
		svc.Router.Handle("/health", http.HandlerFunc(svc.HealthCheck.Handler))
	}
	svc.publicEndpoints(ctx)
}

func (svc *Service) publicEndpoints(ctx context.Context) {
	log.Info(ctx, "enabling public endpoints")

	// Routes
	populationTypes := handler.NewPopulationTypes(svc.responder)
	svc.Router.Get("/population-types", populationTypes.Get)
}
