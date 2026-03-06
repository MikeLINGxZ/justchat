// ChatAgent - 主 Agent，负责接收用户输入并根据需要委托给子 Agent
// 基于 Agent-to-Agent 架构：Chat 本身是 Agent，通过 AgentTool 调用 DateTimeAgent、FruitPriceAgent

package agents

import (
	"context"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

const chatAgentInstruction = `你是一个友好的 AI 助手。你可以直接回答一般性问题，也可以将专业问题转交给子 Agent 处理：

- 日期时间问题（如：今天星期几、现在几点）→ 转交给 DateTimeAgent
- 水果价格问题（如：苹果多少钱、草莓价格）→ 转交给 FruitPriceAgent
- 其他聊天、问候、通用问题 → 直接回答

根据用户问题智能决定是亲自回答还是调用相应的子 Agent。用简洁清晰的中文回复。`

// NewChatAgent 创建主 Chat Agent，将 DateTime 和 Fruit 子 Agent 作为 AgentTool
func NewChatAgent(ctx context.Context, chatModel model.ToolCallingChatModel, subAgents []adk.Agent) (adk.Agent, error) {
	tools := make([]tool.BaseTool, 0, len(subAgents))
	for _, agent := range subAgents {
		tools = append(tools, adk.NewAgentTool(ctx, agent))
	}

	return adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "ChatAgent",
		Description: "主对话助手，可协调日期时间、水果价格等专业子 Agent",
		Instruction: chatAgentInstruction,
		Model:       chatModel,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: tools,
			},
		},
	})
}
