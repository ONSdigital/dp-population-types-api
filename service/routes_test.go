package service_test

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maxcnunes/httpfake"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular/gql"
	"github.com/ONSdigital/dp-net/v2/responder"
	"github.com/ONSdigital/dp-population-types-api/config"
	"github.com/ONSdigital/dp-population-types-api/service"
	"github.com/ONSdigital/dp-population-types-api/service/mock"
)

func TestRoutes(t *testing.T) {
	Convey("Given a service was initialised", t, func() {
		svc := &service.Service{}
		initialiser := buildInitialiserMockWithNilDependencies()
		initialiser.GetResponderFunc = func() service.Responder { return responder.New() }
		initialiser.GetCantabularClientFunc = testCantabularClient
		config := config.Config{}

		fakeServer := httpfake.New()
		config.DatasetAPIURL = fakeServer.ResolveURL("")

		url := fmt.Sprintf(
			`/datasets?offset=0&limit=100&is_based_on=%s`,
			"dataset-id",
		)
		fakeServer.NewHandler().
			Get(url).
			Reply(http.StatusOK).
			BodyString(`{"total_count": 1}`)

		err := svc.Init(context.Background(), &initialiser, &config, "", "", "")
		So(err, ShouldBeNil)

		Convey("Then the http server should respond to the population-types route", func() {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/population-types", nil)
			svc.Router.ServeHTTP(rec, req)
			result := rec.Result()

			So(result.StatusCode, ShouldEqual, http.StatusOK)

			actualResponse, _ := ioutil.ReadAll(rec.Body)
			result.Body.Close()

			So(string(actualResponse), ShouldEqualJSON, `{"items":[]}`)
		})

		Convey("Then the http server should respond to the area-types route", func() {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/population-types/dataset-id/area-types", nil)
			svc.Router.ServeHTTP(rec, req)
			result := rec.Result()

			So(result.StatusCode, ShouldEqual, http.StatusOK)

			actualResponse, _ := ioutil.ReadAll(rec.Body)
			result.Body.Close()

			So(string(actualResponse), ShouldEqualJSON, `{"area-types":null}`)
		})
	})
}

func testCantabularClient(config.CantabularConfig) service.CantabularClient {
	return &mock.CantabularClientMock{
		ListDatasetsFunc: func(ctx context.Context) ([]string, error) {
			return []string{}, nil
		},
		GetGeographyDimensionsFunc: func(ctx context.Context, req cantabular.GetGeographyDimensionsRequest) (*cantabular.GetGeographyDimensionsResponse, error) {
			return &cantabular.GetGeographyDimensionsResponse{Dataset: gql.DatasetRuleBase{}}, nil
		},
	}
}
