# Web Search Tool Implementation Design

**Date:** 2026-05-27  
**Status:** Approved

## Background

`backend/pkg/agent/tools/web_search.go` is currently a stub that returns a placeholder message. This design covers replacing it with a real Bing HTML scraping implementation — no API key required.

## Input / Output

### Input

```go
type webSearchInput struct {
    Query string `json:"query" jsonschema:"description=Search query string,required"`
    Limit int    `json:"limit" jsonschema:"description=Max number of results (default: 5)"`
    Depth string `json:"depth" jsonschema:"description=Result depth: basic (title+url+snippet) or full (also fetch page content). Default: basic"`
}
```

### Output

```go
type webSearchResult struct {
    Title   string `json:"title"`
    URL     string `json:"url"`
    Snippet string `json:"snippet"`
    Content string `json:"content,omitempty"` // only populated when depth=full
}

type webSearchOutput struct {
    Results []webSearchResult `json:"results"`
    Query   string            `json:"query"`
}
```

## Architecture

```
webSearchFunc(ctx, input)
  ├── buildBingURL(query, limit) → URL
  ├── fetchHTML(ctx, url) → html string
  │     └── http.Client + browser User-Agent + 30s timeout
  ├── parseBingResults(html, limit) → []webSearchResult
  │     └── golang.org/x/net/html parses DOM
  │         extracts from <li class="b_algo">:
  │         - <h2><a> → title + url
  │         - .b_caption p / <p> → snippet
  └── [depth=full] concurrent fetchPageContent(ctx, url) → content
          ├── max 3 concurrent goroutines (semaphore)
          ├── fetch target page HTML
          └── extractMainText(html) → extract <p> text, truncate to 2000 chars
```

## Key Implementation Details

- User-Agent is set to a common browser UA string to avoid Bing blocking
- `fetchHTML` and `fetchPageContent` share a package-level `http.Client` (connection pooling)
- `parseBingResults` only includes valid entries (has both title and URL)
- `depth=full`: per-page fetch timeout is 10s; failure leaves `Content` empty without aborting other results
- `Limit` ≤ 0 → default 5; > 10 → capped at 10

## Error Handling

| Scenario | Behavior |
|---|---|
| Bing request fails (network / non-200) | Return `error` to surface to AI |
| HTML parse yields 0 results | Return empty `Results` list (no error) |
| `depth=full` single page fetch fails | `Content` left empty, result still included |
| `Limit` out of range | Normalize silently (0→5, >10→10) |
| Context cancelled / timeout | Propagate via `ctx` |

## Dependencies

- `golang.org/x/net/html` — already in `go.mod` as indirect; will be used directly
- No new dependencies needed

## Files Changed

- `backend/pkg/agent/tools/web_search.go` — full rewrite (all logic in one file)
- `go.mod` / `go.sum` — no changes needed
