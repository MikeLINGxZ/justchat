package data_models

import "time"

// Terminal stores metadata for one PTY-backed command execution.
type Terminal struct {
	OrmModel
	TerminalID    string     `gorm:"type:varchar(80);uniqueIndex" json:"terminal_id"`
	SessionID     uint       `gorm:"index" json:"session_id"`
	MessageID     *uint      `gorm:"index" json:"message_id"`
	ToolCallID    string     `gorm:"type:varchar(255);index" json:"tool_call_id"`
	Title         string     `gorm:"type:varchar(255)" json:"title"`
	Command       string     `gorm:"type:text" json:"command"`
	Args          string     `gorm:"type:text" json:"args"`
	Cwd           string     `gorm:"type:text" json:"cwd"`
	Status        string     `gorm:"type:varchar(32);index" json:"status"`
	Visible       bool       `gorm:"type:bool;default:0;index" json:"visible"`
	PID           int        `json:"pid"`
	ExitCode      *int       `json:"exit_code"`
	CurrentCursor int64      `json:"current_cursor"`
	StartedAt     time.Time  `json:"started_at"`
	EndedAt       *time.Time `json:"ended_at"`
}

// TerminalOutputChunk stores an append-only slice of terminal output.
type TerminalOutputChunk struct {
	OrmModel
	TerminalID  string `gorm:"type:varchar(80);index" json:"terminal_id"`
	Seq         int64  `gorm:"index" json:"seq"`
	CursorStart int64  `gorm:"index" json:"cursor_start"`
	CursorEnd   int64  `gorm:"index" json:"cursor_end"`
	Data        string `gorm:"type:text" json:"data"`
}
