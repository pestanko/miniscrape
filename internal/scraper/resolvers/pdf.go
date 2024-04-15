package resolvers

import (
	"context"

	"github.com/pestanko/miniscrape/internal/models"
)

type pdfResolver struct {
	page models.Page
}

func (u *pdfResolver) Resolve(_ context.Context) models.RunResult {
	return models.RunResult{
		Page:    u.page,
		Content: u.page.URL,
		Status:  models.RunSuccess,
		Kind:    "pdf",
	}
}
