package models

import (
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gorm.io/gorm"
)

type MemoryType string

const (
	MemoryTypeNone  MemoryType = ""      // 无分类记忆
	MemoryTypeSkill MemoryType = "skill" // 技能类记忆
	MemoryTypeEvent MemoryType = "event" // 事件类记忆
	MemoryTypeFlow  MemoryType = "flow"  // 流程类记忆
	MemoryTypePlan  MemoryType = "plan " // 流程类记忆
)

type Memory struct {
	data_models.OrmModel
	Summary       string     `gorm:"type:varchar(500)"`      // 自动摘要（便于检索）
	Content       string     `gorm:"type:text;not null"`     // 记忆内容（原始文本）
	Type          MemoryType `gorm:"type:varchar(50);index"` // 记忆类型
	TimeRangStart *time.Time `gorm:"index"`                  // 发生时间（可选）
	TimeRangeEnd  *time.Time `gorm:"index"`                  // 结束时间（可选）
	Location      *string    `gorm:"type:varchar(500)"`      // 发生地点（可选），多个地点使用","分割
	Characters    *string    `gorm:"type:varchar(500)"`      // 相关人物（可选），多个人物使用","分割
	Context       *string    `gorm:"type:json"`              // 上下文元数据（JSON 格式）

	EmbeddingID      *uint      `gorm:"index"`               // 嵌入id
	Importance       float64    `gorm:"default:0.5"`         // 重要性评分 [0.0 ~ 1.0]
	EmotionalValence float64    `gorm:"default:0.0"`         // 情感极性 [-1.0 ~ +1.0] 负面到正面
	IsForgotten      bool       `gorm:"default:false;index"` // 是否已遗忘（用于模拟遗忘）
	RecallCount      int        `gorm:"default:0"`           // 被回忆的次数（影响强度）
	LastRecalledAt   *time.Time `gorm:"-"`                   // 最后一次被回忆的时间
}

func (m *Memory) Fts() string {
	return `CREATE VIRTUAL TABLE IF NOT EXISTS memory_fts USING fts5(
             summary,
             content,
             location,
             characters,
             content='memories',           -- 关联到实际的数据表
             content_rowid='id'            -- 使用原表的 id 作为 rowid 映射
          );`
}

// AfterCreate 创建后更新 FTS
func (m *Memory) AfterCreate(tx *gorm.DB) error {
	return m.updateMemoryFTS(tx)
}

// AfterUpdate 更新后更新 FTS
func (m *Memory) AfterUpdate(tx *gorm.DB) error {
	return m.updateMemoryFTS(tx)
}

// AfterDelete 删除后清理 FTS
func (m *Memory) AfterDelete(tx *gorm.DB) error {
	return tx.Exec("DELETE FROM memory_fts WHERE rowid = ?", m.ID).Error
}

// updateMemoryFTS 插入或替换 FTS 记录
func (m *Memory) updateMemoryFTS(tx *gorm.DB) error {
	return tx.Exec(`
        INSERT OR REPLACE INTO memory_fts(rowid, summary, content)
        VALUES (?, ?, ?)
    `, m.ID, m.Summary, m.Content).Error
}
