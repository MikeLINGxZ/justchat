package storage

import (
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Storage struct {
	sqliteDB *gorm.DB
}

func NewStorage() (*Storage, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		logger.Errorf("Failed to connect to database: %v", err)
		return nil, err
	}

	err = db.AutoMigrate(&data_models.Model{}, &data_models.Provider{})
	if err != nil {
		logger.Errorf("Failed to migrate models: %v", err)
		return nil, err
	}

	return &Storage{sqliteDB: db}, nil
}
