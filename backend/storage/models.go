package storage

import (
	"context"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"
)

// GetModels 获取所有模型
func (s *Storage) GetModels(ctx context.Context) ([]data_models.Model, error) {
	var res []data_models.Model
	err := s.sqliteDB.Model(&data_models.Model{}).Find(&res).Error
	if err != nil {
		logger.Errorf("failed to get models: %s", err.Error())
		return nil, err
	}
	return res, nil
}
