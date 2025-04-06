// Package deps provides the application dependencies
package deps

import (
	"context"

	"github.com/pestanko/miniscrape/internal/config"
	"github.com/pestanko/miniscrape/internal/instrumentation"
	"github.com/rs/zerolog/log"
)

// Deps Represents an application dependencies
type Deps struct {
	Cfg *config.AppConfig
}

// Close the dependencies
func (d Deps) Close(ctx context.Context) error {
	if d.Cfg.Otel.Enabled {
		if err := instrumentation.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("Failed to shutdown OpenTelemetry")
		}
	}
	return nil
}

// InitAppDeps init the application dependencies
func InitAppDeps() (*Deps, error) {
	cfg := config.GetAppConfig()

	if cfg.Otel.Enabled {
		if _, err := instrumentation.SetupTracing(context.Background(), cfg); err != nil {
			log.Error().Err(err).Msg("Failed to setup tracing")
		}
	}

	return &Deps{
		Cfg: cfg,
	}, nil
}
