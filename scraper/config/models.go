package config

import (
	"io/ioutil"
	"log"
	"path/filepath"

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

type CommandsConfig struct {
	Content CommandConfig `yaml:"content" json:"content"`
}

type CommandConfig struct {
	Name string   `yaml:"name" json:"name"`
	Args []string `yaml:"args" json:"args"`
}

type HtmlFilter struct {
	NoPrettyTables 	bool 	`yaml:"noPrettyTables"`
	TextOnly     	bool 	`yaml:"textOnly"`
	Tables 			string 	`yaml:"tables"`
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
			log.Printf("Loaded category: %v", cat.Name)
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

type RunSelector struct {
	Tags     []string
	Category string
	Page     string
	Force    bool
}
