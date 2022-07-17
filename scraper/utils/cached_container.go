package utils

import (
	"time"
)

// NewCachedContainer create a new instance of the cached container
func NewCachedContainer[T any](
	contentProvider func() *T,
	duration time.Duration,
) CachedContainer[T] {
	return CachedContainer[T]{
		content:         nil,
		contentProvider: contentProvider,
		updateAt:        time.Time{},
		duration:        duration,
	}
}

// CachedContainer is a type of the container which contains a pointer
// to a spec. type
type CachedContainer[T any] struct {
	content *T

	contentProvider func() *T
	updateAt        time.Time
	duration        time.Duration
}

// Get a value inside the container
func (c *CachedContainer[T]) Get() *T {
	now := time.Now()
	updatedPlusDuration := c.updateAt.Add(c.duration)

	if now.After(updatedPlusDuration) {
		c.content = nil
	}

	if c.content == nil {
		c.Update()
	}

	return c.content
}

// Clear the value in the container
func (c *CachedContainer[T]) Clear() {
	c.content = nil
}

// Update the value in the container
func (c *CachedContainer[T]) Update() {
	c.content = c.contentProvider()
	c.updateAt = time.Now()
}
