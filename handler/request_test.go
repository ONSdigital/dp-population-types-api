package handler

import (
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParseRequest(t *testing.T) {

	Convey("Given a valid query param", t, func() {
		req, err := http.NewRequest("GET", "http://localhost/areas?offset=4&limit=10&q=test", nil)
		So(err, ShouldBeNil)
		var ar areaRequest
		parseRequest(req, &ar)
		So(ar, ShouldResemble, areaRequest{
			Limit:    10,
			Offset:   4,
			Category: "test",
		})

		Convey("Given the limit is set to 0 we get dafault", func() {
			req, err := http.NewRequest("GET", "http://localhost/areas?offset=4&limit=0&q=test", nil)
			So(err, ShouldBeNil)
			var ar areaRequest
			parseRequest(req, &ar)
			So(ar, ShouldResemble, areaRequest{
				Limit:    20,
				Offset:   4,
				Category: "test",
			})
		})
	})

	Convey("Given a query param that is not in data structure", t, func() {
		req, err := http.NewRequest("GET", "http://localhost/areas?offset=4&limit=10&X=test", nil)
		So(err, ShouldBeNil)
		var ar areaRequest
		parseRequest(req, &ar)
		So(ar, ShouldResemble, areaRequest{
			Limit:    10,
			Offset:   4,
			Category: "", //only populates if available
		})
	})

}
