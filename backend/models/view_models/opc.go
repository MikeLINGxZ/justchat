package view_models

import (
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
)

// OPCPersonInput 创建/更新人员的输入
type OPCPersonInput struct {
	Uuid   string   `json:"uuid"`
	Name   string   `json:"name"`
	Role   string   `json:"role"`
	Prompt string   `json:"prompt"`
	Tools  []string `json:"tools"`
	Skills []string `json:"skills"`
	Avatar string   `json:"avatar"`
}

// OPCPersonView 人员视图（含最后消息摘要）
type OPCPersonView struct {
	data_models.OPCPerson
	LastMessage   *string    `json:"last_message"`
	LastMessageAt *time.Time `json:"last_message_at"`
}

// OPCGroupInput 创建/更新群组的输入
type OPCGroupInput struct {
	Uuid        string   `json:"uuid"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	MemberUuids []string `json:"member_uuids"`
}

// OPCGroupView 群组视图（含成员和最后消息摘要）
type OPCGroupView struct {
	data_models.OPCGroup
	Members       []OPCPersonView `json:"members"`
	LastMessage   *string         `json:"last_message"`
	LastMessageAt *time.Time      `json:"last_message_at"`
}

// OPCCompletionInput 人员私聊输入
type OPCCompletionInput struct {
	ChatUuid   string `json:"chat_uuid"`
	PersonUuid string `json:"person_uuid"`
	Content    string `json:"content"`
	ModelId    uint   `json:"model_id"`
	ModelName  string `json:"model_name"`
}

// OPCGroupCompletionInput 群聊输入
type OPCGroupCompletionInput struct {
	ChatUuid  string `json:"chat_uuid"`
	GroupUuid string `json:"group_uuid"`
	Content   string `json:"content"`
	ModelId   uint   `json:"model_id"`
	ModelName string `json:"model_name"`
}
