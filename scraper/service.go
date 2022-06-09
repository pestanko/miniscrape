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
}

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

func (s *Service) Scrape(selector config.RunSelector) []RunResult {
	runner := NewAsyncRunner(&s.Cfg, s.GetCategories(), s.getCache())
	return runner.Run(selector)
}

func (s *Service) InvalidateCache(sel config.RunSelector) {
	s.getCache().Invalidate(sel)
}

func (s *Service) GetCategories() []config.Category {
	return *s.categories.Get()
}

func (s *Service) getCache() cache.Cache {
	return cache.NewCache(s.Cfg.Cache, time.Now())
}
