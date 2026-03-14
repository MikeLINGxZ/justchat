package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/wrapper_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/llm_provider"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/llm_provider/tools"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils/event"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils/ierror"
)

// ChatList 聊天列表
func (s *Service) ChatList(ctx context.Context, offset, limit int, keyword *string, isCollection bool) (*view_models.ChatList, error) {
	chats, total, err := s.storage.GetChats(ctx, offset, limit, keyword, isCollection)
	if err != nil {
		return nil, ierror.NewError(err)
	}

	return &view_models.ChatList{
		Lists: chats,
		Total: total,
	}, nil
}

// ChatInfo 对话信息
func (s *Service) ChatInfo(ctx context.Context, chatUuid string) (*view_models.Chat, error) {
	chat, err := s.storage.GetChat(ctx, chatUuid)
	if err != nil {
		return nil, ierror.NewError(err)
	}
	return chat, nil
}

// ChatMessages 聊天消息
func (s *Service) ChatMessages(ctx context.Context, chatUuid string, offset, limit int) (*view_models.MessageList, error) {
	dataMessages, total, err := s.storage.GetMessage(ctx, chatUuid, offset, limit)
	if err != nil {
		return nil, ierror.NewError(err)
	}

	var messages []view_models.Message
	for _, item := range dataMessages {
		messages = append(messages, item)
	}

	return &view_models.MessageList{
		Messages: messages,
		Total:    total,
	}, nil
}

// Completions 聊天
func (s *Service) Completions(ctx context.Context, inputMessage view_models.Message) (*view_models.Completions, error) {
	if inputMessage.UserMessageExtra == nil {
		return nil, ierror.New(ierror.ErrCodeCompletionsParams)
	}
	selectModelId := inputMessage.UserMessageExtra.ModelId
	selectModelName := inputMessage.UserMessageExtra.ModelName
	userMessageUuid := uuid.New().String()
	assistantMessageUuid := uuid.New().String()
	messageKey := event.GenEventsKey(event.EventTypeMsg, assistantMessageUuid)
	chatUuid := inputMessage.ChatUuid
	isNewChat := inputMessage.ChatUuid == ""
	if isNewChat {
		chatUuid = uuid.New().String()
	}
	inputMessage.MessageUuid = userMessageUuid
	inputMessage.ChatUuid = chatUuid

	// 获取模型信息
	providerModel, err := s.storage.GetProviderModel(context.Background(), selectModelId, selectModelName)
	if err != nil {
		return nil, ierror.NewError(err)
	}
	if providerModel == nil {
		return nil, ierror.New(ierror.ErrCodeModelNotFound)
	}

	// 获取工具集
	agentTools, err := tools.ToolRouter.GetToolsByIds(inputMessage.UserMessageExtra.Tools)
	if err != nil {
		return nil, ierror.NewError(err)
	}

	// todo 获取子agent
	var subAgents []adk.Agent
	for _, agentId := range inputMessage.UserMessageExtra.Agents {
		fmt.Println(agentId)
	}

	// 新建供应商（传入工具以支持 tool calling）
	provider, err := llm_provider.NewLlmProvider(ctx, *providerModel, subAgents, agentTools)
	if err != nil {
		return nil, ierror.NewError(err)
	}

	// 如果聊天的uuid为空，则新建一个聊天
	if isNewChat {
		title := inputMessage.Content
		err = s.storage.CreateChat(context.Background(), chatUuid, title)
		if err != nil {
			return nil, ierror.NewError(err)
		}
	}

	// 查找历史消息
	historyMessageData, _, err := s.storage.GetMessage(ctx, chatUuid, 0, 10)
	if err != nil {
		return nil, ierror.NewError(err)
	}

	// 创建用户消息
	_, err = s.storage.CreateMessage(ctx, chatUuid, inputMessage)
	if err != nil {
		return nil, ierror.NewError(err)
	}

	// 合并用户当前消息和历史消息
	historyMessageData = append(historyMessageData, inputMessage)

	// 转换消息
	var schemaMessages []schema.Message
	for _, item := range historyMessageData {
		schemaMessage, err := item.ToSchemaMessage()
		if err != nil {
			continue
		}
		schemaMessages = append(schemaMessages, *schemaMessage)
	}

	// 创建ai消息
	assistantMessage := data_models.Message{
		OrmModel:    data_models.OrmModel{},
		ChatUuid:    chatUuid,
		MessageUuid: assistantMessageUuid,
		Role:        schema.Assistant,
		AssistantMessageExtra: &data_models.AssistantMessageExtra{
			ToolUses: []data_models.ToolUse{},
		},
	}
	assistantMessageId, err := s.storage.CreateMessage(ctx, chatUuid, assistantMessage)
	if err != nil {
		return nil, ierror.NewError(err)
	}
	assistantMessage.ID = assistantMessageId
	// 生成
	agentIter, err := provider.AgentCompletions(context.Background(), schemaMessages)
	if err != nil {
		return nil, ierror.NewError(err)
	}

	// 获取过程数据
	go func() {
		userStop := make(chan struct{})
		s.completionsStopCh[messageKey] = userStop
		var mo *adk.MessageVariant

		defer func() {
			s.app.Event.Emit(messageKey, assistantMessage)
			err := s.storage.SaveOrUpdateMessage(ctx, assistantMessage)
			if err != nil {
				logger.Error("save or update assistant message error", err)
			}
			if isNewChat {
				go func() {
					_, err = s.genChatTitle(context.Background(), chatUuid, *providerModel, true)
					if err != nil {
						logger.Error("gen chat title error", err)
					}
				}()
			}
			if mo != nil {
				mo.MessageStream.Close()
			}
			delete(s.completionsStopCh, messageKey)
			close(userStop)
		}()

		for {
			select {
			case <-userStop:
				assistantMessage.AssistantMessageExtra.FinishReason = "user stop"
				assistantMessage.AssistantMessageExtra.FinishError = ""
				return
			default:
				event, ok := agentIter.Next()
				if !ok {
					assistantMessage.AssistantMessageExtra.FinishReason = "done"
					assistantMessage.AssistantMessageExtra.FinishError = ""
					return
				}
				if event.Err != nil {
					assistantMessage.AssistantMessageExtra.FinishReason = "error"
					assistantMessage.AssistantMessageExtra.FinishError = event.Err.Error()
					return
				}
				if event.Output == nil || event.Output.MessageOutput == nil {
					continue
				}
				mo = event.Output.MessageOutput
				if mo.Role == schema.Tool && !mo.IsStreaming {
					assistantMessage.AssistantMessageExtra.ToolUses = append(assistantMessage.AssistantMessageExtra.ToolUses, data_models.ToolUse{
						ToolName:   mo.Message.ToolName,
						ToolResult: mo.Message.Content,
					})
					err := s.storage.SaveOrUpdateMessage(ctx, assistantMessage)
					if err != nil {
						assistantMessage.AssistantMessageExtra.FinishReason = "error"
						assistantMessage.AssistantMessageExtra.FinishError = err.Error()
						return
					}
					continue
				}
				if !mo.IsStreaming || mo.MessageStream == nil {
					assistantMessage.AssistantMessageExtra.FinishReason = "error"
					assistantMessage.AssistantMessageExtra.FinishError = "streaming fail"
					return
				}
				for {
					msg, err := mo.MessageStream.Recv()
					if err == io.EOF {
						break
					}
					if err != nil {
						assistantMessage.AssistantMessageExtra.FinishReason = "error"
						assistantMessage.AssistantMessageExtra.FinishError = err.Error()
						s.app.Event.Emit(messageKey, assistantMessage)
						return
					}
					if msg != nil {
						marshal, _ := json.Marshal(msg)
						logger.Info(string(marshal))
					}
					if msg != nil && (msg.Content != "" || msg.ReasoningContent != "") {
						assistantMessage.Content = assistantMessage.Content + msg.Content
						assistantMessage.ReasoningContent = assistantMessage.ReasoningContent + msg.ReasoningContent
						err := s.storage.SaveOrUpdateMessage(ctx, assistantMessage)
						if err != nil {
							assistantMessage.AssistantMessageExtra.FinishReason = "error"
							assistantMessage.AssistantMessageExtra.FinishError = err.Error()
							return
						}
						s.app.Event.Emit(messageKey, assistantMessage)
					}
				}
				assistantMessage.AssistantMessageExtra.FinishReason = "done"
				assistantMessage.AssistantMessageExtra.FinishError = ""
			}
		}
	}()

	return &view_models.Completions{
		ChatUuid:    chatUuid,
		MessageUuid: assistantMessageUuid,
		EventKey:    messageKey,
	}, nil
}

func (s *Service) StopCompletions(messageKey string) error {
	stopCh, ok := s.completionsStopCh[messageKey]
	if !ok {
		return nil
	}
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
func (s *Service) GenChatTitle(ctx context.Context, chatUuid string, modelId uint, modelName string, update bool) (string, error) {

	// 获取模型信息
	providerModel, err := s.storage.GetProviderModel(context.Background(), modelId, modelName)
	if err != nil {
		return "", ierror.NewError(err)
	}
	if providerModel == nil {
		return "", ierror.New(ierror.ErrCodeModelNotFound)
	}

	title, err := s.genChatTitle(ctx, chatUuid, *providerModel, update)
	if err != nil {
		return "", ierror.NewError(err)
	}

	return title, nil
}

func (s *Service) genChatTitle(ctx context.Context, chatUuid string, providerModel wrapper_models.ProviderModel, update bool) (string, error) {
	// 新建供应商
	provider, err := llm_provider.NewLlmProvider(ctx, providerModel, []adk.Agent{}, []tool.BaseTool{})
	if err != nil {
		return "", err
	}

	historyMessages, _, err := s.storage.GetMessage(context.Background(), chatUuid, 0, 2)
	if err != nil {
		return "", err
	}
	var messages []schema.Message
	for _, item := range historyMessages {
		schemaMessage, err := item.ToSchemaMessage()
		if err != nil {
			return "", err
		}
		messages = append(messages, *schemaMessage)
	}

	title, err := provider.GenChatTitle(ctx, messages)
	if err != nil {
		return "", err
	}

	if update {
		err := s.storage.RenameChat(context.Background(), chatUuid, title)
		if err != nil {
			return "", err
		}
	}

	return title, nil
}
