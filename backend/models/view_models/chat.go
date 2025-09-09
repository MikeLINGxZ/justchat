package view_models

import "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"

type ChatList struct {
	Lists []Chat `json:"lists"`
	Total int    `json:"total"`
}

type Chat struct {
	data_models.Chat
	Content          []MatchMessage `json:"content"`
	ReasoningContent []MatchMessage `json:"reasoning_content"`
}

type MatchMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Completions struct {
	ChatUuid    string `json:"chat_uuid"`
	MessageUuid string `json:"message_uuid"`
}
