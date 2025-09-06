package service

import (
	"errors"
	"io"

	"github.com/cloudwego/eino/schema"
	"github.com/gofrs/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils/llm"
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

	providerModel, err := s.storage.GetProviderModel(s.ctx, model)
	if err != nil {
		return "", err
	}
	if providerModel == nil {
		return "", errors.New("model not found")
	}

	provider := llm.NewLlmProvider(providerModel.BaseUrl, providerModel.ApiKey, providerModel.Model)
	stream, err := provider.Completions(s.ctx, []schema.Message{message})
	if err != nil {
		return "", err
	}

	//messageUuid := uuid.Must(uuid.NewV4()).String()
	msgChan := make(chan *schema.Message)
	errChan := make(chan error)
	var msgIndex int32

	go func() {
		defer close(msgChan)
		defer close(errChan)

		for {
			message, err := stream.Recv()
			if err == io.EOF { // 流式输出结束
				return
			}
			if err != nil && err != io.EOF {
				errChan <- err
				return
			}
			msgChan <- message
		}
	}()

	go func() {
		for {
			select {
			case msg, ok := <-msgChan:
				if !ok {
					// 流结束，发送最终的完成信号
					return
				}
				runtime.EventsEmit(s.ctx, chatUuid, msg)

				//case err := <-errChan:

			}
			msgIndex++

		}
	}()

	return chatUuid, nil
}
