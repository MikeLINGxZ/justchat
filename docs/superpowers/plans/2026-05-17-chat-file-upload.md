# 对话附件上传（多模态）实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 让用户在 ChatInput 通过 Paperclip 选择本地文件（图片/PDF/文本等），与文字一并发送给模型，并在历史会话中正确重放。

**Architecture:** 仅存绝对路径于消息的新列 `Attachments`（JSON 文本，gorm AutoMigrate）；发送和重放都经 `pkg/agent.BuildUserMessage` 把附件读盘后注入 `model.Message.ContentParts`。图片走 `model.AddImageData`，其他走 `model.AddFileData`（绕开 SDK `AddFilePath` 仅支持白名单扩展的限制）。前端用本地态收集附件、chip 形式展示，发送时透传到后端。

**Tech Stack:** Go (Wails v3, GORM SQLite, `trpc-agent-go` v1.9.1), TypeScript (React 18, Zustand, Tiptap, Vitest + RTL).

**Spec:** `docs/superpowers/specs/2026-05-17-chat-file-upload-design.md`

**Important deviation from spec wording:**
spec 第 3 节描述用 `AddImageFilePath` / `AddFilePath`。实际实现用 `AddImageData` / `AddFileData`：`AddFilePath` 仅支持 `.txt/.md/.pdf/.py/.json/...` 等固定扩展，而 spec 承诺“任意类型”。改用 *Data 变体后我们自己做 mime 推断（`mime.TypeByExtension` + `http.DetectContentType` 兜底），即可支持 `.go/.yaml/.toml` 等任意文件。

---

## 文件结构

| 类型 | 路径                                                         | 职责                                                                                                      |
| ---- | ------------------------------------------------------------ | --------------------------------------------------------------------------------------------------------- |
| 修改 | `backend/models/data_models/message.go`                    | 增加 `Attachments string` 字段                                                                          |
| 新建 | `backend/pkg/agent/attachment.go`                          | `Attachment` 类型 + 序列化/反序列化 + `BuildUserMessage` + mime/kind 推断                             |
| 新建 | `backend/pkg/agent/attachment_test.go`                     | 上述函数的单元测试                                                                                        |
| 修改 | `backend/pkg/agent/chat_handler.go`                        | `SendMessageParams` 增加 `Attachments`；`SendMessage` 与 `toModelMessage` 走 `BuildUserMessage` |
| 修改 | `backend/pkg/agent/chat_handler_test.go`                   | 新增 `SendMessage` 带附件的回归                                                                         |
| 新建 | `backend/service/agent/agent_dto/attachment_input.go`      | DTO 类型 `AttachmentInput`                                                                              |
| 修改 | `backend/service/agent/agent_dto/send_message.go`          | 增加 `Attachments []AttachmentInput`                                                                    |
| 修改 | `backend/service/agent/agent_dto/load_session_messages.go` | `MessageItem` 增加 `Attachments string`                                                               |
| 修改 | `backend/service/agent/agent_internal.go`                  | `toMessageItem` 复制 `Attachments`                                                                    |
| 修改 | `backend/service/agent/agent.go`                           | `SendMessage` 校验 + 转换附件                                                                           |
| 修改 | `backend/service/agent/agent_test.go`                      | 校验路径回归                                                                                              |
| 修改 | `backend/service/file/file.go`                             | 修复 `SelectFile` 三个 bug                                                                              |
| 修改 | `frontend/src/types/index.ts`                              | `Attachment` / `AttachmentKind` / `Message.attachments`                                             |
| 新建 | `frontend/src/lib/attachments.ts`                          | `inferAttachmentMeta` + 常量 + mime 判别                                                                |
| 新建 | `frontend/src/__tests__/attachmentsLib.test.ts`            | 上述工具函数测试                                                                                          |
| 新建 | `frontend/src/components/chat/AttachmentChips.tsx`         | chip 列表组件                                                                                             |
| 新建 | `frontend/src/__tests__/attachmentChips.test.tsx`          | 组件渲染测试                                                                                              |
| 修改 | `frontend/src/store/chatStore.ts`                          | `sendMessage` 透传 `attachments`；`loadMessages` 反序列化                                           |
| 修改 | `frontend/src/components/chat/ChatInput.tsx`               | Paperclip onClick + chips 上方渲染                                                                        |
| 新建 | `frontend/src/__tests__/chatInputAttachments.test.tsx`     | ChatInput 附件交互测试                                                                                    |
| 修改 | `frontend/src/components/chat/MessageItem.tsx`             | user/text 气泡内渲染只读 chips                                                                            |
| 修改 | `frontend/src/i18n/locales/zh-CN.ts`                       | 增加 `input.attach*` 与 `chat.attachment*`                                                            |
| 修改 | `frontend/src/i18n/locales/en.ts`                          | 同上                                                                                                      |

---

## Phase 1 — 后端数据与协议

### Task 1: 在 `Message` 数据模型上增加 `Attachments` 字段

**Files:**

- Modify: `backend/models/data_models/message.go`

- [ ] **Step 1: 修改 Message 结构体**

把 `backend/models/data_models/message.go` 改为：

```go
package data_models

type Message struct {
	OrmModel
	SessionID   uint   `gorm:"index" json:"session_id"`
	ParentID    *uint  `gorm:"index" json:"parent_id"`
	Role        string `gorm:"type:varchar(50)" json:"role"`
	ContentType string `gorm:"type:varchar(50)" json:"content_type"`
	Content     string `gorm:"type:text" json:"content"`
	ModelName   string `gorm:"type:varchar(255)" json:"model_name"`
	AgentName   string `gorm:"type:varchar(255)" json:"agent_name"`
	TokensIn    int    `json:"tokens_in"`
	TokensOut   int    `json:"tokens_out"`
	Extra       string `gorm:"type:text" json:"extra"`
	Attachments string `gorm:"type:text" json:"attachments"`
}
```

- [ ] **Step 2: 确认编译通过**

Run: `go build ./...`
Expected: 无错误，零输出。

- [ ] **Step 3: 提交**

```bash
git add backend/models/data_models/message.go
git commit -m "feat(data): add Attachments column on Message"
```

---

### Task 2: 在 `pkg/agent` 增加附件类型与 `BuildUserMessage`

**Files:**

- Create: `backend/pkg/agent/attachment.go`
- Create: `backend/pkg/agent/attachment_test.go`

- [ ] **Step 1: 写失败测试 `attachment_test.go`**

```go
package agent

import (
	"encoding/base64"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"trpc.group/trpc-go/trpc-agent-go/model"
)

// 67-byte minimal valid PNG (1x1 transparent)
const minimalPNGBase64 = "iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR4nGNgYGD4DwABBAEAfbLI3wAAAABJRU5ErkJggg=="

func writeTempFile(t *testing.T, name string, data []byte) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write %s: %v", name, err)
	}
	return path
}

func TestMarshalUnmarshalAttachmentsRoundTrip(t *testing.T) {
	in := []Attachment{
		{Name: "a.png", Path: "/abs/a.png", Mime: "image/png", Kind: "image"},
		{Name: "b.pdf", Path: "/abs/b.pdf", Mime: "application/pdf", Kind: "file"},
	}

	s, err := MarshalAttachments(in)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	if s == "" {
		t.Fatalf("expected non-empty json, got empty")
	}

	out, err := UnmarshalAttachments(s)
	if err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(out) != 2 || out[0].Name != "a.png" || out[1].Kind != "file" {
		t.Fatalf("round trip mismatch: %+v", out)
	}
}

func TestMarshalEmptyReturnsEmptyString(t *testing.T) {
	got, err := MarshalAttachments(nil)
	if err != nil || got != "" {
		t.Fatalf("expected empty,nil; got %q,%v", got, err)
	}
}

func TestUnmarshalEmptyReturnsNil(t *testing.T) {
	got, err := UnmarshalAttachments("")
	if err != nil || got != nil {
		t.Fatalf("expected nil,nil; got %v,%v", got, err)
	}
}

func TestNormalizeAttachmentFillsMissingFields(t *testing.T) {
	a := NormalizeAttachment(Attachment{Path: "/tmp/x.PNG"})
	if a.Name != "x.PNG" {
		t.Fatalf("name = %q", a.Name)
	}
	if a.Mime != "image/png" {
		t.Fatalf("mime = %q", a.Mime)
	}
	if a.Kind != "image" {
		t.Fatalf("kind = %q", a.Kind)
	}

	b := NormalizeAttachment(Attachment{Path: "/tmp/data.bin"})
	if b.Mime == "" || b.Kind != "file" {
		t.Fatalf("unknown ext should be file, got mime=%q kind=%q", b.Mime, b.Kind)
	}
}

func TestBuildUserMessageNoAttachments(t *testing.T) {
	msg := BuildUserMessage("hello", nil)
	if msg.Role != model.RoleUser {
		t.Fatalf("role = %v", msg.Role)
	}
	if msg.Content != "hello" {
		t.Fatalf("content = %q", msg.Content)
	}
	if len(msg.ContentParts) != 0 {
		t.Fatalf("expected no ContentParts, got %d", len(msg.ContentParts))
	}
}

func TestBuildUserMessageWithImage(t *testing.T) {
	png, _ := base64.StdEncoding.DecodeString(minimalPNGBase64)
	path := writeTempFile(t, "tiny.png", png)

	msg := BuildUserMessage("see image", []Attachment{
		NormalizeAttachment(Attachment{Path: path}),
	})

	if msg.Content != "see image" {
		t.Fatalf("content = %q", msg.Content)
	}
	if len(msg.ContentParts) != 1 || msg.ContentParts[0].Type != model.ContentTypeImage {
		t.Fatalf("expected one image part, got %+v", msg.ContentParts)
	}
	if msg.ContentParts[0].Image == nil || len(msg.ContentParts[0].Image.Data) == 0 {
		t.Fatalf("image data empty")
	}
}

func TestBuildUserMessageWithArbitraryFile(t *testing.T) {
	path := writeTempFile(t, "notes.go", []byte("package x"))

	msg := BuildUserMessage("review please", []Attachment{
		NormalizeAttachment(Attachment{Path: path}),
	})

	if len(msg.ContentParts) != 1 || msg.ContentParts[0].Type != model.ContentTypeFile {
		t.Fatalf("expected file part, got %+v", msg.ContentParts)
	}
	if msg.ContentParts[0].File == nil || msg.ContentParts[0].File.Name != "notes.go" {
		t.Fatalf("file part incorrect: %+v", msg.ContentParts[0].File)
	}
}

func TestBuildUserMessageMissingFileDegrades(t *testing.T) {
	msg := BuildUserMessage("with missing", []Attachment{
		{Name: "ghost.png", Path: "/nonexistent/ghost.png", Mime: "image/png", Kind: "image"},
	})

	if len(msg.ContentParts) != 0 {
		t.Fatalf("expected no ContentParts, got %d", len(msg.ContentParts))
	}
	if !strings.Contains(msg.Content, "[Missing attachment: ghost.png]") {
		t.Fatalf("content missing placeholder: %q", msg.Content)
	}
}

func TestBuildUserMessageMixedSomeMissing(t *testing.T) {
	png, _ := base64.StdEncoding.DecodeString(minimalPNGBase64)
	good := writeTempFile(t, "ok.png", png)

	msg := BuildUserMessage("mixed", []Attachment{
		NormalizeAttachment(Attachment{Path: good}),
		{Name: "lost.pdf", Path: "/nope/lost.pdf", Mime: "application/pdf", Kind: "file"},
	})

	if len(msg.ContentParts) != 1 {
		t.Fatalf("expected 1 ContentPart, got %d", len(msg.ContentParts))
	}
	if !strings.Contains(msg.Content, "[Missing attachment: lost.pdf]") {
		t.Fatalf("content missing placeholder: %q", msg.Content)
	}
}
```

- [ ] **Step 2: 运行测试确认失败**

Run: `go test ./backend/pkg/agent/ -run Attachment -count=1`
Expected: 编译失败，未定义 `Attachment`、`MarshalAttachments` 等符号。

- [ ] **Step 3: 实现 `attachment.go`**

创建 `backend/pkg/agent/attachment.go`：

```go
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
// 空串/空白返回 nil 切片以等同 “无附件”。
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
	switch ext {
	// mime 包在部分平台对常见扩展没有内建映射，兜底显式列举。
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".webp":
		return "image/webp"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	case ".md", ".markdown":
		return "text/markdown"
	case ".txt", ".log":
		return "text/plain"
	}
	return "application/octet-stream"
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
				mimeType = detected
			}
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
```

- [ ] **Step 4: 运行测试确认通过**

Run: `go test ./backend/pkg/agent/ -run Attachment -count=1 -v`
Expected: 全部 PASS。

- [ ] **Step 5: 提交**

```bash
git add backend/pkg/agent/attachment.go backend/pkg/agent/attachment_test.go
git commit -m "feat(agent): add Attachment helpers and BuildUserMessage"
```

---

### Task 3: 增加 `AttachmentInput` DTO

**Files:**

- Create: `backend/service/agent/agent_dto/attachment_input.go`

- [ ] **Step 1: 创建文件**

```go
package agent_dto

// AttachmentInput 是前端发送消息时携带的附件元数据。
// Name/Mime 可空，后端通过 path 推断后写入存储与 model.Message。
type AttachmentInput struct {
	Path string `json:"path"`
	Name string `json:"name,omitempty"`
	Mime string `json:"mime,omitempty"`
}
```

- [ ] **Step 2: 编译验证**

Run: `go build ./backend/...`
Expected: 无错误。

- [ ] **Step 3: 提交**

```bash
git add backend/service/agent/agent_dto/attachment_input.go
git commit -m "feat(agent_dto): add AttachmentInput type"
```

---

### Task 4: 把 `Attachments` 加入 `SendMessageInput` 与 `MessageItem`

**Files:**

- Modify: `backend/service/agent/agent_dto/send_message.go`
- Modify: `backend/service/agent/agent_dto/load_session_messages.go`
- Modify: `backend/service/agent/agent_internal.go`

- [ ] **Step 1: 修改 `send_message.go`**

```go
package agent_dto

type SendMessageInput struct {
	SessionID        uint              `json:"session_id"`
	Content          string            `json:"content"`
	BaseURL          string            `json:"base_url"`
	ApiKey           string            `json:"api_key"`
	ModelName        string            `json:"model_name"`
	EnabledUserTools []string          `json:"enabled_user_tools"`
	Attachments      []AttachmentInput `json:"attachments,omitempty"`
}

type SendMessageOutput struct {
}
```

- [ ] **Step 2: 修改 `load_session_messages.go`，给 `MessageItem` 增加 `Attachments`**

```go
type MessageItem struct {
	ID          uint   `json:"id"`
	SessionID   uint   `json:"session_id"`
	ParentID    *uint  `json:"parent_id"`
	Role        string `json:"role"`
	ContentType string `json:"content_type"`
	Content     string `json:"content"`
	ModelName   string `json:"model_name"`
	AgentName   string `json:"agent_name"`
	TokensIn    int    `json:"tokens_in"`
	TokensOut   int    `json:"tokens_out"`
	Extra       string `json:"extra"`
	Attachments string `json:"attachments"`
	Created     string `json:"created"`
}
```

(其余字段顺序保留不变；仅新增 `Attachments` 一行在 `Extra` 之后。)

- [ ] **Step 3: 修改 `agent_internal.go` 的 `toMessageItem`，复制新字段**

把 `toMessageItem` 改为：

```go
func toMessageItem(m data_models.Message) agent_dto.MessageItem {
	return agent_dto.MessageItem{
		ID:          m.ID,
		SessionID:   m.SessionID,
		ParentID:    m.ParentID,
		Role:        m.Role,
		ContentType: m.ContentType,
		Content:     m.Content,
		ModelName:   m.ModelName,
		AgentName:   m.AgentName,
		TokensIn:    m.TokensIn,
		TokensOut:   m.TokensOut,
		Extra:       m.Extra,
		Attachments: m.Attachments,
		Created:     m.CreatedAt.Format(timeFormat),
	}
}
```

- [ ] **Step 4: 编译验证**

Run: `go build ./backend/...`
Expected: 无错误。

- [ ] **Step 5: 提交**

```bash
git add backend/service/agent/agent_dto/send_message.go backend/service/agent/agent_dto/load_session_messages.go backend/service/agent/agent_internal.go
git commit -m "feat(agent_dto): wire Attachments into SendMessage and MessageItem"
```

---

## Phase 2 — 后端业务逻辑

### Task 5: `ChatHandler` 接受附件并持久化

**Files:**

- Modify: `backend/pkg/agent/chat_handler.go`

- [ ] **Step 1: 给 `SendMessageParams` 加字段**

把 `chat_handler.go` 里 `SendMessageParams` 改为：

```go
// SendMessageParams defines the inputs required to run a chat turn.
type SendMessageParams struct {
	SessionID        uint
	Content          string
	BaseURL          string
	ApiKey           string
	ModelName        string
	EnabledUserTools []string
	Attachments      []Attachment
}
```

- [ ] **Step 2: 改 `SendMessage` 使用 `BuildUserMessage` 并持久化附件**

把 `SendMessage` 方法体替换为下面版本（仅有两处变化：`CreateMessage` 的 `Attachments` 字段，和 `r.Run` 的 user message 改用 `BuildUserMessage`）：

```go
// SendMessage starts a streaming run and consumes events asynchronously.
func (ch *ChatHandler) SendMessage(ctx context.Context, params SendMessageParams) error {
	stor := ch.manager.Storage()
	history, err := ch.loadSessionHistory(params.SessionID)
	if err != nil {
		return err
	}

	attachmentsJSON, err := MarshalAttachments(params.Attachments)
	if err != nil {
		return fmt.Errorf("marshal attachments: %w", err)
	}

	_, err = stor.CreateMessage(data_models.Message{
		SessionID:   params.SessionID,
		Role:        "user",
		ContentType: "text",
		Content:     params.Content,
		Attachments: attachmentsJSON,
	})
	if err != nil {
		return fmt.Errorf("save user message: %w", err)
	}

	_ = stor.TouchSession(params.SessionID)
	ch.updateSessionStatus(params.SessionID, "loading")

	r, err := ch.manager.GetOrCreateRunner(params.BaseURL, params.ApiKey, params.ModelName, params.EnabledUserTools)
	if err != nil {
		ch.updateSessionStatus(params.SessionID, "error-unread")
		return fmt.Errorf("get runner: %w", err)
	}

	streamCtx, cancel := ch.manager.Streams().Start(context.WithoutCancel(ctx), params.SessionID)
	userID := "local"
	sessionIDStr := strconv.FormatUint(uint64(params.SessionID), 10)

	runOptions := make([]agentpkg.RunOption, 0, 1)
	if len(history) > 0 {
		runOptions = append(runOptions, agentpkg.WithMessages(history))
	}

	runCtx := withRunContext(streamCtx, runContext{
		SessionID: params.SessionID,
		ModelName: params.ModelName,
	})

	userMsg := BuildUserMessage(params.Content, params.Attachments)

	events, err := r.Run(
		runCtx,
		userID,
		sessionIDStr,
		userMsg,
		runOptions...,
	)
	if err != nil {
		cancel()
		ch.manager.Streams().Remove(params.SessionID)
		ch.updateSessionStatus(params.SessionID, "error-unread")
		return fmt.Errorf("runner.Run: %w", err)
	}

	go ch.consumeEvents(streamCtx, params.SessionID, params.ModelName, events)
	return nil
}
```

- [ ] **Step 3: 改 `toModelMessage` 在 user/text 分支看附件**

把 `toModelMessage` 中 `case "text":` 分支替换为：

```go
	case "text":
		if stored.Role == "user" {
			atts, err := UnmarshalAttachments(stored.Attachments)
			if err == nil && len(atts) > 0 {
				return BuildUserMessage(stored.Content, atts), true
			}
			return model.NewUserMessage(stored.Content), true
		}
		return model.NewAssistantMessage(stored.Content), true
```

(其他 case 不动。)

- [ ] **Step 4: 编译验证**

Run: `go build ./backend/...`
Expected: 无错误。

- [ ] **Step 5: 运行已有 chat_handler 单测，确保未回归**

Run: `go test ./backend/pkg/agent/ -count=1`
Expected: 已有测试全 PASS（含 Task 2 新测）。

- [ ] **Step 6: 提交**

```bash
git add backend/pkg/agent/chat_handler.go
git commit -m "feat(agent): persist attachments and build multimodal user message"
```

---

### Task 6: 给 `ChatHandler` 加一条带附件的回归测试

**Files:**

- Modify: `backend/pkg/agent/chat_handler_test.go`

- [ ] **Step 1: 写新测试（使用现有的 `manager.runners` 注入模式）**

在 `chat_handler_test.go` 末尾追加：

```go
func TestSendMessagePersistsAttachmentsAndBuildsMultimodal(t *testing.T) {
	manager := NewManager(newTestStorage(t))
	handler := NewChatHandler(manager)

	session, err := manager.Storage().CreateSession(data_models.Session{
		Title:  "Attach",
		Status: "idle",
	})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	r := newDelayedRunner()
	// 注入 runner，绕过真实 GetOrCreateRunner。Key 格式见 manager.GetOrCreateRunner。
	manager.runners["https://api.example.com/m"] = r

	tmpDir := t.TempDir()
	imgPath := filepath.Join(tmpDir, "tiny.png")
	pngBytes, _ := base64.StdEncoding.DecodeString(minimalPNGBase64)
	if err := os.WriteFile(imgPath, pngBytes, 0o644); err != nil {
		t.Fatalf("write png: %v", err)
	}

	err = handler.SendMessage(context.Background(), SendMessageParams{
		SessionID: session.ID,
		Content:   "describe",
		BaseURL:   "https://api.example.com",
		ApiKey:    "x",
		ModelName: "m",
		Attachments: []Attachment{
			NormalizeAttachment(Attachment{Path: imgPath}),
		},
	})
	if err != nil {
		t.Fatalf("send: %v", err)
	}

	// 1. 持久化的 user 消息 Attachments 是合法 JSON 且能解出 1 个图片项。
	msgs, err := manager.Storage().ListMessagesForSession(session.ID, 0, 10)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(msgs) != 1 {
		t.Fatalf("want 1 message, got %d", len(msgs))
	}
	got, err := UnmarshalAttachments(msgs[0].Attachments)
	if err != nil || len(got) != 1 || got[0].Kind != "image" {
		t.Fatalf("attachments persisted incorrectly: %q -> %+v err=%v", msgs[0].Attachments, got, err)
	}

	// 2. 传给 runner 的 user message 含 ContentParts。
	if len(r.lastMessage.ContentParts) != 1 ||
		r.lastMessage.ContentParts[0].Type != model.ContentTypeImage {
		t.Fatalf("runner did not receive multimodal user message: %+v", r.lastMessage)
	}
}
```

并在文件顶部 import 块补足（保留已有项，新增 `encoding/base64`、`os`、`path/filepath`，其余如 `model` 已存在不要重复）：

```go
import (
	"context"
	"encoding/base64"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	agentpkg "trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/model"
)
```

(`minimalPNGBase64` 与 `writeTempFile` 已在 Task 2 的 `attachment_test.go` 里定义，同 package 可直接复用——不要在此文件重复定义。)

- [ ] **Step 2: 运行测试**

Run: `go test ./backend/pkg/agent/ -run TestSendMessagePersistsAttachments -count=1 -v`
Expected: PASS。

- [ ] **Step 3: 提交**

```bash
git add backend/pkg/agent/chat_handler_test.go
git commit -m "test(agent): cover SendMessage attachments persistence and multimodal build"
```

---

### Task 7: `agent` service 校验 + 转换附件

**Files:**

- Modify: `backend/service/agent/agent.go`

- [ ] **Step 1: 在 `agent.go` 顶部 import 区追加**

```go
import (
	"context"
	"fmt"
	"os"
	"strings"

	pkgAgent "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/agent"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/agent/agent_dto"
	// ... 其余按原文件保留
)
```

(原 `context`, `fmt`, `strings`, `pkgAgent`, `agent_dto` 已存在；仅新增 `os`。其余 import 不动。)

- [ ] **Step 2: 在文件末尾追加常量与转换函数**

```go
const (
	maxAttachmentCount = 10
	maxAttachmentBytes = 20 * 1024 * 1024
)

// convertAttachments 校验数量、大小、存在性，并把 DTO 转换为 pkgAgent.Attachment。
func convertAttachments(in []agent_dto.AttachmentInput) ([]pkgAgent.Attachment, error) {
	if len(in) == 0 {
		return nil, nil
	}
	if len(in) > maxAttachmentCount {
		return nil, fmt.Errorf("too many attachments: %d (max %d)", len(in), maxAttachmentCount)
	}
	out := make([]pkgAgent.Attachment, 0, len(in))
	for _, item := range in {
		path := strings.TrimSpace(item.Path)
		if path == "" {
			return nil, fmt.Errorf("attachment path is empty")
		}
		info, err := os.Stat(path)
		if err != nil {
			return nil, fmt.Errorf("attachment not found: %s", path)
		}
		if info.Size() > maxAttachmentBytes {
			return nil, fmt.Errorf("attachment %q exceeds %d bytes", info.Name(), maxAttachmentBytes)
		}
		out = append(out, pkgAgent.NormalizeAttachment(pkgAgent.Attachment{
			Path: path,
			Name: item.Name,
			Mime: item.Mime,
		}))
	}
	return out, nil
}
```

- [ ] **Step 3: 在 `Agent.SendMessage` 里调用校验**

把 `Agent.SendMessage` 方法替换为：

```go
func (a *Agent) SendMessage(ctx context.Context, input agent_dto.SendMessageInput) (*agent_dto.SendMessageOutput, error) {
	atts, err := convertAttachments(input.Attachments)
	if err != nil {
		return nil, fmt.Errorf("send message: %w", err)
	}

	handler := a.manager.NewChatHandler()
	err = handler.SendMessage(ctx, pkgAgent.SendMessageParams{
		SessionID:        input.SessionID,
		Content:          input.Content,
		BaseURL:          input.BaseURL,
		ApiKey:           input.ApiKey,
		ModelName:        input.ModelName,
		EnabledUserTools: input.EnabledUserTools,
		Attachments:      atts,
	})
	if err != nil {
		return nil, fmt.Errorf("send message: %w", err)
	}
	return &agent_dto.SendMessageOutput{}, nil
}
```

- [ ] **Step 4: 编译**

Run: `go build ./backend/...`
Expected: 无错误。

- [ ] **Step 5: 提交**

```bash
git add backend/service/agent/agent.go
git commit -m "feat(agent): validate and forward attachments to chat handler"
```

---

### Task 8: 给 `agent` service 写校验回归测试

**Files:**

- Modify: `backend/service/agent/agent_test.go`

- [ ] **Step 1: 在 `agent_test.go` 中追加**

```go
func TestConvertAttachmentsRejectsTooMany(t *testing.T) {
	in := make([]agent_dto.AttachmentInput, maxAttachmentCount+1)
	for i := range in {
		in[i] = agent_dto.AttachmentInput{Path: "/tmp/x"}
	}
	_, err := convertAttachments(in)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "too many") {
		t.Fatalf("error = %v", err)
	}
}

func TestConvertAttachmentsRejectsMissingFile(t *testing.T) {
	_, err := convertAttachments([]agent_dto.AttachmentInput{
		{Path: "/definitely/not/here.png"},
	})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "not found") {
		t.Fatalf("error = %v", err)
	}
}

func TestConvertAttachmentsRejectsOversize(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "big.bin")
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := f.Truncate(maxAttachmentBytes + 1); err != nil {
		t.Fatalf("truncate: %v", err)
	}
	_ = f.Close()

	_, err = convertAttachments([]agent_dto.AttachmentInput{{Path: path}})
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), "exceeds") {
		t.Fatalf("error = %v", err)
	}
}

func TestConvertAttachmentsNormalizesAndPasses(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "hello.go")
	if err := os.WriteFile(path, []byte("package x"), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	got, err := convertAttachments([]agent_dto.AttachmentInput{{Path: path}})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 1 || got[0].Name != "hello.go" || got[0].Kind != "file" {
		t.Fatalf("normalize failed: %+v", got)
	}
}
```

import 部分需补 `os`、`path/filepath`、`strings`、`testing`（保留已存在的）；以及

```go
"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/agent/agent_dto"
```

如已存在则忽略。

- [ ] **Step 2: 运行测试**

Run: `go test ./backend/service/agent/ -run TestConvertAttachments -count=1 -v`
Expected: 全 PASS。

- [ ] **Step 3: 提交**

```bash
git add backend/service/agent/agent_test.go
git commit -m "test(agent): cover attachment validation paths"
```

---

### Task 9: 修复 `SelectFile` 的三处 bug

**Files:**

- Modify: `backend/service/file/file.go`

- [ ] **Step 1: 把 `SelectFile` 替换为修复版**

```go
// SelectFile select a file
func (f *File) SelectFile(ctx context.Context, input file_dto.SelectFileInput) (*file_dto.SelectFileOutput, error) {
	if input.DefaultFolderPath == "" {
		input.DefaultFolderPath = dir.HomeDir()
	}
	path, err := f.wailsApp.Dialog.OpenFile().
		SetTitle(i18n.TCurrent("select.file", nil)).
		SetDirectory(input.DefaultFolderPath).
		CanChooseFiles(true).
		PromptForSingleSelection()
	if err != nil || path == "" {
		return nil, err
	}

	return &file_dto.SelectFileOutput{FilePath: path}, nil
}
```

修了三处：双重 `i18n.TCurrent` → 单层；key `select.folder` → `select.file`；`CanChooseFiles(false)` → `true`。

`SelectFolder` 保留原样（双重 t 也是 bug 但不在本次范围）。

- [ ] **Step 2: 编译**

Run: `go build ./backend/...`
Expected: 无错误。

- [ ] **Step 3: 提交**

```bash
git add backend/service/file/file.go
git commit -m "fix(file): correct SelectFile title key and CanChooseFiles flag"
```

---

## Phase 3 — 前端类型与工具

### Task 10: 前端类型扩展

**Files:**

- Modify: `frontend/src/types/index.ts`

- [ ] **Step 1: 在 `frontend/src/types/index.ts` 增加附件类型并扩展 `Message`**

在 `Message` 定义之前加入：

```ts
export type AttachmentKind = 'image' | 'file'

export type Attachment = {
  path: string
  name: string
  mime: string
  kind: AttachmentKind
}
```

并把 `Message` 改为：

```ts
export type Message = {
  id: number
  sessionId: number
  parentId: number | null
  role: MessageRole
  contentType: MessageContentType
  content: string
  modelName: string
  agentName: string
  tokensIn: number
  tokensOut: number
  extra: string
  attachments?: Attachment[]
  createdAt: string
}
```

- [ ] **Step 2: 编译检查**

Run: `cd frontend && cnpm run typecheck`
Expected: 无错误。（若仓库没有 typecheck script，用 `cnpm exec tsc --noEmit` 替代。）

- [ ] **Step 3: 提交**

```bash
git add frontend/src/types/index.ts
git commit -m "feat(types): add Attachment and Message.attachments"
```

---

### Task 11: 前端附件工具函数

**Files:**

- Create: `frontend/src/lib/attachments.ts`
- Create: `frontend/src/__tests__/attachmentsLib.test.ts`

- [ ] **Step 1: 写失败测试**

`frontend/src/__tests__/attachmentsLib.test.ts`：

```ts
import { describe, expect, it } from 'vitest'
import {
  inferAttachmentMeta,
  ATTACHMENT_MAX_COUNT,
  ATTACHMENT_MAX_BYTES,
  isImageMime,
  isPdfMime,
} from '@/lib/attachments'

describe('inferAttachmentMeta', () => {
  it('classifies png as image', () => {
    const meta = inferAttachmentMeta('/foo/bar/baz.PNG')
    expect(meta.name).toBe('baz.PNG')
    expect(meta.mime).toBe('image/png')
    expect(meta.kind).toBe('image')
  })

  it('classifies pdf as file with application/pdf', () => {
    const meta = inferAttachmentMeta('/x/spec.pdf')
    expect(meta.kind).toBe('file')
    expect(meta.mime).toBe('application/pdf')
  })

  it('classifies unknown extension as octet-stream file', () => {
    const meta = inferAttachmentMeta('/x/data.xyz')
    expect(meta.kind).toBe('file')
    expect(meta.mime).toBe('application/octet-stream')
  })

  it('handles windows-style backslash path', () => {
    const meta = inferAttachmentMeta('C:\\Users\\me\\notes.md')
    expect(meta.name).toBe('notes.md')
    expect(meta.kind).toBe('file')
  })
})

describe('mime helpers', () => {
  it('isImageMime true for image/* false otherwise', () => {
    expect(isImageMime('image/png')).toBe(true)
    expect(isImageMime('application/pdf')).toBe(false)
  })
  it('isPdfMime checks pdf', () => {
    expect(isPdfMime('application/pdf')).toBe(true)
    expect(isPdfMime('text/plain')).toBe(false)
  })
})

describe('limits', () => {
  it('exports sane defaults', () => {
    expect(ATTACHMENT_MAX_COUNT).toBe(10)
    expect(ATTACHMENT_MAX_BYTES).toBe(20 * 1024 * 1024)
  })
})
```

- [ ] **Step 2: 运行测试确认失败**

Run: `cd frontend && cnpm exec vitest run src/__tests__/attachmentsLib.test.ts`
Expected: 模块不存在错误。

- [ ] **Step 3: 实现 `frontend/src/lib/attachments.ts`**

```ts
import type { Attachment, AttachmentKind } from '@/types'

export const ATTACHMENT_MAX_COUNT = 10
export const ATTACHMENT_MAX_BYTES = 20 * 1024 * 1024

const EXT_MIME: Record<string, string> = {
  '.png': 'image/png',
  '.jpg': 'image/jpeg',
  '.jpeg': 'image/jpeg',
  '.webp': 'image/webp',
  '.gif': 'image/gif',
  '.pdf': 'application/pdf',
  '.txt': 'text/plain',
  '.log': 'text/plain',
  '.md': 'text/markdown',
  '.markdown': 'text/markdown',
  '.json': 'application/json',
  '.html': 'text/html',
  '.css': 'text/css',
  '.js': 'text/javascript',
  '.ts': 'application/typescript',
  '.tsx': 'application/typescript',
  '.go': 'text/x-go',
  '.py': 'text/x-python',
  '.rs': 'text/x-rust',
  '.yaml': 'application/yaml',
  '.yml': 'application/yaml',
  '.toml': 'application/toml',
}

function basename(path: string): string {
  // 支持 unix 与 windows 分隔符。
  const idx = Math.max(path.lastIndexOf('/'), path.lastIndexOf('\\'))
  return idx >= 0 ? path.slice(idx + 1) : path
}

function extOf(path: string): string {
  const name = basename(path)
  const dot = name.lastIndexOf('.')
  return dot >= 0 ? name.slice(dot).toLowerCase() : ''
}

export function inferAttachmentMeta(path: string): Attachment {
  const name = basename(path)
  const mime = EXT_MIME[extOf(path)] ?? 'application/octet-stream'
  const kind: AttachmentKind = mime.startsWith('image/') ? 'image' : 'file'
  return { path, name, mime, kind }
}

export function isImageMime(mime: string): boolean {
  return mime.startsWith('image/')
}

export function isPdfMime(mime: string): boolean {
  return mime === 'application/pdf'
}
```

- [ ] **Step 4: 运行测试**

Run: `cd frontend && cnpm exec vitest run src/__tests__/attachmentsLib.test.ts`
Expected: 全部 PASS。

- [ ] **Step 5: 提交**

```bash
git add frontend/src/lib/attachments.ts frontend/src/__tests__/attachmentsLib.test.ts
git commit -m "feat(frontend): add attachment metadata inference helpers"
```

---

### Task 12: i18n 增量 key

**Files:**

- Modify: `frontend/src/i18n/locales/zh-CN.ts`
- Modify: `frontend/src/i18n/locales/en.ts`

- [ ] **Step 1: 在 zh-CN 的 `input` 节点下追加，且新增 `chat.attachmentMissing`、`chat.attachmentRemove`**

修改 zh-CN.ts 的 `input` 对象，在末尾追加：

```ts
attach: '附加文件',
attachLimitCount: '最多附加 {{max}} 个文件',
attachLimitSize: '单个文件不超过 {{mb}}MB',
```

并在 `chat` 对象末尾追加：

```ts
attachmentMissing: '附件已丢失',
attachmentRemove: '移除附件',
```

- [ ] **Step 2: 对应改 en.ts**

`input` 对象末尾：

```ts
attach: 'Attach file',
attachLimitCount: 'Up to {{max}} files per message',
attachLimitSize: 'File must be at most {{mb}}MB',
```

`chat` 对象末尾：

```ts
attachmentMissing: 'Attachment missing',
attachmentRemove: 'Remove attachment',
```

- [ ] **Step 3: typecheck**

Run: `cd frontend && cnpm exec tsc --noEmit`
Expected: 无错误。

- [ ] **Step 4: 提交**

```bash
git add frontend/src/i18n/locales/zh-CN.ts frontend/src/i18n/locales/en.ts
git commit -m "feat(i18n): add attachment-related strings"
```

---

## Phase 4 — 前端组件

### Task 13: `AttachmentChips` 组件

**Files:**

- Create: `frontend/src/components/chat/AttachmentChips.tsx`
- Create: `frontend/src/__tests__/attachmentChips.test.tsx`

- [ ] **Step 1: 写失败测试**

`frontend/src/__tests__/attachmentChips.test.tsx`：

```tsx
import { render, screen } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { describe, expect, it, vi } from 'vitest'
import { I18nextProvider } from 'react-i18next'
import i18n from '@/i18n'
import { AttachmentChips } from '@/components/chat/AttachmentChips'
import type { Attachment } from '@/types'

const items: Attachment[] = [
  { path: '/a/foo.png', name: 'foo.png', mime: 'image/png', kind: 'image' },
  { path: '/a/spec.pdf', name: 'spec.pdf', mime: 'application/pdf', kind: 'file' },
]

function renderWithI18n(node: React.ReactNode) {
  return render(<I18nextProvider i18n={i18n}>{node}</I18nextProvider>)
}

describe('AttachmentChips', () => {
  it('renders one chip per attachment', () => {
    renderWithI18n(<AttachmentChips items={items} variant="message" />)
    expect(screen.getByText('foo.png')).toBeInTheDocument()
    expect(screen.getByText('spec.pdf')).toBeInTheDocument()
  })

  it('invokes onRemove when × clicked in input variant', async () => {
    const onRemove = vi.fn()
    const user = userEvent.setup()
    renderWithI18n(
      <AttachmentChips items={items} variant="input" onRemove={onRemove} />
    )
    const removes = screen.getAllByLabelText(/remove attachment|移除附件/i)
    expect(removes).toHaveLength(2)
    await user.click(removes[1])
    expect(onRemove).toHaveBeenCalledWith(1)
  })

  it('does not render remove buttons in message variant', () => {
    renderWithI18n(<AttachmentChips items={items} variant="message" />)
    expect(screen.queryByLabelText(/remove attachment|移除附件/i)).toBeNull()
  })

  it('renders nothing when items empty', () => {
    const { container } = renderWithI18n(
      <AttachmentChips items={[]} variant="input" />
    )
    expect(container.firstChild).toBeNull()
  })
})
```

- [ ] **Step 2: 运行确认失败**

Run: `cd frontend && cnpm exec vitest run src/__tests__/attachmentChips.test.tsx`
Expected: 模块不存在错误。

- [ ] **Step 3: 实现组件**

`frontend/src/components/chat/AttachmentChips.tsx`：

```tsx
import { useTranslation } from 'react-i18next'
import { FileText, Image as ImageIcon, Paperclip, X } from 'lucide-react'
import { cn } from '@/lib/utils'
import { isImageMime, isPdfMime } from '@/lib/attachments'
import type { Attachment } from '@/types'

interface Props {
  items: Attachment[]
  variant: 'input' | 'message'
  onRemove?: (index: number) => void
}

function chipIcon(att: Attachment) {
  if (isImageMime(att.mime)) return <ImageIcon size={12} />
  if (isPdfMime(att.mime)) return <FileText size={12} />
  return <Paperclip size={12} />
}

export function AttachmentChips({ items, variant, onRemove }: Props) {
  const { t } = useTranslation()

  if (items.length === 0) return null

  return (
    <div className="flex flex-wrap gap-1.5 px-3 py-2">
      {items.map((att, index) => (
        <span
          key={`${att.path}-${index}`}
          className={cn(
            'inline-flex items-center gap-1 max-w-[18rem] rounded-md border border-border bg-muted/60 px-2 py-1 text-xs text-foreground',
            variant === 'message' && 'bg-background/60'
          )}
        >
          {chipIcon(att)}
          <span className="truncate" title={att.name}>{att.name}</span>
          {variant === 'input' && onRemove && (
            <button
              type="button"
              aria-label={t('chat.attachmentRemove')}
              onClick={() => onRemove(index)}
              className="ml-0.5 rounded text-muted-foreground hover:text-foreground hover:bg-accent p-0.5"
            >
              <X size={10} />
            </button>
          )}
        </span>
      ))}
    </div>
  )
}
```

- [ ] **Step 4: 运行测试**

Run: `cd frontend && cnpm exec vitest run src/__tests__/attachmentChips.test.tsx`
Expected: 全部 PASS。

- [ ] **Step 5: 提交**

```bash
git add frontend/src/components/chat/AttachmentChips.tsx frontend/src/__tests__/attachmentChips.test.tsx
git commit -m "feat(chat): add AttachmentChips component"
```

---

### Task 14: chatStore 透传 `attachments`，loadMessages 反序列化

**Files:**

- Modify: `frontend/src/store/chatStore.ts`

- [ ] **Step 1: 拓宽 `sendMessage` 签名**

把 `ChatStore` 接口里的 `sendMessage` 改为：

```ts
sendMessage: (params: {
  sessionId: number
  content: string
  baseUrl: string
  apiKey: string
  modelName: string
  enabledUserTools: string[]
  attachments: Attachment[]
}) => Promise<void>
```

并在文件顶部 import：

```ts
import type {
  Attachment,
  ConfirmRequestEvent,
  Conversation,
  ConversationStatus,
  Message,
  SessionStatusEvent,
  StreamChunkEvent,
  StreamDoneEvent,
  StreamErrorEvent,
} from '../types'
```

- [ ] **Step 2: 改 `loadMessages` 把 `attachments` 反序列化**

在 `loadMessages` 中 `messages` 数组构造处替换为：

```ts
const messages: Message[] = (result.messages ?? []).map((message) => {
  let attachments: Attachment[] | undefined
  const raw = (message as { attachments?: string }).attachments
  if (raw && raw.trim() !== '') {
    try {
      const parsed = JSON.parse(raw) as Attachment[]
      if (Array.isArray(parsed) && parsed.length > 0) attachments = parsed
    } catch {
      // ignore malformed payload
    }
  }
  return {
    id: message.id,
    sessionId: message.session_id,
    parentId: message.parent_id ?? null,
    role: message.role as Message['role'],
    contentType: message.content_type as Message['contentType'],
    content: message.content,
    modelName: message.model_name,
    agentName: message.agent_name,
    tokensIn: message.tokens_in,
    tokensOut: message.tokens_out,
    extra: message.extra,
    attachments,
    createdAt: message.created,
  }
})
```

- [ ] **Step 3: 改 `sendMessage` 实现以透传**

把 `sendMessage` 函数体里乐观写入的 user message 加上 `attachments: params.attachments.length > 0 ? params.attachments : undefined`，并在最后 `AgentBinding.SendMessage(...)` 调用增加 `attachments` 字段，整体如下（仅展示需要替换的 user message 字面量与 AgentBinding 调用，其余不动）：

```ts
{
  id: tempUserId,
  sessionId: params.sessionId,
  parentId: null,
  role: 'user',
  contentType: 'text',
  content: params.content,
  modelName: '',
  agentName: '',
  tokensIn: 0,
  tokensOut: 0,
  extra: '',
  attachments: params.attachments.length > 0 ? params.attachments : undefined,
  createdAt: now,
},
```

```ts
await AgentBinding.SendMessage({
  session_id: params.sessionId,
  content: params.content,
  base_url: params.baseUrl,
  api_key: params.apiKey,
  model_name: params.modelName,
  enabled_user_tools: params.enabledUserTools,
  attachments: params.attachments.map(att => ({
    path: att.path,
    name: att.name,
    mime: att.mime,
  })),
})
```

- [ ] **Step 4: typecheck**

Run: `cd frontend && cnpm exec tsc --noEmit`
Expected: 无错误。

- [ ] **Step 5: 已有 chatStore 单测回归**

Run: `cd frontend && cnpm exec vitest run src/__tests__/chatStore.test.ts`
Expected: 全 PASS（若已有用例调用 `sendMessage` 但未传 `attachments`，会因 typecheck 失败而被发现；按 `attachments: []` 修复后 PASS）。

- [ ] **Step 6: 提交**

```bash
git add frontend/src/store/chatStore.ts frontend/src/__tests__/chatStore.test.ts
git commit -m "feat(store): pass attachments through sendMessage and loadMessages"
```

---

### Task 15: `ChatInput` 接 `Paperclip` + chips 渲染

**Files:**

- Modify: `frontend/src/components/chat/ChatInput.tsx`
- Create: `frontend/src/__tests__/chatInputAttachments.test.tsx`

- [ ] **Step 1: 写失败测试**

`frontend/src/__tests__/chatInputAttachments.test.tsx`：

```tsx
import { render, screen, waitFor } from '@testing-library/react'
import userEvent from '@testing-library/user-event'
import { beforeEach, describe, expect, it, vi } from 'vitest'
import { I18nextProvider } from 'react-i18next'
import i18n from '@/i18n'
import { ChatInput } from '@/components/chat/ChatInput'
import { useChatStore } from '@/store/chatStore'

vi.mock('@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/file', () => ({
  File: {
    SelectFile: vi.fn(),
  },
}))

vi.mock('@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/provider', () => ({
  Provider: {
    ProviderAndModelList: vi.fn().mockResolvedValue({
      provider_models: [
        {
          provider: {
            id: 1, provider_name: 'P', enabled: true, is_default: true,
            base_url: 'http://x', api_key: 'k',
          },
          models: [{ id: 1, model: 'm', alias: 'm', is_default: true, enable: true }],
        },
      ],
    }),
  },
}))

import { File as FileBinding } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/file'

beforeEach(() => {
  useChatStore.setState({
    conversations: [{
      id: 1, title: 't', createdAt: '', updatedAt: '',
      starred: false, status: 'idle',
    }],
    currentConversationId: 1,
  })
})

function renderInput() {
  return render(<I18nextProvider i18n={i18n}><ChatInput /></I18nextProvider>)
}

describe('ChatInput attachments', () => {
  it('adds a chip when SelectFile returns a path', async () => {
    ;(FileBinding.SelectFile as any).mockResolvedValueOnce({ file_path: '/x/foo.png' })
    const user = userEvent.setup()
    renderInput()

    await user.click(screen.getByLabelText(/attach file|附加文件/i))
    await waitFor(() => expect(screen.getByText('foo.png')).toBeInTheDocument())
  })

  it('does nothing when SelectFile returns empty', async () => {
    ;(FileBinding.SelectFile as any).mockResolvedValueOnce(null)
    const user = userEvent.setup()
    renderInput()

    await user.click(screen.getByLabelText(/attach file|附加文件/i))
    // 无 chip 出现，则 foo.png 不存在
    expect(screen.queryByText('foo.png')).toBeNull()
  })

  it('removes a chip via × button', async () => {
    ;(FileBinding.SelectFile as any).mockResolvedValueOnce({ file_path: '/x/foo.png' })
    const user = userEvent.setup()
    renderInput()

    await user.click(screen.getByLabelText(/attach file|附加文件/i))
    await screen.findByText('foo.png')

    await user.click(screen.getByLabelText(/remove attachment|移除附件/i))
    expect(screen.queryByText('foo.png')).toBeNull()
  })
})
```

- [ ] **Step 2: 运行确认失败**

Run: `cd frontend && cnpm exec vitest run src/__tests__/chatInputAttachments.test.tsx`
Expected: 失败（Paperclip 无 onClick，无 chip 渲染）。

- [ ] **Step 3: 改 `ChatInput.tsx`**

在文件顶部 import 区追加：

```ts
import { File as FileBinding } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/file'
import { AttachmentChips } from './AttachmentChips'
import { ATTACHMENT_MAX_COUNT, inferAttachmentMeta } from '@/lib/attachments'
import type { Attachment } from '@/types'
```

在组件 state 区追加：

```ts
const [attachments, setAttachments] = useState<Attachment[]>([])
```

`useEffect(() => { ... clearContent ... }, [currentConversationId, editor])` 内也清空附件：

```ts
useEffect(() => {
  if (!editor) return
  editor.commands.clearContent()
  setIsEmpty(true)
  setAttachments([])
}, [currentConversationId, editor])
```

把 `handleSend` 的 `isEmpty` 判定与发送参数改为：

```ts
const handleSend = async () => {
  if (!editor || isStreaming) return
  const content = ((editor.storage as { markdown?: { getMarkdown?: () => string } }).markdown?.getMarkdown?.() ?? '').trim()
  if (!content && attachments.length === 0) return

  const selectedProvider = providerModels.find(pm => pm.models.some(m => m.id === selectedModelId))
  const selectedModel = selectedProvider?.models.find(m => m.id === selectedModelId)
  if (!selectedProvider || !selectedModel) return

  let sessionId = currentConversationId
  if (!sessionId) {
    sessionId = await createConversation()
    if (!sessionId) return
    setCurrentConversation(sessionId)
  }

  const outboundAttachments = attachments
  editor.commands.clearContent()
  setIsEmpty(true)
  setAttachments([])

  await sendMessage({
    sessionId,
    content,
    baseUrl: selectedProvider.provider.base_url,
    apiKey: selectedProvider.provider.api_key,
    modelName: selectedModel.model,
    enabledUserTools: tools.filter(tool => tool.enabled).map(tool => tool.id),
    attachments: outboundAttachments,
  })
}
```

新增 `handleAttach`：

```ts
const handleAttach = async () => {
  if (attachments.length >= ATTACHMENT_MAX_COUNT) return
  const result = await FileBinding.SelectFile({})
  if (!result?.file_path) return
  setAttachments(prev => [...prev, inferAttachmentMeta(result.file_path)])
}
```

把 Paperclip 按钮加上 onClick 与 disabled：

```tsx
<button
  aria-label={t('input.attach')}
  onClick={handleAttach}
  disabled={attachments.length >= ATTACHMENT_MAX_COUNT}
  className="p-1.5 rounded-lg hover:bg-accent text-muted-foreground hover:text-foreground transition-colors disabled:opacity-40"
>
  <Paperclip size={16} />
</button>
```

在 `<EditorContent editor={editor} />` 上方插入 chips：

```tsx
<AttachmentChips
  items={attachments}
  variant="input"
  onRemove={(idx) => setAttachments(prev => prev.filter((_, i) => i !== idx))}
/>
<EditorContent editor={editor} />
```

把 Send 按钮的 disabled 改为：

```tsx
disabled={!isStreaming && isEmpty && attachments.length === 0}
```

并把 `className` 里的 `!isEmpty` 也调整为 `!isEmpty || attachments.length > 0`：

```tsx
className={cn(
  'p-1.5 rounded-lg transition-colors',
  isStreaming
    ? 'bg-destructive/10 text-destructive hover:bg-destructive/20'
    : (!isEmpty || attachments.length > 0)
    ? 'bg-primary text-primary-foreground hover:bg-primary/90'
    : 'text-muted-foreground cursor-not-allowed'
)}
```

- [ ] **Step 4: 运行新测试**

Run: `cd frontend && cnpm exec vitest run src/__tests__/chatInputAttachments.test.tsx`
Expected: 全部 PASS。

- [ ] **Step 5: typecheck**

Run: `cd frontend && cnpm exec tsc --noEmit`
Expected: 无错误。

- [ ] **Step 6: 提交**

```bash
git add frontend/src/components/chat/ChatInput.tsx frontend/src/__tests__/chatInputAttachments.test.tsx
git commit -m "feat(chat): wire Paperclip to SelectFile and render attachment chips"
```

---

### Task 16: `MessageItem` 在 user 气泡内显示只读 chips

**Files:**

- Modify: `frontend/src/components/chat/MessageItem.tsx`

- [ ] **Step 1: 改 user/text 分支**

把 `MessageItem` 中 `if (role === 'user' && contentType === 'text')` 那段替换为：

```tsx
if (role === 'user' && contentType === 'text') {
  const attachments = message.attachments ?? []
  return (
    <div className="flex justify-end">
      <div className="max-w-[70%] select-text rounded-2xl rounded-br-sm bg-primary text-primary-foreground px-4 py-2.5 text-sm dark:bg-muted dark:text-foreground">
        {content && <div className="whitespace-pre-wrap break-words">{content}</div>}
        {attachments.length > 0 && (
          <AttachmentChips items={attachments} variant="message" />
        )}
      </div>
    </div>
  )
}
```

在文件顶部 import 区追加：

```ts
import { AttachmentChips } from './AttachmentChips'
```

- [ ] **Step 2: 已有 messageItem 测试回归**

Run: `cd frontend && cnpm exec vitest run src/__tests__/messageItem.test.tsx`
Expected: 全 PASS（旧用例不传 `attachments`，分支保持空数组，UI 仅渲染文字）。

- [ ] **Step 3: typecheck**

Run: `cd frontend && cnpm exec tsc --noEmit`
Expected: 无错误。

- [ ] **Step 4: 提交**

```bash
git add frontend/src/components/chat/MessageItem.tsx
git commit -m "feat(chat): show attachment chips in user message bubble"
```

---

## Phase 5 — 验收

### Task 17: 整体测试套件

- [ ] **Step 1: 跑后端全部测试**

Run: `go test ./backend/... -count=1`
Expected: 全 PASS。

- [ ] **Step 2: 跑前端全部测试**

Run: `cd frontend && cnpm exec vitest run`
Expected: 全 PASS。

- [ ] **Step 3: typecheck**

Run: `cd frontend && cnpm exec tsc --noEmit`
Expected: 无错误。

- [ ] **Step 4: 若失败，回到对应 Task 修复后再跑**

---

### Task 18: 手动验收

启动 `task dev` 后逐项验证（无需 commit）：

- [ ] **a. 浅色 / 深色主题** × **极小 / 小 / 标准 / 大 / 超大** 字号下，ChatInput 中加入 1 个 chip，确认 chip 颜色、间距、文字与气泡风格一致，不溢出。
- [ ] **b. 中英文** 切换语言：Paperclip aria-label、移除按钮 aria-label、可能的错误提示文案均有翻译，无 fallback key。
- [ ] **c. 单图片消息：** 选 PNG → 发送 → 助手能回复；刷新会话后 user 气泡的 chip 仍存在。
- [ ] **d. 单 PDF 消息：** 同上，chip 显示 FileText 图标。
- [ ] **e. 混合多附件：** 一次性选 3 个不同类型（PNG/PDF/.go），都正常进入 chip 列表与历史回放。
- [ ] **f. 多次添加：** 连续点 Paperclip 10 次（达到上限）后按钮 disabled。
- [ ] **g. 缺失文件降级：** 发送后在终端 `rm` 掉源文件，再次打开会话，模型再次发送时不崩溃，且日志/网络面板可以看到模型接收的 user message 含 `[Missing attachment: <name>]`。
- [ ] **h. 错误路径：** 用 `dd` 造一个 25MB 文件 → 选中并发送 → 前端在 `agent:stream:error` 或 `sendMessage` 返回错误处看到“exceeds”相关消息（最低限度：流式状态回到 idle，控制台无未捕获错误）。
- [ ] **i. 取消选择：** 文件对话框点取消 → 不新增 chip，无错误。
- [ ] 全部通过后，本任务到此完成；分支可以推送或合并。

---

## 自检（写完计划后人工核对一次）

**Spec 覆盖：**

- 仅存路径（方案 A） → Task 1/2/5。
- 新增 `Attachments` 字段 → Task 1。
- `BuildUserMessage` 统一构造 → Task 2/5。
- 缺失文件降级 → Task 2（单测）/5（重放）。
- 数量 ≤ 10，大小 ≤ 20MB 校验 → Task 7/8。
- 修复 `SelectFile` 三个 bug → Task 9。
- 前端 `Attachment` 类型与 Message 字段 → Task 10。
- 前端工具函数 → Task 11。
- i18n key → Task 12。
- AttachmentChips 组件 + 测试 → Task 13。
- chatStore 透传与反序列化 → Task 14。
- ChatInput 接 Paperclip → Task 15。
- MessageItem 渲染 chips → Task 16。
- 主题 × 字号 × 中英文手测 → Task 18。

**Placeholder 扫描：**无 TBD / TODO。Task 6 中第 2 步显式标注了“先读 `manager.go` 再选其一”的决策点，给出了两条具体路径，不算占位。

**类型一致性：**

- `Attachment{Name, Path, Mime, Kind}`（Go） ↔ `Attachment{path, name, mime, kind}`（TS）：字段名/语义一致；后端 JSON tag 与前端 JSON 字段对齐（go 用驼峰转下划线？— 不，go json tag 已写为 `path/name/mime/kind`，前端类型同名小写匹配 ✅）。
- `AttachmentInput` 仅含 `path/name/mime`，后端 `NormalizeAttachment` 补 `kind`，前端发送时也只传 `path/name/mime`（Task 14 step 3 一致 ✅）。
- `SendMessageInput.Attachments` (`json:"attachments"`) ↔ 前端 `attachments` ✅。
- `MessageItem.Attachments string` (`json:"attachments"`) ↔ 前端 `message.attachments`（字符串 JSON，前端 parse） ✅。
- `BuildUserMessage(content, atts)` 在 Task 2 定义、Task 5 调用、Task 6 测试中一致。
- `ATTACHMENT_MAX_COUNT = 10`、`ATTACHMENT_MAX_BYTES = 20MB` 在 Task 7（Go）与 Task 11（TS）值一致。

---
