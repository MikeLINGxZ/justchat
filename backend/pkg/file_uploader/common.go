package file_uploader

import (
	"github.com/cloudwego/eino/schema"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
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

func (d *Converter) ConvertMessageUserInputMultiContent(message *view_models.MessagePkg) error {
	if message == nil {
		return nil
	}

	uploader := d.getUploader()
	if uploader == nil {
		return nil
	}

	var paths []string
	path2files := make(map[string]view_models.File)
	for _, file := range message.Files {
		paths = append(paths, file.FilePath)
		path2files[file.FilePath] = file
	}

	path2url, err := uploader.Upload(paths)
	if err != nil {
		return err
	}

	for _, file := range message.Files {
		fileItem := path2files[file.FilePath]
		fileUrl := path2url[file.FilePath]
		if message.Message.UserInputMultiContent == nil {
			message.Message.UserInputMultiContent = []schema.MessageInputPart{}
		}
		switch file.ChatMessagePartType {
		case schema.ChatMessagePartTypeText:
			message.Message.UserInputMultiContent = append(message.Message.UserInputMultiContent, schema.MessageInputPart{
				Type: schema.ChatMessagePartTypeText,
				Text: "",
			})
		case schema.ChatMessagePartTypeImageURL:
			message.Message.UserInputMultiContent = append(message.Message.UserInputMultiContent, schema.MessageInputPart{
				Type: schema.ChatMessagePartTypeImageURL,
				Image: &schema.MessageInputImage{
					MessagePartCommon: schema.MessagePartCommon{
						URL:        &fileUrl,
						Base64Data: nil,
						MIMEType:   fileItem.MineType,
						Extra:      map[string]interface{}{},
					},
					Detail: schema.ImageURLDetailHigh,
				},
			})
		case schema.ChatMessagePartTypeAudioURL:
		case schema.ChatMessagePartTypeVideoURL:
		case schema.ChatMessagePartTypeFileURL:
		default:
			continue
		}
	}

	return nil
}

// getUploader 根据 provider 类型创建对应的上传器
func (d *Converter) getUploader() IFileUploader {
	switch d.providerModel.ProviderType {
	case data_models.ProviderTypeAliyuns:
		return &Aliyuns{
			providerModel: d.providerModel,
		}
	default:
		return nil
	}
}

// getFilePathFromPart 从不同类型的 part 中提取文件路径
func getFilePathFromPart(part *schema.MessageInputPart) (string, bool) {
	var extra map[string]interface{}

	switch part.Type {
	case schema.ChatMessagePartTypeImageURL:
		if part.Image == nil {
			return "", false
		}
		extra = part.Image.Extra
	case schema.ChatMessagePartTypeAudioURL:
		if part.Audio == nil {
			return "", false
		}
		extra = part.Audio.Extra
	case schema.ChatMessagePartTypeVideoURL:
		if part.Video == nil {
			return "", false
		}
		extra = part.Video.Extra
	case schema.ChatMessagePartTypeFileURL:
		if part.File == nil {
			return "", false
		}
		extra = part.File.Extra
	default:
		return "", false
	}

	if extra == nil {
		return "", false
	}

	pathValue, ok := extra[ExtraKeyFilePath.String()]
	if !ok || pathValue == nil {
		return "", false
	}

	path, ok := pathValue.(string)
	if !ok || path == "" {
		return "", false
	}

	return path, true
}
