package pkg

import (
	"time"

	"github.com/pestanko/miniscrape/pkg/cache"
	"github.com/pestanko/miniscrape/pkg/config"
	"github.com/pestanko/miniscrape/pkg/utils"
)

// Service main service representation
type Service struct {
	Cfg        config.AppConfig
	categories utils.CachedContainer[[]config.Category]
}

// NewService create a new instance of the service
func NewService(cfg *config.AppConfig) *Service {
	categoriesLoader := func() *[]config.Category {
		categories := config.LoadCategories(cfg)
		return &categories
	}

	return &Service{
		*cfg,
		utils.NewCachedContainer(categoriesLoader, 10*time.Minute),
	}
}

// Scrape the pages based on selector
func (s *Service) Scrape(selector config.RunSelector) []config.RunResult {
	runner := NewAsyncRunner(&s.Cfg, s.GetCategories(), s.getCache())
	return runner.Run(selector)
}

// InvalidateCache for the provided selector
func (s *Service) InvalidateCache(sel config.RunSelector) {
	s.getCache().Invalidate(sel)
}

// GetCategories get all categories
func (s *Service) GetCategories() []config.Category {
	return *s.categories.Get()
}

func (s *Service) getCache() cache.Cache {
	return cache.NewCache(s.Cfg.Cache, time.Now())
}
