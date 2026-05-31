package tools

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// --- toMarkdown tests ---

func TestToMarkdown_Headings(t *testing.T) {
	html := `<html><body>
		<h1>Title One</h1>
		<h2>Section Two</h2>
		<h3>Subsection Three</h3>
	</body></html>`
	got := toMarkdown(html)
	if !strings.Contains(got, "# Title One") {
		t.Errorf("expected '# Title One', got: %q", got)
	}
	if !strings.Contains(got, "## Section Two") {
		t.Errorf("expected '## Section Two', got: %q", got)
	}
	if !strings.Contains(got, "### Subsection Three") {
		t.Errorf("expected '### Subsection Three', got: %q", got)
	}
}

func TestToMarkdown_Paragraphs(t *testing.T) {
	html := `<html><body>
		<p>First paragraph with enough content to be included.</p>
		<p>Second paragraph also has enough content to be included.</p>
	</body></html>`
	got := toMarkdown(html)
	if !strings.Contains(got, "First paragraph") {
		t.Errorf("expected first paragraph, got: %q", got)
	}
	if !strings.Contains(got, "Second paragraph") {
		t.Errorf("expected second paragraph, got: %q", got)
	}
}

func TestToMarkdown_Links(t *testing.T) {
	html := `<html><body>
		<p>Visit <a href="https://example.com">Example Site</a> for more.</p>
	</body></html>`
	got := toMarkdown(html)
	if !strings.Contains(got, "[Example Site](https://example.com)") {
		t.Errorf("expected markdown link, got: %q", got)
	}
}

func TestToMarkdown_ListItems(t *testing.T) {
	html := `<html><body>
		<ul>
			<li>First item in the list</li>
			<li>Second item in the list</li>
		</ul>
	</body></html>`
	got := toMarkdown(html)
	if !strings.Contains(got, "- First item") {
		t.Errorf("expected '- First item', got: %q", got)
	}
	if !strings.Contains(got, "- Second item") {
		t.Errorf("expected '- Second item', got: %q", got)
	}
}

func TestToMarkdown_SkipsScriptAndStyle(t *testing.T) {
	html := `<html><body>
		<script>var x = secret;</script>
		<style>.foo { color: red; }</style>
		<p>Real content that should appear in the output.</p>
	</body></html>`
	got := toMarkdown(html)
	if strings.Contains(got, "secret") || strings.Contains(got, ".foo") {
		t.Errorf("toMarkdown should skip script/style, got: %q", got)
	}
	if !strings.Contains(got, "Real content") {
		t.Errorf("expected real content, got: %q", got)
	}
}

func TestToMarkdown_TruncatesAt50000(t *testing.T) {
	var sb strings.Builder
	sb.WriteString("<html><body>")
	for range 700 {
		sb.WriteString("<p>This paragraph has enough text to be included and will push us well past the fifty thousand character limit.</p>")
	}
	sb.WriteString("</body></html>")
	got := toMarkdown(sb.String())
	if len(got) > 50003 {
		t.Errorf("toMarkdown should truncate at 50000 chars, got len=%d", len(got))
	}
}

// --- extractMainText tests ---

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
		t.Errorf("missing paragraph text, got %q", got)
	}
}

func TestExtractMainText_IncludesListItems(t *testing.T) {
	htmlStr := `<html><body>
		<ul>
			<li>Stock price: 60.20 HKD</li>
			<li>Change: -0.80 (-1.31%)</li>
		</ul>
	</body></html>`
	got := extractMainText(htmlStr)
	if !strings.Contains(got, "Stock price") {
		t.Errorf("extractMainText should include <li> content, got: %q", got)
	}
}

func TestExtractMainText_IncludesHeadings(t *testing.T) {
	htmlStr := `<html><body>
		<h1>Main Title</h1>
		<h2>Section Title</h2>
	</body></html>`
	got := extractMainText(htmlStr)
	if !strings.Contains(got, "Main Title") {
		t.Errorf("extractMainText should include heading content, got: %q", got)
	}
}

func TestExtractMainText_FallsBackToBodyText(t *testing.T) {
	// No <p>, <li>, or <h*> — should still return body text
	htmlStr := `<html><body>
		<div><span>Stock: 60.20</span></div>
	</body></html>`
	got := extractMainText(htmlStr)
	if !strings.Contains(got, "60.20") {
		t.Errorf("extractMainText should fall back to body text, got: %q", got)
	}
}

func TestExtractMainText_TruncatesAt20000(t *testing.T) {
	var sb strings.Builder
	sb.WriteString("<html><body>")
	for range 300 {
		sb.WriteString("<p>This paragraph has enough text to be included in the extraction result and push past the limit.</p>")
	}
	sb.WriteString("</body></html>")
	got := extractMainText(sb.String())
	if len(got) > 20003 {
		t.Errorf("should truncate at 20000, got len=%d", len(got))
	}
}

// --- webFetchFunc tests ---

func TestWebFetchFunc_TextFormat(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<html><body>
			<p>This is enough real content to be extracted by the text formatter.</p>
		</body></html>`))
	}))
	defer srv.Close()

	origClient := fetchHTTPClient
	fetchHTTPClient = srv.Client()
	defer func() { fetchHTTPClient = origClient }()

	out, err := webFetchFunc(context.Background(), webFetchInput{URL: srv.URL, Format: "text"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Error != "" {
		t.Fatalf("unexpected error field: %q", out.Error)
	}
	if !strings.Contains(out.Content, "real content") {
		t.Errorf("expected content, got: %q", out.Content)
	}
}

func TestWebFetchFunc_MarkdownFormat(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<html><body>
			<h1>Main Title</h1>
			<p>Some content under the heading.</p>
			<a href="https://example.com">A Link</a>
		</body></html>`))
	}))
	defer srv.Close()

	origClient := fetchHTTPClient
	fetchHTTPClient = srv.Client()
	defer func() { fetchHTTPClient = origClient }()

	out, err := webFetchFunc(context.Background(), webFetchInput{URL: srv.URL, Format: "markdown"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.Content, "# Main Title") {
		t.Errorf("expected markdown heading, got: %q", out.Content)
	}
}

func TestWebFetchFunc_DefaultFormatIsText(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<html><body>
			<p>Default format content that should be returned as plain text.</p>
		</body></html>`))
	}))
	defer srv.Close()

	origClient := fetchHTTPClient
	fetchHTTPClient = srv.Client()
	defer func() { fetchHTTPClient = origClient }()

	out, err := webFetchFunc(context.Background(), webFetchInput{URL: srv.URL})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Error != "" {
		t.Fatalf("unexpected error field: %q", out.Error)
	}
	if out.Content == "" {
		t.Error("expected non-empty content with default format")
	}
}

func TestWebFetchFunc_JSONContentReturnedDirectly(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"price":60.20,"currency":"HKD","symbol":"2015.HK"}`))
	}))
	defer srv.Close()

	origClient := fetchHTTPClient
	fetchHTTPClient = srv.Client()
	defer func() { fetchHTTPClient = origClient }()

	out, err := webFetchFunc(context.Background(), webFetchInput{URL: srv.URL})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Error != "" {
		t.Fatalf("unexpected error field: %q", out.Error)
	}
	if !strings.Contains(out.Content, "60.20") {
		t.Errorf("expected JSON content passed through, got: %q", out.Content)
	}
}

func TestWebFetchFunc_SetsAcceptHeader(t *testing.T) {
	var gotAccept string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAccept = r.Header.Get("Accept")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<html><body><p>content</p></body></html>`))
	}))
	defer srv.Close()

	origClient := fetchHTTPClient
	fetchHTTPClient = srv.Client()
	defer func() { fetchHTTPClient = origClient }()

	_, _ = webFetchFunc(context.Background(), webFetchInput{URL: srv.URL})
	if !strings.Contains(gotAccept, "text/html") {
		t.Errorf("expected Accept header to include text/html, got: %q", gotAccept)
	}
}

func TestWebFetchFunc_InvalidURLReturnsError(t *testing.T) {
	out, err := webFetchFunc(context.Background(), webFetchInput{URL: "not-a-url"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Error == "" {
		t.Error("expected error for invalid URL")
	}
}

func TestWebFetchFunc_Non200ReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	origClient := fetchHTTPClient
	fetchHTTPClient = srv.Client()
	defer func() { fetchHTTPClient = origClient }()

	out, err := webFetchFunc(context.Background(), webFetchInput{URL: srv.URL})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.Error == "" {
		t.Error("expected error for non-200 status")
	}
}

func TestWebFetchFunc_SetsURL(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<html><body><p>Content here is long enough to include.</p></body></html>`))
	}))
	defer srv.Close()

	origClient := fetchHTTPClient
	fetchHTTPClient = srv.Client()
	defer func() { fetchHTTPClient = origClient }()

	out, err := webFetchFunc(context.Background(), webFetchInput{URL: srv.URL})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out.URL != srv.URL {
		t.Errorf("expected URL %q, got %q", srv.URL, out.URL)
	}
}
