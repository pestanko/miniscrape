package utils

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCachedContainerInstanceForInteger(t *testing.T) {
	assert := assert.New(t)

	dur := 100 * time.Second
	value := 0

	ctx := t.Context()

	fun := func(_ context.Context) *int {
		value++
		return &value
	}

	cont := NewCachedContainer(fun, dur)

	assert.NotNil(cont.Get(ctx))
	assert.Equal(*cont.Get(ctx), 1)
	assert.Equal(*cont.Get(ctx), 1)
	assert.Equal(*cont.Get(ctx), 1)
}

func TestCachedContainerInstanceExpiration(t *testing.T) {
	assert := assert.New(t)

	dur := 100 * time.Millisecond
	value := 0

	fun := func(_ context.Context) *int {
		value++
		return &value
	}

	cont := NewCachedContainer(fun, dur)
	assert.Equal(*cont.Get(t.Context()), 1)

	time.Sleep(dur * 2)

	assert.Equal(*cont.Get(t.Context()), 2)
}

func TestCachedContainerClear(t *testing.T) {
	assert := assert.New(t)

	dur := 100 * time.Second
	value := 0

	fun := func(_ context.Context) *int {
		value++
		return &value
	}

	cont := NewCachedContainer(fun, dur)
	assert.Equal(*cont.Get(t.Context()), 1)
	cont.Clear()
	assert.Equal(*cont.Get(t.Context()), 2)
}

func TestCachedContainerUpdate(t *testing.T) {
	assert := assert.New(t)

	dur := 100 * time.Second
	value := 0

	fun := func(_ context.Context) *int {
		value++
		return &value
	}

	cont := NewCachedContainer(fun, dur)
	assert.Equal(*cont.Get(t.Context()), 1)
	cont.Update(t.Context())
	assert.Equal(*cont.Get(t.Context()), 2)
}
