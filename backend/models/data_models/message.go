package data_models

import (
	"encoding/json"

	"github.com/cloudwego/eino/schema"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils/ierror"
	"gorm.io/gorm"
)

type Message struct {
	OrmModel
	ChatUuid                     string                 `gorm:"index" json:"chat_uuid"`
	MessageUuid                  string                 `grom:"unique;index" json:"message_uuid"`
	Role                         schema.RoleType        `gorm:"index" json:"role"`
	Content                      string                 `gorm:"type:text" json:"content"`
	ReasoningContent             string                 `gorm:"type:text" json:"reasoning_content"`
	UserMessageExtraContent      string                 `gorm:"type:text" json:"user_message_extra_content"`
	AssistantMessageExtraContent string                 `gorm:"type:text" json:"assistant_message_extra_content"`
	UserMessageExtra             *UserMessageExtra      `gorm:"-" json:"user_message_extra"`
	AssistantMessageExtra        *AssistantMessageExtra `gorm:"-" json:"assistant_message_extra"`
}

func (m *Message) ToSchemaMessage() (*schema.Message, error) {
	schemaMessage := &schema.Message{
		Role:             m.Role,
		Content:          m.Content,
		ReasoningContent: m.ReasoningContent,
	}
	// MultiContent
	if m.Role == schema.User && m.UserMessageExtra != nil && len(m.UserMessageExtra.Files) > 0 {
		schemaMessage.Content = ""
		schemaMessage.ReasoningContent = ""
		var userInputMultiContent []schema.MessageInputPart
		if m.Content != "" {
			userInputMultiContent = append(userInputMultiContent, schema.MessageInputPart{
				Type: schema.ChatMessagePartTypeText,
				Text: m.Content,
			})
		}

		for _, item := range m.UserMessageExtra.Files {
			base64Data, err := utils.ReadFile2Base64Data(item.Path)
			if err != nil {
				return nil, err
			}
			// 通过mineType获取消息类型
			chatMessagePartType, err := utils.MimeType2ChatMessagePartType(item.MineType)
			if err != nil {
				return nil, ierror.NewError(err)
			}

			var text string
			var img *schema.MessageInputImage
			var audio *schema.MessageInputAudio
			var video *schema.MessageInputVideo
			var file *schema.MessageInputFile
			messagePartCommon := schema.MessagePartCommon{
				Base64Data: &base64Data,
				MIMEType:   item.MineType,
				Extra: map[string]interface{}{
					"name":                   item.Name,
					"path":                   item.Path,
					"mime_type":              item.MineType,
					"chat_message_part_type": chatMessagePartType,
					"size":                   item.Size,
				},
			}
			switch chatMessagePartType {
			case schema.ChatMessagePartTypeText, schema.ChatMessagePartTypeFileURL:
				continue
			case schema.ChatMessagePartTypeImageURL:
				img = &schema.MessageInputImage{
					MessagePartCommon: messagePartCommon,
					Detail:            schema.ImageURLDetailHigh,
				}
			case schema.ChatMessagePartTypeAudioURL:
				audio = &schema.MessageInputAudio{
					MessagePartCommon: messagePartCommon,
				}
			case schema.ChatMessagePartTypeVideoURL:
				video = &schema.MessageInputVideo{
					MessagePartCommon: messagePartCommon,
				}
			}
			if img == nil && audio == nil && video == nil {
				continue
			}
			userInputMultiContent = append(userInputMultiContent, schema.MessageInputPart{
				Type:  chatMessagePartType,
				Text:  text,
				Image: img,
				Audio: audio,
				Video: video,
				File:  file,
			})
		}
		schemaMessage.UserInputMultiContent = userInputMultiContent
	}

	return schemaMessage, nil
}

type UserMessageExtra struct {
	ModelId   uint     `json:"model_id"`   // 模型id
	ModelName string   `json:"model_name"` // 模型名称
	Files     []File   `json:"files"`      // 文件路径
	Tools     []string `json:"tools"`      // 工具id
}

type AssistantMessageExtra struct {
	FinishReason string `json:"finish_reason"`
	FinishError  string `json:"finish_error"`
}

type File struct {
	Name     string  `json:"name"`
	Path     string  `json:"path"`
	Preview  *string `json:"preview"` // 如果是图像的话，生成60x60的预览图
	MineType string  `json:"mine_type"`
	Size     int64   `json:"size"`
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

	if m.UserMessageExtraContent != "" {
		var userMessageExtra UserMessageExtra
		err = json.Unmarshal([]byte(m.UserMessageExtraContent), &userMessageExtra)
		if err != nil {
			return err
		}
		m.UserMessageExtra = &userMessageExtra
	}

	if m.AssistantMessageExtraContent != "" {
		var assistantMessageExtra AssistantMessageExtra
		err = json.Unmarshal([]byte(m.AssistantMessageExtraContent), &assistantMessageExtra)
		if err != nil {
			return err
		}
		m.AssistantMessageExtra = &assistantMessageExtra
	}

	return
}

func (m *Message) before(tx *gorm.DB) (err error) {
	if m.UserMessageExtra != nil {
		bytes, err := json.Marshal(m.UserMessageExtra)
		if err != nil {
			return err
		}
		m.UserMessageExtraContent = string(bytes)
	}
	if m.AssistantMessageExtra != nil {
		bytes, err := json.Marshal(m.AssistantMessageExtra)
		if err != nil {
			return err
		}
		m.AssistantMessageExtraContent = string(bytes)
	}
	return
}
