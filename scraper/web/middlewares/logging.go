package middlewares

import (
	"fmt"
	"net/http"
	"time"

	"github.com/pestanko/miniscrape/scraper/config"
	"github.com/pestanko/miniscrape/scraper/utils"
)

func RequestLogger(targetMux http.Handler, cfg *config.AppConfig) http.Handler {

	accessLog := utils.MakeAccessLog(&cfg.Log)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		o := &responseObserver{ResponseWriter: w}

		targetMux.ServeHTTP(o, r)

		// log request by who(IP address)
		accessLog.Info().
			Str("method", r.Method).
			Str("requestUri", r.RequestURI).
			Str("remoteAddr", r.RemoteAddr).
			Str("duration", fmt.Sprintf("%v", time.Since(start))).
			Int("statusCode", o.status).
			Msg("Incoming request")
	})
}

type responseObserver struct {
	http.ResponseWriter
	status      int
	written     int64
	wroteHeader bool
}

func (o *responseObserver) Write(p []byte) (n int, err error) {
	if !o.wroteHeader {
		o.WriteHeader(http.StatusOK)
	}
	n, err = o.ResponseWriter.Write(p)
	o.written += int64(n)
	return
}

func (o *responseObserver) WriteHeader(code int) {
	if o.wroteHeader {
		return
	}
	o.ResponseWriter.WriteHeader(code)
	o.wroteHeader = true
	o.status = code
}
