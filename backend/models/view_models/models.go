package view_models

import (
	"time"

	"gorm.io/gorm"
)

//type Model = data_models.Model

type Model struct {
	ID         uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	ProviderId uint           `gorm:"index" json:"provider_id"` // 提供方id
	Model      string         `gorm:"index" json:"model"`
	OwnedBy    string         `gorm:"type:varchar(255)" json:"owned_by"`
	Object     string         `gorm:"type:varchar(255)" json:"object"`
	Enable     bool           `gorm:"index;type:bool;default:1" json:"enable"`
	Alias      *string        `gorm:"index;type:varchar(255)" json:"alias"`
}
