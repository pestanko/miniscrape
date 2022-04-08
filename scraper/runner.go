package scraper

import (
	"context"
	"github.com/pestanko/miniscrape/scraper/config"
	"log"
)

type RunResultStatus string

const (
	RunSuccess RunResultStatus = "ok"
	RunError   RunResultStatus = "error"
)

type RunSelector struct {
	Tags     []string
	Category string
}

type RunResult struct {
	Page    config.Page
	Content string
	Status  RunResultStatus
}

func NewAsyncRunner(cfg *config.AppConfig, categories []config.Category) Runner {
	return &asyncRunner{
		cfg:        cfg,
		categories: categories,
	}
}

type Runner interface {
	Run(selector RunSelector) []RunResult
}

type asyncRunner struct {
	cfg        *config.AppConfig
	categories []config.Category
}

func (a *asyncRunner) Run(selector RunSelector) []RunResult {
	log.Println("Runner Started!")
	pages := a.filterPages(selector)
	numberOfPages := len(pages)
	if numberOfPages == 0 {
		log.Fatalln("No pages available")
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
			resolver := NewGetPageResolver(page)
			resChan <- resolver.Resolve(ctx)
		}()
	}
}

func (a *asyncRunner) filterPages(sel RunSelector) []config.Page {
	var result []config.Page
	tagsResolver := MakeTagsResolver(sel.Tags)
	for _, category := range a.categories {
		// filter out the category
		if sel.Category != "" && category.Name != sel.Category {
			continue
		}

		for _, page := range category.Pages {
			if tagsResolver.isMatch(page.Tags) {
				result = append(result, page)
			}
		}
	}

	return result
}
