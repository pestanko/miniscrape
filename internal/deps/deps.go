// Package deps provides the application dependencies
package deps

import (
	"context"

	"github.com/pestanko/miniscrape/internal/config"
	"github.com/pestanko/miniscrape/pkg/instrument"
	"github.com/rs/zerolog/log"
)

// Deps Represents an application dependencies
type Deps struct {
	Cfg *config.AppConfig
}

// Close the dependencies
func (d Deps) Close(ctx context.Context) error {
	if d.Cfg.Otel.Enabled {
		if err := instrument.Shutdown(ctx); err != nil {
			log.Error().Err(err).Msg("Failed to shutdown OpenTelemetry")
		}
	}
	return nil
}

// InitAppDeps init the application dependencies
func InitAppDeps() (*Deps, error) {
	cfg := config.GetAppConfig()

	if cfg.Otel.Enabled {
		if _, err := instrument.SetupOTEL(context.Background(), &cfg.Otel); err != nil {
			log.Error().Err(err).Msg("Failed to setup tracing")
		}
	}

	return &Deps{
		Cfg: cfg,
	}, nil
}
