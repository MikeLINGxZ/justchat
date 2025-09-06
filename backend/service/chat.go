package service

import (
	"github.com/cloudwego/eino/schema"
	"github.com/gofrs/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
)

func (s *Service) ChatList(offset, limit int, keyword *string) (*view_models.ChatList, error) {
	chats, total, err := s.storage.GetChats(s.ctx, offset, limit, keyword)
	if err != nil {
		return nil, err
	}

	res := view_models.ChatList{
		Lists: chats,
		Total: total,
	}

	return &res, nil
}

func (s *Service) Completions(chatUuid, model string, message schema.Message) (string, error) {
	if chatUuid == "" {
		uv4, err := uuid.NewV4()
		if err != nil {
			return "", err
		}
		chatUuid = uv4.String()
	}
	go func() {
		contents := []string{
			"你好",
			"，这里是模拟ai回复,",
			"这里是ai助手",
		}
		for i := 0; i < len(contents); i++ {
			finishReason := ""
			if i == len(contents)-1 {
				finishReason = "finish"
			}
			runtime.EventsEmit(s.ctx, chatUuid, schema.Message{
				Role:         "assistant",
				Content:      contents[i],
				MultiContent: nil,
				Name:         "",
				ToolCalls:    nil,
				ToolCallID:   "",
				ToolName:     "",
				ResponseMeta: &schema.ResponseMeta{
					FinishReason: finishReason,
				},
				ReasoningContent: "",
				Extra:            nil,
			})
		}
	}()
	return chatUuid, nil
}
