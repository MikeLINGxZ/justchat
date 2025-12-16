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
	if message == nil {
		return nil
	}

	uploader := d.getUploader()
	if uploader == nil {
		return nil
	}

	paths := d.getPaths(message)
	if len(paths) == 0 {
		return nil
	}

	path2url, err := uploader.Upload(paths)
	if err != nil {
		return err
	}

	for i := range message.UserInputMultiContent {
		part := &message.UserInputMultiContent[i]
		if !isFilePartType(part.Type) {
			continue
		}

		path, ok := getFilePathFromPart(part)
		if !ok {
			continue
		}

		url := path2url[path]
		setURLToPart(part, url)
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

// setURLToPart 设置不同类型的 part 的 URL
func setURLToPart(part *schema.MessageInputPart, url string) bool {
	if url == "" {
		return false
	}

	switch part.Type {
	case schema.ChatMessagePartTypeImageURL:
		if part.Image != nil {
			part.Image.URL = &url
			return true
		}
	case schema.ChatMessagePartTypeAudioURL:
		if part.Audio != nil {
			part.Audio.URL = &url
			return true
		}
	case schema.ChatMessagePartTypeVideoURL:
		if part.Video != nil {
			part.Video.URL = &url
			return true
		}
	case schema.ChatMessagePartTypeFileURL:
		if part.File != nil {
			part.File.URL = &url
			return true
		}
	}

	return false
}

// isFilePartType 检查是否为需要上传的文件类型
func isFilePartType(partType schema.ChatMessagePartType) bool {
	switch partType {
	case schema.ChatMessagePartTypeImageURL,
		schema.ChatMessagePartTypeAudioURL,
		schema.ChatMessagePartTypeVideoURL,
		schema.ChatMessagePartTypeFileURL:
		return true
	default:
		return false
	}
}

func (d *Converter) getPaths(message *schema.Message) []string {
	if message == nil {
		return nil
	}

	var paths []string
	for i := range message.UserInputMultiContent {
		part := &message.UserInputMultiContent[i]
		if !isFilePartType(part.Type) {
			continue
		}

		path, ok := getFilePathFromPart(part)
		if ok {
			paths = append(paths, path)
		}
	}

	return paths
}
