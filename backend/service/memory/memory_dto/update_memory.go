package memory_dto

type UpdateMemoryInput struct {
	ID         uint   `json:"id"`
	Summary    string `json:"summary"`
	Content    string `json:"content"`
	Type       string `json:"type"`
	Target     string `json:"target"`
	Source     string `json:"source"`
	Importance int    `json:"importance"`
	Confidence int    `json:"confidence"`
	Pinned     bool   `json:"pinned"`
}

type UpdateMemoryOutput struct {
	Memory MemoryItem `json:"memory"`
}
