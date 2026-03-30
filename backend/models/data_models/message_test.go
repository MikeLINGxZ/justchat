package data_models

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/cloudwego/eino/schema"
)

func TestToSchemaMessage_ImageFileBecomesImagePart(t *testing.T) {
	dir := t.TempDir()
	imagePath := filepath.Join(dir, "sample.png")
	if err := os.WriteFile(imagePath, []byte("fake-image-bytes"), 0o644); err != nil {
		t.Fatalf("write image fixture: %v", err)
	}

	msg := Message{
		Role:    schema.User,
		Content: "这张图里有什么",
		UserMessageExtra: &UserMessageExtra{
			Files: []File{{
				Name:     "sample.png",
				Path:     imagePath,
				MineType: "image/png",
			}},
		},
	}

	got, err := msg.ToSchemaMessage()
	if err != nil {
		t.Fatalf("ToSchemaMessage() error = %v", err)
	}
	if got.Content != "" {
		t.Fatalf("content = %q, want empty for multi-content user message", got.Content)
	}
	if len(got.UserInputMultiContent) != 2 {
		t.Fatalf("parts len = %d, want 2", len(got.UserInputMultiContent))
	}
	if got.UserInputMultiContent[0].Type != schema.ChatMessagePartTypeText {
		t.Fatalf("first part type = %q, want text", got.UserInputMultiContent[0].Type)
	}
	if got.UserInputMultiContent[1].Type != schema.ChatMessagePartTypeImageURL {
		t.Fatalf("second part type = %q, want image_url", got.UserInputMultiContent[1].Type)
	}
	if got.UserInputMultiContent[1].Image == nil {
		t.Fatalf("image part is nil")
	}
}

func TestToSchemaMessage_FileBecomesFilePart(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "sample.pdf")
	if err := os.WriteFile(filePath, []byte("%PDF-1.4 fake"), 0o644); err != nil {
		t.Fatalf("write file fixture: %v", err)
	}

	msg := Message{
		Role:    schema.User,
		Content: "这个文件主要讲什么",
		UserMessageExtra: &UserMessageExtra{
			Files: []File{{
				Name:     "sample.pdf",
				Path:     filePath,
				MineType: "application/pdf",
			}},
		},
	}

	got, err := msg.ToSchemaMessage()
	if err != nil {
		t.Fatalf("ToSchemaMessage() error = %v", err)
	}
	if len(got.UserInputMultiContent) != 2 {
		t.Fatalf("parts len = %d, want 2", len(got.UserInputMultiContent))
	}
	if got.UserInputMultiContent[1].Type != schema.ChatMessagePartTypeFileURL {
		t.Fatalf("second part type = %q, want file_url", got.UserInputMultiContent[1].Type)
	}
	if got.UserInputMultiContent[1].File == nil {
		t.Fatalf("file part is nil")
	}
	if got.UserInputMultiContent[1].File.Name != "sample.pdf" {
		t.Fatalf("file name = %q, want sample.pdf", got.UserInputMultiContent[1].File.Name)
	}
}

func TestToSchemaMessage_PreservesLeadingTextWithAttachment(t *testing.T) {
	dir := t.TempDir()
	filePath := filepath.Join(dir, "sample.json")
	if err := os.WriteFile(filePath, []byte(`{"hello":"world"}`), 0o644); err != nil {
		t.Fatalf("write file fixture: %v", err)
	}

	msg := Message{
		Role:    schema.User,
		Content: "帮我看看这个文件",
		UserMessageExtra: &UserMessageExtra{
			Files: []File{{
				Name:     "sample.json",
				Path:     filePath,
				MineType: "application/json",
			}},
		},
	}

	got, err := msg.ToSchemaMessage()
	if err != nil {
		t.Fatalf("ToSchemaMessage() error = %v", err)
	}
	if len(got.UserInputMultiContent) == 0 {
		t.Fatalf("parts len = 0, want > 0")
	}
	if got.UserInputMultiContent[0].Type != schema.ChatMessagePartTypeText {
		t.Fatalf("first part type = %q, want text", got.UserInputMultiContent[0].Type)
	}
	if got.UserInputMultiContent[0].Text != "帮我看看这个文件" {
		t.Fatalf("first part text = %q, want original content", got.UserInputMultiContent[0].Text)
	}
}

func TestAssistantMessageExtraJSONRoundTripPreservesPrefaceFields(t *testing.T) {
	extra := AssistantMessageExtra{
		PrefaceContent:          "先给你一个结论",
		PrefaceReasoningContent: "我先快速分析一下",
		CurrentStage:            "任务交付",
	}

	raw, err := json.Marshal(extra)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	var decoded AssistantMessageExtra
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if decoded.PrefaceContent != extra.PrefaceContent {
		t.Fatalf("preface_content = %q, want %q", decoded.PrefaceContent, extra.PrefaceContent)
	}
	if decoded.PrefaceReasoningContent != extra.PrefaceReasoningContent {
		t.Fatalf("preface_reasoning_content = %q, want %q", decoded.PrefaceReasoningContent, extra.PrefaceReasoningContent)
	}
}
