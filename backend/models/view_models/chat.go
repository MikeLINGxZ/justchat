package view_models

import "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"

type Chat struct {
	data_models.Chat
	Content          []MatchMessage `json:"content"`
	ReasoningContent []MatchMessage `json:"reasoning_content"`
}

type MatchMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
