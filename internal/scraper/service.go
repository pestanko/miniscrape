package scraper

import (
	"context"
	"github.com/pestanko/miniscrape/internal/cache"
	"github.com/pestanko/miniscrape/internal/config"
	"github.com/pestanko/miniscrape/internal/models"
	"time"

	"github.com/pestanko/miniscrape/pkg/utils"
)

// Service main service representation
type Service struct {
	Cfg        config.AppConfig
	categories utils.CachedContainer[[]models.Category]
}

// NewService create a new instance of the service
func NewService(cfg *config.AppConfig) *Service {
	categoriesLoader := func() *[]models.Category {
		categories := models.LoadCategories(cfg)
		return &categories
	}

	return &Service{
		*cfg,
		utils.NewCachedContainer(categoriesLoader, 10*time.Minute),
	}
}

// Scrape the pages based on selector
func (s *Service) Scrape(ctx context.Context, selector models.RunSelector) []models.RunResult {
	runner := NewAsyncRunner(&s.Cfg, s.GetCategories(), s.getCache())
	return runner.Run(ctx, selector)
}

// InvalidateCache for the provided selector
func (s *Service) InvalidateCache(sel models.RunSelector) {
	s.getCache().Invalidate(sel)
}

// GetCategories get all categories
func (s *Service) GetCategories() []models.Category {
	return *s.categories.Get()
}

func (s *Service) getCache() cache.Cache {
	return cache.NewCache(s.Cfg.Cache, time.Now())
}
