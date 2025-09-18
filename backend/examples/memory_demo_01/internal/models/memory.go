package models

import (
	"time"

	"gorm.io/gorm"
)

type Memory struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Title        string         `gorm:"not null;size:255" json:"title"`    // 记忆标题
	Content      string         `gorm:"type:text" json:"content"`          // 内容
	DateOccurred *time.Time     `json:"date_occurred,omitempty"`           // 发生日期（可选）
	CreatedAt    time.Time      `json:"created_at"`                        // 创建时间
	UpdatedAt    time.Time      `json:"updated_at"`                        // 更新时间
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // 软删除支持
}
