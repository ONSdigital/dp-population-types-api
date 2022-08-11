package contract

// Area is an area model with ID, Label, and area-type
type Area struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	AreaType string `json:"area_type"`
}

// GetAreasRequest defines the schema for the GET /areas query parameter
type GetAreasRequest struct {
	Dataset  string `schema:"dataset"`
	AreaType string `schema:"area_type"`
	Category string `schema:"q"`
}

// GetAreasResponse is the response object for GET /areas
type GetAreasResponse struct {
	Areas []Area `json:"areas"`
}

// GetAreaResponse is the response object for GET /areas
type GetAreaResponse struct {
	Area Area `json:"area"`
}
