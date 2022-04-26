package scraper

import (
	"github.com/pestanko/miniscrape/scraper/cache"
	"github.com/pestanko/miniscrape/scraper/config"
	"time"
)

type Service struct {
	Cfg        config.AppConfig
	Categories []config.Category
	cache      cache.Cache
}

func NewService(cfg *config.AppConfig) *Service {
	categories := config.LoadCategories(cfg)
	cache := cache.NewCache(cfg.Cache, time.Now())
	return &Service{
		*cfg,
		categories,
		cache,
	}
}

func (s *Service) Scrape(selector config.RunSelector) []RunResult {
	runner := NewAsyncRunner(&s.Cfg, s.Categories, s.cache)
	return runner.Run(selector)
}

func (s *Service) InvalidateCache(sel config.RunSelector) {
	s.cache.Invalidate(sel)
}
