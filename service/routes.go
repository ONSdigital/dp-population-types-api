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

	println("************** BUILD ROUTES", svc.Config.EnablePrivateEndpoints)
	if svc.Config.EnablePrivateEndpoints {
		svc.privateEndpoints(ctx)
	} else {

		svc.publicEndpoints(ctx)
	}

}

func (svc *Service) publicEndpoints(ctx context.Context) {
	log.Info(ctx, "enabling public endpoints")

	// Routes
	populationTypes := handler.NewPopulationTypes(svc.responder, svc.cantabularClient, svc.identityClient, svc.DatasetAPIClient)
	svc.Router.Get("/population-types", populationTypes.Get)
	svc.Router.Get("/population-types/{population-type}/area-types", populationTypes.GetAreaTypes)
}

func (svc *Service) privateEndpoints(ctx context.Context) {
	log.Info(ctx, "enabling private endpoints")

	populationTypes := handler.NewPopulationTypes(svc.responder, svc.cantabularClient, svc.identityClient, svc.DatasetAPIClient)

	r := chi.NewRouter()

	permissions := middleware.NewPermissions(svc.Config.ZebedeeURL, svc.Config.EnablePermissionsAuth)
	checkIdentity := dphandlers.IdentityWithHTTPClient(svc.identityClient)

	r.Use(checkIdentity)
	r.Use(middleware.LogIdentity())
	r.Use(permissions.Require(auth.Permissions{Read: true}))

	svc.Router.Get("/population-types", populationTypes.Get)
	svc.Router.Get("/population-types/{population-type}/area-types", populationTypes.GetAreaTypes)

	// added from cant filt flex
	svc.Router.Mount("/", r)
}
