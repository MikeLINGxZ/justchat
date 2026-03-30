package storage

import (
	"context"
	"errors"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gorm.io/gorm"
)

func (s *Storage) CreateToolApproval(ctx context.Context, approval data_models.ToolApproval) error {
	return s.sqliteDB.WithContext(ctx).Create(&approval).Error
}

func (s *Storage) SaveToolApproval(ctx context.Context, approval data_models.ToolApproval) error {
	var existing data_models.ToolApproval
	err := s.sqliteDB.WithContext(ctx).Where("approval_id = ?", approval.ApprovalID).First(&existing).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.sqliteDB.WithContext(ctx).Create(&approval).Error
		}
		return err
	}
	return s.sqliteDB.WithContext(ctx).Where("approval_id = ?", approval.ApprovalID).Updates(&approval).Error
}

func (s *Storage) GetToolApprovalByApprovalID(ctx context.Context, approvalID string) (*data_models.ToolApproval, error) {
	var approval data_models.ToolApproval
	err := s.sqliteDB.WithContext(ctx).Where("approval_id = ?", approvalID).First(&approval).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &approval, nil
}

func (s *Storage) ListToolApprovalsByTaskUUID(ctx context.Context, taskUUID string) ([]data_models.ToolApproval, error) {
	var approvals []data_models.ToolApproval
	err := s.sqliteDB.WithContext(ctx).
		Where("task_uuid = ?", taskUUID).
		Order("created_at ASC").
		Find(&approvals).Error
	if err != nil {
		return nil, err
	}
	return approvals, nil
}
