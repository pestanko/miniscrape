// Package instrumentation provides OpenTelemetry instrumentation for the application.
package instrumentation

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/pestanko/miniscrape/internal/config"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	otellog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

// TracerProvider Global tracer provider
var TracerProvider *sdktrace.TracerProvider

// LoggerProvider Global logger provider
var LoggerProvider *otellog.LoggerProvider

// SetupTracing initializes OpenTelemetry tracing
func SetupTracing(ctx context.Context, cfg *config.AppConfig) (func(), error) {
	otelCfg := cfg.Otel
	ll := log.With().Interface("otel", otelCfg).Logger()
	// If OTEL is not enabled, return a no-op function
	if !otelCfg.Enabled {
		log.Info().Msg("OpenTelemetry tracing disabled")
		return func() {}, nil
	}

	ll.Debug().Msg("OpenTelemetry tracing enabled")

	if os.Getenv("ENV_NAME") == "" {
		ll.Warn().Msg("ENV_NAME is not set, using 'prod' as default")
		os.Setenv("ENV_NAME", "prod")
	}

	// Create a resource with the service name
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("miniscrape"),
			semconv.ServiceVersionKey.String("v1.0.0"),
			semconv.DeploymentEnvironmentKey.String(os.Getenv("ENV_NAME")),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create OTLP exporter
	var exporter *otlptrace.Exporter
	if otelCfg.Protocol == "http" {
		// HTTP exporter
		opts := []otlptracehttp.Option{
			otlptracehttp.WithEndpoint(otelCfg.Endpoint),
		}
		if otelCfg.Insecure {
			opts = append(opts, otlptracehttp.WithInsecure())
		}
		exporter, err = otlptracehttp.New(ctx, opts...)
		ll.Debug().Msg("OTLP HTTP exporter created")
	} else {
		// Default to gRPC exporter
		opts := []otlptracegrpc.Option{
			otlptracegrpc.WithEndpoint(otelCfg.Endpoint),
		}
		if otelCfg.Insecure {
			opts = append(opts, otlptracegrpc.WithInsecure())
		}
		exporter, err = otlptracegrpc.New(ctx, opts...)
		ll.Debug().Msg("OTLP gRPC exporter created")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	// Create TracerProvider
	TracerProvider = sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
	)

	// Set global TracerProvider
	otel.SetTracerProvider(TracerProvider)

	// Set global propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Setup the logging
	if err := setupOTELLogging(ctx, &otelCfg); err != nil {
		return nil, fmt.Errorf("failed to setup OpenTelemetry logging: %w", err)
	}

	ll.Info().
		Msg("OpenTelemetry tracing initialized")

	// Return a cleanup function
	return func() {
		if err := TracerProvider.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Msg("Failed to shutdown tracer provider")
		}
	}, nil
}

// Tracer returns a tracer from the global provider
func Tracer(name string) trace.Tracer {
	return otel.Tracer(name)
}

// WrapHandler wraps an HTTP handler with OpenTelemetry instrumentation
func WrapHandler(handler http.Handler, operation string) http.Handler {
	return otelhttp.NewHandler(handler, operation)
}

// TracingMiddleware returns middleware for tracing HTTP requests
func TracingMiddleware(next http.Handler, serviceName string) http.Handler {
	return otelhttp.NewHandler(next, serviceName)
}

// setupOTELLogging initializes OpenTelemetry logging
func setupOTELLogging(ctx context.Context, otelCfg *config.OtelConfig) error {
	logExporter, err := otlploggrpc.New(
		ctx,
		otlploggrpc.WithEndpoint(otelCfg.Endpoint),
		otlploggrpc.WithInsecure(),
	)

	if err != nil {
		return fmt.Errorf("failed to create log exporter: %w", err)
	}

	// Create the logger provider
	LoggerProvider := otellog.NewLoggerProvider(
		otellog.WithProcessor(
			otellog.NewBatchProcessor(logExporter),
		),
	)

	// Set the logger provider globally
	global.SetLoggerProvider(LoggerProvider)

	// Instantiate a new slog logger
	logger := otelslog.NewLogger("miniscrape")

	logger.InfoContext(ctx, "OpenTelemetry logging initialized")

	return nil
}
