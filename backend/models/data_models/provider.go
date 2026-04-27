package data_models

type Provider struct {
	OrmModel
	ProviderName      string       `gorm:"type:varchar(255)" json:"provider_name"`        // 提供方名称
	ProviderType      ProviderType `gorm:"type:varchar(255)" json:"provider_type"`        // 提供方类型：aliyuns、deepseek...
	BaseUrl           string       `gorm:"type:varchar(255)" json:"base_url"`             // 基础url
	FileUploadBaseUrl *string      `gorm:"type:varchar(255)" json:"file_upload_base_url"` // 文件上传地址
	ApiKey            string       `gorm:"type:varchar(255)" json:"api_key"`              // api key
	Enable            bool         `gorm:"index;type:bool;default:1" json:"enable"`       // 启用
	IsDefault         bool         `gorm:"index;type:bool;default:0" json:"is_default"`   // 默认供应商
}

type ProviderDefaultModel struct {
	OrmModel
	ProviderID uint `gorm:"index:idx_provider_model" json:"provider_id"`
	ModelId    uint `gorm:"index:idx_provider_model" json:"model_id"`
}

type ProviderType string

const (
	ProviderTypeDeepseek   ProviderType = "deepseek"
	ProviderTypeAliyuns    ProviderType = "aliyuns"
	ProviderTypeOpenrouter ProviderType = "openrouter"
	ProviderTypeOllama     ProviderType = "ollama"
	ProviderTypeOther      ProviderType = "other"
)

func (p ProviderType) String() string {
	return string(p)
}
