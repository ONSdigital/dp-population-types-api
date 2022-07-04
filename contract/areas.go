package contract

// GetAreasRequest defines the schema for the GET /areas query parameter
type GetAreasRequest struct {
	Dataset  string `schema:"dataset"`
	AreaType string `schema:"area-type"`
	Category string `schema:"q"`
}

// GetAreasResponse is the response object for GET /areas
type GetAreasResponse struct {
	Areas []Areas `json:"areas"`
}

// Areas is an area model with ID, Label, and area-type
type Areas struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	AreaType string `json:"area-type"`
}
