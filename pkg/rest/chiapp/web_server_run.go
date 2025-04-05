package chiapp

import (
	"context"
	"net/http"
	"time"

	"github.com/pestanko/miniscrape/pkg/apprun"
	"github.com/rs/zerolog"
)

// RunOps defines runtime options
type RunOps struct {
	ListenAddr       string
	ReadTimeout      time.Duration
	GraceFullTimeout time.Duration
}

// RunWebServer runs a web server
func RunWebServer(appCtx context.Context, handler http.Handler, ops RunOps) (chan error, error) {
	server := http.Server{
		Addr:        ops.ListenAddr,
		Handler:     handler,
		ReadTimeout: ops.ReadTimeout,
	}

	params := apprun.StartParams{
		Start: func(ctx context.Context) error {
			zerolog.Ctx(ctx).Info().
				Str("listen_addr", ops.ListenAddr).
				Msg("Starting web server")
			return server.ListenAndServe()
		},
		Stop: func(ctx context.Context) error {
			ll := zerolog.Ctx(ctx).With().Str("listen_addr", ops.ListenAddr).Logger()
			ll.Info().Msg("Stopping web server")
			if err := server.Shutdown(ctx); err != nil {
				ll.Error().Err(err).Msg("Failed to stop web server")
			} else {
				ll.Info().Msg("Web server stopped")
			}
			return nil
		},
		GraceTimeout: ops.GraceFullTimeout,
	}
	return apprun.Start(appCtx, params)
}
