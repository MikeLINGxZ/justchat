package service

import (
	"os"
	"path/filepath"
	"testing"
)

func TestListPromptFilesOrder(t *testing.T) {
	t.Setenv("LEMONTEA_DATA_PATH", t.TempDir())

	svc := NewService()
	items, err := svc.ListPromptFiles()
	if err != nil {
		t.Fatalf("ListPromptFiles() error = %v", err)
	}
	if len(items) != 12 {
		t.Fatalf("ListPromptFiles() len = %d, want 12", len(items))
	}

	want := []string{
		"system.entry.md",
		"system.planner.md",
		"user.planning.md",
		"system.worker.general.md",
		"system.worker.tool.md",
		"system.synthesizer.md",
		"user.synthesis.md",
		"system.reviewer.md",
		"user.review.md",
		"system.title.md",
		"system.memory.md",
		"system.main_agent.md",
	}
	for i, item := range items {
		if item.Name != want[i] {
			t.Fatalf("ListPromptFiles()[%d] = %q, want %q", i, item.Name, want[i])
		}
	}
}

func TestGetPromptFileCreatesMissingDefault(t *testing.T) {
	t.Setenv("LEMONTEA_DATA_PATH", t.TempDir())

	svc := NewService()
	detail, err := svc.GetPromptFile("system.entry.md")
	if err != nil {
		t.Fatalf("GetPromptFile() error = %v", err)
	}
	if detail == nil || detail.Content == "" {
		t.Fatalf("GetPromptFile() = %#v, want non-empty content", detail)
	}

	path := filepath.Join(os.Getenv("LEMONTEA_DATA_PATH"), "prompt", "system.entry.md")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("prompt file not created: %v", err)
	}
}

func TestUpdatePromptFileRefreshesPromptCache(t *testing.T) {
	t.Setenv("LEMONTEA_DATA_PATH", t.TempDir())

	svc := NewService()
	if err := svc.reloadPromptSet(); err != nil {
		t.Fatalf("reloadPromptSet() error = %v", err)
	}

	content := "你是新的主入口提示词。"
	detail, err := svc.UpdatePromptFile("system.entry.md", content)
	if err != nil {
		t.Fatalf("UpdatePromptFile() error = %v", err)
	}
	if detail.Content != content {
		t.Fatalf("UpdatePromptFile().Content = %q, want %q", detail.Content, content)
	}
	if svc.prompts.EntrySystem != content {
		t.Fatalf("svc.prompts.EntrySystem = %q, want %q", svc.prompts.EntrySystem, content)
	}
}

func TestUpdatePromptFileRejectsUnknownOrEmpty(t *testing.T) {
	t.Setenv("LEMONTEA_DATA_PATH", t.TempDir())

	svc := NewService()
	if _, err := svc.UpdatePromptFile("unknown.md", "hello"); err == nil {
		t.Fatal("UpdatePromptFile() unknown file error = nil, want non-nil")
	}
	if _, err := svc.UpdatePromptFile("system.entry.md", "   "); err == nil {
		t.Fatal("UpdatePromptFile() empty content error = nil, want non-nil")
	}
}

func TestResetPromptFileRestoresDefault(t *testing.T) {
	t.Setenv("LEMONTEA_DATA_PATH", t.TempDir())

	svc := NewService()
	if _, err := svc.UpdatePromptFile("system.entry.md", "临时内容"); err != nil {
		t.Fatalf("UpdatePromptFile() error = %v", err)
	}

	detail, err := svc.ResetPromptFile("system.entry.md")
	if err != nil {
		t.Fatalf("ResetPromptFile() error = %v", err)
	}
	if detail == nil || detail.Content == "" {
		t.Fatalf("ResetPromptFile() = %#v, want non-empty content", detail)
	}
	if detail.Content == "临时内容" {
		t.Fatalf("ResetPromptFile() content not restored, got %q", detail.Content)
	}
}
