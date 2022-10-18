package handler

import (
	"strconv"

	"github.com/ONSdigital/dp-population-types-api/contract"
)

func manualPagination(limit, offset string, populationTypes []string) (contract.PopulationTypes, int, int) {

	// no need to check for errors as
	// if non integer, defaults to 0
	limitInt, _ := strconv.Atoi(limit)
	offsetInt, _ := strconv.Atoi(offset)

	if offsetInt > len(populationTypes) {
		offsetInt = 0
	}

	if limitInt == 0 {
		limitInt = 20
	}
	if limitInt > len(populationTypes) {
		limitInt = len(populationTypes)
	}
	println(limitInt)
	println(offsetInt)
	subset := populationTypes[offsetInt:limitInt]

	return contract.NewPopulationTypes(subset), limitInt, offsetInt
}
