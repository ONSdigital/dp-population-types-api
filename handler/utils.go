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

type CantabularError string

const (
	MaxVariableError        = "Maximum variables at MSOA and above is 5"
	MaxVariableGraphqlError = "Maximum variables in query is 5"
)

func handleError(errString CantabularError) string {
	switch errString {
	case MaxVariableError, MaxVariableGraphqlError:
		return "More than 5 variables selected, query failed"
	default:
		return "Unexpected error"
	}
}
