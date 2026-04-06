package service

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/llm_opc"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/llm_provider"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/llm_provider/agents"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils/ierror"
)

// ==================== OPC Avatar Selection ====================

func (s *Service) OPCSelectAvatar() (string, error) {
	pattern := "*.png;*.jpg;*.jpeg;*.gif;*.webp"
	path, err := s.app.Dialog.OpenFile().SetTitle("选择头像").AddFilter("图片文件", pattern).PromptForSingleSelection()
	if err != nil {
		return "", ierror.NewError(fmt.Errorf("failed to open file dialog: %w", err))
	}
	if path == "" {
		return "", nil
	}

	dataPath, err := utils.GetDataPath()
	if err != nil {
		return "", ierror.NewError(fmt.Errorf("failed to get data path: %w", err))
	}

	avatarDir := filepath.Join(dataPath, "opc_avatars")
	if err := os.MkdirAll(avatarDir, 0755); err != nil {
		return "", ierror.NewError(fmt.Errorf("failed to create avatar directory: %w", err))
	}

	ext := filepath.Ext(path)
	destName := uuid.New().String() + ext
	destPath := filepath.Join(avatarDir, destName)

	srcFile, err := os.Open(path)
	if err != nil {
		return "", ierror.NewError(fmt.Errorf("failed to open source file: %w", err))
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return "", ierror.NewError(fmt.Errorf("failed to create destination file: %w", err))
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return "", ierror.NewError(fmt.Errorf("failed to copy file: %w", err))
	}

	return "image:" + destPath, nil
}

// ==================== OPC Person CRUD ====================

func (s *Service) OPCCreatePerson(ctx context.Context, input view_models.OPCPersonInput) (*view_models.OPCPersonView, error) {
	if strings.TrimSpace(input.Name) == "" {
		return nil, ierror.NewError(fmt.Errorf("person name cannot be empty"))
	}

	personUuid := uuid.New().String()
	agentID := fmt.Sprintf("opc_person_%s", personUuid)
	chatUuid := uuid.New().String()

	// 创建 CustomAgentDef
	now := time.Now().Format(time.RFC3339)
	toolIDs := input.Tools
	if toolIDs == nil {
		toolIDs = []string{}
	}
	skillIDs := input.Skills
	if skillIDs == nil {
		skillIDs = []string{}
	}

	agentDef := agents.CustomAgentDef{
		ID_:         agentID,
		DisplayName: input.Name,
		Description: fmt.Sprintf("OPC Person: %s - %s", input.Name, input.Role),
		PromptText:  input.Prompt,
		ToolIDs:     toolIDs,
		SkillIDs:    skillIDs,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	if err := agents.SaveCustomAgent(agentDef); err != nil {
		return nil, ierror.NewError(err)
	}
	agents.SyncCustomAgentsToRegistry()

	// 创建 Chat
	if err := s.storage.CreateChatWithType(ctx, chatUuid, input.Name, data_models.ChatTypeOPCPerson); err != nil {
		return nil, ierror.NewError(err)
	}

	// 创建 OPCPerson
	person := &data_models.OPCPerson{
		Uuid:     personUuid,
		Name:     input.Name,
		Role:     input.Role,
		AgentID:  agentID,
		Avatar:   input.Avatar,
		ChatUuid: chatUuid,
	}
	if err := s.storage.CreateOPCPerson(ctx, person); err != nil {
		return nil, ierror.NewError(err)
	}

	return s.buildPersonView(ctx, person)
}

func (s *Service) OPCListPersons(ctx context.Context) ([]view_models.OPCPersonView, error) {
	persons, err := s.storage.GetOPCPersons(ctx)
	if err != nil {
		return nil, ierror.NewError(err)
	}

	var views []view_models.OPCPersonView
	for i := range persons {
		view, err := s.buildPersonView(ctx, &persons[i])
		if err != nil {
			continue
		}
		views = append(views, *view)
	}
	return views, nil
}

func (s *Service) OPCGetPerson(ctx context.Context, personUuid string) (*view_models.OPCPersonView, error) {
	person, err := s.storage.GetOPCPerson(ctx, personUuid)
	if err != nil {
		return nil, ierror.NewError(err)
	}
	if person == nil {
		return nil, ierror.NewError(fmt.Errorf("person not found: %s", personUuid))
	}
	return s.buildPersonView(ctx, person)
}

func (s *Service) OPCUpdatePerson(ctx context.Context, input view_models.OPCPersonInput) (*view_models.OPCPersonView, error) {
	if strings.TrimSpace(input.Name) == "" {
		return nil, ierror.NewError(fmt.Errorf("person name cannot be empty"))
	}

	person, err := s.storage.GetOPCPerson(ctx, input.Uuid)
	if err != nil {
		return nil, ierror.NewError(err)
	}
	if person == nil {
		return nil, ierror.NewError(fmt.Errorf("person not found: %s", input.Uuid))
	}

	// 更新 CustomAgentDef
	oldDef, loadErr := agents.LoadCustomAgent(person.AgentID)
	if loadErr != nil {
		return nil, ierror.NewError(loadErr)
	}

	toolIDs := input.Tools
	if toolIDs == nil {
		toolIDs = []string{}
	}
	skillIDs := input.Skills
	if skillIDs == nil {
		skillIDs = []string{}
	}

	updatedDef := agents.CustomAgentDef{
		ID_:         person.AgentID,
		DisplayName: input.Name,
		Description: fmt.Sprintf("OPC Person: %s - %s", input.Name, input.Role),
		PromptText:  input.Prompt,
		ToolIDs:     toolIDs,
		SkillIDs:    skillIDs,
		CreatedAt:   oldDef.CreatedAt,
		UpdatedAt:   time.Now().Format(time.RFC3339),
	}
	if err := agents.SaveCustomAgent(updatedDef); err != nil {
		return nil, ierror.NewError(err)
	}
	agents.SyncCustomAgentsToRegistry()

	// 更新 OPCPerson
	person.Name = input.Name
	person.Role = input.Role
	person.Avatar = input.Avatar
	if err := s.storage.UpdateOPCPerson(ctx, person); err != nil {
		return nil, ierror.NewError(err)
	}

	// 更新 Chat 标题
	_ = s.storage.RenameChat(ctx, person.ChatUuid, input.Name)

	return s.buildPersonView(ctx, person)
}

func (s *Service) OPCDeletePerson(ctx context.Context, personUuid string) error {
	person, err := s.storage.GetOPCPerson(ctx, personUuid)
	if err != nil {
		return ierror.NewError(err)
	}
	if person == nil {
		return ierror.NewError(fmt.Errorf("person not found: %s", personUuid))
	}

	// 从所有群组中移除
	if err := s.storage.DeleteOPCGroupMembersByPerson(ctx, personUuid); err != nil {
		return ierror.NewError(err)
	}

	// 删除关联的 agent
	_ = agents.DeleteCustomAgent(person.AgentID)
	agents.SyncCustomAgentsToRegistry()

	// 删除 1v1 聊天和消息
	if err := s.storage.DeleteChatAndMessages(ctx, person.ChatUuid); err != nil {
		return ierror.NewError(err)
	}

	// 删除人员记录
	if err := s.storage.DeleteOPCPerson(ctx, personUuid); err != nil {
		return ierror.NewError(err)
	}

	return nil
}

func (s *Service) OPCTogglePinPerson(ctx context.Context, personUuid string, pinned bool) error {
	return s.storage.TogglePinOPCPerson(ctx, personUuid, pinned)
}

// OPCClearConversation 清除某个聊天的所有消息（不删除联系人或群聊）
func (s *Service) OPCClearConversation(ctx context.Context, chatUuid string) error {
	if chatUuid == "" {
		return ierror.NewError(fmt.Errorf("chat_uuid cannot be empty"))
	}
	return s.storage.ClearChatMessages(ctx, chatUuid)
}

// ==================== OPC Group CRUD ====================

func (s *Service) OPCCreateGroup(ctx context.Context, input view_models.OPCGroupInput) (*view_models.OPCGroupView, error) {
	if strings.TrimSpace(input.Name) == "" {
		return nil, ierror.NewError(fmt.Errorf("group name cannot be empty"))
	}

	groupUuid := uuid.New().String()
	chatUuid := uuid.New().String()

	// 创建 Chat
	if err := s.storage.CreateChatWithType(ctx, chatUuid, input.Name, data_models.ChatTypeOPCGroup); err != nil {
		return nil, ierror.NewError(err)
	}

	// 创建 OPCGroup
	group := &data_models.OPCGroup{
		Uuid:        groupUuid,
		ChatUuid:    chatUuid,
		Name:        input.Name,
		Description: input.Description,
	}
	if err := s.storage.CreateOPCGroup(ctx, group); err != nil {
		return nil, ierror.NewError(err)
	}

	// 创建成员关系
	if len(input.MemberUuids) > 0 {
		var members []data_models.OPCGroupMember
		for _, personUuid := range input.MemberUuids {
			members = append(members, data_models.OPCGroupMember{
				GroupUuid:  groupUuid,
				PersonUuid: personUuid,
			})
		}
		if err := s.storage.CreateOPCGroupMembers(ctx, members); err != nil {
			return nil, ierror.NewError(err)
		}
	}

	return s.buildGroupView(ctx, group)
}

func (s *Service) OPCListGroups(ctx context.Context) ([]view_models.OPCGroupView, error) {
	groups, err := s.storage.GetOPCGroups(ctx)
	if err != nil {
		return nil, ierror.NewError(err)
	}

	var views []view_models.OPCGroupView
	for i := range groups {
		view, err := s.buildGroupView(ctx, &groups[i])
		if err != nil {
			continue
		}
		views = append(views, *view)
	}
	return views, nil
}

func (s *Service) OPCGetGroup(ctx context.Context, groupUuid string) (*view_models.OPCGroupView, error) {
	group, err := s.storage.GetOPCGroup(ctx, groupUuid)
	if err != nil {
		return nil, ierror.NewError(err)
	}
	if group == nil {
		return nil, ierror.NewError(fmt.Errorf("group not found: %s", groupUuid))
	}
	return s.buildGroupView(ctx, group)
}

func (s *Service) OPCUpdateGroup(ctx context.Context, input view_models.OPCGroupInput) (*view_models.OPCGroupView, error) {
	if strings.TrimSpace(input.Name) == "" {
		return nil, ierror.NewError(fmt.Errorf("group name cannot be empty"))
	}

	group, err := s.storage.GetOPCGroup(ctx, input.Uuid)
	if err != nil {
		return nil, ierror.NewError(err)
	}
	if group == nil {
		return nil, ierror.NewError(fmt.Errorf("group not found: %s", input.Uuid))
	}

	group.Name = input.Name
	group.Description = input.Description
	if err := s.storage.UpdateOPCGroup(ctx, group); err != nil {
		return nil, ierror.NewError(err)
	}

	// 更新 Chat 标题
	_ = s.storage.RenameChat(ctx, group.ChatUuid, input.Name)

	// 更新成员关系
	if err := s.storage.DeleteOPCGroupMembersByGroup(ctx, group.Uuid); err != nil {
		return nil, ierror.NewError(err)
	}
	if len(input.MemberUuids) > 0 {
		var members []data_models.OPCGroupMember
		for _, personUuid := range input.MemberUuids {
			members = append(members, data_models.OPCGroupMember{
				GroupUuid:  group.Uuid,
				PersonUuid: personUuid,
			})
		}
		if err := s.storage.CreateOPCGroupMembers(ctx, members); err != nil {
			return nil, ierror.NewError(err)
		}
	}

	return s.buildGroupView(ctx, group)
}

func (s *Service) OPCDeleteGroup(ctx context.Context, groupUuid string) error {
	group, err := s.storage.GetOPCGroup(ctx, groupUuid)
	if err != nil {
		return ierror.NewError(err)
	}
	if group == nil {
		return ierror.NewError(fmt.Errorf("group not found: %s", groupUuid))
	}

	// 删除成员关系
	if err := s.storage.DeleteOPCGroupMembersByGroup(ctx, groupUuid); err != nil {
		return ierror.NewError(err)
	}

	// 删除关联聊天和消息
	if err := s.storage.DeleteChatAndMessages(ctx, group.ChatUuid); err != nil {
		return ierror.NewError(err)
	}

	// 删除群组记录
	if err := s.storage.DeleteOPCGroup(ctx, groupUuid); err != nil {
		return ierror.NewError(err)
	}

	return nil
}

func (s *Service) OPCTogglePinGroup(ctx context.Context, groupUuid string, pinned bool) error {
	return s.storage.TogglePinOPCGroup(ctx, groupUuid, pinned)
}

// ==================== OPC Person Chat ====================

// OPCPersonChat 人员私聊：发送消息给指定人员，非流式返回完整回复
func (s *Service) OPCPersonChat(ctx context.Context, input view_models.OPCCompletionInput) (*view_models.Completions, error) {
	person, err := s.storage.GetOPCPerson(ctx, input.PersonUuid)
	if err != nil {
		return nil, ierror.NewError(err)
	}
	if person == nil {
		return nil, ierror.NewError(fmt.Errorf("person not found: %s", input.PersonUuid))
	}

	chatUuid := input.ChatUuid
	if chatUuid == "" {
		chatUuid = person.ChatUuid
	}

	// 获取模型信息
	providerModel, err := s.storage.GetProviderModel(ctx, input.ModelId, input.ModelName)
	if err != nil {
		return nil, ierror.NewError(err)
	}
	if providerModel == nil {
		return nil, ierror.New(ierror.ErrCodeModelNotFound)
	}

	// 保存用户消息
	userMsgUuid := uuid.New().String()
	userMsg := data_models.Message{
		ChatUuid:    chatUuid,
		MessageUuid: userMsgUuid,
		Role:        schema.User,
		Content:     input.Content,
		UserMessageExtra: &data_models.UserMessageExtra{
			ModelId:   input.ModelId,
			ModelName: input.ModelName,
		},
	}
	if _, err := s.storage.CreateMessage(ctx, chatUuid, userMsg); err != nil {
		return nil, ierror.NewError(err)
	}

	// 创建助手消息（占位）
	assistantMsgUuid := uuid.New().String()
	taskUuid := uuid.New().String()
	eventKey := fmt.Sprintf("event:opc_person:%s", taskUuid)

	assistantMsg := data_models.Message{
		ChatUuid:         chatUuid,
		MessageUuid:      assistantMsgUuid,
		Role:             schema.Assistant,
		SenderPersonUuid: person.Uuid,
		AssistantMessageExtra: &data_models.AssistantMessageExtra{
			CurrentStage: "opc.stage.thinking",
		},
	}
	assistantMsgId, err := s.storage.CreateMessage(ctx, chatUuid, assistantMsg)
	if err != nil {
		return nil, ierror.NewError(err)
	}
	assistantMsg.ID = assistantMsgId

	// 创建 task
	task := data_models.Task{
		TaskUuid:             taskUuid,
		ChatUuid:             chatUuid,
		AssistantMessageUuid: assistantMsgUuid,
		Status:               data_models.TaskStatusPending,
		EventKey:             eventKey,
	}
	if err := s.storage.CreateTask(ctx, task); err != nil {
		return nil, ierror.NewError(err)
	}

	// 启动异步生成任务
	personCopy := *person
	go func() {
		ctx := context.Background()
		now := time.Now()
		task.Status = data_models.TaskStatusRunning
		task.StartedAt = &now
		_ = s.storage.SaveTask(ctx, task)

		// 通知前端：开始输入
		s.app.Event.Emit(eventKey, map[string]interface{}{
			"type":        "opc:typing",
			"person_uuid": personCopy.Uuid,
			"person_name": personCopy.Name,
		})

		err := s.opcGeneratePersonReply(ctx, &personCopy, chatUuid, &assistantMsg, providerModel.ModelId, providerModel.Model)
		finishNow := time.Now()
		if err != nil {
			logger.Errorf("OPC person chat failed: %v", err)
			task.Status = data_models.TaskStatusFailed
			task.FinishError = err.Error()
		} else {
			task.Status = data_models.TaskStatusCompleted
		}
		task.FinishedAt = &finishNow
		_ = s.storage.SaveTask(ctx, task)

		// 通知前端：消息完成
		s.app.Event.Emit(eventKey, map[string]interface{}{
			"type":    "opc:message",
			"message": assistantMsg,
		})
		s.app.Event.Emit(eventKey, map[string]interface{}{
			"type": "opc:complete",
		})
	}()

	return &view_models.Completions{
		ChatUuid:    chatUuid,
		TaskUuid:    taskUuid,
		MessageUuid: assistantMsgUuid,
		EventKey:    eventKey,
	}, nil
}

// opcGeneratePersonReply 为指定人员生成回复，委托给 llm_opc 包
func (s *Service) opcGeneratePersonReply(ctx context.Context, person *data_models.OPCPerson, chatUuid string, assistantMsg *data_models.Message, modelId uint, modelName string) error {
	// 加载历史消息
	historyMessages, _, err := s.storage.GetMessage(ctx, chatUuid, 0, 20)
	if err != nil {
		return err
	}

	// 构建 IOpcPerson
	opcPerson, err := s.buildOpcPerson(person)
	if err != nil {
		return err
	}

	// 获取模型
	providerModel, err := s.storage.GetProviderModel(ctx, modelId, modelName)
	if err != nil || providerModel == nil {
		return fmt.Errorf("model not available")
	}

	chatModel, err := llm_provider.NewToolCallingChatModel(ctx, *providerModel)
	if err != nil {
		return err
	}

	// 委托给 llm_opc 生成回复
	result, err := llm_opc.GeneratePersonReply(ctx, llm_opc.PersonChatParams{
		Person:          opcPerson,
		ChatModel:       chatModel,
		HistoryMessages: historyMessages,
		SkipMessageUuid: assistantMsg.MessageUuid,
	})
	if err != nil {
		return err
	}

	// 更新助手消息
	assistantMsg.Content = result.Content
	assistantMsg.ReasoningContent = result.ReasoningContent
	if assistantMsg.AssistantMessageExtra == nil {
		assistantMsg.AssistantMessageExtra = &data_models.AssistantMessageExtra{}
	}
	assistantMsg.AssistantMessageExtra.CurrentStage = "opc.stage.completed"

	return s.storage.SaveOrUpdateMessage(ctx, *assistantMsg)
}

// ==================== OPC Group Chat ====================

// OPCGroupChat 群聊：发送消息，agents 自主决定是否回复
func (s *Service) OPCGroupChat(ctx context.Context, input view_models.OPCGroupCompletionInput) (*view_models.Completions, error) {
	group, err := s.storage.GetOPCGroup(ctx, input.GroupUuid)
	if err != nil {
		return nil, ierror.NewError(err)
	}
	if group == nil {
		return nil, ierror.NewError(fmt.Errorf("group not found: %s", input.GroupUuid))
	}

	chatUuid := input.ChatUuid
	if chatUuid == "" {
		chatUuid = group.ChatUuid
	}

	// 获取模型信息
	providerModel, err := s.storage.GetProviderModel(ctx, input.ModelId, input.ModelName)
	if err != nil {
		return nil, ierror.NewError(err)
	}
	if providerModel == nil {
		return nil, ierror.New(ierror.ErrCodeModelNotFound)
	}

	// 保存用户消息
	userMsgUuid := uuid.New().String()
	userMsg := data_models.Message{
		ChatUuid:    chatUuid,
		MessageUuid: userMsgUuid,
		Role:        schema.User,
		Content:     input.Content,
		UserMessageExtra: &data_models.UserMessageExtra{
			ModelId:   input.ModelId,
			ModelName: input.ModelName,
		},
	}
	if _, err := s.storage.CreateMessage(ctx, chatUuid, userMsg); err != nil {
		return nil, ierror.NewError(err)
	}

	taskUuid := uuid.New().String()
	eventKey := fmt.Sprintf("event:opc_group:%s", taskUuid)

	// 创建 task
	task := data_models.Task{
		TaskUuid: taskUuid,
		ChatUuid: chatUuid,
		Status:   data_models.TaskStatusPending,
		EventKey: eventKey,
	}
	if err := s.storage.CreateTask(ctx, task); err != nil {
		return nil, ierror.NewError(err)
	}

	// 异步执行群聊逻辑（使用 background context，避免请求 context 取消导致失败）
	go s.opcGroupChatAsync(context.Background(), group, chatUuid, taskUuid, eventKey, input)

	return &view_models.Completions{
		ChatUuid: chatUuid,
		TaskUuid: taskUuid,
		EventKey: eventKey,
	}, nil
}

func (s *Service) opcGroupChatAsync(ctx context.Context, group *data_models.OPCGroup, chatUuid, taskUuid, eventKey string, input view_models.OPCGroupCompletionInput) {
	// 更新 task 状态
	now := time.Now()
	runningTask := data_models.Task{TaskUuid: taskUuid, Status: data_models.TaskStatusRunning, StartedAt: &now}
	_ = s.storage.SaveTask(ctx, runningTask)

	// 获取群成员并构建 IOpcPerson 列表
	members, err := s.resolveGroupMembers(ctx, group.Uuid)
	if err != nil {
		logger.Errorf("OPC group chat: get members failed: %v", err)
		s.opcGroupChatFinish(ctx, taskUuid, eventKey, "error", err.Error())
		return
	}

	if len(members) == 0 {
		s.opcGroupChatFinish(ctx, taskUuid, eventKey, "completed", "")
		return
	}

	// 获取历史消息（最后10条作为上下文）
	historyMessages, _, err := s.storage.GetMessage(ctx, chatUuid, 0, 10)
	if err != nil {
		logger.Errorf("OPC group chat: get history failed: %v", err)
		s.opcGroupChatFinish(ctx, taskUuid, eventKey, "error", err.Error())
		return
	}

	// 获取模型
	providerModel, err := s.storage.GetProviderModel(ctx, input.ModelId, input.ModelName)
	if err != nil || providerModel == nil {
		logger.Errorf("OPC group chat: get model failed: %v", err)
		s.opcGroupChatFinish(ctx, taskUuid, eventKey, "error", "model not available")
		return
	}

	chatModel, err := llm_provider.NewToolCallingChatModel(ctx, *providerModel)
	if err != nil {
		logger.Errorf("OPC group chat: create model failed: %v", err)
		s.opcGroupChatFinish(ctx, taskUuid, eventKey, "error", err.Error())
		return
	}

	// 构建上下文
	memberRoster := llm_opc.BuildMemberRoster(members)
	contextSummary := llm_opc.BuildContextSummary(historyMessages, llm_opc.MemberLookupFunc(members))

	// 决策轮：并行判断哪些成员应该回复
	decisions, err := llm_opc.DecideResponders(ctx, llm_opc.GroupDecisionParams{
		Members:        members,
		ChatModel:      chatModel,
		ContextSummary: contextSummary,
		MemberRoster:   memberRoster,
	})
	if err != nil {
		logger.Errorf("OPC group chat: decision round failed: %v", err)
		s.opcGroupChatFinish(ctx, taskUuid, eventKey, "error", err.Error())
		return
	}

	// 回复轮：按优先级顺序生成回复
	for _, decision := range decisions {
		person := decision.Person

		// 通知前端：某人开始输入
		s.app.Event.Emit(eventKey, map[string]interface{}{
			"type":        "opc:typing",
			"person_uuid": person.Uuid(),
			"person_name": person.Name(),
		})

		// 重新获取最新历史消息（包括其他成员的回复）
		latestMessages, _, _ := s.storage.GetMessage(ctx, chatUuid, 0, 20)

		replyModel, err := llm_provider.NewToolCallingChatModel(ctx, *providerModel)
		if err != nil {
			continue
		}

		result, err := llm_opc.GenerateGroupReply(ctx, llm_opc.GroupReplyParams{
			Person:          person,
			GroupName:       group.Name,
			GroupDesc:       group.Description,
			AllMembers:      members,
			ChatModel:       replyModel,
			HistoryMessages: latestMessages,
		})
		if err != nil {
			logger.Warm("OPC group chat: generate reply failed:", err)
			continue
		}

		// 保存回复消息
		replyMsgUuid := uuid.New().String()
		replyMsg := data_models.Message{
			ChatUuid:         chatUuid,
			MessageUuid:      replyMsgUuid,
			Role:             schema.Assistant,
			Content:          result.Content,
			ReasoningContent: result.ReasoningContent,
			SenderPersonUuid: person.Uuid(),
		}
		if _, err := s.storage.CreateMessage(ctx, chatUuid, replyMsg); err != nil {
			logger.Warm("OPC group chat: save reply failed:", err)
			continue
		}

		// 通知前端：消息完成
		s.app.Event.Emit(eventKey, map[string]interface{}{
			"type":    "opc:message",
			"message": replyMsg,
		})
	}

	s.opcGroupChatFinish(ctx, taskUuid, eventKey, "completed", "")
}

func (s *Service) opcGroupChatFinish(ctx context.Context, taskUuid, eventKey, status, errMsg string) {
	now := time.Now()
	taskStatus := data_models.TaskStatusCompleted
	finishErr := ""
	if status == "error" {
		taskStatus = data_models.TaskStatusFailed
		finishErr = errMsg
	}
	_ = s.storage.SaveTask(ctx, data_models.Task{
		TaskUuid:    taskUuid,
		Status:      taskStatus,
		FinishedAt:  &now,
		FinishError: finishErr,
	})

	s.app.Event.Emit(eventKey, map[string]interface{}{
		"type":  "opc:complete",
		"error": errMsg,
	})
}

// ==================== OPC Search ====================

func (s *Service) OPCSearch(ctx context.Context, keyword string) ([]view_models.OPCPersonView, []view_models.OPCGroupView, error) {
	persons, err := s.storage.SearchOPCPersons(ctx, keyword)
	if err != nil {
		return nil, nil, ierror.NewError(err)
	}

	groups, err := s.storage.SearchOPCGroups(ctx, keyword)
	if err != nil {
		return nil, nil, ierror.NewError(err)
	}

	var personViews []view_models.OPCPersonView
	for i := range persons {
		view, err := s.buildPersonView(ctx, &persons[i])
		if err != nil {
			continue
		}
		personViews = append(personViews, *view)
	}

	var groupViews []view_models.OPCGroupView
	for i := range groups {
		view, err := s.buildGroupView(ctx, &groups[i])
		if err != nil {
			continue
		}
		groupViews = append(groupViews, *view)
	}

	return personViews, groupViews, nil
}

// ==================== Helper methods ====================

func (s *Service) buildPersonView(ctx context.Context, person *data_models.OPCPerson) (*view_models.OPCPersonView, error) {
	view := &view_models.OPCPersonView{
		OPCPerson: *person,
	}

	// 获取最后一条消息
	lastMsg, err := s.storage.GetLastMessage(ctx, person.ChatUuid)
	if err == nil && lastMsg != nil {
		content := truncateString(lastMsg.Content, 50)
		view.LastMessage = &content
		view.LastMessageAt = &lastMsg.CreatedAt
	}

	return view, nil
}

func (s *Service) buildGroupView(ctx context.Context, group *data_models.OPCGroup) (*view_models.OPCGroupView, error) {
	view := &view_models.OPCGroupView{
		OPCGroup: *group,
	}

	// 获取成员
	memberRecords, err := s.storage.GetOPCGroupMembers(ctx, group.Uuid)
	if err == nil {
		for _, m := range memberRecords {
			person, err := s.storage.GetOPCPerson(ctx, m.PersonUuid)
			if err != nil || person == nil {
				continue
			}
			personView, err := s.buildPersonView(ctx, person)
			if err != nil {
				continue
			}
			view.Members = append(view.Members, *personView)
		}
	}

	// 获取最后一条消息
	lastMsg, err := s.storage.GetLastMessage(ctx, group.ChatUuid)
	if err == nil && lastMsg != nil {
		content := truncateString(lastMsg.Content, 50)
		view.LastMessage = &content
		view.LastMessageAt = &lastMsg.CreatedAt
	}

	return view, nil
}

func truncateString(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	return string(runes[:maxLen]) + "..."
}

// buildOpcPerson 从存储中加载人员并构建 IOpcPerson
func (s *Service) buildOpcPerson(person *data_models.OPCPerson) (llm_opc.IOpcPerson, error) {
	agentDef, ok := agents.FindAgent(person.AgentID)
	if !ok {
		return nil, fmt.Errorf("agent not found: %s", person.AgentID)
	}
	customDef, ok := agentDef.(*agents.CustomAgentDef)
	if !ok {
		return nil, fmt.Errorf("agent is not a custom agent: %s", person.AgentID)
	}
	return llm_opc.NewOpcPerson(person, customDef), nil
}

// resolveGroupMembers 获取群组所有成员并构建 IOpcPerson 列表
func (s *Service) resolveGroupMembers(ctx context.Context, groupUuid string) ([]llm_opc.IOpcPerson, error) {
	memberRecords, err := s.storage.GetOPCGroupMembers(ctx, groupUuid)
	if err != nil {
		return nil, err
	}

	var members []llm_opc.IOpcPerson
	for _, m := range memberRecords {
		person, err := s.storage.GetOPCPerson(ctx, m.PersonUuid)
		if err != nil || person == nil {
			continue
		}
		opcPerson, err := s.buildOpcPerson(person)
		if err != nil {
			continue
		}
		members = append(members, opcPerson)
	}
	return members, nil
}
