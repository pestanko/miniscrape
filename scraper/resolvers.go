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
	"jaytaylor.com/html2text"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246",
	`Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.27 Safari/537.36`,
	`Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36`,
}

var normPattern = regexp.MustCompile("\n\n")

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

type pageResolvedGet struct {
	page   config.Page
	client http.Client
}

func (r *pageResolvedGet) Resolve(_ context.Context) RunResult {
	var bodyContent []byte
	var err error
	if r.page.Command.Content.Name != "" {
		bodyContent, err = r.getContentByCommand()
	} else {
		bodyContent, err = r.getContentByRequest()
	}

	bodyContent = transformEncoding(bodyContent)

	if err != nil {
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
	}
}

func (r *pageResolvedGet) getContentByCommand() ([]byte, error) {
	// Use command
	log.Printf("Using command: '%s' with args %v", r.page.Command.Content.Name, r.page.Command.Content.Args)
	var outb, errb bytes.Buffer
	cmd := exec.Command(r.page.Command.Content.Name, r.page.Command.Content.Args...)
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		log.Printf("Command error[%v]: %s", err, errb.String())
	}

	return outb.Bytes(), err
}

func (r *pageResolvedGet) getContentByRequest() ([]byte, error) {
	req, err := http.NewRequest("GET", r.page.Url, nil)
	if err != nil {
		log.Printf("Request creation failed for (url: \"%s\"): %v\n", r.page.Url, err)
		return []byte{}, err
	}
	randomUserAgent := userAgents[rand.Intn(len(userAgents))]
	req.Header.Add("User-Agent", randomUserAgent)
	res, err := r.client.Do(req)
	if err != nil {
		log.Printf("Request failed for (url: \"%s\"): %v\n", r.page.Url, err)
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
		log.Printf("Failed to read a body for (url: \"%s\"): %v\n", r.page.Url, err)
		return []byte{}, err
	}

	return bodyContent, err
}

func (r *pageResolvedGet) parseUsingXPathQuery(content []byte) ([]string, error) {
	log.Printf("Parse using the the XPath: %s", r.page.XPath)
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
	log.Printf("Parse using the the CSS Query: %s", r.page.Query)
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
	log.Printf("Found content, converting: %s", htmlContent)

	if r.page.Filters.Html.Tables == "custom" {
		htmlContent = useCustomHTMLTablesConverter(htmlContent)
	}

	text, err := html2text.FromString(htmlContent, html2text.Options{
		PrettyTables: !r.page.Filters.Html.NoPrettyTables,
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

func (u *urlOnlyResolver) Resolve(_ context.Context) RunResult {
	return RunResult{
		Page:    u.page,
		Content: fmt.Sprintf("Url for %s menu: %s", u.page.Name, u.page.Url),
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
	return normPattern.ReplaceAllString(content, "\n")
}

func transformEncoding(content []byte) []byte {
	bytesReader := bytes.NewReader(content);

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


func useCustomHTMLTablesConverter(content string) string {
	if content == "" {
		return ""
	}

	content = strings.ReplaceAll(content, "</tr>", "<br/>")
	
	return strings.ReplaceAll(content, "</TR>", "<br/>")
}