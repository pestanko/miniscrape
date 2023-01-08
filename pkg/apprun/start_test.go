package apprun

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const deadlineTimeout = 5 * time.Second

func TestStart(t *testing.T) {
	t.Run("no start and stop function, close asap", func(t *testing.T) {
		params := StartParams{
			Start:        nil,
			Stop:         nil,
			GraceTimeout: 0,
		}

		appCtx, cancelCallback := context.WithCancel(context.Background())
		cancelCallback()

		errC, err := Start(appCtx, params)
		assert.NoError(t, err)

		select {
		case err := <-errC:
			assert.NoError(t, err)
		case <-time.After(deadlineTimeout):
			assert.Fail(t, "test timeout")
		}
	})

	t.Run("start and stop, with channels", func(t *testing.T) {
		appCtx, cancelCallback := context.WithCancel(context.Background())
		unlockStartChan := make(chan bool)
		unlockStopChan := make(chan bool)
		params := StartParams{
			Start:        startWithTimeout(t, unlockStartChan),
			Stop:         stopWithTimeout(t, unlockStopChan, false),
			GraceTimeout: 1 * time.Second,
		}

		go func() {
			select {
			case <-unlockStartChan:
				assert.True(t, true, "start channel should unblock")
				cancelCallback()
			case <-time.After(deadlineTimeout):
				assert.Fail(t, "test timeout for start")
			}
		}()

		go func() {
			select {
			case <-unlockStopChan:
				assert.True(t, true, "stop channel should unblock")
			case <-time.After(deadlineTimeout):
				assert.Fail(t, "test timeout for stop")
			}
		}()

		errC, err := Start(appCtx, params)
		assert.NoError(t, err)

		select {
		case err := <-errC:
			assert.NoError(t, err)
		case <-time.After(deadlineTimeout):
			assert.Fail(t, "test timeout")
		}
	})

	t.Run("start context timeout", func(t *testing.T) {
		appCtx, cancelCallback := context.WithTimeout(context.Background(), 500*time.Millisecond)
		unlockStartChan := make(chan bool)
		unlockStopChan := make(chan bool)
		params := StartParams{
			Start:        startWithTimeout(t, unlockStartChan),
			Stop:         stopWithTimeout(t, unlockStopChan, false),
			GraceTimeout: 1 * time.Second,
		}

		go func() {
			select {
			case <-unlockStartChan:
				assert.True(t, true, "start channel should unblock")
				cancelCallback()
			case <-time.After(deadlineTimeout):
				assert.Fail(t, "test timeout for start")
			}
		}()

		go func() {
			select {
			case <-unlockStopChan:
				assert.True(t, true, "stop channel should unblock")
			case <-time.After(deadlineTimeout):
				assert.Fail(t, "test timeout for stop")
			}
		}()

		errC, err := Start(appCtx, params)
		assert.NoError(t, err)

		select {
		case err := <-errC:
			assert.NoError(t, err)
		case <-time.After(deadlineTimeout):
			assert.Fail(t, "test timeout")
		}
	})

	t.Run("stop graceful timeout", func(t *testing.T) {
		appCtx, cancelCallback := context.WithTimeout(context.Background(), 500*time.Millisecond)
		unlockStartChan := make(chan bool)
		unlockStopChan := make(chan bool)
		params := StartParams{
			Start:        startWithTimeout(t, unlockStartChan),
			Stop:         stopWithTimeout(t, unlockStopChan, true),
			GraceTimeout: 800 * time.Millisecond,
		}

		go func() {
			select {
			case <-unlockStartChan:
				assert.True(t, true, "start channel should unblock")
				cancelCallback()
			case <-time.After(deadlineTimeout):
				assert.Fail(t, "test timeout for start")
			}
		}()

		go func() {
			select {
			case <-unlockStopChan:
				assert.True(t, true, "stop channel should unblock")
			case <-time.After(deadlineTimeout):
				assert.Fail(t, "test timeout for stop")
			}
		}()

		errC, err := Start(appCtx, params)
		assert.NoError(t, err)

		select {
		case err := <-errC:
			assert.Error(t, err)
		case <-time.After(deadlineTimeout):
			assert.Fail(t, "test timeout")
		}
	})

	t.Run("start returns error", func(t *testing.T) {
		appCtx, cancelCallback := context.WithTimeout(context.Background(), 500*time.Millisecond)
		defer cancelCallback()

		unlockStopChan := make(chan bool)
		startErr := fmt.Errorf("some start error")
		params := StartParams{
			Start: func(ctx context.Context) error {
				return startErr
			},
			Stop:         stopWithTimeout(t, unlockStopChan, true),
			GraceTimeout: 800 * time.Millisecond,
		}

		go func() {
			select {
			case <-unlockStopChan:
				assert.True(t, true, "stop channel should unblock")
			case <-time.After(deadlineTimeout):
				assert.Fail(t, "test timeout for stop")
			}
		}()

		errC, err := Start(appCtx, params)
		assert.NoError(t, err)

		select {
		case err := <-errC:
			assert.ErrorIs(t, err, startErr)
		case <-time.After(deadlineTimeout):
			assert.Fail(t, "test timeout")
		}
	})

	t.Run("stop error", func(t *testing.T) {
		appCtx, cancelCallback := context.WithTimeout(context.Background(), 500*time.Millisecond)
		unlockStartChan := make(chan bool)
		unlockStopChan := make(chan bool)
		params := StartParams{
			Start:        startWithTimeout(t, unlockStartChan),
			Stop:         stopWithTimeout(t, unlockStartChan, true),
			GraceTimeout: 800 * time.Millisecond,
		}

		go func() {
			select {
			case <-unlockStartChan:
				assert.True(t, true, "start channel should unblock")
				cancelCallback()
			case <-time.After(deadlineTimeout):
				assert.Fail(t, "test timeout for start")
			}
		}()

		go func() {
			select {
			case <-unlockStopChan:
				assert.True(t, true, "stop channel should unblock")
			case <-time.After(deadlineTimeout):
				assert.Fail(t, "test timeout for stop")
			}
		}()

		errC, err := Start(appCtx, params)
		assert.NoError(t, err)

		select {
		case err := <-errC:
			assert.Error(t, err)
		case <-time.After(deadlineTimeout):
			assert.Fail(t, "test timeout")
		}
	})
}

func startWithTimeout(t *testing.T, unlock chan bool) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		unlock <- true
		time.Sleep(10 * time.Second)
		assert.Fail(t, "Start overtime!")
		return nil
	}
}

func stopWithTimeout(t *testing.T, unlock chan bool, sleep bool) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		unlock <- true
		if sleep {
			time.Sleep(2 * time.Second)
			assert.Fail(t, "Stop overtime!")
		}
		return nil
	}
}
