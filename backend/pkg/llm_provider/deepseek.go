package llm_provider

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino/schema"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/wrapper_models"
)

type Deepseek struct {
	providerModel wrapper_models.ProviderModel
}

func NewDeepseek(providerModel wrapper_models.ProviderModel) IProvider {
	return &Deepseek{
		providerModel: providerModel,
	}
}
func (d *Deepseek) Completions(ctx context.Context, messages []schema.Message) (*schema.StreamReader[*schema.Message], error) {
	chatModel, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		BaseURL: d.providerModel.BaseUrl,
		Model:   d.providerModel.Model,
		APIKey:  d.providerModel.ApiKey,
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

func (d *Deepseek) BuildUserMessage(ctx context.Context, message view_models.MessagePkg) (*schema.Message, error) {
	return &schema.Message{
		Role:    schema.User,
		Content: message.Content,
	}, nil
}
