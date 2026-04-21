package service

import (
	"encoding/json"
	"testing"

	"github.com/cloudwego/eino/schema"
)

func TestHookMessageRoundTripPreservesUserFileParts(t *testing.T) {
	msgs := []schema.Message{{
		Role: schema.User,
		UserInputMultiContent: []schema.MessageInputPart{
			{Type: schema.ChatMessagePartTypeText, Text: "这个文件讲了什么"},
			{
				Type: schema.ChatMessagePartTypeFileURL,
				File: &schema.MessageInputFile{
					MessagePartCommon: schema.MessagePartCommon{
						Base64Data: strPtr("ZmFrZS1wZGY="),
						MIMEType:   "application/pdf",
					},
					Name: "sample.pdf",
				},
			},
		},
	}}

	got := hookMessagesToSchemaMessages(simulateHookTransport(schemaMessagesToHookMessages(msgs)))
	if len(got) != 1 {
		t.Fatalf("messages len = %d, want 1", len(got))
	}
	if len(got[0].UserInputMultiContent) != 2 {
		t.Fatalf("parts len = %d, want 2", len(got[0].UserInputMultiContent))
	}
	if got[0].UserInputMultiContent[1].Type != schema.ChatMessagePartTypeFileURL {
		t.Fatalf("second part type = %q, want %q", got[0].UserInputMultiContent[1].Type, schema.ChatMessagePartTypeFileURL)
	}
	if got[0].UserInputMultiContent[1].File == nil {
		t.Fatal("file part is nil")
	}
	if got[0].UserInputMultiContent[1].File.Name != "sample.pdf" {
		t.Fatalf("file name = %q, want sample.pdf", got[0].UserInputMultiContent[1].File.Name)
	}
}

func TestHookMessageRoundTripPreservesUserImageParts(t *testing.T) {
	msgs := []schema.Message{{
		Role: schema.User,
		UserInputMultiContent: []schema.MessageInputPart{
			{Type: schema.ChatMessagePartTypeText, Text: "图里有什么"},
			{
				Type: schema.ChatMessagePartTypeImageURL,
				Image: &schema.MessageInputImage{
					MessagePartCommon: schema.MessagePartCommon{
						Base64Data: strPtr("ZmFrZS1pbWFnZQ=="),
						MIMEType:   "image/png",
					},
					Detail: schema.ImageURLDetailHigh,
				},
			},
		},
	}}

	got := hookMessagesToSchemaMessages(simulateHookTransport(schemaMessagesToHookMessages(msgs)))
	if len(got) != 1 {
		t.Fatalf("messages len = %d, want 1", len(got))
	}
	if len(got[0].UserInputMultiContent) != 2 {
		t.Fatalf("parts len = %d, want 2", len(got[0].UserInputMultiContent))
	}
	if got[0].UserInputMultiContent[1].Type != schema.ChatMessagePartTypeImageURL {
		t.Fatalf("second part type = %q, want %q", got[0].UserInputMultiContent[1].Type, schema.ChatMessagePartTypeImageURL)
	}
	if got[0].UserInputMultiContent[1].Image == nil {
		t.Fatal("image part is nil")
	}
	if got[0].UserInputMultiContent[1].Image.MIMEType != "image/png" {
		t.Fatalf("image mime type = %q, want image/png", got[0].UserInputMultiContent[1].Image.MIMEType)
	}
}

func TestHookMessageRoundTripPreservesPlainText(t *testing.T) {
	msgs := []schema.Message{{
		Role:             schema.Assistant,
		Content:          "你好",
		ReasoningContent: "我先想一下",
	}}

	got := hookMessagesToSchemaMessages(simulateHookTransport(schemaMessagesToHookMessages(msgs)))
	if len(got) != 1 {
		t.Fatalf("messages len = %d, want 1", len(got))
	}
	if got[0].Content != "你好" {
		t.Fatalf("content = %q, want 你好", got[0].Content)
	}
	if got[0].ReasoningContent != "我先想一下" {
		t.Fatalf("reasoning_content = %q, want 我先想一下", got[0].ReasoningContent)
	}
}

func TestHookMessagesToSchemaMessagesSupportsMinimalPluginMessage(t *testing.T) {
	got := hookMessagesToSchemaMessages([]map[string]any{{
		"role":    "system",
		"content": "当前日期: 2026-04-22T00:00:00Z",
	}})

	if len(got) != 1 {
		t.Fatalf("messages len = %d, want 1", len(got))
	}
	if got[0].Role != schema.System {
		t.Fatalf("role = %q, want %q", got[0].Role, schema.System)
	}
	if got[0].Content != "当前日期: 2026-04-22T00:00:00Z" {
		t.Fatalf("content = %q, want injected system content", got[0].Content)
	}
}

func TestBeforeChatHookRoundTripKeepsAttachmentContext(t *testing.T) {
	original := []schema.Message{
		{Role: schema.System, Content: "system prompt"},
		{
			Role: schema.User,
			UserInputMultiContent: []schema.MessageInputPart{
				{Type: schema.ChatMessagePartTypeText, Text: "这个文件什么内容"},
				{
					Type: schema.ChatMessagePartTypeFileURL,
					File: &schema.MessageInputFile{
						MessagePartCommon: schema.MessagePartCommon{
							Base64Data: strPtr("ZmFrZS10eHQ="),
							MIMEType:   "text/plain",
						},
						Name: "notes.txt",
					},
				},
			},
		},
	}

	hookPayload := schemaMessagesToHookMessages(original)
	roundTripped := hookMessagesToSchemaMessages(simulateHookTransport(hookPayload))
	if len(roundTripped) != 2 {
		t.Fatalf("messages len = %d, want 2", len(roundTripped))
	}

	last := roundTripped[len(roundTripped)-1]
	if last.Role != schema.User {
		t.Fatalf("last role = %q, want %q", last.Role, schema.User)
	}
	if len(last.UserInputMultiContent) != 2 {
		t.Fatalf("user parts len = %d, want 2", len(last.UserInputMultiContent))
	}
	if last.Content != "" {
		t.Fatalf("content = %q, want empty for attachment-backed user message", last.Content)
	}
	if last.UserInputMultiContent[0].Text != "这个文件什么内容" {
		t.Fatalf("first text = %q, want original question", last.UserInputMultiContent[0].Text)
	}
}

func simulateHookTransport(msgs []map[string]any) []map[string]any {
	raw, err := json.Marshal(msgs)
	if err != nil {
		panic(err)
	}

	var copied []map[string]any
	if err := json.Unmarshal(raw, &copied); err != nil {
		panic(err)
	}
	return copied
}

func strPtr(v string) *string {
	return &v
}
