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
	Category    string        `yaml:"category"`
	Tags        []string      `yaml:"tags"`
	Filters     FiltersConfig `yaml:"filters"`
}

type FiltersConfig struct {
	Cut     CutFilter     `yaml:"cut"`
	CutLine CutLineFilter `yaml:"cutLine"`
	Day     DayFilter     `yaml:"day"`
}

type CutFilter struct {
	Before string `yaml:"before"`
	After  string `yaml:"after"`
}

type CutLineFilter struct {
	StartsWith string `yaml:"startsWith"`
	Contains   string `yaml:"contains"`
	CutAfter   string `yaml:"cutAfter"`
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
