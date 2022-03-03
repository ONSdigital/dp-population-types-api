package handler_test

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/ONSdigital/dp-population-types-api/handler"
)

func TestEndpointRoot(t *testing.T) {

	Convey("Given a population-types handler", t, func() {

		responder := fakeResponder{}
		cantabularClient := fakeCantabularClient{
			listDatasetsReturnData: []string{"Thing one", "Thing two"},
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

		Convey("But the cantabular client is returning errors", func() {
			cantabularClient.listDatasetsReturnError = errors.New("oh no")

			Convey("When I get population types", func() {

				recorder := httptest.NewRecorder()
				req := httptest.NewRequest("the-method", "https://the-host/the-path", nil)
				handler.Get(recorder, req)

				result := recorder.Result()
				SoMsg("Then the response should be InternalServerError", result.StatusCode, ShouldEqual, http.StatusInternalServerError)

				actual, err := ioutil.ReadAll(result.Body)
				So(err, ShouldBeNil)
				SoMsg("And the response should reflect the top-level error message",
					strings.Contains(string(actual), "failed to fetch population types"), ShouldBeTrue)
				SoMsg("And the response should reflect the error returned by the cantabular client",
					strings.Contains(string(actual), cantabularClient.listDatasetsReturnError.Error()), ShouldBeTrue)
			})
		})
	})

}

type fakeResponder struct {
}

func (r *fakeResponder) Error(_ context.Context, w http.ResponseWriter, i int, err error) {
	w.WriteHeader(i)
	w.Write([]byte(err.Error()))
}

func (r *fakeResponder) JSON(_ context.Context, w http.ResponseWriter, status int, resp interface{}) {
	w.WriteHeader(status)
	bytes, err := json.Marshal(resp)
	if err != nil {
		panic(err)
	}
	w.Write(bytes)
}

type fakeCantabularClient struct {
	listDatasetsReturnData  []string
	listDatasetsReturnError error
}

func (t *fakeCantabularClient) ListDatasets(ctx context.Context) ([]string, error) {
	return t.listDatasetsReturnData, t.listDatasetsReturnError
}
