package file

import (
	"context"
	"encoding/base64"
	"os"
	"path/filepath"
	"testing"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/file/file_dto"
)

func TestSaveTempFile_Valid(t *testing.T) {
	f := &File{}
	payload := []byte("hello world")
	encoded := base64.StdEncoding.EncodeToString(payload)

	out, err := f.SaveTempFile(context.Background(), file_dto.SaveTempFileInput{
		Name: "test.txt",
		Data: encoded,
		Mime: "text/plain",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out == nil || out.FilePath == "" {
		t.Fatal("expected non-empty file path")
	}
	t.Cleanup(func() { os.Remove(out.FilePath) })

	got, err := os.ReadFile(out.FilePath)
	if err != nil {
		t.Fatalf("read temp file: %v", err)
	}
	if string(got) != "hello world" {
		t.Fatalf("expected 'hello world', got %q", string(got))
	}
}

func TestSaveTempFile_InvalidBase64(t *testing.T) {
	f := &File{}
	_, err := f.SaveTempFile(context.Background(), file_dto.SaveTempFileInput{
		Name: "test.txt",
		Data: "!!!not-valid-base64!!!",
		Mime: "text/plain",
	})
	if err == nil {
		t.Fatal("expected error for invalid base64, got nil")
	}
}

func TestCleanTempDir_RemovesFiles(t *testing.T) {
	tmpDir := t.TempDir() // isolated, auto-cleaned after test
	for _, name := range []string{"a.png", "b.txt"} {
		path := filepath.Join(tmpDir, name)
		if err := os.WriteFile(path, []byte("data"), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}

	cleanTempDir(tmpDir)

	entries, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("readdir: %v", err)
	}
	for _, e := range entries {
		if !e.IsDir() {
			t.Errorf("expected file %s to be removed", e.Name())
		}
	}
}
