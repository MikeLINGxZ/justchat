package agents

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/flow/agent/multiagent/host"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/agents/memory/storage"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/agents/memory/tools"
)

//go:embed agent.prompt.v0.md
var memoryAgentPrompt string

func NewMemoryAgent(ctx context.Context, baseURL, apiKey, model string, storage *storage.Storage) (*host.Specialist, error) {

	arkModel, err := ark.NewChatModel(ctx, &ark.ChatModelConfig{
		BaseURL: baseURL,
		APIKey:  apiKey,
		Model:   model,
	})
	if err != nil {
		slog.Error(fmt.Sprintf("new memory agent: %v", err))
		return nil, err
	}

	writeMemoryTool, err := tools.NewWriteMemoryTool(storage)
	if err != nil {
		return nil, err
	}
	readMemoryTool, err := tools.NewReadMemoryTool(storage)
	if err != nil {
		return nil, err
	}
	editMemoryTool, err := tools.NewEditMemoryTool(storage)
	if err != nil {
		return nil, err
	}
	getCurrentTimeTool, err := tools.NewGetCurrentTimeTool()
	if err != nil {
		return nil, err
	}

	ragent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: arkModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{writeMemoryTool, readMemoryTool, editMemoryTool, getCurrentTimeTool},
		},
	})
	if err != nil {
		slog.Error(fmt.Sprintf("new memory agent: %v", err))
		return nil, err
	}

	return &host.Specialist{
		AgentMeta: host.AgentMeta{
			Name:        "Memory_Weaver",
			IntendedUse: "Responsible for storing, organizing, and managing the user's personal memories with context awareness and emotional sensitivity. Simulates human-like long-term memory behaviors including encoding, association, and retrieval. Must actively recall past experiences when queried, support personalized interaction, experiential accumulation, and answer questions about the user's history accurately and naturally.",
		},
		Streamable: func(ctx context.Context, input []*schema.Message, opts ...agent.AgentOption) (output *schema.StreamReader[*schema.Message], err error) {
			var agentMessages []*schema.Message
			agentMessages = append(agentMessages, &schema.Message{
				Role:    schema.System,
				Content: memoryAgentPrompt,
			})
			agentMessages = append(agentMessages, input...)
			return ragent.Stream(ctx, agentMessages, opts...)
		},
	}, nil
}

func NewMemory(ctx context.Context, baseURL, apiKey, model string, storage *storage.Storage) (*host.Specialist, error) {

	// 1. 创建对话模板
	chatTpl := prompt.FromMessages(schema.FString,
		schema.SystemMessage(memoryAgentPrompt),
		schema.MessagesPlaceholder("message_histories", true),
		schema.UserMessage("{user_query}"),
	)

	// 2. 创建对话模型
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: baseURL,
		Model:   model,
		APIKey:  apiKey,
	})
	if err != nil {
		slog.Error(fmt.Sprintf("new chat model err: %v", err))
		return nil, err
	}

	// 3. 绑定工具
	writeMemoryTool, err := tools.NewWriteMemoryTool(storage)
	if err != nil {
		return nil, err
	}
	readMemoryTool, err := tools.NewReadMemoryTool(storage)
	if err != nil {
		return nil, err
	}
	getCurrentTimeTool, err := tools.NewGetCurrentTimeTool()
	if err != nil {
		return nil, err
	}
	toolInfos, err := newMemoryTools(ctx, []tool.InvokableTool{
		writeMemoryTool,
		readMemoryTool,
		getCurrentTimeTool,
	})
	if err != nil {
		slog.Error(fmt.Sprintf("new memory tools error: %v", err))
		return nil, err
	}
	err = chatModel.BindForcedTools(toolInfos)
	if err != nil {
		slog.Error(fmt.Sprintf("bind forced tools error: %v", err))
		return nil, err
	}

	// 4. 创建node实例
	toolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
		Tools: []tool.BaseTool{writeMemoryTool, readMemoryTool, getCurrentTimeTool},
	})
	if err != nil {
		slog.Error(fmt.Sprintf("new tool err: %v", err))
		return nil, err
	}

	const (
		nodeKeyOfTemplate  = "template"
		nodeKeyOfChatModel = "chat_model"
		nodeKeyOfTools     = "tools"
	)

	// 5. 创建一个graph
	g := compose.NewGraph[map[string]any, []*schema.Message]()

	_ = g.AddChatTemplateNode(nodeKeyOfTemplate, chatTpl)
	_ = g.AddChatModelNode(nodeKeyOfChatModel, chatModel)
	_ = g.AddToolsNode(nodeKeyOfTools, toolsNode)
	_ = g.AddEdge(compose.START, nodeKeyOfTemplate)
	_ = g.AddEdge(nodeKeyOfTemplate, nodeKeyOfChatModel)
	_ = g.AddEdge(nodeKeyOfChatModel, nodeKeyOfTools)
	_ = g.AddEdge(nodeKeyOfTools, compose.END)

	r, err := g.Compile(ctx)
	if err != nil {
		slog.Error(fmt.Sprintf("compile err: %v", err))
		return nil, err
	}

	return &host.Specialist{
		AgentMeta: host.AgentMeta{
			Name:        "Memory_Weaver",
			IntendedUse: "Responsible for storing, organizing, and managing the user's personal memories, with context awareness and emotional sensitivity, simulating human-like long-term memory cognitive behaviors, and supporting personalized interaction and experiential accumulation.",
		},
		Invokable: func(ctx context.Context, input []*schema.Message, opts ...agent.AgentOption) (output *schema.Message, err error) {
			var historyMessages []*schema.Message
			if len(input) > 1 {
				historyMessages = input[1:]
			}
			out, err := r.Invoke(ctx, map[string]any{
				"message_histories": historyMessages,
				"user_query":        input[len(input)-1],
			})
			if err != nil {
				return
			}
			return out[len(out)-1], nil
		},
		Streamable: nil,
	}, nil
}

func newMemoryTools(ctx context.Context, invokableTools []tool.InvokableTool) ([]*schema.ToolInfo, error) {
	var toolInfos []*schema.ToolInfo

	for _, invokableTool := range invokableTools {
		toolInfo, err := invokableTool.Info(ctx)
		if err != nil {
			return nil, err
		}
		toolInfos = append(toolInfos, toolInfo)
	}

	return toolInfos, nil
}
