// ChatAgent - 主 Agent，负责接收用户输入并根据需要委托给子 Agent
// 基于 Agent-to-Agent 架构：Chat 本身是 Agent，通过 AgentTool 调用 DateTimeAgent、FruitPriceAgent

package main

import (
	"context"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

const chatAgentInstruction = `你是一个友好的 AI 助手。你可以直接回答一般性问题，也可以调用子 Agent 获取专业信息：

- 日期时间（今天几号、星期几、现在几点）→ 必须调用 DateTimeAgent 获取准确信息
- 水果价格（苹果多少钱、草莓价格等）→ 必须调用 FruitPriceAgent 获取准确价格
- 其他聊天、问候、通用问题 → 直接回答

重要规则：
1. 当需要日期/时间/价格信息时，必须实际调用相应子 Agent，并将工具返回的结果完整呈现给用户
2. 不要只说「我将帮你查看」而不调用工具，必须调用并展示实际结果
3. 遇到条件判断（如「如果今天6号则查苹果，如果星期五则查香蕉」），先调 DateTimeAgent 获知日期，再根据条件调 FruitPriceAgent，最后综合回答
4. 用简洁清晰的中文回复。`

// newChatAgent 创建主 Chat Agent，将 DateTime 和 Fruit 子 Agent 作为 AgentTool
func newChatAgent(ctx context.Context, chatModel model.ToolCallingChatModel, subAgents []adk.Agent) (adk.Agent, error) {
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
