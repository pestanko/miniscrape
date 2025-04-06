package cmd

import (
	"fmt"

	"github.com/pestanko/miniscrape/internal/config"
	"github.com/pestanko/miniscrape/internal/models"
	"github.com/pestanko/miniscrape/internal/scraper"
	"github.com/pestanko/miniscrape/pkg/applog"

	"github.com/spf13/cobra"
)

var (
	selector    models.RunSelector
	noCache     bool
	noContent   bool
	updateCache bool
)

// scrapeCmd represents the scrape command
var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "Scrape the pages",
	Long:  ``,
	RunE: func(cmd *cobra.Command, _ []string) error {
		cfg := config.GetAppConfig()
		applog.InitGlobalLogger(&cfg.Log)
		if noCache {
			cfg.Cache.Enabled = false
		}
		if updateCache {
			cfg.Cache.Update = true
		}

		scrapeService := scraper.NewService(cfg)
		results := scrapeService.Scrape(cmd.Context(), selector)
		for _, r := range results {
			fmt.Printf("Result[%s] for  \"%s (%s)\" (url: \"%s\")\n",
				r.Status,
				r.Page.Name,
				r.Page.CodeName,
				r.Page.Homepage)
			if !noContent {
				fmt.Printf("%s\n\n", r.Content)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(scrapeCmd)

	// Here you will define your flags and configuration settings.

	// and all subcommands, e.g.:
	// scrapeCmd.PersistentFlags().String("foo", "", "A help for foo")
	scrapeCmd.PersistentFlags().StringVarP(&selector.Category, "category", "C", "",
		"Scrape pages based on the category")
	scrapeCmd.PersistentFlags().StringSliceVarP(&selector.Tags, "tags", "T", []string{},
		"Select pages based on provided tags")

	scrapeCmd.PersistentFlags().BoolVar(&noCache, "no-cache", false,
		"Disable caching")
	scrapeCmd.PersistentFlags().BoolVarP(&updateCache, "update-cache", "U", false,
		"Update cache")

	scrapeCmd.PersistentFlags().StringVarP(&selector.Page, "name", "N", "",
		"Select by codename")

	scrapeCmd.PersistentFlags().BoolVarP(&selector.Force, "force", "f", false,
		"Force scrape - ignore disabled")

	scrapeCmd.PersistentFlags().BoolVar(&noContent, "no-content", false,
		"Do not print out the content")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scrapeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
