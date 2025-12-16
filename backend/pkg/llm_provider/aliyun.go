package llm_provider

import (
	"context"

	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/schema"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/wrapper_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils"
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

func (a *Aliyun) BuildUserMessage(ctx context.Context, message view_models.MessagePkg) (*schema.Message, error) {
	var paths []string
	path2base64data := make(map[string]string)
	for _, file := range message.Files {
		paths = append(paths, file.FilePath)
	}

	for _, path := range paths {
		data, err := utils.ReadFile2Base64Data(path)
		if err != nil {
			return nil, err
		}
		path2base64data[path] = data
	}

	var userInputMultiContent []schema.MessageInputPart
	if message.Content != "" {
		userInputMultiContent = append(userInputMultiContent, schema.MessageInputPart{
			Type: schema.ChatMessagePartTypeText,
			Text: message.Content,
		})
	}

	for _, item := range message.Files {
		var text string
		var img *schema.MessageInputImage
		var audio *schema.MessageInputAudio
		var video *schema.MessageInputVideo
		var file *schema.MessageInputFile
		base64Data := path2base64data[item.FilePath]
		messagePartCommon := schema.MessagePartCommon{
			Base64Data: &base64Data,
			MIMEType:   item.MineType,
			Extra: map[string]interface{}{
				"name":                   item.Name,
				"path":                   item.FilePath,
				"mime_type":              item.MineType,
				"chat_message_part_type": item.ChatMessagePartType,
				"size":                   item.Size,
			},
		}
		switch item.ChatMessagePartType {
		case schema.ChatMessagePartTypeText, schema.ChatMessagePartTypeFileURL:
			continue
		case schema.ChatMessagePartTypeImageURL:
			img = &schema.MessageInputImage{
				MessagePartCommon: messagePartCommon,
				Detail:            schema.ImageURLDetailHigh,
			}
		case schema.ChatMessagePartTypeAudioURL:
			audio = &schema.MessageInputAudio{
				MessagePartCommon: messagePartCommon,
			}
		case schema.ChatMessagePartTypeVideoURL:
			video = &schema.MessageInputVideo{
				MessagePartCommon: messagePartCommon,
			}
		}
		if img == nil && audio == nil && video == nil {
			continue
		}
		userInputMultiContent = append(userInputMultiContent, schema.MessageInputPart{
			Type:  item.ChatMessagePartType,
			Text:  text,
			Image: img,
			Audio: audio,
			Video: video,
			File:  file,
		})
	}

	return &schema.Message{
		Role:                  schema.User,
		Content:               "",
		UserInputMultiContent: userInputMultiContent,
	}, nil
}
