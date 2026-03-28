package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"
	"time"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/wrapper_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/llm_provider"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/llm_provider/tools"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/tasker"
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
	taskUuid := uuid.New().String()
	eventKey := event.GenEventsKey(event.EventTypeTask, taskUuid)
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
	agentTools, toolMetaByID, cleanupTools, err := s.resolveSelectedTools(ctx, inputMessage.UserMessageExtra.Tools)
	if err != nil {
		return nil, ierror.NewError(err)
	}

	// todo 获取子agent
	var subAgents []adk.Agent
	for _, agentId := range inputMessage.UserMessageExtra.Agents {
		fmt.Println(agentId)
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

	task := data_models.Task{
		TaskUuid:             taskUuid,
		ChatUuid:             chatUuid,
		AssistantMessageUuid: assistantMessageUuid,
		Status:               data_models.TaskStatusPending,
		EventKey:             eventKey,
	}
	if err = s.storage.CreateTask(ctx, task); err != nil {
		return nil, ierror.NewError(err)
	}

	var assistantMu sync.Mutex

	cloneAssistantMessage := func(src data_models.Message) data_models.Message {
		clone := src
		if src.UserMessageExtra != nil {
			userExtra := *src.UserMessageExtra
			if len(src.UserMessageExtra.Files) > 0 {
				userExtra.Files = append([]data_models.File(nil), src.UserMessageExtra.Files...)
			}
			if len(src.UserMessageExtra.Tools) > 0 {
				userExtra.Tools = append([]string(nil), src.UserMessageExtra.Tools...)
			}
			if len(src.UserMessageExtra.Agents) > 0 {
				userExtra.Agents = append([]string(nil), src.UserMessageExtra.Agents...)
			}
			clone.UserMessageExtra = &userExtra
		}
		if src.AssistantMessageExtra != nil {
			assistantExtra := *src.AssistantMessageExtra
			if len(src.AssistantMessageExtra.ToolUses) > 0 {
				assistantExtra.ToolUses = append([]data_models.ToolUse(nil), src.AssistantMessageExtra.ToolUses...)
			}
			clone.AssistantMessageExtra = &assistantExtra
		}
		return clone
	}

	persistAssistantSnapshotLocked := func(updateTask bool) error {
		if err := s.storage.SaveOrUpdateMessage(context.Background(), assistantMessage); err != nil {
			return err
		}
		if updateTask {
			if err := s.storage.SaveTask(context.Background(), task); err != nil {
				return err
			}
		}
		s.emitTaskEvent(task, cloneAssistantMessage(assistantMessage))
		return nil
	}

	findToolUseIndexLocked := func(callID string) int {
		if assistantMessage.AssistantMessageExtra == nil {
			return -1
		}
		for idx := range assistantMessage.AssistantMessageExtra.ToolUses {
			if assistantMessage.AssistantMessageExtra.ToolUses[idx].CallID == callID {
				return idx
			}
		}
		return -1
	}

	startToolUseLocked := func(callID, toolName string) error {
		if assistantMessage.AssistantMessageExtra == nil {
			assistantMessage.AssistantMessageExtra = &data_models.AssistantMessageExtra{}
		}
		now := time.Now()
		contentPos := len([]rune(strings.TrimRight(assistantMessage.Content, "\n")))
		toolID := toolName
		displayName := toolName
		description := ""
		if meta, ok := toolMetaByID[toolName]; ok {
			toolID = meta.ID
			displayName = meta.Name
			description = meta.Description
		} else if registeredTool, ok := tools.ToolRouter.GetToolByID(toolName); ok {
			toolID = registeredTool.Id()
			displayName = registeredTool.Name()
			description = registeredTool.Description()
		}
		idx := findToolUseIndexLocked(callID)
		if idx == -1 {
			assistantMessage.AssistantMessageExtra.ToolUses = append(assistantMessage.AssistantMessageExtra.ToolUses, data_models.ToolUse{
				Index:           len(assistantMessage.AssistantMessageExtra.ToolUses) + 1,
				CallID:          callID,
				ContentPos:      contentPos,
				ToolID:          toolID,
				ToolName:        displayName,
				ToolDescription: description,
				Status:          data_models.ToolUseStatusRunning,
				StartedAt:       &now,
			})
		} else {
			toolUse := &assistantMessage.AssistantMessageExtra.ToolUses[idx]
			toolUse.ToolID = toolID
			toolUse.ToolName = displayName
			toolUse.ToolDescription = description
			toolUse.Status = data_models.ToolUseStatusRunning
			if toolUse.Index == 0 {
				toolUse.Index = idx + 1
			}
			if toolUse.StartedAt == nil {
				toolUse.StartedAt = &now
			}
			if toolUse.ContentPos == 0 {
				toolUse.ContentPos = contentPos
			}
			toolUse.FinishedAt = nil
		}
		task.LastOutputAt = &now
		return persistAssistantSnapshotLocked(true)
	}

	finishToolUseLocked := func(callID, toolName, toolResult string, runErr error) error {
		if assistantMessage.AssistantMessageExtra == nil {
			assistantMessage.AssistantMessageExtra = &data_models.AssistantMessageExtra{}
		}
		now := time.Now()
		idx := findToolUseIndexLocked(callID)
		if idx == -1 {
			assistantMessage.AssistantMessageExtra.ToolUses = append(assistantMessage.AssistantMessageExtra.ToolUses, data_models.ToolUse{
				Index:      len(assistantMessage.AssistantMessageExtra.ToolUses) + 1,
				CallID:     callID,
				ContentPos: len([]rune(strings.TrimRight(assistantMessage.Content, "\n"))),
			})
			idx = len(assistantMessage.AssistantMessageExtra.ToolUses) - 1
		}
		toolUse := &assistantMessage.AssistantMessageExtra.ToolUses[idx]
		if toolUse.Index == 0 {
			toolUse.Index = idx + 1
		}
		toolUse.CallID = callID
		if toolName != "" {
			if meta, ok := toolMetaByID[toolName]; ok {
				toolUse.ToolID = meta.ID
				toolUse.ToolName = meta.Name
				toolUse.ToolDescription = meta.Description
			} else if registeredTool, ok := tools.ToolRouter.GetToolByID(toolName); ok {
				toolUse.ToolID = registeredTool.Id()
				toolUse.ToolName = registeredTool.Name()
				toolUse.ToolDescription = registeredTool.Description()
			} else {
				if toolUse.ToolID == "" {
					toolUse.ToolID = toolName
				}
				if toolUse.ToolName == "" {
					toolUse.ToolName = toolName
				}
			}
		}
		if toolUse.StartedAt == nil {
			toolUse.StartedAt = &now
		}
		toolUse.FinishedAt = &now
		toolUse.ToolResult = toolResult
		toolUse.ElapsedMs = now.Sub(*toolUse.StartedAt).Milliseconds()
		if runErr != nil {
			toolUse.Status = data_models.ToolUseStatusError
			if toolUse.ToolResult == "" {
				toolUse.ToolResult = runErr.Error()
			}
		} else {
			toolUse.Status = data_models.ToolUseStatusDone
		}
		task.LastOutputAt = &now
		return persistAssistantSnapshotLocked(true)
	}

	finalizeRunningToolUsesLocked := func() {
		if assistantMessage.AssistantMessageExtra == nil {
			return
		}
		now := time.Now()
		for idx := range assistantMessage.AssistantMessageExtra.ToolUses {
			toolUse := &assistantMessage.AssistantMessageExtra.ToolUses[idx]
			if toolUse.Status != data_models.ToolUseStatusRunning && toolUse.Status != data_models.ToolUseStatusPending {
				continue
			}
			if toolUse.StartedAt == nil {
				toolUse.StartedAt = &now
			}
			toolUse.FinishedAt = &now
			toolUse.ElapsedMs = now.Sub(*toolUse.StartedAt).Milliseconds()
			toolUse.Status = data_models.ToolUseStatusError
		}
	}

	toolMiddleware := compose.ToolMiddleware{
		Invokable: func(next compose.InvokableToolEndpoint) compose.InvokableToolEndpoint {
			return func(toolCtx context.Context, input *compose.ToolInput) (*compose.ToolOutput, error) {
				assistantMu.Lock()
				err := startToolUseLocked(input.CallID, input.Name)
				assistantMu.Unlock()
				if err != nil {
					return nil, err
				}

				output, runErr := next(toolCtx, input)
				result := ""
				if output != nil {
					result = output.Result
				}

				assistantMu.Lock()
				err = finishToolUseLocked(input.CallID, input.Name, result, runErr)
				assistantMu.Unlock()
				if err != nil {
					return nil, err
				}
				if runErr != nil {
					return nil, runErr
				}
				return output, nil
			}
		},
	}

	// 新建供应商（传入工具以支持 tool calling）
	provider, err := llm_provider.NewLlmProvider(ctx, *providerModel, subAgents, agentTools, toolMiddleware)
	if err != nil {
		return nil, ierror.NewError(err)
	}
	// 生成
	agentIter, err := provider.AgentCompletions(context.Background(), schemaMessages)
	if err != nil {
		return nil, ierror.NewError(err)
	}

	// 获取过程数据
	tasker.Manager.StartTask(tasker.Runtime{
		TaskUUID:             taskUuid,
		ChatUUID:             chatUuid,
		AssistantMessageUUID: assistantMessageUuid,
		EventKey:             eventKey,
	}, func(userStop <-chan struct{}) {
		var mo *adk.MessageVariant
		now := time.Now()
		assistantMu.Lock()
		task.Status = data_models.TaskStatusRunning
		task.StartedAt = &now
		saveErr := s.storage.SaveTask(context.Background(), task)
		assistantMu.Unlock()
		if saveErr != nil {
			logger.Error("save task running status error", saveErr)
		}

		defer func() {
			if cleanupTools != nil {
				cleanupTools()
			}
			assistantMu.Lock()
			finalizeRunningToolUsesLocked()
			finishedAt := time.Now()
			task.FinishReason = assistantMessage.AssistantMessageExtra.FinishReason
			task.FinishError = assistantMessage.AssistantMessageExtra.FinishError
			task.FinishedAt = &finishedAt
			switch assistantMessage.AssistantMessageExtra.FinishReason {
			case "done":
				task.Status = data_models.TaskStatusCompleted
			case "user stop":
				task.Status = data_models.TaskStatusStopped
			default:
				task.Status = data_models.TaskStatusFailed
			}
			saveErr := persistAssistantSnapshotLocked(true)
			assistantMu.Unlock()
			if saveErr != nil {
				logger.Error("save task finished status error", saveErr)
			}
			if isNewChat {
				go func() {
					_, titleErr := s.genChatTitle(context.Background(), chatUuid, *providerModel, true)
					if titleErr != nil {
						logger.Error("gen chat title error", titleErr)
					}
				}()
			}
			if mo != nil && mo.MessageStream != nil {
				mo.MessageStream.Close()
			}
		}()

		for {
			nextCh := make(chan struct {
				ev *adk.AgentEvent
				ok bool
			}, 1)
			go func() {
				ev, ok := agentIter.Next()
				nextCh <- struct {
					ev *adk.AgentEvent
					ok bool
				}{ev, ok}
			}()
			var event *adk.AgentEvent
			var ok bool
			select {
			case <-userStop:
				assistantMu.Lock()
				assistantMessage.AssistantMessageExtra.FinishReason = "user stop"
				assistantMessage.AssistantMessageExtra.FinishError = ""
				assistantMu.Unlock()
				return
			case nr := <-nextCh:
				event, ok = nr.ev, nr.ok
			}
			if !ok {
				assistantMu.Lock()
				assistantMessage.AssistantMessageExtra.FinishReason = "done"
				assistantMessage.AssistantMessageExtra.FinishError = ""
				assistantMu.Unlock()
				return
			}
			if event.Err != nil {
				assistantMu.Lock()
				assistantMessage.AssistantMessageExtra.FinishReason = "error"
				assistantMessage.AssistantMessageExtra.FinishError = event.Err.Error()
				assistantMu.Unlock()
				return
			}
			if event.Output == nil || event.Output.MessageOutput == nil {
				continue
			}
			mo = event.Output.MessageOutput
			if mo.Role == schema.Tool && !mo.IsStreaming {
				assistantMu.Lock()
				saveErr := finishToolUseLocked(mo.Message.ToolCallID, mo.Message.ToolName, mo.Message.Content, nil)
				if saveErr != nil {
					assistantMessage.AssistantMessageExtra.FinishReason = "error"
					assistantMessage.AssistantMessageExtra.FinishError = saveErr.Error()
					assistantMu.Unlock()
					return
				}
				assistantMu.Unlock()
				continue
			}
			if !mo.IsStreaming || mo.MessageStream == nil {
				assistantMu.Lock()
				assistantMessage.AssistantMessageExtra.FinishReason = "error"
				assistantMessage.AssistantMessageExtra.FinishError = "streaming fail"
				assistantMu.Unlock()
				return
			}
		streamLoop:
			for {
				recvCh := make(chan struct {
					msg adk.Message
					err error
				}, 1)
				go func() {
					m, e := mo.MessageStream.Recv()
					recvCh <- struct {
						msg adk.Message
						err error
					}{m, e}
				}()
				select {
				case <-userStop:
					assistantMu.Lock()
					assistantMessage.AssistantMessageExtra.FinishReason = "user stop"
					assistantMessage.AssistantMessageExtra.FinishError = ""
					assistantMu.Unlock()
					if mo.MessageStream != nil {
						mo.MessageStream.Close()
						mo.MessageStream = nil
					}
					return
				case rr := <-recvCh:
					msg, err := rr.msg, rr.err
					if err == io.EOF {
						break streamLoop
					}
					if err != nil {
						assistantMu.Lock()
						assistantMessage.AssistantMessageExtra.FinishReason = "error"
						assistantMessage.AssistantMessageExtra.FinishError = err.Error()
						assistantMu.Unlock()
						return
					}
					if msg != nil {
						marshal, _ := json.Marshal(msg)
						logger.Info(string(marshal))
					}
					if msg != nil && (msg.Content != "" || msg.ReasoningContent != "") {
						assistantMu.Lock()
						assistantMessage.Content = assistantMessage.Content + msg.Content
						assistantMessage.ReasoningContent = assistantMessage.ReasoningContent + msg.ReasoningContent
						lastOutputAt := time.Now()
						task.LastOutputAt = &lastOutputAt
						saveErr := persistAssistantSnapshotLocked(true)
						if saveErr != nil {
							assistantMessage.AssistantMessageExtra.FinishReason = "error"
							assistantMessage.AssistantMessageExtra.FinishError = saveErr.Error()
							assistantMu.Unlock()
							return
						}
						assistantMu.Unlock()
					}
				}
			}
			assistantMu.Lock()
			assistantMessage.AssistantMessageExtra.FinishReason = "done"
			assistantMessage.AssistantMessageExtra.FinishError = ""
			assistantMu.Unlock()
		}
	})

	return &view_models.Completions{
		ChatUuid:    chatUuid,
		TaskUuid:    taskUuid,
		MessageUuid: assistantMessageUuid,
		EventKey:    eventKey,
	}, nil
}

func (s *Service) StopCompletions(messageKey string) error {
	tasker.Manager.StopByEventKey(messageKey)
	return nil
}

func (s *Service) StopTask(taskUuid string) error {
	tasker.Manager.StopTask(taskUuid)
	return nil
}

func (s *Service) GetTask(ctx context.Context, taskUuid string) (*view_models.Task, error) {
	task, err := s.storage.GetTask(ctx, taskUuid)
	if err != nil {
		return nil, ierror.NewError(err)
	}
	if task == nil {
		return nil, nil
	}
	viewTask := view_models.Task(*task)
	return &viewTask, nil
}

func (s *Service) GetChatActiveTask(ctx context.Context, chatUuid string) (*view_models.Task, error) {
	task, err := s.storage.GetChatActiveTask(ctx, chatUuid)
	if err != nil {
		return nil, ierror.NewError(err)
	}
	if task == nil {
		return nil, nil
	}
	viewTask := view_models.Task(*task)
	return &viewTask, nil
}

func (s *Service) GetRunningTasks(ctx context.Context) (*view_models.TaskList, error) {
	tasks, err := s.storage.GetRunningTasks(ctx)
	if err != nil {
		return nil, ierror.NewError(err)
	}
	viewTasks := make([]view_models.Task, 0, len(tasks))
	for _, task := range tasks {
		viewTasks = append(viewTasks, view_models.Task(task))
	}
	return &view_models.TaskList{Tasks: viewTasks}, nil
}

func (s *Service) emitTaskEvent(task data_models.Task, assistantMessage data_models.Message) {
	s.app.Event.Emit(task.EventKey, view_models.TaskStreamEvent{
		TaskUuid:         task.TaskUuid,
		ChatUuid:         task.ChatUuid,
		EventKey:         task.EventKey,
		Status:           task.Status,
		FinishReason:     task.FinishReason,
		FinishError:      task.FinishError,
		AssistantMessage: assistantMessage,
	})
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
	provider, err := llm_provider.NewLlmProvider(ctx, providerModel, []adk.Agent{}, []tool.BaseTool{}, compose.ToolMiddleware{})
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
