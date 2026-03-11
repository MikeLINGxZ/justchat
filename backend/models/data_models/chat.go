package data_models

type Chat struct {
	OrmModel
	Uuid         string `grom:"unique;index" json:"uuid"`
	Title        string `gorm:"type:varchar(255)" json:"title"`
	Prompt       string `gorm:"type:text" json:"prompt"`
	IsCollection bool   `gorm:"index;default:false" json:"is_collection"`
}
