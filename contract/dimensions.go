package contract

// Dimension is an area-type model with ID and Label
type Dimension struct {
	ID         string `json:"id"`
	Label      string `json:"label"`
	TotalCount int    `json:"total_count"`
}

type GetDimensionsRequest struct {
	QueryParams
	PopulationType string
	SearchText     string `schema:"q"`
}

// GetAreaTypesResponse is the response object for GET /dimensions
type GetDimensionsResponse struct {
	PaginationResponse
	Dimensions []Dimension `json:"dimensions"`
}
