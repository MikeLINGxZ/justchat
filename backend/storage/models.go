package storage

import (
	"context"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/do"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"
)

func (s *Storage) GetModels(ctx context.Context) ([]do.Model, error) {
	var res []do.Model
	err := s.sqliteDB.Model(&do.Model{}).Find(&res).Error
	if err != nil {
		logger.Errorf("failed to get models: %s", err.Error())
		return nil, err
	}
	return res, nil
}
