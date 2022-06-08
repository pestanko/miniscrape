package scraper

import (
	"context"
	"log"
	"strings"

	"github.com/pestanko/miniscrape/scraper/cache"
	"github.com/pestanko/miniscrape/scraper/config"
	"github.com/pestanko/miniscrape/scraper/utils"
)

type RunResultStatus string

const (
	RunSuccess RunResultStatus = "ok"
	RunError   RunResultStatus = "error"
	RunEmpty   RunResultStatus = "empty"
)

type RunResult struct {
	Page    config.Page
	Content string
	Status  RunResultStatus
}

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

type Runner interface {
	Run(selector config.RunSelector) []RunResult
}

type asyncRunner struct {
	cfg        *config.AppConfig
	categories []config.Category
	cache      cache.Cache
}

func (a *asyncRunner) Run(selector config.RunSelector) []RunResult {
	log.Println("Runner Started!")
	pages := a.filterPages(selector)
	numberOfPages := len(pages)
	if numberOfPages == 0 {
		log.Println("No pages available")
		return []RunResult{}
	}

	log.Printf("Processing number of pages: %d", numberOfPages)
	channelWithResults := make(chan RunResult, numberOfPages)
	ctx := context.Background()
	// start async tasks
	a.startAsyncRequests(channelWithResults, ctx, pages)
	// collect results
	resultsCollection := a.collectResults(channelWithResults, numberOfPages)
	log.Println("Runner Ended")

	return resultsCollection
}

func (a *asyncRunner) collectResults(channelWithResults chan RunResult, numberOfPages int) []RunResult {
	var resultsCollection []RunResult
	for res := range channelWithResults {
		resultsCollection = append(resultsCollection, res)
		if len(resultsCollection) == numberOfPages {
			break
		}
	}

	return resultsCollection
}

func (a *asyncRunner) startAsyncRequests(resChan chan<- RunResult, ctx context.Context, pages []config.Page) {
	for idx, page := range pages {
		idx := idx
		page := page
		go func() {
			log.Printf("%03d. Starting to Resolve \"%s\"", idx, page.CodeName)
			resolver := NewGetCachedPageResolver(page, a.cache)
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
