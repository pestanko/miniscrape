// Package apprun represents a simple application runner with application dependencies
package apprun

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"github.com/pestanko/miniscrape/internal/config"
	"github.com/pestanko/miniscrape/internal/instrumentation"
	"github.com/pestanko/miniscrape/pkg/utils/collut"
	"github.com/rs/zerolog/log"
)

// CtxCloser is a type that implements the Close method with a context parameter
type CtxCloser interface {
	Close(ctx context.Context) error
}

// NewAppRunner create a new application runner instance
// Datadog Tracing can be enabled/disabled in 2 ways, you can either:
// - set TRACING_ENABLED env variable to true/false
// - directly call WithForceTracingEnabled(true/false)
// function WithForceTracingEnabled - has precedence and overrides env variable
// it can be useful for testing
func NewAppRunner[D CtxCloser](ops ...func(a *AppRunner[D])) AppRunner[D] {
	a := AppRunner[D]{}
	a.isTracingEnabled, _ = strconv.ParseBool(os.Getenv("TRACING_ENABLED"))
	collut.OpsApplyAllRef(&a, ops...)

	return a
}

// WithDepProvider set the dependency provider func for the App Runner
// Dependency provided, provides dependencies for the application
// dependencies are for example database connection pools, redis, mongo clients ...
// The dependencies are then injected to the function body provided to the Run method
func WithDepProvider[D CtxCloser](dp func() (D, error)) func(a *AppRunner[D]) {
	return func(a *AppRunner[D]) {
		a.DependencyProvider = func(_ context.Context) (D, error) {
			return dp() // nolint:wrapcheck
		}
	}
}

// WithDepProviderCtx set the dependency provider func for the App Runner with context
// for more info see WithDepProvider
func WithDepProviderCtx[D CtxCloser](
	dp func(ctx context.Context) (D, error),
) func(a *AppRunner[D]) {
	return func(a *AppRunner[D]) {
		a.DependencyProvider = dp
	}
}

// WithForceTracingEnabled force tracing enabled/disabled
// It overrides default behavior, where whether tracing is enabled
// is based on TRACING_ENABLED env var
func WithForceTracingEnabled[D CtxCloser](isEnabled bool) func(a *AppRunner[D]) {
	return func(a *AppRunner[D]) {
		a.isTracingEnabled = isEnabled
	}
}

// AppRunner represents a runner for the application
// generic parameter D (as dependencies) must implement closer,
// the Run method closes the "dependencies" after the execution
type AppRunner[D CtxCloser] struct {
	DependencyProvider func(ctx context.Context) (D, error)
	isTracingEnabled   bool
}

// Run the application
// The method will:
// - enable/disable datadog tracing
// - creates dependencies using the DependencyProvider callback
// - handles the application shutdown - (app receives the signal)
// - provides the initialized dependencies and context to the function body callback
// - do cleanup - closes dependencies, stops the application
func (a *AppRunner[D]) Run(
	ctx context.Context,
	body func(ctx context.Context, d D) error,
) error {

	// we need to create dependency provider
	// after that we can close them
	deps, err := a.DependencyProvider(ctx)

	if err != nil {
		log.Error().Err(err).
			Msg("unable to initialize the application")
		return fmt.Errorf("unable to initialize the application: %w", err)
	}

	ctx, stopCallback := signal.NotifyContext(context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	defer func() {
		stopCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := deps.Close(stopCtx); err != nil {
			log.Error().Err(err).
				Msg("unable to de-initialize (close) the application")
		}
		log.Info().Msg("application dependencies has been cleaned-up")

		// stop the execution
		stopCallback()
	}()

	// run the inner application
	return body(ctx, deps)
}

func initTracing(ctx context.Context, cfg *config.AppConfig) {
	if cfg.Otel.Enabled {
		if _, err := instrumentation.SetupTracing(ctx, cfg); err != nil {
			log.Error().Err(err).Msg("Failed to setup tracing")
		}
	}

	propagator := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(propagator)
}
