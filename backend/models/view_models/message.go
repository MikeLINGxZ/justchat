package view_models

import (
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
)

type MessageList struct {
	Messages []Message `json:"messages"`
	Total    int       `json:"total"`
}

type Message = data_models.Message

type FileInfo = data_models.File
