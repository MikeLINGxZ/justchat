package models

import (
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
)

type MemoryTarget string

const (
	MemoryTargetUser  MemoryTarget = "user"
	MemoryTargetAgent MemoryTarget = "memory"
)

type MemoryType string

// 新版三类型语义：
//
//	fact         不变的客观事实（用户出生地、家庭成员、过敏源等）
//	information  可变的偏好/状态/习惯（当前住所、常用设备、最近兴趣等）
//	event        带时间锚点的具体事件或计划
const (
	MemoryTypeNone  MemoryType = ""
	MemoryTypeFact  MemoryType = "fact"
	MemoryTypeInfo  MemoryType = "information"
	MemoryTypeEvent MemoryType = "event"
)

type Memory struct {
	data_models.OrmModel
	Summary string     `gorm:"type:varchar(500)"`      // 标题
	Content string     `gorm:"type:text;not null"`     // 内容（含时间/地点/人物等所有信息）
	Type    MemoryType `gorm:"type:varchar(50);index"` // 类型：fact / information / event

	Target    MemoryTarget `gorm:"type:varchar(32);default:user;index"` // user / memory
	Source    string       `gorm:"type:varchar(64);default:agent"`      // agent / manual / legacy
	CharCount int          `gorm:"default:0"`

	EmbeddingID *uint `gorm:"index"`               // 嵌入 id
	IsForgotten bool  `gorm:"default:false;index"` // 是否已遗忘
	RecallCount int   `gorm:"default:0"`           // 召回次数

	LastRecalledAt *time.Time `gorm:"column:last_recalled_at"`
	LastUsedAt     *time.Time `gorm:"column:last_used_at"`

	// === 已废弃字段（保留为 nullable，仅供历史迁移读取，不再写入新数据）===
	TimeRangStart    *time.Time `gorm:"index"`
	TimeRangeEnd     *time.Time `gorm:"index"`
	Location         *string    `gorm:"type:varchar(500)"`
	Characters       *string    `gorm:"type:varchar(500)"`
	Context          *string    `gorm:"type:json"`
	Importance       float64    `gorm:"default:0.5"`
	EmotionalValence float64    `gorm:"default:0.0"`
	TrustScore       float64    `gorm:"default:0.5"`
}
