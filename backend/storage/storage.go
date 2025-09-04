package storage

import (
	"fmt"

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

	err = db.AutoMigrate(&data_models.Model{}, &data_models.Provider{}, &data_models.Chat{}, &data_models.Message{})
	if err != nil {
		logger.Errorf("Failed to migrate models: %v", err)
		return nil, err
	}

	// 创建 Storage 实例
	storage := &Storage{
		sqliteDB: db,
	}

	// 创建 FTS 索引
	if err := storage.createFTSIndex(); err != nil {
		logger.Errorf("Failed to create FTS index: %v", err)
		return nil, err
	}

	return storage, nil
}

func (s *Storage) createFTSIndex() error {
	// 检查是否已存在FTS表
	var count int64
	err := s.sqliteDB.Raw(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='messages_fts'`).Scan(&count).Error
	if err != nil {
		return err
	}

	if count > 0 {
		return nil // FTS表已存在
	}

	// 创建FTS5虚拟表
	ftsSQL := `
	CREATE VIRTUAL TABLE messages_fts USING fts5(
		id UNINDEXED,
		chat_id UNINDEXED, 
		content,
		searchable_content,
		searchable_reasoning_content,
		created_at UNINDEXED
	);
	`

	if err := s.sqliteDB.Exec(ftsSQL).Error; err != nil {
		return fmt.Errorf("failed to create FTS table: %w", err)
	}

	// 创建触发器自动同步数据
	triggers := []string{
		// INSERT 触发器
		`CREATE TRIGGER messages_fts_insert AFTER INSERT ON messages BEGIN
			INSERT INTO messages_fts(id, chat_id, content, searchable_content, searchable_reasoning_content, created_at)
			VALUES (new.id, new.chat_id, new.message_json, new.searchable_content, new.searchable_reasoning_content, new.created_at);
		END;`,

		// UPDATE 触发器
		`CREATE TRIGGER messages_fts_update AFTER UPDATE ON messages BEGIN
			UPDATE messages_fts 
			SET content = new.message_json, 
				searchable_content = new.searchable_content,
				searchable_reasoning_content = new.searchable_reasoning_content
			WHERE id = new.id;
		END;`,

		// DELETE 触发器
		`CREATE TRIGGER messages_fts_delete AFTER DELETE ON messages BEGIN
			DELETE FROM messages_fts WHERE id = old.id;
		END;`,
	}

	for _, trigger := range triggers {
		if err := s.sqliteDB.Exec(trigger).Error; err != nil {
			return fmt.Errorf("failed to create trigger: %w", err)
		}
	}

	return nil
}
