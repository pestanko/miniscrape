package resolvers

import (
	"context"
	"net/http"

	"github.com/pestanko/miniscrape/internal/models"
	"github.com/rs/zerolog"
)

type imageResolver struct {
	page   models.Page
	client http.Client
}

// Resolve implements PageResolver
func (r *imageResolver) Resolve(ctx context.Context) models.RunResult {

	ll := zerolog.Ctx(ctx).With().
		Interface("page",
			zerolog.Dict().
				Str("codename", r.page.CodeName).
				Str("url", r.page.URL).
				Str("namespace", r.page.Namespace()).
				Str("resolver", r.page.Resolver),
		).
		Logger()

	ll.Debug().Msg("Resolving manu")

	bodyContent, err := getContentForWebPage(ctx, &r.page)
	if err != nil {
		return makeErrorResult(r.page, err)
	}

	contentArray, err := ParseWebPageContent(ctx, &r.page, bodyContent)
	if err != nil {
		zerolog.Ctx(ctx).
			Warn().
			Err(err).
			Str("page", r.page.Namespace()).
			Str("pageUrl", r.page.URL).
			Msg("Content parsing failed")

		return makeErrorResult(r.page, err)
	}

	if len(contentArray) == 0 {
		ll.Warn().Msg("No content found")
		return makeEmptyResult(r.page, "img")
	}

	// Pick the first image
	content := getAttrValue(contentArray[0].Attrs, "src")

	return models.RunResult{
		Page:    r.page,
		Content: content,
		Status:  models.RunSuccess,
		Kind:    "img",
	}
}
