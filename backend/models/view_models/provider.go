package view_models

import (
	"time"

	"gorm.io/gorm"
)

type Provider struct {
	ID             uint           `json:"id"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `json:"deleted_at"`
	ProviderName   string         `json:"provider_name"` // 提供方名称
	BaseUrl        string         `json:"base_url"`      // 基础url
	ApiKey         string         `json:"api_key"`       // api key
	Enable         bool           `json:"enable"`        // 启用
	DefaultModelId *uint          `json:"default_model_id"`
	Models         []Model        `json:"models"`
}
