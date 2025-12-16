package llm_provider

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/wrapper_models"
)

type LlmProvider struct {
	providerModel *wrapper_models.ProviderModel
}

func NewLlmProvider(providerModel *wrapper_models.ProviderModel) *LlmProvider {
	return &LlmProvider{
		providerModel: providerModel,
	}
}

func (l *LlmProvider) Completions(ctx context.Context, messages []schema.Message) (*schema.StreamReader[*schema.Message], error) {
	var chatModel model.BaseChatModel
	var err error

	// 创建llm模型实例
	switch l.providerModel.ProviderType {
	case data_models.ProviderTypeDeepseek:
		chatModel, err = deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
			BaseURL: l.providerModel.BaseUrl,
			Model:   l.providerModel.Model,
			APIKey:  l.providerModel.ApiKey,
		})
	case data_models.ProviderTypeAliyuns:
		chatModel, err = qwen.NewChatModel(ctx, &qwen.ChatModelConfig{
			BaseURL: l.providerModel.BaseUrl,
			Model:   l.providerModel.Model,
			APIKey:  l.providerModel.ApiKey,
		})
	default:
		chatModel, err = openai.NewChatModel(ctx, &openai.ChatModelConfig{
			BaseURL: l.providerModel.BaseUrl,
			Model:   l.providerModel.Model,
			APIKey:  l.providerModel.ApiKey,
		})
	}
	if err != nil {
		return nil, err
	}

	// 处理文件上传
	messages, err = l.processFile(messages)
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

func (l *LlmProvider) processFile(messages []schema.Message) ([]schema.Message, error) {
	return messages, nil
}
