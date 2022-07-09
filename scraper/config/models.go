package config

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

type Category struct {
	Pages []Page `yaml:"pages" json:"pages"`
	Name  string `yaml:"name" json:"name"`
}

type Page struct {
	CodeName    string         `yaml:"codename" json:"codename"`
	Name        string         `yaml:"name" json:"name"`
	Homepage    string         `yaml:"homepage" json:"homepage"`
	Url         string         `yaml:"url" json:"url"`
	Query       string         `yaml:"query" json:"query"`
	XPath       string         `yaml:"xpath" json:"xpath"`
	CachePolicy string         `yaml:"cache_policy" json:"cachePolicy"`
	Resolver    string         `yaml:"resolver" json:"resolver"`
	Category    string         `yaml:"category" json:"category"`
	Disabled    bool           `yaml:"disabled" json:"disabled"`
	Tags        []string       `yaml:"tags" json:"tags"`
	Filters     FiltersConfig  `yaml:"filters" json:"filters"`
	Command     CommandsConfig `yaml:"command" json:"command"`
}

func (p Page) Namespace() string {
	return fmt.Sprintf("%s/%s", p.Category, p.CodeName)
}

type CommandsConfig struct {
	Content CommandConfig `yaml:"content" json:"content"`
}

type CommandConfig struct {
	Name string   `yaml:"name" json:"name"`
	Args []string `yaml:"args" json:"args"`
}

type HtmlFilter struct {
	TextOnly bool   `yaml:"textOnly"`
	Tables   string `yaml:"tables"`
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
	baseDir := "config/categories"
	var categories []Category

	for _, catName := range cfg.Categories {
		ok, cat := loadCategoryFile(baseDir, catName)
		if ok {
			log.Info().Str("category", cat.Name).Msg("Loaded category:")
			categories = append(categories, cat)
		}
	}

	return categories
}

func loadCategoryFile(baseDir string, catName string) (bool, Category) {
	fp := filepath.Join(baseDir, catName+".yml")
	log.Info().Str("file", fp).Msg("Loading file")

	content, err := ioutil.ReadFile(fp)
	if err != nil {
		log.Error().Err(err).Str("file", fp).Msg("Unable to open file")
		return false, Category{}
	}

	var cat Category
	if err = yaml.Unmarshal(content, &cat); err != nil {
		log.Error().Err(err).Str("file", fp).Msg("Unable to load file")
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

type RunSelector struct {
	Tags     []string
	Category string
	Page     string
	Force    bool
}
