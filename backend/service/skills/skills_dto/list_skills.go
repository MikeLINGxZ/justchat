package skills_dto

// ListSkillsInput is the request payload for listing all skills.
type ListSkillsInput struct{}

// ListSkillsOutput returns all known skills.
type ListSkillsOutput struct {
	Skills []SkillItem `json:"skills"`
}

// SkillItem is the view-model representation of one skill returned to the frontend.
type SkillItem struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Body        string `json:"body,omitempty"`
	Source      string `json:"source"`
	Disabled    bool   `json:"disabled"`
}
