package middlewares

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/pestanko/miniscrape/pkg/applog"
	"go.opentelemetry.io/otel/trace"

	"github.com/rs/zerolog"
)

// LogParams represents a logger params
type LogParams struct {
	LogCfg applog.LogConfig
	Log    zerolog.Logger
}

// Logger log all requests
func Logger(params LogParams) func(targetMux http.Handler) http.Handler {

	return func(targetMux http.Handler) http.Handler {
		accessLog := applog.MakeAccessLog(&params.LogCfg)

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			o := &responseObserver{ResponseWriter: w}

			ctx := r.Context()
			span := trace.SpanContextFromContext(ctx)

			spanDict := zerolog.Dict().
				Str("trace_id", span.TraceID().String()).
				Str("span_id", span.SpanID().String())

			ll := zerolog.Ctx(ctx).With().
				Dict("otel", spanDict).
				Logger()

			accessLog := accessLog.With().Dict("otel", spanDict).Logger()

			ctx = ll.WithContext(ctx)
			r = r.WithContext(ctx)

			targetMux.ServeHTTP(o, r)

			// log request by who(IP address)
			accessLog.Info().
				Interface("req", map[string]any{
					"method":      r.Method,
					"uri":         r.RequestURI,
					"source_addr": r.RemoteAddr,
				}).
				Interface("res", map[string]any{
					"duration": time.Since(start),
					"status":   o.status,
				}).
				Msg("Incoming request")

			slog.Default().InfoContext(ctx, "Incoming request",
				slog.String("htt_method", r.Method),
				slog.String("request_uri", r.RequestURI),
				slog.String("source_addr", r.RemoteAddr),
				slog.Duration("duration", time.Since(start)),
				slog.Int("http_status", o.status),
			)
		})
	}
}

type responseObserver struct {
	http.ResponseWriter
	status      int
	written     int64
	wroteHeader bool
}

// Write using observer
func (o *responseObserver) Write(p []byte) (n int, err error) {
	if !o.wroteHeader {
		o.WriteHeader(http.StatusOK)
	}
	n, err = o.ResponseWriter.Write(p)
	o.written += int64(n)
	return
}

// WriteHeader using observer
func (o *responseObserver) WriteHeader(code int) {
	if o.wroteHeader {
		return
	}
	o.ResponseWriter.WriteHeader(code)
	o.wroteHeader = true
	o.status = code
}
