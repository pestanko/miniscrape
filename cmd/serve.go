package cmd

import (
	"github.com/pestanko/miniscrape/pkg/config"
	"github.com/pestanko/miniscrape/pkg/utils"
	"github.com/pestanko/miniscrape/pkg/web"

	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve a simple API for the pkg",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.GetAppConfig()
		utils.InitGlobalLogger(&cfg.Log)

		server := web.NewServer(cfg)

		server.Serve()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
