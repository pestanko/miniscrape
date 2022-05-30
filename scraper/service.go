package scraper

import (
	"time"

	"github.com/pestanko/miniscrape/scraper/cache"
	"github.com/pestanko/miniscrape/scraper/config"
	"github.com/pestanko/miniscrape/scraper/utils"
)

type Service struct {
	Cfg        config.AppConfig
	categories utils.CachedContainer[[]config.Category]
	cache      cache.Cache
}

func NewService(cfg *config.AppConfig) *Service {
	cache := cache.NewCache(cfg.Cache, time.Now())
	categoriesLoader := func() *[]config.Category {
		categories := config.LoadCategories(cfg)
		return &categories
	}

	return &Service{
		*cfg,
		utils.NewCachedContainer(categoriesLoader, 10*time.Minute),
		cache,
	}
}

func (s *Service) Scrape(selector config.RunSelector) []RunResult {
	runner := NewAsyncRunner(&s.Cfg, s.GetCategories(), s.cache)
	return runner.Run(selector)
}

func (s *Service) InvalidateCache(sel config.RunSelector) {
	s.cache.Invalidate(sel)
}

func (s *Service) GetCategories() []config.Category {
	return *s.categories.Get()
}
