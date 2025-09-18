package agents

import (
	"context"
	_ "embed"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/flow/agent/multiagent/host"
	"github.com/cloudwego/eino/schema"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/examples/memory_demo_01/internal/storage"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/examples/memory_demo_01/internal/tools"
)

//go:embed memory.prompt.md
var memoryAgentPrompt string

func NewMemory(ctx context.Context, baseURL, apiKey, model string, storage *storage.Storage) (*host.Specialist, error) {
	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		BaseURL: baseURL,
		Model:   model,
		APIKey:  apiKey,
	})
	if err != nil {
		return nil, err
	}

	baseTools, err := newMemoryTools(storage)
	if err != nil {
		return nil, err
	}

	// 创建 tools 节点
	toolsNode, err := compose.NewToolNode(context.Background(), &compose.ToolsNodeConfig{
		Tools: baseTools,
	})

	chain := compose.NewChain[[]*schema.Message, []*schema.Message]()
	// 添加一个系统提示词
	chain.
		AppendLambda(compose.InvokableLambda(func(ctx context.Context, input []*schema.Message) ([]*schema.Message, error) {
			systemMsg := &schema.Message{
				Role:    schema.System,
				Content: memoryAgentPrompt,
			}
			return append([]*schema.Message{systemMsg}, input...), nil
		})).
		AppendChatModel(chatModel).
		AppendToolsNode(toolsNode)

	r, err := chain.Compile(ctx)
	if err != nil {
		return nil, err
	}

	return &host.Specialist{
		AgentMeta: host.AgentMeta{
			Name:        "Memory_Weaver",
			IntendedUse: "Responsible for storing, organizing, and managing the user's personal memories, with context awareness and emotional sensitivity, simulating human-like long-term memory cognitive behaviors, and supporting personalized interaction and experiential accumulation.",
		},
		Invokable: func(ctx context.Context, input []*schema.Message, opts ...agent.AgentOption) (output *schema.Message, err error) {
			invoke, err := r.Invoke(ctx, input, agent.GetComposeOptions(opts...)...)
			if err != nil {
				return nil, err
			}
			return invoke[0], nil
		},
	}, nil
}

func newMemoryTools(storage *storage.Storage) ([]tool.BaseTool, error) {
	writerMemory, err := utils.InferTool(
		"write_memory",
		"Record a personal memory or life event with title, content, and optional date",
		tools.NewWriteMemoryTool(storage))
	if err != nil {
		return nil, err
	}

	readMemory, err := utils.InferTool(
		"read_memory",
		"Search and retrieve personal memories based on a keyword and/or optional date range. "+
			"Use this when the user asks about past events, experiences, or wants to recall something they've shared before.",
		tools.NewReadMemoryTool(storage),
	)
	if err != nil {
		return nil, err
	}

	getCurrentTime, err := utils.InferTool(
		"get_current_time",
		"Get the current date and time in UTC. "+
			"Use this whenever the user asks about the current time, today's date, or needs time context for reasoning (e.g., 'What's the date today?' or 'When did this happen?').",
		tools.NewGetCurrentTimeTool(),
	)
	if err != nil {
		return nil, err
	}

	return []tool.BaseTool{
		writerMemory,
		readMemory,
		getCurrentTime,
	}, nil
}
