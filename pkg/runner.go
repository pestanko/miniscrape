package pkg

import (
	"context"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/pestanko/miniscrape/pkg/cache"
	"github.com/pestanko/miniscrape/pkg/config"
	"github.com/pestanko/miniscrape/pkg/resolvers"
	"github.com/pestanko/miniscrape/pkg/utils"
)

// NewAsyncRunner instance of the new asynchronous runner
func NewAsyncRunner(
	cfg *config.AppConfig,
	cats []config.Category,
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
	// Run the runner to get pages contant
	// based on the selector
	Run(selector config.RunSelector) []config.RunResult
}

type asyncRunner struct {
	cfg        *config.AppConfig
	categories []config.Category
	cache      cache.Cache
}

func (a *asyncRunner) Run(selector config.RunSelector) []config.RunResult {
	log.Debug().Msg("Runner Started!")
	pages := a.filterPages(selector)
	numberOfPages := len(pages)
	if numberOfPages == 0 {
		log.Warn().Msg("No pages available")
		return []config.RunResult{}
	}

	log.Debug().Int("numberOfPages", numberOfPages).Msg("Processing number of pages")
	channelWithResults := make(chan config.RunResult, numberOfPages)
	ctx := context.Background()
	// start async tasks
	a.startAsyncRequests(ctx, channelWithResults, pages)
	// collect results
	resultsCollection := a.collectResults(channelWithResults, numberOfPages)
	log.Debug().Msg("Runner Ended")

	return resultsCollection
}

func (a *asyncRunner) collectResults(channelWithResults chan config.RunResult, numberOfPages int) []config.RunResult {
	var resultsCollection []config.RunResult
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
	resChan chan<- config.RunResult,
	pages []config.Page,
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

func (a *asyncRunner) filterPages(sel config.RunSelector) []config.Page {
	var result []config.Page
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
