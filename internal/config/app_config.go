// Package config provides the configuration for the application.
package config

import (
	"os"

	"github.com/pestanko/miniscrape/pkg/utils/applog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// AppConfig defines the main application configuration
type AppConfig struct {
	// Categories list of all available categories
	Categories []string `json:"categories"`
	// Cache configuration
	Cache CacheCfg `json:"cache"`
	// Web configuration
	Web WebCfg `json:"web"`
	// Log configuration
	Log applog.LogConfig `json:"log"`
	// Otel OpenTelemetry configuration
	Otel OtelConfig `json:"otel" yaml:"otel"`
}

// CacheCfg defines the configuration for the cache
type CacheCfg struct {
	// Enabled whether the cache is enabled
	Enabled bool `json:"enabled"`
	// Update whether the cache should be updated
	Update bool `json:"update"`
	// Root directory for the cache
	Root string `json:"root"`
}

// WebCfg web config
type WebCfg struct {
	// Addr where the server should be running
	Addr string `json:"addr" yaml:"addr"`
	// Users list of available users
	Users []User `json:"user" yaml:"user"`
}

// OtelConfig holds the OpenTelemetry configuration
type OtelConfig struct {
	Enabled  bool   `env:"OTEL_ENABLED,default=true" json:"enabled" yaml:"enabled"`
	Endpoint string `env:"OTEL_EXPORTER_OTLP_ENDPOINT,default=localhost:4317" json:"endpoint" yaml:"endpoint"`
	Protocol string `env:"OTEL_EXPORTER_OTLP_PROTOCOL,default=grpc" json:"protocol" yaml:"protocol"`
	Insecure bool   `env:"OTEL_INSECURE,default=true" json:"insecure" yaml:"insecure"`
}

// User definition in the system
type User struct {
	// Username of the user
	Username string `json:"username"`
	// Password of the user
	Password string `json:"password"`
}

// GetAppConfig - Unmarshal the app configuration using the viper
func GetAppConfig() *AppConfig {
	var config AppConfig

	if err := viper.Unmarshal(&config); err != nil {
		log.Info().
			Str("file", viper.ConfigFileUsed()).
			Err(err).
			Msg("Unable to load configuration")
	}

	log.Info().Interface("config", config).Msg("loaded config")

	if noCache := os.Getenv("APP_NO_CACHE"); noCache == "true" {
		log.Info().Msg("Cache is explictelly disabled")
		config.Cache.Enabled = false
	}

	return &config
}
