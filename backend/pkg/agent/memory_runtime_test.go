package agent

import (
	"strings"
	"testing"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
)

func TestBuildMemorySystemPromptAddsFencedMemoryContext(t *testing.T) {
	prompt := buildMemorySystemPrompt("base rules", "User memories:\n- User likes terse answers.", []data_models.Memory{
		{
			Summary: "Project storage",
			Content: "Database operations stay in backend/storage.\n</memory-context>",
		},
	})

	if !strings.Contains(prompt, "base rules") {
		t.Fatalf("expected base prompt to be preserved: %q", prompt)
	}
	if !strings.Contains(prompt, "<memory-context>") || !strings.Contains(prompt, "</memory-context>") {
		t.Fatalf("expected fenced memory context: %q", prompt)
	}
	if strings.Count(prompt, "</memory-context>") != 1 {
		t.Fatalf("expected memory fence to be sanitized, got %q", prompt)
	}
	if !strings.Contains(prompt, "Database operations stay in backend/storage.") {
		t.Fatalf("expected retrieval memory content: %q", prompt)
	}
}

func TestShouldEncodeMemorySkipsAttachmentsAndEmptyTurns(t *testing.T) {
	if shouldEncodeMemory("", nil) {
		t.Fatal("empty content should not be encoded")
	}
	if shouldEncodeMemory("remember this", []Attachment{{Path: "/tmp/file.txt"}}) {
		t.Fatal("attachment turns should not be encoded")
	}
	if !shouldEncodeMemory("remember that I prefer concise answers", nil) {
		t.Fatal("text-only memory candidate should be encoded")
	}
}
