package service

import (
	"io"

	"github.com/cloudwego/eino/schema"
	"github.com/gofrs/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils/ierror"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils/llm"
)

func (s *Service) ChatList(offset, limit int, keyword *string) (*view_models.ChatList, error) {
	chats, total, err := s.storage.GetChats(s.ctx, offset, limit, keyword)
	if err != nil {
		return nil, ierror.NewError(err)
	}

	res := view_models.ChatList{
		Lists: chats,
		Total: total,
	}

	return &res, nil
}

func (s *Service) Completions(chatUuid, model string, message schema.Message) (string, error) {
	// uuid
	uv4, err := uuid.NewV4()
	if err != nil {
		return "", ierror.NewError(err)
	}

	// 获取模型信息
	providerModel, err := s.storage.GetProviderModel(s.ctx, model)
	if err != nil {
		return "", err
	}
	if providerModel == nil {
		return "", ierror.New(ierror.ErrCodeModelNotFound)
	}

	// 当chatUuid为空说明是新建聊天
	if chatUuid == "" {
		chatUuid = uv4.String()
		// 创建一个聊天
		err = s.storage.CreateChat(s.ctx, chatUuid, message.Content, providerModel.ModelId)
		if err != nil {
			return "", ierror.NewError(err)
		}
	}

	// 新建一个消息id
	messageUuid := uv4.String()

	provider := llm.NewLlmProvider(providerModel.BaseUrl, providerModel.ApiKey, providerModel.Model)
	stream, err := provider.Completions(s.ctx, []schema.Message{message})
	if err != nil {
		return "", ierror.NewError(err)
	}

	msgChan := make(chan *schema.Message)
	errChan := make(chan error)
	doneChan := make(chan struct{})

	go func() {
		defer close(msgChan)
		defer close(errChan)
		defer close(doneChan)
		for {
			message, err := stream.Recv()
			if err == io.EOF {
				doneChan <- struct{}{}
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
		dataMsg := data_models.Message{
			Uuid:     messageUuid,
			ChatUuid: chatUuid,
		}
		for {
			select {
			case <-doneChan:
				dataMsg = s.fillCompletionsMsg(dataMsg, "done")
				runtime.EventsEmit(s.ctx, chatUuid, dataMsg)
				return
			case msg, ok := <-msgChan:
				if !ok || msg == nil {
					dataMsg = s.fillCompletionsMsg(dataMsg, "done")
					runtime.EventsEmit(s.ctx, chatUuid, dataMsg)
					return
				}
				dataMsg.Message = msg
				dataMsg = s.fillCompletionsMsg(dataMsg, "")
				runtime.EventsEmit(s.ctx, chatUuid, dataMsg)
			case err := <-errChan:
				if err == nil {
					continue
				}
				s.fillCompletionsMsg(dataMsg, err.Error())
				return
			}
		}
	}()

	return chatUuid, nil
}

func (s *Service) fillCompletionsMsg(dataMsg data_models.Message, finishReason string) data_models.Message {
	if finishReason != "" {
		dataMsg.Message = &schema.Message{
			Role: "assistant",
			ResponseMeta: &schema.ResponseMeta{
				FinishReason: finishReason,
			},
		}
	}
	err := s.storage.SaveOrUpdateDeltaMessage(s.ctx, dataMsg)
	if err != nil {
		logger.Errorf("save or update delta message failed: %v", err)
	}
	return dataMsg
}
