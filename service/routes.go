package service

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ONSdigital/log.go/v2/log"
)

func (svc *Service) buildRoutes(ctx context.Context) {
	svc.router = chi.NewRouter()
	svc.router.Handle("/health", http.HandlerFunc(svc.HealthCheck.Handler))
	svc.publicEndpoints(ctx)
}

func (svc *Service) publicEndpoints(ctx context.Context) {
	log.Info(ctx, "enabling public endpoints")

	// Routes
	// areaTypes := handler.NewAreaTypes(svc.responder, svc.cantabularClient)

	// svc.router.Get("/area-types", areaTypes.Get)
}
