package tools

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWebSearchFuncParsesHTMLResults(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("q") != "coffee report" {
			t.Fatalf("unexpected query: %s", r.URL.RawQuery)
		}
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`<html><body>
			<li class="b_algo">
				<h2><a href="https://example.com/a">Coffee Market</a></h2>
				<div class="b_caption"><p>Market summary and current trend.</p></div>
			</li>
			<li class="b_algo">
				<h2><a href="https://example.com/b">Coffee Data</a></h2>
				<p>Useful data for a research report.</p>
			</li>
		</body></html>`))
	}))
	defer srv.Close()

	origEndpoint := webSearchEndpoint
	origClient := webSearchHTTPClient
	webSearchEndpoint = srv.URL
	webSearchHTTPClient = srv.Client()
	defer func() {
		webSearchEndpoint = origEndpoint
		webSearchHTTPClient = origClient
	}()

	out, err := webSearchFunc(context.Background(), webSearchInput{Query: "coffee report", Limit: 2})
	if err != nil {
		t.Fatalf("web search: %v", err)
	}
	if len(out.Results) != 2 {
		t.Fatalf("expected 2 results, got %+v", out.Results)
	}
	if out.Results[0].Title != "Coffee Market" || out.Results[0].URL != "https://example.com/a" {
		t.Fatalf("unexpected first result: %+v", out.Results[0])
	}
	if !strings.Contains(out.Results[0].Snippet, "Market summary") {
		t.Fatalf("expected snippet, got %+v", out.Results[0])
	}
}

func TestWebSearchFuncRequiresQuery(t *testing.T) {
	_, err := webSearchFunc(context.Background(), webSearchInput{})
	if err == nil {
		t.Fatal("expected query to be required")
	}
}
