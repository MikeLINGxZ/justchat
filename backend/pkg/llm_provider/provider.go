package llm_provider

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino-ext/components/model/ollama"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/wrapper_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/llm_provider/agents"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/prompts"
)

type Provider struct {
	chatModel     model.BaseChatModel
	toolChatModel model.ToolCallingChatModel
	mainAgent     adk.Agent
	tools         []tool.BaseTool
	prompts       prompts.PromptSet
}

// NewToolCallingChatModel 根据供应商配置创建 ToolCallingChatModel 实例。
func NewToolCallingChatModel(ctx context.Context, providerModel wrapper_models.ProviderModel) (model.ToolCallingChatModel, error) {
	switch providerModel.ProviderType {
	case data_models.ProviderTypeDeepseek:
		return deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
			BaseURL: providerModel.BaseUrl,
			Model:   providerModel.Model,
			APIKey:  providerModel.ApiKey,
		})
	case data_models.ProviderTypeAliyuns:
		return qwen.NewChatModel(ctx, &qwen.ChatModelConfig{
			BaseURL: providerModel.BaseUrl,
			Model:   providerModel.Model,
			APIKey:  providerModel.ApiKey,
		})
	case data_models.ProviderTypeOpenrouter:
		return openai.NewChatModel(ctx, &openai.ChatModelConfig{
			BaseURL: providerModel.BaseUrl,
			Model:   providerModel.Model,
			APIKey:  providerModel.ApiKey,
		})
	case data_models.ProviderTypeOllama:
		return ollama.NewChatModel(ctx, &ollama.ChatModelConfig{
			BaseURL: providerModel.BaseUrl,
			Model:   providerModel.Model,
		})
	default:
		return openai.NewChatModel(ctx, &openai.ChatModelConfig{
			BaseURL: providerModel.BaseUrl,
			Model:   providerModel.Model,
			APIKey:  providerModel.ApiKey,
		})
	}
}

// NewLlmProvider 创建 LLM 供应商，tools 为可选参数，传入时会将工具绑定到模型以支持 tool calling
func NewLlmProvider(ctx context.Context, providerModel wrapper_models.ProviderModel, subAgents []adk.Agent, tools []tool.BaseTool, toolMiddleware compose.ToolMiddleware, promptSet prompts.PromptSet, skillSummary string) (*Provider, error) {
	chatModel, err := NewToolCallingChatModel(ctx, providerModel)
	if err != nil {
		return nil, err
	}

	instruction := promptSet.MainAgentSystem
	if skillSummary != "" {
		instruction = instruction + "\n\n" + skillSummary
	}

	mainAgent, err := agents.NewMainAgent(ctx, chatModel, subAgents, tools, toolMiddleware, instruction)
	if err != nil {
		return nil, err
	}

	return &Provider{chatModel: chatModel, toolChatModel: chatModel, tools: tools, mainAgent: mainAgent, prompts: promptSet}, nil
}

func (p *Provider) Completions(ctx context.Context, messages []schema.Message) (*schema.StreamReader[*schema.Message], error) {
	messagesPoint := make([]*schema.Message, len(messages))
	for i := range messages {
		messagesPoint[i] = &messages[i]
	}

	streamResult, err := p.chatModel.Stream(ctx, messagesPoint)
	if err != nil {
		return nil, err
	}

	return streamResult, nil
}

func (p *Provider) AgentCompletions(ctx context.Context, messages []schema.Message) (*adk.AsyncIterator[*adk.AgentEvent], error) {
	messagesPoint := make([]*schema.Message, len(messages))
	for i := range messages {
		messagesPoint[i] = &messages[i]
	}

	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           p.mainAgent,
		EnableStreaming: true,
	})

	iter := runner.Run(ctx, messagesPoint)
	return iter, nil
}

func (p *Provider) Generate(ctx context.Context, messages []schema.Message) (*schema.Message, error) {
	messagesPoint := make([]*schema.Message, len(messages))
	for i := range messages {
		messagesPoint[i] = &messages[i]
	}
	return p.chatModel.Generate(ctx, messagesPoint)
}

func (p *Provider) Stream(ctx context.Context, messages []schema.Message) (*schema.StreamReader[*schema.Message], error) {
	return p.Completions(ctx, messages)
}

func (p *Provider) ToolCallingModel() model.ToolCallingChatModel {
	return p.toolChatModel
}

func (p *Provider) Prompts() prompts.PromptSet {
	return p.prompts
}

// GenChatTitle 生成一个聊天的标题
func (p *Provider) GenChatTitle(ctx context.Context, messages []schema.Message) (string, error) {
	genTitle := ""
	contextMessages := []schema.Message{
		{
			Role:    schema.System,
			Content: p.prompts.TitleSystem,
		},
	}
	contextMessages = append(contextMessages, messages...)

	var messagePoint []*schema.Message
	for _, item := range contextMessages {
		messagePoint = append(messagePoint, &item)
	}

	generateMsg, err := p.chatModel.Generate(ctx, messagePoint)
	if err != nil {
		return "", err
	}
	genTitle += generateMsg.Content

	return genTitle, nil
}
