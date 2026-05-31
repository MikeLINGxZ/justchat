// backend/pkg/prompt/prompt_test.go
package prompt

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_ReturnsDefault_WhenNoCustomFile(t *testing.T) {
	registry = make(map[string]string)
	Register("test_agent", "You are a helpful assistant.")

	t.Setenv("LEMONTEA_DATA_DIR", t.TempDir())

	got, err := Load("test_agent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "You are a helpful assistant." {
		t.Fatalf("expected default prompt, got: %s", got)
	}
}

func TestLoad_ReturnsCustom_WhenFileExists(t *testing.T) {
	registry = make(map[string]string)
	Register("test_agent", "default prompt")

	tmpDir := t.TempDir()
	t.Setenv("LEMONTEA_DATA_DIR", tmpDir)

	promptDir := filepath.Join(tmpDir, "prompt", "test_agent")
	if err := os.MkdirAll(promptDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(promptDir, "index.md"), []byte("custom prompt content"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := Load("test_agent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "custom prompt content" {
		t.Fatalf("expected custom prompt, got: %s", got)
	}
}

func TestLoad_ErrorsOnUnregistered(t *testing.T) {
	registry = make(map[string]string)

	t.Setenv("LEMONTEA_DATA_DIR", t.TempDir())

	_, err := Load("nonexistent")
	if err == nil {
		t.Fatal("expected error for unregistered prompt")
	}
}
