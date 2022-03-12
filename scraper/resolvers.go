package scraper

import (
	"bytes"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
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

func NewGetPageResolver(page *config.Page) PageResolver {
	return &pageResolvedGet{
		page: page,
		client: http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

type pageResolvedGet struct {
	page   *config.Page
	client http.Client
}

func (r pageResolvedGet) Resolve(ctx context.Context) RunResult {
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

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(bodyContent))
	if err != nil {
		log.Printf("Parsing failed for (url: \"%s\"): %v\n", r.page.Url, err)
		return makeErrorResult(r.page, err)
	}

	contentArray := r.selectByQueryString(doc)
	content := strings.Join(contentArray, "\n")
	content = r.applyFilters(content)
	log.Printf("Alvin resolved!")

	return RunResult{
		Page:    r.page,
		Status:  RunSuccess,
		Content: content,
	}
}

func (r pageResolvedGet) selectByQueryString(doc *goquery.Document) []string {
	var content []string
	doc.Find(r.page.Query).Each(func(idx int, selection *goquery.Selection) {
		text, err := r.htmlToText(selection)
		if err != nil {
			log.Printf("Text extraction failed for (url: \"%s\"): %v\n", r.page.Url, err)
			return
		}
		content = append(content, normalizeString(text))
	})

	return content
}

func (r pageResolvedGet) htmlToText(selection *goquery.Selection) (string, error) {
	htmlContent, err := selection.Html()
	if err != nil {
		return "", err
	}
	text, err := html2text.FromString(htmlContent)
	if err != nil {
		log.Printf("Text extraction failed for (url: \"%s\"): %v\n", r.page.Url, err)
		return "", err
	}
	return text, err
}

func (r pageResolvedGet) applyFilters(content string) string {
	filters := []func(*config.Page) PageFilter{
		NewCutFilter,
		NewDayFilter,
	}
	newContent := content
	for _, newFilter := range filters {
		filter := newFilter(r.page)
		newContent, _ = filter.Filter(newContent)
	}
	return newContent
}

func makeErrorResult(page *config.Page, err error) RunResult {
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
