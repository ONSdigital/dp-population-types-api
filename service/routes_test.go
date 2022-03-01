package service_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/ONSdigital/dp-net/v2/responder"
	"github.com/ONSdigital/dp-population-types-api/config"
	"github.com/ONSdigital/dp-population-types-api/service"
)

func TestRoutes(t *testing.T) {
	Convey("Given a service was initialised", t, func() {
		svc := &service.Service{}
		initialiser := buildInitialiserMockWithNilDependencies()
		initialiser.GetResponderFunc = func() service.Responder { return responder.New() }
		config := config.Config{}
		err := svc.Init(context.Background(), &initialiser, &config, "", "", "")
		So(err, ShouldBeNil)

		Convey("Then the http server should respond to the population-types route", func() {
			rec := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/population-types", nil)
			svc.Router.ServeHTTP(rec, req)
			result := rec.Result()

			So(result.StatusCode, ShouldEqual, 200)

			actualResponse, _ := ioutil.ReadAll(rec.Body)
			result.Body.Close()

			So(string(actualResponse), ShouldEqualJSON, `{"items":[]}`)
		})
	})
}
