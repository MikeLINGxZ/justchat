package skills

import (
	"strings"
	"testing"
)

func TestParseSkill_Valid(t *testing.T) {
	raw := strings.Join([]string{
		"---",
		"name: install-cli-from-docs",
		"description: Install a CLI plugin from official documentation",
		"---",
		"",
		"Body line 1",
		"Body line 2",
	}, "\n")

	got, err := Parse([]byte(raw))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got.Name != "install-cli-from-docs" {
		t.Errorf("name: got %q", got.Name)
	}
	if got.Description == "" {
		t.Errorf("description must not be empty")
	}
	if !strings.Contains(got.Body, "Body line 1") {
		t.Errorf("body lost: %q", got.Body)
	}
}

func TestParseSkill_MissingFrontmatter(t *testing.T) {
	_, err := Parse([]byte("no frontmatter here"))
	if err == nil {
		t.Fatal("expected error for missing frontmatter")
	}
}

func TestParseSkill_MissingName(t *testing.T) {
	raw := "---\ndescription: x\n---\nbody"
	_, err := Parse([]byte(raw))
	if err == nil {
		t.Fatal("expected error for missing name")
	}
}

func TestParseSkill_InvalidNameChars(t *testing.T) {
	raw := "---\nname: Bad Name!\ndescription: x\n---\nbody"
	_, err := Parse([]byte(raw))
	if err == nil {
		t.Fatal("expected error for invalid name chars")
	}
}
