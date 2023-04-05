package handler

import (
	"context"
	"net/http"
	"strings"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-cantabular-filter-flex-api/model"
	"github.com/ONSdigital/dp-population-types-api/config"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

const DATASET_TYPE_MICRODATA = "microdata"

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
	BlockedAreas      int                      `json:"blocked_areas"`
	TotalAreas        int                      `json:"total_areas"`
	AreasReturned     int                      `json:"areas_returned"`
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

func (c *CensusObservations) toGetDatasetObservationsResponse(query *cantabular.StaticDatasetQuery, ctx context.Context) (*GetObservationsResponse, error) {
	log.Info(ctx, "Starting to process response", log.Data{"area-type": query.Dataset.Table.Dimensions[0].Variable.Name})

	var getObservationResponse []GetObservationResponse

	dimLength := make([]int, 0)
	dimIndices := make([]int, 0)

	for k := 0; k < len(query.Dataset.Table.Dimensions); k++ {
		dimLength = append(dimLength, len(query.Dataset.Table.Dimensions[k].Categories))
		dimIndices = append(dimIndices, 0)
	}

	log.Info(ctx, "Created the arrays to hold dimension and categorisation information.  About to begin the processing loop for the results.", log.Data{"area-type": query.Dataset.Table.Dimensions[0].Variable.Name})

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

	log.Info(ctx, "Process observation response", log.Data{"observation-response-size": len(getObservationResponse)})

	var getObservationsResponse GetObservationsResponse
	getObservationsResponse.Observations = getObservationResponse
	getObservationsResponse.TotalObservations = len(query.Dataset.Table.Values)

	getObservationsResponse.BlockedAreas = query.Dataset.Table.Rules.Blocked.Count
	getObservationsResponse.TotalAreas = query.Dataset.Table.Rules.Total.Count
	getObservationsResponse.AreasReturned = query.Dataset.Table.Rules.Total.Count

	log.Info(ctx, "Complete response processed, about to return to user", log.Data{"total-response-length": len(getObservationsResponse.Observations)})

	return &getObservationsResponse, nil
}

func (c *CensusObservations) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	logData := log.Data{
		"method": http.MethodGet,
	}

	cReq := cantabular.StaticDatasetQueryRequest{
		Dataset:   chi.URLParam(r, "population-type"),
		Variables: strings.Split(r.URL.Query().Get("dimensions"), ","),
	}

	log.Info(ctx, "handling census-observations request", log.Data{"request": cReq})

	//check if the dataset is of type microdata
	dataset, err := c.ctblr.StaticDatasetType(ctx, cReq.Dataset)
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

	if dataset.Type != DATASET_TYPE_MICRODATA {
		c.respond.Error(
			ctx,
			w,
			http.StatusBadRequest,
			Error{
				err:     errors.New("Only supports dataset of type microdata"),
				logData: logData,
			},
		)
		return
	}

	areaType := ""
	if r.URL.Query().Get("area-type") != "" {
		//the first value in this collection is variable and rest are codes
		vals := strings.Split(r.URL.Query().Get("area-type"), ",")
		areaType = vals[0]

		if len(vals[1:]) > 0 { //only populate filters if codes are available
			cReq.Filters = []cantabular.Filter{{
				Variable: vals[0],
				Codes:    vals[1:],
			}}
		}
	}

	log.Info(ctx, "census-observations - determined area-type and areas values", log.Data{"updated-request": cReq})

	//check if the dimensions has the area-type (variable) in it, else append
	//e.g. /population-types/UR/census-observations?dimensions=resident_age_7b&area-type=ltla,E06000008
	//in this case `dimensions=resident_age_7` is considered as `ltla,dimensions=resident_age_7`
	addVaraible := true
	for _, v := range cReq.Variables {
		if v == areaType {
			addVaraible = false
			break
		}
	}

	if addVaraible {
		cReq.Variables = append([]string{areaType}, cReq.Variables...)
	}

	logData["population_type"] = cReq.Dataset
	logData["variables"] = cReq.Variables
	logData["filters"] = cReq.Filters

	log.Info(ctx, "handling census-observations - all parameters now set.  Sending query to cantabular", logData)

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
	if len(qRes.Dataset.Table.Error) > 0 {
		c.respond.Error(
			ctx,
			w,
			http.StatusBadRequest,
			Error{
				err:     errors.New(handleError(CantabularError(qRes.Dataset.Table.Error))),
				logData: logData,
			},
		)
		return
	}

	log.Info(ctx, "Response received from Cantabular - no errors have been returned", log.Data{"value-count": len(qRes.Dataset.Table.Values)})

	response, err := c.toGetDatasetObservationsResponse(qRes, ctx)
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
	}

	//special handling for self link
	response.Links.Self.HREF = r.URL.String()
	c.respond.JSON(ctx, w, http.StatusOK, response)

}
