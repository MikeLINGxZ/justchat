package service

import (
	"context"
	"fmt"
	"io"

	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/llm_provider"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils/ierror"
)

// ChatList 聊天列表
func (s *Service) ChatList(offset, limit int, keyword *string, isCollection bool) (*view_models.ChatList, error) {
	chats, total, err := s.storage.GetChats(context.Background(), offset, limit, keyword, isCollection)
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
	dataMessages, total, err := s.storage.GetMessage(context.Background(), chatUuid, offset, limit)
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
func (s *Service) Completions(message view_models.MessagePkg) (*view_models.Completions, error) {

	// 获取模型信息
	providerModel, err := s.storage.GetProviderModel(context.Background(), message.Model)
	if err != nil {
		return nil, ierror.NewError(err)
	}
	if providerModel == nil {
		return nil, ierror.New(ierror.ErrCodeModelNotFound)
	}

	// 新建供应商
	provider := llm_provider.NewLlmProvider(*providerModel)

	// 如果聊天的uuid为空，则代表新建一个聊天
	isNewChat := false
	if message.ChatUuid == "" {
		isNewChat = true
		message.ChatUuid = uuid.New().String()
		title := message.Content
		// 创建一个聊天
		err = s.storage.CreateChat(context.Background(), message.ChatUuid, title)
		if err != nil {
			return nil, ierror.NewError(err)
		}
	}
	genChatTitle := func() {
		if isNewChat {
			_, err = s.GenChatTitle(message.ChatUuid, message.Model, true)
			if err != nil {
				logger.Error("gen chat title error", err)
			}
		}
	}

	// 查找历史消息
	historyMessageData, _, err := s.storage.GetMessage(context.Background(), message.ChatUuid, 0, 10)
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

	// 转换消息内容
	schemaMessage, err := provider.BuildUserMessage(context.Background(), message)
	if err != nil {
		return nil, ierror.NewError(err)
	}

	// 创建用户消息
	err = s.storage.CreateMessage(context.Background(), message.ChatUuid, data_models.Message{
		Uuid:     uuid.New().String(),
		ChatUuid: message.ChatUuid,
		Message:  schemaMessage,
	})
	if err != nil {
		return nil, ierror.NewError(err)
	}

	// 生成
	stream, err := provider.Completions(context.Background(), append(historyMessages, *schemaMessage))
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

		userStopCh := make(chan struct{})
		s.completionsStopCh[message.ChatUuid] = userStopCh
		defer close(userStopCh)
		defer delete(s.completionsStopCh, message.ChatUuid)

		for {
			select {
			case <-userStopCh:
				stream.Close()
				return
			default:
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
		}
	}()
	messageUuid := uuid.New().String()
	eventsKey := utils.GenEventsKey(messageUuid)
	go func() {
		dataModelMsg := data_models.Message{
			Uuid:     messageUuid,
			ChatUuid: message.ChatUuid,
		}
		for {
			select {
			case <-doneChan:
				dataModelMsg = s.fillCompletionsMsg(dataModelMsg, "done")
				genChatTitle()
				s.app.Event.Emit(eventsKey, dataModelMsg.Message)
				return
			case msg, ok := <-msgChan:
				if !ok || msg == nil {
					dataModelMsg = s.fillCompletionsMsg(dataModelMsg, "done")
					genChatTitle()
					s.app.Event.Emit(eventsKey, dataModelMsg.Message)
					return
				}
				dataModelMsg.Message = msg
				dataModelMsg = s.fillCompletionsMsg(dataModelMsg, "")
				if msg.ResponseMeta != nil && msg.ResponseMeta.FinishReason != "" {
					genChatTitle()
					s.app.Event.Emit(eventsKey, dataModelMsg.Message)
					return
				}
				s.app.Event.Emit(eventsKey, dataModelMsg.Message)
			case err := <-errChan:
				if err == nil {
					continue
				}
				s.fillCompletionsMsg(dataModelMsg, err.Error())
				s.app.Event.Emit(eventsKey, dataModelMsg.Message)
				return
			}
		}
	}()

	return &view_models.Completions{
		ChatUuid:    message.ChatUuid,
		MessageUuid: messageUuid,
	}, nil
}

func (s *Service) StopCompletions(chatUuid string) error {
	stopCh, ok := s.completionsStopCh[chatUuid]
	if !ok {
		return nil
	}
	fmt.Println("stop completions:", chatUuid)
	if stopCh != nil {
		stopCh <- struct{}{}
	}
	return nil
}

// DeleteChat 删除聊天
func (s *Service) DeleteChat(chatUuid string) error {
	err := s.storage.DeleteChat(context.Background(), chatUuid)
	if err != nil {
		return ierror.NewError(err)
	}
	return nil
}

// RenameChat 重命名聊天
func (s *Service) RenameChat(chatUuid, title string) error {
	err := s.storage.RenameChat(context.Background(), chatUuid, title)
	if err != nil {
		return ierror.NewError(err)
	}
	return nil
}

// CollectionChat 收藏/取消收藏对话
func (s *Service) CollectionChat(chatUuid string, isCollection bool) error {
	err := s.storage.CollectionChat(context.Background(), chatUuid, isCollection)
	if err != nil {
		return ierror.NewError(err)
	}

	return nil
}

// GenChatTitle 创建聊天标题
func (s *Service) GenChatTitle(chatUuid, model string, update bool) (string, error) {

	// 获取模型信息
	providerModel, err := s.storage.GetProviderModel(context.Background(), model)
	if err != nil {
		return "", ierror.NewError(err)
	}
	if providerModel == nil {
		return "", ierror.New(ierror.ErrCodeModelNotFound)
	}

	// 新建供应商
	provider := llm_provider.NewLlmProvider(*providerModel)

	historyMessages, _, err := s.storage.GetMessage(context.Background(), chatUuid, 0, 2)
	if err != nil {
		return "", ierror.NewError(err)
	}
	var messages []schema.Message
	for _, item := range historyMessages {
		if item.Message == nil {
			continue
		}
		fmt.Println("msg", item.Message)
		messages = append(messages, *item.Message)
	}

	title, err := llm_provider.GenChatTitle(provider, messages)
	if err != nil {
		return "", ierror.NewError(err)
	}

	if update {
		err := s.storage.RenameChat(context.Background(), chatUuid, title)
		if err != nil {
			return "", ierror.NewError(err)
		}
	}

	return title, nil
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
	err := s.storage.SaveOrUpdateDeltaMessage(context.Background(), dataMsg)
	if err != nil {
		logger.Errorf("save or update delta message failed: %v", err)
	}
	return dataMsg
}
