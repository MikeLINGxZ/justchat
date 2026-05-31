package data_models

import "time"

// Memory stores a compact long-term memory that can be injected or retrieved.
type Memory struct {
	OrmModel
	Summary        string     `gorm:"type:varchar(255);index" json:"summary"`
	Content        string     `gorm:"type:text" json:"content"`
	Type           string     `gorm:"type:varchar(32);index" json:"type"`
	Target         string     `gorm:"type:varchar(32);index" json:"target"`
	Source         string     `gorm:"type:varchar(32);index" json:"source"`
	CharCount      int        `json:"char_count"`
	EmbeddingID    *uint      `gorm:"index" json:"embedding_id"`
	IsForgotten    bool       `gorm:"index" json:"is_forgotten"`
	RecallCount    int        `json:"recall_count"`
	LastRecalledAt *time.Time `json:"last_recalled_at"`
	LastUsedAt     *time.Time `json:"last_used_at"`
	Importance     int        `gorm:"index" json:"importance"`
	Confidence     int        `gorm:"index" json:"confidence"`
	Pinned         bool       `gorm:"index" json:"pinned"`
}

// MemoryEmbedding stores optional vector data for a memory.
type MemoryEmbedding struct {
	OrmModel
	MemoryID      uint   `gorm:"index" json:"memory_id"`
	Vector        string `gorm:"type:text" json:"vector"`
	ModelName     string `gorm:"type:varchar(255)" json:"model_name"`
	Dimensions    int    `json:"dimensions"`
	SchemaVersion int    `json:"schema_version"`
}
