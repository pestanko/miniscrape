// Package asyncrun contains tha asynchronous runner,
// it's responsibilities are:
// - run the callback asynchronously, with the Stop/Shutdown functionality
package asyncrun

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

const defaultTimeout = 30 * time.Second

// ErrShutdownTimeout this error will be returned if the graceful timeout period is exceeded
var ErrShutdownTimeout = fmt.Errorf("runtime context timeout: %w", context.DeadlineExceeded)

// Params parameters for the AsyncRun function;
// Run defines the main callback - function that will be executed async way
// Shutdown defines a stop/shutdown callback - function that will do a cleanup/shutdown
// of the Run function
type Params struct {
	Run          func(ctx context.Context) error
	Shutdown     func(ctx context.Context) error
	GraceTimeout time.Duration
}

// AsyncRun will start the async execution of the params.Run callback provided through the
// params argument
// The function returns an error channel - the caller should be checking for the runtime errors
// and wait until there is some error message or the channel is closed
// most general usage:
// ```go
// errC := AsyncRun(ctx, params)
//
//	if err = <-errC; err != nil {
//		   return err
//	}
//
// / ```
func AsyncRun(appCtx context.Context, params Params) chan error {
	if params.GraceTimeout == 0 {
		params.GraceTimeout = defaultTimeout
	}

	if params.Run == nil {
		params.Run = noopFn
	}

	if params.Shutdown == nil {
		params.Shutdown = noopFn
	}

	errC := make(chan error)
	runtimeCtx, runtimeCancelCallback := context.WithCancel(context.Background())

	// run the Run function
	go func() {
		if err := params.Run(runtimeCtx); err != nil {
			errC <- err
		}
	}()

	// run the Stop function
	go func() {
		// if the application context is closed, the function will continue
		// we start with the shutdown process
		<-appCtx.Done()

		// Either wait for shutdown to complete
		// or timeout in params.GraceTimeout seconds (default: 30)
		go func() {
			defer close(errC)

			select {
			// shutdown completed
			case <-runtimeCtx.Done():
				log.Debug().Msg("graceful shutdown completed in time")

			// graceful shutdown timeout
			case <-time.After(params.GraceTimeout):
				defer runtimeCancelCallback()
				err := ErrShutdownTimeout
				log.Error().Err(err).
					Msg("graceful shutdown timed out.. forcing exit.")
				errC <- err
			}
		}()

		if err := params.Shutdown(runtimeCtx); err != nil {
			errC <- err
		}

		defer runtimeCancelCallback()
	}()

	return errC
}

// noopFn represents a no operation function
func noopFn(ctx context.Context) error {
	return nil
}
