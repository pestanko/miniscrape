package chiapp

import (
	"context"
	"net/http"
	"time"

	"github.com/pestanko/miniscrape/pkg/apprun"
)

// RunOps defines runtime options
type RunOps struct {
	ListenAddr       string
	ReadTimeout      time.Duration
	GraceFullTimeout time.Duration
}

func RunWebServer(appCtx context.Context, handler http.Handler, ops RunOps) (chan error, error) {
	server := http.Server{
		Addr:        ops.ListenAddr,
		Handler:     handler,
		ReadTimeout: ops.ReadTimeout,
	}

	params := apprun.StartParams{
		Start: func(ctx context.Context) error {
			return server.ListenAndServe()
		},
		Stop: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
		GraceTimeout: ops.GraceFullTimeout,
	}
	return apprun.Start(appCtx, params)
}
