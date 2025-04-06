// Package scraper represents a main application logic
package scraper

import (
	"context"
	"strings"

	"github.com/pestanko/miniscrape/internal/cache"
	"github.com/pestanko/miniscrape/internal/config"
	"github.com/pestanko/miniscrape/internal/models"
	"github.com/pestanko/miniscrape/internal/scraper/resolvers"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

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
	// create a new span
	span := trace.SpanFromContext(ctx)
	span.AddEvent("Runner Started")
	span.SetAttributes(attribute.String("category", selector.Category))
	span.SetAttributes(attribute.String("page", selector.Page))
	span.SetAttributes(attribute.String("tags", strings.Join(selector.Tags, ",")))

	defer func() {
		span.AddEvent("Runner Ended")
		span.End()
	}()

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
			span := trace.SpanFromContext(ctx)
			span.AddEvent("start page resolve")
			span.SetAttributes(
				attribute.String("page", page.CodeName),
				attribute.String("namespace", page.Namespace()),
				attribute.String("resolver", page.Resolver),
				attribute.String("url", page.URL),
			)

			defer func() {
				span.AddEvent("end page resolve")
				span.End()
			}()

			llPage := zerolog.Dict().Str("codename", page.CodeName).Str("namespace", page.Namespace()).Str("url", page.URL)
			ll := zerolog.Ctx(ctx).With().
				Dict("page", llPage).Logger()
			ll.Debug().Int("idx", idx).Msg("Starting to Resolve")
			ctx := ll.WithContext(ctx)
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
