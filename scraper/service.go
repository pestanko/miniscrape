package scraper

import "github.com/pestanko/miniscrape/scraper/config"

type Service struct {
	Cfg        config.AppConfig
	Categories []config.Category
}

func NewService(cfg *config.AppConfig) *Service {
	categories := config.LoadCategories(cfg)

	return &Service{
		*cfg,
		categories,
	}
}

func (s *Service) Scrape(selector RunSelector) []RunResult {
	runner := NewAsyncRunner(&s.Cfg, s.Categories)
	return runner.Run(selector)
}
