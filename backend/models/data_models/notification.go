package data_models

import "time"

// Notification stores a user-facing notification emitted by background tasks.
type Notification struct {
	OrmModel
	SessionID  uint       `gorm:"index" json:"session_id"`
	Kind       string     `gorm:"type:varchar(40);index" json:"kind"`
	Title      string     `gorm:"type:varchar(255)" json:"title"`
	Message    string     `gorm:"type:text" json:"message"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
}
