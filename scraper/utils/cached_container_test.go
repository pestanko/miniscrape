package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCachedContainerInstanceForInteger(t *testing.T) {
	assert := assert.New(t)

	dur := 100 * time.Second
	value := 0

	fun := func() *int {
		value++
		return &value
	}

	cont := NewCachedContainer(fun, dur)

	assert.NotNil(cont.Get())
	assert.Equal(*cont.Get(), 1)
	assert.Equal(*cont.Get(), 1)
	assert.Equal(*cont.Get(), 1)
}

func TestCachedContainerInstanceExpiration(t *testing.T) {
	assert := assert.New(t)

	dur := 100 * time.Millisecond
	value := 0

	fun := func() *int {
		value++
		return &value
	}

	cont := NewCachedContainer(fun, dur)
	assert.Equal(*cont.Get(), 1)

	time.Sleep(dur * 2)

	assert.Equal(*cont.Get(), 2)
}

func TestCachedContainerClear(t *testing.T) {
	assert := assert.New(t)

	dur := 100 * time.Second
	value := 0

	fun := func() *int {
		value++
		return &value
	}

	cont := NewCachedContainer(fun, dur)
	assert.Equal(*cont.Get(), 1)
	cont.Clear()
	assert.Equal(*cont.Get(), 2)
}

func TestCachedContainerUpdate(t *testing.T) {
	assert := assert.New(t)

	dur := 100 * time.Second
	value := 0

	fun := func() *int {
		value++
		return &value
	}

	cont := NewCachedContainer(fun, dur)
	assert.Equal(*cont.Get(), 1)
	cont.Update()
	assert.Equal(*cont.Get(), 2)
}
