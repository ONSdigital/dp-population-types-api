package contract

type PopulationType struct {
	Name string `json:"name"`
}

type PopulationTypes struct {
	Items []PopulationType `json:"items"`
}

type GetPopulationTypesRequest struct {
	QueryParams
}

func NewPopulationTypes(names []string) PopulationTypes {
	items := make([]PopulationType, len(names))
	for i, name := range names {
		items[i] = PopulationType{name}
	}
	return PopulationTypes{Items: items}
}

func (r *GetPopulationTypesRequest) Paginate(types []string) *PopulationTypes {
	endInt := r.Offset + *r.Limit

	if r.Offset > len(types) {
		r.Offset = 0
	}

	if r.Limit == nil {
		r.Limit = &defaultLimit
	}

	if endInt > len(types) {
		endInt = len(types)
	}

	subset := types[r.Offset:endInt]

	response := NewPopulationTypes(subset)

	return &response
}

type GetPopulationTypesResponse struct {
	PaginationResponse
	PopulationTypes
}
