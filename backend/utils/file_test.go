package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/cloudwego/eino/schema"
)

func TestMimeType2ChatMessagePartType_TextApplicationTypes(t *testing.T) {
	cases := map[string]schema.ChatMessagePartType{
		"text/markdown":    schema.ChatMessagePartTypeText,
		"application/json": schema.ChatMessagePartTypeText,
		"application/yaml": schema.ChatMessagePartTypeText,
		"application/xml":  schema.ChatMessagePartTypeText,
		"application/pdf":  schema.ChatMessagePartTypeFileURL,
		"image/png":        schema.ChatMessagePartTypeImageURL,
		"audio/mpeg":       schema.ChatMessagePartTypeAudioURL,
		"video/mp4":        schema.ChatMessagePartTypeVideoURL,
	}

	for mimeType, want := range cases {
		got, err := MimeType2ChatMessagePartType(mimeType)
		if err != nil {
			t.Fatalf("MimeType2ChatMessagePartType(%q) error = %v", mimeType, err)
		}
		if got != want {
			t.Fatalf("MimeType2ChatMessagePartType(%q) = %q, want %q", mimeType, got, want)
		}
	}
}

func TestDetectMimeType_TextFallbackForUnknownExtension(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "notes.customext")
	if err := os.WriteFile(path, []byte("plain utf8 text"), 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	got, err := DetectMimeType(path)
	if err != nil {
		t.Fatalf("DetectMimeType() error = %v", err)
	}
	if got != "text/plain" {
		t.Fatalf("DetectMimeType() = %q, want %q", got, "text/plain")
	}
}
