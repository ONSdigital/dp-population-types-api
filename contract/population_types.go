package contract

type PopulationType struct {
	Name string `json:"name"`
}

type PopulationTypes struct {
	Items []PopulationType `json:"items"`
}

func NewPopulationTypes(names []string) PopulationTypes {
	items := make([]PopulationType, len(names))
	for i, name := range names {
		items[i] = PopulationType{name}
	}
	return PopulationTypes{Items: items}
}

type GetPopulationTypesResponse struct {
	PaginationResponse
	PopulationTypes
}
