package resolvers

import (
	"context"
	"fmt"

	"github.com/pestanko/miniscrape/scraper/config"
)

type urlOnlyResolver struct {
	page config.Page
}

func (u *urlOnlyResolver) Resolve(_ context.Context) config.RunResult {
	return config.RunResult{
		Page:    u.page,
		Content: fmt.Sprintf("URL for %s menu: %s", u.page.Name, u.page.URL),
		Status:  config.RunSuccess,
		Kind:    "url",
	}
}

func makeErrorResult(page config.Page, err error) config.RunResult {
	return config.RunResult{
		Page:    page,
		Content: fmt.Sprintf("Error: %v\n", err),
		Status:  config.RunError,
		Kind:    "error",
	}
}
