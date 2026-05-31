package agent

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	pkgProvider "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/provider"
)

// TestOneshotCompleteReturnsTextAndUsage verifies the helper decodes a non-streaming OpenAI-compatible response.
func TestOneshotCompleteReturnsTextAndUsage(t *testing.T) {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/chat/completions" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"id":"chatcmpl-test",
			"object":"chat.completion",
			"created":1730000000,
			"model":"qwen-test",
			"choices":[{"index":0,"message":{"role":"assistant","content":"manifest text"}}],
			"usage":{"prompt_tokens":12,"completion_tokens":7,"total_tokens":19}
		}`)
	}))
	defer server.Close()

	resp, err := OneshotComplete(context.Background(), OneshotRequest{
		BaseURL:      server.URL,
		APIKey:       "test-key",
		ModelName:    "qwen-test",
		ProviderType: pkgProvider.OpenAiCompatibility,
		System:       "system prompt",
		User:         "user prompt",
	})
	if err != nil {
		t.Fatalf("OneshotComplete returned error: %v", err)
	}
	if resp.Text != "manifest text" {
		t.Fatalf("unexpected text: %q", resp.Text)
	}
	if resp.PromptTokens != 12 || resp.CompletionTokens != 7 {
		t.Fatalf("unexpected usage: %+v", resp)
	}
}

// TestOneshotCompleteReturnsHTTPError verifies API-level failures become a Go error with useful details.
func TestOneshotCompleteReturnsHTTPError(t *testing.T) {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprint(w, `{"error":{"message":"invalid api key"}}`)
	}))
	defer server.Close()

	_, err := OneshotComplete(context.Background(), OneshotRequest{
		BaseURL:      server.URL,
		APIKey:       "bad-key",
		ModelName:    "qwen-test",
		ProviderType: pkgProvider.OpenAiCompatibility,
		System:       "system prompt",
		User:         "user prompt",
	})
	if err == nil {
		t.Fatal("expected OneshotComplete to return an error")
	}
	if !strings.Contains(err.Error(), "invalid api key") {
		t.Fatalf("expected error to contain API detail, got %v", err)
	}
}

// TestOneshotCompleteReturnsTimeout verifies the helper applies its timeout guard to slow providers.
func TestOneshotCompleteReturnsTimeout(t *testing.T) {
	t.Helper()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(150 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{
			"id":"chatcmpl-test",
			"object":"chat.completion",
			"created":1730000000,
			"model":"qwen-test",
			"choices":[{"index":0,"message":{"role":"assistant","content":"late"}}]
		}`)
	}))
	defer server.Close()

	_, err := OneshotComplete(context.Background(), OneshotRequest{
		BaseURL:      server.URL,
		APIKey:       "test-key",
		ModelName:    "qwen-test",
		ProviderType: pkgProvider.OpenAiCompatibility,
		System:       "system prompt",
		User:         "user prompt",
		Timeout:      50 * time.Millisecond,
	})
	if err == nil {
		t.Fatal("expected timeout error")
	}
	if !strings.Contains(err.Error(), context.DeadlineExceeded.Error()) {
		t.Fatalf("expected deadline exceeded, got %v", err)
	}
}
