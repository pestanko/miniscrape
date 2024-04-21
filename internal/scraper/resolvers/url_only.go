package resolvers

import (
	"context"
	"fmt"

	"github.com/pestanko/miniscrape/internal/models"
)

type urlOnlyResolver struct {
	page models.Page
}

func (u *urlOnlyResolver) Resolve(_ context.Context) models.RunResult {
	return models.RunResult{
		Page:    u.page,
		Content: u.page.URL,
		Status:  models.RunSuccess,
		Kind:    "url",
	}
}

func makeErrorResult(page models.Page, err error) models.RunResult {
	return models.RunResult{
		Page:    page,
		Content: fmt.Sprintf("Error: %v\n", err),
		Status:  models.RunError,
		Kind:    "error",
	}
}

func makeEmptyResult(page models.Page, kind string) models.RunResult {
	return models.RunResult{
		Page:    page,
		Content: "",
		Status:  models.RunEmpty,
		Kind:    kind,
	}
}
