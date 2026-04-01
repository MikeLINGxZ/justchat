package view_models

type SkillSummary struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Version     string   `json:"version"`
	Tags        []string `json:"tags"`
}

type SkillDetail struct {
	SkillSummary
	Content string `json:"content"`
}
