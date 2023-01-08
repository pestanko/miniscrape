package apprun

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
)

const defaultTimeout = 30 * time.Second

// StartParams parameters for the Start function;
// Start defines the start callback - function that will be executed async way
// Stop defines a stop/shutdown callback - function that will executed to cleanup/stop the Start
// function
type StartParams struct {
	Start        func(ctx context.Context) error
	Stop         func(ctx context.Context) error
	GraceTimeout time.Duration
}

// Start executing the StartParams.Start function asynchronously and then clean up/shutdown the
// execution
func Start(appCtx context.Context, params StartParams) (chan error, error) {
	if params.GraceTimeout == 0 {
		params.GraceTimeout = defaultTimeout
	}

	if params.Start == nil {
		params.Start = noopFn
	}

	if params.Stop == nil {
		params.Stop = noopFn
	}

	errC := make(chan error)
	runtimeConnect := context.Background()

	// run the Start function
	go func() {
		if err := params.Start(runtimeConnect); err != nil {
			errC <- err
		}
	}()

	// run the Stop function
	go func() {
		// if the application context is closed, the function will continue
		<-appCtx.Done()
		// Stop signal with grace period of 30 seconds
		shutdownCtx, cancelCallback := context.WithTimeout(runtimeConnect, params.GraceTimeout)
		defer cancelCallback()

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Err(shutdownCtx.Err()).Msg("graceful shutdown timed out.. forcing exit.")
				errC <- shutdownCtx.Err()
			}
		}()

		if err := params.Stop(runtimeConnect); err != nil {
			errC <- err
		}

		close(errC)

	}()

	return errC, nil
}

func noopFn(ctx context.Context) error {
	return nil
}
