package handler_test

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/ONSdigital/dp-population-types-api/handler"
)

func TestEndpointRoot(t *testing.T) {

	Convey("Given a population-types handler", t, func() {

		responder := fakeResponder{}
		cantabularClient := fakeCantabularClient{
			bakedResponse: []string{"Thing one", "Thing two"},
		}
		handler := handler.NewPopulationTypes(&responder, &cantabularClient)

		Convey("When I get population types", func() {

			recorder := httptest.NewRecorder()
			req := httptest.NewRequest("the-method", "https://the-host/the-path", nil)
			handler.Get(recorder, req)

			expectedJSON := `{"items":[{"name": "Thing one"},{"name": "Thing two"}]}`

			result := recorder.Result()
			SoMsg("Then the response should be OK", result.StatusCode, ShouldEqual, http.StatusOK)

			actual, err := ioutil.ReadAll(result.Body)
			So(err, ShouldBeNil)
			SoMsg("And the response should match expected", string(actual), ShouldEqualJSON, expectedJSON)
		})
	})

}

type fakeResponder struct {
}

func (r *fakeResponder) JSON(ctx context.Context, w http.ResponseWriter, status int, resp interface{}) {
	w.WriteHeader(status)
	bytes, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	w.Write(bytes)
}

type fakeCantabularClient struct {
	bakedResponse []string
}

func (t *fakeCantabularClient) ListDatasets(ctx context.Context) ([]string, error) {
	return t.bakedResponse, nil
}
