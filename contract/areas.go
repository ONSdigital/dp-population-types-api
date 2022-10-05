package contract

import (
	"net/url"

	"github.com/pkg/errors"
)

// Area is an area model with ID, Label, and area-type
type Area struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	AreaType string `json:"area_type"`
}

// GetAreasRequest defines the schema for the GET /areas query parameter
type GetAreasRequest struct {
	QueryParams
	Category string `schema:"q"`
}

func (r *GetAreasRequest) Valid() error {
	var err error
	if r.Category, err = url.QueryUnescape(r.Category); err != nil {
		return errors.New("invalid query string")
	}

	if err := r.QueryParams.Valid(); err != nil {
		return errors.Wrap(err, "invalid query pameters")
	}

	return nil
}

// GetAreasResponse is the response object for GET /areas
type GetAreasResponse struct {
	PaginationResponse
	Areas []Area `json:"items"`
}

// GetAreaResponse is the response object for GET /areas
type GetAreaResponse struct {
	Area Area `json:"area"`
}
