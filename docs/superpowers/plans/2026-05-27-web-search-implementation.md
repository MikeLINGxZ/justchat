# Web Search Tool Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Replace the `web_search` stub with a real Bing HTML scraping implementation that returns title, URL, snippet, and optionally full page content.

**Architecture:** Scrape `https://www.bing.com/search?q=...` using Go's `net/http`, parse the HTML DOM with `golang.org/x/net/html` to extract result entries from `<li class="b_algo">` nodes. When `depth=full`, concurrently fetch each result page and extract main text paragraphs.

**Tech Stack:** Go stdlib (`net/http`, `net/url`, `strings`, `sync`), `golang.org/x/net/html` (already in go.mod)

---

## File Map

- Modify: `backend/pkg/agent/tools/web_search.go` — full rewrite, all logic here
- Create: `backend/pkg/agent/tools/web_search_test.go` — unit tests using `net/http/httptest`

---

### Task 1: Pure helper functions with tests

**Files:**
- Modify: `backend/pkg/agent/tools/web_search.go`
- Create: `backend/pkg/agent/tools/web_search_test.go`

- [ ] **Step 1: Write the failing tests**

Create `backend/pkg/agent/tools/web_search_test.go`:

```go
package tools

import (
	"strings"
	"testing"
)

func TestNormalizeLimit(t *testing.T) {
	cases := []struct{ in, want int }{
		{0, 5},
		{-1, 5},
		{3, 3},
		{10, 10},
		{11, 10},
		{100, 10},
	}
	for _, c := range cases {
		got := normalizeLimit(c.in)
		if got != c.want {
			t.Errorf("normalizeLimit(%d) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestBuildBingURL(t *testing.T) {
	u := buildBingURL("golang concurrency", 5)
	if !strings.Contains(u, "bing.com/search") {
		t.Errorf("expected bing.com/search in URL, got %q", u)
	}
	if !strings.Contains(u, "golang+concurrency") && !strings.Contains(u, "golang%20concurrency") {
		t.Errorf("expected query in URL, got %q", u)
	}
}

func TestExtractText(t *testing.T) {
	htmlStr := `<p>Hello <b>World</b></p>`
	doc, _ := parseHTML(htmlStr)
	// find the <p> node
	var p *htmlNode
	walkHTML(doc, func(n *htmlNode) bool {
		if n.Type == elementNode && n.Data == "p" {
			p = n
			return false
		}
		return true
	})
	if p == nil {
		t.Fatal("no <p> found")
	}
	got := extractText(p)
	if got != "Hello World" {
		t.Errorf("extractText = %q, want %q", got, "Hello World")
	}
}

func TestExtractMainText_SkipsScriptStyle(t *testing.T) {
	htmlStr := `<html><body>
		<script>var x = 1;</script>
		<style>.foo{color:red}</style>
		<p>This is real content that should appear.</p>
		<p>More content here for the test.</p>
	</body></html>`
	got := extractMainText(htmlStr)
	if strings.Contains(got, "var x") || strings.Contains(got, ".foo") {
		t.Errorf("extractMainText should skip script/style, got %q", got)
	}
	if !strings.Contains(got, "real content") {
		t.Errorf("extractMainText should include paragraph text, got %q", got)
	}
}

func TestExtractMainText_TruncatesAt2000(t *testing.T) {
	var sb strings.Builder
	sb.WriteString("<html><body>")
	for i := 0; i < 50; i++ {
		sb.WriteString("<p>This paragraph has enough text to be included in the extraction result for testing purposes indeed.</p>")
	}
	sb.WriteString("</body></html>")
	got := extractMainText(sb.String())
	if len(got) > 2003 { // 2000 + "..."
		t.Errorf("extractMainText should truncate to 2000 chars, got len=%d", len(got))
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop && go test ./backend/pkg/agent/tools/ -run "TestNormalizeLimit|TestBuildBingURL|TestExtractText|TestExtractMainText" -v 2>&1 | tail -20
```

Expected: compile errors (`normalizeLimit`, `buildBingURL`, `parseHTML`, etc. undefined)

- [ ] **Step 3: Write the helper implementations**

Replace the entire content of `backend/pkg/agent/tools/web_search.go` with:

```go
package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/html"
	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

const (
	bingBaseURL            = "https://www.bing.com/search"
	defaultLimit           = 5
	maxLimit               = 10
	fullDepthPageTimeout   = 10 * time.Second
	maxFullDepthConcurrent = 3
)

// type aliases for testability
type htmlNode = html.Node

const elementNode = html.ElementNode

var bingHTTPClient = &http.Client{Timeout: 30 * time.Second}

type webSearchInput struct {
	Query string `json:"query" jsonschema:"description=Search query string,required"`
	Limit int    `json:"limit" jsonschema:"description=Max number of results (default: 5)"`
	Depth string `json:"depth" jsonschema:"description=Result depth: basic (title+url+snippet) or full (also fetches page content). Default: basic"`
}

type webSearchResult struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet"`
	Content string `json:"content,omitempty"`
}

type webSearchOutput struct {
	Results []webSearchResult `json:"results"`
	Query   string            `json:"query"`
}

func normalizeLimit(limit int) int {
	if limit <= 0 {
		return defaultLimit
	}
	if limit > maxLimit {
		return maxLimit
	}
	return limit
}

func buildBingURL(query string, limit int) string {
	params := url.Values{}
	params.Set("q", query)
	params.Set("count", fmt.Sprintf("%d", limit))
	return bingBaseURL + "?" + params.Encode()
}

func parseHTML(htmlStr string) (*html.Node, error) {
	return html.Parse(strings.NewReader(htmlStr))
}

// walkHTML calls fn on each node; if fn returns false, stops descending that branch.
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

func extractMainText(htmlStr string) string {
	doc, err := parseHTML(htmlStr)
	if err != nil {
		return ""
	}
	var paragraphs []string
	walkHTML(doc, func(n *html.Node) bool {
		if n.Type == html.ElementNode {
			switch n.Data {
			case "script", "style", "nav", "header", "footer":
				return false
			case "p":
				text := strings.TrimSpace(extractText(n))
				if len(text) > 30 {
					paragraphs = append(paragraphs, text)
				}
			}
		}
		return true
	})
	content := strings.Join(paragraphs, "\n\n")
	if len(content) > 2000 {
		content = content[:2000] + "..."
	}
	return content
}

func fetchHTML(ctx context.Context, client *http.Client, targetURL string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, targetURL, nil)
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/124.0.0.0 Safari/537.36")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("fetch %s: %w", targetURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status %d from %s", resp.StatusCode, targetURL)
	}

	var sb strings.Builder
	buf := make([]byte, 32*1024)
	for {
		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			sb.Write(buf[:n])
		}
		if readErr != nil {
			break
		}
	}
	return sb.String(), nil
}

func parseBingResults(htmlStr string, limit int) []webSearchResult {
	doc, err := parseHTML(htmlStr)
	if err != nil {
		return nil
	}
	var results []webSearchResult
	walkHTML(doc, func(n *html.Node) bool {
		if len(results) >= limit {
			return false
		}
		if n.Type == html.ElementNode && n.Data == "li" {
			for _, attr := range n.Attr {
				if attr.Key == "class" && strings.Contains(attr.Val, "b_algo") {
					if r, ok := extractBingResult(n); ok {
						results = append(results, r)
					}
					return false
				}
			}
		}
		return true
	})
	return results
}

func extractBingResult(li *html.Node) (webSearchResult, bool) {
	var r webSearchResult
	walkHTML(li, func(n *html.Node) bool {
		if n.Type == html.ElementNode && n.Data == "h2" {
			walkHTML(n, func(a *html.Node) bool {
				if a.Type == html.ElementNode && a.Data == "a" {
					for _, attr := range a.Attr {
						if attr.Key == "href" {
							r.URL = attr.Val
						}
					}
					r.Title = strings.TrimSpace(extractText(a))
					return false
				}
				return true
			})
			return false
		}
		return true
	})
	if r.Snippet == "" {
		walkHTML(li, func(n *html.Node) bool {
			if r.Snippet != "" {
				return false
			}
			if n.Type == html.ElementNode && n.Data == "p" {
				text := strings.TrimSpace(extractText(n))
				if text != "" {
					r.Snippet = text
					return false
				}
			}
			return true
		})
	}
	if r.Title == "" || r.URL == "" {
		return webSearchResult{}, false
	}
	return r, true
}

func webSearchFunc(ctx context.Context, input webSearchInput) (webSearchOutput, error) {
	limit := normalizeLimit(input.Limit)
	targetURL := buildBingURL(input.Query, limit)

	htmlStr, err := fetchHTML(ctx, bingHTTPClient, targetURL)
	if err != nil {
		return webSearchOutput{}, fmt.Errorf("fetch bing results: %w", err)
	}

	results := parseBingResults(htmlStr, limit)

	if input.Depth == "full" && len(results) > 0 {
		sem := make(chan struct{}, maxFullDepthConcurrent)
		var mu sync.Mutex
		var wg sync.WaitGroup
		for i := range results {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()
				fetchCtx, cancel := context.WithTimeout(ctx, fullDepthPageTimeout)
				defer cancel()
				pageHTML, fetchErr := fetchHTML(fetchCtx, bingHTTPClient, results[i].URL)
				if fetchErr != nil {
					return
				}
				content := extractMainText(pageHTML)
				mu.Lock()
				results[i].Content = content
				mu.Unlock()
			}(i)
		}
		wg.Wait()
	}

	return webSearchOutput{Results: results, Query: input.Query}, nil
}

func NewWebSearchTool() *function.FunctionTool[webSearchInput, webSearchOutput] {
	return function.NewFunctionTool(
		webSearchFunc,
		function.WithName("web_search"),
		function.WithDescription("Search the web using Bing and return results with title, URL, and snippet. Use depth=full to also fetch page content."),
	)
}

func WebSearchMeta() ToolMeta {
	return ToolMeta{
		Name:            "web_search",
		Description:     "Search the web for information",
		Category:        CategoryUser,
		RequiresConfirm: false,
		FormatPurpose: func(args json.RawMessage) string {
			var input webSearchInput
			_ = json.Unmarshal(args, &input)
			return fmt.Sprintf("Search the web: %s", input.Query)
		},
	}
}
```

- [ ] **Step 4: Run the tests to verify they pass**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop && go test ./backend/pkg/agent/tools/ -run "TestNormalizeLimit|TestBuildBingURL|TestExtractText|TestExtractMainText" -v 2>&1 | tail -20
```

Expected: all 4 test functions PASS

- [ ] **Step 5: Commit**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop && git add backend/pkg/agent/tools/web_search.go backend/pkg/agent/tools/web_search_test.go && git commit -m "feat(web_search): implement helper functions and HTML extraction"
```

---

### Task 2: parseBingResults with mock HTML tests

**Files:**
- Modify: `backend/pkg/agent/tools/web_search_test.go`

- [x] **Step 1: Write the failing test**

Append to `backend/pkg/agent/tools/web_search_test.go`:

```go
func TestParseBingResults_ExtractsResults(t *testing.T) {
	// Minimal Bing-like HTML with two b_algo entries
	htmlStr := `<html><body><ol id="b_results">
		<li class="b_algo">
			<h2><a href="https://example.com/page1">First Result Title</a></h2>
			<div class="b_caption"><p>First result snippet text here.</p></div>
		</li>
		<li class="b_algo">
			<h2><a href="https://example.com/page2">Second Result Title</a></h2>
			<div class="b_caption"><p>Second result snippet text here.</p></div>
		</li>
		<li class="b_algo">
			<h2><a href="https://example.com/page3">Third Result Title</a></h2>
			<p>Third result snippet text here.</p>
		</li>
	</ol></body></html>`

	results := parseBingResults(htmlStr, 5)
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}
	if results[0].Title != "First Result Title" {
		t.Errorf("results[0].Title = %q, want %q", results[0].Title, "First Result Title")
	}
	if results[0].URL != "https://example.com/page1" {
		t.Errorf("results[0].URL = %q, want %q", results[0].URL, "https://example.com/page1")
	}
	if results[0].Snippet == "" {
		t.Error("results[0].Snippet should not be empty")
	}
}

func TestParseBingResults_RespectsLimit(t *testing.T) {
	htmlStr := `<html><body><ol id="b_results">
		<li class="b_algo"><h2><a href="https://a.com">A</a></h2><p>snippet a</p></li>
		<li class="b_algo"><h2><a href="https://b.com">B</a></h2><p>snippet b</p></li>
		<li class="b_algo"><h2><a href="https://c.com">C</a></h2><p>snippet c</p></li>
	</ol></body></html>`

	results := parseBingResults(htmlStr, 2)
	if len(results) != 2 {
		t.Fatalf("expected 2 results with limit=2, got %d", len(results))
	}
}

func TestParseBingResults_SkipsEntriesWithoutURL(t *testing.T) {
	htmlStr := `<html><body><ol id="b_results">
		<li class="b_algo"><h2><a>No href here</a></h2><p>snippet</p></li>
		<li class="b_algo"><h2><a href="https://valid.com">Valid</a></h2><p>snippet</p></li>
	</ol></body></html>`

	results := parseBingResults(htmlStr, 5)
	if len(results) != 1 {
		t.Fatalf("expected 1 result (skip entry without URL), got %d", len(results))
	}
	if results[0].URL != "https://valid.com" {
		t.Errorf("results[0].URL = %q, want https://valid.com", results[0].URL)
	}
}
```

- [ ] **Step 2: Run tests to verify they pass** (they should pass since parseBingResults is already implemented)

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop && go test ./backend/pkg/agent/tools/ -run "TestParseBingResults" -v 2>&1 | tail -20
```

Expected: all 3 TestParseBingResults tests PASS

- [ ] **Step 3: Commit**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop && git add backend/pkg/agent/tools/web_search_test.go && git commit -m "test(web_search): add parseBingResults unit tests"
```

---

### Task 3: fetchHTML with mock HTTP server tests

**Files:**
- Modify: `backend/pkg/agent/tools/web_search_test.go`

- [ ] **Step 1: Update imports and write the failing tests**

First, update the import block at the top of `backend/pkg/agent/tools/web_search_test.go` to:

```go
import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)
```

Then append to `backend/pkg/agent/tools/web_search_test.go`:

```go
func TestFetchHTML_ReturnsBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("<html><body>hello</body></html>"))
	}))
	defer srv.Close()

	client := srv.Client()
	got, err := fetchHTML(context.Background(), client, srv.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(got, "hello") {
		t.Errorf("expected body to contain 'hello', got %q", got)
	}
}

func TestFetchHTML_ErrorOnNon200(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	client := srv.Client()
	_, err := fetchHTML(context.Background(), client, srv.URL)
	if err == nil {
		t.Fatal("expected error for non-200 status")
	}
}

func TestFetchHTML_SetsUserAgent(t *testing.T) {
	var gotUA string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotUA = r.Header.Get("User-Agent")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	defer srv.Close()

	client := srv.Client()
	_, _ = fetchHTML(context.Background(), client, srv.URL)
	if !strings.Contains(gotUA, "Mozilla") {
		t.Errorf("expected browser User-Agent, got %q", gotUA)
	}
}
```

Note: The test file already imports `"strings"` and `"testing"`. Add `"context"`, `"net/http"`, `"net/http/httptest"` to the import block at the top of the test file.

- [x] **Step 2: Run tests to verify they pass**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop && go test ./backend/pkg/agent/tools/ -run "TestFetchHTML" -v 2>&1 | tail -20
```

Expected: all 3 TestFetchHTML tests PASS

- [ ] **Step 3: Commit**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop && git add backend/pkg/agent/tools/web_search_test.go && git commit -m "test(web_search): add fetchHTML unit tests with httptest"
```

---

### Task 4: webSearchFunc integration test (basic depth)

**Files:**
- Modify: `backend/pkg/agent/tools/web_search_test.go`

- [ ] **Step 1: Write the failing test**

Append to `backend/pkg/agent/tools/web_search_test.go`:

```go
func TestWebSearchFunc_BasicDepth(t *testing.T) {
	bingHTML := `<html><body><ol id="b_results">
		<li class="b_algo">
			<h2><a href="https://example.com/go">Go Language</a></h2>
			<p>Go is an open source programming language.</p>
		</li>
		<li class="b_algo">
			<h2><a href="https://example.com/tour">A Tour of Go</a></h2>
			<p>An interactive introduction to Go.</p>
		</li>
	</ol></body></html>`

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(bingHTML))
	}))
	defer srv.Close()

	// Override bingHTTPClient and bingBaseURL for this test
	origClient := bingHTTPClient
	origURL := bingBaseURL
	bingHTTPClient = srv.Client()
	bingBaseURL = srv.URL + "/search"
	defer func() {
		bingHTTPClient = origClient
		bingBaseURL = origURL
	}()

	out, err := webSearchFunc(context.Background(), webSearchInput{Query: "golang", Limit: 5})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Query != "golang" {
		t.Errorf("out.Query = %q, want golang", out.Query)
	}
	if len(out.Results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(out.Results))
	}
	if out.Results[0].Title != "Go Language" {
		t.Errorf("results[0].Title = %q", out.Results[0].Title)
	}
	if out.Results[0].Content != "" {
		t.Error("basic depth should not populate Content")
	}
}
```

Note: This test overrides `bingHTTPClient` and `bingBaseURL` (which are package-level vars). Change `bingBaseURL` from a `const` to a `var` in `web_search.go`:

In `web_search.go`, change:
```go
const (
    bingBaseURL            = "https://www.bing.com/search"
    ...
)
```
to:
```go
var bingBaseURL = "https://www.bing.com/search"

const (
    defaultLimit           = 5
    maxLimit               = 10
    fullDepthPageTimeout   = 10 * time.Second
    maxFullDepthConcurrent = 3
)
```

- [ ] **Step 2: Run tests to verify they pass**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop && go test ./backend/pkg/agent/tools/ -run "TestWebSearchFunc_BasicDepth" -v 2>&1 | tail -20
```

Expected: PASS

- [ ] **Step 3: Commit**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop && git add backend/pkg/agent/tools/web_search.go backend/pkg/agent/tools/web_search_test.go && git commit -m "test(web_search): add webSearchFunc basic-depth integration test"
```

---

### Task 5: webSearchFunc full-depth test

**Files:**
- Modify: `backend/pkg/agent/tools/web_search_test.go`

- [x] **Step 1: Write the failing test**

Append to `backend/pkg/agent/tools/web_search_test.go`:

```go
func TestWebSearchFunc_FullDepth_FetchesContent(t *testing.T) {
	pageContent := `<html><body>
		<p>This is the main content of the result page with enough text to be included.</p>
		<p>A second paragraph with additional information about the topic in detail.</p>
	</body></html>`

	// Server handles both the Bing search page and the result pages
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if strings.Contains(r.URL.Path, "/search") {
			bingHTML := `<html><body><ol id="b_results">
				<li class="b_algo">
					<h2><a href="` + "http://" + r.Host + `/page1` + `">Page One</a></h2>
					<p>Page one snippet here for testing purposes.</p>
				</li>
			</ol></body></html>`
			_, _ = w.Write([]byte(bingHTML))
		} else {
			_, _ = w.Write([]byte(pageContent))
		}
	}))
	defer srv.Close()

	origClient := bingHTTPClient
	origURL := bingBaseURL
	bingHTTPClient = srv.Client()
	bingBaseURL = srv.URL + "/search"
	defer func() {
		bingHTTPClient = origClient
		bingBaseURL = origURL
	}()

	out, err := webSearchFunc(context.Background(), webSearchInput{Query: "test", Limit: 1, Depth: "full"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out.Results) == 0 {
		t.Fatal("expected at least 1 result")
	}
	if out.Results[0].Content == "" {
		t.Error("full depth should populate Content field")
	}
	if !strings.Contains(out.Results[0].Content, "main content") {
		t.Errorf("Content should include page text, got %q", out.Results[0].Content)
	}
}

func TestWebSearchFunc_FullDepth_ContinuesOnPageFetchError(t *testing.T) {
	// Bing returns one result pointing to a URL that 404s
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/search") {
			w.WriteHeader(http.StatusOK)
			bingHTML := `<html><body><ol id="b_results">
				<li class="b_algo">
					<h2><a href="` + "http://" + r.Host + `/will-404` + `">Some Page</a></h2>
					<p>Some snippet about this page here.</p>
				</li>
			</ol></body></html>`
			_, _ = w.Write([]byte(bingHTML))
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	origClient := bingHTTPClient
	origURL := bingBaseURL
	bingHTTPClient = srv.Client()
	bingBaseURL = srv.URL + "/search"
	defer func() {
		bingHTTPClient = origClient
		bingBaseURL = origURL
	}()

	out, err := webSearchFunc(context.Background(), webSearchInput{Query: "test", Limit: 1, Depth: "full"})
	if err != nil {
		t.Fatalf("unexpected error even when page fetch fails: %v", err)
	}
	if len(out.Results) == 0 {
		t.Fatal("should still return results even when page content fetch fails")
	}
	// Content is empty because page 404'd — that's fine
	if out.Results[0].Title != "Some Page" {
		t.Errorf("result title = %q", out.Results[0].Title)
	}
}
```

- [x] **Step 2: Run tests to verify they pass**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop && go test ./backend/pkg/agent/tools/ -run "TestWebSearchFunc_FullDepth" -v 2>&1 | tail -20
```

Expected: both PASS

- [ ] **Step 3: Commit**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop && git add backend/pkg/agent/tools/web_search_test.go && git commit -m "test(web_search): add full-depth and error-resilience tests"
```

---

### Task 6: Run full test suite and verify

**Files:** none

- [x] **Step 1: Run all web_search tests**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop && go test ./backend/pkg/agent/tools/ -v 2>&1 | grep -E "^(=== RUN|--- PASS|--- FAIL|FAIL|ok)"
```

Expected: all tests PASS, no FAIL lines

- [x] **Step 2: Run the broader backend tests to catch regressions**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop && go test ./backend/... 2>&1 | tail -20
```

Expected: all packages pass

- [x] **Step 3: Verify the tool builds cleanly**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop && go build ./backend/... 2>&1
```

Expected: no output (clean build)
