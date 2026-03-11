package agents

import (
	"context"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

// NewMainAgent 创建主 Chat Agent，将 DateTime 和 Fruit 子 Agent 作为 AgentTool
func NewMainAgent(ctx context.Context, chatModel model.ToolCallingChatModel, subAgents []adk.Agent, tools []tool.BaseTool) (adk.Agent, error) {
	for _, agent := range subAgents {
		tools = append(tools, adk.NewAgentTool(ctx, agent))
	}

	return adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "MainChatAgent",
		Description: "主对话助手，可协调日期时间、水果价格等专业子 Agent",
		Instruction: "你是一个友好的 AI 助手。",
		Model:       chatModel,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: tools,
			},
		},
	})
}
