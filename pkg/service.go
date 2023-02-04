package pkg

import (
	"context"
	config2 "github.com/pestanko/miniscrape/internal/config"
	"time"

	"github.com/pestanko/miniscrape/pkg/cache"
	"github.com/pestanko/miniscrape/pkg/utils"
)

// Service main service representation
type Service struct {
	Cfg        config2.AppConfig
	categories utils.CachedContainer[[]config2.Category]
}

// NewService create a new instance of the service
func NewService(cfg *config2.AppConfig) *Service {
	categoriesLoader := func() *[]config2.Category {
		categories := config2.LoadCategories(cfg)
		return &categories
	}

	return &Service{
		*cfg,
		utils.NewCachedContainer(categoriesLoader, 10*time.Minute),
	}
}

// Scrape the pages based on selector
func (s *Service) Scrape(ctx context.Context, selector config2.RunSelector) []config2.RunResult {
	runner := NewAsyncRunner(&s.Cfg, s.GetCategories(), s.getCache())
	return runner.Run(ctx, selector)
}

// InvalidateCache for the provided selector
func (s *Service) InvalidateCache(sel config2.RunSelector) {
	s.getCache().Invalidate(sel)
}

// GetCategories get all categories
func (s *Service) GetCategories() []config2.Category {
	return *s.categories.Get()
}

func (s *Service) getCache() cache.Cache {
	return cache.NewCache(s.Cfg.Cache, time.Now())
}
