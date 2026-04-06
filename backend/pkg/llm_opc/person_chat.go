package llm_opc

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/llm_provider/tools"
)

// PersonChatParams 个人聊天生成参数
type PersonChatParams struct {
	Person          IOpcPerson
	ChatModel       model.ToolCallingChatModel
	HistoryMessages []data_models.Message
	SkipMessageUuid string // 跳过正在填充的 assistant 消息
}

// ChatResult 聊天生成结果
type ChatResult struct {
	Content          string
	ReasoningContent string
}

// GeneratePersonReply 为指定人员生成回复
func GeneratePersonReply(ctx context.Context, params PersonChatParams) (*ChatResult, error) {
	systemPrompt := BuildPersonSystemPrompt(params.Person)
	msgs := BuildSchemaMessages(systemPrompt, params.HistoryMessages, params.SkipMessageUuid)

	// 加载用户配置的工具
	personTools, err := resolvePersonTools(params.Person)
	if err != nil {
		return nil, err
	}

	if len(personTools) > 0 {
		return runAgentWithTools(ctx, params.ChatModel, params.Person, systemPrompt, msgs, personTools)
	}

	// 无工具时直接调用模型
	result, err := params.ChatModel.Generate(ctx, msgs)
	if err != nil {
		return nil, err
	}

	return &ChatResult{
		Content:          result.Content,
		ReasoningContent: result.ReasoningContent,
	}, nil
}

// resolvePersonTools 从 ToolRouter 加载人员配置的工具
func resolvePersonTools(person IOpcPerson) ([]tool.BaseTool, error) {
	toolIDs := person.Tools()
	if len(toolIDs) == 0 {
		return nil, nil
	}
	return tools.ToolRouter.GetToolsByIds(toolIDs)
}

// runAgentWithTools 使用 adk agent 运行带工具的对话
func runAgentWithTools(ctx context.Context, chatModel model.ToolCallingChatModel, person IOpcPerson, instruction string, messages []*schema.Message, agentTools []tool.BaseTool) (*ChatResult, error) {
	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        fmt.Sprintf("opc_person_%s", person.Uuid()),
		Description: person.Desc(),
		Instruction: instruction,
		Model:       chatModel,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: agentTools,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           agent,
		EnableStreaming: false,
	})

	iter := runner.Run(ctx, messages)
	return collectAgentResult(iter)
}

// collectAgentResult 从 agent 迭代器中收集最终结果
func collectAgentResult(iter *adk.AsyncIterator[*adk.AgentEvent]) (*ChatResult, error) {
	result := &ChatResult{}

	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			return nil, event.Err
		}
		if event.Output == nil || event.Output.MessageOutput == nil {
			continue
		}
		mo := event.Output.MessageOutput
		if mo.Role != schema.Assistant {
			continue
		}
		if mo.Message != nil && mo.Message.Content != "" {
			result.Content = mo.Message.Content
			result.ReasoningContent = mo.Message.ReasoningContent
		}
	}

	return result, nil
}
