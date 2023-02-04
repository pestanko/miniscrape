package resolvers

import (
	"context"
	"fmt"
	config2 "github.com/pestanko/miniscrape/internal/config"
)

type urlOnlyResolver struct {
	page config2.Page
}

func (u *urlOnlyResolver) Resolve(_ context.Context) config2.RunResult {
	return config2.RunResult{
		Page:    u.page,
		Content: fmt.Sprintf("URL for %s menu: %s", u.page.Name, u.page.URL),
		Status:  config2.RunSuccess,
		Kind:    "url",
	}
}

func makeErrorResult(page config2.Page, err error) config2.RunResult {
	return config2.RunResult{
		Page:    page,
		Content: fmt.Sprintf("Error: %v\n", err),
		Status:  config2.RunError,
		Kind:    "error",
	}
}
