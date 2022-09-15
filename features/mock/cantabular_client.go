package mock

import (
	"context"
	"errors"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"

	dperrors "github.com/ONSdigital/dp-api-clients-go/v2/errors"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/log.go/v2/log"
)

var (
	errFailedToRespond = errors.New("cantabular failed to respond")
)

type CantabularClient struct {
	Healthy                        bool
	BadRequest                     bool
	NotFound                       bool
	GetGeographyDimensionsResponse *cantabular.GetGeographyDimensionsResponse
	GetDimensionsResponse          *cantabular.GetDimensionsResponse
	GetAreasResponse               *cantabular.GetAreasResponse
	GetAreaResponse                *cantabular.GetAreaResponse
	GetParentsResponse             *cantabular.GetParentsResponse
	GetParentAreaCountResult       *cantabular.GetParentAreaCountResult
	GetCategorisationsResponse     *cantabular.GetCategorisationsResponse
	ListDatasetsResponse           []string
}

func (c *CantabularClient) Checker(_ context.Context, _ *healthcheck.CheckState) error {
	return nil
}

func (c *CantabularClient) CheckerAPIExt(_ context.Context, _ *healthcheck.CheckState) error {
	return nil
}

func (c *CantabularClient) ListDatasets(_ context.Context) ([]string, error) {
	if !c.Healthy {
		return nil, errFailedToRespond
	}

	return c.ListDatasetsResponse, nil
}

func (c *CantabularClient) GetGeographyDimensions(_ context.Context, _ cantabular.GetGeographyDimensionsRequest) (*cantabular.GetGeographyDimensionsResponse, error) {
	if !c.Healthy {
		return nil, dperrors.New(
			errors.New("error(s) returned by graphQL query"),
			http.StatusNotFound,
			log.Data{"errors": map[string]string{"message": "404 Not Found: dataset not loaded in this server"}},
		)
	}

	return c.GetGeographyDimensionsResponse, nil
}

func (c *CantabularClient) GetDimensions(_ context.Context, _ cantabular.GetDimensionsRequest) (*cantabular.GetDimensionsResponse, error) {
	if !c.Healthy {
		return nil, dperrors.New(
			errors.New("error(s) returned by graphQL query"),
			http.StatusNotFound,
			log.Data{"errors": map[string]string{"message": "404 Not Found: dataset not loaded in this server"}},
		)
	}

	return c.GetDimensionsResponse, nil
}

func (c *CantabularClient) StatusCode(_ error) int {
	if c.BadRequest {
		return http.StatusBadRequest
	}

	if c.NotFound {
		return http.StatusNotFound
	}

	return http.StatusNotFound
}

func (c *CantabularClient) GetAreas(_ context.Context, _ cantabular.GetAreasRequest) (*cantabular.GetAreasResponse, error) {
	if !c.Healthy {
		return nil, dperrors.New(
			errors.New("failed to get areas"),
			http.StatusNotFound,
			log.Data{"errors": map[string]string{"message": "404 Not Found: dataset not loaded in this server"}},
		)
	}
	if c.BadRequest {
		return nil, dperrors.New(
			errors.New("bad request"),
			http.StatusBadRequest,
			log.Data{"errors": map[string]string{"message": "400 Bad Request: dataset not loaded in this server"}},
		)
	}
	if c.NotFound {
		return nil, dperrors.New(
			errors.New("not found"),
			http.StatusNotFound,
			log.Data{"errors": map[string]string{"message": "404 Not Found: dataset not loaded in this server"}},
		)
	}

	// like is this below an accurate representation of its behaviour?
	if c.GetAreasResponse == nil {
		return &cantabular.GetAreasResponse{}, nil
	}

	return c.GetAreasResponse, nil
}

func (c *CantabularClient) GetArea(_ context.Context, _ cantabular.GetAreaRequest) (*cantabular.GetAreaResponse, error) {
	if !c.Healthy {
		return nil, dperrors.New(
			errors.New("failed to get area"),
			http.StatusNotFound,
			log.Data{"errors": map[string]string{"message": "404 Not Found: dataset not loaded in this server"}},
		)
	}
	if c.BadRequest {
		return nil, dperrors.New(
			errors.New("bad request"),
			http.StatusBadRequest,
			log.Data{"errors": map[string]string{"message": "400 Bad Request: dataset not loaded in this server"}},
		)
	}
	if c.NotFound {
		return nil, dperrors.New(
			errors.New("not found"),
			http.StatusNotFound,
			log.Data{"errors": map[string]string{"message": "404 Not Found: dataset not loaded in this server"}},
		)
	}

	if c.GetAreaResponse == nil {
		return &cantabular.GetAreaResponse{}, nil
	}

	return c.GetAreaResponse, nil
}

func (c *CantabularClient) GetParents(_ context.Context, _ cantabular.GetParentsRequest) (*cantabular.GetParentsResponse, error) {
	if !c.Healthy {
		return nil, dperrors.New(
			errors.New("test error response"),
			http.StatusNotFound,
			nil,
		)
	}
	return c.GetParentsResponse, nil
}
func (c *CantabularClient) GetParentAreaCount(_ context.Context, _ cantabular.GetParentAreaCountRequest) (*cantabular.GetParentAreaCountResult, error) {
	if !c.Healthy {
		return nil, dperrors.New(
			errors.New("test error response"),
			http.StatusNotFound,
			nil,
		)
	}

	return c.GetParentAreaCountResult, nil
}

func (c *CantabularClient) GetCategorisations(_ context.Context, _ cantabular.GetCategorisationsRequest) (*cantabular.GetCategorisationsResponse, error) {
	if !c.Healthy {
		return nil, dperrors.New(
			errors.New("test error response"),
			http.StatusNotFound,
			nil,
		)
	}
	return c.GetCategorisationsResponse, nil
}
