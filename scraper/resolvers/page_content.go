package resolvers

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os/exec"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/antchfx/htmlquery"
	"github.com/pestanko/miniscrape/scraper/config"
	"github.com/pestanko/miniscrape/scraper/filters"
	"github.com/rs/zerolog/log"
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
	page    config.Page
	client  http.Client
	filters []func(*config.Page) filters.PageFilter
}

func (r *pageContentResolver) Resolve(_ context.Context) config.RunResult {
	bodyContent, err := getContentForWebPage(&r.page)
	if err != nil {
		return makeErrorResult(r.page, err)
	}

	contentArray, err := parseWebPageContent(&r.page, bodyContent)
	if err != nil {
		log.Error().
			Err(err).
			Str("url", r.page.Url).
			Str("page", r.page.Namespace()).
			Msg("Content parsing failed")
		return makeErrorResult(r.page, err)
	}

	content := strings.Join(contentArray, "\n")
	content = r.applyFilters(content)

	var status = config.RunSuccess
	if content == "" {
		log.Debug().
			Str("page", r.page.Namespace()).
			Msg("Content resolved but the content is empty")
		status = config.RunEmpty
	} else {
		log.Debug().
			Str("page", r.page.Namespace()).
			Msg("Content resolved")
	}

	return config.RunResult{
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
	log.Debug().
		Str("cmdName", page.Command.Content.Name).
		Strs("cmdArgs", page.Command.Content.Args).
		Msg("Resolve using command")
	var outb, errb bytes.Buffer
	cmd := exec.Command(page.Command.Content.Name, page.Command.Content.Args...)
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	err := cmd.Run()
	if err != nil {
		log.Error().
			Err(err).
			Str("page", page.Namespace()).
			Str("cmdName", page.Command.Content.Name).
			Strs("cmdArgs", page.Command.Content.Args).
			Msg("Command error")
		log.Trace().
			Str("page", page.Namespace()).
			Str("cmdName", page.Command.Content.Name).
			Strs("cmdArgs", page.Command.Content.Args).
			Str("stderr", errb.String()).
			Msg("Command error trace")
	}

	return outb.Bytes(), err
}

func getContentByRequest(page *config.Page) ([]byte, error) {
	req, err := http.NewRequest("GET", page.Url, nil)
	if err != nil {
		log.Error().
			Err(err).
			Str("pageUrl", page.Url).
			Str("page", page.Namespace()).
			Msg("Request initialization failed")
		return []byte{}, err
	}

	randomUserAgent := userAgents[rand.Intn(len(userAgents))]
	req.Header.Add("User-Agent", randomUserAgent)

	res, err := httpClient.Do(req)
	if err != nil {
		log.Error().
			Err(err).
			Str("pageUrl", page.Url).
			Str("page", page.Namespace()).
			Int("status", res.StatusCode).
			Msg("Request failed")
		log.Trace().
			Stack().
			Err(err).
			Str("page", page.Namespace()).
			Str("content", fmt.Sprintf("%v", res)).
			Msg("Error reponse content")
		return []byte{}, err
	}

	defer func() {
		if err := res.Body.Close(); err != nil {
			log.Error().
				Err(err).
				Str("pageUrl", page.Url).
				Str("page", page.Namespace()).
				Int("status", res.StatusCode).
				Msg("Unable to close body")
		}
	}()

	bodyContent, err := io.ReadAll(res.Body)
	if err != nil {
		log.Error().
			Err(err).
			Str("pageUrl", page.Url).
			Str("page", page.Namespace()).
			Int("status", res.StatusCode).
			Msg("Failed to read a body")

		return []byte{}, err
	}

	return bodyContent, err
}

func parseUsingXPathQuery(content []byte, xpath string) ([]string, error) {
	log.Trace().Str("xpath", xpath).Msg("Parse using the the XPath")
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
	log.Trace().Str("query", query).Msg("Parse using the the CSS query")
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(bodyContent))
	if err != nil {
		return []string{}, err
	}

	var content []string
	doc.Find(query).Each(func(idx int, selection *goquery.Selection) {
		htmlContent, err := selection.Html()
		if err != nil {
			log.Warn().
				Err(err).
				Msg("Text extraction failed")
			return
		}
		content = append(content, htmlContent)
	})

	return content, nil
}

func (r *pageContentResolver) applyFilters(content string) string {
	if strings.TrimSpace(content) == "" {
		return ""
	}

	var err error

	for _, newFilter := range r.filters {
		filter := newFilter(&r.page)

		if !filter.IsEnabled() {
			continue
		}

		log.Trace().
			Err(err).
			Str("filter", filter.Name()).
			Str("content", content).
			Msg("Appling filter")

		content, err = filter.Filter(content)

		if err != nil {
			log.Warn().
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

func transformEncoding(content []byte) []byte {
	bytesReader := bytes.NewReader(content)

	e, name, _, err := DetermineEncodingFromReader(bytes.NewReader(content))
	if err != nil {
		log.Warn().
			Err(err).
			Msg("Unable to determine the encoding")
		return content
	}

	log.Trace().
		Str("encoding", name).
		Msg("Found encoding")

	reader := transform.NewReader(bytesReader, e.NewDecoder())
	result, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Warn().
			Err(err).
			Msg("Unable to read from reader")
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
