package agent_dto

import (
	"encoding/json"
	"testing"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/provider"
)

func TestSendMessageInputJSONTags(t *testing.T) {
	payload, err := json.Marshal(SendMessageInput{
		SessionID:        7,
		Content:          "hello",
		BaseURL:          "https://api.example.com/v1",
		ApiKey:           "secret",
		ModelName:        "gpt-test",
		ProviderType:     provider.OpenAiCompatibility,
		EnabledUserTools: []string{"web_fetch", "code_exec"},
	})
	if err != nil {
		t.Fatalf("marshal send message input: %v", err)
	}

	got := string(payload)
	expectedFragments := []string{
		`"session_id":7`,
		`"content":"hello"`,
		`"base_url":"https://api.example.com/v1"`,
		`"api_key":"secret"`,
		`"model_name":"gpt-test"`,
		`"provider_type":"openai_compatibility"`,
		`"enabled_user_tools":["web_fetch","code_exec"]`,
	}

	for _, fragment := range expectedFragments {
		if !containsJSONFragment(got, fragment) {
			t.Fatalf("expected JSON payload to contain %s, got %s", fragment, got)
		}
	}
}

func TestListSessionsOutputJSONTags(t *testing.T) {
	payload, err := json.Marshal(ListSessionsOutput{
		Sessions: []SessionItem{
			{
				ID:      1,
				Title:   "Tea Chat",
				Starred: true,
				Status:  "idle",
				Created: "2026-05-15T00:00:00Z",
				Updated: "2026-05-15T00:00:01Z",
			},
		},
		NextCursor: 1,
		HasMore:    true,
	})
	if err != nil {
		t.Fatalf("marshal list sessions output: %v", err)
	}

	got := string(payload)
	expectedFragments := []string{
		`"sessions":[`,
		`"id":1`,
		`"title":"Tea Chat"`,
		`"starred":true`,
		`"status":"idle"`,
		`"created":"2026-05-15T00:00:00Z"`,
		`"updated":"2026-05-15T00:00:01Z"`,
		`"next_cursor":1`,
		`"has_more":true`,
	}

	for _, fragment := range expectedFragments {
		if !containsJSONFragment(got, fragment) {
			t.Fatalf("expected JSON payload to contain %s, got %s", fragment, got)
		}
	}
}

func containsJSONFragment(payload string, fragment string) bool {
	return len(payload) >= len(fragment) && func() bool {
		for i := 0; i+len(fragment) <= len(payload); i++ {
			if payload[i:i+len(fragment)] == fragment {
				return true
			}
		}
		return false
	}()
}
