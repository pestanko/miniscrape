package config

import (
	"github.com/spf13/viper"
	"log"
)

type AppConfig struct {
	Pages []Page `yaml:"pages"`
}

type Page struct {
	CodeName    string        `yaml:"codename"`
	Name        string        `yaml:"name"`
	Homepage    string        `yaml:"homepage"`
	Url         string        `yaml:"url"`
	Query       string        `yaml:"query"`
	CachePolicy string        `yaml:"cache_policy"`
	Filters     FiltersConfig `yaml:"filters"`
	Tags        []string      `yaml:"tags"`
}

type FiltersConfig struct {
	Cut CutFilter `yaml:"cut"`
	Day DayFilter `yaml:"day"`
}

type CutFilter struct {
	Before string `yaml:"before"`
	After  string `yaml:"after"`
}

type DayFilter struct {
	Days    []string `yaml:"days"`
	Enabled bool     `yaml:"enabled"`
}

// GetAppConfig - Unmarshal the app configuration using the viper
func GetAppConfig() *AppConfig {
	var config AppConfig

	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Unable to load configuration (file used: %s): %v", viper.ConfigFileUsed(), err)
	}
	return &config
}
