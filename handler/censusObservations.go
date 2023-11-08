package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

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

func (c *CensusObservations) toGetDatasetObservationsResponse(r io.Reader, ctx context.Context, w http.ResponseWriter) (string, error) {
	log.Info(ctx, "Starting to process response")

	buf := new(strings.Builder)

	writ, err := io.Copy(buf, r)
	if err != nil {
		fmt.Println(writ)
		fmt.Println("an error occurred")
		fmt.Println(err)
	}

	//var datasetResponse interface{}

	//c.respond.JSON(ctx, w, http.StatusOK, buf.String())
	return buf.String(), nil
}

func (c *CensusObservations) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	cancelContext, cancel := context.WithTimeout(ctx, time.Second*300)
	defer cancel()
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

	// // stream consumer/uploader for encrypted private files
	// consume = func(ctx context.Context, file io.Reader) error {
	// 	// if file == nil {
	// 	// 	return errors.New("no file content has been provided")
	// 	// }
	// 	// log.Info(ctx, "uploading encrypted private file to S3", logData)
	var consume cantabular.Consumer
	consume = func(ctx context.Context, file io.Reader) error {
		if file == nil {
			return errors.New("no file content has been provided")
		}

		response, err := c.toGetDatasetObservationsResponse(file, cancelContext, w)
		if err != nil {
			panic("That didn't go well")
		}
		fmt.Println("THE RESPONSE IS")
		if len(response) == 0 {
			fmt.Println("there was an error")
		}
		fmt.Println(response)

		// var trimmed string
		// trimmed := strings.TrimLeft(response, "{")
		// newtrimmed := strings.TrimLeft(trimmed, ":")
		// // //fmt.Println(newtrimmed)

		// if err != nil {
		// 	fmt.Println(err.Error())
		// }
		w.Write([]byte(response))
		//c.respond.JSON(ctx, w, http.StatusOK, response)

		//fmt.Println(data)
		return nil
	}

	//qRes, err := c.ctblr.StaticDatasetQuery(ctx, cReq)
	qRes, err := c.ctblr.StaticDatasetQueryStreamJson(cancelContext, cReq, consume)
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

	fmt.Println("THE COUNT IS ")
	fmt.Println(len(qRes.Observations))

	//c.respond.JSON(ctx, w, http.StatusOK, b)

}
