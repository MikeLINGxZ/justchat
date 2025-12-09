package llm

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
)

type ProviderType string

const (
	ProviderTypeDeepseek   ProviderType = "deepseek"
	ProviderTypeAliyuns    ProviderType = "aliyuns"
	ProviderTypeOpenrouter ProviderType = "openrouter"
	ProviderTypeOther      ProviderType = "other"
)

func (p ProviderType) String() string {
	return string(p)
}

type LlmProvider struct {
	providerType ProviderType
	baseURL      string
	apiKey       string
	model        string
}

func NewLlmProvider(providerType ProviderType, baseUrl, apiKey, model string) *LlmProvider {
	return &LlmProvider{
		providerType: providerType,
		baseURL:      baseUrl,
		apiKey:       apiKey,
		model:        model,
	}
}

func (l *LlmProvider) Completions(ctx context.Context, messages []schema.Message) (*schema.StreamReader[*schema.Message], error) {
	var chatModel model.BaseChatModel
	var err error

	// 创建llm模型实例
	switch l.providerType {
	case ProviderTypeDeepseek:
		chatModel, err = deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
			BaseURL: l.baseURL,
			Model:   l.model,
			APIKey:  l.apiKey,
		})
	case ProviderTypeAliyuns:
		chatModel, err = qwen.NewChatModel(ctx, &qwen.ChatModelConfig{
			BaseURL: l.baseURL,
			Model:   l.model,
			APIKey:  l.apiKey,
		})
	default:
		chatModel, err = openai.NewChatModel(ctx, &openai.ChatModelConfig{
			BaseURL: l.baseURL,
			Model:   l.model,
			APIKey:  l.apiKey,
		})
	}
	if err != nil {
		return nil, err
	}

	var messagesPoint []*schema.Message
	for _, item := range messages {
		messagesPoint = append(messagesPoint, &item)
	}
	// 调用LLM服务
	streamResult, err := chatModel.Stream(ctx, messagesPoint)
	if err != nil {
		return nil, err
	}

	return streamResult, nil
}
