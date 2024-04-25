package resolvers

import (
	"context"

	"github.com/pestanko/miniscrape/internal/models"
	"github.com/pestanko/miniscrape/internal/scraper/filters"
)

// PageResolver is a main interface for page content resolvers
type PageResolver interface {
	Resolve(ctx context.Context) models.RunResult
}

// NewPageResolver creates a new instance of the page resovler
func NewPageResolver(page models.Page) PageResolver {
	switch page.Resolver {
	case "url_only", "urlonly", "url-only":
		return &urlOnlyResolver{
			page: page,
		}
	case "url", "iframe":
		return &iframeResolver{
			page: page,
		}
	case "image", "img":
		return &imageResolver{
			page:   page,
			client: httpClient,
		}
	case "pdf":
		return &pdfResolver{
			page: page,
		}
	case "get", "default":
		fallthrough
	default:
		return &pageContentResolver{
			page:   page,
			client: httpClient,
			filters: []func(*models.Page) filters.PageFilter{
				filters.NewHTMLToMdConverter,
				filters.NewNewLineTrimConverter,
				filters.NewCutFilter,
				filters.NewDayFilter,
				filters.NewCutLineFilter,
			},
		}
	}
}
