package handler

import "github.com/ONSdigital/dp-population-types-api/contract"

func filterPopulationTypes(coundition []string, toFilter []contract.PopulationType) []contract.PopulationType {
	items := make(map[string]int)
	filtered := make([]contract.PopulationType, 0)

	for _, item := range coundition {
		items[item] = 1
	}

	for _, item := range toFilter {
		if _, ok := items[item.Name]; ok {
			filtered = append(filtered, item)
		}
	}

	return filtered
}
