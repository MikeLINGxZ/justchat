package skills

import (
	"context"
	"errors"
	"os"
	"strings"

	"github.com/wailsapp/wails/v3/pkg/application"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/dir"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/ierror"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/skills"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/skills/skills_dto"
)

// Skills exposes skill management operations to the Wails frontend via RPC.
type Skills struct {
	manager  *skills.Manager
	wailsApp *application.App
}

// NewSkills creates a Skills service backed by the given skills manager.
func NewSkills(manager *skills.Manager) *Skills {
	return &Skills{manager: manager}
}

// Manager returns the underlying skills manager, used by the agent service
// to wire the Skill meta-tool provider.
func (s *Skills) Manager() *skills.Manager {
	return s.manager
}

// ListSkills returns all known skills sorted by name.
func (s *Skills) ListSkills(ctx context.Context, input skills_dto.ListSkillsInput) (*skills_dto.ListSkillsOutput, error) {
	all := s.manager.List()
	items := make([]skills_dto.SkillItem, 0, len(all))
	for _, sk := range all {
		items = append(items, toSkillItem(sk))
	}
	return &skills_dto.ListSkillsOutput{Skills: items}, nil
}

// GetSkill returns a single skill by name.
func (s *Skills) GetSkill(ctx context.Context, input skills_dto.GetSkillInput) (*skills_dto.GetSkillOutput, error) {
	sk, ok := s.manager.Get(input.Name)
	if !ok {
		return nil, ierror.Error(ierror.ErrSkillsNotFound, os.ErrNotExist)
	}
	return &skills_dto.GetSkillOutput{Skill: toSkillItem(sk)}, nil
}

// CreateSkill creates a new user-owned skill and persists it to disk.
func (s *Skills) CreateSkill(ctx context.Context, input skills_dto.CreateSkillInput) (*skills_dto.CreateSkillOutput, error) {
	sk := skills.Skill{
		Name:        input.Name,
		Description: input.Description,
		Body:        input.Body,
		Source:      skills.SourceUser,
	}
	created, err := s.manager.Create(sk)
	if err != nil {
		if errors.Is(err, skills.ErrInvalidSkillName) {
			return nil, ierror.Error(ierror.ErrSkillsInvalidName, err)
		}
		if errors.Is(err, skills.ErrInvalidSkillContent) {
			return nil, ierror.Error(ierror.ErrSkillsInvalidContent, err)
		}
		return nil, ierror.Error(ierror.ErrSkillsWriteFailed, err)
	}
	// Re-sync the in-memory disabled state from persisted config.
	s.refreshDisabledFromConfig()
	return &skills_dto.CreateSkillOutput{Skill: toSkillItem(created)}, nil
}

// UpdateSkill updates an existing non-builtin skill's description and body.
func (s *Skills) UpdateSkill(ctx context.Context, input skills_dto.UpdateSkillInput) (*skills_dto.UpdateSkillOutput, error) {
	existing, ok := s.manager.Get(input.Name)
	if !ok {
		return nil, ierror.Error(ierror.ErrSkillsNotFound, os.ErrNotExist)
	}
	if existing.Source == skills.SourceBuiltin {
		return nil, ierror.Error(ierror.ErrSkillsBuiltinLocked, os.ErrPermission)
	}
	targetName := strings.TrimSpace(input.NewName)
	if targetName == "" {
		targetName = input.Name
	}
	update := skills.Skill{
		Name:        targetName,
		Description: input.Description,
		Body:        input.Body,
	}
	updated, err := s.manager.Update(input.Name, update)
	if err != nil {
		if errors.Is(err, skills.ErrInvalidSkillName) {
			return nil, ierror.Error(ierror.ErrSkillsInvalidName, err)
		}
		if errors.Is(err, skills.ErrInvalidSkillContent) {
			return nil, ierror.Error(ierror.ErrSkillsInvalidContent, err)
		}
		if errors.Is(err, skills.ErrSkillNameTaken) {
			return nil, ierror.Error(ierror.ErrSkillsNameTaken, err)
		}
		return nil, ierror.Error(ierror.ErrSkillsWriteFailed, err)
	}
	if updated.Name != input.Name {
		if err := s.renameDisabledSkill(input.Name, updated.Name); err != nil {
			return nil, ierror.Error(ierror.ErrSettingsSaveConfig, err)
		}
	}
	// Restore the pre-update disabled flag that the manager's Refresh cleared.
	_ = s.manager.SetDisabled(updated.Name, existing.Disabled)
	s.refreshDisabledFromConfig()
	return &skills_dto.UpdateSkillOutput{Skill: toSkillItem(updated)}, nil
}

// DeleteSkill removes a non-builtin skill from disk and cache.
func (s *Skills) DeleteSkill(ctx context.Context, input skills_dto.DeleteSkillInput) (*skills_dto.DeleteSkillOutput, error) {
	existing, ok := s.manager.Get(input.Name)
	if !ok {
		return nil, ierror.Error(ierror.ErrSkillsNotFound, os.ErrNotExist)
	}
	if existing.Source == skills.SourceBuiltin {
		return nil, ierror.Error(ierror.ErrSkillsBuiltinLocked, os.ErrPermission)
	}
	if err := s.manager.Delete(input.Name); err != nil {
		return nil, ierror.Error(ierror.ErrSkillsDeleteFailed, err)
	}
	s.refreshDisabledFromConfig()
	return &skills_dto.DeleteSkillOutput{}, nil
}

// ToggleSkill enables or disables a skill and persists the change to config.json.
func (s *Skills) ToggleSkill(ctx context.Context, input skills_dto.ToggleSkillInput) (*skills_dto.ToggleSkillOutput, error) {
	cfg, err := s.loadConfig()
	if err != nil {
		return nil, ierror.Error(ierror.ErrSettingsLoadConfig, err)
	}

	// Build the updated disabled names list.
	disabled := make(map[string]bool, len(cfg.DisabledSkills))
	for _, name := range cfg.DisabledSkills {
		disabled[name] = true
	}
	if input.Disabled {
		disabled[input.Name] = true
	} else {
		delete(disabled, input.Name)
	}
	names := make([]string, 0, len(disabled))
	for name := range disabled {
		names = append(names, name)
	}
	cfg.DisabledSkills = names

	// Persist the updated config before touching the manager.
	if err := s.saveConfig(cfg); err != nil {
		return nil, ierror.Error(ierror.ErrSettingsSaveConfig, err)
	}

	// Apply the disabled flag to the in-memory cache.
	if err := s.manager.SetDisabled(input.Name, input.Disabled); err != nil {
		return nil, ierror.Error(ierror.ErrSkillsNotFound, err)
	}

	sk, ok := s.manager.Get(input.Name)
	if !ok {
		return nil, ierror.Error(ierror.ErrSkillsNotFound, os.ErrNotExist)
	}
	return &skills_dto.ToggleSkillOutput{Skill: toSkillItem(sk)}, nil
}

// ImportSkill reads a SKILL.md from a local path, writes it to the skills directory, and reloads.
func (s *Skills) ImportSkill(ctx context.Context, input skills_dto.ImportSkillInput) (*skills_dto.ImportSkillOutput, error) {
	raw, err := os.ReadFile(input.Path)
	if err != nil {
		return nil, ierror.Error(ierror.ErrSkillsLoadFailed, err)
	}
	sk, err := skills.Parse(raw)
	if err != nil {
		return nil, ierror.Error(ierror.ErrSkillsInvalidContent, err)
	}
	sk.Source = skills.SourceUser

	// Write the parsed skill to the skills root directory.
	dataDir, err := dir.GetDataDir()
	if err != nil {
		return nil, ierror.Error(ierror.ErrSkillsWriteFailed, err)
	}
	skillsRoot := dir.SkillsRoot(dataDir)
	if err := skills.WriteSkill(skillsRoot, sk); err != nil {
		return nil, ierror.Error(ierror.ErrSkillsWriteFailed, err)
	}

	// Reload the manager with the current disabled names from config.
	s.refreshDisabledFromConfig()

	imported, ok := s.manager.Get(sk.Name)
	if !ok {
		return nil, ierror.Error(ierror.ErrSkillsNotFound, os.ErrNotExist)
	}
	return &skills_dto.ImportSkillOutput{Skill: toSkillItem(imported)}, nil
}

// toSkillItem converts an internal Skill to the DTO SkillItem for frontend consumption.
func toSkillItem(sk skills.Skill) skills_dto.SkillItem {
	return skills_dto.SkillItem{
		Name:        sk.Name,
		Description: sk.Description,
		Body:        sk.Body,
		Source:      string(sk.Source),
		Disabled:    sk.Disabled,
	}
}

// refreshDisabledFromConfig reloads the manager's in-memory cache using disabled names from config.json.
// This keeps the disabled state consistent between the on-disk config and the manager after mutations.
func (s *Skills) refreshDisabledFromConfig() {
	cfg, err := s.loadConfig()
	if err != nil {
		// Fall back to no disabled skills when config cannot be read.
		_ = s.manager.Refresh(nil)
		return
	}
	_ = s.manager.Refresh(cfg.DisabledSkills)
}
