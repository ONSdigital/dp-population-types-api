package contract

type PopulationType struct {
	Name        string `json:"name"`
	Label       string `json:"label"`
	Description string `json:"description"`
}

type GetPopulationTypesRequest struct {
	DefaultDatasets bool `schema:"require-default-dataset"`
	QueryParams
}

type GetPopulationTypesResponse struct {
	PaginationResponse
	Items []PopulationType `json:"items"`
}

func (r *GetPopulationTypesResponse) Paginate() {
	endInt := r.Offset + r.Limit

	if r.Offset > len(r.Items) {
		r.Offset = 0
	}

	if r.Limit == 0 {
		r.Limit = defaultLimit
	}

	if endInt > len(r.Items) {
		endInt = len(r.Items)
	}

	r.Items = r.Items[r.Offset:endInt]
	r.Count = len(r.Items)
}
