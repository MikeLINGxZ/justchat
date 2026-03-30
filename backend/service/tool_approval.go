package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/tool_approval"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils/ierror"
)

func (s *Service) RespondToolApproval(input data_models.ToolApprovalResponse) error {
	if strings.TrimSpace(input.ApprovalID) == "" {
		return ierror.NewError(fmt.Errorf("approval_id is required"))
	}
	if input.Decision != data_models.ToolApprovalDecisionAllow &&
		input.Decision != data_models.ToolApprovalDecisionReject &&
		input.Decision != data_models.ToolApprovalDecisionCustom {
		return ierror.NewError(fmt.Errorf("unsupported decision: %s", input.Decision))
	}
	if input.Decision == data_models.ToolApprovalDecisionCustom && strings.TrimSpace(input.Comment) == "" {
		return ierror.NewError(fmt.Errorf("comment is required for custom approval response"))
	}
	if !tool_approval.Manager.HasWaiter(input.ApprovalID) {
		return ierror.NewError(fmt.Errorf("approval %s is not active", input.ApprovalID))
	}

	approval, err := s.storage.GetToolApprovalByApprovalID(context.Background(), input.ApprovalID)
	if err != nil {
		return ierror.NewError(err)
	}
	if approval == nil {
		return ierror.NewError(fmt.Errorf("approval %s not found", input.ApprovalID))
	}
	if approval.Status != data_models.ToolApprovalStatusPending {
		return ierror.NewError(fmt.Errorf("approval %s is not pending", input.ApprovalID))
	}

	now := time.Now()
	approval.Status = data_models.ToolApprovalStatusResolved
	approval.Decision = input.Decision
	approval.ResponseComment = input.Comment
	approval.RespondedAt = &now
	if err := s.storage.SaveToolApproval(context.Background(), *approval); err != nil {
		return ierror.NewError(err)
	}

	if err := tool_approval.Manager.Resolve(tool_approval.WaitResult{
		ApprovalID:  input.ApprovalID,
		Decision:    input.Decision,
		Comment:     input.Comment,
		RespondedAt: now,
	}); err != nil {
		return ierror.NewError(err)
	}

	return nil
}
