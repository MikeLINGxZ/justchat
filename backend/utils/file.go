package utils

import (
	"encoding/base64"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/cloudwego/eino/schema"
)

const maxChatTextFileBytes = 128 * 1024

func MimeType2ChatMessagePartType(mimeType string) (schema.ChatMessagePartType, error) {
	if mimeType == "" {
		return schema.ChatMessagePartTypeText, errors.New("empty MIME type")
	}

	// 标准化：trim 空格，转小写（MIME 类型不区分大小写）
	mimeType = normalizeMimeType(mimeType)

	if IsTextMimeType(mimeType) {
		return schema.ChatMessagePartTypeText, nil
	}

	// 按层级匹配：先精确匹配常用特例（可选），再通配符前缀匹配
	switch mimeType {
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

func IsTextMimeType(mimeType string) bool {
	mimeType = normalizeMimeType(mimeType)
	if mimeType == "" {
		return false
	}
	if strings.HasPrefix(mimeType, "text/") {
		return true
	}

	switch mimeType {
	case "application/json",
		"application/ld+json",
		"application/x-ndjson",
		"application/xml",
		"application/x-yaml",
		"application/yaml",
		"application/toml",
		"application/x-toml",
		"application/csv",
		"application/javascript",
		"application/x-javascript",
		"application/sql",
		"application/x-sh",
		"application/x-shellscript",
		"application/x-httpd-php",
		"application/x-env",
		"application/x-ini",
		"application/x-properties",
		"application/x-empty":
		return true
	}

	return strings.HasSuffix(mimeType, "+json") || strings.HasSuffix(mimeType, "+xml")
}

func DetectMimeType(path string) (string, error) {
	extMimeType := normalizeMimeType(mime.TypeByExtension(strings.ToLower(filepath.Ext(path))))
	if extMimeType != "" && extMimeType != "application/octet-stream" {
		return extMimeType, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	if len(data) == 0 {
		return "text/plain", nil
	}

	sniffed := normalizeMimeType(http.DetectContentType(data))
	if sniffed != "" && sniffed != "application/octet-stream" {
		if strings.HasPrefix(sniffed, "text/") || utf8.Valid(data) {
			return sniffed, nil
		}
		return sniffed, nil
	}

	if utf8.Valid(data) {
		return "text/plain", nil
	}

	return "application/octet-stream", nil
}

func ReadFile2Base64Data(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err // 保留原始错误（含路径、权限、不存在等信息）
	}
	return base64.StdEncoding.EncodeToString(data), nil
}

func ReadTextFileForChat(path, name, mimeType string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	truncated := false
	if len(data) > maxChatTextFileBytes {
		data = data[:maxChatTextFileBytes]
		truncated = true
	}

	content := string(data)
	note := ""
	if truncated {
		note = fmt.Sprintf("\n[truncated: showing first %d bytes]", maxChatTextFileBytes)
	}

	return fmt.Sprintf(
		"<attached_text_file name=%q mime_type=%q>\n%s%s\n</attached_text_file>",
		name,
		normalizeMimeType(mimeType),
		content,
		note,
	), nil
}

func normalizeMimeType(mimeType string) string {
	mimeType = strings.TrimSpace(strings.ToLower(mimeType))
	if mimeType == "" {
		return ""
	}
	mediaType, _, err := mime.ParseMediaType(mimeType)
	if err == nil {
		return mediaType
	}
	return mimeType
}
