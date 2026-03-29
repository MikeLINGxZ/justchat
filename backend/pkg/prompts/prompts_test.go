package prompts

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadPromptCreatesMissingFile(t *testing.T) {
	t.Setenv("LEMONTEA_DATA_PATH", t.TempDir())

	got, err := LoadPrompt("system.test.md", "hello world")
	if err != nil {
		t.Fatalf("LoadPrompt() error = %v", err)
	}
	if got != "hello world" {
		t.Fatalf("LoadPrompt() = %q, want %q", got, "hello world")
	}

	dir, err := PromptDir()
	if err != nil {
		t.Fatalf("PromptDir() error = %v", err)
	}
	content, err := os.ReadFile(filepath.Join(dir, "system.test.md"))
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(content) != "hello world\n" {
		t.Fatalf("file content = %q, want %q", string(content), "hello world\n")
	}
}

func TestLoadPromptPreservesExistingFile(t *testing.T) {
	t.Setenv("LEMONTEA_DATA_PATH", t.TempDir())

	dir, err := PromptDir()
	if err != nil {
		t.Fatalf("PromptDir() error = %v", err)
	}
	path := filepath.Join(dir, "system.custom.md")
	if err := os.WriteFile(path, []byte("custom prompt\n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	got, err := LoadPrompt("system.custom.md", "fallback prompt")
	if err != nil {
		t.Fatalf("LoadPrompt() error = %v", err)
	}
	if got != "custom prompt" {
		t.Fatalf("LoadPrompt() = %q, want %q", got, "custom prompt")
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(content) != "custom prompt\n" {
		t.Fatalf("file content = %q, want %q", string(content), "custom prompt\n")
	}
}

func TestLoadPromptRewritesEmptyFileWithFallback(t *testing.T) {
	t.Setenv("LEMONTEA_DATA_PATH", t.TempDir())

	dir, err := PromptDir()
	if err != nil {
		t.Fatalf("PromptDir() error = %v", err)
	}
	path := filepath.Join(dir, "user.empty.md")
	if err := os.WriteFile(path, []byte("   \n"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	got, err := LoadPrompt("user.empty.md", "fallback body")
	if err == nil {
		t.Fatalf("LoadPrompt() error = nil, want non-nil fallback warning")
	}
	if got != "fallback body" {
		t.Fatalf("LoadPrompt() = %q, want %q", got, "fallback body")
	}

	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(content) != "fallback body\n" {
		t.Fatalf("file content = %q, want %q", string(content), "fallback body\n")
	}
}
