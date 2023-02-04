package resolvers

import (
	"context"
	config2 "github.com/pestanko/miniscrape/internal/config"

	"github.com/pestanko/miniscrape/pkg/cache"
	"github.com/rs/zerolog/log"
)

// NewGetCachedPageResolver a new instance of the cached resolver
func NewGetCachedPageResolver(page config2.Page, cacheInstance cache.Cache) PageResolver {
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
	page     config2.Page
}

func (c *cachedPageResolver) Resolve(ctx context.Context) config2.RunResult {
	namespace := cache.NewNamespace(c.page.Category, c.page.CodeName)
	if c.cache.IsPageCached(namespace) {
		log.Debug().Str("page", c.page.Namespace()).Msg("Loading content from cache")

		content := string(c.cache.GetContent(cache.Item{
			Namespace: namespace,
		}))
		return config2.RunResult{
			Page:    c.page,
			Content: content,
			Status:  config2.RunSuccess,
		}
	}

	res := c.resolver.Resolve(ctx)
	if res.Status != config2.RunSuccess {
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
