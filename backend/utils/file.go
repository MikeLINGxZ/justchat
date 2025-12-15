package utils

import (
	"errors"
	"strings"

	"github.com/cloudwego/eino/schema"
)

func MimeType2ChatMessagePartType(mimeType string) (schema.ChatMessagePartType, error) {
	if mimeType == "" {
		return schema.ChatMessagePartTypeText, errors.New("empty MIME type")
	}

	// 标准化：trim 空格，转小写（MIME 类型不区分大小写）
	mimeType = strings.TrimSpace(strings.ToLower(mimeType))

	// 按层级匹配：先精确匹配常用特例（可选），再通配符前缀匹配
	switch mimeType {
	// Text
	case "text/plain", "text/markdown", "text/html", "text/css", "text/javascript", "text/typescript":
		return schema.ChatMessagePartTypeText, nil

	// Image
	case "image/svg+xml":
		return schema.ChatMessagePartTypeImageURL, nil
	default:
		if strings.HasPrefix(mimeType, "image/") {
			return schema.ChatMessagePartTypeImageURL, nil
		}
	}

	// Audio
	if strings.HasPrefix(mimeType, "audio/") {
		return schema.ChatMessagePartTypeAudioURL, nil
	}

	// Video
	if strings.HasPrefix(mimeType, "video/") {
		return schema.ChatMessagePartTypeVideoURL, nil
	}

	// All others → file_url (e.g., application/pdf, application/json, model/gltf-binary, etc.)
	return schema.ChatMessagePartTypeFileURL, nil
}
