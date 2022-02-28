package service_test

import (
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	dphttp "github.com/ONSdigital/dp-net/http"
	"github.com/ONSdigital/dp-population-types-api/config"
	"testing"
	"time"

	"github.com/ONSdigital/dp-population-types-api/service"
	. "github.com/smartystreets/goconvey/convey"
)

func TestInitialiser(t *testing.T) {

	Convey("Given an initialiser", t, func() {

		var actualConfig cantabular.Config
		var actualUserAgent dphttp.Clienter

		initialiser := service.NewInit()
		initialiser.CantabularClientFactory = func(cfg cantabular.Config, ua dphttp.Clienter) *cantabular.Client {
			actualConfig = cfg
			actualUserAgent = ua
			return nil
		}

		Convey("When a cantabular client is built", func() {
			cantabularConfig := config.CantabularConfig{
				CantabularURL:         "CantabularURL",
				CantabularExtURL:      "CantabularExtURL",
				DefaultRequestTimeout: 42 * time.Hour,
			}
			_ = initialiser.GetCantabularClient(cantabularConfig)
			Convey("Then the cantabular client factory should be called with the expectd configuration and user agent", func() {
				So(actualConfig, ShouldResemble, cantabular.Config{
					Host:           "CantabularURL",
					ExtApiHost:     "CantabularExtURL",
					GraphQLTimeout: 42 * time.Hour,
				})
				So(actualUserAgent, ShouldNotBeNil)
			})
		})
	})
}
