package service

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ONSdigital/dp-authorisation/auth"
	dphandlers "github.com/ONSdigital/dp-net/v2/handlers"
	"github.com/ONSdigital/dp-population-types-api/handler"
	"github.com/ONSdigital/dp-population-types-api/middleware"

	"github.com/ONSdigital/log.go/v2/log"
)

func (svc *Service) buildRoutes(ctx context.Context) {
	svc.Router = chi.NewRouter()
	if svc.HealthCheck != nil {
		svc.Router.Handle("/health", http.HandlerFunc(svc.HealthCheck.Handler))
	}

	if svc.Config.EnablePrivateEndpoints {
		svc.privateEndpoints(ctx)
	} else {
		svc.publicEndpoints(ctx)
	}
}

func (svc *Service) publicEndpoints(ctx context.Context) {
	log.Info(ctx, "enabling public endpoints")

	// Routes
	populationTypes := handler.NewPopulationTypes(
		svc.Config,
		svc.responder,
		svc.cantabularClient,
		svc.datasetAPIClient,
	)
	svc.Router.Get("/population-types", populationTypes.Get)
	svc.Router.Get("/population-types/{population-type}/area-types", populationTypes.GetAreaTypes)

	areas := handler.NewAreas(svc.Config, svc.datasetAPIClient, svc.responder, svc.cantabularClient)
	svc.Router.Get("/population-types/{population-type}/area-types/{area-type}/areas", areas.Get)
	svc.Router.Get("/population-types/{population-type}/area-types/{area-type}/parents", populationTypes.GetAreaTypeParents)
}

func (svc *Service) privateEndpoints(ctx context.Context) {
	log.Info(ctx, "enabling private endpoints")

	r := chi.NewRouter()

	// Middleware
	permissions := middleware.NewPermissions(svc.Config.ZebedeeURL, svc.Config.EnablePermissionsAuth)
	checkIdentity := dphandlers.IdentityWithHTTPClient(svc.identityClient)

	// Routes
	populationTypes := handler.NewPopulationTypes(
		svc.Config,
		svc.responder,
		svc.cantabularClient,
		svc.datasetAPIClient,
	)

	r.Use(checkIdentity)
	r.Use(middleware.LogIdentity())
	r.Use(permissions.Require(auth.Permissions{Read: true}))

	r.Get("/population-types", populationTypes.Get)
	r.Get("/population-types/{population-type}/area-types", populationTypes.GetAreaTypes)
	r.Get("/population-types/{population-type}/area-types/{area-type}/parents", populationTypes.GetAreaTypeParents)

	areas := handler.NewAreas(svc.Config, svc.datasetAPIClient, svc.responder, svc.cantabularClient)
	r.Get("/population-types/{population-type}/area-types/{area-type}/areas", areas.Get)

	svc.Router.Mount("/", r)
}
