package skills

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeSkill(t *testing.T, root, name, description, body string, sidecar *SidecarManifest) string {
	t.Helper()
	dir := filepath.Join(root, name)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	skill := Skill{Name: name, Description: description, Body: body}
	out, err := Render(skill)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "SKILL.md"), out, 0o644); err != nil {
		t.Fatal(err)
	}
	if sidecar != nil {
		if err := WriteSidecar(dir, *sidecar); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func TestLoadFromDir_UserDefault(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "foo", "Foo skill", "Hello", nil)

	skills, err := LoadFromDir(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) != 1 {
		t.Fatalf("got %d skills", len(skills))
	}
	if skills[0].Source != SourceUser {
		t.Errorf("expected SourceUser, got %q", skills[0].Source)
	}
}

func TestLoadFromDir_AISourceFromSidecar(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "ai-thing", "AI skill", "x", &SidecarManifest{
		Source:    SourceAI,
		CreatedAt: time.Now(),
	})

	skills, err := LoadFromDir(root)
	if err != nil {
		t.Fatal(err)
	}
	if skills[0].Source != SourceAI {
		t.Errorf("expected SourceAI, got %q", skills[0].Source)
	}
}

func TestLoadFromDir_SkipsInvalid(t *testing.T) {
	root := t.TempDir()
	// Valid skill
	writeSkill(t, root, "good", "Good", "x", nil)
	// Invalid: no SKILL.md
	if err := os.MkdirAll(filepath.Join(root, "broken"), 0o755); err != nil {
		t.Fatal(err)
	}

	skills, err := LoadFromDir(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(skills) != 1 || skills[0].Name != "good" {
		t.Fatalf("expected only 'good', got %+v", skills)
	}
}
