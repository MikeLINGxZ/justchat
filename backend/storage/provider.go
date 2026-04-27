package storage

import (
	"context"
	"errors"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gorm.io/gorm"
)

// GetProviders 获取所有供应商
func (s *Storage) GetProviders(ctx context.Context) ([]data_models.Provider, error) {
	var res []data_models.Provider
	err := s.sqliteDB.Model(&data_models.Provider{}).Find(&res).Error
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetProviderByID 获取指定 ID 的供应商
func (s *Storage) GetProviderByID(ctx context.Context, id uint) (*data_models.Provider, error) {
	var res data_models.Provider
	err := s.sqliteDB.Model(&data_models.Provider{}).Where("id = ?", id).First(&res).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &res, nil
}

// AddProvider 添加供应商
func (s *Storage) AddProvider(ctx context.Context, provider data_models.Provider) (uint, error) {
	return provider.ID, s.sqliteDB.Create(&provider).Error
}

// UpdateProvider 更新供应商
func (s *Storage) UpdateProvider(ctx context.Context, provider data_models.Provider) error {
	return s.sqliteDB.Save(&provider).Error
}

// ClearDefaultProviders 清空默认供应商标记
func (s *Storage) ClearDefaultProviders(ctx context.Context) error {
	return s.sqliteDB.Model(&data_models.Provider{}).Where("is_default = ?", true).Update("is_default", false).Error
}

// SetProviderDefault 设置指定供应商为默认供应商
func (s *Storage) SetProviderDefault(ctx context.Context, providerId uint) error {
	return s.sqliteDB.Model(&data_models.Provider{}).Where("id = ?", providerId).Update("is_default", true).Error
}

// DeleteProvider 删除供应商
func (s *Storage) DeleteProvider(ctx context.Context, id uint) error {
	return s.sqliteDB.Where("id = ?", id).Delete(&data_models.Provider{}).Error
}

// UpdateProviderDefaultModel 更新供应商默认模型
func (s *Storage) UpdateProviderDefaultModel(ctx context.Context, providerId uint, modelId uint) error {
	var count int64
	err := s.sqliteDB.Model(&data_models.ProviderDefaultModel{}).
		Where("provider_id = ?", providerId).
		Count(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		// 更新
		return s.sqliteDB.Model(&data_models.ProviderDefaultModel{}).
			Where("provider_id = ?", providerId).
			Update("model_id", modelId).Error
	} else {
		// 插入
		return s.sqliteDB.Create(&data_models.ProviderDefaultModel{
			ProviderID: providerId,
			ModelId:    modelId,
		}).Error
	}
}
