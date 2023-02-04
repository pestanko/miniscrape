package asyncrun

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

const testDeadline = 5 * time.Second
const gracefulTimeout = 1 * time.Second

func TestAsyncRun(t *testing.T) {
	t.Run("all params are empty or 0, the context is closed before calling", func(t *testing.T) {
		// all params are empty
		params := Params{
			Run:          nil,
			Shutdown:     nil,
			GraceTimeout: 0,
		}

		appCtx, cancelCallback := context.WithCancel(context.Background())
		// we close the context before calling
		cancelCallback()

		// we start the execution
		errC := AsyncRun(appCtx, params)

		// assert there was no execution error
		// the start and shutdown will do nothing
		assertErrChan(t, errC, nil)
	})

	t.Run(
		"Run the params.Run and shutdown functions with direct ctx cancel",
		func(t *testing.T) {
			// It Tests: it uses the direct context cancellation the assertion
			appCtx, cancelCallback := context.WithCancel(context.Background())
			defer cancelCallback()
			// this channel is used whether the params.Run has been called
			unlockRunC := make(chan bool)
			// this channel is used whether the params.Shutdown has been called
			unlockShutdownC := make(chan bool)
			params := Params{
				// run func is an example function that will do unlockRunC<-true when called,
				// then sleep
				Run: runFunc(t, unlockRunC),
				// run func is an example function that will do unlockShutdownC<-true when called,
				// since it should not timeout (3-rd parameter) there is no sleep,
				// it would end immediately
				Shutdown:     shutdownFunc(t, unlockShutdownC, false),
				GraceTimeout: gracefulTimeout,
			}

			// assert that unlock channel contains value a.k.a the Run has been called
			assertRunUnlock(t, unlockRunC, cancelCallback)

			// assert that unlock channel contains value a.k.a the Shutdown has been called
			assertShutdownUnlock(t, unlockShutdownC)

			// Execute the function
			errC := AsyncRun(appCtx, params)

			// assert there was no execution error
			assertErrChan(t, errC, nil)
		},
	)

	t.Run("Run the params.Run function context timeout", func(t *testing.T) {
		// It Tests: it uses context.WithTimeout instead of direct context cancellation

		appCtx, cancelCallback := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancelCallback()

		unlockRunC := make(chan bool)
		unlockShutdownC := make(chan bool)
		params := Params{
			Run:          runFunc(t, unlockRunC),
			Shutdown:     shutdownFunc(t, unlockShutdownC, false),
			GraceTimeout: gracefulTimeout,
		}

		// do not call the cancel callback, it should timeout
		assertRunUnlock(t, unlockRunC, nil)

		assertShutdownUnlock(t, unlockShutdownC)

		errC := AsyncRun(appCtx, params)

		// assert there was no execution error
		assertErrChan(t, errC, nil)
	})

	t.Run("Shutdown function will exceed the graceful timeout period", func(t *testing.T) {
		// It Tests: it tests the graceful timeout for the shutdown function
		// the Shutdown function will wait (using sleep) until the Graceful period ends

		appCtx, cancelCallback := context.WithTimeout(context.Background(), 300*time.Millisecond)
		defer cancelCallback()

		unlockRunC := make(chan bool)
		unlockShutdownC := make(chan bool)
		params := Params{
			Run:          runFunc(t, unlockRunC),
			Shutdown:     shutdownFunc(t, unlockShutdownC, true),
			GraceTimeout: gracefulTimeout,
		}

		// do not call the cancel callback, it should timeout
		assertRunUnlock(t, unlockRunC, nil)

		assertShutdownUnlock(t, unlockShutdownC)

		errC := AsyncRun(appCtx, params)

		// assert there was execution error: shutdown timeout
		assertErrChan(t, errC, func(err error) {
			assert.ErrorIs(t, err, ErrShutdownTimeout)
		})
	})

	t.Run("The Run function returns error", func(t *testing.T) {
		// It tests: The Run function returns an error
		// this error is sent to `errC<-` channel and we are asserting it

		appCtx, cancelCallback := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancelCallback()

		// only the shutdown is called
		unlockShutdownC := make(chan bool)
		// error that will be returned after the Run starts
		startErr := fmt.Errorf("some run error")
		params := Params{
			// the Run will return an error immediately after call
			Run: func(ctx context.Context) error {
				return startErr
			},
			Shutdown:     shutdownFunc(t, unlockShutdownC, true),
			GraceTimeout: gracefulTimeout,
		}

		assertShutdownUnlock(t, unlockShutdownC)

		errC := AsyncRun(appCtx, params)

		// assert there was execution error: startErr
		assertErrChan(t, errC, func(err error) {
			assert.ErrorIs(t, err, startErr)
		})
	})

	t.Run("The Shutdown function will produce an error", func(t *testing.T) {
		// It tests: The Shutdown function returns an error
		// this error is sent to `errC<-` channel and we are asserting it
		appCtx, cancelCallback := context.WithTimeout(context.Background(), 500*time.Millisecond)
		unlockRunC := make(chan bool)
		shutdownErr := fmt.Errorf("unable to shutdown")
		params := Params{
			Run: runFunc(t, unlockRunC),
			Shutdown: func(ctx context.Context) error {
				return shutdownErr
			},
			GraceTimeout: gracefulTimeout,
		}

		assertRunUnlock(t, unlockRunC, cancelCallback)

		errC := AsyncRun(appCtx, params)

		// assert there was execution error: shutdownErr
		assertErrChan(t, errC, func(err error) {
			assert.ErrorIs(t, err, shutdownErr)
		})
	})
}

// assertErrChan assert whether the errC has been set and to what it has been set
// otherwise the test would timeout
// match - if nil we assert that there is no error, otherwise it takes a match function
func assertErrChan(t *testing.T, errC chan error, match func(err error)) {
	if match == nil {
		match = func(err error) {
			assert.NoError(t, err)
		}
	}

	select {
	case err := <-errC:
		match(err)
	// If there is no error or error channel is not closed, test will timeout
	case <-time.After(testDeadline):
		assert.Fail(t, "test timeout")
	}
}

// assertRunUnlock check whether the Run function channel has been unlocked
// it has to be unlocked - it means that function has been called
// after it was called, we can call the `cancelCallback` that would manually stop the execution
func assertRunUnlock(t *testing.T, unlockRunC chan bool, cancelCallback context.CancelFunc) {
	go func() {
		select {
		// wait until the Run channel unlocks
		case <-unlockRunC:
			assert.True(t, true, "run channel should unblock")
			if cancelCallback != nil {
				cancelCallback()
			}
		// this case will happen only if the Run function is not called
		// Run function should be called always
		case <-time.After(testDeadline):
			assert.Fail(t, "test timeout for run callback - THIS SHOULD NEVER HAPPEN")
		}
	}()
}

func assertShutdownUnlock(t *testing.T, unlockShutdownC chan bool) {
	go func() {
		select {
		case <-unlockShutdownC:
			assert.True(t, true, "shutdown channel should unblock")
		// this case will happen only if the Shutdown function is not called
		// Shutdown function should be called always
		case <-time.After(testDeadline):
			assert.Fail(t, "test timeout for shutdown - THIS SHOULD NEVER HAPPEN")
		}
	}()
}

// runFunc test implementation of the Run function
// unlock channel is set in order to check whether the function has been called
func runFunc(t *testing.T, unlock chan bool) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		unlock <- true

		select {
		case <-ctx.Done():
			log.Info().Msg("context is done run shutdown")
			// execution ended successfully
		case <-time.After(10 * time.Second):
			log.Info().Msg("Should not happen - for run func")
			assert.Fail(t, "Run function overtime!")
		}
		return nil
	}
}

// shutdownFunc test implementation of the shutdown function
// unlock channel is set in order to check whether the function has been called
// shouldTimeout if true - the function should timeout and the exec. should end with GrateTimeout
func shutdownFunc(
	t *testing.T,
	unlock chan bool,
	shouldTimeout bool,
) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		unlock <- true

		if shouldTimeout {
			select {
			case <-ctx.Done(): // this should happen after `gracefulTimeout` -> which is 1 second
				log.Info().Msg("context is done for shutdown")
				assert.True(t, true, "Shutdown context should be done")
			case <-time.After(5 * time.Second): // 5 second is >> then 1 second (gracefulTimeout)
				log.Info().Msg("Should not happen - for shutdown func")
				assert.Fail(t, "Shutdown function overtime!")
			}
		}

		return nil
	}
}
