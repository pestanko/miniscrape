package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

// Category representation
type Category struct {
	// Pages list of all pages
	Pages []Page `yaml:"pages" json:"pages"`
	// Name of the category
	Name string `yaml:"name" json:"name"`
}

// Page representation
type Page struct {
	// CodeName of the page
	CodeName string `yaml:"codename" json:"codename"`
	// Name of the page
	Name string `yaml:"name" json:"name"`
	// Homepage of the page
	Homepage string `yaml:"homepage" json:"homepage"`
	// URL of the page where the lunch menu is
	URL string `yaml:"url" json:"url"`
	// Query css query to use for element extraction
	Query string `yaml:"query" json:"query"`
	// XPath query to use for element extraction
	XPath string `yaml:"xpath" json:"xpath"`
	// CachePolicy for the webpage
	CachePolicy string `yaml:"cache_policy" json:"cachePolicy"`
	// Resolver to be used
	Resolver string `yaml:"resolver" json:"resolver"`
	// Category name
	Category string `yaml:"category" json:"category"`
	// Disabled whether the page has been disabled
	Disabled bool `yaml:"disabled" json:"disabled"`
	// Tags list of all tags for the page
	Tags []string `yaml:"tags" json:"tags"`
	// Filters for the page
	Filters FiltersConfig `yaml:"filters" json:"filters"`
	// Command config for cmd to be executed to get webpage content
	Command CommandsConfig `yaml:"command" json:"command"`
}

// Namespace for the page
func (p Page) Namespace() string {
	return fmt.Sprintf("%s/%s", p.Category, p.CodeName)
}

// CommandsConfig wrapper for command configuration for the page
type CommandsConfig struct {
	// Content command configuration content
	Content CommandConfig `yaml:"content" json:"content"`
}

// CommandConfig command configuration for the page
type CommandConfig struct {
	// Name of the command
	Name string `yaml:"name" json:"name"`
	// Args of the command
	Args []string `yaml:"args" json:"args"`
}

// FiltersConfig for the webpage
type FiltersConfig struct {
	// Cut filter configuration
	Cut CutFilter `yaml:"cut"`
	// CutLine filter configuration
	CutLine CutLineFilter `yaml:"cutLine"`
	// Day filter configuration
	Day DayFilter `yaml:"day"`
	// HTML filter configuration
	HTML HTMLFilter `yaml:"html"`
}

// HTMLFilter for the webpage
type HTMLFilter struct {
	// TextOnly - whether it should parse only the text
	TextOnly bool `yaml:"textOnly"`
	// Tables resolver
	Tables string `yaml:"tables"`
}

// CutFilter for the webpage
type CutFilter struct {
	// Before which text the content should be cut
	Before string `yaml:"before"`
	// After which text the content should be cut
	After string `yaml:"after"`
}

// CutLineFilter for the webpage
type CutLineFilter struct {
	// StartsWith remove the line if it starts with the text
	StartsWith string `yaml:"startsWith"`
	// Contains remove the line if it contains the text
	Contains string `yaml:"contains"`
	// CutAfter Cut the line after provided text
	CutAfter string `yaml:"cutAfter"`
}

// DayFilter for the webpage
type DayFilter struct {
	// List of days to be used as separators, if empty - use default
	Days []string `yaml:"days"`
	// Whether the filter is enabled
	Enabled bool `yaml:"enabled"`
}

// LoadCategories Load all categories from the app config
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

	content, err := os.ReadFile(filepath.Clean(fp))
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

// RunSelector represents which pages should be selected
type RunSelector struct {
	// Tags list of all tags to be selected
	Tags []string
	// Category name of the category
	Category string
	// Page codename for the page
	Page string
	// Force load even if disabled
	Force bool
}
