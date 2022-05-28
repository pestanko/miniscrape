package config

import (
	"log"

	"github.com/spf13/viper"
)

// AppConfig defines the main application configuration
type AppConfig struct {
	Categories []string
	Cache      CacheCfg
	Web        WebCfg
}

// CacheCfg defines the configuration for the cache
type CacheCfg struct {
	Enabled bool
	Update  bool
	Root    string
}

// WebCfg web config
type WebCfg struct {
	Addr string `json:"addr" yaml:"addr"`
}

// GetAppConfig - Unmarshal the app configuration using the viper
func GetAppConfig() *AppConfig {
	var config AppConfig

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to load configuration (file used: %s): %v", viper.ConfigFileUsed(), err)
	}
	return &config
}
