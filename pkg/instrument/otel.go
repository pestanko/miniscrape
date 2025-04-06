// Package instrument provides OpenTelemetry instrumentation for the application.
package instrument

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/pestanko/miniscrape/pkg/utils"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	otellog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
)

// TracerProvider Global tracer provider
var TracerProvider *sdktrace.TracerProvider

// LoggerProvider Global logger provider
var LoggerProvider *otellog.LoggerProvider

// MeterProvider Global meter provider
var MeterProvider *metric.MeterProvider

// SetupOTEL initializes OTEL
func SetupOTEL(ctx context.Context, otelCfg *OtelConfig) (func(), error) {
	ll := log.With().Interface("otel", otelCfg).Logger()
	// If OTEL is not enabled, return a no-op function
	if !otelCfg.Enabled {
		log.Info().Msg("OpenTelemetry tracing disabled")
		return func() {}, nil
	}

	ll.Debug().Msg("OpenTelemetry tracing enabled")

	serviceInfo := getServiceInfo(ctx, otelCfg)

	// Create a resource with the service name
	res, err := getResourceFromServiceInfo(ctx, serviceInfo)

	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Setup the tracing
	if err := SetupOTELTracing(ctx, otelCfg, res); err != nil {
		return nil, fmt.Errorf("failed to setup OpenTelemetry tracing: %w", err)
	}

	// Setup the logging
	if err := SetupOTELLogging(ctx, otelCfg, serviceInfo); err != nil {
		return nil, fmt.Errorf("failed to setup OpenTelemetry logging: %w", err)
	}

	// Setup the metrics
	if err := SetupOTELMetrics(ctx, otelCfg, res); err != nil {
		return nil, fmt.Errorf("failed to setup OpenTelemetry metrics: %w", err)
	}

	ll.Info().
		Interface("service_info", serviceInfo).
		Msg("OpenTelemetry tracing initialized")

	// Return a cleanup function
	return func() {
		if err := TracerProvider.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Msg("Failed to shutdown tracer provider")
		}
	}, nil
}

// Shutdown shuts down the OpenTelemetry tracing
func Shutdown(ctx context.Context) error {
	var err error
	if TracerProvider != nil {
		if err = TracerProvider.Shutdown(ctx); err != nil {
			err = errors.Join(err, fmt.Errorf("failed to shutdown tracer provider: %w", err))
		}
	}

	if LoggerProvider != nil {
		if err = LoggerProvider.Shutdown(ctx); err != nil {
			err = errors.Join(err, fmt.Errorf("failed to shutdown logger provider: %w", err))
		}
	}

	if MeterProvider != nil {
		if err = MeterProvider.Shutdown(ctx); err != nil {
			err = errors.Join(err, fmt.Errorf("failed to shutdown meter provider: %w", err))
		}
	}

	return err
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

// SetupOTELTracing initializes OTEL tracing
func SetupOTELTracing(ctx context.Context, otelCfg *OtelConfig, res *resource.Resource) error {
	ll := log.With().Interface("otel", otelCfg).Logger()
	var err error
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
		return fmt.Errorf("failed to create OTLP exporter: %w", err)
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
	return nil
}

// SetupOTELLogging initializes OpenTelemetry logging
func SetupOTELLogging(ctx context.Context, otelCfg *OtelConfig, info *ServiceInfo) error {
	logExporter, err := otlploggrpc.New(
		ctx,
		otlploggrpc.WithEndpoint(otelCfg.Endpoint),
		otlploggrpc.WithCompressor("gzip"),
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
	logger := otelslog.NewLogger(
		info.Name,
		otelslog.WithSource(true),
	)

	slog.SetDefault(logger)

	logger.InfoContext(ctx, "OpenTelemetry logging initialized")

	return nil
}

// SetupOTELMetrics initializes OpenTelemetry metrics
func SetupOTELMetrics(ctx context.Context, otelCfg *OtelConfig, res *resource.Resource) error {
	// Interval which the metrics will be reported to the collector
	interval := 5 * time.Second
	metricExporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint(otelCfg.Endpoint),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		return fmt.Errorf("failed to create metric exporter: %w", err)
	}

	periodicReader := metric.NewPeriodicReader(
		metricExporter,
		metric.WithInterval(interval),
	)

	MeterProvider = metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(periodicReader),
	)

	otel.SetMeterProvider(MeterProvider)

	slog.Default().InfoContext(ctx, "OpenTelemetry metrics initialized")

	return nil
}

// OtelConfig holds the OpenTelemetry configuration
type OtelConfig struct {
	Enabled  bool   `env:"OTEL_ENABLED,default=true" json:"enabled" yaml:"enabled"`
	Endpoint string `env:"OTEL_EXPORTER_OTLP_ENDPOINT,default=localhost:4317" json:"endpoint" yaml:"endpoint"`
	Protocol string `env:"OTEL_EXPORTER_OTLP_PROTOCOL,default=grpc" json:"protocol" yaml:"protocol"`
	Insecure bool   `env:"OTEL_INSECURE,default=true" json:"insecure" yaml:"insecure"`
}

// ServiceInfo holds the service information
type ServiceInfo struct {
	Name    string `json:"name" yaml:"name"`
	Version string `json:"version" yaml:"version"`
	Env     string `json:"env" yaml:"env"`
}

// getServiceInfo returns the service information
func getServiceInfo(_ context.Context, _ *OtelConfig) *ServiceInfo {
	serviceName := utils.GetEnvOrDefault("SERVICE_NAME", "service")
	serviceVersion := utils.GetEnvOrDefault("SERVICE_VERSION", "v1.0.0")
	envName := utils.GetEnvOrDefault("ENV_NAME", "dev")

	return &ServiceInfo{
		Name:    serviceName,
		Version: serviceVersion,
		Env:     envName,
	}
}

// getResourceFromServiceInfo creates a resource from the service info
func getResourceFromServiceInfo(ctx context.Context, info *ServiceInfo) (*resource.Resource, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(info.Name),
			semconv.ServiceVersionKey.String(info.Version),
			semconv.DeploymentEnvironmentKey.String(info.Env),
		),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	return res, nil
}
