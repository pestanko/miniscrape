package utils

import (
	"time"
)

type CachedContainer[T any] struct {
	content *T

	contentProvider func() *T
	updateAt        time.Time
	duration        time.Duration
}

func NewCachedContainer[T any](contentProvider func() *T, duration time.Duration) CachedContainer[T] {
	return CachedContainer[T]{
		content:         nil,
		contentProvider: contentProvider,
		updateAt:        time.Time{},
		duration:        duration,
	}
}

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

func (c *CachedContainer[T]) Clear() {
	c.content = nil
}

func (c *CachedContainer[T]) Update() {
	c.content = c.contentProvider()
	c.updateAt = time.Now()
}
