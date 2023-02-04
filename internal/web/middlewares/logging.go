package middlewares

import (
	"fmt"
	"github.com/pestanko/miniscrape/pkg/utils/applog"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"time"

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

			requestDict := zerolog.Dict().
				Str("method", r.Method).
				Str("requestUri", r.RequestURI).
				Str("remoteAddr", r.RemoteAddr)

			responseDict := zerolog.Dict().
				Str("duration", fmt.Sprintf("%v", time.Since(start))).
				Int("statusCode", o.status)

			// log request by who(IP address)
			accessLog.Info().
				Dict("request", requestDict).
				Dict("response", responseDict).
				Msg("Incoming request")
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
