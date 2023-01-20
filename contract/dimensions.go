package contract

import (
	"net/url"

	"github.com/pkg/errors"
)

// Dimension is an area-type model with ID and Label
type Dimension struct {
	ID                   string `json:"id"`
	Label                string `json:"label"`
	Description          string `json:"description"`
	TotalCount           int    `json:"total_count"`
	QualityStatementText string `json:"quality_statement_text"`
}

type GetDimensionsRequest struct {
	QueryParams
	PopulationType string
	SearchText     string `schema:"q"`
}

type GetDimensionsDescriptionRequest struct {
	QueryParams
	PopulationType string
	DimensionNames []string `schema:"q"`
}

// GetAreaTypesResponse is the response object for GET /dimensions
type GetDimensionsResponse struct {
	PaginationResponse
	Dimensions []Dimension `json:"items"`
}

type GetCategorisationsRequest struct {
	QueryParams
	PopulationType string
	Variable       string
}

func (r *GetDimensionsRequest) Valid() error {
	var err error
	if r.SearchText, err = url.QueryUnescape(r.SearchText); err != nil {
		return errors.New("invalid query string")
	}

	if err := r.QueryParams.Valid(); err != nil {
		return errors.Wrap(err, "invalid query pameters")
	}

	return nil
}

type GetCategorisationsResponse struct {
	PaginationResponse
	Items []Category `json:"items"`
}

type Category struct {
	Id         string              `json:"id"`
	Label      string              `json:"label"`
	Categories []DimensionCategory `json:"categories"`
}

type DimensionCategory struct {
	ID    string `json:"id"`
	Label string `json:"label"`
}

type GetBaseVariableRequest struct {
	PopulationType string
	Variable       string
}

type GetBaseVariableResponse struct {
	Name  string `json:"name"`
	Label string `json:"label"`
}

type GetDimensionCategoriesRequest struct {
	QueryParams
	PopulationType string
	Variables      string `schema:"dims"`
}

type GetDimensionCategoriesResponse struct {
	PaginationResponse
	Categories []Category `json:"items"`
}
