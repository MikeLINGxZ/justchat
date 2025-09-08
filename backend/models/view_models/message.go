package view_models

import "github.com/cloudwego/eino/schema"

type MessageList struct {
	Messages []schema.Message `json:"messages"`
	Total    int              `json:"total"`
}
