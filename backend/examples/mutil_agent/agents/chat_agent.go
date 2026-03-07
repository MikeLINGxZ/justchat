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
- 日期时间（今天几号、星期几、现在几点）→ DateTimeAgent
- 水果价格（苹果多少钱、草莓价格等）→ FruitPriceAgent
- MemoryAgent：以下情况必须调用 MemoryAgent，不得跳过：
  · 记录：用户提及发生的事、消费、安排等想留存时，必须调用 record_daily。例如：「昨天我买了一台电脑花了5200」「今天去了公园」「记一下午饭吃了面」
  · 回忆：用户问过去或今天的安排/消费/做过的事时，必须先调用 read_daily 查询。例如：「今天要做什么」「昨天我花了多少钱」「昨天做了什么」「上周买过什么」
- 其他聊天、问候、通用问题 → 直接回答
重要：涉及「记」「花了」「买了」「做了」等表达，或询问过去/今日安排、消费、事件时，务必调用 MemoryAgent，不得凭空回答。用简洁清晰的中文回复。`

// NewChatAgent 创建主 Chat Agent，将 DateTime 和 Fruit 子 Agent 作为 AgentTool
func NewChatAgent(ctx context.Context, chatModel model.ToolCallingChatModel, subAgents []adk.Agent) (adk.Agent, error) {
	tools := make([]tool.BaseTool, 0, len(subAgents))
	for _, agent := range subAgents {
		tools = append(tools, adk.NewAgentTool(ctx, agent))
	}

	return adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "ChatAgent",
		Description: "主对话助手，可协调日期时间、水果价格、记忆日常等专业子 Agent",
		Instruction: chatAgentInstruction,
		Model:       chatModel,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: tools,
			},
		},
	})
}
