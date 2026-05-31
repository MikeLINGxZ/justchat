package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/charset"
	"golang.org/x/net/publicsuffix"
	"golang.org/x/text/transform"
	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

const (
	fetchTimeout = 60 * time.Second

	// maxHTMLBodySize caps the HTTP response body read to prevent hanging on
	// huge pages (matches Claude Code's 10MB limit).
	maxHTMLBodySize = 10 * 1024 * 1024

	// maxMarkdownLength is the character limit after HTML→Markdown conversion,
	// keeping token usage reasonable (matches Claude Code's 100k cap).
	maxMarkdownLength = 50_000

	// maxTextLength is the character limit for plain-text extraction.
	maxTextLength = 20_000
)

// type aliases for testability
type htmlNode = html.Node

const elementNode = html.ElementNode

var fetchHTTPClient = newFetchHTTPClient()

func newFetchHTTPClient() *http.Client {
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	return &http.Client{
		Timeout: fetchTimeout,
		Jar:     jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}
}

type webFetchInput struct {
	URL    string `json:"url"    jsonschema:"description=The URL to fetch (must start with http:// or https://),required"`
	Format string `json:"format" jsonschema:"description=Output format: text (default) or markdown"`
}

type webFetchOutput struct {
	URL     string `json:"url"`
	Content string `json:"content"`
	Error   string `json:"error,omitempty"`
}

// fetchContent fetches the URL and returns the decoded body and Content-Type.
// It caps the body to maxHTMLBodySize and converts non-UTF-8 charsets to UTF-8.
func fetchContent(ctx context.Context, client *http.Client, targetURL string) (body string, contentType string, err error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return "", "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/markdown, text/html, application/xhtml+xml, application/xml;q=0.9, */*;q=0.8")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	resp, err := client.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("fetch %s: %w", targetURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("status %d from %s", resp.StatusCode, targetURL)
	}

	ct := resp.Header.Get("Content-Type")
	rawBytes, err := io.ReadAll(io.LimitReader(resp.Body, maxHTMLBodySize))
	if err != nil {
		return "", "", fmt.Errorf("read body: %w", err)
	}

	// Decode charset for HTML responses (handles GBK, GB2312, etc.)
	if strings.Contains(ct, "text/html") {
		enc, _, _ := charset.DetermineEncoding(rawBytes, ct)
		decoded, _, decErr := transform.Bytes(enc.NewDecoder(), rawBytes)
		if decErr == nil {
			return string(decoded), ct, nil
		}
	}

	return string(rawBytes), ct, nil
}

// --- HTML parsing helpers ---

func parseHTML(htmlStr string) (*html.Node, error) {
	return html.Parse(strings.NewReader(htmlStr))
}

func walkHTML(n *html.Node, fn func(*html.Node) bool) {
	if !fn(n) {
		return
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walkHTML(c, fn)
	}
}

func extractText(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var sb strings.Builder
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		sb.WriteString(extractText(c))
	}
	return sb.String()
}

// isSkippedElement returns true for elements whose subtrees should be ignored.
func isSkippedElement(tag string) bool {
	switch tag {
	case "script", "style", "nav", "header", "footer", "aside", "noscript":
		return true
	}
	return false
}

// extractMainText extracts visible text from an HTML page as plain text.
// It includes headings, paragraphs, list items, and table cells.
// Falls back to all body text if no structured elements are found.
func extractMainText(htmlStr string) string {
	doc, err := parseHTML(htmlStr)
	if err != nil {
		return ""
	}

	var parts []string
	walkHTML(doc, func(n *html.Node) bool {
		if n.Type != html.ElementNode {
			return true
		}
		if isSkippedElement(n.Data) {
			return false
		}
		switch n.Data {
		case "p", "li", "td", "th", "dt", "dd", "h1", "h2", "h3", "h4", "h5", "h6":
			text := strings.TrimSpace(extractText(n))
			if text != "" {
				parts = append(parts, text)
			}
			return false
		}
		return true
	})

	content := strings.Join(parts, "\n")

	// Fallback: if no structured elements found, return all body text
	if content == "" {
		content = strings.TrimSpace(extractText(doc))
	}

	if len(content) > maxTextLength {
		content = content[:maxTextLength] + "..."
	}
	return content
}

// --- Markdown conversion ---

// toMarkdown converts an HTML string to Markdown.
// Headings, paragraphs, lists, and inline links are preserved.
// Script, style, nav, header, footer, aside subtrees are skipped.
func toMarkdown(htmlStr string) string {
	doc, err := parseHTML(htmlStr)
	if err != nil {
		return ""
	}
	var sb strings.Builder
	renderMarkdownNode(doc, &sb)
	result := strings.TrimSpace(sb.String())
	if len(result) > maxMarkdownLength {
		result = result[:maxMarkdownLength] + "..."
	}
	return result
}

func renderMarkdownNode(n *htmlNode, sb *strings.Builder) {
	if n.Type == elementNode {
		if isSkippedElement(n.Data) {
			return
		}
		switch n.Data {
		case "h1":
			if text := strings.TrimSpace(extractText(n)); text != "" {
				sb.WriteString("# " + text + "\n\n")
			}
			return
		case "h2":
			if text := strings.TrimSpace(extractText(n)); text != "" {
				sb.WriteString("## " + text + "\n\n")
			}
			return
		case "h3":
			if text := strings.TrimSpace(extractText(n)); text != "" {
				sb.WriteString("### " + text + "\n\n")
			}
			return
		case "h4", "h5", "h6":
			if text := strings.TrimSpace(extractText(n)); text != "" {
				sb.WriteString("#### " + text + "\n\n")
			}
			return
		case "p":
			var inline strings.Builder
			renderInlineNode(n, &inline)
			if text := strings.TrimSpace(inline.String()); text != "" {
				sb.WriteString(text + "\n\n")
			}
			return
		case "li":
			var inline strings.Builder
			renderInlineNode(n, &inline)
			if text := strings.TrimSpace(inline.String()); text != "" {
				sb.WriteString("- " + text + "\n")
			}
			return
		case "ul", "ol":
			sb.WriteString("\n")
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				renderMarkdownNode(c, sb)
			}
			sb.WriteString("\n")
			return
		case "a":
			var inline strings.Builder
			renderInlineNode(n, &inline)
			if text := strings.TrimSpace(inline.String()); text != "" {
				sb.WriteString(text + "\n\n")
			}
			return
		case "br":
			sb.WriteString("\n")
			return
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		renderMarkdownNode(c, sb)
	}
}

// renderInlineNode renders a node's content inline, converting <a> to markdown links.
func renderInlineNode(n *htmlNode, sb *strings.Builder) {
	if n.Type == html.TextNode {
		sb.WriteString(n.Data)
		return
	}
	if n.Type == elementNode {
		if isSkippedElement(n.Data) {
			return
		}
		if n.Data == "a" {
			var href string
			for _, attr := range n.Attr {
				if attr.Key == "href" {
					href = attr.Val
					break
				}
			}
			text := strings.TrimSpace(extractText(n))
			if href != "" && text != "" {
				sb.WriteString("[" + text + "](" + href + ")")
			} else {
				sb.WriteString(text)
			}
			return
		}
		if n.Data == "br" {
			sb.WriteString("\n")
			return
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		renderInlineNode(c, sb)
	}
}

// --- Main function ---

func webFetchFunc(ctx context.Context, input webFetchInput) (webFetchOutput, error) {
	if !strings.HasPrefix(input.URL, "http://") && !strings.HasPrefix(input.URL, "https://") {
		return webFetchOutput{
			URL:   input.URL,
			Error: "URL must start with http:// or https://",
		}, nil
	}

	body, contentType, err := fetchContent(ctx, fetchHTTPClient, input.URL)
	if err != nil {
		return webFetchOutput{URL: input.URL, Error: err.Error()}, nil
	}

	// Non-HTML responses (JSON, plain text, markdown) are returned directly.
	if !strings.Contains(contentType, "text/html") {
		if len(body) > maxMarkdownLength {
			body = body[:maxMarkdownLength] + "..."
		}
		return webFetchOutput{URL: input.URL, Content: body}, nil
	}

	var content string
	if strings.ToLower(input.Format) == "markdown" {
		content = toMarkdown(body)
	} else {
		content = extractMainText(body)
	}

	return webFetchOutput{URL: input.URL, Content: content}, nil
}

func NewWebFetchTool() *function.FunctionTool[webFetchInput, webFetchOutput] {
	return function.NewFunctionTool(
		webFetchFunc,
		function.WithName("web_fetch"),
		function.WithDescription("Fetch a URL and return its content. HTML is converted to readable text (default) or markdown. JSON and plain text are returned as-is. The model decides which URL to visit."),
	)
}

func WebFetchMeta() ToolMeta {
	return ToolMeta{
		Name:            "web_fetch",
		Description:     "Fetch any URL and read its content",
		Category:        CategoryBuiltin,
		RequiresConfirm: false,
		FormatPurpose: func(args json.RawMessage) string {
			var input webFetchInput
			_ = json.Unmarshal(args, &input)
			return fmt.Sprintf("Fetch: %s", input.URL)
		},
	}
}
