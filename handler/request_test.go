package handler

import (
	"net/http"
	"testing"

	"github.com/ONSdigital/dp-population-types-api/contract"
	. "github.com/smartystreets/goconvey/convey"
)

func TestParseRequest(t *testing.T) {
	limit := 10
	Convey("Given a query param that is not in data structure", t, func() {
		req, err := http.NewRequest("GET", "http://localhost/areas?offset=4&limit=10&X=test", nil)
		So(err, ShouldBeNil)
		var ar contract.GetAreasRequest
		err = parseRequest(req, &ar)
		So(err, ShouldNotBeNil)
		So(ar, ShouldResemble, contract.GetAreasRequest{
			QueryParams: contract.QueryParams{
				Limit:  &limit,
				Offset: 4,
			},
			Category: "", //only populates if available
		})
	})

}
