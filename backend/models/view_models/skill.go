package view_models

type SkillSummary struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	When        string   `json:"when"`
	Version     string   `json:"version"`
	Tags        []string `json:"tags"`
}

type SkillDetail struct {
	SkillSummary
	Content string `json:"content"`
}
