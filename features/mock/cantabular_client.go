package mock

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular/gql"

	dperrors "github.com/ONSdigital/dp-api-clients-go/v2/errors"
	"github.com/ONSdigital/dp-healthcheck/healthcheck"
	"github.com/ONSdigital/log.go/v2/log"
)

var (
	errFailedToRespond = errors.New("cantabular failed to respond")
)

type CantabularClient struct {
	Healthy                          bool
	BadRequest                       bool
	NotFound                         bool
	BadGateway                       bool
	ResponseTooLarge                 bool
	GetGeographyDimensionsResponse   *cantabular.GetGeographyDimensionsResponse
	GetDimensionsResponse            *cantabular.GetDimensionsResponse
	GetDimensionsDescriptionResponse *cantabular.GetDimensionsResponse
	GetAreasResponse                 *cantabular.GetAreasResponse
	GetAreaResponse                  *cantabular.GetAreaResponse
	GetParentsResponse               *cantabular.GetParentsResponse
	GetParentAreaCountResult         *cantabular.GetParentAreaCountResult
	GetCategorisationsResponse       *cantabular.GetCategorisationsResponse
	GetBaseVariableResponse          *cantabular.GetBaseVariableResponse
	GetDimensionCategoriesRespnse    *cantabular.GetDimensionCategoriesResponse
	ListDatasetsResponse             *cantabular.ListDatasetsResponse
	GetBlockedAreaCountResult        *cantabular.GetBlockedAreaCountResult
	GetStaticDatasetQuery            *cantabular.StaticDatasetQuery
	GetStaticDataset                 *gql.Dataset
	GetObservationsResponse          *cantabular.GetObservationsResponse
}

// CheckQueryCount implements service.CantabularClient.
func (c *CantabularClient) CheckQueryCount(context.Context, cantabular.StaticDatasetQueryRequest) (int, error) {
	if c.BadRequest {
		return 0, dperrors.New(
			errors.New("Maximum variables at MSOA and above is 5"),
			http.StatusBadRequest,
			log.Data{"errors": map[string]string{"message": "Maximum variables at MSOA and above is 5"}},
		)
	}
	if c.ResponseTooLarge {
		return 500000, nil
	}
	return c.GetObservationsResponse.TotalObservations, nil
}

func (c *CantabularClient) Checker(_ context.Context, _ *healthcheck.CheckState) error {
	return nil
}

func (c *CantabularClient) CheckerAPIExt(_ context.Context, _ *healthcheck.CheckState) error {
	return nil
}

func (c *CantabularClient) ListDatasets(_ context.Context) (*cantabular.ListDatasetsResponse, error) {
	if !c.Healthy {
		return nil, errFailedToRespond
	}

	return c.ListDatasetsResponse, nil
}

func (c *CantabularClient) GetDimensionCategories(_ context.Context, _ cantabular.GetDimensionCategoriesRequest) (
	*cantabular.GetDimensionCategoriesResponse, error) {
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
	if !c.Healthy {
		return nil, dperrors.New(
			errors.New("error(s) returned by graphQL query"),
			http.StatusNotFound,
			log.Data{"errors": map[string]string{"message": "404 Not Found: dataset not loaded in this server"}},
		)
	}
	return c.GetDimensionCategoriesRespnse, nil
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

func (c *CantabularClient) GetDimensionsDescription(_ context.Context, _ cantabular.GetDimensionsDescriptionRequest) (*cantabular.GetDimensionsResponse, error) {
	if !c.Healthy {
		return nil, dperrors.New(
			errors.New("error(s) returned by graphQL query"),
			http.StatusNotFound,
			log.Data{"errors": map[string]string{"message": "404 Not Found: dataset not loaded in this server"}},
		)
	}

	return c.GetDimensionsDescriptionResponse, nil
}

func (c *CantabularClient) StatusCode(_ error) int {
	if c.BadGateway {
		return http.StatusBadGateway
	}
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

	if c.GetAreasResponse == nil {
		return &cantabular.GetAreasResponse{}, nil
	}

	return c.GetAreasResponse, nil
}
func (c *CantabularClient) GetAreasTotalCount(_ context.Context, _ cantabular.GetAreasRequest) (int, error) {
	if !c.Healthy {
		return 0, dperrors.New(
			errors.New("failed to get areas"),
			http.StatusNotFound,
			log.Data{"errors": map[string]string{"message": "404 Not Found: dataset not loaded in this server"}},
		)
	}
	if c.BadRequest {
		return 0, dperrors.New(
			errors.New("bad request"),
			http.StatusBadRequest,
			log.Data{"errors": map[string]string{"message": "400 Bad Request: dataset not loaded in this server"}},
		)
	}
	if c.NotFound {
		return 0, dperrors.New(
			errors.New("not found"),
			http.StatusNotFound,
			log.Data{"errors": map[string]string{"message": "404 Not Found: dataset not loaded in this server"}},
		)
	}

	return c.GetAreasResponse.PaginationResponse.TotalCount, nil
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

func (c *CantabularClient) GetBlockedAreaCount(_ context.Context, _ cantabular.GetBlockedAreaCountRequest) (*cantabular.GetBlockedAreaCountResult, error) {
	if !c.Healthy {
		return nil, dperrors.New(
			errors.New("test error response"),
			http.StatusNotFound,
			nil,
		)
	}
	return c.GetBlockedAreaCountResult, nil
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

func (c *CantabularClient) GetBaseVariable(_ context.Context, _ cantabular.GetBaseVariableRequest) (*cantabular.GetBaseVariableResponse, error) {
	if !c.Healthy {
		return nil, dperrors.New(
			errors.New("test error response"),
			http.StatusNotFound,
			nil,
		)
	}

	if c.BadGateway {
		return nil, dperrors.New(
			errors.New("bad gateway"),
			http.StatusBadGateway,
			log.Data{"errors": map[string]string{"message": "variable at position 1 does not exist"}},
		)
	}

	if c.NotFound {
		return nil, dperrors.New(
			errors.New("not found"),
			http.StatusNotFound,
			log.Data{"errors": map[string]string{"message": "404 Not Found: dataset not loaded in this server"}},
		)
	}

	return c.GetBaseVariableResponse, nil
}

func (c *CantabularClient) StaticDatasetQuery(context.Context, cantabular.StaticDatasetQueryRequest) (*cantabular.StaticDatasetQuery, error) {
	if c.BadRequest {
		return nil, dperrors.New(
			errors.New("bad request"),
			http.StatusBadRequest,
			log.Data{"errors": map[string]string{"message": "400 Bad Request: codes not found for variable on filter for ltla"}},
		)
	}

	return c.GetStaticDatasetQuery, nil
}

// StaticDatasetQueryStreamJson implements service.CantabularClient.
func (c *CantabularClient) StaticDatasetQueryStreamJSON(context.Context, cantabular.StaticDatasetQueryRequest, func(ctx context.Context, r io.Reader) error) (cantabular.GetObservationsResponse, error) {
	if c.BadRequest {
		return *c.GetObservationsResponse, dperrors.New(
			errors.New("bad request"),
			http.StatusBadRequest,
			log.Data{"errors": map[string]string{"message": "400 Bad Request: codes not found for variable on filter for ltla"}},
		)
	}

	return *c.GetObservationsResponse, nil
}

func (c *CantabularClient) StaticDatasetType(ctx context.Context, datasetName string) (*gql.Dataset, error) {
	return c.GetStaticDataset, nil
}
