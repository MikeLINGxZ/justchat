package agent

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"trpc.group/trpc-go/trpc-agent-go/model"
)

var mimeTypes = map[string]string{
	".c":        "text/x-c",
	".cpp":      "text/x-c++",
	".cs":       "text/x-csharp",
	".css":      "text/css",
	".doc":      "application/msword",
	".docx":     "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
	".gif":      "image/gif",
	".html":     "text/html",
	".java":     "text/x-java",
	".jpg":      "image/jpeg",
	".jpeg":     "image/jpeg",
	".js":       "text/javascript",
	".json":     "application/json",
	".log":      "text/plain",
	".markdown": "text/markdown",
	".md":       "text/markdown",
	".pdf":      "application/pdf",
	".php":      "text/x-php",
	".png":      "image/png",
	".pptx":     "application/vnd.openxmlformats-officedocument.presentationml.presentation",
	".py":       "text/x-python",
	".rb":       "text/x-ruby",
	".sh":       "application/x-sh",
	".tex":      "text/x-tex",
	".ts":       "application/typescript",
	".txt":      "text/plain",
	".webp":     "image/webp",
}

// Attachment 表示一条用户消息携带的本地文件附件元数据。
// 只持久化路径与轻量元信息，文件本身不入库。
type Attachment struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Mime string `json:"mime"`
	Kind string `json:"kind"` // "image" or "file"
}

// MarshalAttachments 将附件切片序列化为存库的 JSON 字符串。
// 空切片返回空串，便于在 DB 中保持自然空状态。
func MarshalAttachments(atts []Attachment) (string, error) {
	if len(atts) == 0 {
		return "", nil
	}
	b, err := json.Marshal(atts)
	if err != nil {
		return "", fmt.Errorf("marshal attachments: %w", err)
	}
	return string(b), nil
}

// UnmarshalAttachments 从存库 JSON 字符串解析附件切片。
// 空串/空白返回 nil 切片以等同 "无附件"。
func UnmarshalAttachments(s string) ([]Attachment, error) {
	if strings.TrimSpace(s) == "" {
		return nil, nil
	}
	var atts []Attachment
	if err := json.Unmarshal([]byte(s), &atts); err != nil {
		return nil, fmt.Errorf("unmarshal attachments: %w", err)
	}
	return atts, nil
}

// NormalizeAttachment 用 Path 补齐缺失的 Name/Mime/Kind。
// Name 取自 filepath.Base；Mime 先按扩展名推断，再以 application/octet-stream 兜底；
// Kind 根据 mime 前缀判定 image 还是 file。
func NormalizeAttachment(a Attachment) Attachment {
	if a.Name == "" && a.Path != "" {
		a.Name = filepath.Base(a.Path)
	}
	if a.Mime == "" {
		a.Mime = inferMimeFromPath(a.Path)
	}
	if a.Kind == "" {
		a.Kind = kindFromMime(a.Mime)
	}
	return a
}

func inferMimeFromPath(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	if m := mime.TypeByExtension(ext); m != "" {
		// mime.TypeByExtension 可能返回带 charset 的形式，截到分号即可。
		if idx := strings.Index(m, ";"); idx >= 0 {
			return strings.TrimSpace(m[:idx])
		}
		return m
	}

	mimeType, ok := mimeTypes[ext]
	if !ok {
		return "application/octet-stream"
	}

	return mimeType
}

func kindFromMime(m string) string {
	if strings.HasPrefix(m, "image/") {
		return "image"
	}
	return "file"
}

// BuildUserMessage 把文本与附件构造成 trpc model.Message。
// 附件按 Kind 分发到 AddImageData 或 AddFileData；单个文件读取失败时
// 跳过该 ContentPart，并在 Content 末尾追加 "[Missing attachment: <name>]"
// 占位文本，让模型知道附件曾经存在但已不可用。
func BuildUserMessage(content string, atts []Attachment) model.Message {
	msg := model.NewUserMessage(content)
	if len(atts) == 0 {
		return msg
	}

	for _, raw := range atts {
		a := NormalizeAttachment(raw)
		data, err := os.ReadFile(a.Path)
		if err != nil {
			msg.Content = appendMissingPlaceholder(msg.Content, a.Name)
			continue
		}
		if a.Kind == "image" {
			format := imageFormatFromMime(a.Mime)
			msg.AddImageData(data, "auto", format)
			continue
		}
		// 检测 mime 兜底：若声明的 mime 为空/octet-stream，用 http.DetectContentType 复核。
		mimeType := a.Mime
		if mimeType == "" || mimeType == "application/octet-stream" {
			head := data
			if len(head) > 512 {
				head = head[:512]
			}
			detected := http.DetectContentType(head)
			if detected != "" {
				// Strip optional "; charset=..." from detected mime.
				if idx := strings.Index(detected, ";"); idx >= 0 {
					mimeType = strings.TrimSpace(detected[:idx])
				} else {
					mimeType = detected
				}
			}
		}
		if strings.HasPrefix(mimeType, "image/") {
			msg.AddImageData(data, "auto", imageFormatFromMime(mimeType))
			continue
		}
		msg.AddFileData(a.Name, data, mimeType)
	}
	return msg
}

func appendMissingPlaceholder(content, name string) string {
	placeholder := fmt.Sprintf("[Missing attachment: %s]", name)
	if content == "" {
		return placeholder
	}
	return content + "\n" + placeholder
}

// imageFormatFromMime maps a mime type to a format hint string for the provider SDK.
// The format string is purely advisory — the SDK accepts any non-empty string and also
// handles "" gracefully. Unrecognized image mimes (e.g., "image/avif") return "".
func imageFormatFromMime(m string) string {
	switch m {
	case "image/png":
		return "png"
	case "image/jpeg":
		return "jpeg"
	case "image/webp":
		return "webp"
	case "image/gif":
		return "gif"
	}
	return ""
}
