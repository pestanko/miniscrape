package cmd

import (
	"context"
	"github.com/pestanko/miniscrape/internal/web"
	"github.com/pestanko/miniscrape/pkg/rest/chiapp"
	"github.com/pestanko/miniscrape/pkg/utils/applog"
	"time"

	"github.com/pestanko/miniscrape/internal/deps"
	"github.com/pestanko/miniscrape/pkg/apprun"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve a simple API for the pkg",
	Long:  `Serve a simple API for the pkg`,
	RunE: func(cmd *cobra.Command, args []string) error {
		run := apprun.NewAppRunner(
			apprun.WithDepProvider(deps.InitAppDeps),
		)

		return run.Run(cmd.Context(), func(ctx context.Context, d *deps.Deps) error {
			applog.InitGlobalLogger(&d.Cfg.Log)

			server := web.NewServer(d.Cfg)

			ops := chiapp.RunOps{
				ListenAddr:       ":8080",
				ReadTimeout:      5 * time.Second,
				GraceFullTimeout: 30 * time.Second,
			}

			errC, err := chiapp.RunWebServer(ctx, server, ops)
			if err != nil {
				return err
			}

			if err = <-errC; errC != nil {
				return err
			}

			return nil
		})

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
