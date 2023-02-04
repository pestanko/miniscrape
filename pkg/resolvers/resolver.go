package resolvers

import (
	"context"
	"github.com/pestanko/miniscrape/internal/config"

	"github.com/pestanko/miniscrape/pkg/filters"
)

// PageResolver is a main interface for page content resolvers
type PageResolver interface {
	Resolve(ctx context.Context) config.RunResult
}

// NewPageResolver creates a new instance of the page resovler
func NewPageResolver(page config.Page) PageResolver {
	switch page.Resolver {
	case "url_only", "urlonly", "url-only":
		return &urlOnlyResolver{
			page: page,
		}
	case "image", "img":
		return &imageResolver{
			page:   page,
			client: httpClient,
		}
	case "get", "default":
		fallthrough
	default:
		return &pageContentResolver{
			page:   page,
			client: httpClient,
			filters: []func(*config.Page) filters.PageFilter{
				filters.NewHTMLToMdConverter,
				filters.NewNewLineTrimConverter,
				filters.NewCutFilter,
				filters.NewDayFilter,
				filters.NewCutLineFilter,
			},
		}
	}
}
