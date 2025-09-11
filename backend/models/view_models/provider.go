package view_models

import (
	"time"

	"gorm.io/gorm"
)

//type Provider = data_models.Provider

type Provider struct {
	ID           uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	ProviderName string         `gorm:"type:varchar(255)" json:"provider_name"`  // 提供方名称
	BaseUrl      string         `gorm:"type:varchar(255)" json:"base_url"`       // 基础url
	ApiKey       string         `gorm:"type:varchar(255)" json:"api_key"`        // api key
	Enable       bool           `gorm:"index;type:bool;default:1" json:"enable"` // 启用
	Alias        *string        `gorm:"type:varchar(255)" json:"alias"`          // 别名
}
