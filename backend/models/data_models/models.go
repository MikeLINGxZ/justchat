package data_models

type Model struct {
	OrmModel
	ProviderId uint    `gorm:"index" json:"provider_id"` // 提供方id
	Model      string  `gorm:"index" json:"model"`
	OwnedBy    string  `gorm:"type:varchar(255)" json:"owned_by"`
	Object     string  `gorm:"type:varchar(255)" json:"object"`
	Enable     bool    `gorm:"index;type:bool;default:1" json:"enable"`
	Alias      *string `gorm:"index;type:varchar(255)" json:"alias"`
}
