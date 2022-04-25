package config

import (
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path/filepath"
)

type Category struct {
	Pages []Page `yaml:"pages" json:"pages"`
	Name  string `yaml:"name" json:"name"`
}

type Page struct {
	CodeName    string        `yaml:"codename" json:"code_name"`
	Name        string        `yaml:"name" json:"name"`
	Homepage    string        `yaml:"homepage" json:"homepage"`
	Url         string        `yaml:"url" json:"url"`
	Query       string        `yaml:"query" json:"query"`
	XPath       string        `yaml:"xpath" json:"xpath"`
	CachePolicy string        `yaml:"cache_policy" json:"cache_policy"`
	Resolver    string        `yaml:"resolver" json:"resolver"`
	Category    string        `yaml:"category" json:"category"`
	Disabled    bool          `yaml:"disabled" json:"disabled"`
	Tags        []string      `yaml:"tags" json:"tags"`
	Filters     FiltersConfig `yaml:"filters" json:"filters"`
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
