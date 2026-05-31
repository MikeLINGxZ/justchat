package memory_dto

type CreateMemoryInput struct {
	Summary    string `json:"summary"`
	Content    string `json:"content"`
	Type       string `json:"type"`
	Target     string `json:"target"`
	Source     string `json:"source"`
	Importance int    `json:"importance"`
	Confidence int    `json:"confidence"`
	Pinned     bool   `json:"pinned"`
}

type CreateMemoryOutput struct {
	Memory MemoryItem `json:"memory"`
}
