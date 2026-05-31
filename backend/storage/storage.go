package storage

import (
	"path/filepath"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/dir"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Storage struct {
	sqliteDB *gorm.DB
}

// NewStorage opens the application SQLite database and auto-migrates all models.
func NewStorage() (*Storage, error) {
	dataDir, err := dir.GetDataDir()
	if err != nil {
		return nil, err
	}

	dbPath := filepath.Join(dataDir, dir.DataBaseFileName)
	gormDB, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return NewStorageFromDB(gormDB)
}

// NewStorageFromDB wraps an existing gorm.DB handle after running auto-migration.
// Used in tests to inject an in-memory database.
func NewStorageFromDB(db *gorm.DB) (*Storage, error) {
	if err := db.AutoMigrate(
		&data_models.Provider{},
		&data_models.ProviderDefaultModel{},
		&data_models.Model{},
		&data_models.Session{},
		&data_models.TaskState{},
		&data_models.Notification{},
		&data_models.Message{},
		&data_models.Terminal{},
		&data_models.TerminalOutputChunk{},
		&data_models.Memory{},
		&data_models.MemoryEmbedding{},
	); err != nil {
		return nil, err
	}
	return &Storage{sqliteDB: db}, nil
}
