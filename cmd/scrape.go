/*
Copyright © 2022 Peter Stanko <peter.stanko0@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
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
)

// scrapeCmd represents the scrape command
var scrapeCmd = &cobra.Command{
	Use:   "scrape",
	Short: "Scrape the pages",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetAppConfig()
		categories := config.LoadCategories(cfg)
		runner := scraper.NewAsyncRunner(cfg, categories)
		selector := scraper.RunSelector{
			Tags:     tagsArg,
			Category: categoryArg,
		}
		results := runner.Run(selector)
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

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// scrapeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
