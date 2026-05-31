package storage

import (
	"errors"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gorm.io/gorm"
)

// CreateProvider inserts a provider record and returns the saved row with its auto-generated ID populated.
func (s *Storage) CreateProvider(p data_models.Provider) (*data_models.Provider, error) {
	if err := s.sqliteDB.Create(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

// CreateModels bulk-inserts model records; IDs are populated on the pointed-to slice in place by GORM.
func (s *Storage) CreateModels(models *[]data_models.Model) error {
	if len(*models) == 0 {
		return nil
	}
	return s.sqliteDB.Create(models).Error
}

// ListProviders returns all provider rows (soft-delete aware).
func (s *Storage) ListProviders() ([]data_models.Provider, error) {
	var providers []data_models.Provider
	if err := s.sqliteDB.Find(&providers).Error; err != nil {
		return nil, err
	}
	return providers, nil
}

// ListModelsForProvider returns all model rows belonging to providerID.
func (s *Storage) ListModelsForProvider(providerID uint) ([]data_models.Model, error) {
	var models []data_models.Model
	if err := s.sqliteDB.Where("provider_id = ?", providerID).Find(&models).Error; err != nil {
		return nil, err
	}
	return models, nil
}

// UpsertDefaultModel creates or updates the default-model record for providerID.
func (s *Storage) UpsertDefaultModel(providerID, modelID uint) error {
	var row data_models.ProviderDefaultModel
	err := s.sqliteDB.Where("provider_id = ?", providerID).First(&row).Error
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}
		return s.sqliteDB.Create(&data_models.ProviderDefaultModel{
			ProviderID: providerID,
			ModelId:    modelID,
		}).Error
	}
	return s.sqliteDB.Model(&row).Update("model_id", modelID).Error
}

// GetDefaultModel returns the default model record for providerID, or nil if unset.
func (s *Storage) GetDefaultModel(providerID uint) (*data_models.ProviderDefaultModel, error) {
	var row data_models.ProviderDefaultModel
	err := s.sqliteDB.Where("provider_id = ?", providerID).First(&row).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &row, nil
}

// UpdateProvider updates the name, base URL, API key, and enable flag for the provider identified by p.ID.
func (s *Storage) UpdateProvider(p data_models.Provider) error {
	return s.sqliteDB.Model(&p).Updates(map[string]any{
		"provider_name": p.ProviderName,
		"base_url":      p.BaseUrl,
		"api_key":       p.ApiKey,
		"enable":        p.Enable,
	}).Error
}

// DeleteProvider soft-deletes the provider record identified by providerID.
func (s *Storage) DeleteProvider(providerID uint) error {
	return s.sqliteDB.Delete(&data_models.Provider{}, providerID).Error
}

// DeleteModelsForProvider soft-deletes all model records belonging to providerID.
func (s *Storage) DeleteModelsForProvider(providerID uint) error {
	return s.sqliteDB.Where("provider_id = ?", providerID).Delete(&data_models.Model{}).Error
}

// DeleteDefaultModelRecord soft-deletes the default-model record for providerID.
func (s *Storage) DeleteDefaultModelRecord(providerID uint) error {
	return s.sqliteDB.Where("provider_id = ?", providerID).Delete(&data_models.ProviderDefaultModel{}).Error
}

// DeleteModel soft-deletes a single model record by its ID.
func (s *Storage) DeleteModel(modelID uint) error {
	return s.sqliteDB.Delete(&data_models.Model{}, modelID).Error
}

// DeleteProviderCascade removes the provider, all its models, and the default-model record in a single transaction.
func (s *Storage) DeleteProviderCascade(providerID uint) error {
	return s.sqliteDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("provider_id = ?", providerID).Delete(&data_models.Model{}).Error; err != nil {
			return err
		}
		if err := tx.Where("provider_id = ?", providerID).Delete(&data_models.ProviderDefaultModel{}).Error; err != nil {
			return err
		}
		return tx.Delete(&data_models.Provider{}, providerID).Error
	})
}
