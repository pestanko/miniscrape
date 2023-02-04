// Package scraper represents a main application logic
package scraper

import (
	"context"
	"github.com/pestanko/miniscrape/internal/cache"
	"github.com/pestanko/miniscrape/internal/config"
	"github.com/pestanko/miniscrape/internal/models"
	"github.com/pestanko/miniscrape/internal/scraper/resolvers"
	"github.com/rs/zerolog"
	"strings"

	"github.com/pestanko/miniscrape/pkg/utils"
)

// NewAsyncRunner instance of the new asynchronous runner
func NewAsyncRunner(
	cfg *config.AppConfig,
	cats []models.Category,
	cache cache.Cache,
) Runner {
	return &asyncRunner{
		cfg:        cfg,
		categories: cats,
		cache:      cache,
	}
}

// Runner interface
type Runner interface {
	// Run the runner to get pages content
	// based on the selector
	Run(ctx context.Context, selector models.RunSelector) []models.RunResult
}

type asyncRunner struct {
	cfg        *config.AppConfig
	categories []models.Category
	cache      cache.Cache
}

func (a *asyncRunner) Run(ctx context.Context, selector models.RunSelector) []models.RunResult {
	zerolog.Ctx(ctx).Debug().Msg("Runner Started!")
	pages := a.filterPages(selector)
	numberOfPages := len(pages)
	ll := zerolog.Ctx(ctx).
		With().
		Int("number_of_pages", numberOfPages).
		Interface("selector", selector).
		Logger()

	if numberOfPages == 0 {
		ll.Warn().Msg("No pages available")
		return []models.RunResult{}
	}

	ll.Debug().Msg("Processing number of pages")
	channelWithResults := make(chan models.RunResult, numberOfPages)
	// start async tasks
	a.startAsyncRequests(ctx, channelWithResults, pages)
	// collect results
	resultsCollection := a.collectResults(channelWithResults, numberOfPages)
	ll.Debug().Msg("Runner Ended")

	return resultsCollection
}

func (a *asyncRunner) collectResults(channelWithResults chan models.RunResult, numberOfPages int) []models.RunResult {
	var resultsCollection []models.RunResult
	for res := range channelWithResults {
		resultsCollection = append(resultsCollection, res)
		if len(resultsCollection) == numberOfPages {
			break
		}
	}

	return resultsCollection
}

func (a *asyncRunner) startAsyncRequests(
	ctx context.Context,
	resChan chan<- models.RunResult,
	pages []models.Page,
) {
	for idx, page := range pages {
		idx := idx
		page := page
		go func() {
			zerolog.Ctx(ctx).
				Debug().
				Int("idx", idx).
				Str("codename", page.CodeName).
				Msg("Starting to Resolve")
			resolver := resolvers.NewGetCachedPageResolver(page, a.cache)
			resChan <- resolver.Resolve(ctx)
		}()
	}
}

func (a *asyncRunner) filterPages(sel models.RunSelector) []models.Page {
	var result []models.Page
	tagsResolver := utils.MakeTagsResolver(sel.Tags)
	for _, category := range a.categories {
		// filter out the category

		if sel.Category != "" && category.Name != sel.Category {
			continue
		}

		for _, page := range category.Pages {
			if page.Disabled && !sel.Force {
				continue
			}

			if sel.Page != "" && !strings.Contains(page.CodeName, sel.Page) {
				continue
			}

			if tagsResolver.IsMatch(page.Tags) {
				result = append(result, page)
			}
		}
	}

	return result
}
