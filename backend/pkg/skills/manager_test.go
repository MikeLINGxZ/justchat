package skills

import (
	"os"
	"path/filepath"
	"testing"
)

func TestManager_DiskShadowsBuiltinIsNoOpWhenBuiltinAbsent(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "user-a", "User A", "body-a", nil)
	writeSkill(t, root, "user-b", "User B", "body-b", nil)

	m := NewManager(root)
	if err := m.Refresh(nil); err != nil {
		t.Fatal(err)
	}
	skills := m.List()
	// 2 user skills + builtin skills (install-cli-from-docs, etc.)
	userSkills := 0
	for _, s := range skills {
		if s.Source == SourceUser {
			userSkills++
		}
	}
	if userSkills != 2 {
		t.Fatalf("expected 2 user skills, got %d", userSkills)
	}
}

func TestManager_Disabled(t *testing.T) {
	root := t.TempDir()
	writeSkill(t, root, "alpha", "Alpha", "x", nil)

	m := NewManager(root)
	// Disable alpha AND all builtin skills so Enabled() returns zero.
	builtin, _ := LoadBuiltin()
	disabled := []string{"alpha"}
	for _, b := range builtin {
		disabled = append(disabled, b.Name)
	}
	if err := m.Refresh(disabled); err != nil {
		t.Fatal(err)
	}
	got, ok := m.Get("alpha")
	if !ok {
		t.Fatal("missing skill alpha")
	}
	if !got.Disabled {
		t.Fatal("expected alpha to be disabled")
	}
	if len(m.Enabled()) != 0 {
		t.Fatalf("expected zero enabled, got %d", len(m.Enabled()))
	}
}

func TestManager_Create_Update_Delete(t *testing.T) {
	root := t.TempDir()
	m := NewManager(root)
	if err := m.Refresh(nil); err != nil {
		t.Fatal(err)
	}

	created, err := m.Create(Skill{
		Name: "x-one", Description: "desc", Body: "body", Source: SourceUser,
	})
	if err != nil {
		t.Fatal(err)
	}
	if created.Name != "x-one" {
		t.Fatalf("unexpected name: %q", created.Name)
	}
	if _, err := os.Stat(filepath.Join(root, "x-one", SkillFileName)); err != nil {
		t.Fatalf("SKILL.md missing: %v", err)
	}

	if _, err := m.Update("x-one", Skill{Name: "x-one", Description: "new", Body: "new-body"}); err != nil {
		t.Fatal(err)
	}
	got, _ := m.Get("x-one")
	if got.Description != "new" {
		t.Fatalf("update did not persist description: %q", got.Description)
	}

	if err := m.Delete("x-one"); err != nil {
		t.Fatal(err)
	}
	if _, ok := m.Get("x-one"); ok {
		t.Fatal("delete did not remove from cache")
	}
}
