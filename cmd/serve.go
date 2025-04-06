package cmd

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"

	"github.com/pestanko/miniscrape/internal/web"
	"github.com/pestanko/miniscrape/pkg/rest/chiapp"
	"github.com/pestanko/miniscrape/pkg/utils/applog"

	"github.com/pestanko/miniscrape/internal/deps"
	"github.com/pestanko/miniscrape/pkg/apprun"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve a simple API for the miniscrape",
	Long:  `Serve a simple API for the miniscrape`,
	RunE: func(cmd *cobra.Command, _ []string) error {
		if err := godotenv.Load(); err != nil {
			log.Debug().Msg("No .env file found, using environment variables only")
		}

		run := apprun.NewAppRunner(
			apprun.WithDepProvider(deps.InitAppDeps),
		)

		return run.Run(cmd.Context(), func(ctx context.Context, d *deps.Deps) error {
			applog.InitGlobalLogger(&d.Cfg.Log)

			server := web.NewServer(d.Cfg)

			listenAddr := d.Cfg.Web.Addr
			if listenAddr == "" {
				listenAddr = ":8080"
			}

			ops := chiapp.RunOps{
				ListenAddr:       listenAddr,
				ReadTimeout:      5 * time.Second,
				GraceFullTimeout: 30 * time.Second,
			}

			errC, err := chiapp.RunWebServer(ctx, server, ops)
			if err != nil {
				log.Error().Err(err).Msg("Failed to run web server - received error from the server")
				return err
			}

			if err = <-errC; errC != nil {
				if errors.Is(err, http.ErrServerClosed) {
					log.Info().Msg("Web server closed")
					return nil
				}

				log.Error().Err(err).Msg("Failed to run web server - received error from the channel")
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
