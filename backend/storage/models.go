package storage

import (
	"context"
	"errors"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/wrapper_models"
	"gorm.io/gorm"
)

// GetModels 获取所有模型
func (s *Storage) GetModels(ctx context.Context) ([]data_models.Model, error) {
	var res []data_models.Model
	err := s.sqliteDB.Model(&data_models.Model{}).Find(&res).Error
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Storage) GetModel(ctx context.Context, model string) (*data_models.Model, error) {
	var models data_models.Model
	err := s.sqliteDB.Model(&data_models.Model{}).Where("model = ?", model).First(&models).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &models, nil
}

func (s *Storage) GetProviderModel(ctx context.Context, model string) (*wrapper_models.ProviderModel, error) {
	modelInfo, err := s.GetModel(ctx, model)
	if err != nil {
		return nil, err
	}
	if modelInfo == nil {
		return nil, nil
	}

	provider, err := s.GetProviderByID(ctx, modelInfo.ProviderId)
	if err != nil {
		return nil, err
	}
	if provider == nil {
		return nil, nil
	}

	return &wrapper_models.ProviderModel{
		BaseUrl: provider.BaseUrl,
		ApiKey:  provider.ApiKey,
		Model:   modelInfo.Model,
		ModelId: modelInfo.ID,
	}, nil
}

// AddProviderModel 添加供应商模型
func (s *Storage) AddProviderModel(ctx context.Context, model data_models.Model) error {
	return s.sqliteDB.Create(&model).Error
}

// DeleteAllProviderModel 删除供应商模型
func (s *Storage) DeleteAllProviderModel(ctx context.Context, providerId uint) error {
	return s.sqliteDB.Model(&data_models.Model{}).Where("provider_id = ?", providerId).Delete(&data_models.Model{}).Error
}
