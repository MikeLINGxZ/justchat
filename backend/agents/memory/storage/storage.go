package storage

import (
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/agents/memory/models"
	"gorm.io/gorm"
)

type Storage struct {
	sqliteDb *gorm.DB
}

// NewStorage 创建 Memory Storage，接受外部传入的 *gorm.DB 实例（主数据库连接）。
func NewStorage(db *gorm.DB) (*Storage, error) {
	err := db.AutoMigrate(&models.Memory{})
	if err != nil {
		return nil, err
	}
	return &Storage{sqliteDb: db}, nil
}

// DB 返回底层 *gorm.DB 实例，供 FTS5 初始化等场景使用。
func (s *Storage) DB() *gorm.DB {
	return s.sqliteDb
}
