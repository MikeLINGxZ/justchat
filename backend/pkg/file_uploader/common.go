package file_uploader

import (
	"github.com/cloudwego/eino/schema"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/wrapper_models"
)

type IFileUploader interface {
	Upload(paths []string) (map[string]string, error)
}

type ExtraKey string

const (
	ExtraKeyName                ExtraKey = "name"
	ExtraKeyFilePath            ExtraKey = "file_path"
	ExtraKeyMimeType            ExtraKey = "mime_type"
	ExtraKeySize                ExtraKey = "size"
	ExtraKeyChatMessagePartType ExtraKey = "chat_message_part_type"
	ExtraKeyPreviewImg          ExtraKey = "preview_img"
)

func (key ExtraKey) String() string {
	return string(key)
}

type Converter struct {
	providerModel *wrapper_models.ProviderModel
}

func NewConverter(providerModel *wrapper_models.ProviderModel) Converter {
	return Converter{
		providerModel: providerModel,
	}
}

func (d *Converter) ConvertMessageUserInputMultiContent(message *schema.Message) error {

	var uploader IFileUploader

	switch d.providerModel.ProviderType {
	case data_models.ProviderTypeAliyuns:
		uploader = &Aliyuns{
			providerModel: d.providerModel,
		}
	default:
		return nil
	}

	var paths []string
	for i, part := range message.UserInputMultiContent {
		switch part.Type {
		case schema.ChatMessagePartTypeImageURL:
			if message.UserInputMultiContent[i].Image == nil || message.UserInputMultiContent[i].Image.Extra == nil {
				continue
			}
			path := message.UserInputMultiContent[i].Image.Extra[ExtraKeyFilePath.String()]
			if path == nil || path == "" {
				continue
			}
			paths = append(paths, path.(string))
		case schema.ChatMessagePartTypeAudioURL:
			if message.UserInputMultiContent[i].Audio == nil || message.UserInputMultiContent[i].Audio.Extra == nil {
				continue
			}
			path := message.UserInputMultiContent[i].Audio.Extra[ExtraKeyFilePath.String()]
			if path == nil || path == "" {
				continue
			}
			paths = append(paths, path.(string))
		case schema.ChatMessagePartTypeVideoURL:
			if message.UserInputMultiContent[i].Video == nil || message.UserInputMultiContent[i].Video.Extra == nil {
				continue
			}
			path := message.UserInputMultiContent[i].Video.Extra[ExtraKeyFilePath.String()]
			if path == nil || path == "" {
				continue
			}
			paths = append(paths, path.(string))
		case schema.ChatMessagePartTypeFileURL:
			if message.UserInputMultiContent[i].File == nil || message.UserInputMultiContent[i].File.Extra == nil {
				continue
			}
			path := message.UserInputMultiContent[i].File.Extra[ExtraKeyFilePath.String()]
			if path == nil || path == "" {
				continue
			}
			paths = append(paths, path.(string))
		default:
			continue
		}
	}

	path2url, err := uploader.Upload(paths)
	if err != nil {
		return err
	}

	for i, part := range message.UserInputMultiContent {
		switch part.Type {
		case schema.ChatMessagePartTypeImageURL:
			if message.UserInputMultiContent[i].Image == nil || message.UserInputMultiContent[i].Image.Extra == nil {
				continue
			}
			path := message.UserInputMultiContent[i].Image.Extra[ExtraKeyFilePath.String()]
			if path == nil || path == "" {
				continue
			}
			url := path2url[path.(string)]
			if url != "" {
				message.UserInputMultiContent[i].Image.URL = &url
			}
		case schema.ChatMessagePartTypeAudioURL:
			if message.UserInputMultiContent[i].Audio == nil || message.UserInputMultiContent[i].Audio.Extra == nil {
				continue
			}
			path := message.UserInputMultiContent[i].Audio.Extra[ExtraKeyFilePath.String()]
			if path == nil || path == "" {
				continue
			}
			url := path2url[path.(string)]
			if url != "" {
				message.UserInputMultiContent[i].Audio.URL = &url
			}
		case schema.ChatMessagePartTypeVideoURL:
			if message.UserInputMultiContent[i].Video == nil || message.UserInputMultiContent[i].Video.Extra == nil {
				continue
			}
			path := message.UserInputMultiContent[i].Video.Extra[ExtraKeyFilePath.String()]
			if path == nil || path == "" {
				continue
			}
			url := path2url[path.(string)]
			if url != "" {
				message.UserInputMultiContent[i].Video.URL = &url
			}
		case schema.ChatMessagePartTypeFileURL:
			if message.UserInputMultiContent[i].File == nil || message.UserInputMultiContent[i].File.Extra == nil {
				continue
			}
			path := message.UserInputMultiContent[i].File.Extra[ExtraKeyFilePath.String()]
			if path == nil || path == "" {
				continue
			}
			url := path2url[path.(string)]
			if url != "" {
				message.UserInputMultiContent[i].File.URL = &url
			}
		default:
			continue
		}
	}

	return nil
}
