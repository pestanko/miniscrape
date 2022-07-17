package resolvers

import (
	"context"

	"github.com/pestanko/miniscrape/scraper/config"
	"github.com/pestanko/miniscrape/scraper/filters"
)

type PageResolver interface {
	Resolve(ctx context.Context) config.RunResult
}

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
