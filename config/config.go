package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-population-types-api
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	EnablePrivateEndpoints     bool          `envconfig:"ENABLE_PRIVATE_ENDPOINTS"`

	EnablePermissionsAuth bool   `envconfig:"ENABLE_PERMISSIONS_AUTH"`
	ZebedeeURL            string `envconfig:"ZEBEDEE_URL"`
	ServiceAuthToken      string `envconfig:"SERVICE_AUTH_TOKEN"`
	DatasetAPIURL         string `envconfig:"DATASET_API_URL"`
	CantabularConfig
}

type CantabularConfig struct {
	CantabularURL                string        `envconfig:"CANTABULAR_URL"`
	CantabularExtURL             string        `envconfig:"CANTABULAR_EXT_API_URL"`
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
		BindAddr:                   "localhost:12900",
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		EnablePrivateEndpoints:     true,
		ZebedeeURL:                 "http://localhost:8082",
		ServiceAuthToken:           "",
		EnablePermissionsAuth:      true,
		DatasetAPIURL:              "http://localhost:22000",
		CantabularConfig: CantabularConfig{
			CantabularURL:         "http://localhost:8491",
			CantabularExtURL:      "http://localhost:8492",
			DefaultRequestTimeout: 10 * time.Second,
		},
	}

	return cfg, envconfig.Process("", cfg)
}
