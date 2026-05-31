package agent

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"trpc.group/trpc-go/trpc-agent-go/model"
)

// 67-byte minimal valid PNG (1x1 transparent)
const minimalPNGBase64 = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR4nGNgYGD4DwABBAEAfbLI3wAAAABJRU5ErkJggg=="

func writeTempFile(t *testing.T, name string, data []byte) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
	return path
}

func TestMarshalUnmarshalAttachmentsRoundTrip(t *testing.T) {
	in := []Attachment{
		{Name: "a.png", Path: "/abs/a.png", Mime: "image/png", Kind: "image"},
		{Name: "b.pdf", Path: "/abs/b.pdf", Mime: "application/pdf", Kind: "file"},
	}

	s, err := MarshalAttachments(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if s == "" {
		t.Fatalf("expected non-empty json, got empty")
	}

	out, err := UnmarshalAttachments(s)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(out) != 2 || out[0].Name != "a.png" || out[1].Kind != "file" {
		t.Fatalf("round trip mismatch: %+v", out)
	}
}

func TestMarshalEmptyReturnsEmptyString(t *testing.T) {
	got, err := MarshalAttachments(nil)
	if err != nil || got != "" {
		t.Fatalf("expected empty,nil; got %q,%v", got, err)
	}
}

func TestUnmarshalEmptyReturnsNil(t *testing.T) {
	got, err := UnmarshalAttachments("")
	if err != nil || got != nil {
		t.Fatalf("expected nil,nil; got %v,%v", got, err)
	}
}

func TestNormalizeAttachmentFillsMissingFields(t *testing.T) {
	a := NormalizeAttachment(Attachment{Path: "/tmp/x.PNG"})
	if a.Name != "x.PNG" {
		t.Fatalf("name = %q", a.Name)
	}
	if a.Mime != "image/png" {
		t.Fatalf("mime = %q", a.Mime)
	}
	if a.Kind != "image" {
		t.Fatalf("kind = %q", a.Kind)
	}

	b := NormalizeAttachment(Attachment{Path: "/tmp/data.bin"})
	if b.Mime == "" || b.Kind != "file" {
		t.Fatalf("unknown ext should be file, got mime=%q kind=%q", b.Mime, b.Kind)
	}
}

func TestBuildUserMessageNoAttachments(t *testing.T) {
	msg := BuildUserMessage("hello", nil)
	if msg.Role != model.RoleUser {
		t.Fatalf("role = %v", msg.Role)
	}
	if msg.Content != "hello" {
		t.Fatalf("content = %q", msg.Content)
	}
	if len(msg.ContentParts) != 0 {
		t.Fatalf("expected no ContentParts, got %d", len(msg.ContentParts))
	}
}

func TestBuildUserMessageWithImage(t *testing.T) {
	png, _ := base64.StdEncoding.DecodeString(minimalPNGBase64)
	path := writeTempFile(t, "tiny.png", png)

	msg := BuildUserMessage("see image", []Attachment{
		NormalizeAttachment(Attachment{Path: path}),
	})

	if msg.Content != "see image" {
		t.Fatalf("content = %q", msg.Content)
	}
	if len(msg.ContentParts) != 1 || msg.ContentParts[0].Type != model.ContentTypeImage {
		t.Fatalf("expected one image part, got %+v", msg.ContentParts)
	}
	if msg.ContentParts[0].Image == nil || len(msg.ContentParts[0].Image.Data) == 0 {
		t.Fatalf("image data empty")
	}
}

func TestBuildUserMessageWithArbitraryFile(t *testing.T) {
	path := writeTempFile(t, "notes.go", []byte("package x"))

	msg := BuildUserMessage("review please", []Attachment{
		NormalizeAttachment(Attachment{Path: path}),
	})

	if len(msg.ContentParts) != 1 || msg.ContentParts[0].Type != model.ContentTypeFile {
		t.Fatalf("expected file part, got %+v", msg.ContentParts)
	}
	if msg.ContentParts[0].File == nil || msg.ContentParts[0].File.Name != "notes.go" {
		t.Fatalf("file part incorrect: %+v", msg.ContentParts[0].File)
	}
}

func TestBuildUserMessageMissingFileDegrades(t *testing.T) {
	msg := BuildUserMessage("with missing", []Attachment{
		{Name: "ghost.png", Path: "/nonexistent/ghost.png", Mime: "image/png", Kind: "image"},
	})

	if len(msg.ContentParts) != 0 {
		t.Fatalf("expected no ContentParts, got %d", len(msg.ContentParts))
	}
	if !strings.Contains(msg.Content, "[Missing attachment: ghost.png]") {
		t.Fatalf("content missing placeholder: %q", msg.Content)
	}
}

func TestBuildUserMessageReclassifiesOctetStreamImage(t *testing.T) {
	png, _ := base64.StdEncoding.DecodeString(minimalPNGBase64)
	path := writeTempFile(t, "blob.bin", png)

	// Force the file branch: explicitly set Mime to octet-stream and Kind to file.
	msg := BuildUserMessage("ok", []Attachment{
		{Name: "blob.bin", Path: path, Mime: "application/octet-stream", Kind: "file"},
	})

	if len(msg.ContentParts) != 1 || msg.ContentParts[0].Type != model.ContentTypeImage {
		t.Fatalf("expected reclassified image part, got %+v", msg.ContentParts)
	}
}

func TestBuildUserMessageMixedSomeMissing(t *testing.T) {
	png, _ := base64.StdEncoding.DecodeString(minimalPNGBase64)
	good := writeTempFile(t, "ok.png", png)

	msg := BuildUserMessage("mixed", []Attachment{
		NormalizeAttachment(Attachment{Path: good}),
		{Name: "lost.pdf", Path: "/nope/lost.pdf", Mime: "application/pdf", Kind: "file"},
	})

	if len(msg.ContentParts) != 1 {
		t.Fatalf("expected 1 ContentPart, got %d", len(msg.ContentParts))
	}
	if !strings.Contains(msg.Content, "[Missing attachment: lost.pdf]") {
		t.Fatalf("content missing placeholder: %q", msg.Content)
	}
}
