package llm_provider

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/schema"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/wrapper_models"
)

type Aliyun struct {
	providerModel wrapper_models.ProviderModel
}

func NewAliyun(providerModel wrapper_models.ProviderModel) IProvider {
	return &Aliyun{
		providerModel: providerModel,
	}
}

func (a *Aliyun) Completions(ctx context.Context, messages []schema.Message) (*schema.StreamReader[*schema.Message], error) {
	chatModel, err := qwen.NewChatModel(ctx, &qwen.ChatModelConfig{
		BaseURL: a.providerModel.BaseUrl,
		Model:   a.providerModel.Model,
		APIKey:  a.providerModel.ApiKey,
	})

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
