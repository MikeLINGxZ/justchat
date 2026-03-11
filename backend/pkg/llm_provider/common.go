package llm_provider

import (
	"context"
	"fmt"
	"io"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/wrapper_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/llm_provider/agents"
)

type Provider struct {
	chatModel model.BaseChatModel
	mainAgent adk.Agent
	tools     []tool.BaseTool
}

// NewLlmProvider 创建 LLM 供应商，tools 为可选参数，传入时会将工具绑定到模型以支持 tool calling
func NewLlmProvider(ctx context.Context, providerModel wrapper_models.ProviderModel, tools []tool.BaseTool) (*Provider, error) {
	var chatModel model.ToolCallingChatModel
	var err error
	switch providerModel.ProviderType {
	case data_models.ProviderTypeDeepseek:
		chatModel, err = deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
			BaseURL: providerModel.BaseUrl,
			Model:   providerModel.Model,
			APIKey:  providerModel.ApiKey,
		})
		if err != nil {
			return nil, err
		}
	case data_models.ProviderTypeAliyuns:
		chatModel, err = qwen.NewChatModel(ctx, &qwen.ChatModelConfig{
			BaseURL: providerModel.BaseUrl,
			Model:   providerModel.Model,
			APIKey:  providerModel.ApiKey,
		})
		if err != nil {
			return nil, err
		}
	case data_models.ProviderTypeOpenrouter:
		chatModel, err = openai.NewChatModel(ctx, &openai.ChatModelConfig{
			BaseURL: providerModel.BaseUrl,
			Model:   providerModel.Model,
			APIKey:  providerModel.ApiKey,
		})
		if err != nil {
			return nil, err
		}
	default:
		chatModel, err = openai.NewChatModel(ctx, &openai.ChatModelConfig{
			BaseURL: providerModel.BaseUrl,
			Model:   providerModel.Model,
			APIKey:  providerModel.ApiKey,
		})
		if err != nil {
			return nil, err
		}
	}

	mainAgent, err := agents.NewMainAgent(ctx, chatModel, nil, tools)
	if err != nil {
		return nil, err
	}

	return &Provider{chatModel: chatModel, tools: tools, mainAgent: mainAgent}, nil
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

//
//func (p *Provider) BuildUserMessage(ctx context.Context, message view_models.Message) (*schema.Message, error) {
//	var paths []string
//	path2base64data := make(map[string]string)
//	for _, file := range message.Files {
//		paths = append(paths, file.Path)
//	}
//
//	for _, path := range paths {
//		data, err := utils.ReadFile2Base64Data(path)
//		if err != nil {
//			return nil, err
//		}
//		path2base64data[path] = data
//	}
//
//	var userInputMultiContent []schema.MessageInputPart
//	if message.Content != "" {
//		userInputMultiContent = append(userInputMultiContent, schema.MessageInputPart{
//			Type: schema.ChatMessagePartTypeText,
//			Text: message.Content,
//		})
//	}
//
//	for _, item := range message.Files {
//		var text string
//		var img *schema.MessageInputImage
//		var audio *schema.MessageInputAudio
//		var video *schema.MessageInputVideo
//		var file *schema.MessageInputFile
//		base64Data := path2base64data[item.Path]
//		messagePartCommon := schema.MessagePartCommon{
//			Base64Data: &base64Data,
//			MIMEType:   item.MineType,
//			Extra: map[string]interface{}{
//				"name":                   item.Name,
//				"path":                   item.Path,
//				"mime_type":              item.MineType,
//				"chat_message_part_type": item.ChatMessagePartType,
//				"size":                   item.Size,
//			},
//		}
//		switch item.ChatMessagePartType {
//		case schema.ChatMessagePartTypeText, schema.ChatMessagePartTypeFileURL:
//			continue
//		case schema.ChatMessagePartTypeImageURL:
//			img = &schema.MessageInputImage{
//				MessagePartCommon: messagePartCommon,
//				Detail:            schema.ImageURLDetailHigh,
//			}
//		case schema.ChatMessagePartTypeAudioURL:
//			audio = &schema.MessageInputAudio{
//				MessagePartCommon: messagePartCommon,
//			}
//		case schema.ChatMessagePartTypeVideoURL:
//			video = &schema.MessageInputVideo{
//				MessagePartCommon: messagePartCommon,
//			}
//		}
//		if img == nil && audio == nil && video == nil {
//			continue
//		}
//		userInputMultiContent = append(userInputMultiContent, schema.MessageInputPart{
//			Type:  item.ChatMessagePartType,
//			Text:  text,
//			Image: img,
//			Audio: audio,
//			Video: video,
//			File:  file,
//		})
//	}
//
//	return &schema.Message{
//		Role:                  schema.User,
//		Content:               "",
//		UserInputMultiContent: userInputMultiContent,
//	}, nil
//}

// GenChatTitle 生成一个聊天的标题
func (p *Provider) GenChatTitle(ctx context.Context, messages []schema.Message) (string, error) {
	genTitle := ""
	contextMessages := []schema.Message{
		{
			Role: schema.System,
			Content: `
					你是一位专业的对话摘要与标题提炼专家。请根据我提供的聊天记录，生成1个最合适的标题，要求满足以下所有条件：
					✅ 准确概括核心主题：抓住双方讨论的实质焦点（如问题、决策、情感、事件或共识），而非罗列细节；
					✅ 简洁有力：控制在8–15个汉字以内，避免标点（除必要顿号）、英文和冗余修饰；
					✅ 中性客观，不带主观判断或情绪渲染（除非聊天本身是明确的情感倾诉，此时可适度体现温度，如“深夜倾诉：关于成长的迷茫与自我接纳”）；
					✅ 适配通用场景：标题应便于归档、检索或快速理解，不依赖上下文即可读懂；
					✅ 直接输出标题，不需要其他内容；
					❌ 不要解释、不要复述对话、不要添加额外信息、不要输出任何说明文字——只输出标题本身，且仅一行。
					
					请严格遵循以上规则。现在，我的聊天记录如下：
					`,
		},
	}
	contextMessages = append(contextMessages, messages...)

	resp, err := p.Completions(context.Background(), contextMessages)
	if err == nil {
		for {
			recv, err := resp.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				return "", err
			}
			genTitle += recv.Content
		}
	}
	if genTitle == "" {
		return "", fmt.Errorf("failed to generate title")
	}
	return genTitle, nil
}
