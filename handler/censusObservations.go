package handler

import (
	"net/http"
	"strings"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	"github.com/ONSdigital/dp-population-types-api/config"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

type GetObservationResponse struct {
	Dimensions  []ObservationDimension `bson:"dimensions"           json:"dimensions"`
	Observation float32                `bson:"observation"   json:"observation"`
}

type DatasetJSONLinks struct {
	Self model.Link `json:"self"`
}

type GetObservationsResponse struct {
	Observations      []GetObservationResponse `bson:"observations"           json:"observations"`
	Links             DatasetJSONLinks         `json:"links"`
	TotalObservations int                      `json:"total_observations"`
}

type ObservationDimension struct {
	Dimension   string `bson:"dimension"           json:"dimension"`
	DimensionID string `bson:"dimension_id"           json:"dimension_id"`
	Option      string `bson:"option"           json:"option"`
	OptionID    string `bson:"option_id"           json:"option_id"`
}

type GetFilterDimensionOptionsItem struct {
	Option    string     `json:"option"`
	Self      model.Link `json:"self"`
	Filter    model.Link `json:"filter"`
	Dimension model.Link `json:"Dimension"`
}

type CensusObservations struct {
	cfg     *config.Config
	respond responder
	ctblr   cantabularClient
}

// NewCensusObservations returns a new census observations
func NewCensusObservations(cfg *config.Config, r responder, c cantabularClient) *CensusObservations {
	return &CensusObservations{
		cfg:     cfg,
		respond: r,
		ctblr:   c,
	}
}

func getDimensionRow(query *cantabular.StaticDatasetQuery, dimIndices []int, dimIndex int) (value []ObservationDimension) {

	var observationDimensions []ObservationDimension

	for index, element := range dimIndices {
		dimension := query.Dataset.Table.Dimensions[index]

		observationDimensions = append(observationDimensions, ObservationDimension{
			Dimension:   dimension.Variable.Label,
			DimensionID: dimension.Variable.Name,
			Option:      dimension.Categories[element].Label,
			OptionID:    dimension.Categories[element].Code,
		})
	}

	return observationDimensions
}

func (c *CensusObservations) toGetDatasetObservationsResponse(query *cantabular.StaticDatasetQuery) (*GetObservationsResponse, error) {
	var getObservationResponse []GetObservationResponse

	dimLength := make([]int, 0)
	dimIndices := make([]int, 0)

	for k := 0; k < len(query.Dataset.Table.Dimensions); k++ {
		dimLength = append(dimLength, len(query.Dataset.Table.Dimensions[k].Categories))
		dimIndices = append(dimIndices, 0)
	}

	for v := 0; v < len(query.Dataset.Table.Values); v++ {
		dimension := getDimensionRow(query, dimIndices, v)
		getObservationResponse = append(getObservationResponse, GetObservationResponse{
			Dimensions:  dimension,
			Observation: query.Dataset.Table.Values[v],
		})

		i := len(dimIndices) - 1
		for i >= 0 {
			dimIndices[i] += 1
			if dimIndices[i] < dimLength[i] {
				break
			}
			dimIndices[i] = 0
			i -= 1
		}

	}

	var getObservationsResponse GetObservationsResponse
	getObservationsResponse.Observations = getObservationResponse
	getObservationsResponse.TotalObservations = len(query.Dataset.Table.Values)

	return &getObservationsResponse, nil
}

func (c *CensusObservations) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cReq := cantabular.StaticDatasetQueryRequest{
		Dataset:   chi.URLParam(r, "population-type"),
		Variables: strings.Split(r.URL.Query().Get("dimensions"), ","),
	}

	if r.URL.Query().Get("area-type") != "" {
		//the first value in this collection is variable and rest are codes
		vals := strings.Split(r.URL.Query().Get("area-type"), ",")
		cReq.Filters = []cantabular.Filter{{
			Variable: vals[0],
			Codes:    vals[1:],
		}}
	}

	//check if the dimensions has the area-type (variable) in it, else append
	//e.g. /population-types/UR/census-observations?dimensions=resident_age_7b&area-type=ltla,E06000008
	//in this case `dimensions=resident_age_7` is considered as `ltla,dimensions=resident_age_7`
	addVaraible := true
	for _, v := range cReq.Variables {
		if v == cReq.Filters[0].Variable {
			addVaraible = false
			break
		}
	}

	if addVaraible {
		cReq.Variables = append([]string{cReq.Filters[0].Variable}, cReq.Variables...)
	}

	logData := log.Data{
		"population_type": cReq.Dataset,
		"variables":       cReq.Variables,
		"filters":         cReq.Filters,
	}

	qRes, err := c.ctblr.StaticDatasetQuery(ctx, cReq)
	if err != nil {
		c.respond.Error(
			ctx,
			w,
			statusCode(err),
			Error{
				err:     err,
				logData: logData,
			},
		)
		return
	}

	response, err := c.toGetDatasetObservationsResponse(qRes)
	if err != nil {
		c.respond.Error(
			ctx,
			w,
			statusCode(err),
			Error{
				err:     errors.Wrap(err, "failed to generate response"),
				logData: logData,
			},
		)
	} else {
		//special handling for self link
		response.Links.Self.HREF = r.URL.String()
		c.respond.JSON(ctx, w, http.StatusOK, response)
	}
}
