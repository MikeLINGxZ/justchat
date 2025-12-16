package view_models

import "github.com/cloudwego/eino/schema"

type MessageList struct {
	Messages []schema.Message `json:"messages"`
	Total    int              `json:"total"`
}

type MessagePkg struct {
	ChatUuid string `json:"chatUuid"`
	Model    string `json:"model"`
	Content  string `json:"content"`
	Files    []File `json:"files"`
}
