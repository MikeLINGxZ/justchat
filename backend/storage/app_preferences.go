package storage

import (
	"context"
	"errors"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gorm.io/gorm"
)

func (s *Storage) GetAppPreferences(ctx context.Context) (*data_models.AppPreferences, error) {
	var prefs data_models.AppPreferences
	err := s.sqliteDB.WithContext(ctx).
		Model(&data_models.AppPreferences{}).
		Where("singleton_id = ?", data_models.AppPreferencesSingletonID).
		First(&prefs).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &data_models.AppPreferences{
			SingletonID:         data_models.AppPreferencesSingletonID,
			Language:            data_models.AppLanguageZhCN,
			Region:              data_models.AppRegionAsia,
			MemorySystemEnabled: false,
			VectorSearchEnabled: false,
			EmbeddingProvider:   "ollama",
			EmbeddingBaseURL:    "http://localhost:11434",
			EmbeddingAPIKey:     "",
			EmbeddingModel:      "bge-m3",
		}, nil
	}
	if err != nil {
		return nil, err
	}
	return &prefs, nil
}

func (s *Storage) SaveAppPreferences(ctx context.Context, prefs data_models.AppPreferences) error {
	prefs.SingletonID = data_models.AppPreferencesSingletonID
	var existing data_models.AppPreferences
	err := s.sqliteDB.WithContext(ctx).
		Model(&data_models.AppPreferences{}).
		Where("singleton_id = ?", prefs.SingletonID).
		First(&existing).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return s.sqliteDB.WithContext(ctx).Create(&prefs).Error
	}
	if err != nil {
		return err
	}

	existing.Language = prefs.Language
	existing.Region = prefs.Region
	existing.MemorySystemEnabled = prefs.MemorySystemEnabled
	existing.VectorSearchEnabled = prefs.VectorSearchEnabled
	existing.EmbeddingProvider = prefs.EmbeddingProvider
	existing.EmbeddingBaseURL = prefs.EmbeddingBaseURL
	existing.EmbeddingAPIKey = prefs.EmbeddingAPIKey
	existing.EmbeddingModel = prefs.EmbeddingModel
	return s.sqliteDB.WithContext(ctx).Save(&existing).Error
}
