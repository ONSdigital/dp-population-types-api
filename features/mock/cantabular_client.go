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
	GetGeographyDimensionsResponse *cantabular.GetGeographyDimensionsResponse
	GetAreasResponse               *cantabular.GetAreasResponse
	GetParentsResponse             *cantabular.GetParentsResponse
	ListDatasetsResponse           []string
}

func (c *CantabularClient) Checker(ctx context.Context, state *healthcheck.CheckState) error {
	return nil
}

func (c *CantabularClient) ListDatasets(ctx context.Context) ([]string, error) {
	if !c.Healthy {
		return nil, errFailedToRespond
	}

	return c.ListDatasetsResponse, nil
}

func (c *CantabularClient) GetGeographyDimensions(ctx context.Context, _ cantabular.GetGeographyDimensionsRequest) (*cantabular.GetGeographyDimensionsResponse, error) {
	if !c.Healthy {
		return nil, dperrors.New(
			errors.New("error(s) returned by graphQL query"),
			http.StatusNotFound,
			log.Data{"errors": map[string]string{"message": "404 Not Found: dataset not loaded in this server"}},
		)
	}

	return c.GetGeographyDimensionsResponse, nil
}

func (c *CantabularClient) StatusCode(_ error) int {
	return http.StatusNotFound
}

func (c *CantabularClient) GetAreas(ctx context.Context, _ cantabular.GetAreasRequest) (*cantabular.GetAreasResponse, error) {
	if !c.Healthy {
		return nil, dperrors.New(
			errors.New("failed to get areas"),
			http.StatusNotFound,
			log.Data{"errors": map[string]string{"message": "404 Not Found: dataset not loaded in this server"}},
		)
	}
	if c.GetAreasResponse == nil {
		return &cantabular.GetAreasResponse{}, nil
	}

	return c.GetAreasResponse, nil
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
