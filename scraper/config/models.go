package config

import (
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path/filepath"
)

type AppConfig struct {
	Categories []string
	Cache      CacheCfg
}

type CacheCfg struct {
	Enabled bool
	Update  bool
	Root    string
}

type Category struct {
	Pages []Page `yaml:"pages"`
	Name  string `yaml:"name"`
}

type Page struct {
	CodeName    string        `yaml:"codename"`
	Name        string        `yaml:"name"`
	Homepage    string        `yaml:"homepage"`
	Url         string        `yaml:"url"`
	Query       string        `yaml:"query"`
	CachePolicy string        `yaml:"cache_policy"`
	Resolver    string        `yaml:"resolver"`
	Category    string        `yaml:"category"`
	Disabled    bool          `yaml:"disabled"`
	Tags        []string      `yaml:"tags"`
	Filters     FiltersConfig `yaml:"filters"`
}

type HtmlFilter struct {
	PrettyTables bool `yaml:"pretty_tables"`
	TextOnly     bool `yaml:"text_only"`
}

type FiltersConfig struct {
	Cut     CutFilter     `yaml:"cut"`
	CutLine CutLineFilter `yaml:"cutLine"`
	Day     DayFilter     `yaml:"day"`
	Html    HtmlFilter    `yaml:"html"`
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

func LoadCategories(cfg *AppConfig) []Category {
	cfgPathUsed := viper.ConfigFileUsed()
	baseDir := filepath.Dir(cfgPathUsed)
	var categories []Category

	for _, catName := range cfg.Categories {
		ok, cat := loadCategoryFile(baseDir, catName)
		if ok {
			categories = append(categories, cat)
		}
	}

	return categories
}

func loadCategoryFile(baseDir string, catName string) (bool, Category) {
	fp := filepath.Join(baseDir, catName+".yml")
	log.Printf("Loading file: %s", fp)

	content, err := ioutil.ReadFile(fp)
	if err != nil {
		log.Printf("Unable to open file \"%s\": %v", fp, err)
		return false, Category{}
	}

	var cat Category
	if err = yaml.Unmarshal(content, &cat); err != nil {
		log.Printf("Unable to load file \"%s\": %v", fp, err)
		return false, Category{}
	}

	if cat.Name == "" {
		cat.Name = catName
	}

	// Normalize the pages
	for idx := range cat.Pages {
		if cat.Pages[idx].Category == "" {
			cat.Pages[idx].Category = cat.Name
		}
		if cat.Pages[idx].Resolver == "" {
			cat.Pages[idx].Resolver = "default"
		}
	}

	return true, cat
}
