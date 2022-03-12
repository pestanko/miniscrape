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

type RunResult struct {
	Page    *config.Page
	Content string
	Status  RunResultStatus
}

func NewAsyncRunner(cfg *config.AppConfig) Runner {
	return &asyncRunner{cfg: cfg}
}

type Runner interface {
	Run() []RunResult
}

type asyncRunner struct {
	cfg *config.AppConfig
}

func (a asyncRunner) Run() []RunResult {
	log.Println("Runner Started!")
	numberOfPages := len(a.cfg.Pages)
	if numberOfPages == 0 {
		log.Fatalln("No pages available")
	}
	log.Printf("Processing number of pages: %d", numberOfPages)
	channelWithResults := make(chan RunResult, numberOfPages)
	ctx := context.Background()
	// start async tasks
	a.startAsyncRequests(channelWithResults, ctx)
	// collect results
	resultsCollection := a.collectResults(channelWithResults)
	log.Println("Runner Ended")

	return resultsCollection
}

func (a asyncRunner) collectResults(channelWithResults chan RunResult) []RunResult {
	var resultsCollection []RunResult
	for res := range channelWithResults {
		resultsCollection = append(resultsCollection, res)
		if len(resultsCollection) == len(a.cfg.Pages) {
			break
		}
	}

	return resultsCollection
}

func (a asyncRunner) startAsyncRequests(resChan chan<- RunResult, ctx context.Context) {
	for idx, page := range a.cfg.Pages {
		idx := idx
		page := page
		go func() {
			log.Printf("%03d. Starting to Resolve \"%s\"", idx, page.CodeName)
			resolver := NewGetPageResolver(&page)
			resChan <- resolver.Resolve(ctx)
		}()
	}
}
