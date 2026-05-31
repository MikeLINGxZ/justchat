package memory_dto

// MemoryItem is the frontend-facing representation of a long-term memory.
type MemoryItem struct {
	ID             uint    `json:"id"`
	Summary        string  `json:"summary"`
	Content        string  `json:"content"`
	Type           string  `json:"type"`
	Target         string  `json:"target"`
	Source         string  `json:"source"`
	CharCount      int     `json:"char_count"`
	EmbeddingID    *uint   `json:"embedding_id"`
	IsForgotten    bool    `json:"is_forgotten"`
	RecallCount    int     `json:"recall_count"`
	LastRecalledAt *string `json:"last_recalled_at"`
	LastUsedAt     *string `json:"last_used_at"`
	Importance     int     `json:"importance"`
	Confidence     int     `json:"confidence"`
	Pinned         bool    `json:"pinned"`
	Created        string  `json:"created"`
	Updated        string  `json:"updated"`
}

// MemoryStatsItem contains aggregate memory counts.
type MemoryStatsItem struct {
	Total     int64            `json:"total"`
	Active    int64            `json:"active"`
	Forgotten int64            `json:"forgotten"`
	ByTarget  map[string]int64 `json:"by_target"`
	ByType    map[string]int64 `json:"by_type"`
}
