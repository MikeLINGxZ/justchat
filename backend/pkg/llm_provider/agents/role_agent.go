package agents

import (
	"context"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
)

func NewRoleAgent(ctx context.Context, chatModel model.ToolCallingChatModel, name, description, instruction string, tools []tool.BaseTool, toolMiddleware compose.ToolMiddleware) (adk.Agent, error) {
	return adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        name,
		Description: description,
		Instruction: instruction,
		Model:       chatModel,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: tools,
			},
		},
		Middlewares: []adk.AgentMiddleware{
			{
				WrapToolCall: toolMiddleware,
			},
		},
	})
}
