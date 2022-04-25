package config

import (
	"github.com/spf13/viper"
	"log"
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
	Addr string
}

// GetAppConfig - Unmarshal the app configuration using the viper
func GetAppConfig() *AppConfig {
	var config AppConfig

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to load configuration (file used: %s): %v", viper.ConfigFileUsed(), err)
	}
	return &config
}
