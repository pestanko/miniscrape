/*
	Copyright © 2022 Peter Stanko <peter.stanko0@gmail.com>
*/
package cmd

import (
	"fmt"
	"github.com/pestanko/miniscrape/scraper"
	"github.com/pestanko/miniscrape/scraper/config"
	"github.com/spf13/cobra"
)

var (
	categoryArg string
	tagsArg     []string
	noCache     bool
	updateCache bool
	nameArg     string
)

// scrapeCmd represents the scrape command
var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "Scrape the pages",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetAppConfig()
		if noCache {
			cfg.Cache.Enabled = false
		}
		if updateCache {
			cfg.Cache.Update = true
		}
		selector := config.RunSelector{
			Tags:     tagsArg,
			Category: categoryArg,
			Page:     nameArg,
		}

		scrapeService := scraper.NewService(cfg)
		results := scrapeService.Scrape(selector)
		for _, r := range results {
			fmt.Printf("Result[%s] for  \"%s (%s)\" (url: \"%s\")\n",
				r.Status,
				r.Page.Name,
				r.Page.CodeName,
				r.Page.Homepage)
			fmt.Printf("%s\n\n", r.Content)
		}
	},
}

func init() {
	rootCmd.AddCommand(scrapeCmd)

	// Here you will define your flags and configuration settings.

	// and all subcommands, e.g.:
	// scrapeCmd.PersistentFlags().String("foo", "", "A help for foo")
	scrapeCmd.PersistentFlags().StringVarP(&categoryArg, "category", "C", "",
		"Scrape pages based on the category")
	scrapeCmd.PersistentFlags().StringSliceVarP(&tagsArg, "tags", "T", []string{},
		"Select pages based on provided tags")

	scrapeCmd.PersistentFlags().BoolVar(&noCache, "no-cache", false,
		"Disable caching")
	scrapeCmd.PersistentFlags().BoolVarP(&updateCache, "update-cache", "U", false,
		"Update cache")

	scrapeCmd.PersistentFlags().StringVarP(&nameArg, "nameArg", "N", "",
		"Select by codename")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scrapeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
