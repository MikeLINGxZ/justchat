package storage

import (
	"context"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/examples/memory_demo_01/internal/models"
)

func (s *Storage) WriterMemory(ctx context.Context, title, content string, date *time.Time) error {
	memory := models.Memory{
		Title:        title,
		Content:      content,
		DateOccurred: date,
	}

	result := s.sqliteDb.WithContext(ctx).Create(&memory)
	return result.Error
}

func (s *Storage) ReadMemory(ctx context.Context, keyword string, startAt, endAt *time.Time) ([]models.Memory, error) {
	var memories []models.Memory

	if keyword == "" {
		err := s.sqliteDb.WithContext(ctx).Find(&memories).Error
		if err != nil {
			return nil, err
		}
		return memories, nil
	}

	// 开始构建查询
	db := s.sqliteDb.WithContext(ctx).Model(&models.Memory{})

	// 模糊匹配 Title 或 Content
	likeKeyword := "%" + keyword + "%"
	db = db.Where("title LIKE ? OR content LIKE ?", likeKeyword, likeKeyword)

	// 时间范围过滤：DateOccurred 在 [startAt, endAt] 之间
	if startAt != nil {
		db = db.Where("date_occurred >= ?", *startAt)
	}
	if endAt != nil {
		db = db.Where("date_occurred <= ?", *endAt)
	}

	// 按发生日期倒序排列（最新的记忆在前），若无日期则按创建时间
	db = db.Order("date_occurred DESC NULLS LAST, created_at DESC")

	// 执行查询
	err := db.Find(&memories).Error
	return memories, err
}
