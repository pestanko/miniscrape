package apprun

import (
	"context"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

// initTracerProvider initialize the Otel Tracer provider
func initTracerProvider(ctx context.Context) *sdktrace.TracerProvider {
	//exporter, err := stdout.New(stdout.WithPrettyPrint())
	//if err != nil {
	//	log.Err(err).Msg("unable to init trace provider")
	//	return nil
	//}
	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			// the service name used to display traces in backends
			semconv.ServiceNameKey.String("mux-server"),
		),
	)
	if err != nil {
		log.Err(err).Msg("unable to initialize resource")
		return nil
	}

	return sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		//sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
	)
}
