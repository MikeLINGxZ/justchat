package skills_dto

// CreateSkillInput provides the fields needed to create a new user skill.
type CreateSkillInput struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Body        string `json:"body"`
}

// CreateSkillOutput returns the newly created skill.
type CreateSkillOutput struct {
	Skill SkillItem `json:"skill"`
}
