package pkg

import (
	"context"
	config2 "github.com/pestanko/miniscrape/internal/config"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/pestanko/miniscrape/pkg/cache"
	"github.com/pestanko/miniscrape/pkg/resolvers"
	"github.com/pestanko/miniscrape/pkg/utils"
)

// NewAsyncRunner instance of the new asynchronous runner
func NewAsyncRunner(
	cfg *config2.AppConfig,
	cats []config2.Category,
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
	Run(ctx context.Context, selector config2.RunSelector) []config2.RunResult
}

type asyncRunner struct {
	cfg        *config2.AppConfig
	categories []config2.Category
	cache      cache.Cache
}

func (a *asyncRunner) Run(ctx context.Context, selector config2.RunSelector) []config2.RunResult {
	log.Debug().Msg("Runner Started!")
	pages := a.filterPages(selector)
	numberOfPages := len(pages)
	ll := log.With().
		Int("number_of_pages", numberOfPages).
		Interface("selector", selector).
		Logger()

	if numberOfPages == 0 {
		ll.Warn().Msg("No pages available")
		return []config2.RunResult{}
	}

	log.Debug().Int("numberOfPages", numberOfPages).Msg("Processing number of pages")
	channelWithResults := make(chan config2.RunResult, numberOfPages)
	// start async tasks
	a.startAsyncRequests(ctx, channelWithResults, pages)
	// collect results
	resultsCollection := a.collectResults(channelWithResults, numberOfPages)
	log.Debug().Msg("Runner Ended")

	return resultsCollection
}

func (a *asyncRunner) collectResults(channelWithResults chan config2.RunResult, numberOfPages int) []config2.RunResult {
	var resultsCollection []config2.RunResult
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
	resChan chan<- config2.RunResult,
	pages []config2.Page,
) {
	for idx, page := range pages {
		idx := idx
		page := page
		go func() {
			log.Debug().Int("idx", idx).Str("codename", page.CodeName).Msg("Starting to Resolve")
			resolver := resolvers.NewGetCachedPageResolver(page, a.cache)
			resChan <- resolver.Resolve(ctx)
		}()
	}
}

func (a *asyncRunner) filterPages(sel config2.RunSelector) []config2.Page {
	var result []config2.Page
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
