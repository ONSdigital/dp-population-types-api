package service

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/ONSdigital/dp-authorisation/auth"
	dphandlers "github.com/ONSdigital/dp-net/v2/handlers"
	"github.com/ONSdigital/dp-population-types-api/config"
	"github.com/ONSdigital/dp-population-types-api/handler"
	"github.com/ONSdigital/dp-population-types-api/middleware"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/riandyrn/otelchi"
)

func (svc *Service) buildRoutes(ctx context.Context) {
	svc.Router = chi.NewRouter()

	cfg, _ := config.Get()
	svc.Router.Use(otelchi.Middleware(cfg.OTServiceName))

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
		svc.mongoClient,
	)
	svc.Router.Get("/population-types", populationTypes.GetAll)
	svc.Router.Get("/population-types/{population-type}", populationTypes.Get)

	dimensions := handler.NewDimensions(
		svc.Config,
		svc.responder,
		svc.cantabularClient,
		svc.datasetAPIClient,
	)
	svc.Router.Get("/population-types/{population-type}/dimensions", dimensions.GetAll)
	svc.Router.Get("/population-types/{population-type}/dimensions/{dimension}/categorisations", dimensions.GetCategorisations)
	svc.Router.Get("/population-types/{population-type}/dimension-categories", dimensions.GetDimensionCategories)

	svc.Router.Get("/population-types/{population-type}/dimensions/{dimension}/base", dimensions.GetBaseVariable)
	svc.Router.Get("/population-types/{population-type}/dimensions-description", dimensions.GetDescription)
	svc.Router.Get("/population-types/{population-type}/blocked-areas-count", dimensions.GetBlockedAreaCount)

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

	metadata := handler.NewMetadata(svc.Config, svc.responder, svc.mongoClient)
	svc.Router.Get("/population-types/{population-type}/metadata", metadata.Get)

	if svc.Config.CensusObservationsFF {
		censusObservations := handler.NewCensusObservations(svc.Config, svc.responder, svc.cantabularClient)
		svc.Router.Get("/population-types/{population-type}/census-observations", censusObservations.Get)
	}
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
		svc.mongoClient,
	)
	r.Get("/population-types", populationTypes.GetAll)
	r.Get("/population-types/{population-type}", populationTypes.Get)

	dimensions := handler.NewDimensions(
		svc.Config,
		svc.responder,
		svc.cantabularClient,
		svc.datasetAPIClient,
	)
	r.Get("/population-types/{population-type}/dimensions", dimensions.GetAll)
	r.Get("/population-types/{population-type}/dimensions/{dimension}/categorisations", dimensions.GetCategorisations)
	r.Get("/population-types/{population-type}/dimensions/{dimension}/base", dimensions.GetBaseVariable)
	r.Get("/population-types/{population-type}/dimensions-description", dimensions.GetDescription)
	r.Get("/population-types/{population-type}/blocked-areas-count", dimensions.GetBlockedAreaCount)
	r.Get("/population-types/{population-type}/dimension-categories", dimensions.GetDimensionCategories)

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

	metadata := handler.NewMetadata(svc.Config, svc.responder, svc.mongoClient)

	r.Get("/population-types/{population-type}/metadata", metadata.Get)
	r.Put("/population-types/{population-type}/metadata", metadata.Put)

	if svc.Config.CensusObservationsFF {
		censusObservations := handler.NewCensusObservations(svc.Config, svc.responder, svc.cantabularClient)
		r.Get("/population-types/{population-type}/census-observations", censusObservations.Get)
	}

	svc.Router.Mount("/", r)
}
