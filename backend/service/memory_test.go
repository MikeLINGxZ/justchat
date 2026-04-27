package service

import (
	"strings"
	"testing"

	"github.com/cloudwego/eino/schema"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/prompts"
)

func TestExtractMemoryEncodingInputBlocksAttachmentQuestion(t *testing.T) {
	messages := []schema.Message{{
		Role: schema.User,
		UserInputMultiContent: []schema.MessageInputPart{
			{Type: schema.ChatMessagePartTypeText, Text: "这张图里有什么"},
			{
				Type: schema.ChatMessagePartTypeImageURL,
				Image: &schema.MessageInputImage{
					MessagePartCommon: schema.MessagePartCommon{MIMEType: "image/png"},
					Detail:            schema.ImageURLDetailHigh,
				},
			},
		},
	}}

	input := extractMemoryEncodingInput(messages)
	if input.UserMessage != "这张图里有什么" {
		t.Fatalf("UserMessage = %q, want %q", input.UserMessage, "这张图里有什么")
	}
	if !input.HasNonTextPart {
		t.Fatal("HasNonTextPart = false, want true")
	}
	if !input.ExternalQuestion {
		t.Fatal("ExternalQuestion = false, want true")
	}
	if shouldEncodeMemoryInput(input) {
		t.Fatal("shouldEncodeMemoryInput() = true, want false")
	}
}

func TestExtractMemoryEncodingInputKeepsUserDisclosureWithAttachment(t *testing.T) {
	messages := []schema.Message{{
		Role: schema.User,
		UserInputMultiContent: []schema.MessageInputPart{
			{Type: schema.ChatMessagePartTypeText, Text: "帮我写个请假消息，我明天去医院复查"},
			{
				Type: schema.ChatMessagePartTypeFileURL,
				File: &schema.MessageInputFile{
					MessagePartCommon: schema.MessagePartCommon{MIMEType: "application/pdf"},
					Name:              "note.pdf",
				},
			},
		},
	}}

	input := extractMemoryEncodingInput(messages)
	if input.UserMessage != "帮我写个请假消息，我明天去医院复查" {
		t.Fatalf("UserMessage = %q", input.UserMessage)
	}
	if !input.HasNonTextPart {
		t.Fatal("HasNonTextPart = false, want true")
	}
	if input.ExternalQuestion {
		t.Fatal("ExternalQuestion = true, want false")
	}
	if !shouldEncodeMemoryInput(input) {
		t.Fatal("shouldEncodeMemoryInput() = false, want true")
	}
}

func TestExtractMemoryEncodingInputPlainUserFactStillEncodes(t *testing.T) {
	messages := []schema.Message{{
		Role:    schema.User,
		Content: "我对花生过敏",
	}}

	input := extractMemoryEncodingInput(messages)
	if input.UserMessage != "我对花生过敏" {
		t.Fatalf("UserMessage = %q, want %q", input.UserMessage, "我对花生过敏")
	}
	if input.HasNonTextPart {
		t.Fatal("HasNonTextPart = true, want false")
	}
	if input.ExternalQuestion {
		t.Fatal("ExternalQuestion = true, want false")
	}
	if !shouldEncodeMemoryInput(input) {
		t.Fatal("shouldEncodeMemoryInput() = false, want true")
	}
}

func TestDefaultMemoryPromptMatchesCurrentToolContract(t *testing.T) {
	content, ok := prompts.DefaultPromptContent("system.memory.md")
	if !ok {
		t.Fatal("DefaultPromptContent(system.memory.md) = not found, want found")
	}

	required := []string{
		"Hermes 风格的有界长期记忆",
		"target=user：用户画像",
		"target=memory：助手/环境笔记",
		"没有 read 动作",
		"图片内容本身",
		"文件、PDF、网页、代码、日志、表格本身的内容",
		"临时任务请求",
	}
	for _, item := range required {
		if !strings.Contains(content, item) {
			t.Fatalf("memory prompt missing required content %q", item)
		}
	}

	forbidden := []string{
		"memory_type`（可选：skill | event | plan）",
		"可选：`time_range_start`、`time_range_end`",
		"可更新：`title`、`content`、`time_range_start`",
		"write_memory",
		"read_memory",
	}
	for _, item := range forbidden {
		if strings.Contains(content, item) {
			t.Fatalf("memory prompt contains legacy token %q", item)
		}
	}
}
