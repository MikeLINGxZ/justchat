package storage

import (
	"context"
	"path/filepath"

	memory_models "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/agents/memory/models"
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

	err = db.AutoMigrate(
		&data_models.AppPreferences{},
		&data_models.Model{},
		&data_models.Provider{},
		&data_models.Chat{},
		&data_models.Message{},
		&data_models.ProviderDefaultModel{},
		&data_models.Task{},
		&data_models.CustomMCPServer{},
		&data_models.ToolApproval{},
		&data_models.OPCPerson{},
		&data_models.OPCGroup{},
		&data_models.OPCGroupMember{},
		&memory_models.Memory{},
	)
	if err != nil {
		logger.Errorf("Failed to migrate models: %v", err)
		return nil, err
	}

	// 尝试初始化 FTS5 全文检索索引（用于记忆系统）。
	// 需要 SQLite 编译时启用 FTS5 模块（build tag: sqlite_fts5）。
	// 初始化失败时静默降级为 LIKE 检索，不阻塞启动。
	if ftsErr := initMemoryFTS(db); ftsErr != nil {
		logger.Warm("memory FTS5 not available, falling back to LIKE search:", ftsErr)
	}

	// 创建 Storage 实例
	storage := &Storage{
		sqliteDB: db,
	}

	return storage, nil
}

// DB 返回底层 *gorm.DB 实例。
func (s *Storage) DB() *gorm.DB {
	return s.sqliteDB
}

// initMemoryFTS 创建记忆系统的 FTS5 虚拟表和自动同步触发器。
func initMemoryFTS(db *gorm.DB) error {
	stmts := []string{
		// FTS5 虚拟表：对 summary 和 content 建立全文索引
		`CREATE VIRTUAL TABLE IF NOT EXISTS memories_fts USING fts5(
			summary,
			content,
			content='memories',
			content_rowid='id',
			tokenize='unicode61'
		)`,
		// 触发器：INSERT 时同步
		`CREATE TRIGGER IF NOT EXISTS memories_ai AFTER INSERT ON memories BEGIN
			INSERT INTO memories_fts(rowid, summary, content) VALUES (new.id, new.summary, new.content);
		END`,
		// 触发器：DELETE 时同步
		`CREATE TRIGGER IF NOT EXISTS memories_ad AFTER DELETE ON memories BEGIN
			INSERT INTO memories_fts(memories_fts, rowid, summary, content) VALUES('delete', old.id, old.summary, old.content);
		END`,
		// 触发器：UPDATE 时同步
		`CREATE TRIGGER IF NOT EXISTS memories_au AFTER UPDATE ON memories BEGIN
			INSERT INTO memories_fts(memories_fts, rowid, summary, content) VALUES('delete', old.id, old.summary, old.content);
			INSERT INTO memories_fts(rowid, summary, content) VALUES (new.id, new.summary, new.content);
		END`,
	}
	for _, stmt := range stmts {
		if err := db.Exec(stmt).Error; err != nil {
			return err
		}
	}

	// 重建索引：将已有记忆数据同步到 FTS5 索引中。
	// 使用 FTS5 rebuild 命令从 content table (memories) 重新读取全部数据。
	if err := db.Exec(`INSERT INTO memories_fts(memories_fts) VALUES('rebuild')`).Error; err != nil {
		logger.Warm("FTS5 rebuild failed (non-fatal):", err)
	}

	return nil
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
