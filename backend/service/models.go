package service

import (
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
)

// GetModels 获取所有模型
func (s *Service) GetModels() ([]view_models.Model, error) {
	models, err := s.storage.GetModels(s.ctx)
	if err != nil {
		return nil, err
	}
	res := make([]view_models.Model, len(models))
	for i, model := range models {
		res[i] = model
	}

	return res, nil
}
