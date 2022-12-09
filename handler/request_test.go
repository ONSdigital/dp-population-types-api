package handler

import (
	"net/http"
	"testing"

	"github.com/ONSdigital/dp-population-types-api/contract"
	. "github.com/smartystreets/goconvey/convey"
)

func TestParseRequest(t *testing.T) {

	Convey("Given a query param that is not in data structure", t, func() {
		req, err := http.NewRequest("GET", "http://localhost/areas?offset=4&limit=10&X=test", nil)
		So(err, ShouldBeNil)
		var ar contract.GetAreasRequest
		err = parseRequest(req, &ar)
		So(err, ShouldNotBeNil)
		So(ar, ShouldResemble, contract.GetAreasRequest{
			QueryParams: contract.QueryParams{
				Limit:  10,
				Offset: 4,
			},
			Category: "", //only populates if available
		})
	})

	Convey("Given a dimensions request with search query param", t, func() {
		req, err := http.NewRequest("GET", "http://localhost/population-types/dummy_data_households/dimensions?limit=1000&offset=0&q=household%2Bdeprivation", nil)
		So(err, ShouldBeNil)
		var ar contract.GetDimensionsRequest
		err = parseRequest(req, &ar)
		So(err, ShouldBeNil)
		So(ar, ShouldResemble, contract.GetDimensionsRequest{
			QueryParams: contract.QueryParams{
				Limit:  1000,
				Offset: 0,
			},
			SearchText: "household deprivation", //only populates if available
		})
	})
}
