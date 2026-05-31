package skills_dto

// UpdateSkillInput provides the fields needed to update an existing user skill.
type UpdateSkillInput struct {
	Name        string `json:"name"`
	NewName     string `json:"new_name,omitempty"`
	Description string `json:"description"`
	Body        string `json:"body"`
}

// UpdateSkillOutput returns the updated skill.
type UpdateSkillOutput struct {
	Skill SkillItem `json:"skill"`
}
