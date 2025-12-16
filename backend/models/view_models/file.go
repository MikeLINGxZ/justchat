package view_models

import (
	"github.com/cloudwego/eino/schema"
)

type File struct {
	PreviewImg          *string                    `json:"preview"`
	Name                string                     `json:"name"`
	FilePath            string                     `json:"file_path"`
	MineType            string                     `json:"mine_type"`
	ChatMessagePartType schema.ChatMessagePartType `json:"chat_message_part_type"`
	Size                int64                      `json:"size"`
}
