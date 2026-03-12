package view_models

import (
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gorm.io/gorm"
)

type SupportProvider struct {
	ProviderType      data_models.ProviderType `json:"provider_type"`
	Icon              string                   `json:"icon"`
	Name              string                   `json:"name"`
	BaseUrl           string                   `json:"base_url"`
	FileUploadBaseUrl *string                  `json:"file_upload_base_url"`
	Description       string                   `json:"description"`
}

type Provider struct {
	ID                uint                     `json:"id"`
	CreatedAt         time.Time                `json:"created_at"`
	UpdatedAt         time.Time                `json:"updated_at"`
	DeletedAt         gorm.DeletedAt           `json:"deleted_at"`
	ProviderName      string                   `json:"provider_name"`        // 提供方名称
	BaseUrl           string                   `json:"base_url"`             // 基础url
	FileUploadBaseUrl *string                  `json:"file_upload_base_url"` // 文件上传url
	ApiKey            string                   `json:"api_key"`              // api key
	Enable            bool                     `json:"enable"`               // 启用
	DefaultModelId    *uint                    `json:"default_model_id"`
	Models            []Model                  `json:"models"`
	CustomModels      []data_models.Model      `json:"custom_models"`
	ProviderType      data_models.ProviderType `json:"provider_type"`
}
