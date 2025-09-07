package data_models

import (
	"encoding/json"

	"github.com/cloudwego/eino/schema"
	"gorm.io/gorm"
)

type Chat struct {
	OrmModel
	Uuid    string `grom:"unique;index" json:"uuid"`
	ModelID uint   `gorm:"index" json:"model_id"`
	Title   string `gorm:"type:varchar(255)" json:"title"`
	Prompt  string `gorm:"type:text" json:"prompt"`
}

type Message struct {
	OrmModel
	Uuid                       string          `grom:"unique;index" json:"uuid"`
	ChatUuid                   string          `gorm:"index" json:"chat_uuid"`
	SearchableContent          string          `gorm:"type:text" json:"search_content"`
	SearchableReasoningContent string          `gorm:"type:text" json:"search_reasoning_content"`
	MessageJson                string          `gorm:"type:text" json:"message_json"`
	Message                    *schema.Message `gorm:"-" json:"-"`
}

func (m *Message) BeforeCreate(tx *gorm.DB) (err error) {
	return m.before(tx)
}

func (m *Message) BeforeUpdate(tx *gorm.DB) (err error) {
	return m.before(tx)
}

func (m *Message) BeforeSave(tx *gorm.DB) (err error) {
	return m.before(tx)
}

func (m *Message) AfterFind(tx *gorm.DB) (err error) {
	var message schema.Message
	if m.MessageJson == "" {
		return
	}
	err = json.Unmarshal([]byte(m.MessageJson), &message)
	if err != nil {
		return err
	}
	m.Message = &message
	return
}

func (m *Message) before(tx *gorm.DB) (err error) {
	messageBytes, err := json.Marshal(m.Message)
	if err != nil {
		return err
	}
	m.MessageJson = string(messageBytes)
	m.SearchableContent = m.Message.Content
	m.SearchableReasoningContent = m.Message.ReasoningContent
	return
}
