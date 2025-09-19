package storage

import (
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/agents/memory/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Storage struct {
	sqliteDb *gorm.DB
}

func NewStorage() (*Storage, error) {
	db, err := gorm.Open(sqlite.Open("./memory_demo.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&models.Memory{})
	if err != nil {
		return nil, err
	}

	err = db.Exec((&models.Memory{}).Fts()).Error
	if err != nil {
		return nil, err
	}

	return &Storage{sqliteDb: db}, nil
}
