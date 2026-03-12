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

// GetModelsByProviderIds 通过供应商ids获取模型
func (s *Storage) GetModelsByProviderIds(ctx context.Context, providerIds []uint) (map[uint][]data_models.Model, error) {
	var res []data_models.Model
	err := s.sqliteDB.Model(&data_models.Model{}).Where("provider_id IN (?)", providerIds).Find(&res).Error
	if err != nil {
		return nil, err
	}
	resMap := make(map[uint][]data_models.Model)
	for _, model := range res {
		resMap[model.ProviderId] = append(resMap[model.ProviderId], model)
	}
	return resMap, nil
}

// GetProviderDefaultModelByProviderIds 通过供应商ids获取供应商默认模型
func (s *Storage) GetProviderDefaultModelByProviderIds(ctx context.Context, providerIds []uint) (map[uint]data_models.Model, error) {
	var providerModels []data_models.ProviderDefaultModel
	err := s.sqliteDB.Model(&data_models.ProviderDefaultModel{}).Where("provider_id IN (?)", providerIds).Find(&providerModels).Error
	if err != nil {
		return nil, err
	}
	resMap := make(map[uint]data_models.Model)

	// 遍历每个供应商默认模型配置
	for _, pm := range providerModels {
		// 根据模型ID获取完整的模型信息
		var model data_models.Model
		err := s.sqliteDB.Model(&data_models.Model{}).Where("id = ?", pm.ModelId).First(&model).Error
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// 如果模型不存在，跳过这个配置
				continue
			}
			return nil, err
		}
		// 将模型信息添加到结果映射中
		resMap[pm.ProviderID] = model
	}

	return resMap, nil
}

func (s *Storage) GetModel(ctx context.Context, model string) (*data_models.Model, error) {
	var models data_models.Model
	err := s.sqliteDB.Model(&data_models.Model{}).Joins("join providers p on p.id = provider_id AND p.deleted_at IS NULL").Where("model = ?", model).First(&models).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &models, nil
}

func (s *Storage) GetModelByIdAndName(ctx context.Context, modelId uint, modelName string) (*data_models.Model, error) {
	var models data_models.Model
	err := s.sqliteDB.Model(&data_models.Model{}).Joins("join providers p on p.id = provider_id AND p.deleted_at IS NULL").Where("models.id = ? AND models.model = ?", modelId, modelName).First(&models).Error
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &models, nil
}

func (s *Storage) GetProviderModel(ctx context.Context, modelId uint, modelName string) (*wrapper_models.ProviderModel, error) {
	modelInfo, err := s.GetModelByIdAndName(ctx, modelId, modelName)
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
		ProviderType:      provider.ProviderType,
		BaseUrl:           provider.BaseUrl,
		FileUploadBaseUrl: provider.FileUploadBaseUrl,
		ApiKey:            provider.ApiKey,
		Model:             modelInfo.Model,
		ModelId:           modelInfo.ID,
	}, nil
}

// AddProviderModel 添加供应商模型
func (s *Storage) AddProviderModel(ctx context.Context, model data_models.Model) error {
	return s.sqliteDB.Create(&model).Error
}

// DeleteAllProviderModel 删除供应商模型
func (s *Storage) DeleteAllProviderModel(ctx context.Context, providerId uint) error {
	return s.sqliteDB.Model(&data_models.Model{}).Where("provider_id = ? AND is_custom = 0", providerId).Delete(&data_models.Model{}).Error
}

// DeleteProviderModel 删除供应商模型
func (s *Storage) DeleteProviderModel(ctx context.Context, providerId uint, modelName string) error {
	return s.sqliteDB.Model(&data_models.Model{}).Where("provider_id = ? AND model = ?", providerId, modelName).Delete(&data_models.Model{}).Error
}
