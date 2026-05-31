package tools

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	pkgskills "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/skills"
)

// stubProvider is a minimal SkillProvider for unit tests.
type stubProvider struct{ items []pkgskills.Skill }

// Enabled returns only non-disabled skills from the stub list.
func (s stubProvider) Enabled() []pkgskills.Skill {
	out := make([]pkgskills.Skill, 0, len(s.items))
	for _, sk := range s.items {
		if !sk.Disabled {
			out = append(out, sk)
		}
	}
	return out
}

// Get returns a skill by name from the stub list.
func (s stubProvider) Get(name string) (pkgskills.Skill, bool) {
	for _, sk := range s.items {
		if sk.Name == name {
			return sk, true
		}
	}
	return pkgskills.Skill{}, false
}

func TestBuildSkillTool_ListsEnabled(t *testing.T) {
	p := stubProvider{items: []pkgskills.Skill{
		{Name: "foo", Description: "Foo desc"},
		{Name: "bar", Description: "Bar desc", Disabled: true},
	}}
	meta := BuildSkillTool(p)
	if meta.Name != SkillToolName {
		t.Fatalf("name: got %q, want %q", meta.Name, SkillToolName)
	}
	if !strings.Contains(meta.Description, "foo: Foo desc") {
		t.Fatalf("description missing enabled skill 'foo': %s", meta.Description)
	}
	if strings.Contains(meta.Description, "bar:") {
		t.Fatalf("description should not list disabled skill 'bar': %s", meta.Description)
	}
}

func TestBuildSkillTool_EmptyProvider(t *testing.T) {
	p := stubProvider{}
	meta := BuildSkillTool(p)
	if !strings.Contains(meta.Description, "(none)") {
		t.Fatalf("expected '(none)' in description for empty provider: %s", meta.Description)
	}
}

func TestBuildSkillTool_FormatPurpose(t *testing.T) {
	p := stubProvider{items: []pkgskills.Skill{{Name: "foo", Description: "Foo"}}}
	meta := BuildSkillTool(p)
	if meta.FormatPurpose == nil {
		t.Fatal("FormatPurpose should not be nil")
	}
	args := json.RawMessage(`{"name":"foo"}`)
	purpose := meta.FormatPurpose(args)
	if purpose != "Loading skill: foo" {
		t.Fatalf("FormatPurpose: got %q, want %q", purpose, "Loading skill: foo")
	}
}

func TestInvokeSkill_Success(t *testing.T) {
	p := stubProvider{items: []pkgskills.Skill{
		{Name: "foo", Body: "# Foo instructions"},
	}}
	body, err := InvokeSkill(context.Background(), p, json.RawMessage(`{"name":"foo"}`))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if body != "# Foo instructions" {
		t.Fatalf("body: got %q", body)
	}
}

func TestInvokeSkill_NotFound(t *testing.T) {
	p := stubProvider{}
	_, err := InvokeSkill(context.Background(), p, json.RawMessage(`{"name":"missing"}`))
	if err == nil {
		t.Fatal("expected error for missing skill")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Fatalf("error should mention 'not found': %v", err)
	}
}

func TestInvokeSkill_Disabled(t *testing.T) {
	p := stubProvider{items: []pkgskills.Skill{
		{Name: "foo", Body: "body", Disabled: true},
	}}
	_, err := InvokeSkill(context.Background(), p, json.RawMessage(`{"name":"foo"}`))
	if err == nil {
		t.Fatal("expected error for disabled skill")
	}
	if !strings.Contains(err.Error(), "disabled") {
		t.Fatalf("error should mention 'disabled': %v", err)
	}
}

func TestInvokeSkill_EmptyName(t *testing.T) {
	p := stubProvider{}
	_, err := InvokeSkill(context.Background(), p, json.RawMessage(`{}`))
	if err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestInvokeSkill_InvalidJSON(t *testing.T) {
	p := stubProvider{}
	_, err := InvokeSkill(context.Background(), p, json.RawMessage(`not-json`))
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
