package data_models

import (
	"encoding/json"

	"github.com/cloudwego/eino/schema"
)

type Chat struct {
	OrmModel
	ModelID uint   `gorm:"index" json:"model_id"`
	Title   string `gorm:"type:varchar(255)" json:"title"`
	Prompt  string `gorm:"type:text" json:"prompt"`
}

type Message struct {
	OrmModel

	// 关联字段
	ChatID           uint   `gorm:"index;not null" json:"chat_id"`
	ConversationUUID string `gorm:"index;size:64" json:"conversation_uuid,omitempty"` // 会话UUID，可选

	// schema.Message 的核心字段
	Role    string `gorm:"size:20;index;not null" json:"role"` // user, assistant, system, tool
	Content string `gorm:"type:text" json:"content"`           // 主要文本内容，支持全文搜索

	// 多模态内容（JSON存储）
	MultiContent string `gorm:"type:json" json:"multi_content,omitempty"` // schema.ChatMessagePart 数组的JSON

	// 扩展字段
	Name             string `gorm:"size:100" json:"name,omitempty"`               // 消息发送者名称
	ReasoningContent string `gorm:"type:text" json:"reasoning_content,omitempty"` // 模型推理过程，支持全文搜索

	// Assistant 消息专用字段
	ToolCalls string `gorm:"type:json" json:"tool_calls,omitempty"` // schema.ToolCall 数组的JSON

	// Tool 消息专用字段
	ToolCallID string `gorm:"size:100;index" json:"tool_call_id,omitempty"`
	ToolName   string `gorm:"size:100;index" json:"tool_name,omitempty"`

	// 响应元信息（JSON存储）
	ResponseMeta string `gorm:"type:json" json:"response_meta,omitempty"` // schema.ResponseMeta 的JSON

	// 扩展信息（JSON存储）
	Extra string `gorm:"type:json" json:"extra,omitempty"` // map[string]any 的JSON

	// 搜索优化字段
	SearchableContent string `gorm:"type:text;index" json:"-"` // 组合的可搜索内容（content + reasoning_content + name）

	// 业务扩展字段
	MessageUUID string `gorm:"uniqueIndex;size:64" json:"message_uuid"` // 消息唯一标识
	ParentID    *uint  `gorm:"index" json:"parent_id,omitempty"`        // 父消息ID（用于消息链）
	TokenCount  int    `gorm:"default:0" json:"token_count"`            // token数量
	Status      string `gorm:"size:20;default:'sent'" json:"status"`    // sent, failed, pending, streaming
	ErrorMsg    string `gorm:"type:text" json:"error_msg,omitempty"`    // 错误信息
	Metadata    string `gorm:"type:json" json:"metadata,omitempty"`     // 额外元数据
}

// updateSearchableContent 更新可搜索内容字段
func (m *Message) updateSearchableContent() {
	searchContent := ""
	if m.Content != "" {
		searchContent += m.Content + " "
	}
	if m.ReasoningContent != "" {
		searchContent += m.ReasoningContent + " "
	}
	if m.Name != "" {
		searchContent += m.Name + " "
	}
	m.SearchableContent = searchContent
}

// ToSchemaMessage 转换为 schema.Message
func (m *Message) ToSchemaMessage() (*schema.Message, error) {
	msg := &schema.Message{
		Role:             schema.RoleType(m.Role),
		Content:          m.Content,
		Name:             m.Name,
		ReasoningContent: m.ReasoningContent,
	}

	// 解析 MultiContent
	if m.MultiContent != "" {
		var multiContent []schema.ChatMessagePart
		if err := json.Unmarshal([]byte(m.MultiContent), &multiContent); err != nil {
			return nil, err
		}
		msg.MultiContent = multiContent
	}

	// 解析 ToolCalls
	if m.ToolCalls != "" {
		var toolCalls []schema.ToolCall
		if err := json.Unmarshal([]byte(m.ToolCalls), &toolCalls); err != nil {
			return nil, err
		}
		msg.ToolCalls = toolCalls
	}

	// 设置 Tool 消息字段
	msg.ToolCallID = m.ToolCallID
	msg.ToolName = m.ToolName

	// 解析 ResponseMeta
	if m.ResponseMeta != "" {
		var responseMeta schema.ResponseMeta
		if err := json.Unmarshal([]byte(m.ResponseMeta), &responseMeta); err != nil {
			return nil, err
		}
		msg.ResponseMeta = &responseMeta
	}

	// 解析 Extra
	if m.Extra != "" {
		var extra map[string]any
		if err := json.Unmarshal([]byte(m.Extra), &extra); err != nil {
			return nil, err
		}
		msg.Extra = extra
	}

	return msg, nil
}

// FromSchemaMessage 从 schema.Message 创建
func (m *Message) FromSchemaMessage(msg *schema.Message) error {
	m.Role = string(msg.Role)
	m.Content = msg.Content
	m.Name = msg.Name
	m.ReasoningContent = msg.ReasoningContent
	m.ToolCallID = msg.ToolCallID
	m.ToolName = msg.ToolName

	// 序列化 MultiContent
	if len(msg.MultiContent) > 0 {
		data, err := json.Marshal(msg.MultiContent)
		if err != nil {
			return err
		}
		m.MultiContent = string(data)
	}

	// 序列化 ToolCalls
	if len(msg.ToolCalls) > 0 {
		data, err := json.Marshal(msg.ToolCalls)
		if err != nil {
			return err
		}
		m.ToolCalls = string(data)
	}

	// 序列化 ResponseMeta
	if msg.ResponseMeta != nil {
		data, err := json.Marshal(msg.ResponseMeta)
		if err != nil {
			return err
		}
		m.ResponseMeta = string(data)
	}

	// 序列化 Extra
	if msg.Extra != nil {
		data, err := json.Marshal(msg.Extra)
		if err != nil {
			return err
		}
		m.Extra = string(data)
	}

	return nil
}
