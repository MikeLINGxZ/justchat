package skills_dto

// ImportSkillInput provides a local file path to a SKILL.md for importing.
type ImportSkillInput struct {
	Path string `json:"path"`
}

// ImportSkillOutput returns the imported skill.
type ImportSkillOutput struct {
	Skill SkillItem `json:"skill"`
}
