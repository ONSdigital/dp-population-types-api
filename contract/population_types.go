package contract

type PopulationType struct {
	Name string `json:"name"`
}

func NewPopulationType(name string) *PopulationType {
	return &PopulationType{Name: name}
}

type GetPopulationTypesResponse struct {
	PopulationTypes
}
type PopulationTypes struct {
	Items []PopulationType `json:"items"`
}

func NewPopulationTypes(names []string) PopulationTypes {
	items := make([]PopulationType, len(names))
	for i, name := range names {
		items[i] = *NewPopulationType(name)
	}
	return PopulationTypes{Items: items}
}
