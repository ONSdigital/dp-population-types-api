package config

import (
	"time"

	mongo "github.com/ONSdigital/dp-mongodb/v3/mongodb"
	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-population-types-api
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	ServiceAuthToken           string        `envconfig:"SERVICE_AUTH_TOKEN"    json:"-"`
	EnablePrivateEndpoints     bool          `envconfig:"ENABLE_PRIVATE_ENDPOINTS"`
	EnablePermissionsAuth      bool          `envconfig:"ENABLE_PERMISSIONS_AUTH"`
	ZebedeeURL                 string        `envconfig:"ZEBEDEE_URL"`
	DatasetAPIURL              string        `envconfig:"DATASET_API_URL"`
	MetadataCollection         string        `envconfig:"METADATA_COLLECTION"`
	Mongo                      mongo.MongoDriverConfig
	CantabularConfig
}

type CantabularConfig struct {
	CantabularURL                string        `envconfig:"CANTABULAR_URL"`
	CantabularExtURL             string        `envconfig:"CANTABULAR_API_EXT_URL"`
	DefaultRequestTimeout        time.Duration `envconfig:"DEFAULT_REQUEST_TIMEOUT"`
	CantabularHealthcheckEnabled bool          `envconfig:"CANTABULAR_HEALTHCHECK_ENABLED"`
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		BindAddr:                   "localhost:27300",
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		EnablePrivateEndpoints:     false,
		ZebedeeURL:                 "http://localhost:8082",
		ServiceAuthToken:           "",
		EnablePermissionsAuth:      true,
		DatasetAPIURL:              "http://localhost:22000",
		CantabularConfig: CantabularConfig{
			CantabularURL:                "http://localhost:8491",
			CantabularExtURL:             "http://localhost:8492",
			DefaultRequestTimeout:        10 * time.Second,
			CantabularHealthcheckEnabled: false,
		},
		MetadataCollection: "defaultDatasetMetadata",
		Mongo: mongo.MongoDriverConfig{
			ClusterEndpoint: "localhost:27017",
			Username:        "",
			Password:        "",
			Database:        "filters",
			Collections: map[string]string{
				"defaultDatasetMetadata": "defaultDatasetMetadata",
			},
			ReplicaSet:                    "",
			IsStrongReadConcernEnabled:    false,
			IsWriteConcernMajorityEnabled: true,
			ConnectTimeout:                5 * time.Second,
			QueryTimeout:                  15 * time.Second,
			TLSConnectionConfig: mongo.TLSConnectionConfig{
				IsSSL: false,
			},
		},
	}

	return cfg, envconfig.Process("", cfg)
}
