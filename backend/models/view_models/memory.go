package view_models

import "time"

type Memory struct {
	ID           uint      `json:"id"`
	Summary      string    `json:"summary"`
	Content      string    `json:"content"`
	Type         string    `json:"type"`
	IsForgotten  bool      `json:"is_forgotten"`
	RecallCount  int       `json:"recall_count"`
	HasEmbedding bool      `json:"has_embedding"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type MemoryUpdateInput struct {
	Summary string `json:"summary"`
	Content string `json:"content"`
	Type    string `json:"type"`
}

type MemoryListQuery struct {
	Offset      int    `json:"offset"`
	Limit       int    `json:"limit"`
	Keyword     string `json:"keyword"`
	Type        string `json:"type"`
	IsForgotten bool   `json:"is_forgotten"`
}

type MemoryListResponse struct {
	Memories []Memory `json:"memories"`
	Total    int64    `json:"total"`
}

type MemoryStats struct {
	Total     int64 `json:"total"`
	WeekNew   int64 `json:"week_new"`
	Forgotten int64 `json:"forgotten"`
}
