//nolint:typecheck
package config_test

import (
	"reflect"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/ONSdigital/dp-mongodb/v3/mongodb"
	"github.com/ONSdigital/dp-population-types-api/config"
)

func TestConfig(t *testing.T) {

	Convey("Given an environment with no environment variables set", t, func() {

		Convey("When the config values are retrieved", func() {

			Convey("Then there should be no error returned, and values are as expected", func() {
				configuration, err := config.Get() // This Get() is only called once, when inside this function
				// set as a env var, setting manually so not clashing.
				configuration.ServiceAuthToken = ""

				So(err, ShouldBeNil)
				So(configuration, ShouldResemble, &config.Config{
					BindAddr:                   "localhost:27300",
					GracefulShutdownTimeout:    5 * time.Second,
					HealthCheckInterval:        30 * time.Second,
					HealthCheckCriticalTimeout: 90 * time.Second,
					EnablePrivateEndpoints:     false,
					ZebedeeURL:                 "http://localhost:8082",
					EnablePermissionsAuth:      true,
					DatasetAPIURL:              "http://localhost:22000",
					OTExporterOTLPEndpoint:     "localhost:4317",
					OTServiceName:              "dp-population-types-api",
					OTBatchTimeout:				5 * time.Second,
					CantabularConfig: config.CantabularConfig{
						CantabularURL:                "http://localhost:8491",
						CantabularExtURL:             "http://localhost:8492",
						DefaultRequestTimeout:        10 * time.Second,
						CantabularHealthcheckEnabled: false,
					},
					MetadataCollection: "defaultDatasetMetadata",
					Mongo: mongodb.MongoDriverConfig{
						Username:        "",
						Password:        "",
						ClusterEndpoint: "localhost:27017",
						Database:        "filters",
						Collections: map[string]string{
							"defaultDatasetMetadata": "defaultDatasetMetadata",
						},
						ReplicaSet:                    "",
						IsStrongReadConcernEnabled:    false,
						IsWriteConcernMajorityEnabled: true,
						ConnectTimeout:                time.Duration(5000000000),
						QueryTimeout:                  time.Duration(15000000000),
						TLSConnectionConfig: mongodb.TLSConnectionConfig{
							IsSSL:              false,
							VerifyCert:         false,
							CACertChain:        "",
							RealHostnameForSSH: "",
						},
					},
					CensusObservationsFF: false,
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
			assertTagEnvConfig(cantabularConfigMetadata, "CantabularExtURL", "CANTABULAR_API_EXT_URL")
			assertTagEnvConfig(cantabularConfigMetadata, "DefaultRequestTimeout", "DEFAULT_REQUEST_TIMEOUT")
		})
	})
}

func assertTagEnvConfig(configMetadata reflect.Type, fieldName string, expected string) {
	field, found := configMetadata.FieldByName(fieldName)
	So(found, ShouldBeTrue)
	So(field.Tag.Get("envconfig"), ShouldEqual, expected)
}
