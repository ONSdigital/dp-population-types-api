package contract

// AreaType is an area-type model with ID and Label
type AreaType struct {
	ID                      string `json:"id"`
	Label                   string `json:"label"`
	Description             string `json:"description"`
	TotalCount              int    `json:"total_count"`
	GeographyHierarchyOrder int    `json:"geography_hierarchy_order"`
}

type GetAreaTypesRequest struct {
	QueryParams
	PopulationType string
}

// GetAreaTypesResponse is the response object for GET /area-types
type GetAreaTypesResponse struct {
	PaginationResponse
	AreaTypes []AreaType `json:"items"`
}

type GetAreaTypeParentsRequest struct {
	QueryParams
	PopulationType string
	AreaType       string
}

type GetAreaTypeParentsResponse struct {
	PaginationResponse
	AreaTypes []AreaType `json:"items"`
}
