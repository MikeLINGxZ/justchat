package service

import (
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/vo"
)

func (s *Service) GetModels() ([]vo.Model, error) {
	models, err := s.storage.GetModels(s.ctx)
	if err != nil {
		return nil, err
	}
	res := make([]vo.Model, len(models))
	for i, model := range models {
		res[i] = model
	}
	return res, nil
}
