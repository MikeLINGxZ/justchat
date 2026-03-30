package service

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"
	"sync/atomic"
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
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/tool_approval"
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

func preserveWorkflowPreface(message *data_models.Message) {
	if message == nil {
		return
	}
	if message.AssistantMessageExtra == nil {
		message.AssistantMessageExtra = &data_models.AssistantMessageExtra{}
	}
	if message.AssistantMessageExtra.PrefaceContent == "" && strings.TrimSpace(message.Content) != "" {
		message.AssistantMessageExtra.PrefaceContent = message.Content
	}
	if message.AssistantMessageExtra.PrefaceReasoningContent == "" && strings.TrimSpace(message.ReasoningContent) != "" {
		message.AssistantMessageExtra.PrefaceReasoningContent = message.ReasoningContent
	}
}

func resetDirectAssistantState(message *data_models.Message) {
	if message == nil {
		return
	}
	if message.AssistantMessageExtra == nil {
		message.AssistantMessageExtra = &data_models.AssistantMessageExtra{}
	}
	message.Content = ""
	message.ReasoningContent = ""
	message.AssistantMessageExtra.RouteType = ""
	message.AssistantMessageExtra.CurrentStage = ""
	message.AssistantMessageExtra.CurrentAgent = ""
	message.AssistantMessageExtra.PendingApprovals = nil
	message.AssistantMessageExtra.ExecutionTrace = data_models.ExecutionTrace{Steps: []data_models.TraceStep{}}
	message.AssistantMessageExtra.FinishError = ""
}

func buildApprovalDecisionToolResult(toolName string, decision data_models.ToolApprovalDecision, comment string) string {
	switch decision {
	case data_models.ToolApprovalDecisionAllow:
		return fmt.Sprintf("用户已允许执行工具：%s。", toolName)
	case data_models.ToolApprovalDecisionCustom:
		if strings.TrimSpace(comment) == "" {
			return fmt.Sprintf("用户没有直接批准执行工具：%s，并要求先补充更多说明。请根据用户反馈调整方案。", toolName)
		}
		return fmt.Sprintf("用户没有直接批准执行工具：%s，并提供了意见：%s。请根据该意见调整方案，必要时重新发起更合适的工具请求。", toolName, comment)
	default:
		return fmt.Sprintf("用户拒绝了工具：%s 的本次调用。请不要执行该操作，并改用无需该操作的方案。", toolName)
	}
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

	// 获取工作流工具集
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
			ToolUses:       []data_models.ToolUse{},
			ExecutionTrace: data_models.ExecutionTrace{Steps: []data_models.TraceStep{}},
			RouteType:      "",
			CurrentStage:   "等待执行",
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
			if len(src.AssistantMessageExtra.PendingApprovals) > 0 {
				assistantExtra.PendingApprovals = append([]data_models.ToolApprovalSummary(nil), src.AssistantMessageExtra.PendingApprovals...)
			}
			if len(src.AssistantMessageExtra.ExecutionTrace.Steps) > 0 {
				assistantExtra.ExecutionTrace.Steps = make([]data_models.TraceStep, 0, len(src.AssistantMessageExtra.ExecutionTrace.Steps))
				for _, step := range src.AssistantMessageExtra.ExecutionTrace.Steps {
					clonedStep := step
					if len(step.Metadata) > 0 {
						clonedStep.Metadata = make(map[string]interface{}, len(step.Metadata))
						for key, value := range step.Metadata {
							clonedStep.Metadata[key] = value
						}
					}
					if len(step.DetailBlocks) > 0 {
						clonedStep.DetailBlocks = append([]data_models.TraceDetailBlock(nil), step.DetailBlocks...)
					}
					assistantExtra.ExecutionTrace.Steps = append(assistantExtra.ExecutionTrace.Steps, clonedStep)
				}
			}
			clone.AssistantMessageExtra = &assistantExtra
		}
		return clone
	}

	var pendingTraceDelta []data_models.TraceStep
	var lastSnapshotPersistAt time.Time
	var lastSnapshotPersistContentLen int

	emitAssistantSnapshotLocked := func() {
		traceDelta := append([]data_models.TraceStep(nil), pendingTraceDelta...)
		pendingTraceDelta = nil
		s.emitTaskEvent(task, cloneAssistantMessage(assistantMessage), traceDelta)
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
		lastSnapshotPersistAt = time.Now()
		lastSnapshotPersistContentLen = len([]rune(assistantMessage.Content))
		emitAssistantSnapshotLocked()
		return nil
	}

	persistAssistantSnapshotThrottledLocked := func(updateTask bool) error {
		const minPersistInterval = 350 * time.Millisecond
		const minPersistContentDelta = 48

		currentContentLen := len([]rune(assistantMessage.Content))
		shouldPersist := lastSnapshotPersistAt.IsZero() ||
			time.Since(lastSnapshotPersistAt) >= minPersistInterval ||
			currentContentLen-lastSnapshotPersistContentLen >= minPersistContentDelta
		if shouldPersist {
			return persistAssistantSnapshotLocked(updateTask)
		}
		emitAssistantSnapshotLocked()
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

	findTraceStepIndexLocked := func(stepID string) int {
		if assistantMessage.AssistantMessageExtra == nil {
			return -1
		}
		for idx := range assistantMessage.AssistantMessageExtra.ExecutionTrace.Steps {
			if assistantMessage.AssistantMessageExtra.ExecutionTrace.Steps[idx].StepID == stepID {
				return idx
			}
		}
		return -1
	}

	findPendingApprovalIndexLocked := func(approvalID string) int {
		if assistantMessage.AssistantMessageExtra == nil {
			return -1
		}
		for idx := range assistantMessage.AssistantMessageExtra.PendingApprovals {
			if assistantMessage.AssistantMessageExtra.PendingApprovals[idx].ApprovalID == approvalID {
				return idx
			}
		}
		return -1
	}

	upsertPendingApprovalLocked := func(summary data_models.ToolApprovalSummary) {
		if assistantMessage.AssistantMessageExtra == nil {
			assistantMessage.AssistantMessageExtra = &data_models.AssistantMessageExtra{}
		}
		idx := findPendingApprovalIndexLocked(summary.ApprovalID)
		if idx == -1 {
			assistantMessage.AssistantMessageExtra.PendingApprovals = append(assistantMessage.AssistantMessageExtra.PendingApprovals, summary)
			return
		}
		assistantMessage.AssistantMessageExtra.PendingApprovals[idx] = summary
	}

	removePendingApprovalLocked := func(approvalID string) {
		if assistantMessage.AssistantMessageExtra == nil {
			return
		}
		idx := findPendingApprovalIndexLocked(approvalID)
		if idx == -1 {
			return
		}
		assistantMessage.AssistantMessageExtra.PendingApprovals = append(
			assistantMessage.AssistantMessageExtra.PendingApprovals[:idx],
			assistantMessage.AssistantMessageExtra.PendingApprovals[idx+1:]...,
		)
	}

	appendTraceStepLocked := func(step data_models.TraceStep) error {
		if assistantMessage.AssistantMessageExtra == nil {
			assistantMessage.AssistantMessageExtra = &data_models.AssistantMessageExtra{}
		}
		if step.StartedAt == nil {
			now := time.Now()
			step.StartedAt = &now
		}
		assistantMessage.AssistantMessageExtra.ExecutionTrace.Steps = append(assistantMessage.AssistantMessageExtra.ExecutionTrace.Steps, step)
		pendingTraceDelta = append(pendingTraceDelta, step)
		task.LastOutputAt = step.StartedAt
		return persistAssistantSnapshotLocked(true)
	}

	updateTraceStepLocked := func(step data_models.TraceStep) error {
		if assistantMessage.AssistantMessageExtra == nil {
			assistantMessage.AssistantMessageExtra = &data_models.AssistantMessageExtra{}
		}
		idx := findTraceStepIndexLocked(step.StepID)
		if idx == -1 {
			return appendTraceStepLocked(step)
		}
		assistantMessage.AssistantMessageExtra.ExecutionTrace.Steps[idx] = step
		pendingTraceDelta = append(pendingTraceDelta, step)
		now := time.Now()
		task.LastOutputAt = &now
		return persistAssistantSnapshotLocked(true)
	}

	startTraceStepLocked := func(stepID, parentStepID string, stepType data_models.TraceStepType, title, summary, inputPreview, stage, agentName string, detailBlocks []data_models.TraceDetailBlock, metadata map[string]interface{}) error {
		now := time.Now()
		if assistantMessage.AssistantMessageExtra == nil {
			assistantMessage.AssistantMessageExtra = &data_models.AssistantMessageExtra{}
		}
		assistantMessage.AssistantMessageExtra.CurrentStage = stage
		assistantMessage.AssistantMessageExtra.CurrentAgent = agentName
		return appendTraceStepLocked(data_models.TraceStep{
			StepID:       stepID,
			ParentStepID: parentStepID,
			Type:         stepType,
			Title:        title,
			Summary:      summary,
			InputPreview: inputPreview,
			Status:       data_models.TraceStepStatusRunning,
			AgentName:    agentName,
			StartedAt:    &now,
			DetailBlocks: detailBlocks,
			Metadata:     metadata,
		})
	}

	finishTraceStepLocked := func(stepID, summary, outputPreview, stage, agentName string, status data_models.TraceStepStatus, detailBlocks []data_models.TraceDetailBlock, metadata map[string]interface{}) error {
		if assistantMessage.AssistantMessageExtra == nil {
			assistantMessage.AssistantMessageExtra = &data_models.AssistantMessageExtra{}
		}
		idx := findTraceStepIndexLocked(stepID)
		if idx == -1 {
			return nil
		}
		step := assistantMessage.AssistantMessageExtra.ExecutionTrace.Steps[idx]
		now := time.Now()
		if step.StartedAt == nil {
			step.StartedAt = &now
		}
		step.Status = status
		step.Summary = summary
		step.OutputPreview = outputPreview
		step.FinishedAt = &now
		step.ElapsedMs = now.Sub(*step.StartedAt).Milliseconds()
		if detailBlocks != nil {
			step.DetailBlocks = detailBlocks
		}
		if metadata != nil {
			step.Metadata = metadata
		}
		assistantMessage.AssistantMessageExtra.CurrentStage = stage
		assistantMessage.AssistantMessageExtra.CurrentAgent = agentName
		return updateTraceStepLocked(step)
	}

	updateCurrentStageLocked := func(stage, agentName string) {
		if assistantMessage.AssistantMessageExtra == nil {
			assistantMessage.AssistantMessageExtra = &data_models.AssistantMessageExtra{}
		}
		assistantMessage.AssistantMessageExtra.CurrentStage = stage
		assistantMessage.AssistantMessageExtra.CurrentAgent = agentName
	}

	startToolUseLocked := func(toolCtx context.Context, callID, toolName, toolArgs string) error {
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
		parentStepID, _ := toolCtx.Value(traceParentStepIDContextKey).(string)
		agentName, _ := toolCtx.Value(traceAgentNameContextKey).(string)
		assistantMessage.AssistantMessageExtra.CurrentStage = "子任务执行"
		assistantMessage.AssistantMessageExtra.CurrentAgent = agentName
		return appendTraceStepLocked(data_models.TraceStep{
			StepID:       callID,
			ParentStepID: parentStepID,
			Type:         data_models.TraceStepTypeToolCall,
			Title:        fmt.Sprintf("调用工具：%s", displayName),
			Summary:      description,
			Status:       data_models.TraceStepStatusRunning,
			AgentName:    agentName,
			ToolName:     displayName,
			StartedAt:    &now,
			InputPreview: compactText(toolArgs, 180),
			DetailBlocks: []data_models.TraceDetailBlock{
				{
					Kind:    "tool_args",
					Title:   "工具参数",
					Content: toolArgs,
					Format:  data_models.TraceDetailFormatJSON,
				},
			},
			Metadata: map[string]interface{}{
				"tool_id": toolID,
			},
		})
	}

	finishToolUseWithStatusLocked := func(toolCtx context.Context, callID, toolName, toolResult string, toolStatus data_models.ToolUseStatus, traceStatus data_models.TraceStepStatus, runErr error) error {
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
		toolUse.Status = toolStatus
		if runErr != nil && toolUse.ToolResult == "" {
			toolUse.ToolResult = runErr.Error()
		}
		task.LastOutputAt = &now
		traceIdx := findTraceStepIndexLocked(callID)
		if traceIdx != -1 {
			traceStep := assistantMessage.AssistantMessageExtra.ExecutionTrace.Steps[traceIdx]
			if traceStep.StartedAt == nil {
				traceStep.StartedAt = &now
			}
			traceStep.Status = traceStatus
			traceStep.OutputPreview = compactText(toolUse.ToolResult, 240)
			traceStep.Summary = toolUse.ToolDescription
			traceStep.ToolName = toolUse.ToolName
			traceStep.AgentName, _ = toolCtx.Value(traceAgentNameContextKey).(string)
			traceStep.FinishedAt = &now
			traceStep.ElapsedMs = now.Sub(*traceStep.StartedAt).Milliseconds()
			traceStep.DetailBlocks = append(traceStep.DetailBlocks[:0:0], traceStep.DetailBlocks...)
			traceStep.DetailBlocks = append(traceStep.DetailBlocks, data_models.TraceDetailBlock{
				Kind:    "tool_result",
				Title:   "工具结果",
				Content: toolUse.ToolResult,
				Format:  data_models.TraceDetailFormatText,
			})
			if runErr != nil {
				traceStep.DetailBlocks = append(traceStep.DetailBlocks, data_models.TraceDetailBlock{
					Kind:    "tool_result",
					Title:   "错误信息",
					Content: runErr.Error(),
					Format:  data_models.TraceDetailFormatText,
				})
			}
			return updateTraceStepLocked(traceStep)
		}
		return persistAssistantSnapshotLocked(true)
	}

	finishToolUseLocked := func(toolCtx context.Context, callID, toolName, toolResult string, runErr error) error {
		toolStatus := data_models.ToolUseStatusDone
		traceStatus := data_models.TraceStepStatusDone
		if runErr != nil {
			toolStatus = data_models.ToolUseStatusError
			traceStatus = data_models.TraceStepStatusError
		}
		return finishToolUseWithStatusLocked(toolCtx, callID, toolName, toolResult, toolStatus, traceStatus, runErr)
	}

	setToolApprovalPendingLocked := func(toolCtx context.Context, callID string, approval data_models.ToolApproval) error {
		if assistantMessage.AssistantMessageExtra == nil {
			assistantMessage.AssistantMessageExtra = &data_models.AssistantMessageExtra{}
		}

		summary := approval.Summary()
		upsertPendingApprovalLocked(summary)
		task.Status = data_models.TaskStatusWaitingApproval
		assistantMessage.AssistantMessageExtra.CurrentStage = "等待用户确认"
		assistantMessage.AssistantMessageExtra.CurrentAgent, _ = toolCtx.Value(traceAgentNameContextKey).(string)

		if idx := findToolUseIndexLocked(callID); idx != -1 {
			assistantMessage.AssistantMessageExtra.ToolUses[idx].Status = data_models.ToolUseStatusAwaitingApproval
			assistantMessage.AssistantMessageExtra.ToolUses[idx].ToolResult = approval.Message
		}

		if traceIdx := findTraceStepIndexLocked(callID); traceIdx != -1 {
			traceStep := assistantMessage.AssistantMessageExtra.ExecutionTrace.Steps[traceIdx]
			traceStep.Status = data_models.TraceStepStatusAwaitingApproval
			traceStep.Summary = approval.Message
			traceStep.DetailBlocks = append(traceStep.DetailBlocks[:0:0], traceStep.DetailBlocks...)
			traceStep.DetailBlocks = append(traceStep.DetailBlocks, data_models.TraceDetailBlock{
				Kind:    "approval_request",
				Title:   "确认请求",
				Content: approval.Message,
				Format:  data_models.TraceDetailFormatMarkdown,
			})
			if traceStep.Metadata == nil {
				traceStep.Metadata = map[string]interface{}{}
			}
			traceStep.Metadata["approval_id"] = approval.ApprovalID
			traceStep.Metadata["approval_status"] = approval.Status
			traceStep.Metadata["approval_decision_required"] = true
			traceStep.Metadata["approval_title"] = approval.Title
			traceStep.Metadata["approval_message"] = approval.Message
			traceStep.Metadata["approval_scope"] = approval.Scope
			assistantMessage.AssistantMessageExtra.ExecutionTrace.Steps[traceIdx] = traceStep
			pendingTraceDelta = append(pendingTraceDelta, traceStep)
		}
		return persistAssistantSnapshotLocked(true)
	}

	resumeApprovedToolLocked := func(toolCtx context.Context, callID string, approval data_models.ToolApproval) error {
		if assistantMessage.AssistantMessageExtra == nil {
			assistantMessage.AssistantMessageExtra = &data_models.AssistantMessageExtra{}
		}

		removePendingApprovalLocked(approval.ApprovalID)
		task.Status = data_models.TaskStatusRunning
		assistantMessage.AssistantMessageExtra.CurrentStage = "子任务执行"
		assistantMessage.AssistantMessageExtra.CurrentAgent, _ = toolCtx.Value(traceAgentNameContextKey).(string)

		if idx := findToolUseIndexLocked(callID); idx != -1 {
			toolUse := &assistantMessage.AssistantMessageExtra.ToolUses[idx]
			toolUse.Status = data_models.ToolUseStatusRunning
			toolUse.ToolResult = ""
			toolUse.FinishedAt = nil
		}
		if traceIdx := findTraceStepIndexLocked(callID); traceIdx != -1 {
			traceStep := assistantMessage.AssistantMessageExtra.ExecutionTrace.Steps[traceIdx]
			traceStep.Status = data_models.TraceStepStatusRunning
			traceStep.FinishedAt = nil
			traceStep.OutputPreview = ""
			if traceStep.Metadata == nil {
				traceStep.Metadata = map[string]interface{}{}
			}
			traceStep.Metadata["approval_id"] = approval.ApprovalID
			traceStep.Metadata["approval_status"] = approval.Status
			traceStep.Metadata["approval_decision_required"] = false
			traceStep.Metadata["approval_decision"] = approval.Decision
			traceStep.DetailBlocks = append(traceStep.DetailBlocks[:0:0], traceStep.DetailBlocks...)
			traceStep.DetailBlocks = append(traceStep.DetailBlocks, data_models.TraceDetailBlock{
				Kind:    "approval_response",
				Title:   "确认结果",
				Content: "用户已允许继续执行",
				Format:  data_models.TraceDetailFormatText,
			})
			assistantMessage.AssistantMessageExtra.ExecutionTrace.Steps[traceIdx] = traceStep
			pendingTraceDelta = append(pendingTraceDelta, traceStep)
		}
		return persistAssistantSnapshotLocked(true)
	}

	finalizeRunningToolUsesLocked := func() {
		if assistantMessage.AssistantMessageExtra == nil {
			return
		}
		now := time.Now()
		for idx := range assistantMessage.AssistantMessageExtra.ToolUses {
			toolUse := &assistantMessage.AssistantMessageExtra.ToolUses[idx]
			if toolUse.Status != data_models.ToolUseStatusRunning &&
				toolUse.Status != data_models.ToolUseStatusPending &&
				toolUse.Status != data_models.ToolUseStatusAwaitingApproval {
				continue
			}
			if toolUse.StartedAt == nil {
				toolUse.StartedAt = &now
			}
			toolUse.FinishedAt = &now
			toolUse.ElapsedMs = now.Sub(*toolUse.StartedAt).Milliseconds()
			toolUse.Status = data_models.ToolUseStatusError
		}
		for idx := range assistantMessage.AssistantMessageExtra.ExecutionTrace.Steps {
			step := &assistantMessage.AssistantMessageExtra.ExecutionTrace.Steps[idx]
			if step.Status != data_models.TraceStepStatusRunning &&
				step.Status != data_models.TraceStepStatusPending &&
				step.Status != data_models.TraceStepStatusAwaitingApproval {
				continue
			}
			if step.StartedAt == nil {
				step.StartedAt = &now
			}
			step.FinishedAt = &now
			step.ElapsedMs = now.Sub(*step.StartedAt).Milliseconds()
			step.Status = data_models.TraceStepStatusError
		}
		assistantMessage.AssistantMessageExtra.PendingApprovals = nil
	}

	var handoffMu sync.Mutex
	var workflowHandoffDecision *workflowHandoff
	setWorkflowHandoff := func(handoff workflowHandoff) {
		handoffMu.Lock()
		defer handoffMu.Unlock()
		cloned := handoff
		workflowHandoffDecision = &cloned
	}
	getWorkflowHandoff := func() *workflowHandoff {
		handoffMu.Lock()
		defer handoffMu.Unlock()
		if workflowHandoffDecision == nil {
			return nil
		}
		cloned := *workflowHandoffDecision
		return &cloned
	}

	toolMiddleware := compose.ToolMiddleware{
		Invokable: func(next compose.InvokableToolEndpoint) compose.InvokableToolEndpoint {
			return func(toolCtx context.Context, input *compose.ToolInput) (*compose.ToolOutput, error) {
				finishWithMiddlewareError := func(err error) (*compose.ToolOutput, error) {
					if input.Name == workflowHandoffToolName {
						return nil, err
					}
					assistantMu.Lock()
					finishErr := finishToolUseLocked(toolCtx, input.CallID, input.Name, "", err)
					assistantMu.Unlock()
					if finishErr != nil {
						return nil, finishErr
					}
					return nil, err
				}

				if input.Name != workflowHandoffToolName {
					assistantMu.Lock()
					err := startToolUseLocked(toolCtx, input.CallID, input.Name, input.Arguments)
					assistantMu.Unlock()
					if err != nil {
						return nil, err
					}

					if registeredTool, ok := tools.ToolRouter.GetToolByID(input.Name); ok {
						if approvalTool, ok := registeredTool.(tool_approval.ApprovalAwareTool); ok {
							prompt, err := approvalTool.BuildApprovalPrompt(toolCtx, input.Arguments)
							if err != nil {
								return finishWithMiddlewareError(err)
							}

							requestedAt := time.Now()
							approval := data_models.ToolApproval{
								ApprovalID:           uuid.NewString(),
								TaskUuid:             task.TaskUuid,
								ChatUuid:             task.ChatUuid,
								AssistantMessageUuid: task.AssistantMessageUuid,
								ToolCallID:           input.CallID,
								ToolID:               registeredTool.Id(),
								ToolName:             registeredTool.Name(),
								Status:               data_models.ToolApprovalStatusPending,
								Title:                prompt.Title,
								Message:              prompt.Message,
								Scope:                prompt.Scope,
								ArgumentsJSON:        input.Arguments,
								RequestedAt:          &requestedAt,
							}

							if err := tool_approval.Manager.Register(approval.ApprovalID); err != nil {
								return finishWithMiddlewareError(err)
							}
							if err := s.storage.CreateToolApproval(context.Background(), approval); err != nil {
								tool_approval.Manager.Cancel(approval.ApprovalID)
								return finishWithMiddlewareError(err)
							}

							assistantMu.Lock()
							err = setToolApprovalPendingLocked(toolCtx, input.CallID, approval)
							assistantMu.Unlock()
							if err != nil {
								tool_approval.Manager.Cancel(approval.ApprovalID)
								return finishWithMiddlewareError(err)
							}

							waitResult, err := tool_approval.Manager.Wait(toolCtx, approval.ApprovalID)
							if err != nil {
								expiredAt := time.Now()
								approval.Status = data_models.ToolApprovalStatusExpired
								approval.Decision = data_models.ToolApprovalDecisionReject
								approval.ResponseComment = err.Error()
								approval.RespondedAt = &expiredAt
								_ = s.storage.SaveToolApproval(context.Background(), approval)
								return finishWithMiddlewareError(err)
							}

							approval.Status = data_models.ToolApprovalStatusResolved
							approval.Decision = waitResult.Decision
							approval.ResponseComment = waitResult.Comment
							approval.RespondedAt = &waitResult.RespondedAt

							switch waitResult.Decision {
							case data_models.ToolApprovalDecisionAllow:
								assistantMu.Lock()
								err = resumeApprovedToolLocked(toolCtx, input.CallID, approval)
								assistantMu.Unlock()
								if err != nil {
									return finishWithMiddlewareError(err)
								}
							case data_models.ToolApprovalDecisionReject, data_models.ToolApprovalDecisionCustom:
								result := buildApprovalDecisionToolResult(approval.ToolName, waitResult.Decision, waitResult.Comment)
								assistantMu.Lock()
								removePendingApprovalLocked(approval.ApprovalID)
								task.Status = data_models.TaskStatusRunning
								agentName, _ := toolCtx.Value(traceAgentNameContextKey).(string)
								updateCurrentStageLocked("子任务执行", agentName)
								err = finishToolUseWithStatusLocked(
									toolCtx,
									input.CallID,
									input.Name,
									result,
									data_models.ToolUseStatusRejected,
									data_models.TraceStepStatusRejected,
									nil,
								)
								assistantMu.Unlock()
								if err != nil {
									return nil, err
								}
								return &compose.ToolOutput{Result: result}, nil
							default:
								return finishWithMiddlewareError(fmt.Errorf("unknown approval decision: %s", waitResult.Decision))
							}
						}
					}
				}

				output, runErr := next(toolCtx, input)
				result := ""
				if output != nil {
					result = output.Result
				}

				if input.Name != workflowHandoffToolName {
					assistantMu.Lock()
					err = finishToolUseLocked(toolCtx, input.CallID, input.Name, result, runErr)
					assistantMu.Unlock()
					if err != nil {
						return nil, err
					}
				}
				if runErr != nil {
					return nil, runErr
				}
				return output, nil
			}
		},
	}

	directTools := append([]tool.BaseTool{newWorkflowHandoffTool(setWorkflowHandoff)}, agentTools...)

	// 新建供应商（入口暴露 workflow handoff 工具 + 用户已选工具，工作流执行阶段仍使用 agentTools）
	provider, err := llm_provider.NewLlmProvider(ctx, *providerModel, subAgents, directTools, toolMiddleware, s.prompts)
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
		now := time.Now()
		runCtx, cancel := context.WithCancel(context.Background())
		defer cancel()
		var userStopped atomic.Bool
		go func() {
			<-userStop
			userStopped.Store(true)
			cancel()
		}()

		assistantMu.Lock()
		task.Status = data_models.TaskStatusRunning
		task.StartedAt = &now
		updateCurrentStageLocked("准备执行", "")
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
		}()

		failWithError := func(err error) {
			assistantMu.Lock()
			defer assistantMu.Unlock()
			if userStopped.Load() || errors.Is(runCtx.Err(), context.Canceled) {
				assistantMessage.AssistantMessageExtra.FinishReason = "user stop"
				assistantMessage.AssistantMessageExtra.FinishError = ""
				return
			}
			if err != nil {
				assistantMessage.AssistantMessageExtra.FinishReason = "error"
				assistantMessage.AssistantMessageExtra.FinishError = err.Error()
			}
		}

		appendAssistantContentLocked := func(content, reasoning string) error {
			assistantMessage.Content += content
			assistantMessage.ReasoningContent += reasoning
			lastOutputAt := time.Now()
			task.LastOutputAt = &lastOutputAt
			return persistAssistantSnapshotThrottledLocked(true)
		}

		resetDirectAssistantStateLocked := func() {
			resetDirectAssistantState(&assistantMessage)
		}

		runSinglePassEntry := func(userRequest string) error {
			entryCtx, entryCancel := context.WithCancel(runCtx)
			defer entryCancel()
			entryMessages := append([]schema.Message{{
				Role:    schema.System,
				Content: s.prompts.EntrySystem,
			}}, schemaMessages...)
			iter, err := provider.AgentCompletions(entryCtx, entryMessages)
			if err != nil {
				return err
			}
			for {
				if getWorkflowHandoff() != nil {
					entryCancel()
					break
				}

				event, ok := iter.Next()
				if !ok {
					break
				}
				if event.Err != nil {
					return event.Err
				}
				if event.Output == nil || event.Output.MessageOutput == nil {
					continue
				}
				mo := event.Output.MessageOutput
				if mo.Role == schema.Tool {
					if getWorkflowHandoff() != nil {
						entryCancel()
						break
					}
					continue
				}
				if mo.Role != schema.Assistant {
					continue
				}
				handleChunk := func(msg *schema.Message) error {
					if msg == nil {
						return nil
					}
					if strings.TrimSpace(msg.Content) == "" && strings.TrimSpace(msg.ReasoningContent) == "" {
						return nil
					}
					assistantMu.Lock()
					defer assistantMu.Unlock()
					return appendAssistantContentLocked(msg.Content, msg.ReasoningContent)
				}
				if mo.IsStreaming && mo.MessageStream != nil {
					for {
						msg, streamErr := mo.MessageStream.Recv()
						if streamErr == io.EOF {
							break
						}
						if streamErr != nil {
							return streamErr
						}
						if err := handleChunk(msg); err != nil {
							return err
						}
						if getWorkflowHandoff() != nil {
							entryCancel()
							break
						}
					}
					mo.MessageStream.Close()
					continue
				}
				if err := handleChunk(mo.Message); err != nil {
					return err
				}
			}
			return nil
		}

		executeWorkflow := func(handoff *workflowHandoff) error {
			userRequest := inputMessage.Content
			hasAttachedFiles := inputMessage.UserMessageExtra != nil && len(inputMessage.UserMessageExtra.Files) > 0
			originalUserMessage := findLatestUserContextMessage(schemaMessages)
			planningSummary := "PlannerAgent 正在生成执行计划"
			planningInputLabel := "任务请求"
			synthesisSummary := "SynthesizerAgent 正在整合所有子任务产出"
			reviewSummaryText := "ReviewerAgent 正在检查答案是否满足目标"
			if hasAttachedFiles {
				planningSummary = "PlannerAgent 正在基于原始多模态输入生成执行计划"
				planningInputLabel = "任务请求（沿用原始多模态输入）"
				synthesisSummary = "SynthesizerAgent 正在结合原始多模态输入整合结果"
				reviewSummaryText = "ReviewerAgent 正在结合原始多模态输入审核答案"
			}
			logWorkflowHandoff := func(stepID string, handoff workflowHandoff) error {
				assistantMessage.AssistantMessageExtra.RouteType = data_models.RouteTypeWorkflow
				summary := strings.TrimSpace(handoff.Summary)
				if summary == "" {
					summary = handoff.Reason
				}
				return startTraceStepLocked(stepID, "", data_models.TraceStepTypeClassify, "主模型交付 Workflow", summary, userRequest, "任务交付", "MainRouterAgent", []data_models.TraceDetailBlock{
					{Kind: "input", Title: "用户输入", Content: userRequest, Format: data_models.TraceDetailFormatText},
					{Kind: "review", Title: "交付结果", Content: fmt.Sprintf("{\"reason\":%q,\"summary\":%q,\"rule_name\":%q}", handoff.Reason, handoff.Summary, handoff.RuleName), Format: data_models.TraceDetailFormatJSON},
				}, map[string]interface{}{
					"route_source": handoff.Source,
					"rule_name":    handoff.RuleName,
				})
			}
			finishWorkflowHandoff := func(stepID string, handoff workflowHandoff) error {
				return finishTraceStepLocked(stepID, handoff.Reason, handoff.Summary, "任务交付", "MainRouterAgent", data_models.TraceStepStatusDone, []data_models.TraceDetailBlock{
					{Kind: "input", Title: "用户输入", Content: userRequest, Format: data_models.TraceDetailFormatText},
					{Kind: "review", Title: "交付结果", Content: fmt.Sprintf("{\"reason\":%q,\"summary\":%q,\"rule_name\":%q}", handoff.Reason, handoff.Summary, handoff.RuleName), Format: data_models.TraceDetailFormatJSON},
				}, map[string]interface{}{
					"route_source": handoff.Source,
					"rule_name":    handoff.RuleName,
				})
			}

			if handoff != nil {
				assistantMu.Lock()
				preserveWorkflowPreface(&assistantMessage)
				resetDirectAssistantStateLocked()
				err := persistAssistantSnapshotLocked(false)
				if err == nil {
					err = logWorkflowHandoff("workflow_handoff", *handoff)
				}
				if err == nil {
					err = finishWorkflowHandoff("workflow_handoff", *handoff)
				}
				assistantMu.Unlock()
				if err != nil {
					return err
				}
			}

			assistantMu.Lock()
			if err := startTraceStepLocked("plan", "", data_models.TraceStepTypePlan, "拆解任务", planningSummary, userRequest, "任务拆解", "PlannerAgent", []data_models.TraceDetailBlock{
				{Kind: "input", Title: planningInputLabel, Content: userRequest, Format: data_models.TraceDetailFormatText},
			}, nil); err != nil {
				assistantMu.Unlock()
				return err
			}
			assistantMu.Unlock()
			plan, err := generateWorkflowPlan(runCtx, provider, userRequest, schemaMessages, agentTools)
			if err != nil {
				return err
			}
			taskTitles := make([]string, 0, len(plan.Tasks))
			for _, item := range plan.Tasks {
				taskTitles = append(taskTitles, item.Title)
			}
			assistantMu.Lock()
			err = finishTraceStepLocked("plan", fmt.Sprintf("共拆分 %d 个子任务", len(plan.Tasks)), strings.Join(taskTitles, " | "), "任务拆解", "PlannerAgent", data_models.TraceStepStatusDone, []data_models.TraceDetailBlock{
				{Kind: "plan", Title: "完整计划", Content: formatWorkflowPlanForTrace(plan), Format: data_models.TraceDetailFormatMarkdown},
			}, map[string]interface{}{
				"goal":                plan.Goal,
				"completion_criteria": plan.CompletionCriteria,
			})
			assistantMu.Unlock()
			if err != nil {
				return err
			}

			results := map[string]workflowTaskResult{}
			executeBatches := func(filterIDs map[string]struct{}, retryInstructions string) error {
				batches := batchTasksByDependencies(plan.Tasks, filterIDs)
				assistantMu.Lock()
				retryCount := assistantMessage.AssistantMessageExtra.RetryCount
				assistantMu.Unlock()
				for batchIndex, batch := range batches {
					priorResults := make(map[string]workflowTaskResult, len(results))
					for key, value := range results {
						priorResults[key] = value
					}
					type taskRun struct {
						task   workflowPlanTask
						result workflowTaskResult
						err    error
					}
					resultCh := make(chan taskRun, len(batch))
					var wg sync.WaitGroup
					for _, baseTask := range batch {
						taskForRun := baseTask
						if strings.TrimSpace(retryInstructions) != "" {
							taskForRun.Description += "\n补充修正要求：" + retryInstructions
						}
						wg.Add(1)
						go func(taskItem workflowPlanTask, batchNo int) {
							defer wg.Done()
							dispatchStepID := fmt.Sprintf("dispatch_%s_%d", taskItem.ID, retryCount)
							assistantMu.Lock()
							dispatchErr := startTraceStepLocked(dispatchStepID, "", data_models.TraceStepTypeDispatch, "分派子任务", fmt.Sprintf("第 %d 批：%s", batchNo+1, taskItem.Title), taskItem.Description, "子任务执行", "MainRouterAgent", []data_models.TraceDetailBlock{
								{Kind: "plan", Title: "分派内容", Content: formatDispatchedTaskForTrace(taskItem, batchNo+1), Format: data_models.TraceDetailFormatMarkdown},
							}, map[string]interface{}{
								"task_id": taskItem.ID,
							})
							if dispatchErr == nil {
								dispatchErr = finishTraceStepLocked(dispatchStepID, "任务已分派", taskItem.Description, "子任务执行", "MainRouterAgent", data_models.TraceStepStatusDone, []data_models.TraceDetailBlock{
									{Kind: "plan", Title: "分派内容", Content: formatDispatchedTaskForTrace(taskItem, batchNo+1), Format: data_models.TraceDetailFormatMarkdown},
								}, map[string]interface{}{
									"task_id": taskItem.ID,
								})
							}
							assistantMu.Unlock()
							if dispatchErr != nil {
								resultCh <- taskRun{task: taskItem, err: dispatchErr}
								return
							}

							agentStepID := fmt.Sprintf("agent_%s_retry_%d", taskItem.ID, retryCount)
							assistantMu.Lock()
							agentErr := startTraceStepLocked(agentStepID, "", data_models.TraceStepTypeAgentRun, taskItem.Title, taskItem.Description, buildWorkerPrompt(plan, taskItem, priorResults), "子任务执行", taskItem.SuggestedAgent, buildAgentTraceDetails(plan, taskItem, priorResults, retryInstructions), map[string]interface{}{
								"task_id":         taskItem.ID,
								"expected_output": taskItem.ExpectedOutput,
							})
							assistantMu.Unlock()
							if agentErr != nil {
								resultCh <- taskRun{task: taskItem, err: agentErr}
								return
							}

							execResult, execErr := executePlanTask(runCtx, provider, taskItem, plan, priorResults, originalUserMessage, agentTools, toolMiddleware, agentStepID)
							assistantMu.Lock()
							finishStatus := data_models.TraceStepStatusDone
							summary := "子任务执行完成"
							outputPreview := compactText(execResult.Output, 240)
							metadata := map[string]interface{}{
								"task_id":    taskItem.ID,
								"used_tools": execResult.UsedTools,
							}
							if execErr != nil {
								finishStatus = data_models.TraceStepStatusError
								summary = execErr.Error()
							}
							agentErr = finishTraceStepLocked(agentStepID, summary, outputPreview, "子任务执行", taskItem.SuggestedAgent, finishStatus, buildAgentResultTraceDetails(plan, taskItem, priorResults, retryInstructions, execResult, execErr), metadata)
							assistantMu.Unlock()
							if agentErr != nil && execErr == nil {
								execErr = agentErr
							}
							resultCh <- taskRun{task: taskItem, result: execResult, err: execErr}
						}(taskForRun, batchIndex)
					}
					wg.Wait()
					close(resultCh)
					for item := range resultCh {
						if item.err != nil {
							return item.err
						}
						results[item.task.ID] = item.result
					}
				}
				return nil
			}

			if err := executeBatches(nil, ""); err != nil {
				return err
			}

			reviewFeedback := ""
			draft := ""
			review := reviewDecision{}
			for attempt := 0; attempt < 2; attempt++ {
				assistantMu.Lock()
				err := startTraceStepLocked(fmt.Sprintf("synthesize_%d", attempt), "", data_models.TraceStepTypeSynthesize, "汇总子任务结果", synthesisSummary, userRequest, "结果汇总", "SynthesizerAgent", []data_models.TraceDetailBlock{
					{Kind: "input", Title: planningInputLabel, Content: userRequest, Format: data_models.TraceDetailFormatText},
				}, nil)
				assistantMu.Unlock()
				if err != nil {
					return err
				}
				draft, err = synthesizeWorkflowAnswer(runCtx, provider, userRequest, originalUserMessage, plan, results, reviewFeedback)
				if err != nil {
					return err
				}
				assistantMu.Lock()
				err = finishTraceStepLocked(fmt.Sprintf("synthesize_%d", attempt), "已生成候选答案", compactText(draft, 240), "结果汇总", "SynthesizerAgent", data_models.TraceStepStatusDone, []data_models.TraceDetailBlock{
					{Kind: "output", Title: "候选答案", Content: draft, Format: data_models.TraceDetailFormatMarkdown},
				}, nil)
				assistantMu.Unlock()
				if err != nil {
					return err
				}

				assistantMu.Lock()
				err = startTraceStepLocked(fmt.Sprintf("review_%d", attempt), "", data_models.TraceStepTypeReview, "审核候选答案", reviewSummaryText, draft, "结果审核", "ReviewerAgent", []data_models.TraceDetailBlock{
					{Kind: "output", Title: "待审核答案", Content: draft, Format: data_models.TraceDetailFormatMarkdown},
				}, nil)
				assistantMu.Unlock()
				if err != nil {
					return err
				}
				review, err = reviewWorkflowAnswer(runCtx, provider, userRequest, originalUserMessage, plan, results, draft)
				if err != nil {
					return err
				}
				reviewSummary := "审核通过"
				reviewStatus := data_models.TraceStepStatusDone
				if !review.Approved {
					reviewSummary = strings.Join(review.Issues, "；")
					reviewStatus = data_models.TraceStepStatusError
				}
				assistantMu.Lock()
				err = finishTraceStepLocked(fmt.Sprintf("review_%d", attempt), reviewSummary, compactText(review.RetryInstructions, 240), "结果审核", "ReviewerAgent", reviewStatus, []data_models.TraceDetailBlock{
					{Kind: "review", Title: "审核结果", Content: formatReviewDecisionForTrace(review), Format: data_models.TraceDetailFormatJSON},
				}, map[string]interface{}{
					"approved":           review.Approved,
					"affected_task_ids":  review.AffectedTaskIDs,
					"retry_instructions": review.RetryInstructions,
				})
				assistantMu.Unlock()
				if err != nil {
					return err
				}
				if review.Approved {
					break
				}
				if attempt == 1 {
					break
				}

				retryStepID := fmt.Sprintf("retry_%d", attempt+1)
				assistantMu.Lock()
				assistantMessage.AssistantMessageExtra.RetryCount = attempt + 1
				err = startTraceStepLocked(retryStepID, "", data_models.TraceStepTypeRetry, fmt.Sprintf("第 %d 次修正", attempt+1), strings.Join(review.Issues, "；"), review.RetryInstructions, "重新生成", "MainRouterAgent", []data_models.TraceDetailBlock{
					{Kind: "retry", Title: "重试原因", Content: strings.Join(review.Issues, "\n"), Format: data_models.TraceDetailFormatMarkdown},
					{Kind: "retry", Title: "修正指令", Content: review.RetryInstructions, Format: data_models.TraceDetailFormatText},
				}, map[string]interface{}{
					"retry_instructions": review.RetryInstructions,
				})
				assistantMu.Unlock()
				if err != nil {
					return err
				}
				filterIDs := map[string]struct{}{}
				if len(review.AffectedTaskIDs) > 0 {
					for _, taskID := range review.AffectedTaskIDs {
						filterIDs[taskID] = struct{}{}
					}
				}
				if len(filterIDs) == 0 {
					filterIDs = nil
				}
				if err := executeBatches(filterIDs, review.RetryInstructions); err != nil {
					return err
				}
				reviewFeedback = strings.Join(review.Issues, "；") + "\n" + review.RetryInstructions
				assistantMu.Lock()
				err = finishTraceStepLocked(retryStepID, "已完成定向重试", review.RetryInstructions, "重新生成", "MainRouterAgent", data_models.TraceStepStatusDone, []data_models.TraceDetailBlock{
					{Kind: "retry", Title: "重试结果", Content: formatRetrySummaryForTrace(review), Format: data_models.TraceDetailFormatMarkdown},
				}, map[string]interface{}{
					"affected_task_ids": review.AffectedTaskIDs,
				})
				assistantMu.Unlock()
				if err != nil {
					return err
				}
			}

			finalSummary := "答案已通过审核"
			if !review.Approved {
				finalSummary = "达到最大重试次数，输出当前最佳答案"
				draft = strings.TrimSpace(draft + "\n\n注意：系统已进行一次自动修正，但仍建议你根据上面的内容做最终确认。")
			}

			assistantMu.Lock()
			err = startTraceStepLocked("finalize_workflow", "", data_models.TraceStepTypeFinalize, "输出最终答案", finalSummary, draft, "已完成", "MainRouterAgent", []data_models.TraceDetailBlock{
				{Kind: "output", Title: "最终答案草稿", Content: draft, Format: data_models.TraceDetailFormatMarkdown},
			}, nil)
			if err == nil {
				assistantMessage.Content = draft
				err = persistAssistantSnapshotLocked(true)
			}
			if err == nil {
				err = finishTraceStepLocked("finalize_workflow", finalSummary, compactText(draft, 240), "已完成", "MainRouterAgent", data_models.TraceStepStatusDone, []data_models.TraceDetailBlock{
					{Kind: "output", Title: "最终答案", Content: draft, Format: data_models.TraceDetailFormatMarkdown},
				}, map[string]interface{}{
					"approved": review.Approved,
				})
			}
			assistantMu.Unlock()
			return err
		}

		if guard := shouldForceWorkflow(inputMessage); guard.Force {
			if err := executeWorkflow(&workflowHandoff{
				Reason:   guard.Reason,
				Summary:  guard.Reason,
				Source:   routeSourceGuardRule,
				RuleName: guard.RuleName,
			}); err != nil {
				failWithError(err)
				return
			}
		} else {
			assistantMu.Lock()
			resetDirectAssistantStateLocked()
			saveErr := persistAssistantSnapshotLocked(true)
			assistantMu.Unlock()
			if saveErr != nil {
				failWithError(saveErr)
				return
			}

			if err := runSinglePassEntry(inputMessage.Content); err != nil {
				failWithError(err)
				return
			}

			handoff := getWorkflowHandoff()
			if handoff != nil {
				if err := executeWorkflow(handoff); err != nil {
					failWithError(err)
					return
				}
			} else {
				assistantMu.Lock()
				assistantMessage.Content = strings.TrimSpace(assistantMessage.Content)
				if assistantMessage.Content == "" {
					assistantMessage.Content = "抱歉，我暂时没有生成内容，请重试。"
				}
				saveErr = persistAssistantSnapshotLocked(true)
				assistantMu.Unlock()
				if saveErr != nil {
					failWithError(saveErr)
					return
				}
			}
		}
		assistantMu.Lock()
		assistantMessage.AssistantMessageExtra.FinishReason = "done"
		assistantMessage.AssistantMessageExtra.FinishError = ""
		assistantMu.Unlock()
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
	for {
		task, err := s.storage.GetChatActiveTask(ctx, chatUuid)
		if err != nil {
			return nil, ierror.NewError(err)
		}
		if task == nil {
			return nil, nil
		}
		task, err = s.repairStaleActiveTask(ctx, task)
		if err != nil {
			return nil, ierror.NewError(err)
		}
		if task == nil {
			continue
		}
		viewTask := view_models.Task(*task)
		return &viewTask, nil
	}
}

func (s *Service) GetRunningTasks(ctx context.Context) (*view_models.TaskList, error) {
	tasks, err := s.storage.GetRunningTasks(ctx)
	if err != nil {
		return nil, ierror.NewError(err)
	}
	viewTasks := make([]view_models.Task, 0, len(tasks))
	for _, task := range tasks {
		liveTask, err := s.repairStaleActiveTask(ctx, &task)
		if err != nil {
			return nil, ierror.NewError(err)
		}
		if liveTask == nil {
			continue
		}
		viewTasks = append(viewTasks, view_models.Task(*liveTask))
	}
	return &view_models.TaskList{Tasks: viewTasks}, nil
}

func (s *Service) emitTaskEvent(task data_models.Task, assistantMessage data_models.Message, traceDelta []data_models.TraceStep) {
	s.app.Event.Emit(task.EventKey, view_models.TaskStreamEvent{
		TaskUuid:         task.TaskUuid,
		ChatUuid:         task.ChatUuid,
		EventKey:         task.EventKey,
		Status:           task.Status,
		FinishReason:     task.FinishReason,
		FinishError:      task.FinishError,
		ExecutionTrace:   assistantMessage.AssistantMessageExtra.ExecutionTrace,
		TraceDelta:       traceDelta,
		CurrentStage:     assistantMessage.AssistantMessageExtra.CurrentStage,
		CurrentAgent:     assistantMessage.AssistantMessageExtra.CurrentAgent,
		RetryCount:       assistantMessage.AssistantMessageExtra.RetryCount,
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
	provider, err := llm_provider.NewLlmProvider(ctx, providerModel, []adk.Agent{}, []tool.BaseTool{}, compose.ToolMiddleware{}, s.prompts)
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
		payload := struct {
			ChatUuid string `json:"chat_uuid"`
			Title    string `json:"title"`
		}{
			ChatUuid: chatUuid,
			Title:    title,
		}
		s.app.Event.Emit(event.GenEventsKey(event.EventTypeChatTitle, chatUuid), payload)
		s.app.Event.Emit(event.GenEventsKey(event.EventTypeChatTitle, "all"), payload)
	}

	return title, nil
}
