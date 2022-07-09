package config

import (
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

// AppConfig defines the main application configuration
type AppConfig struct {
	Categories []string
	Cache      CacheCfg
	Web        WebCfg
	Log        LogConfig
}

// LogConfig
type LogConfig struct {
	Dir                   string
	ConsoleLoggingEnabled bool
}

// CacheCfg defines the configuration for the cache
type CacheCfg struct {
	Enabled bool
	Update  bool
	Root    string
}

// WebCfg web config
type WebCfg struct {
	Addr  string `json:"addr" yaml:"addr"`
	Users []User `json:"user" yaml:"user"`
}

// User definition in the system
type User struct {
	Username string
	Password string
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
	return &config
}
