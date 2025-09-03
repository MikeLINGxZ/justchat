package data_models

import "github.com/cloudwego/eino/schema"

type Chat struct {
	OrmModel
	ModelID uint   `gorm: "index" json:"model_id"`
	Title   string `gorm: "type:varchar(255)" json:"title"`
	Prompt  string `gorm: "type:text" json:"prompt"`
}

type Message struct {
	OrmModel
	ChatID uint 
	schema.Message
}
