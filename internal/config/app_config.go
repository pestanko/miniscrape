package config

import (
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
	Log LogConfig `json:"log"`
}

// LogConfig logger configuration
type LogConfig struct {
	// Dir where to store log files
	Dir string `json:"dir"`
	// ConsoleLoggingEnabled whether logger should use console logging
	ConsoleLoggingEnabled bool `json:"console_logging_enabled"`
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

	return &config
}
