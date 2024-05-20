package resolvers

import (
	"context"

	"github.com/pestanko/miniscrape/internal/cache"
	"github.com/pestanko/miniscrape/internal/models"

	"github.com/rs/zerolog/log"
)

// NewGetCachedPageResolver a new instance of the cached resolver
func NewGetCachedPageResolver(page models.Page, cacheInstance cache.Cache) PageResolver {
	inner := NewPageResolver(page)
	if cacheInstance == nil {
		return inner
	}
	return &cachedPageResolver{
		resolver: inner,
		cache:    cacheInstance,
		page:     page,
	}
}

type cachedPageResolver struct {
	resolver PageResolver
	cache    cache.Cache
	page     models.Page
}

func (c *cachedPageResolver) Resolve(ctx context.Context) models.RunResult {
	namespace := cache.NewNamespace(c.page.Category, c.page.CodeName)
	if c.cache.IsPageCached(namespace) {
		log.Debug().Str("pageNamespace", c.page.Namespace()).Msg("Loading content from cache")

		content := string(c.cache.GetContent(cache.Item{
			Namespace: namespace,
		}))
		return models.RunResult{
			Page:    c.page,
			Content: content,
			Status:  models.RunSuccess,
		}
	}

	res := c.resolver.Resolve(ctx)
	if res.Status != models.RunSuccess {
		return res
	}

	err := c.cache.Store(cache.Item{
		Namespace:   namespace,
		CachePolicy: c.page.CachePolicy,
	}, []byte(res.Content))
	if err != nil {
		return makeErrorResult(c.page, err)
	}

	return res
}
