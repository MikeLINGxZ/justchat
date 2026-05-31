package skills_dto

// GetSkillInput requests a single skill by name.
type GetSkillInput struct {
	Name string `json:"name"`
}

// GetSkillOutput returns the requested skill.
type GetSkillOutput struct {
	Skill SkillItem `json:"skill"`
}
