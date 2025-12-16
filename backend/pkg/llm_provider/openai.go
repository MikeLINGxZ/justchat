package llm_provider

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/schema"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/wrapper_models"
)

type Openai struct {
	providerModel wrapper_models.ProviderModel
}

func NewOpenai(providerModel wrapper_models.ProviderModel) IProvider {
	return &Openai{providerModel: providerModel}
}

func (o *Openai) Completions(ctx context.Context, messages []schema.Message) (*schema.StreamReader[*schema.Message], error) {
	chatModel, err := qwen.NewChatModel(ctx, &qwen.ChatModelConfig{
		BaseURL: o.providerModel.BaseUrl,
		Model:   o.providerModel.Model,
		APIKey:  o.providerModel.ApiKey,
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

func (o *Openai) UploadFile(paths []string) (map[string]string, error) {
	//TODO implement me
	panic("implement me")
}

func (o *Openai) BuildUserMessage(ctx context.Context, message view_models.MessagePkg) (*schema.Message, error) {
	//TODO implement me
	panic("implement me")
}
