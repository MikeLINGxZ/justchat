package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"golang.org/x/net/html"
	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

const WebSearchToolName = "web_search"

var (
	webSearchEndpoint   = "https://www.bing.com/search"
	webSearchHTTPClient = newFetchHTTPClient()
)

type webSearchInput struct {
	Query string `json:"query" jsonschema:"description=Search query,required"`
	Limit int    `json:"limit" jsonschema:"description=Maximum results to return, default 5, max 10"`
}

type webSearchResult struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet"`
}

type webSearchOutput struct {
	Query   string            `json:"query"`
	Results []webSearchResult `json:"results"`
}

// BuildWebSearchTool returns registry metadata for the web search tool.
func BuildWebSearchTool() ToolMeta {
	return ToolMeta{
		Name:            WebSearchToolName,
		Description:     "Search the web for public information and return concise result titles, URLs, and snippets.",
		Category:        CategoryBuiltin,
		RequiresConfirm: false,
		FormatPurpose: func(args json.RawMessage) string {
			var input webSearchInput
			_ = json.Unmarshal(args, &input)
			return "Search web: " + input.Query
		},
	}
}

// webSearchFunc searches the configured search endpoint and parses visible results.
func webSearchFunc(ctx context.Context, input webSearchInput) (webSearchOutput, error) {
	query := strings.TrimSpace(input.Query)
	if query == "" {
		return webSearchOutput{}, errors.New("query is required")
	}
	limit := input.Limit
	if limit <= 0 {
		limit = 5
	}
	if limit > 10 {
		limit = 10
	}
	searchURL, err := buildWebSearchURL(webSearchEndpoint, query, limit)
	if err != nil {
		return webSearchOutput{}, err
	}
	body, _, err := fetchContent(ctx, webSearchHTTPClient, searchURL)
	if err != nil {
		return webSearchOutput{}, err
	}
	results := parseWebSearchResults(body, limit)
	return webSearchOutput{Query: query, Results: results}, nil
}

// NewWebSearchTool creates the function tool used for public web searches.
func NewWebSearchTool() *function.FunctionTool[webSearchInput, webSearchOutput] {
	meta := BuildWebSearchTool()
	return function.NewFunctionTool(
		webSearchFunc,
		function.WithName(WebSearchToolName),
		function.WithDescription(meta.Description),
	)
}

// buildWebSearchURL appends query parameters to the configured search endpoint.
func buildWebSearchURL(endpoint string, query string, limit int) (string, error) {
	parsed, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("parse search endpoint: %w", err)
	}
	q := parsed.Query()
	q.Set("q", query)
	q.Set("count", fmt.Sprintf("%d", limit))
	parsed.RawQuery = q.Encode()
	return parsed.String(), nil
}

// parseWebSearchResults extracts Bing-style result cards from an HTML page.
func parseWebSearchResults(rawHTML string, limit int) []webSearchResult {
	doc, err := html.Parse(strings.NewReader(rawHTML))
	if err != nil {
		return nil
	}
	results := make([]webSearchResult, 0, limit)
	walkHTML(doc, func(n *html.Node) bool {
		if len(results) >= limit {
			return false
		}
		if n.Type == html.ElementNode && n.Data == "li" && hasClass(n, "b_algo") {
			if result, ok := parseSearchResultNode(n); ok {
				results = append(results, result)
			}
			return false
		}
		return true
	})
	return results
}

// parseSearchResultNode extracts title, URL, and snippet from one result node.
func parseSearchResultNode(n *html.Node) (webSearchResult, bool) {
	var result webSearchResult
	walkHTML(n, func(current *html.Node) bool {
		if current.Type == html.ElementNode && current.Data == "a" && result.URL == "" {
			result.Title = strings.TrimSpace(extractText(current))
			result.URL = attrValue(current, "href")
			return false
		}
		return true
	})
	walkHTML(n, func(current *html.Node) bool {
		if current.Type == html.ElementNode && current.Data == "p" && result.Snippet == "" {
			result.Snippet = strings.TrimSpace(extractText(current))
			return false
		}
		return true
	})
	return result, result.Title != "" && result.URL != ""
}

// hasClass reports whether an HTML node has a class token.
func hasClass(n *html.Node, className string) bool {
	classes := strings.Fields(attrValue(n, "class"))
	for _, current := range classes {
		if current == className {
			return true
		}
	}
	return false
}

// attrValue returns the value of a named HTML attribute.
func attrValue(n *html.Node, name string) string {
	for _, attr := range n.Attr {
		if attr.Key == name {
			return attr.Val
		}
	}
	return ""
}
