package llm_opc

import (
	"context"
	"encoding/json"
	"sort"
	"strings"
	"sync"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"
)

// GroupDecisionParams 群聊决策轮参数
type GroupDecisionParams struct {
	Members        []IOpcPerson
	ChatModel      model.ToolCallingChatModel
	ContextSummary string
	MemberRoster   string
}

// DecisionResult 群聊决策结果
type DecisionResult struct {
	Person   IOpcPerson
	Priority int
	Reason   string
}

// GroupReplyParams 群聊回复生成参数
type GroupReplyParams struct {
	Person          IOpcPerson
	GroupName       string
	GroupDesc       string
	AllMembers      []IOpcPerson
	ChatModel       model.ToolCallingChatModel
	HistoryMessages []data_models.Message
}

// DecideResponders 并行执行决策轮，判断哪些成员应该回复
// 返回按优先级降序排列的响应者列表
func DecideResponders(ctx context.Context, params GroupDecisionParams) ([]DecisionResult, error) {
	type rawDecision struct {
		ShouldRespond bool   `json:"should_respond"`
		Priority      int    `json:"priority"`
		Reason        string `json:"reason"`
	}

	var decisions []DecisionResult
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, m := range params.Members {
		wg.Add(1)
		go func(member IOpcPerson) {
			defer wg.Done()

			prompt := BuildDecisionPrompt(member, params.MemberRoster, params.ContextSummary)

			chatModel, err := cloneChatModel(ctx, params.ChatModel)
			if err != nil {
				chatModel = params.ChatModel
			}

			result, err := chatModel.Generate(ctx, []*schema.Message{
				{Role: schema.User, Content: prompt},
			})
			if err != nil {
				logger.Warm("OPC decision: %s generate failed: %v", member.Name(), err)
				return
			}

			// 解析决策 JSON
			content := strings.TrimSpace(result.Content)
			content = strings.TrimPrefix(content, "```json")
			content = strings.TrimPrefix(content, "```")
			content = strings.TrimSuffix(content, "```")
			content = strings.TrimSpace(content)

			var decision rawDecision
			if err := json.Unmarshal([]byte(content), &decision); err != nil {
				logger.Warm("OPC decision: %s parse failed: %v", member.Name(), err)
				return
			}

			if decision.ShouldRespond {
				mu.Lock()
				decisions = append(decisions, DecisionResult{
					Person:   member,
					Priority: decision.Priority,
					Reason:   decision.Reason,
				})
				mu.Unlock()
			}
		}(m)
	}
	wg.Wait()

	// 按优先级降序排列
	sort.Slice(decisions, func(i, j int) bool {
		return decisions[i].Priority > decisions[j].Priority
	})

	return decisions, nil
}

// GenerateGroupReply 为群聊中的一个成员生成回复
// 使用 GroupInfoTool + PersonInfoTool 替代提示词注入成员名单
func GenerateGroupReply(ctx context.Context, params GroupReplyParams) (*ChatResult, error) {
	systemPrompt := BuildGroupSystemPrompt(params.Person)
	msgs := BuildSchemaMessages(systemPrompt, params.HistoryMessages, "")

	// 创建会话级群聊工具
	groupTools := NewGroupChatTools(params.GroupName, params.GroupDesc, params.AllMembers)

	// 加载人员自身配置的工具
	personTools, err := resolvePersonTools(params.Person)
	if err != nil {
		personTools = nil // 加载失败不影响群聊工具
	}

	// 合并所有工具
	allTools := make([]tool.BaseTool, 0, len(groupTools)+len(personTools))
	allTools = append(allTools, groupTools...)
	allTools = append(allTools, personTools...)

	return runAgentWithTools(ctx, params.ChatModel, params.Person, systemPrompt, msgs, allTools)
}

// cloneChatModel 为群聊成员创建独立的 ChatModel 实例
// 目前直接复用传入的模型，后续可根据需要为每个成员创建独立实例
func cloneChatModel(_ context.Context, chatModel model.ToolCallingChatModel) (model.ToolCallingChatModel, error) {
	return chatModel, nil
}

// MemberLookupFunc 创建一个成员查找函数，用于 BuildContextSummary
func MemberLookupFunc(members []IOpcPerson) func(uuid string) string {
	return func(uuid string) string {
		for _, m := range members {
			if m.Uuid() == uuid {
				return m.Name()
			}
		}
		return ""
	}
}
