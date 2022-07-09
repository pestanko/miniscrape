package scraper

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/antchfx/htmlquery"
	"github.com/pestanko/miniscrape/scraper/cache"
	"github.com/pestanko/miniscrape/scraper/config"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246",
	`Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.27 Safari/537.36`,
	`Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36`,
}

var normPattern = regexp.MustCompile("\n\n")

var httpClient = http.Client{
	Timeout: 30 * time.Second,
}

type PageResolver interface {
	Resolve(ctx context.Context) RunResult
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
		return &pageResolverContent{
			page:   page,
			client: httpClient,
			filters: []func(*config.Page) PageFilter{
				NewHTMLToMdConverter,
				NewNewLineTrimConverter,
				NewCutFilter,
				NewDayFilter,
				NewCutLineFilter,
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
	namespace := cache.NewNamespace(c.page.Category, c.page.CodeName)
	if c.cache.IsPageCached(namespace) {
		log.Printf("Loading content from cache '%s'", c.page.CodeName)
		content := string(c.cache.GetContent(cache.Item{
			Namespace: namespace,
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
		Namespace:   namespace,
		CachePolicy: c.page.CachePolicy,
	}, []byte(res.Content))
	if err != nil {
		return makeErrorResult(c.page, err)
	}

	return res
}

type pageResolverContent struct {
	page    config.Page
	client  http.Client
	filters []func(*config.Page) PageFilter
}

func (r *pageResolverContent) Resolve(_ context.Context) RunResult {
	bodyContent, err := getContentForWebPage(&r.page)
	if err != nil {
		return makeErrorResult(r.page, err)
	}

	contentArray, err := parseWebPageContent(&r.page, bodyContent)
	if err != nil {
		log.Printf("Parsing failed for (url: \"%s\"): %v\n", r.page.Url, err)
		return makeErrorResult(r.page, err)
	}

	content := strings.Join(contentArray, "\n")
	content = r.applyFilters(content)

	var status = RunSuccess
	if content == "" {
		log.Printf("%s resolved but the content is empty", r.page.CodeName)
		status = RunEmpty
	} else {
		log.Printf("%s resolved!", r.page.CodeName)
	}

	return RunResult{
		Page:    r.page,
		Status:  status,
		Content: content,
		Kind:    "content",
	}
}

func getContentForWebPage(page *config.Page) (bodyContent []byte, err error) {
	if page.Command.Content.Name != "" {
		bodyContent, err = getContentByCommand(page)
	} else {
		bodyContent, err = getContentByRequest(page)
	}

	if err == nil {
		bodyContent = transformEncoding(bodyContent)
	}

	return
}

func getContentByCommand(page *config.Page) ([]byte, error) {
	// Use command
	log.Printf("Using command: '%s' with args %v", page.Command.Content.Name, page.Command.Content.Args)
	var outb, errb bytes.Buffer
	cmd := exec.Command(page.Command.Content.Name, page.Command.Content.Args...)
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		log.Printf("Command error[%v]: %s", err, errb.String())
	}

	return outb.Bytes(), err
}

func getContentByRequest(page *config.Page) ([]byte, error) {
	req, err := http.NewRequest("GET", page.Url, nil)
	if err != nil {
		log.Printf("Request creation failed for (url: \"%s\"): %v\n", page.Url, err)
		return []byte{}, err
	}

	randomUserAgent := userAgents[rand.Intn(len(userAgents))]
	req.Header.Add("User-Agent", randomUserAgent)

	res, err := httpClient.Do(req)
	if err != nil {
		log.Printf("Request failed for (url: \"%s\"): %v\n", page.Url, err)
		log.Printf("Error[%d]: %v", res.StatusCode, res)
		return []byte{}, err
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Printf("Unable to close body: %v", err)
		}
	}()

	bodyContent, err := io.ReadAll(res.Body)
	if err != nil {
		log.Printf("Failed to read a body for (url: \"%s\"): %v\n", page.Url, err)
		return []byte{}, err
	}

	return bodyContent, err
}

func parseUsingXPathQuery(content []byte, xpath string) ([]string, error) {
	log.Printf("Parse using the the XPath: %s", xpath)
	root, err := htmlquery.Parse(bytes.NewReader(content))
	if err != nil {
		return []string{}, err
	}
	nodes, err := htmlquery.QueryAll(root, xpath)
	if err != nil {
		return []string{}, err
	}

	var result []string

	for _, node := range nodes {
		html := htmlquery.OutputHTML(node, true)
		result = append(result, html)
	}

	return result, nil
}

func parseWebPageContent(page *config.Page, bodyContent []byte) (contentArray []string, err error) {
	if page.Query != "" {
		contentArray, err = parseUsingCssQuery(bodyContent, page.Query)
	} else {
		contentArray, err = parseUsingXPathQuery(bodyContent, page.XPath)
	}
	return
}

func parseUsingCssQuery(bodyContent []byte, query string) ([]string, error) {
	log.Printf("Parse using the the CSS Query: %s", query)
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(bodyContent))
	if err != nil {
		return []string{}, err
	}

	var content []string
	doc.Find(query).Each(func(idx int, selection *goquery.Selection) {
		htmlContent, err := selection.Html()
		if err != nil {
			log.Printf("Text extraction failed %v\n", err)
			return
		}
		content = append(content, htmlContent)
	})

	return content, nil
}

func (r *pageResolverContent) applyFilters(content string) string {
	if strings.TrimSpace(content) == "" {
		return ""
	}

	var err error

	for _, newFilter := range r.filters {
		filter := newFilter(&r.page)

		if !filter.IsEnabled() {
			continue
		}

		log.Printf("Appling filter \"%s\": %s", filter.Name(), content)
		content, err = filter.Filter(content)

		if err != nil {
			log.Printf("Unable to apply filter: %v", err)
		}

		if content == "" {
			return ""
		}
	}

	return strings.TrimSpace(content)
}

type urlOnlyResolver struct {
	page config.Page
}

func (u *urlOnlyResolver) Resolve(_ context.Context) RunResult {
	return RunResult{
		Page:    u.page,
		Content: fmt.Sprintf("Url for %s menu: %s", u.page.Name, u.page.Url),
		Status:  RunSuccess,
		Kind:    "url",
	}
}

func makeErrorResult(page config.Page, err error) RunResult {
	return RunResult{
		Page:    page,
		Content: fmt.Sprintf("Error: %v\n", err),
		Status:  RunError,
		Kind:    "error",
	}
}

func transformEncoding(content []byte) []byte {
	bytesReader := bytes.NewReader(content)

	e, name, _, err := DetermineEncodingFromReader(bytes.NewReader(content))
	if err != nil {
		log.Printf("Unable to determine the encoding: %v", err)
		return content
	}

	log.Printf("Found encoding: %s", name)

	reader := transform.NewReader(bytesReader, e.NewDecoder())
	result, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Printf("Unable to read from reader: %v", err)
	}

	return result
}

func DetermineEncodingFromReader(r io.Reader) (e encoding.Encoding, name string, certain bool, err error) {
	b, err := bufio.NewReader(r).Peek(1024)
	if err != nil {
		return
	}

	e, name, certain = charset.DetermineEncoding(b, "")
	return
}

type imageResolver struct {
	page   config.Page
	client http.Client
}

// Resolve implements PageResolver
func (r *imageResolver) Resolve(ctx context.Context) RunResult {
	bodyContent, err := getContentForWebPage(&r.page)
	if err != nil {
		return makeErrorResult(r.page, err)
	}

	contentArray, err := parseWebPageContent(&r.page, bodyContent)
	if err != nil {
		log.Printf("Parsing failed for (url: \"%s\"): %v\n", r.page.Url, err)
		return makeErrorResult(r.page, err)
	}

	content := strings.Join(contentArray, "")

	return RunResult{
		Page:    r.page,
		Content: content,
		Status:  RunSuccess,
		Kind:    "img",
	}
}
