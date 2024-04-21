package resolvers

import (
	"bytes"
	"context"

	"github.com/PuerkitoBio/goquery"
	"github.com/antchfx/htmlquery"
	"github.com/pestanko/miniscrape/internal/models"
	"github.com/rs/zerolog"
	"golang.org/x/net/html"
)

// HTMLPageNode represents a node in the HTML page
type HTMLPageNode struct {
	Content string
	Attrs   []html.Attribute
}

// ParseWebPageContent parses the web page content
func ParseWebPageContent(
	ctx context.Context,
	page *models.Page,
	bodyContent []byte,
) (contentArray []HTMLPageNode, err error) {
	if page.Query != "" {
		contentArray, err = parseUsingCSSQuery(ctx, bodyContent, page.Query)
	} else {
		contentArray, err = parseUsingXPathQuery(ctx, bodyContent, page.XPath)
	}
	return
}

func parseUsingXPathQuery(ctx context.Context, content []byte, xpath string) ([]HTMLPageNode, error) {
	zerolog.Ctx(ctx).Trace().
		Str("xpath", xpath).
		Msg("Parse using the the XPath")

	root, err := htmlquery.Parse(bytes.NewReader(content))
	if err != nil {
		return []HTMLPageNode{}, err
	}
	nodes, err := htmlquery.QueryAll(root, xpath)
	if err != nil {
		return []HTMLPageNode{}, err
	}

	var result []HTMLPageNode

	for _, node := range nodes {
		if node == nil {
			continue
		}
		htmlContent := htmlquery.OutputHTML(node, true)
		result = append(result, HTMLPageNode{
			Content: htmlContent,
			Attrs:   node.Attr,
		})
	}

	return result, nil
}

func parseUsingCSSQuery(ctx context.Context, bodyContent []byte, query string) ([]HTMLPageNode, error) {
	ll := zerolog.Ctx(ctx).With().Str("css_query", query).Logger()
	ll.Trace().Msg("Parse using the the CSS query")
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(bodyContent))
	if err != nil {
		return []HTMLPageNode{}, err
	}

	var content []HTMLPageNode
	doc.Find(query).Each(func(idx int, selection *goquery.Selection) {
		htmlContent, err := selection.Html()
		if err != nil {
			ll.Warn().
				Err(err).
				Msg("Text extraction failed")
			return
		}

		attrs := getAttributesFromSelection(selection)

		content = append(content, HTMLPageNode{
			Content: htmlContent,
			Attrs:   attrs,
		})
	})

	if len(content) == 0 {
		ll.Warn().Msg("No content found")
	}

	return content, nil
}

func getAttributesFromSelection(selection *goquery.Selection) []html.Attribute {
	if selection == nil || len(selection.Nodes) == 0 {
		return []html.Attribute{}
	}

	return selection.Nodes[0].Attr
}

// getAttrValue returns the value of the attribute
func getAttrValue(attrs []html.Attribute, name string) string {
	for _, attr := range attrs {
		if attr.Key == name {
			return attr.Val
		}
	}
	return ""
}

// concatContent concats the content
func concatContent(contentArray []HTMLPageNode) string {
	var content string
	for _, node := range contentArray {
		content += node.Content + "\n"
	}
	return content
}
