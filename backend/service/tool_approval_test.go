package service

import (
	"context"
	"testing"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/tool_approval"
)

func TestRespondToolApprovalResolvesPendingApproval(t *testing.T) {
	svc, st := newTaskRecoveryTestService(t)
	now := time.Now()
	approval := data_models.ToolApproval{
		ApprovalID:           "approval-live-1",
		TaskUuid:             "task-live-approval",
		ChatUuid:             "chat-live-approval",
		AssistantMessageUuid: "assistant-live-approval",
		ToolCallID:           "call-live-approval",
		ToolID:               "file_tool",
		ToolName:             "文件工具",
		Status:               data_models.ToolApprovalStatusPending,
		Title:                "读取文件",
		Message:              "等待用户确认",
		RequestedAt:          &now,
	}
	if err := st.CreateToolApproval(context.Background(), approval); err != nil {
		t.Fatalf("CreateToolApproval() error = %v", err)
	}
	if err := tool_approval.Manager.Register(approval.ApprovalID); err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	defer tool_approval.Manager.Cancel(approval.ApprovalID)

	if err := svc.RespondToolApproval(data_models.ToolApprovalResponse{
		ApprovalID: approval.ApprovalID,
		Decision:   data_models.ToolApprovalDecisionCustom,
		Comment:    "先说明要读哪个文件",
	}); err != nil {
		t.Fatalf("RespondToolApproval() error = %v", err)
	}

	result, err := tool_approval.Manager.Wait(context.Background(), approval.ApprovalID)
	if err != nil {
		t.Fatalf("Wait() error = %v", err)
	}
	if result.Decision != data_models.ToolApprovalDecisionCustom {
		t.Fatalf("decision = %q, want %q", result.Decision, data_models.ToolApprovalDecisionCustom)
	}
	if result.Comment != "先说明要读哪个文件" {
		t.Fatalf("comment = %q, want custom comment", result.Comment)
	}

	stored, err := st.GetToolApprovalByApprovalID(context.Background(), approval.ApprovalID)
	if err != nil {
		t.Fatalf("GetToolApprovalByApprovalID() error = %v", err)
	}
	if stored == nil {
		t.Fatal("stored approval = nil")
	}
	if stored.Status != data_models.ToolApprovalStatusResolved {
		t.Fatalf("status = %q, want resolved", stored.Status)
	}
	if stored.Decision != data_models.ToolApprovalDecisionCustom {
		t.Fatalf("decision = %q, want custom", stored.Decision)
	}
	if stored.ResponseComment != "先说明要读哪个文件" {
		t.Fatalf("response_comment = %q, want custom comment", stored.ResponseComment)
	}
}
