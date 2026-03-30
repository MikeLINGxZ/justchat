package data_models

import "time"

type ToolApprovalStatus string

const (
	ToolApprovalStatusPending  ToolApprovalStatus = "pending"
	ToolApprovalStatusResolved ToolApprovalStatus = "resolved"
	ToolApprovalStatusExpired  ToolApprovalStatus = "expired"
)

type ToolApprovalDecision string

const (
	ToolApprovalDecisionAllow  ToolApprovalDecision = "allow"
	ToolApprovalDecisionReject ToolApprovalDecision = "reject"
	ToolApprovalDecisionCustom ToolApprovalDecision = "custom"
)

type ToolApproval struct {
	OrmModel
	ApprovalID           string               `gorm:"uniqueIndex;type:varchar(255)" json:"approval_id"`
	TaskUuid             string               `gorm:"index;type:varchar(255)" json:"task_uuid"`
	ChatUuid             string               `gorm:"index;type:varchar(255)" json:"chat_uuid"`
	AssistantMessageUuid string               `gorm:"index;type:varchar(255)" json:"assistant_message_uuid"`
	ToolCallID           string               `gorm:"index;type:varchar(255)" json:"tool_call_id"`
	ToolID               string               `gorm:"index;type:varchar(255)" json:"tool_id"`
	ToolName             string               `gorm:"type:varchar(255)" json:"tool_name"`
	Status               ToolApprovalStatus   `gorm:"index;type:varchar(32)" json:"status"`
	Decision             ToolApprovalDecision `gorm:"type:varchar(32)" json:"decision"`
	Title                string               `gorm:"type:text" json:"title"`
	Message              string               `gorm:"type:text" json:"message"`
	Scope                string               `gorm:"type:text" json:"scope"`
	ArgumentsJSON        string               `gorm:"type:text" json:"arguments_json"`
	ResponseComment      string               `gorm:"type:text" json:"response_comment"`
	RequestedAt          *time.Time           `json:"requested_at"`
	RespondedAt          *time.Time           `json:"responded_at"`
}

type ToolApprovalSummary struct {
	ApprovalID  string               `json:"approval_id"`
	ToolCallID  string               `json:"tool_call_id"`
	ToolID      string               `json:"tool_id"`
	ToolName    string               `json:"tool_name"`
	Status      ToolApprovalStatus   `json:"status"`
	Decision    ToolApprovalDecision `json:"decision"`
	Title       string               `json:"title"`
	Message     string               `json:"message"`
	Scope       string               `json:"scope"`
	RequestedAt *time.Time           `json:"requested_at"`
	RespondedAt *time.Time           `json:"responded_at"`
}

type ToolApprovalResponse struct {
	ApprovalID string               `json:"approval_id"`
	Decision   ToolApprovalDecision `json:"decision"`
	Comment    string               `json:"comment"`
}

func (t ToolApproval) Summary() ToolApprovalSummary {
	return ToolApprovalSummary{
		ApprovalID:  t.ApprovalID,
		ToolCallID:  t.ToolCallID,
		ToolID:      t.ToolID,
		ToolName:    t.ToolName,
		Status:      t.Status,
		Decision:    t.Decision,
		Title:       t.Title,
		Message:     t.Message,
		Scope:       t.Scope,
		RequestedAt: t.RequestedAt,
		RespondedAt: t.RespondedAt,
	}
}
