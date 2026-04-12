package service

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/i18n"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/skills"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils/ierror"
)

var skillNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]+$`)

// ListSkills returns a summary list of all available skills.
func (s *Service) ListSkills() ([]view_models.SkillSummary, error) {
	metas, err := skills.ListSkills()
	if err != nil {
		return nil, ierror.NewError(err)
	}

	result := make([]view_models.SkillSummary, 0, len(metas))
	for _, m := range metas {
		tags := m.Tags
		if tags == nil {
			tags = []string{}
		}
		result = append(result, view_models.SkillSummary{
			Name:        m.Name,
			Description: m.Description,
			When:        m.When,
			Version:     m.Version,
			Tags:        tags,
		})
	}
	return result, nil
}

// GetSkill returns the full detail of a skill by name.
func (s *Service) GetSkill(name string) (*view_models.SkillDetail, error) {
	skill, err := skills.LoadSkill(name)
	if err != nil {
		return nil, ierror.NewError(err)
	}

	tags := skill.Tags
	if tags == nil {
		tags = []string{}
	}

	return &view_models.SkillDetail{
		SkillSummary: view_models.SkillSummary{
			Name:        skill.Name,
			Description: skill.Description,
			When:        skill.When,
			Version:     skill.Version,
			Tags:        tags,
		},
		Content: skill.Content,
	}, nil
}

// CreateSkill creates a new skill after validation.
func (s *Service) CreateSkill(input view_models.SkillDetail) (*view_models.SkillDetail, error) {
	if err := validateSkillInput(input); err != nil {
		return nil, ierror.NewError(err)
	}

	if skills.SkillExists(input.Name) {
		return nil, ierror.NewError(fmt.Errorf("skill already exists: %s", input.Name))
	}

	skill := viewModelToSkill(input)
	if err := skills.SaveSkill(skill); err != nil {
		return nil, ierror.NewError(err)
	}

	return s.GetSkill(input.Name)
}

// UpdateSkill updates an existing skill.
func (s *Service) UpdateSkill(name string, input view_models.SkillDetail) (*view_models.SkillDetail, error) {
	if !skills.SkillExists(name) {
		return nil, ierror.NewError(fmt.Errorf("skill not found: %s", name))
	}

	if err := validateSkillInput(input); err != nil {
		return nil, ierror.NewError(err)
	}

	// Use the URL name as the canonical name.
	input.Name = name
	skill := viewModelToSkill(input)
	if err := skills.SaveSkill(skill); err != nil {
		return nil, ierror.NewError(err)
	}

	return s.GetSkill(name)
}

// DeleteSkill removes a skill by name.
func (s *Service) DeleteSkill(name string) error {
	if !skills.SkillExists(name) {
		return ierror.NewError(fmt.Errorf("skill not found: %s", name))
	}

	if err := skills.DeleteSkill(name); err != nil {
		return ierror.NewError(err)
	}
	return nil
}

// SelectSkillFolder opens a folder selection dialog for importing skills.
func (s *Service) SelectSkillFolder() (string, error) {
	path, err := s.app.Dialog.OpenFile().
		CanChooseDirectories(true).
		CanChooseFiles(false).
		SetTitle(i18n.TCurrent("app.dialog.select_skill_folder", nil)).
		PromptForSingleSelection()
	if err != nil {
		return "", ierror.NewError(err)
	}
	return path, nil
}

// ImportSkillsFromFolder scans a folder for .md skill files and imports them.
func (s *Service) ImportSkillsFromFolder(folderPath string) ([]view_models.SkillSummary, error) {
	if folderPath == "" {
		return nil, ierror.NewError(fmt.Errorf("folder path is empty"))
	}

	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return nil, ierror.NewError(err)
	}

	var imported []view_models.SkillSummary
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		raw, err := os.ReadFile(filepath.Join(folderPath, entry.Name()))
		if err != nil {
			continue
		}

		skill, err := skills.ParseSkillContent(raw)
		if err != nil {
			continue
		}

		if !skillNamePattern.MatchString(skill.Name) || skills.SkillExists(skill.Name) {
			continue
		}

		if err := skills.SaveSkill(*skill); err != nil {
			continue
		}

		tags := skill.Tags
		if tags == nil {
			tags = []string{}
		}
		imported = append(imported, view_models.SkillSummary{
			Name:        skill.Name,
			Description: skill.Description,
			When:        skill.When,
			Version:     skill.Version,
			Tags:        tags,
		})
	}

	return imported, nil
}

func validateSkillInput(input view_models.SkillDetail) error {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return fmt.Errorf("skill name cannot be empty")
	}
	if !skillNamePattern.MatchString(name) {
		return fmt.Errorf("skill name must match ^[a-zA-Z0-9_-]+$")
	}
	if strings.TrimSpace(input.Description) == "" {
		return fmt.Errorf("skill description cannot be empty")
	}
	if strings.TrimSpace(input.Content) == "" {
		return fmt.Errorf("skill content cannot be empty")
	}
	return nil
}

func viewModelToSkill(input view_models.SkillDetail) skills.Skill {
	tags := input.Tags
	if tags == nil {
		tags = []string{}
	}
	return skills.Skill{
		SkillMeta: skills.SkillMeta{
			Name:        strings.TrimSpace(input.Name),
			Description: strings.TrimSpace(input.Description),
			When:        strings.TrimSpace(input.When),
			Version:     strings.TrimSpace(input.Version),
			Tags:        tags,
		},
		Content: strings.TrimSpace(input.Content),
	}
}
