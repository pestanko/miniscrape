package resolvers

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"github.com/pestanko/miniscrape/internal/models"
	"github.com/pestanko/miniscrape/internal/scraper/filters"
	"io"
	"math/rand"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/antchfx/htmlquery"
	"github.com/rs/zerolog"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/transform"
)

var userAgents = []string{
	"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/42.0.2311.135 Safari/537.36 Edge/12.246",
	`Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.27 Safari/537.36`,
	`Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/47.0.2526.111 Safari/537.36`,
}

var httpClient = http.Client{
	Timeout: 30 * time.Second,
}

type pageContentResolver struct {
	page    models.Page
	client  http.Client
	filters []func(*models.Page) filters.PageFilter
}

func (r *pageContentResolver) Resolve(ctx context.Context) models.RunResult {
	bodyContent, err := getContentForWebPage(ctx, &r.page)
	if err != nil {
		return makeErrorResult(r.page, err)
	}

	ll := zerolog.Ctx(ctx).With().
		Dict("page", zerolog.Dict().
			Str("namespace", r.page.Namespace()).
			Str("url", r.page.URL).
			Str("codename", r.page.CodeName)).
		Logger()

	ll.Trace().Bytes("body", bodyContent).Msg("page body")

	contentArray, err := parseWebPageContent(ctx, &r.page, bodyContent)
	if err != nil {
		ll.
			Err(err).
			Str("url", r.page.URL).
			Msg("Content parsing failed")
		return makeErrorResult(r.page, err)
	}

	content := strings.Join(contentArray, "\n")
	content = r.applyFilters(ctx, content)

	var status = models.RunSuccess
	if content == "" {
		ll.Warn().
			Msg("Content resolved but the content is empty")
		status = models.RunEmpty
	} else {
		ll.Debug().
			Msg("Content resolved")
	}

	return models.RunResult{
		Page:    r.page,
		Status:  status,
		Content: content,
		Kind:    "content",
	}
}

func getContentForWebPage(ctx context.Context, page *models.Page) (bodyContent []byte, err error) {
	if page.Command.Content.Name != "" {
		bodyContent, err = getContentByCommand(ctx, page)
	} else {
		bodyContent, err = getContentByRequest(ctx, page)
	}

	if err == nil {
		bodyContent = transformEncoding(ctx, bodyContent)
	}

	return
}

func getContentByCommand(ctx context.Context, page *models.Page) ([]byte, error) {
	// Use command
	cmdContent := page.Command.Content

	ll := zerolog.Ctx(ctx).
		With().
		Str("page", page.Namespace()).
		Str("cmdName", cmdContent.Name).
		Strs("cmdArgs", cmdContent.Args).
		Logger()

	ll.Debug().Msg("Resolve using command")

	var outb, errb bytes.Buffer
	cmd := exec.Command(cmdContent.Name, cmdContent.Args...) // #nosec G204
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		ll.Error().Msg("Command error")
		ll.Trace().
			Err(err).
			Str("stderr", errb.String()).
			Msg("Command error trace")
	}

	return outb.Bytes(), err
}

func getContentByRequest(ctx context.Context, page *models.Page) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, page.URL, nil)
	ll := zerolog.Ctx(ctx).With().
		Str("pageUrl", page.URL).
		Str("page", page.Namespace()).
		Logger()

	if err != nil {
		ll.Err(err).
			Msg("Request initialization failed")
		return []byte{}, err
	}

	randomUserAgent := userAgents[rand.Intn(len(userAgents))] // #nosec G404
	req.Header.Add("User-Agent", randomUserAgent)

	res, err := httpClient.Do(req)
	if err != nil {
		ll.Error().
			Err(err).
			Int("status", res.StatusCode).
			Msg("Request failed")
		ll.Trace().
			Stack().
			Err(err).
			Str("content", fmt.Sprintf("%v", res)).
			Msg("Error response content")
		return []byte{}, err
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			ll.Error().
				Err(err).
				Int("status", res.StatusCode).
				Msg("Unable to close body")
		}
	}()

	bodyContent, err := io.ReadAll(res.Body)
	if err != nil {
		ll.Error().
			Err(err).
			Int("status", res.StatusCode).
			Msg("Failed to read a body")

		return []byte{}, err
	}

	return bodyContent, err
}

func parseUsingXPathQuery(ctx context.Context, content []byte, xpath string) ([]string, error) {
	zerolog.Ctx(ctx).Trace().
		Str("xpath", xpath).
		Msg("Parse using the the XPath")

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

func parseWebPageContent(
	ctx context.Context,
	page *models.Page,
	bodyContent []byte,
) (contentArray []string, err error) {
	if page.Query != "" {
		contentArray, err = parseUsingCSSQuery(ctx, bodyContent, page.Query)
	} else {
		contentArray, err = parseUsingXPathQuery(ctx, bodyContent, page.XPath)
	}
	return
}

func parseUsingCSSQuery(ctx context.Context, bodyContent []byte, query string) ([]string, error) {
	ll := zerolog.Ctx(ctx).With().Str("css_query", query).Logger()
	ll.Trace().Msg("Parse using the the CSS query")
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(bodyContent))
	if err != nil {
		return []string{}, err
	}

	var content []string
	doc.Find(query).Each(func(idx int, selection *goquery.Selection) {
		htmlContent, err := selection.Html()
		if err != nil {
			ll.Warn().
				Err(err).
				Msg("Text extraction failed")
			return
		}
		content = append(content, htmlContent)
	})

	return content, nil
}

func (r *pageContentResolver) applyFilters(ctx context.Context, content string) string {
	if strings.TrimSpace(content) == "" {
		return ""
	}

	var err error

	ll := zerolog.Ctx(ctx)

	for _, newFilter := range r.filters {
		filter := newFilter(&r.page)

		if !filter.IsEnabled() {
			continue
		}

		ll.Trace().
			Err(err).
			Str("filter", filter.Name()).
			Str("content", content).
			Msg("Appling filter")

		content, err = filter.Filter(content)

		if err != nil {
			ll.Warn().
				Err(err).
				Str("filter", filter.Name()).
				Msg("Unable to apply filter")
		}

		if content == "" {
			return ""
		}
	}

	return strings.TrimSpace(content)
}

func transformEncoding(ctx context.Context, content []byte) []byte {
	bytesReader := bytes.NewReader(content)

	ll := zerolog.Ctx(ctx)

	e, name, _, err := DetermineEncodingFromReader(bytes.NewReader(content))
	if err != nil {
		ll.Warn().
			Err(err).
			Msg("Unable to determine the encoding")
		return content
	}

	ll.Trace().
		Str("encoding", name).
		Msg("Found encoding")

	reader := transform.NewReader(bytesReader, e.NewDecoder())
	result, err := io.ReadAll(reader)
	if err != nil {
		ll.Warn().
			Err(err).
			Msg("Unable to read from reader")
	}

	return result
}

// DetermineEncodingFromReader based on the content
func DetermineEncodingFromReader(
	r io.Reader,
) (e encoding.Encoding, name string, certain bool, err error) {
	b, err := bufio.NewReader(r).Peek(1024)
	if err != nil {
		return
	}

	e, name, certain = charset.DetermineEncoding(b, "")
	return
}
