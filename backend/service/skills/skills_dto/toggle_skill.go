package skills_dto

// ToggleSkillInput identifies the skill and desired disabled state.
type ToggleSkillInput struct {
	Name     string `json:"name"`
	Disabled bool   `json:"disabled"`
}

// ToggleSkillOutput returns the skill after toggling its disabled state.
type ToggleSkillOutput struct {
	Skill SkillItem `json:"skill"`
}
