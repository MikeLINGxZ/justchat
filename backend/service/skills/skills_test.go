package skills

import (
	"context"
	"errors"
	"strings"
	"testing"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/ierror"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/skills"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/skills/skills_dto"
)

// newTestSkillsService creates a Skills service with an isolated temp-dir manager for testing.
func newTestSkillsService(t *testing.T) *Skills {
	t.Helper()
	tempDir := t.TempDir()
	t.Setenv("LEMONTEA_DATA_DIR", tempDir)
	mgr := skills.NewManager(tempDir)
	return NewSkills(mgr)
}

// TestListSkillsReturnsEmptyInitially verifies that ListSkills returns an empty list
// when no skills have been created.
func TestListSkillsReturnsEmptyInitially(t *testing.T) {
	svc := newTestSkillsService(t)
	ctx := context.Background()

	out, err := svc.ListSkills(ctx, skills_dto.ListSkillsInput{})
	if err != nil {
		t.Fatalf("ListSkills: %v", err)
	}
	if len(out.Skills) != 0 {
		t.Fatalf("expected 0 skills, got %d", len(out.Skills))
	}
}

// TestCreateSkillAndGetSkill verifies the create-then-read round-trip.
func TestCreateSkillAndGetSkill(t *testing.T) {
	svc := newTestSkillsService(t)
	ctx := context.Background()

	createOut, err := svc.CreateSkill(ctx, skills_dto.CreateSkillInput{
		Name:        "test-skill",
		Description: "A test skill",
		Body:        "Do the test thing.",
	})
	if err != nil {
		t.Fatalf("CreateSkill: %v", err)
	}
	if createOut.Skill.Name != "test-skill" {
		t.Fatalf("expected name 'test-skill', got %q", createOut.Skill.Name)
	}
	if createOut.Skill.Source != "user" {
		t.Fatalf("expected source 'user', got %q", createOut.Skill.Source)
	}

	getOut, err := svc.GetSkill(ctx, skills_dto.GetSkillInput{Name: "test-skill"})
	if err != nil {
		t.Fatalf("GetSkill: %v", err)
	}
	if getOut.Skill.Name != "test-skill" {
		t.Fatalf("expected name 'test-skill', got %q", getOut.Skill.Name)
	}
	if getOut.Skill.Description != "A test skill" {
		t.Fatalf("expected description 'A test skill', got %q", getOut.Skill.Description)
	}
	// Body gains a trailing newline from the Render→Parse round-trip; compare trimmed.
	if strings.TrimRight(getOut.Skill.Body, "\n") != "Do the test thing." {
		t.Fatalf("expected body 'Do the test thing.', got %q", getOut.Skill.Body)
	}
}

// TestCreateSkillRejectsInvalidName verifies invalid names fail before writing an unreadable skill.
func TestCreateSkillRejectsInvalidName(t *testing.T) {
	svc := newTestSkillsService(t)
	ctx := context.Background()

	_, err := svc.CreateSkill(ctx, skills_dto.CreateSkillInput{
		Name:        "Bad Skill",
		Description: "A test skill",
		Body:        "Do the test thing.",
	})
	if err == nil {
		t.Fatalf("expected invalid-name error, got nil")
	}
	if !errors.Is(err, ierror.Error(ierror.ErrSkillsInvalidName, errors.New("x"))) {
		t.Fatalf("expected ErrSkillsInvalidName, got %v", err)
	}
}

// TestDeleteSkillRemovesSkill verifies that a created skill can be deleted and is no longer retrievable.
func TestDeleteSkillRemovesSkill(t *testing.T) {
	svc := newTestSkillsService(t)
	ctx := context.Background()

	_, err := svc.CreateSkill(ctx, skills_dto.CreateSkillInput{
		Name:        "delete-me",
		Description: "Temporary",
		Body:        "Will be deleted.",
	})
	if err != nil {
		t.Fatalf("CreateSkill: %v", err)
	}

	_, err = svc.DeleteSkill(ctx, skills_dto.DeleteSkillInput{Name: "delete-me"})
	if err != nil {
		t.Fatalf("DeleteSkill: %v", err)
	}

	_, err = svc.GetSkill(ctx, skills_dto.GetSkillInput{Name: "delete-me"})
	if err == nil {
		t.Fatalf("expected error after deletion, got nil")
	}
}

// TestToggleSkillChangesDisabledState verifies toggling a skill's disabled flag
// and that the state persists across ListSkills calls.
func TestToggleSkillChangesDisabledState(t *testing.T) {
	svc := newTestSkillsService(t)
	ctx := context.Background()

	_, err := svc.CreateSkill(ctx, skills_dto.CreateSkillInput{
		Name:        "toggle-me",
		Description: "Toggleable",
		Body:        "Can be toggled.",
	})
	if err != nil {
		t.Fatalf("CreateSkill: %v", err)
	}

	// Disable the skill.
	toggleOut, err := svc.ToggleSkill(ctx, skills_dto.ToggleSkillInput{
		Name:     "toggle-me",
		Disabled: true,
	})
	if err != nil {
		t.Fatalf("ToggleSkill disable: %v", err)
	}
	if !toggleOut.Skill.Disabled {
		t.Fatalf("expected disabled=true, got false")
	}

	// Verify the disabled state is reflected in ListSkills.
	listOut, err := svc.ListSkills(ctx, skills_dto.ListSkillsInput{})
	if err != nil {
		t.Fatalf("ListSkills: %v", err)
	}
	found := false
	for _, sk := range listOut.Skills {
		if sk.Name == "toggle-me" {
			if !sk.Disabled {
				t.Fatalf("expected disabled=true in list, got false")
			}
			found = true
		}
	}
	if !found {
		t.Fatalf("skill 'toggle-me' not found in ListSkills output")
	}

	// Re-enable the skill.
	toggleOut, err = svc.ToggleSkill(ctx, skills_dto.ToggleSkillInput{
		Name:     "toggle-me",
		Disabled: false,
	})
	if err != nil {
		t.Fatalf("ToggleSkill enable: %v", err)
	}
	if toggleOut.Skill.Disabled {
		t.Fatalf("expected disabled=false, got true")
	}
}

// TestUpdateSkillModifiesContent verifies that updating a skill changes its description and body.
func TestUpdateSkillModifiesContent(t *testing.T) {
	svc := newTestSkillsService(t)
	ctx := context.Background()

	_, err := svc.CreateSkill(ctx, skills_dto.CreateSkillInput{
		Name:        "update-me",
		Description: "Original",
		Body:        "Original body.",
	})
	if err != nil {
		t.Fatalf("CreateSkill: %v", err)
	}

	updateOut, err := svc.UpdateSkill(ctx, skills_dto.UpdateSkillInput{
		Name:        "update-me",
		Description: "Updated",
		Body:        "Updated body.",
	})
	if err != nil {
		t.Fatalf("UpdateSkill: %v", err)
	}
	if updateOut.Skill.Description != "Updated" {
		t.Fatalf("expected description 'Updated', got %q", updateOut.Skill.Description)
	}
	if strings.TrimRight(updateOut.Skill.Body, "\n") != "Updated body." {
		t.Fatalf("expected body 'Updated body.', got %q", updateOut.Skill.Body)
	}
}

// TestUpdateSkillRenamesUserSkill verifies that editing the skill name moves it
// to the new identifier and removes the old identifier from the cache.
func TestUpdateSkillRenamesUserSkill(t *testing.T) {
	svc := newTestSkillsService(t)
	ctx := context.Background()

	_, err := svc.CreateSkill(ctx, skills_dto.CreateSkillInput{
		Name:        "rename-me",
		Description: "Original",
		Body:        "Original body.",
	})
	if err != nil {
		t.Fatalf("CreateSkill: %v", err)
	}
	_, err = svc.ToggleSkill(ctx, skills_dto.ToggleSkillInput{
		Name:     "rename-me",
		Disabled: true,
	})
	if err != nil {
		t.Fatalf("ToggleSkill: %v", err)
	}

	updateOut, err := svc.UpdateSkill(ctx, skills_dto.UpdateSkillInput{
		Name:        "rename-me",
		NewName:     "renamed-skill",
		Description: "Renamed",
		Body:        "Renamed body.",
	})
	if err != nil {
		t.Fatalf("UpdateSkill: %v", err)
	}
	if updateOut.Skill.Name != "renamed-skill" {
		t.Fatalf("expected renamed skill, got %q", updateOut.Skill.Name)
	}
	if _, err = svc.GetSkill(ctx, skills_dto.GetSkillInput{Name: "rename-me"}); err == nil {
		t.Fatalf("expected old skill name to be removed")
	}
	getOut, err := svc.GetSkill(ctx, skills_dto.GetSkillInput{Name: "renamed-skill"})
	if err != nil {
		t.Fatalf("GetSkill renamed: %v", err)
	}
	if !getOut.Skill.Disabled {
		t.Fatalf("expected disabled state to move with renamed skill")
	}
}

// TestToggleSkillPersistsAcrossServiceInstances verifies that the disabled state
// survives config round-trips by reloading disabled names from config.
func TestToggleSkillPersistsAcrossServiceInstances(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("LEMONTEA_DATA_DIR", tempDir)

	// First service instance: create and disable.
	svc1 := NewSkills(skills.NewManager(tempDir))
	ctx := context.Background()

	_, err := svc1.CreateSkill(ctx, skills_dto.CreateSkillInput{
		Name:        "persist-test",
		Description: "Persistence check",
		Body:        "Testing persistence.",
	})
	if err != nil {
		t.Fatalf("CreateSkill: %v", err)
	}

	_, err = svc1.ToggleSkill(ctx, skills_dto.ToggleSkillInput{
		Name:     "persist-test",
		Disabled: true,
	})
	if err != nil {
		t.Fatalf("ToggleSkill: %v", err)
	}

	// Second service instance: should see the disabled flag from config.
	mgr2 := skills.NewManager(tempDir)
	svc2 := NewSkills(mgr2)
	svc2.refreshDisabledFromConfig()

	getOut, err := svc2.GetSkill(ctx, skills_dto.GetSkillInput{Name: "persist-test"})
	if err != nil {
		t.Fatalf("GetSkill: %v", err)
	}
	if !getOut.Skill.Disabled {
		t.Fatalf("expected disabled=true in new service instance, got false")
	}
}

// TestGetSkillReturnsErrorForMissing verifies that GetSkill returns an error for a non-existent skill.
func TestGetSkillReturnsErrorForMissing(t *testing.T) {
	svc := newTestSkillsService(t)
	ctx := context.Background()

	_, err := svc.GetSkill(ctx, skills_dto.GetSkillInput{Name: "nonexistent"})
	if err == nil {
		t.Fatalf("expected error for nonexistent skill, got nil")
	}
}
