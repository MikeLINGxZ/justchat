package data_models

// TaskState stores resumable key/value state for one automated task session.
type TaskState struct {
	OrmModel
	SessionID uint   `gorm:"index:idx_task_state_session_key,unique" json:"session_id"`
	Key       string `gorm:"type:varchar(120);index:idx_task_state_session_key,unique" json:"key"`
	Value     string `gorm:"type:text" json:"value"`
}
