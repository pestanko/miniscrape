package resolvers

import (
	"context"
	"github.com/pestanko/miniscrape/internal/models"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
)

type imageResolver struct {
	page   models.Page
	client http.Client
}

// Resolve implements PageResolver
func (r *imageResolver) Resolve(ctx context.Context) models.RunResult {
	bodyContent, err := getContentForWebPage(&r.page)
	if err != nil {
		return makeErrorResult(r.page, err)
	}

	contentArray, err := parseWebPageContent(&r.page, bodyContent)
	if err != nil {
		log.Warn().
			Err(err).
			Str("page", r.page.Namespace()).
			Str("pageUrl", r.page.URL).
			Msg("Content parsing failed")

		return makeErrorResult(r.page, err)
	}

	content := strings.Join(contentArray, "")

	return models.RunResult{
		Page:    r.page,
		Content: content,
		Status:  models.RunSuccess,
		Kind:    "img",
	}
}
