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
