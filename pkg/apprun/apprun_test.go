package apprun

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type emptyDeps struct {
	isCloseCalled bool
}

func (e *emptyDeps) Close() error {
	e.isCloseCalled = true
	return nil
}

type emptyDepsErr struct {
	isCloseCalled bool
}

func (e *emptyDepsErr) Close() error {
	e.isCloseCalled = true
	return errors.New("error while closing")
}

func TestCreateAppRunner(t *testing.T) {
	ctx := context.TODO()
	deps := &emptyDeps{isCloseCalled: false}

	t.Setenv("TRACING_ENABLED", "false")

	t.Run("run without error", func(t *testing.T) {
		s := assert.New(t)

		runner := NewAppRunner(WithDepProvider(func() (*emptyDeps, error) {
			return deps, nil
		}))

		err := runner.Run(ctx, func(ctx context.Context, d *emptyDeps) error {
			return nil
		})

		s.NoError(err)
		s.True(deps.isCloseCalled)
		s.False(runner.isTracingEnabled)
	})

	t.Run("run tracing enabled using env", func(t *testing.T) {
		s := assert.New(t)

		t.Setenv("TRACING_ENABLED", "true")

		runner := NewAppRunner(WithDepProvider(func() (*emptyDeps, error) {
			return deps, nil
		}))

		err := runner.Run(ctx, func(ctx context.Context, d *emptyDeps) error {
			return nil
		})

		s.NoError(err)
		s.True(deps.isCloseCalled)
		s.True(runner.isTracingEnabled)
	})

	t.Run("run tracing explicit disabled ignoring env", func(t *testing.T) {
		s := assert.New(t)

		t.Setenv("TRACING_ENABLED", "true")

		runner := NewAppRunner(WithDepProvider(func() (*emptyDeps, error) {
			return deps, nil
		}), WithForceTracingEnabled[*emptyDeps](false))

		err := runner.Run(ctx, func(ctx context.Context, d *emptyDeps) error {
			return nil
		})

		s.NoError(err)
		s.True(deps.isCloseCalled)
		s.False(runner.isTracingEnabled)
	})

	t.Run("run tracing enabled using with", func(t *testing.T) {
		s := assert.New(t)

		runner := NewAppRunner(
			WithDepProvider(func() (*emptyDeps, error) {
				return deps, nil
			}),
			WithForceTracingEnabled[*emptyDeps](true),
		)

		err := runner.Run(ctx, func(ctx context.Context, d *emptyDeps) error {
			return nil
		})

		s.NoError(err)
		s.True(deps.isCloseCalled)
		s.True(runner.isTracingEnabled)
	})

	t.Run("run with error while running it", func(t *testing.T) {
		s := assert.New(t)

		runner := NewAppRunner(WithDepProvider(func() (*emptyDeps, error) {
			return deps, nil
		}))

		err := runner.Run(ctx, func(ctx context.Context, d *emptyDeps) error {
			return errors.New("error while running the app")
		})

		s.Error(err)
		s.True(deps.isCloseCalled)
		s.False(runner.isTracingEnabled)
	})

	t.Run("run with error while creating deps", func(t *testing.T) {
		s := assert.New(t)

		runner := NewAppRunner(WithDepProvider(func() (*emptyDeps, error) {
			return nil, errors.New("error while creating deps")
		}))

		err := runner.Run(ctx, func(ctx context.Context, d *emptyDeps) error {
			return nil
		})

		s.Error(err)
		s.False(runner.isTracingEnabled)
	})

	t.Run("run with error while closing deps", func(t *testing.T) {
		s := assert.New(t)

		runner := NewAppRunner(WithDepProvider(func() (*emptyDepsErr, error) {
			return &emptyDepsErr{}, nil
		}))

		err := runner.Run(ctx, func(ctx context.Context, d *emptyDepsErr) error {
			return nil
		})

		s.NoError(err)
		s.False(runner.isTracingEnabled)
	})
}
