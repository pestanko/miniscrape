package middlewares

import (
	"log"
	"net/http"
	"time"
)

func RequestLogger(targetMux http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		o := &responseObserver{ResponseWriter: w}

		targetMux.ServeHTTP(o, r)

		// log request by who(IP address)
		log.Printf(
			"%-7s%s\t%-6d%s\t %v",
			r.Method,
			r.RequestURI,
			o.status,
			r.RemoteAddr,
			time.Since(start),
		)
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
