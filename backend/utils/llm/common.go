package llm

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
)

type LlmProvider struct {
	providerType      data_models.ProviderType
	baseURL           string
	apiKey            string
	model             string
	fileUploadBaseUrl *string
}

func NewLlmProvider(providerType data_models.ProviderType, fileUploadBaseUrl *string, baseUrl, apiKey, model string) *LlmProvider {
	return &LlmProvider{
		providerType:      providerType,
		baseURL:           baseUrl,
		apiKey:            apiKey,
		model:             model,
		fileUploadBaseUrl: fileUploadBaseUrl,
	}
}

func (l *LlmProvider) Completions(ctx context.Context, messages []schema.Message) (*schema.StreamReader[*schema.Message], error) {
	var chatModel model.BaseChatModel
	var err error

	// 创建llm模型实例
	switch l.providerType {
	case data_models.ProviderTypeDeepseek:
		chatModel, err = deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
			BaseURL: l.baseURL,
			Model:   l.model,
			APIKey:  l.apiKey,
		})
	case data_models.ProviderTypeAliyuns:
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
