package provider_dto

import "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_model"

// EditProviderInput carries all mutable fields for an existing provider.
type EditProviderInput struct {
	ProviderId   int64              `json:"provider_id"`   // 供应商ID
	ProviderName string             `json:"provider_name"` // 供应商名称
	BaseUrl      string             `json:"base_url"`      // 供应商基础URL
	ApiKey       string             `json:"api_key"`       // 供应商API密钥
	Enable       bool               `json:"enable"`        // 是否启用
	DefaultModel *string            `json:"default_model"` // 默认模型名称，nil 表示不修改
	Models       []view_model.Model `json:"models"`        // 需要新增的模型（id=0 的条目）
}

type EditProviderOutput struct{}
