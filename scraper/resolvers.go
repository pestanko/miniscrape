package scraper

import (
	"bytes"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/antchfx/htmlquery"
	"github.com/pestanko/miniscrape/scraper/cache"
	"github.com/pestanko/miniscrape/scraper/config"
	"io"
	"jaytaylor.com/html2text"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

type PageResolver interface {
	Resolve(ctx context.Context) RunResult
}

func NewPageResolver(page config.Page) PageResolver {
	switch page.Resolver {
	case "url_only", "urlonly", "url-only":
		return &urlOnlyResolver{
			page: page,
		}
	case "get", "default":
		fallthrough
	default:
		return &pageResolvedGet{
			page: page,
			client: http.Client{
				Timeout: 30 * time.Second,
			},
		}
	}
}

func NewGetCachedPageResolver(page config.Page, cacheInstance cache.Cache) PageResolver {
	inner := NewPageResolver(page)
	if cacheInstance == nil {
		return inner
	}
	return &cachedPageResolver{
		resolver: inner,
		cache:    cacheInstance,
		page:     page,
	}
}

type cachedPageResolver struct {
	resolver PageResolver
	cache    cache.Cache
	page     config.Page
}

func (c *cachedPageResolver) Resolve(ctx context.Context) RunResult {
	if c.cache.IsPageCached(c.page.CodeName) {
		log.Printf("Loading content from cache '%s'", c.page.CodeName)
		content := string(c.cache.GetContent(cache.Item{
			PageName:     c.page.CodeName,
			CategoryName: c.page.Category,
		}))
		return RunResult{
			Page:    c.page,
			Content: content,
			Status:  RunSuccess,
		}
	}

	res := c.resolver.Resolve(ctx)
	if res.Status != RunSuccess {
		return res
	}

	err := c.cache.Store(cache.Item{
		PageName:     c.page.CodeName,
		CategoryName: c.page.Category,
		CachePolicy:  c.page.CachePolicy,
	}, []byte(res.Content))
	if err != nil {
		return makeErrorResult(c.page, err)
	}

	return res
}

type pageResolvedGet struct {
	page   config.Page
	client http.Client
}

func (r *pageResolvedGet) Resolve(ctx context.Context) RunResult {
	res, err := r.client.Get(r.page.Url)
	if err != nil {
		log.Printf("Request failed for (url: \"%s\"): %v\n", r.page.Url, err)
		log.Printf("Error[%d]: %v", res.StatusCode, res)
		return makeErrorResult(r.page, err)
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Printf("Unable to close body: %v", err)
		}
	}()

	bodyContent, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Failed to read a body for (url: \"%s\"): %v\n", r.page.Url, err)
		return makeErrorResult(r.page, err)
	}

	var contentArray []string

	if r.page.Query != "" {
		contentArray, err = r.parseUsingCssQuery(bodyContent)
	} else {
		contentArray, err = r.parseUsingXPathQuery(bodyContent)
	}

	if err != nil {
		log.Printf("Parsing failed for (url: \"%s\"): %v\n", r.page.Url, err)
		return makeErrorResult(r.page, err)
	}
	content := r.applyFilters(contentArray)
	log.Printf("%s resolved!", r.page.CodeName)

	return RunResult{
		Page:    r.page,
		Status:  RunSuccess,
		Content: content,
	}
}

func (r *pageResolvedGet) parseUsingXPathQuery(content []byte) ([]string, error) {
	root, err := htmlquery.Parse(bytes.NewReader(content))
	if err != nil {
		return []string{}, err
	}
	nodes, err := htmlquery.QueryAll(root, r.page.XPath)
	if err != nil {
		return []string{}, err
	}

	var result []string

	for _, node := range nodes {
		html := htmlquery.OutputHTML(node, true)
		text, err := r.htmlToText(html)
		if err != nil {
			continue
		}
		result = append(result, text)
	}

	return result, nil
}

func (r *pageResolvedGet) parseUsingCssQuery(bodyContent []byte) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(bodyContent))
	if err != nil {
		return []string{}, err
	}

	var content []string
	doc.Find(r.page.Query).Each(func(idx int, selection *goquery.Selection) {
		htmlContent, err := selection.Html()
		if err != nil {
			log.Printf("Text extraction failed for (url: \"%s\"): %v\n", r.page.Url, err)
			return
		}
		text, err := r.htmlToText(htmlContent)
		if err != nil {
			log.Printf("Text extraction failed for (url: \"%s\"): %v\n", r.page.Url, err)
			return
		}
		content = append(content, text)
	})

	return content, nil
}

func (r *pageResolvedGet) htmlToText(htmlContent string) (string, error) {
	log.Printf("Found by query: %s", htmlContent)
	text, err := html2text.FromString(htmlContent, html2text.Options{
		PrettyTables: r.page.Filters.Html.PrettyTables,
		TextOnly:     r.page.Filters.Html.TextOnly,
	})
	if err != nil {
		log.Printf("Text extraction failed for (url: \"%s\"): %v\n", r.page.Url, err)
		return "", err
	}
	return normalizeString(text), nil
}

func (r *pageResolvedGet) applyFilters(contentArray []string) string {
	content := strings.Join(contentArray, "\n")
	if strings.TrimSpace(content) == "" {
		return ""
	}
	filters := []func(*config.Page) PageFilter{
		NewCutFilter,
		NewDayFilter,
		NewCutLineFilter,
	}
	newContent := content
	for _, newFilter := range filters {
		filter := newFilter(&r.page)
		if !filter.IsEnabled() {
			continue
		}
		newContent, _ = filter.Filter(newContent)
		if newContent == "" {
			return ""
		}
	}
	return newContent
}

type urlOnlyResolver struct {
	page config.Page
}

func (u *urlOnlyResolver) Resolve(ctx context.Context) RunResult {
	return RunResult{
		Page:    u.page,
		Content: fmt.Sprintf("Url for menu: %s", u.page.Url),
		Status:  RunSuccess,
	}
}

func makeErrorResult(page config.Page, err error) RunResult {
	return RunResult{
		Page:    page,
		Content: fmt.Sprintf("Error: %v\n", err),
		Status:  RunError,
	}
}

func normalizeString(content string) string {
	var pat = regexp.MustCompile("\n\n")
	return pat.ReplaceAllString(content, "\n")
}
