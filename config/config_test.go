//nolint:typecheck
package config_test

import (
	"os"
	"reflect"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/ONSdigital/dp-population-types-api/config"
)

func TestConfig(t *testing.T) {
	os.Clearenv()

	Convey("Given an environment with no environment variables set", t, func() {

		Convey("When the config values are retrieved", func() {

			Convey("Then there should be no error returned, and values are as expected", func() {
				configuration, err := config.Get() // This Get() is only called once, when inside this function
				So(err, ShouldBeNil)
				So(configuration, ShouldResemble, &config.Config{
					BindAddr:                   "localhost:12900",
					GracefulShutdownTimeout:    5 * time.Second,
					HealthCheckInterval:        30 * time.Second,
					HealthCheckCriticalTimeout: 90 * time.Second,
					CantabularConfig: config.CantabularConfig{
						CantabularURL:         "http://localhost:8491",
						CantabularExtURL:      "http://localhost:8492",
						DefaultRequestTimeout: 10 * time.Second,
					},
				})
			})

		})
	})

	Convey("Configuration variables should be bound to the correct environment variables", t, func() {

		Convey("Top-level", func() {
			configMetadata := reflect.TypeOf(config.Config{})
			assertTagEnvConfig(configMetadata, "BindAddr", "BIND_ADDR")
			assertTagEnvConfig(configMetadata, "GracefulShutdownTimeout", "GRACEFUL_SHUTDOWN_TIMEOUT")
			assertTagEnvConfig(configMetadata, "HealthCheckInterval", "HEALTHCHECK_INTERVAL")
			assertTagEnvConfig(configMetadata, "HealthCheckCriticalTimeout", "HEALTHCHECK_CRITICAL_TIMEOUT")
		})

		Convey("Cantabular config", func() {
			cantabularConfigMetadata := reflect.TypeOf(config.CantabularConfig{})
			assertTagEnvConfig(cantabularConfigMetadata, "CantabularURL", "CANTABULAR_URL")
			assertTagEnvConfig(cantabularConfigMetadata, "CantabularExtURL", "CANTABULAR_EXT_API_URL")
			assertTagEnvConfig(cantabularConfigMetadata, "DefaultRequestTimeout", "DEFAULT_REQUEST_TIMEOUT")
		})
	})
}

func assertTagEnvConfig(configMetadata reflect.Type, fieldName string, expected string) {
	field, found := configMetadata.FieldByName(fieldName)
	So(found, ShouldBeTrue)
	So(field.Tag.Get("envconfig"), ShouldEqual, expected)
}
