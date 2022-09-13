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

	dimensions := handler.NewDimensions(
		svc.Config,
		svc.responder,
		svc.cantabularClient,
		svc.datasetAPIClient,
	)
	svc.Router.Get("/population-types/{population-type}/dimensions", dimensions.GetAll)
	svc.Router.Get("/population-types/{population-type}/dimensions/{dimension}/categorisations", dimensions.GetCategorisations)

	areaTypes := handler.NewAreaTypes(
		svc.Config,
		svc.responder,
		svc.cantabularClient,
		svc.datasetAPIClient,
	)
	svc.Router.Get("/population-types/{population-type}/area-types", areaTypes.Get)
	svc.Router.Get("/population-types/{population-type}/area-types/{area-type}/parents", areaTypes.GetParents)
	svc.Router.Get("/population-types/{population-type}/area-types/{area-type}/parents/{parent-area-type}/areas-count", areaTypes.GetParentAreaCount)

	areas := handler.NewAreas(
		svc.Config,
		svc.responder,
		svc.cantabularClient,
		svc.datasetAPIClient,
	)
	svc.Router.Get("/population-types/{population-type}/area-types/{area-type}/areas", areas.GetAll)
	svc.Router.Get("/population-types/{population-type}/area-types/{area-type}/areas/{area-id}", areas.Get)
}

func (svc *Service) privateEndpoints(ctx context.Context) {
	log.Info(ctx, "enabling private endpoints")

	r := chi.NewRouter()

	// Middleware
	permissions := middleware.NewPermissions(svc.Config.ZebedeeURL, svc.Config.EnablePermissionsAuth)
	checkIdentity := dphandlers.IdentityWithHTTPClient(svc.identityClient)

	r.Use(checkIdentity)
	r.Use(middleware.LogIdentity())
	r.Use(permissions.Require(auth.Permissions{Read: true}))

	// Routes
	populationTypes := handler.NewPopulationTypes(
		svc.Config,
		svc.responder,
		svc.cantabularClient,
		svc.datasetAPIClient,
	)
	r.Get("/population-types", populationTypes.Get)

	dimensions := handler.NewDimensions(
		svc.Config,
		svc.responder,
		svc.cantabularClient,
		svc.datasetAPIClient,
	)
	r.Get("/population-types/{population-type}/dimensions", dimensions.GetAll)
	r.Get("/population-types/{population-type}/dimensions/{dimension}/categorisations", dimensions.GetCategorisations)

	areaTypes := handler.NewAreaTypes(
		svc.Config,
		svc.responder,
		svc.cantabularClient,
		svc.datasetAPIClient,
	)
	r.Get("/population-types/{population-type}/area-types", areaTypes.Get)
	r.Get("/population-types/{population-type}/area-types/{area-type}/parents", areaTypes.GetParents)
	r.Get("/population-types/{population-type}/area-types/{area-type}/parents/{parent-area-type}/areas-count", areaTypes.GetParentAreaCount)

	areas := handler.NewAreas(
		svc.Config,
		svc.responder,
		svc.cantabularClient,
		svc.datasetAPIClient,
	)
	r.Get("/population-types/{population-type}/area-types/{area-type}/areas", areas.GetAll)
	r.Get("/population-types/{population-type}/area-types/{area-type}/areas/{area-id}", areas.Get)

	svc.Router.Mount("/", r)
}
