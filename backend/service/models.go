package service

import (
	"context"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils/ierror"
)

// GetModels 获取所有模型
func (s *Service) GetModels() ([]view_models.Model, error) {
	models, err := s.storage.GetModels(context.Background())
	if err != nil {
		return nil, ierror.NewError(err)
	}
	res := make([]view_models.Model, len(models))
	for i, model := range models {
		res[i] = model
	}

	return res, nil
}
