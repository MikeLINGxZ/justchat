package view_model

type Model struct {
	ID         uint    `json:"id"`
	ProviderId uint    `gorm:"index" json:"provider_id"`
	Model      string  `json:"model"`
	OwnedBy    string  `json:"owned_by"`
	Object     string  `json:"object"`
	Enable     bool    `json:"enable"`
	Alias      *string `json:"alias"`
	IsCustom   bool    `json:"is_custom"`
	IsDefault  bool    `json:"is_default"` // 是否是默认模型
}
