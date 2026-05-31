package tools

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	pkgskills "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/skills"
)

type fakeSkillCreator struct {
	created []pkgskills.Skill
}

// Create records created skills for assertions.
func (f *fakeSkillCreator) Create(skill pkgskills.Skill) (pkgskills.Skill, error) {
	f.created = append(f.created, skill)
	return skill, nil
}

// TestInvokeProposeSkillPersistsConfirmedAISkill verifies confirmed proposals are saved with AI source.
func TestInvokeProposeSkillPersistsConfirmedAISkill(t *testing.T) {
	requester := &fakeAttentionRequester{}
	creator := &fakeSkillCreator{}

	out, err := InvokeProposeSkill(context.Background(), requester, creator, 9, json.RawMessage(`{
		"name":"demo-skill",
		"description":"demo",
		"body":"hello"
	}`))
	if err != nil {
		t.Fatal(err)
	}
	if len(creator.created) != 1 {
		t.Fatalf("expected 1 created skill, got %d", len(creator.created))
	}
	if creator.created[0].Source != pkgskills.SourceAI {
		t.Fatalf("expected ai source, got %q", creator.created[0].Source)
	}
	if !strings.Contains(out, "demo-skill") {
		t.Fatalf("unexpected tool result: %q", out)
	}
}
