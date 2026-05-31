package skills_dto

// DeleteSkillInput identifies the skill to delete by name.
type DeleteSkillInput struct {
	Name string `json:"name"`
}

// DeleteSkillOutput is returned on successful deletion.
type DeleteSkillOutput struct{}
