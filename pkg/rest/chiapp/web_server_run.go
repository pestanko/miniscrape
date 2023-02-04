package chiapp

import (
	"context"
	"github.com/pestanko/miniscrape/pkg/asyncrun"
	"net/http"
	"time"
)

// RunOps defines runtime options
type RunOps struct {
	ListenAddr       string
	ReadTimeout      time.Duration
	GraceFullTimeout time.Duration
}

func RunWebServer(appCtx context.Context, handler http.Handler, ops RunOps) chan error {
	server := http.Server{
		Addr:        ops.ListenAddr,
		Handler:     handler,
		ReadTimeout: ops.ReadTimeout,
	}

	params := asyncrun.Params{
		Run: func(ctx context.Context) error {
			return server.ListenAndServe()
		},
		Shutdown: func(ctx context.Context) error {
			return server.Shutdown(ctx)
		},
		GraceTimeout: ops.GraceFullTimeout,
	}

	return asyncrun.AsyncRun(appCtx, params)
}
