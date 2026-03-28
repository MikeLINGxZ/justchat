package storage

import (
	"context"
	"errors"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gorm.io/gorm"
)

func (s *Storage) ListCustomMCPServers(ctx context.Context) ([]data_models.CustomMCPServer, error) {
	var res []data_models.CustomMCPServer
	err := s.sqliteDB.WithContext(ctx).Model(&data_models.CustomMCPServer{}).Order("id asc").Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Storage) GetCustomMCPServerBySourcePath(ctx context.Context, sourcePath string) (*data_models.CustomMCPServer, error) {
	var res data_models.CustomMCPServer
	err := s.sqliteDB.WithContext(ctx).Model(&data_models.CustomMCPServer{}).Where("source_path = ?", sourcePath).First(&res).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &res, nil
}

func (s *Storage) GetCustomMCPServerBySourcePathUnscoped(ctx context.Context, sourcePath string) (*data_models.CustomMCPServer, error) {
	var res data_models.CustomMCPServer
	err := s.sqliteDB.WithContext(ctx).Unscoped().Model(&data_models.CustomMCPServer{}).Where("source_path = ?", sourcePath).First(&res).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &res, nil
}

func (s *Storage) GetCustomMCPServerByToolID(ctx context.Context, toolID string) (*data_models.CustomMCPServer, error) {
	var res data_models.CustomMCPServer
	err := s.sqliteDB.WithContext(ctx).Model(&data_models.CustomMCPServer{}).Where("tool_id = ?", toolID).First(&res).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &res, nil
}

func (s *Storage) SaveCustomMCPServer(ctx context.Context, server data_models.CustomMCPServer) error {
	return s.sqliteDB.WithContext(ctx).Save(&server).Error
}

func (s *Storage) DeleteCustomMCPServerByToolID(ctx context.Context, toolID string) error {
	return s.sqliteDB.WithContext(ctx).Unscoped().Where("tool_id = ?", toolID).Delete(&data_models.CustomMCPServer{}).Error
}
