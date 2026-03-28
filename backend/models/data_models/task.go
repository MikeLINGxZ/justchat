package data_models

import "time"

type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusRunning   TaskStatus = "running"
	TaskStatusCompleted TaskStatus = "completed"
	TaskStatusFailed    TaskStatus = "failed"
	TaskStatusStopped   TaskStatus = "stopped"
)

type Task struct {
	OrmModel
	TaskUuid             string     `gorm:"uniqueIndex" json:"task_uuid"`
	ChatUuid             string     `gorm:"index" json:"chat_uuid"`
	AssistantMessageUuid string     `gorm:"index" json:"assistant_message_uuid"`
	Status               TaskStatus `gorm:"index;type:varchar(32)" json:"status"`
	EventKey             string     `gorm:"type:varchar(255)" json:"event_key"`
	FinishReason         string     `gorm:"type:varchar(255)" json:"finish_reason"`
	FinishError          string     `gorm:"type:text" json:"finish_error"`
	StartedAt            *time.Time `json:"started_at"`
	FinishedAt           *time.Time `json:"finished_at"`
	LastOutputAt         *time.Time `json:"last_output_at"`
}
