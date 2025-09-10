package service

import (
	"io"

	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils/ierror"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils/llm"
)

// ChatList 聊天列表
func (s *Service) ChatList(offset, limit int, keyword *string, isCollection bool) (*view_models.ChatList, error) {
	chats, total, err := s.storage.GetChats(s.ctx, offset, limit, keyword, isCollection)
	if err != nil {
		return nil, ierror.NewError(err)
	}

	return &view_models.ChatList{
		Lists: chats,
		Total: total,
	}, nil
}

// ChatMessages 聊天消息
func (s *Service) ChatMessages(chatUuid string, offset, limit int) (*view_models.MessageList, error) {
	dataMessages, total, err := s.storage.GetMessage(s.ctx, chatUuid, offset, limit)
	if err != nil {
		return nil, ierror.NewError(err)
	}

	var messages []schema.Message
	for _, item := range dataMessages {
		if item.Message == nil {
			continue
		}
		messages = append(messages, *item.Message)
	}

	return &view_models.MessageList{
		Messages: messages,
		Total:    total,
	}, nil
}

// Completions 聊天
func (s *Service) Completions(chatUuid, model string, message schema.Message) (*view_models.Completions, error) {

	// 获取模型信息
	providerModel, err := s.storage.GetProviderModel(s.ctx, model)
	if err != nil {
		return nil, ierror.NewError(err)
	}
	if providerModel == nil {
		return nil, ierror.New(ierror.ErrCodeModelNotFound)
	}

	// 当chatUuid为空说明是新建聊天
	if chatUuid == "" {
		chatUuid = uuid.New().String()
		// 创建一个聊天
		err = s.storage.CreateChat(s.ctx, chatUuid, message.Content, providerModel.ModelId)
		if err != nil {
			return nil, ierror.NewError(err)
		}
	}

	// 查找历史消息
	historyMessageData, _, err := s.storage.GetMessage(s.ctx, chatUuid, 0, 10)
	if err != nil {
		return nil, ierror.NewError(err)
	}
	var historyMessages []schema.Message
	for _, item := range historyMessageData {
		if item.Message == nil {
			continue
		}
		historyMessages = append(historyMessages, *item.Message)
	}

	// 创建用户消息
	err = s.storage.CreateMessage(s.ctx, chatUuid, data_models.Message{
		Uuid:     uuid.New().String(),
		ChatUuid: chatUuid,
		Message:  &message,
	})
	if err != nil {
		return nil, ierror.NewError(err)
	}

	provider := llm.NewLlmProvider(providerModel.BaseUrl, providerModel.ApiKey, providerModel.Model)
	stream, err := provider.Completions(s.ctx, append(historyMessages, message))
	if err != nil {
		return nil, ierror.NewError(err)
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
	messageUuid := uuid.New().String()
	go func() {
		dataModelMsg := data_models.Message{
			Uuid:     messageUuid,
			ChatUuid: chatUuid,
		}
		for {
			select {
			case <-doneChan:
				dataModelMsg = s.fillCompletionsMsg(dataModelMsg, "done")
				runtime.EventsEmit(s.ctx, messageUuid, dataModelMsg.Message)
				return
			case msg, ok := <-msgChan:
				if !ok || msg == nil {
					dataModelMsg = s.fillCompletionsMsg(dataModelMsg, "done")
					runtime.EventsEmit(s.ctx, messageUuid, dataModelMsg.Message)
					return
				}
				dataModelMsg.Message = msg
				dataModelMsg = s.fillCompletionsMsg(dataModelMsg, "")
				runtime.EventsEmit(s.ctx, messageUuid, dataModelMsg.Message)
				if msg.ResponseMeta != nil && msg.ResponseMeta.FinishReason != "" {
					return
				}
			case err := <-errChan:
				if err == nil {
					continue
				}
				s.fillCompletionsMsg(dataModelMsg, err.Error())
				runtime.EventsEmit(s.ctx, messageUuid, dataModelMsg.Message)
				return
			}
		}
	}()

	return &view_models.Completions{
		ChatUuid:    chatUuid,
		MessageUuid: messageUuid,
	}, nil
}

// DeleteChat 删除聊天
func (s *Service) DeleteChat(chatUuid string) error {
	err := s.storage.DeleteChat(s.ctx, chatUuid)
	if err != nil {
		return ierror.NewError(err)
	}
	return nil
}

// RenameChat 重命名聊天
func (s *Service) RenameChat(chatUuid, title string) error {
	err := s.storage.RenameChat(s.ctx, chatUuid, title)
	if err != nil {
		return ierror.NewError(err)
	}
	return nil
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
