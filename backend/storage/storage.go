package storage

import (
	"context"
	"path/filepath"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Storage struct {
	sqliteDB *gorm.DB
}

func NewStorage() (*Storage, error) {
	dbPath, err := getDbPath()
	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		logger.Errorf("Failed to connect to database: %v", err)
		return nil, err
	}

	err = db.AutoMigrate(&data_models.Model{}, &data_models.Provider{}, &data_models.Chat{}, &data_models.Message{}, &data_models.ProviderDefaultModel{})
	if err != nil {
		logger.Errorf("Failed to migrate models: %v", err)
		return nil, err
	}

	// 创建 Storage 实例
	storage := &Storage{
		sqliteDB: db,
	}

	return storage, nil
}

// NewFnTransaction 开启事物
func (s *Storage) NewFnTransaction(ctx context.Context, fn func(ctx context.Context, s *Storage) error) error {
	tx := s.sqliteDB.Begin()
	txStorage := &Storage{
		sqliteDB: tx,
	}
	err := fn(ctx, txStorage)
	if err != nil {
		tx.Rollback()
		return err
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func getDbPath() (string, error) {
	dbName := "data.db"
	dataPath, err := utils.GetDataPath()
	if err != nil {
		return "", err
	}
	// 构建数据库路径
	dbPath := filepath.Join(dataPath, dbName)
	return dbPath, nil
}
