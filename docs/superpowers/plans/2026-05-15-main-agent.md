# Main Agent Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement the core AI conversation engine with streaming chat, session management, tool calling, and prompt loading based on trpc-agent-go framework.

**Architecture:** trpc-agent-go Runner + LLMAgent handles the conversation loop. Our storage layer persists session metadata and messages for the frontend. Wails Events bridge streaming from Go channels to React. A generic prompt loader and extensible tool registry support future multi-agent expansion.

**Tech Stack:** Go 1.25, Wails v3, trpc-agent-go, GORM/SQLite, React 18, TypeScript, Zustand

---

## File Structure

### New Files

| Path | Responsibility |
|------|---------------|
| `backend/models/data_models/session.go` | Session ORM model |
| `backend/models/data_models/message.go` | Message ORM model |
| `backend/storage/session.go` | Session CRUD + pagination |
| `backend/storage/message.go` | Message write + paginated read |
| `backend/pkg/prompt/prompt.go` | Generic prompt register + load |
| `backend/pkg/prompt/prompt_test.go` | Prompt loading tests |
| `backend/pkg/agent/tools/registry.go` | Tool registry with meta |
| `backend/pkg/agent/tools/registry_test.go` | Registry tests |
| `backend/pkg/agent/tools/datetime.go` | Built-in: datetime tool |
| `backend/pkg/agent/tools/file_rw.go` | Built-in: file read/write |
| `backend/pkg/agent/tools/shell.go` | Built-in: shell command |
| `backend/pkg/agent/tools/web_search.go` | User tool: web search |
| `backend/pkg/agent/tools/code_exec.go` | User tool: code execution |
| `backend/pkg/agent/manager.go` | Agent lifecycle, runner factory |
| `backend/pkg/agent/chat_handler.go` | Chat processing, event consumption |
| `backend/pkg/agent/stream_manager.go` | Active stream tracking, cancellation |
| `backend/service/agent/agent.go` | Wails service public methods |
| `backend/service/agent/agent_implement.go` | ServiceStartup |
| `backend/service/agent/agent_internal.go` | Private helpers |
| `backend/service/agent/agent_dto/send_message.go` | SendMessage DTO |
| `backend/service/agent/agent_dto/stop_generation.go` | StopGeneration DTO |
| `backend/service/agent/agent_dto/respond_to_confirm.go` | RespondToConfirm DTO |
| `backend/service/agent/agent_dto/create_session.go` | CreateSession DTO |
| `backend/service/agent/agent_dto/list_sessions.go` | ListSessions DTO |
| `backend/service/agent/agent_dto/load_session_messages.go` | LoadSessionMessages DTO |
| `backend/service/agent/agent_dto/rename_session.go` | RenameSession DTO |
| `backend/service/agent/agent_dto/delete_session.go` | DeleteSession DTO |
| `backend/service/agent/agent_dto/toggle_star_session.go` | ToggleStarSession DTO |
| `backend/service/agent/agent_dto/generate_title.go` | GenerateTitle DTO |
| `frontend/src/components/chat/ToolCallBlock.tsx` | Tool call display component |
| `frontend/src/components/chat/ToolConfirmCard.tsx` | Sensitive tool confirm card |

### Modified Files

| Path | Change |
|------|--------|
| `backend/storage/storage.go` | Add Session + Message to AutoMigrate |
| `backend/pkg/id/event_id/id.go` | Add agent stream event IDs |
| `main.go` | Register agent service |
| `frontend/src/types/index.ts` | Update Message type, add tool types |
| `frontend/src/store/chatStore.ts` | Replace mock with backend calls + Wails events |
| `frontend/src/components/chat/ChatInput.tsx` | Wire real send/stop |
| `frontend/src/components/chat/ChatMessages.tsx` | Render tool blocks, confirm cards |
| `frontend/src/components/chat/MessageItem.tsx` | Support new content types |
| `frontend/src/components/sidebar/ConversationList.tsx` | Backend pagination |
| `frontend/src/components/sidebar/ConversationItem.tsx` | Wire backend rename/delete/star |
| `frontend/src/i18n/locales/zh-CN.ts` | Add tool/confirm i18n keys |
| `frontend/src/i18n/locales/en.ts` | Add tool/confirm i18n keys |
| `frontend/src/mock/data.ts` | Remove or keep as fallback |

---

### Task 1: Install trpc-agent-go dependency

**Files:**
- Modify: `go.mod`
- Modify: `go.sum`

- [ ] **Step 1: Add trpc-agent-go module**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop && go get trpc.group/trpc-go/trpc-agent-go@latest
```

- [ ] **Step 2: Verify dependency resolves**

```bash
go mod tidy
```

Expected: exit 0, no errors.

- [ ] **Step 3: Commit**

```bash
git add go.mod go.sum
git commit -m "chore: add trpc-agent-go dependency"
```

---

### Task 2: Data models — Session and Message

**Files:**
- Create: `backend/models/data_models/session.go`
- Create: `backend/models/data_models/message.go`
- Modify: `backend/storage/storage.go`

- [ ] **Step 1: Create Session model**

```go
// backend/models/data_models/session.go
package data_models

// Session stores conversation metadata for frontend display.
type Session struct {
	OrmModel
	Title   string `gorm:"type:varchar(255)" json:"title"`
	Starred bool   `gorm:"type:bool;default:0;index" json:"starred"`
	Status  string `gorm:"type:varchar(50);default:'idle'" json:"status"`
}
```

- [ ] **Step 2: Create Message model**

```go
// backend/models/data_models/message.go
package data_models

// Message stores a single message record for frontend display and history.
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
}
```

- [ ] **Step 3: Add models to AutoMigrate in storage.go**

In `backend/storage/storage.go`, add `&data_models.Session{}` and `&data_models.Message{}` to the `AutoMigrate` call:

```go
func NewStorageFromDB(db *gorm.DB) (*Storage, error) {
	if err := db.AutoMigrate(
		&data_models.Provider{},
		&data_models.ProviderDefaultModel{},
		&data_models.Model{},
		&data_models.Session{},
		&data_models.Message{},
	); err != nil {
		return nil, err
	}
	return &Storage{sqliteDB: db}, nil
}
```

- [ ] **Step 4: Verify compilation**

```bash
go build ./...
```

Expected: exit 0.

- [ ] **Step 5: Commit**

```bash
git add backend/models/data_models/session.go backend/models/data_models/message.go backend/storage/storage.go
git commit -m "feat(models): add Session and Message data models"
```

---

### Task 3: Storage layer — session.go and message.go

**Files:**
- Create: `backend/storage/session.go`
- Create: `backend/storage/message.go`

- [ ] **Step 1: Create session storage**

```go
// backend/storage/session.go
package storage

import (
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gorm.io/gorm"
)

// CreateSession inserts a new session record and returns it with the auto-generated ID.
func (s *Storage) CreateSession(session data_models.Session) (*data_models.Session, error) {
	if err := s.sqliteDB.Create(&session).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

// GetSession retrieves a single session by ID.
func (s *Storage) GetSession(id uint) (*data_models.Session, error) {
	var session data_models.Session
	if err := s.sqliteDB.First(&session, id).Error; err != nil {
		return nil, err
	}
	return &session, nil
}

// ListSessions returns sessions ordered by updated_at desc with cursor-based pagination.
// When cursor is 0 it starts from the latest. Returns up to limit rows.
func (s *Storage) ListSessions(cursor uint, limit int, starredOnly bool) ([]data_models.Session, error) {
	var sessions []data_models.Session
	q := s.sqliteDB.Order("updated_at DESC")
	if cursor > 0 {
		q = q.Where("id < ?", cursor)
	}
	if starredOnly {
		q = q.Where("starred = ?", true)
	}
	if err := q.Limit(limit).Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}

// UpdateSessionTitle updates the title for the given session.
func (s *Storage) UpdateSessionTitle(id uint, title string) error {
	return s.sqliteDB.Model(&data_models.Session{}).Where("id = ?", id).Update("title", title).Error
}

// UpdateSessionStatus updates the status for the given session.
func (s *Storage) UpdateSessionStatus(id uint, status string) error {
	return s.sqliteDB.Model(&data_models.Session{}).Where("id = ?", id).Update("status", status).Error
}

// UpdateSessionStarred toggles the starred flag for the given session.
func (s *Storage) UpdateSessionStarred(id uint, starred bool) error {
	return s.sqliteDB.Model(&data_models.Session{}).Where("id = ?", id).Update("starred", starred).Error
}

// DeleteSession soft-deletes a session and all its messages in a transaction.
func (s *Storage) DeleteSession(id uint) error {
	return s.sqliteDB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("session_id = ?", id).Delete(&data_models.Message{}).Error; err != nil {
			return err
		}
		return tx.Delete(&data_models.Session{}, id).Error
	})
}

// TouchSession bumps updated_at to now for the given session.
func (s *Storage) TouchSession(id uint) error {
	return s.sqliteDB.Model(&data_models.Session{}).Where("id = ?", id).Update("updated_at", gorm.Expr("CURRENT_TIMESTAMP")).Error
}
```

- [ ] **Step 2: Create message storage**

```go
// backend/storage/message.go
package storage

import "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"

// CreateMessage inserts a single message record.
func (s *Storage) CreateMessage(msg data_models.Message) (*data_models.Message, error) {
	if err := s.sqliteDB.Create(&msg).Error; err != nil {
		return nil, err
	}
	return &msg, nil
}

// CreateMessages bulk-inserts message records.
func (s *Storage) CreateMessages(msgs *[]data_models.Message) error {
	if len(*msgs) == 0 {
		return nil
	}
	return s.sqliteDB.Create(msgs).Error
}

// ListMessagesForSession returns messages for a session ordered by created_at asc with offset pagination.
func (s *Storage) ListMessagesForSession(sessionID uint, offset int, limit int) ([]data_models.Message, error) {
	var msgs []data_models.Message
	if err := s.sqliteDB.Where("session_id = ?", sessionID).
		Order("created_at ASC").
		Offset(offset).Limit(limit).
		Find(&msgs).Error; err != nil {
		return nil, err
	}
	return msgs, nil
}

// CountMessagesForSession returns the total message count for a session.
func (s *Storage) CountMessagesForSession(sessionID uint) (int64, error) {
	var count int64
	if err := s.sqliteDB.Model(&data_models.Message{}).Where("session_id = ?", sessionID).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// DeleteMessagesForSession soft-deletes all messages belonging to a session.
func (s *Storage) DeleteMessagesForSession(sessionID uint) error {
	return s.sqliteDB.Where("session_id = ?", sessionID).Delete(&data_models.Message{}).Error
}
```

- [ ] **Step 3: Verify compilation**

```bash
go build ./...
```

Expected: exit 0.

- [ ] **Step 4: Commit**

```bash
git add backend/storage/session.go backend/storage/message.go
git commit -m "feat(storage): add session and message storage operations"
```

---

### Task 4: Prompt loading mechanism

**Files:**
- Create: `backend/pkg/prompt/prompt.go`
- Create: `backend/pkg/prompt/prompt_test.go`

- [ ] **Step 1: Write the failing test**

```go
// backend/pkg/prompt/prompt_test.go
package prompt

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad_ReturnsDefault_WhenNoCustomFile(t *testing.T) {
	registry = make(map[string]string)
	Register("test_agent", "You are a helpful assistant.")

	t.Setenv("LEMONTEA_DATA_DIR", t.TempDir())

	got, err := Load("test_agent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "You are a helpful assistant." {
		t.Fatalf("expected default prompt, got: %s", got)
	}
}

func TestLoad_ReturnsCustom_WhenFileExists(t *testing.T) {
	registry = make(map[string]string)
	Register("test_agent", "default prompt")

	tmpDir := t.TempDir()
	t.Setenv("LEMONTEA_DATA_DIR", tmpDir)

	promptDir := filepath.Join(tmpDir, "prompt", "test_agent")
	if err := os.MkdirAll(promptDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(promptDir, "index.md"), []byte("custom prompt content"), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := Load("test_agent")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "custom prompt content" {
		t.Fatalf("expected custom prompt, got: %s", got)
	}
}

func TestLoad_ErrorsOnUnregistered(t *testing.T) {
	registry = make(map[string]string)

	t.Setenv("LEMONTEA_DATA_DIR", t.TempDir())

	_, err := Load("nonexistent")
	if err == nil {
		t.Fatal("expected error for unregistered prompt")
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop && go test ./backend/pkg/prompt/ -v
```

Expected: compilation failure (package doesn't exist yet).

- [ ] **Step 3: Implement prompt.go**

```go
// backend/pkg/prompt/prompt.go
package prompt

import (
	"fmt"
	"os"
	"path/filepath"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/dir"
)

var registry = make(map[string]string)

// Register associates a promptId with its built-in default content.
func Register(promptId string, defaultContent string) {
	registry[promptId] = defaultContent
}

// Load returns the prompt for promptId. It checks {dataDir}/prompt/{promptId}/index.md first;
// if the file does not exist, it falls back to the registered default.
func Load(promptId string) (string, error) {
	defaultContent, ok := registry[promptId]
	if !ok {
		return "", fmt.Errorf("prompt %q is not registered", promptId)
	}

	dataDir, err := dir.GetDataDir()
	if err != nil {
		return defaultContent, nil
	}

	customPath := filepath.Join(dataDir, "prompt", promptId, "index.md")
	content, err := os.ReadFile(customPath)
	if err != nil {
		return defaultContent, nil
	}

	return string(content), nil
}
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
go test ./backend/pkg/prompt/ -v
```

Expected: all 3 tests PASS.

- [ ] **Step 5: Commit**

```bash
git add backend/pkg/prompt/
git commit -m "feat(prompt): add generic prompt loading mechanism"
```

---

### Task 5: Tool registry

**Files:**
- Create: `backend/pkg/agent/tools/registry.go`
- Create: `backend/pkg/agent/tools/registry_test.go`

- [ ] **Step 1: Write the failing test**

```go
// backend/pkg/agent/tools/registry_test.go
package tools

import (
	"encoding/json"
	"testing"
)

type mockTool struct {
	name string
}

func (m *mockTool) Name() string { return m.name }

func TestRegistry_RegisterAndGet(t *testing.T) {
	r := NewRegistry()
	meta := ToolMeta{
		Name:            "test_tool",
		Description:     "A test tool",
		Category:        CategoryBuiltin,
		RequiresConfirm: false,
		FormatPurpose:   func(args json.RawMessage) string { return "testing" },
	}
	r.Register(meta)

	got, ok := r.Get("test_tool")
	if !ok {
		t.Fatal("expected tool to be found")
	}
	if got.Name != "test_tool" {
		t.Fatalf("expected name test_tool, got %s", got.Name)
	}
}

func TestRegistry_BuiltinAndUserTools(t *testing.T) {
	r := NewRegistry()
	r.Register(ToolMeta{Name: "builtin1", Category: CategoryBuiltin})
	r.Register(ToolMeta{Name: "user1", Category: CategoryUser})
	r.Register(ToolMeta{Name: "user2", Category: CategoryUser})

	builtins := r.BuiltinTools()
	if len(builtins) != 1 {
		t.Fatalf("expected 1 builtin, got %d", len(builtins))
	}

	userTools := r.UserTools()
	if len(userTools) != 2 {
		t.Fatalf("expected 2 user tools, got %d", len(userTools))
	}
}

func TestRegistry_EnabledTools(t *testing.T) {
	r := NewRegistry()
	r.Register(ToolMeta{Name: "builtin1", Category: CategoryBuiltin})
	r.Register(ToolMeta{Name: "user1", Category: CategoryUser})
	r.Register(ToolMeta{Name: "user2", Category: CategoryUser})

	enabled := r.EnabledTools([]string{"user1"})
	if len(enabled) != 2 {
		t.Fatalf("expected 2 (1 builtin + 1 user enabled), got %d", len(enabled))
	}

	names := make(map[string]bool)
	for _, m := range enabled {
		names[m.Name] = true
	}
	if !names["builtin1"] || !names["user1"] {
		t.Fatalf("expected builtin1 and user1, got %v", names)
	}
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./backend/pkg/agent/tools/ -v
```

Expected: compilation failure.

- [ ] **Step 3: Implement registry.go**

```go
// backend/pkg/agent/tools/registry.go
package tools

import (
	"encoding/json"
	"sync"
)

const (
	CategoryBuiltin = "builtin"
	CategoryUser    = "user"
)

// ToolMeta describes a tool's metadata including confirmation requirements and purpose formatting.
type ToolMeta struct {
	Name            string
	Description     string
	Category        string
	RequiresConfirm bool
	FormatPurpose   func(args json.RawMessage) string
}

// Registry manages tool registration and lookup.
type Registry struct {
	mu    sync.RWMutex
	tools map[string]ToolMeta
	order []string
}

// NewRegistry creates an empty tool registry.
func NewRegistry() *Registry {
	return &Registry{
		tools: make(map[string]ToolMeta),
	}
}

// Register adds a tool to the registry.
func (r *Registry) Register(meta ToolMeta) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.tools[meta.Name]; !exists {
		r.order = append(r.order, meta.Name)
	}
	r.tools[meta.Name] = meta
}

// Get retrieves a tool's metadata by name.
func (r *Registry) Get(name string) (ToolMeta, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	meta, ok := r.tools[name]
	return meta, ok
}

// BuiltinTools returns all tools in the builtin category.
func (r *Registry) BuiltinTools() []ToolMeta {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []ToolMeta
	for _, name := range r.order {
		if r.tools[name].Category == CategoryBuiltin {
			result = append(result, r.tools[name])
		}
	}
	return result
}

// UserTools returns all tools in the user category.
func (r *Registry) UserTools() []ToolMeta {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []ToolMeta
	for _, name := range r.order {
		if r.tools[name].Category == CategoryUser {
			result = append(result, r.tools[name])
		}
	}
	return result
}

// EnabledTools returns all builtin tools plus the user tools whose names appear in the enabled list.
func (r *Registry) EnabledTools(enabled []string) []ToolMeta {
	enabledSet := make(map[string]bool, len(enabled))
	for _, name := range enabled {
		enabledSet[name] = true
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []ToolMeta
	for _, name := range r.order {
		meta := r.tools[name]
		if meta.Category == CategoryBuiltin || enabledSet[name] {
			result = append(result, meta)
		}
	}
	return result
}
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
go test ./backend/pkg/agent/tools/ -v
```

Expected: all 3 tests PASS.

- [ ] **Step 5: Commit**

```bash
git add backend/pkg/agent/tools/
git commit -m "feat(tools): add extensible tool registry"
```

---

### Task 6: Built-in tools — datetime, file_rw, shell

**Files:**
- Create: `backend/pkg/agent/tools/datetime.go`
- Create: `backend/pkg/agent/tools/file_rw.go`
- Create: `backend/pkg/agent/tools/shell.go`

- [ ] **Step 1: Implement datetime tool**

```go
// backend/pkg/agent/tools/datetime.go
package tools

import (
	"context"
	"encoding/json"
	"time"

	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

type dateTimeInput struct {
	Format   string `json:"format" jsonschema:"description=Go time format string (default: RFC3339)"`
	Timezone string `json:"timezone" jsonschema:"description=IANA timezone name (default: Local)"`
}

type dateTimeOutput struct {
	DateTime string `json:"datetime"`
	Unix     int64  `json:"unix"`
	Timezone string `json:"timezone"`
}

// dateTimeFunc returns the current date and time.
func dateTimeFunc(ctx context.Context, input dateTimeInput) (dateTimeOutput, error) {
	loc := time.Local
	if input.Timezone != "" {
		parsed, err := time.LoadLocation(input.Timezone)
		if err != nil {
			return dateTimeOutput{}, err
		}
		loc = parsed
	}

	now := time.Now().In(loc)
	format := time.RFC3339
	if input.Format != "" {
		format = input.Format
	}

	return dateTimeOutput{
		DateTime: now.Format(format),
		Unix:     now.Unix(),
		Timezone: loc.String(),
	}, nil
}

// NewDateTimeTool creates the built-in datetime tool.
func NewDateTimeTool() *function.FunctionTool[dateTimeInput, dateTimeOutput] {
	return function.NewFunctionTool(
		dateTimeFunc,
		function.WithName("datetime"),
		function.WithDescription("Get the current date and time in a specified format and timezone"),
	)
}

// DateTimeMeta returns the ToolMeta for the datetime tool.
func DateTimeMeta() ToolMeta {
	return ToolMeta{
		Name:            "datetime",
		Description:     "Get the current date and time",
		Category:        CategoryBuiltin,
		RequiresConfirm: false,
		FormatPurpose: func(args json.RawMessage) string {
			return "Get current date and time"
		},
	}
}
```

- [ ] **Step 2: Implement file_rw tool**

```go
// backend/pkg/agent/tools/file_rw.go
package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

type fileReadInput struct {
	Path string `json:"path" jsonschema:"description=Absolute file path to read,required"`
}

type fileReadOutput struct {
	Content string `json:"content"`
	Size    int64  `json:"size"`
}

// fileReadFunc reads the content of a file.
func fileReadFunc(ctx context.Context, input fileReadInput) (fileReadOutput, error) {
	content, err := os.ReadFile(input.Path)
	if err != nil {
		return fileReadOutput{}, fmt.Errorf("read file: %w", err)
	}
	info, _ := os.Stat(input.Path)
	var size int64
	if info != nil {
		size = info.Size()
	}
	return fileReadOutput{Content: string(content), Size: size}, nil
}

// NewFileReadTool creates the built-in file read tool.
func NewFileReadTool() *function.FunctionTool[fileReadInput, fileReadOutput] {
	return function.NewFunctionTool(
		fileReadFunc,
		function.WithName("file_read"),
		function.WithDescription("Read the content of a file at the given path"),
	)
}

type fileWriteInput struct {
	Path    string `json:"path" jsonschema:"description=Absolute file path to write,required"`
	Content string `json:"content" jsonschema:"description=Content to write,required"`
}

type fileWriteOutput struct {
	BytesWritten int `json:"bytes_written"`
}

// fileWriteFunc writes content to a file, creating parent directories as needed.
func fileWriteFunc(ctx context.Context, input fileWriteInput) (fileWriteOutput, error) {
	if err := os.MkdirAll(filepath.Dir(input.Path), 0o755); err != nil {
		return fileWriteOutput{}, fmt.Errorf("create dir: %w", err)
	}
	if err := os.WriteFile(input.Path, []byte(input.Content), 0o644); err != nil {
		return fileWriteOutput{}, fmt.Errorf("write file: %w", err)
	}
	return fileWriteOutput{BytesWritten: len(input.Content)}, nil
}

// NewFileWriteTool creates the built-in file write tool.
func NewFileWriteTool() *function.FunctionTool[fileWriteInput, fileWriteOutput] {
	return function.NewFunctionTool(
		fileWriteFunc,
		function.WithName("file_write"),
		function.WithDescription("Write content to a file at the given path"),
	)
}

// FileReadMeta returns the ToolMeta for the file read tool.
func FileReadMeta() ToolMeta {
	return ToolMeta{
		Name:            "file_read",
		Description:     "Read file content",
		Category:        CategoryBuiltin,
		RequiresConfirm: true,
		FormatPurpose: func(args json.RawMessage) string {
			var input fileReadInput
			_ = json.Unmarshal(args, &input)
			return fmt.Sprintf("Read file: %s", input.Path)
		},
	}
}

// FileWriteMeta returns the ToolMeta for the file write tool.
func FileWriteMeta() ToolMeta {
	return ToolMeta{
		Name:            "file_write",
		Description:     "Write content to file",
		Category:        CategoryBuiltin,
		RequiresConfirm: true,
		FormatPurpose: func(args json.RawMessage) string {
			var input fileWriteInput
			_ = json.Unmarshal(args, &input)
			return fmt.Sprintf("Write to file: %s (%d bytes)", input.Path, len(input.Content))
		},
	}
}
```

- [ ] **Step 3: Implement shell tool**

```go
// backend/pkg/agent/tools/shell.go
package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

const shellTimeout = 30 * time.Second

type shellInput struct {
	Command string `json:"command" jsonschema:"description=Shell command to execute,required"`
	WorkDir string `json:"work_dir" jsonschema:"description=Working directory (default: user home)"`
}

type shellOutput struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exit_code"`
}

// shellFunc executes a shell command with a timeout.
func shellFunc(ctx context.Context, input shellInput) (shellOutput, error) {
	ctx, cancel := context.WithTimeout(ctx, shellTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, "sh", "-c", input.Command)
	if input.WorkDir != "" {
		cmd.Dir = input.WorkDir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			return shellOutput{}, fmt.Errorf("exec: %w", err)
		}
	}

	return shellOutput{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
	}, nil
}

// NewShellTool creates the built-in shell command tool.
func NewShellTool() *function.FunctionTool[shellInput, shellOutput] {
	return function.NewFunctionTool(
		shellFunc,
		function.WithName("shell"),
		function.WithDescription("Execute a shell command and return stdout, stderr, and exit code"),
	)
}

// ShellMeta returns the ToolMeta for the shell tool.
func ShellMeta() ToolMeta {
	return ToolMeta{
		Name:            "shell",
		Description:     "Execute shell commands",
		Category:        CategoryBuiltin,
		RequiresConfirm: true,
		FormatPurpose: func(args json.RawMessage) string {
			var input shellInput
			_ = json.Unmarshal(args, &input)
			return fmt.Sprintf("Execute command: %s", input.Command)
		},
	}
}
```

- [ ] **Step 4: Verify compilation**

```bash
go build ./backend/pkg/agent/tools/...
```

Expected: exit 0.

- [ ] **Step 5: Commit**

```bash
git add backend/pkg/agent/tools/datetime.go backend/pkg/agent/tools/file_rw.go backend/pkg/agent/tools/shell.go
git commit -m "feat(tools): add built-in tools — datetime, file_rw, shell"
```

---

### Task 7: User tools — web_search, code_exec

**Files:**
- Create: `backend/pkg/agent/tools/web_search.go`
- Create: `backend/pkg/agent/tools/code_exec.go`

- [ ] **Step 1: Implement web_search tool**

```go
// backend/pkg/agent/tools/web_search.go
package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

type webSearchInput struct {
	Query string `json:"query" jsonschema:"description=Search query string,required"`
	Limit int    `json:"limit" jsonschema:"description=Max number of results (default: 5)"`
}

type webSearchResult struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Snippet string `json:"snippet"`
}

type webSearchOutput struct {
	Results []webSearchResult `json:"results"`
}

// webSearchFunc performs a web search. Placeholder implementation for Phase 1.
func webSearchFunc(ctx context.Context, input webSearchInput) (webSearchOutput, error) {
	// TODO: integrate real search API (SearXNG / Bing / Google)
	return webSearchOutput{
		Results: []webSearchResult{
			{
				Title:   "Search not yet implemented",
				URL:     "",
				Snippet: fmt.Sprintf("Search for '%s' is not yet available. A search API integration is needed.", input.Query),
			},
		},
	}, nil
}

// NewWebSearchTool creates the user-facing web search tool.
func NewWebSearchTool() *function.FunctionTool[webSearchInput, webSearchOutput] {
	return function.NewFunctionTool(
		webSearchFunc,
		function.WithName("web_search"),
		function.WithDescription("Search the web for information"),
	)
}

// WebSearchMeta returns the ToolMeta for the web search tool.
func WebSearchMeta() ToolMeta {
	return ToolMeta{
		Name:            "web_search",
		Description:     "Search the web for information",
		Category:        CategoryUser,
		RequiresConfirm: false,
		FormatPurpose: func(args json.RawMessage) string {
			var input webSearchInput
			_ = json.Unmarshal(args, &input)
			return fmt.Sprintf("Search the web: %s", input.Query)
		},
	}
}
```

- [ ] **Step 2: Implement code_exec tool**

```go
// backend/pkg/agent/tools/code_exec.go
package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

const codeExecTimeout = 30 * time.Second

type codeExecInput struct {
	Language string `json:"language" jsonschema:"description=Programming language (python/javascript/bash),required"`
	Code     string `json:"code" jsonschema:"description=Source code to execute,required"`
}

type codeExecOutput struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exit_code"`
}

// codeExecFunc executes code in a subprocess.
func codeExecFunc(ctx context.Context, input codeExecInput) (codeExecOutput, error) {
	runners := map[string]struct {
		cmd string
		ext string
	}{
		"python":     {cmd: "python3", ext: ".py"},
		"javascript": {cmd: "node", ext: ".js"},
		"bash":       {cmd: "bash", ext: ".sh"},
	}

	runner, ok := runners[input.Language]
	if !ok {
		return codeExecOutput{}, fmt.Errorf("unsupported language: %s", input.Language)
	}

	tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("lemontea_exec_%d%s", time.Now().UnixNano(), runner.ext))
	if err := os.WriteFile(tmpFile, []byte(input.Code), 0o644); err != nil {
		return codeExecOutput{}, fmt.Errorf("write temp file: %w", err)
	}
	defer os.Remove(tmpFile)

	ctx, cancel := context.WithTimeout(ctx, codeExecTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, runner.cmd, tmpFile)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			return codeExecOutput{}, fmt.Errorf("exec: %w", err)
		}
	}

	return codeExecOutput{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
	}, nil
}

// NewCodeExecTool creates the user-facing code execution tool.
func NewCodeExecTool() *function.FunctionTool[codeExecInput, codeExecOutput] {
	return function.NewFunctionTool(
		codeExecFunc,
		function.WithName("code_exec"),
		function.WithDescription("Execute code in Python, JavaScript, or Bash and return the output"),
	)
}

// CodeExecMeta returns the ToolMeta for the code execution tool.
func CodeExecMeta() ToolMeta {
	return ToolMeta{
		Name:            "code_exec",
		Description:     "Execute code in Python, JavaScript, or Bash",
		Category:        CategoryUser,
		RequiresConfirm: true,
		FormatPurpose: func(args json.RawMessage) string {
			var input codeExecInput
			_ = json.Unmarshal(args, &input)
			preview := input.Code
			if len(preview) > 100 {
				preview = preview[:100] + "..."
			}
			return fmt.Sprintf("Execute %s code:\n%s", input.Language, preview)
		},
	}
}
```

- [ ] **Step 3: Verify compilation**

```bash
go build ./backend/pkg/agent/tools/...
```

Expected: exit 0.

- [ ] **Step 4: Commit**

```bash
git add backend/pkg/agent/tools/web_search.go backend/pkg/agent/tools/code_exec.go
git commit -m "feat(tools): add user tools — web_search, code_exec"
```

---

### Task 8: Event IDs

**Files:**
- Modify: `backend/pkg/id/event_id/id.go`

- [ ] **Step 1: Add agent stream event constants**

Append to `backend/pkg/id/event_id/id.go`:

```go
// AgentStreamChunk is emitted for each streaming text/thinking delta.
const AgentStreamChunk = "agent:stream:chunk"

// AgentStreamToolCall is emitted when the model invokes a tool.
const AgentStreamToolCall = "agent:stream:tool_call"

// AgentStreamConfirmRequest is emitted when a sensitive tool needs user confirmation.
const AgentStreamConfirmRequest = "agent:stream:confirm_request"

// AgentStreamToolResult is emitted after a tool finishes execution.
const AgentStreamToolResult = "agent:stream:tool_result"

// AgentStreamDone is emitted when the conversation turn completes.
const AgentStreamDone = "agent:stream:done"

// AgentStreamError is emitted when the conversation turn fails.
const AgentStreamError = "agent:stream:error"

// AgentSessionStatus is emitted when a session's status changes.
const AgentSessionStatus = "agent:session:status"
```

- [ ] **Step 2: Commit**

```bash
git add backend/pkg/id/event_id/id.go
git commit -m "feat(events): add agent stream event IDs"
```

---

### Task 9: Agent core — stream_manager.go

**Files:**
- Create: `backend/pkg/agent/stream_manager.go`

- [ ] **Step 1: Implement stream manager**

```go
// backend/pkg/agent/stream_manager.go
package agent

import (
	"context"
	"sync"
)

// activeStream tracks a single in-progress conversation stream.
type activeStream struct {
	cancel    context.CancelFunc
	confirmCh chan confirmResponse
}

// confirmResponse carries the user's response to a tool confirmation request.
type confirmResponse struct {
	Approved bool
	Message  string
}

// StreamManager tracks active streams keyed by session ID.
type StreamManager struct {
	mu      sync.RWMutex
	streams map[uint]*activeStream
}

// NewStreamManager creates a new StreamManager.
func NewStreamManager() *StreamManager {
	return &StreamManager{
		streams: make(map[uint]*activeStream),
	}
}

// Start registers a new active stream for the session and returns a cancellable context.
func (sm *StreamManager) Start(parentCtx context.Context, sessionID uint) (context.Context, context.CancelFunc) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if existing, ok := sm.streams[sessionID]; ok {
		existing.cancel()
	}

	ctx, cancel := context.WithCancel(parentCtx)
	sm.streams[sessionID] = &activeStream{
		cancel:    cancel,
		confirmCh: make(chan confirmResponse, 1),
	}
	return ctx, cancel
}

// Stop cancels the stream for the given session and removes it.
func (sm *StreamManager) Stop(sessionID uint) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if s, ok := sm.streams[sessionID]; ok {
		s.cancel()
		close(s.confirmCh)
		delete(sm.streams, sessionID)
	}
}

// Remove removes a stream without cancelling (used when a stream finishes normally).
func (sm *StreamManager) Remove(sessionID uint) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.streams, sessionID)
}

// IsActive returns whether a stream is active for the session.
func (sm *StreamManager) IsActive(sessionID uint) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	_, ok := sm.streams[sessionID]
	return ok
}

// SendConfirmResponse sends a user confirmation response for the given session.
func (sm *StreamManager) SendConfirmResponse(sessionID uint, resp confirmResponse) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	s, ok := sm.streams[sessionID]
	if !ok {
		return false
	}
	select {
	case s.confirmCh <- resp:
		return true
	default:
		return false
	}
}

// WaitForConfirm blocks until the user responds to a confirmation request or the context is cancelled.
func (sm *StreamManager) WaitForConfirm(ctx context.Context, sessionID uint) (confirmResponse, error) {
	sm.mu.RLock()
	s, ok := sm.streams[sessionID]
	sm.mu.RUnlock()

	if !ok {
		return confirmResponse{}, context.Canceled
	}

	select {
	case <-ctx.Done():
		return confirmResponse{}, ctx.Err()
	case resp, chanOpen := <-s.confirmCh:
		if !chanOpen {
			return confirmResponse{}, context.Canceled
		}
		return resp, nil
	}
}
```

- [ ] **Step 2: Verify compilation**

```bash
go build ./backend/pkg/agent/...
```

Expected: exit 0.

- [ ] **Step 3: Commit**

```bash
git add backend/pkg/agent/stream_manager.go
git commit -m "feat(agent): add stream manager for active stream tracking and cancellation"
```

---

### Task 10: Agent core — manager.go

**Files:**
- Create: `backend/pkg/agent/manager.go`

- [ ] **Step 1: Implement manager**

```go
// backend/pkg/agent/manager.go
package agent

import (
	"fmt"
	"sync"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/agent/tools"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/id/prompt_id"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/prompt"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/storage"
	"github.com/wailsapp/wails/v3/pkg/application"
	"trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/agent/llmagent"
	"trpc.group/trpc-go/trpc-agent-go/model"
	"trpc.group/trpc-go/trpc-agent-go/model/openai"
	"trpc.group/trpc-go/trpc-agent-go/runner"
	"trpc.group/trpc-go/trpc-agent-go/session/inmemory"
	toolpkg "trpc.group/trpc-go/trpc-agent-go/tool"
)

const defaultMainAgentPrompt = `You are Lemontea, a helpful AI assistant. You can help users with various tasks including answering questions, writing, coding, and more. Be concise, accurate, and helpful. When using tools, explain what you are doing and why.`

// Manager manages the Agent lifecycle, runner instances, and tool registry.
type Manager struct {
	mu            sync.RWMutex
	istorage      *storage.Storage
	wailsApp      *application.App
	toolRegistry  *tools.Registry
	streamManager *StreamManager
	runners       map[string]runner.Runner
}

// NewManager creates a new agent manager.
func NewManager(istorage *storage.Storage) *Manager {
	prompt.Register(prompt_id.MainAgent, defaultMainAgentPrompt)

	registry := tools.NewRegistry()
	registry.Register(tools.DateTimeMeta())
	registry.Register(tools.FileReadMeta())
	registry.Register(tools.FileWriteMeta())
	registry.Register(tools.ShellMeta())
	registry.Register(tools.WebSearchMeta())
	registry.Register(tools.CodeExecMeta())

	return &Manager{
		istorage:      istorage,
		toolRegistry:  registry,
		streamManager: NewStreamManager(),
		runners:       make(map[string]runner.Runner),
	}
}

// SetApp sets the Wails application reference (called during ServiceStartup).
func (m *Manager) SetApp(app *application.App) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.wailsApp = app
}

// ToolRegistry returns the tool registry.
func (m *Manager) ToolRegistry() *tools.Registry {
	return m.toolRegistry
}

// StreamManager returns the stream manager.
func (m *Manager) Streams() *StreamManager {
	return m.streamManager
}

// Storage returns the storage handle.
func (m *Manager) Storage() *storage.Storage {
	return m.istorage
}

// App returns the Wails application reference.
func (m *Manager) App() *application.App {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.wailsApp
}

// buildAgentTools converts enabled ToolMeta list into trpc-agent-go tool.Tool slice.
func (m *Manager) buildAgentTools(enabledUserTools []string) []toolpkg.Tool {
	metas := m.toolRegistry.EnabledTools(enabledUserTools)
	var agentTools []toolpkg.Tool
	for _, meta := range metas {
		switch meta.Name {
		case "datetime":
			agentTools = append(agentTools, tools.NewDateTimeTool())
		case "file_read":
			agentTools = append(agentTools, tools.NewFileReadTool())
		case "file_write":
			agentTools = append(agentTools, tools.NewFileWriteTool())
		case "shell":
			agentTools = append(agentTools, tools.NewShellTool())
		case "web_search":
			agentTools = append(agentTools, tools.NewWebSearchTool())
		case "code_exec":
			agentTools = append(agentTools, tools.NewCodeExecTool())
		}
	}
	return agentTools
}

// GetOrCreateRunner returns an existing runner or creates a new one for the given provider config.
func (m *Manager) GetOrCreateRunner(baseURL, apiKey, modelName string, enabledUserTools []string) (runner.Runner, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	key := fmt.Sprintf("%s/%s", baseURL, modelName)
	if r, ok := m.runners[key]; ok {
		return r, nil
	}

	instruction, err := prompt.Load(prompt_id.MainAgent)
	if err != nil {
		instruction = defaultMainAgentPrompt
	}

	mdl := openai.New(modelName,
		openai.WithBaseURL(baseURL),
		openai.WithAPIKey(apiKey),
	)

	agentTools := m.buildAgentTools(enabledUserTools)

	ag := llmagent.New("main",
		llmagent.WithModel(mdl),
		llmagent.WithInstruction(instruction),
		llmagent.WithGenerationConfig(model.GenerationConfig{
			Stream: true,
		}),
		llmagent.WithTools(agentTools),
	)

	sessionSvc := inmemory.NewSessionService()

	r := runner.NewRunner("lemontea", ag,
		runner.WithSessionService(sessionSvc),
	)

	m.runners[key] = r
	return r, nil
}

// NewChatHandler creates a ChatHandler bound to this manager.
func (m *Manager) NewChatHandler() *ChatHandler {
	return NewChatHandler(m)
}
```

- [ ] **Step 2: Verify compilation**

```bash
go build ./backend/pkg/agent/...
```

Expected: exit 0.

- [ ] **Step 3: Commit**

```bash
git add backend/pkg/agent/manager.go
git commit -m "feat(agent): add manager for agent lifecycle and runner factory"
```

---

### Task 11: Agent core — chat_handler.go

**Files:**
- Create: `backend/pkg/agent/chat_handler.go`

- [ ] **Step 1: Implement chat handler**

```go
// backend/pkg/agent/chat_handler.go
package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/id/event_id"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/model"
)

// ChatHandler processes chat requests and bridges events to Wails.
type ChatHandler struct {
	manager *Manager
}

// NewChatHandler creates a ChatHandler from the given manager.
func NewChatHandler(m *Manager) *ChatHandler {
	return &ChatHandler{manager: m}
}

// SendMessageParams holds the parameters for sending a message.
type SendMessageParams struct {
	SessionID        uint
	Content          string
	BaseURL          string
	ApiKey           string
	ModelName        string
	EnabledUserTools []string
}

// emitEvent sends a Wails event to all windows.
func (ch *ChatHandler) emitEvent(name string, data any) {
	app := ch.manager.App()
	if app != nil {
		app.EmitEvent(name, data)
	}
}

// updateSessionStatus updates the session status in storage and emits the event.
func (ch *ChatHandler) updateSessionStatus(sessionID uint, status string) {
	_ = ch.manager.Storage().UpdateSessionStatus(sessionID, status)
	ch.emitEvent(event_id.AgentSessionStatus, map[string]any{
		"sessionId": sessionID,
		"status":    status,
	})
}

// SendMessage starts a streaming conversation. It runs in a goroutine and pushes events to the frontend.
func (ch *ChatHandler) SendMessage(ctx context.Context, params SendMessageParams) error {
	stor := ch.manager.Storage()

	// Save user message to storage.
	_, err := stor.CreateMessage(data_models.Message{
		SessionID:   params.SessionID,
		Role:        "user",
		ContentType: "text",
		Content:     params.Content,
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

	streamCtx, cancel := ch.manager.Streams().Start(ctx, params.SessionID)
	userID := "local"
	sessionIDStr := strconv.FormatUint(uint64(params.SessionID), 10)

	events, err := r.Run(streamCtx, userID, sessionIDStr, model.NewUserMessage(params.Content))
	if err != nil {
		cancel()
		ch.manager.Streams().Remove(params.SessionID)
		ch.updateSessionStatus(params.SessionID, "error-unread")
		return fmt.Errorf("runner.Run: %w", err)
	}

	go ch.consumeEvents(streamCtx, params.SessionID, params.ModelName, events)
	return nil
}

// consumeEvents reads from the event channel and emits Wails events.
func (ch *ChatHandler) consumeEvents(ctx context.Context, sessionID uint, modelName string, events <-chan *event.Event) {
	defer ch.manager.Streams().Remove(sessionID)

	stor := ch.manager.Storage()
	var fullContent string
	var fullThinking string
	var tokensIn, tokensOut int

	for evt := range events {
		if evt == nil {
			continue
		}

		select {
		case <-ctx.Done():
			ch.updateSessionStatus(sessionID, "idle")
			return
		default:
		}

		if evt.Error != nil {
			ch.emitEvent(event_id.AgentStreamError, map[string]any{
				"sessionId": sessionID,
				"error":     evt.Error.Message,
			})
			ch.updateSessionStatus(sessionID, "error-unread")
			// Save error message.
			_, _ = stor.CreateMessage(data_models.Message{
				SessionID:   sessionID,
				Role:        "assistant",
				ContentType: "text",
				Content:     fmt.Sprintf("Error: %s", evt.Error.Message),
				ModelName:   modelName,
			})
			return
		}

		// Process streaming choices.
		for _, choice := range evt.Choices {
			delta := choice.Delta

			// Thinking content.
			if delta.ReasoningContent != "" {
				fullThinking += delta.ReasoningContent
				ch.emitEvent(event_id.AgentStreamChunk, map[string]any{
					"sessionId":   sessionID,
					"delta":       delta.ReasoningContent,
					"contentType": "thinking",
				})
			}

			// Text content.
			if delta.Content != "" {
				fullContent += delta.Content
				ch.emitEvent(event_id.AgentStreamChunk, map[string]any{
					"sessionId":   sessionID,
					"delta":       delta.Content,
					"contentType": "text",
				})
			}

			// Tool calls.
			for _, tc := range delta.ToolCalls {
				toolName := tc.ToolName
				argsJSON, _ := json.Marshal(tc)

				meta, hasMeta := ch.manager.ToolRegistry().Get(toolName)
				purpose := ""
				if hasMeta && meta.FormatPurpose != nil {
					purpose = meta.FormatPurpose(argsJSON)
				}

				ch.emitEvent(event_id.AgentStreamToolCall, map[string]any{
					"sessionId": sessionID,
					"toolName":  toolName,
					"args":      string(argsJSON),
					"purpose":   purpose,
				})

				// Save tool_call message.
				_, _ = stor.CreateMessage(data_models.Message{
					SessionID:   sessionID,
					Role:        "assistant",
					ContentType: "tool_call",
					Content:     string(argsJSON),
					ModelName:   modelName,
					Extra:       purpose,
				})

				// Handle confirmation for sensitive tools.
				if hasMeta && meta.RequiresConfirm {
					ch.emitEvent(event_id.AgentStreamConfirmRequest, map[string]any{
						"sessionId": sessionID,
						"requestId": tc.ToolID,
						"toolName":  toolName,
						"args":      string(argsJSON),
						"purpose":   purpose,
					})
					ch.updateSessionStatus(sessionID, "waiting-unread")

					resp, err := ch.manager.Streams().WaitForConfirm(ctx, sessionID)
					if err != nil {
						ch.updateSessionStatus(sessionID, "idle")
						return
					}

					// Save confirm response.
					confirmContent := "approved"
					if !resp.Approved {
						confirmContent = "rejected"
						if resp.Message != "" {
							confirmContent = fmt.Sprintf("rejected: %s", resp.Message)
						}
					}
					_, _ = stor.CreateMessage(data_models.Message{
						SessionID:   sessionID,
						Role:        "user",
						ContentType: "confirm_response",
						Content:     confirmContent,
					})

					ch.updateSessionStatus(sessionID, "loading")
				}
			}
		}

		// Usage info.
		if evt.Usage != nil {
			tokensIn = evt.Usage.PromptTokens
			tokensOut = evt.Usage.CompletionTokens
		}

		if evt.Done {
			break
		}
	}

	// Save thinking message if present.
	if fullThinking != "" {
		_, _ = stor.CreateMessage(data_models.Message{
			SessionID:   sessionID,
			Role:        "assistant",
			ContentType: "thinking",
			Content:     fullThinking,
			ModelName:   modelName,
		})
	}

	// Save assistant text message.
	if fullContent != "" {
		_, _ = stor.CreateMessage(data_models.Message{
			SessionID:   sessionID,
			Role:        "assistant",
			ContentType: "text",
			Content:     fullContent,
			ModelName:   modelName,
			TokensIn:    tokensIn,
			TokensOut:   tokensOut,
		})
	}

	ch.emitEvent(event_id.AgentStreamDone, map[string]any{
		"sessionId": sessionID,
		"usage": map[string]int{
			"input":  tokensIn,
			"output": tokensOut,
		},
	})

	ch.updateSessionStatus(sessionID, "done-unread")
}
```

- [ ] **Step 2: Verify compilation**

```bash
go build ./backend/pkg/agent/...
```

Expected: exit 0.

- [ ] **Step 3: Commit**

```bash
git add backend/pkg/agent/chat_handler.go
git commit -m "feat(agent): add chat handler for streaming conversation processing"
```

---

### Task 12: Service layer — DTOs

**Files:**
- Create: `backend/service/agent/agent_dto/send_message.go`
- Create: `backend/service/agent/agent_dto/stop_generation.go`
- Create: `backend/service/agent/agent_dto/respond_to_confirm.go`
- Create: `backend/service/agent/agent_dto/create_session.go`
- Create: `backend/service/agent/agent_dto/list_sessions.go`
- Create: `backend/service/agent/agent_dto/load_session_messages.go`
- Create: `backend/service/agent/agent_dto/rename_session.go`
- Create: `backend/service/agent/agent_dto/delete_session.go`
- Create: `backend/service/agent/agent_dto/toggle_star_session.go`
- Create: `backend/service/agent/agent_dto/generate_title.go`

- [ ] **Step 1: Create all DTO files**

```go
// backend/service/agent/agent_dto/send_message.go
package agent_dto

type SendMessageInput struct {
	SessionID        uint     `json:"session_id"`
	Content          string   `json:"content"`
	BaseURL          string   `json:"base_url"`
	ApiKey           string   `json:"api_key"`
	ModelName        string   `json:"model_name"`
	EnabledUserTools []string `json:"enabled_user_tools"`
}

type SendMessageOutput struct {
}
```

```go
// backend/service/agent/agent_dto/stop_generation.go
package agent_dto

type StopGenerationInput struct {
	SessionID uint `json:"session_id"`
}

type StopGenerationOutput struct {
}
```

```go
// backend/service/agent/agent_dto/respond_to_confirm.go
package agent_dto

type RespondToConfirmInput struct {
	SessionID uint   `json:"session_id"`
	Approved  bool   `json:"approved"`
	Message   string `json:"message"`
}

type RespondToConfirmOutput struct {
}
```

```go
// backend/service/agent/agent_dto/create_session.go
package agent_dto

type CreateSessionInput struct {
	Title string `json:"title"`
}

type CreateSessionOutput struct {
	SessionID uint   `json:"session_id"`
	Title     string `json:"title"`
}
```

```go
// backend/service/agent/agent_dto/list_sessions.go
package agent_dto

type ListSessionsInput struct {
	Cursor      uint `json:"cursor"`
	Limit       int  `json:"limit"`
	StarredOnly bool `json:"starred_only"`
}

type SessionItem struct {
	ID      uint   `json:"id"`
	Title   string `json:"title"`
	Starred bool   `json:"starred"`
	Status  string `json:"status"`
	Created string `json:"created"`
	Updated string `json:"updated"`
}

type ListSessionsOutput struct {
	Sessions   []SessionItem `json:"sessions"`
	NextCursor uint          `json:"next_cursor"`
	HasMore    bool          `json:"has_more"`
}
```

```go
// backend/service/agent/agent_dto/load_session_messages.go
package agent_dto

type LoadSessionMessagesInput struct {
	SessionID uint `json:"session_id"`
	Offset    int  `json:"offset"`
	Limit     int  `json:"limit"`
}

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
	Created     string `json:"created"`
}

type LoadSessionMessagesOutput struct {
	Messages []MessageItem `json:"messages"`
	Total    int64         `json:"total"`
	HasMore  bool          `json:"has_more"`
}
```

```go
// backend/service/agent/agent_dto/rename_session.go
package agent_dto

type RenameSessionInput struct {
	SessionID uint   `json:"session_id"`
	Title     string `json:"title"`
}

type RenameSessionOutput struct {
}
```

```go
// backend/service/agent/agent_dto/delete_session.go
package agent_dto

type DeleteSessionInput struct {
	SessionID uint `json:"session_id"`
}

type DeleteSessionOutput struct {
}
```

```go
// backend/service/agent/agent_dto/toggle_star_session.go
package agent_dto

type ToggleStarSessionInput struct {
	SessionID uint `json:"session_id"`
	Starred   bool `json:"starred"`
}

type ToggleStarSessionOutput struct {
}
```

```go
// backend/service/agent/agent_dto/generate_title.go
package agent_dto

type GenerateTitleInput struct {
	SessionID uint   `json:"session_id"`
	BaseURL   string `json:"base_url"`
	ApiKey    string `json:"api_key"`
	ModelName string `json:"model_name"`
}

type GenerateTitleOutput struct {
	Title string `json:"title"`
}
```

- [ ] **Step 2: Verify compilation**

```bash
go build ./backend/service/agent/agent_dto/...
```

Expected: exit 0.

- [ ] **Step 3: Commit**

```bash
git add backend/service/agent/agent_dto/
git commit -m "feat(agent): add all agent service DTOs"
```

---

### Task 13: Service layer — agent.go + implement + internal

**Files:**
- Create: `backend/service/agent/agent.go`
- Create: `backend/service/agent/agent_implement.go`
- Create: `backend/service/agent/agent_internal.go`

- [ ] **Step 1: Create agent_implement.go**

```go
// backend/service/agent/agent_implement.go
package agent

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func (a *Agent) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	a.manager.SetApp(application.Get())
	return nil
}
```

- [ ] **Step 2: Create agent_internal.go**

```go
// backend/service/agent/agent_internal.go
package agent

import (
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/agent/agent_dto"
)

const timeFormat = time.RFC3339

// toSessionItem converts a data model Session to a DTO SessionItem.
func toSessionItem(s data_models.Session) agent_dto.SessionItem {
	return agent_dto.SessionItem{
		ID:      s.ID,
		Title:   s.Title,
		Starred: s.Starred,
		Status:  s.Status,
		Created: s.CreatedAt.Format(timeFormat),
		Updated: s.UpdatedAt.Format(timeFormat),
	}
}

// toMessageItem converts a data model Message to a DTO MessageItem.
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
		Created:     m.CreatedAt.Format(timeFormat),
	}
}
```

- [ ] **Step 3: Create agent.go with all public methods**

```go
// backend/service/agent/agent.go
package agent

import (
	"context"
	"fmt"

	pkgAgent "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/agent"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/agent/agent_dto"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/storage"
)

// Agent is the Wails service that exposes agent operations to the frontend.
type Agent struct {
	manager *pkgAgent.Manager
}

// NewAgent creates the agent service bound to the given storage.
func NewAgent(istorage *storage.Storage) *Agent {
	return &Agent{
		manager: pkgAgent.NewManager(istorage),
	}
}

// SendMessage starts a streaming conversation for the given session.
func (a *Agent) SendMessage(ctx context.Context, input agent_dto.SendMessageInput) (*agent_dto.SendMessageOutput, error) {
	handler := a.manager.NewChatHandler()
	err := handler.SendMessage(ctx, pkgAgent.SendMessageParams{
		SessionID:        input.SessionID,
		Content:          input.Content,
		BaseURL:          input.BaseURL,
		ApiKey:           input.ApiKey,
		ModelName:        input.ModelName,
		EnabledUserTools: input.EnabledUserTools,
	})
	if err != nil {
		return nil, fmt.Errorf("send message: %w", err)
	}
	return &agent_dto.SendMessageOutput{}, nil
}

// StopGeneration cancels the active stream for the given session.
func (a *Agent) StopGeneration(ctx context.Context, input agent_dto.StopGenerationInput) (*agent_dto.StopGenerationOutput, error) {
	a.manager.Streams().Stop(input.SessionID)
	_ = a.manager.Storage().UpdateSessionStatus(input.SessionID, "idle")
	return &agent_dto.StopGenerationOutput{}, nil
}

// RespondToConfirm sends the user's confirmation response for a tool call.
func (a *Agent) RespondToConfirm(ctx context.Context, input agent_dto.RespondToConfirmInput) (*agent_dto.RespondToConfirmOutput, error) {
	a.manager.Streams().SendConfirmResponse(input.SessionID, pkgAgent.ConfirmResponse{
		Approved: input.Approved,
		Message:  input.Message,
	})
	return &agent_dto.RespondToConfirmOutput{}, nil
}

// CreateSession creates a new conversation session.
func (a *Agent) CreateSession(ctx context.Context, input agent_dto.CreateSessionInput) (*agent_dto.CreateSessionOutput, error) {
	title := input.Title
	if title == "" {
		title = "New Chat"
	}
	session, err := a.manager.Storage().CreateSession(data_models.Session{
		Title:  title,
		Status: "idle",
	})
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}
	return &agent_dto.CreateSessionOutput{
		SessionID: session.ID,
		Title:     session.Title,
	}, nil
}

// ListSessions returns sessions with cursor-based pagination.
func (a *Agent) ListSessions(ctx context.Context, input agent_dto.ListSessionsInput) (*agent_dto.ListSessionsOutput, error) {
	limit := input.Limit
	if limit <= 0 {
		limit = 20
	}

	sessions, err := a.manager.Storage().ListSessions(input.Cursor, limit+1, input.StarredOnly)
	if err != nil {
		return nil, fmt.Errorf("list sessions: %w", err)
	}

	hasMore := len(sessions) > limit
	if hasMore {
		sessions = sessions[:limit]
	}

	items := make([]agent_dto.SessionItem, 0, len(sessions))
	for _, s := range sessions {
		items = append(items, toSessionItem(s))
	}

	var nextCursor uint
	if hasMore && len(items) > 0 {
		nextCursor = items[len(items)-1].ID
	}

	return &agent_dto.ListSessionsOutput{
		Sessions:   items,
		NextCursor: nextCursor,
		HasMore:    hasMore,
	}, nil
}

// LoadSessionMessages returns messages for a session with offset pagination.
func (a *Agent) LoadSessionMessages(ctx context.Context, input agent_dto.LoadSessionMessagesInput) (*agent_dto.LoadSessionMessagesOutput, error) {
	limit := input.Limit
	if limit <= 0 {
		limit = 50
	}

	msgs, err := a.manager.Storage().ListMessagesForSession(input.SessionID, input.Offset, limit)
	if err != nil {
		return nil, fmt.Errorf("list messages: %w", err)
	}

	total, err := a.manager.Storage().CountMessagesForSession(input.SessionID)
	if err != nil {
		return nil, fmt.Errorf("count messages: %w", err)
	}

	items := make([]agent_dto.MessageItem, 0, len(msgs))
	for _, m := range msgs {
		items = append(items, toMessageItem(m))
	}

	return &agent_dto.LoadSessionMessagesOutput{
		Messages: items,
		Total:    total,
		HasMore:  int64(input.Offset+len(msgs)) < total,
	}, nil
}

// RenameSession updates the title of a session.
func (a *Agent) RenameSession(ctx context.Context, input agent_dto.RenameSessionInput) (*agent_dto.RenameSessionOutput, error) {
	if err := a.manager.Storage().UpdateSessionTitle(input.SessionID, input.Title); err != nil {
		return nil, fmt.Errorf("rename session: %w", err)
	}
	return &agent_dto.RenameSessionOutput{}, nil
}

// DeleteSession removes a session and all its messages.
func (a *Agent) DeleteSession(ctx context.Context, input agent_dto.DeleteSessionInput) (*agent_dto.DeleteSessionOutput, error) {
	a.manager.Streams().Stop(input.SessionID)
	if err := a.manager.Storage().DeleteSession(input.SessionID); err != nil {
		return nil, fmt.Errorf("delete session: %w", err)
	}
	return &agent_dto.DeleteSessionOutput{}, nil
}

// ToggleStarSession sets the starred flag on a session.
func (a *Agent) ToggleStarSession(ctx context.Context, input agent_dto.ToggleStarSessionInput) (*agent_dto.ToggleStarSessionOutput, error) {
	if err := a.manager.Storage().UpdateSessionStarred(input.SessionID, input.Starred); err != nil {
		return nil, fmt.Errorf("toggle star: %w", err)
	}
	return &agent_dto.ToggleStarSessionOutput{}, nil
}

// GenerateTitle uses the LLM to generate a short conversation title from the first messages.
func (a *Agent) GenerateTitle(ctx context.Context, input agent_dto.GenerateTitleInput) (*agent_dto.GenerateTitleOutput, error) {
	msgs, err := a.manager.Storage().ListMessagesForSession(input.SessionID, 0, 4)
	if err != nil || len(msgs) == 0 {
		return &agent_dto.GenerateTitleOutput{Title: "New Chat"}, nil
	}

	var conversationPreview string
	for _, m := range msgs {
		if m.ContentType == "text" {
			preview := m.Content
			if len(preview) > 200 {
				preview = preview[:200]
			}
			conversationPreview += fmt.Sprintf("%s: %s\n", m.Role, preview)
		}
	}

	// Use the LLM to generate a title.
	r, err := a.manager.GetOrCreateRunner(input.BaseURL, input.ApiKey, input.ModelName, nil)
	if err != nil {
		return &agent_dto.GenerateTitleOutput{Title: "New Chat"}, nil
	}

	titlePrompt := fmt.Sprintf("Generate a concise title (under 20 characters, in the same language as the conversation) for this conversation. Return ONLY the title, no quotes or explanation.\n\n%s", conversationPreview)

	events, err := r.Run(ctx, "local", fmt.Sprintf("title-%d", input.SessionID),
		model.NewUserMessage(titlePrompt),
	)
	if err != nil {
		return &agent_dto.GenerateTitleOutput{Title: "New Chat"}, nil
	}

	var title string
	for evt := range events {
		if evt == nil {
			continue
		}
		for _, choice := range evt.Choices {
			title += choice.Delta.Content
		}
	}

	if title == "" {
		title = "New Chat"
	}

	_ = a.manager.Storage().UpdateSessionTitle(input.SessionID, title)

	return &agent_dto.GenerateTitleOutput{Title: title}, nil
}
```

Note: The `GenerateTitle` method uses `model` from `trpc.group/trpc-go/trpc-agent-go/model`. Add this import:

```go
import "trpc.group/trpc-go/trpc-agent-go/model"
```

- [ ] **Step 4: Verify compilation**

```bash
go build ./backend/service/agent/...
```

Expected: exit 0.

- [ ] **Step 5: Commit**

```bash
git add backend/service/agent/
git commit -m "feat(agent): add Wails agent service with all public methods"
```

---

### Task 14: Register agent service in main.go + export ConfirmResponse

**Files:**
- Modify: `main.go`
- Modify: `backend/pkg/agent/stream_manager.go`

- [ ] **Step 1: Export ConfirmResponse type**

The `confirmResponse` struct in `stream_manager.go` is used by the service layer, so it must be exported. Rename `confirmResponse` to `ConfirmResponse` in `backend/pkg/agent/stream_manager.go`:

```go
// ConfirmResponse carries the user's response to a tool confirmation request.
type ConfirmResponse struct {
	Approved bool
	Message  string
}
```

Update all references in the same file: `confirmResponse` → `ConfirmResponse`.

- [ ] **Step 2: Register agent service in main.go**

Add import and service registration:

```go
import (
	// ... existing imports ...
	agentSvc "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/agent"
)
```

In the `Services` slice, add:

```go
application.NewService(agentSvc.NewAgent(istorage)),
```

- [ ] **Step 3: Verify compilation**

```bash
go build ./...
```

Expected: exit 0.

- [ ] **Step 4: Commit**

```bash
git add main.go backend/pkg/agent/stream_manager.go
git commit -m "feat: register agent service in main and export ConfirmResponse"
```

---

### Task 15: Frontend — types and i18n updates

**Files:**
- Modify: `frontend/src/types/index.ts`
- Modify: `frontend/src/i18n/locales/zh-CN.ts`
- Modify: `frontend/src/i18n/locales/en.ts`

- [ ] **Step 1: Update types/index.ts**

Replace the existing `Message` interface and add new types:

```typescript
// frontend/src/types/index.ts
export type Theme = 'auto' | 'light' | 'dark'
export type FontSize = 'xs' | 'sm' | 'md' | 'lg' | 'xl'
export type Language = 'zh-CN' | 'en'
export type ConversationTab = 'chats' | 'favorites'
export type ConversationStatus =
  | 'idle'
  | 'loading'
  | 'done-unread'
  | 'error-unread'
  | 'waiting-unread'

export type MessageContentType =
  | 'text'
  | 'tool_call'
  | 'tool_result'
  | 'thinking'
  | 'confirm_request'
  | 'confirm_response'

export type MessageRole = 'user' | 'assistant' | 'tool' | 'system'

export type Conversation = {
  id: number
  title: string
  createdAt: string
  updatedAt: string
  starred: boolean
  status: ConversationStatus
}

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
  createdAt: string
}

export type StreamChunkEvent = {
  sessionId: number
  delta: string
  contentType: 'text' | 'thinking'
}

export type ToolCallEvent = {
  sessionId: number
  toolName: string
  args: string
  purpose: string
}

export type ConfirmRequestEvent = {
  sessionId: number
  requestId: string
  toolName: string
  args: string
  purpose: string
}

export type StreamDoneEvent = {
  sessionId: number
  usage: { input: number; output: number }
}

export type StreamErrorEvent = {
  sessionId: number
  error: string
}

export type SessionStatusEvent = {
  sessionId: number
  status: ConversationStatus
}

export type Tool = {
  id: string
  name: string
  description: string
  enabled: boolean
}

export type ModelProvider = {
  id: string
  name: string
  models: Model[]
}

export type Model = {
  id: string
  name: string
  providerId: string
}
```

- [ ] **Step 2: Add i18n keys to zh-CN.ts**

Add under the existing keys in `frontend/src/i18n/locales/zh-CN.ts`:

```typescript
  toolCall: {
    executing: '正在执行',
    completed: '执行完成',
    failed: '执行失败',
    confirm: '确认执行',
    reject: '拒绝',
    approved: '已确认',
    rejected: '已拒绝',
    purpose: '用途',
    waitingConfirm: '等待确认...',
    inputPlaceholder: '输入修改意见（可选）',
  },
```

- [ ] **Step 3: Add i18n keys to en.ts**

Add the same structure in `frontend/src/i18n/locales/en.ts`:

```typescript
  toolCall: {
    executing: 'Executing',
    completed: 'Completed',
    failed: 'Failed',
    confirm: 'Confirm',
    reject: 'Reject',
    approved: 'Approved',
    rejected: 'Rejected',
    purpose: 'Purpose',
    waitingConfirm: 'Waiting for confirmation...',
    inputPlaceholder: 'Enter feedback (optional)',
  },
```

- [ ] **Step 4: Commit**

```bash
git add frontend/src/types/index.ts frontend/src/i18n/locales/zh-CN.ts frontend/src/i18n/locales/en.ts
git commit -m "feat(frontend): update types and i18n for agent integration"
```

---

### Task 16: Frontend — ToolCallBlock and ToolConfirmCard components

**Files:**
- Create: `frontend/src/components/chat/ToolCallBlock.tsx`
- Create: `frontend/src/components/chat/ToolConfirmCard.tsx`

- [ ] **Step 1: Create ToolCallBlock**

```tsx
// frontend/src/components/chat/ToolCallBlock.tsx
import { useState } from 'react'
import { ChevronDown, ChevronRight, Wrench, Check, X } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { cn } from '@/lib/utils'

type Props = {
  toolName: string
  purpose: string
  args: string
  result: string
  status: 'executing' | 'completed' | 'failed'
}

export function ToolCallBlock({ toolName, purpose, args, result, status }: Props) {
  const { t } = useTranslation()
  const [expanded, setExpanded] = useState(false)

  const statusIcon = {
    executing: <Wrench size={12} className="animate-spin" />,
    completed: <Check size={12} className="text-green-500" />,
    failed: <X size={12} className="text-red-500" />,
  }[status]

  const statusText = t(`toolCall.${status}`)

  return (
    <div className="my-2 rounded-lg border border-border bg-muted/30 overflow-hidden">
      <button
        onClick={() => setExpanded(v => !v)}
        className="flex items-center gap-2 w-full px-3 py-2 text-xs hover:bg-muted transition-colors"
      >
        {statusIcon}
        <Wrench size={12} className="text-muted-foreground" />
        <span className="font-medium text-foreground">{toolName}</span>
        <span className="text-muted-foreground flex-1 text-left truncate">{purpose}</span>
        <span className={cn(
          'text-muted-foreground',
          status === 'completed' && 'text-green-600',
          status === 'failed' && 'text-red-600',
        )}>
          {statusText}
        </span>
        {expanded ? <ChevronDown size={12} /> : <ChevronRight size={12} />}
      </button>
      {expanded && (
        <div className="border-t border-border px-3 py-2 space-y-2">
          {args && (
            <div>
              <div className="text-[10px] font-semibold text-muted-foreground uppercase mb-1">Args</div>
              <pre className="text-xs bg-muted rounded p-2 overflow-x-auto whitespace-pre-wrap">{args}</pre>
            </div>
          )}
          {result && (
            <div>
              <div className="text-[10px] font-semibold text-muted-foreground uppercase mb-1">Result</div>
              <pre className="text-xs bg-muted rounded p-2 overflow-x-auto whitespace-pre-wrap max-h-48 overflow-y-auto">{result}</pre>
            </div>
          )}
        </div>
      )}
    </div>
  )
}
```

- [ ] **Step 2: Create ToolConfirmCard**

```tsx
// frontend/src/components/chat/ToolConfirmCard.tsx
import { useState } from 'react'
import { ShieldAlert, Check, X } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { cn } from '@/lib/utils'

type Props = {
  toolName: string
  purpose: string
  args: string
  onConfirm: (message: string) => void
  onReject: (message: string) => void
  resolved: boolean
  approved: boolean
}

export function ToolConfirmCard({ toolName, purpose, args, onConfirm, onReject, resolved, approved }: Props) {
  const { t } = useTranslation()
  const [feedback, setFeedback] = useState('')

  if (resolved) {
    return (
      <div className={cn(
        'my-2 flex items-center gap-2 rounded-lg border px-3 py-2 text-xs',
        approved
          ? 'border-green-200 bg-green-50 dark:border-green-800 dark:bg-green-950'
          : 'border-red-200 bg-red-50 dark:border-red-800 dark:bg-red-950'
      )}>
        {approved ? <Check size={14} className="text-green-600" /> : <X size={14} className="text-red-600" />}
        <span className="font-medium">{toolName}</span>
        <span className="text-muted-foreground">
          {approved ? t('toolCall.approved') : t('toolCall.rejected')}
        </span>
      </div>
    )
  }

  return (
    <div className="my-2 rounded-lg border border-amber-200 bg-amber-50 dark:border-amber-800 dark:bg-amber-950 overflow-hidden">
      <div className="flex items-center gap-2 px-3 py-2">
        <ShieldAlert size={14} className="text-amber-600 shrink-0" />
        <span className="text-xs font-medium text-foreground">{toolName}</span>
        <span className="text-xs text-muted-foreground">{t('toolCall.waitingConfirm')}</span>
      </div>
      <div className="px-3 pb-2 text-xs text-muted-foreground">
        <div className="font-medium text-foreground mb-1">{t('toolCall.purpose')}: {purpose}</div>
        {args && (
          <pre className="bg-muted rounded p-2 overflow-x-auto whitespace-pre-wrap mb-2 max-h-32 overflow-y-auto">{args}</pre>
        )}
      </div>
      <div className="border-t border-amber-200 dark:border-amber-800 px-3 py-2 flex items-center gap-2">
        <input
          value={feedback}
          onChange={e => setFeedback(e.target.value)}
          placeholder={t('toolCall.inputPlaceholder')}
          className="flex-1 min-w-0 text-xs bg-transparent border border-input rounded px-2 py-1 outline-none focus:ring-1 focus:ring-ring placeholder:text-muted-foreground"
        />
        <button
          onClick={() => onReject(feedback)}
          className="flex items-center gap-1 px-2 py-1 rounded text-xs bg-red-100 text-red-700 hover:bg-red-200 dark:bg-red-900 dark:text-red-300 dark:hover:bg-red-800 transition-colors"
        >
          <X size={12} />
          {t('toolCall.reject')}
        </button>
        <button
          onClick={() => onConfirm(feedback)}
          className="flex items-center gap-1 px-2 py-1 rounded text-xs bg-green-100 text-green-700 hover:bg-green-200 dark:bg-green-900 dark:text-green-300 dark:hover:bg-green-800 transition-colors"
        >
          <Check size={12} />
          {t('toolCall.confirm')}
        </button>
      </div>
    </div>
  )
}
```

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/chat/ToolCallBlock.tsx frontend/src/components/chat/ToolConfirmCard.tsx
git commit -m "feat(frontend): add ToolCallBlock and ToolConfirmCard components"
```

---

### Task 17: Frontend — chatStore refactor

**Files:**
- Modify: `frontend/src/store/chatStore.ts`

- [ ] **Step 1: Rewrite chatStore to use backend bindings + Wails events**

```typescript
// frontend/src/store/chatStore.ts
import { create } from 'zustand'
import { Events } from '@wailsio/runtime'
import { Agent as AgentBinding } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/agent'
import type {
  Conversation,
  Message,
  ConversationStatus,
  StreamChunkEvent,
  ToolCallEvent,
  ConfirmRequestEvent,
  StreamDoneEvent,
  StreamErrorEvent,
  SessionStatusEvent,
} from '../types'

type StreamingMessage = {
  content: string
  thinking: string
}

type PendingConfirm = {
  requestId: string
  toolName: string
  args: string
  purpose: string
}

type ChatStore = {
  conversations: Conversation[]
  messages: Record<number, Message[]>
  currentConversationId: number | null
  streamingMessages: Record<number, StreamingMessage>
  pendingConfirms: Record<number, PendingConfirm | null>
  sessionsLoading: boolean
  hasMoreSessions: boolean
  sessionsCursor: number

  loadSessions: (reset?: boolean) => Promise<void>
  loadMessages: (sessionId: number) => Promise<void>
  createConversation: () => Promise<number>
  deleteConversation: (id: number) => Promise<void>
  renameConversation: (id: number, title: string) => Promise<void>
  toggleStar: (id: number, starred: boolean) => Promise<void>
  setCurrentConversation: (id: number | null) => void
  sendMessage: (params: {
    sessionId: number
    content: string
    baseUrl: string
    apiKey: string
    modelName: string
    enabledUserTools: string[]
  }) => Promise<void>
  stopGeneration: (sessionId: number) => Promise<void>
  respondToConfirm: (sessionId: number, approved: boolean, message: string) => Promise<void>
  generateTitle: (sessionId: number, baseUrl: string, apiKey: string, modelName: string) => Promise<void>
  initEventListeners: () => () => void
}

export const useChatStore = create<ChatStore>()((set, get) => ({
  conversations: [],
  messages: {},
  currentConversationId: null,
  streamingMessages: {},
  pendingConfirms: {},
  sessionsLoading: false,
  hasMoreSessions: true,
  sessionsCursor: 0,

  loadSessions: async (reset = false) => {
    if (get().sessionsLoading) return
    set({ sessionsLoading: true })
    try {
      const cursor = reset ? 0 : get().sessionsCursor
      const result = await AgentBinding.ListSessions({
        cursor,
        limit: 20,
        starred_only: false,
      })
      if (!result) return
      const newConvs: Conversation[] = (result.sessions ?? []).map(s => ({
        id: s.id,
        title: s.title,
        createdAt: s.created,
        updatedAt: s.updated,
        starred: s.starred,
        status: s.status as ConversationStatus,
      }))
      set(state => ({
        conversations: reset ? newConvs : [...state.conversations, ...newConvs],
        hasMoreSessions: result.has_more,
        sessionsCursor: result.next_cursor,
      }))
    } finally {
      set({ sessionsLoading: false })
    }
  },

  loadMessages: async (sessionId) => {
    const result = await AgentBinding.LoadSessionMessages({
      session_id: sessionId,
      offset: 0,
      limit: 100,
    })
    if (!result) return
    const msgs: Message[] = (result.messages ?? []).map(m => ({
      id: m.id,
      sessionId: m.session_id,
      parentId: m.parent_id ?? null,
      role: m.role as Message['role'],
      contentType: m.content_type as Message['contentType'],
      content: m.content,
      modelName: m.model_name,
      agentName: m.agent_name,
      tokensIn: m.tokens_in,
      tokensOut: m.tokens_out,
      extra: m.extra,
      createdAt: m.created,
    }))
    set(state => ({
      messages: { ...state.messages, [sessionId]: msgs },
    }))
  },

  createConversation: async () => {
    const result = await AgentBinding.CreateSession({ title: '' })
    if (!result) return 0
    const conv: Conversation = {
      id: result.session_id,
      title: result.title,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString(),
      starred: false,
      status: 'idle',
    }
    set(state => ({ conversations: [conv, ...state.conversations] }))
    return result.session_id
  },

  deleteConversation: async (id) => {
    await AgentBinding.DeleteSession({ session_id: id })
    set(state => ({
      conversations: state.conversations.filter(c => c.id !== id),
      currentConversationId: state.currentConversationId === id ? null : state.currentConversationId,
    }))
  },

  renameConversation: async (id, title) => {
    await AgentBinding.RenameSession({ session_id: id, title })
    set(state => ({
      conversations: state.conversations.map(c =>
        c.id === id ? { ...c, title } : c
      ),
    }))
  },

  toggleStar: async (id, starred) => {
    await AgentBinding.ToggleStarSession({ session_id: id, starred })
    set(state => ({
      conversations: state.conversations.map(c =>
        c.id === id ? { ...c, starred } : c
      ),
    }))
  },

  setCurrentConversation: (id) => {
    if (id) {
      set(state => ({
        conversations: state.conversations.map(c =>
          c.id === id && c.status !== 'idle' && c.status !== 'loading'
            ? { ...c, status: 'idle' }
            : c
        ),
      }))
      get().loadMessages(id)
    }
    set({ currentConversationId: id })
  },

  sendMessage: async (params) => {
    set(state => ({
      streamingMessages: {
        ...state.streamingMessages,
        [params.sessionId]: { content: '', thinking: '' },
      },
    }))
    await AgentBinding.SendMessage({
      session_id: params.sessionId,
      content: params.content,
      base_url: params.baseUrl,
      api_key: params.apiKey,
      model_name: params.modelName,
      enabled_user_tools: params.enabledUserTools,
    })
  },

  stopGeneration: async (sessionId) => {
    await AgentBinding.StopGeneration({ session_id: sessionId })
  },

  respondToConfirm: async (sessionId, approved, message) => {
    await AgentBinding.RespondToConfirm({
      session_id: sessionId,
      approved,
      message,
    })
    set(state => ({
      pendingConfirms: { ...state.pendingConfirms, [sessionId]: null },
    }))
  },

  generateTitle: async (sessionId, baseUrl, apiKey, modelName) => {
    const result = await AgentBinding.GenerateTitle({
      session_id: sessionId,
      base_url: baseUrl,
      api_key: apiKey,
      model_name: modelName,
    })
    if (result?.title) {
      set(state => ({
        conversations: state.conversations.map(c =>
          c.id === sessionId ? { ...c, title: result.title } : c
        ),
      }))
    }
  },

  initEventListeners: () => {
    const offs: Array<() => void> = []

    offs.push(Events.On('agent:stream:chunk', (event: { data: StreamChunkEvent[] }) => {
      const data = event.data[0]
      if (!data) return
      set(state => {
        const prev = state.streamingMessages[data.sessionId] ?? { content: '', thinking: '' }
        const key = data.contentType === 'thinking' ? 'thinking' : 'content'
        return {
          streamingMessages: {
            ...state.streamingMessages,
            [data.sessionId]: { ...prev, [key]: prev[key] + data.delta },
          },
        }
      })
    }))

    offs.push(Events.On('agent:stream:confirm_request', (event: { data: ConfirmRequestEvent[] }) => {
      const data = event.data[0]
      if (!data) return
      set(state => ({
        pendingConfirms: {
          ...state.pendingConfirms,
          [data.sessionId]: {
            requestId: data.requestId,
            toolName: data.toolName,
            args: data.args,
            purpose: data.purpose,
          },
        },
      }))
    }))

    offs.push(Events.On('agent:stream:done', (event: { data: StreamDoneEvent[] }) => {
      const data = event.data[0]
      if (!data) return
      set(state => {
        const updated = { ...state.streamingMessages }
        delete updated[data.sessionId]
        return { streamingMessages: updated }
      })
      get().loadMessages(data.sessionId)
    }))

    offs.push(Events.On('agent:stream:error', (event: { data: StreamErrorEvent[] }) => {
      const data = event.data[0]
      if (!data) return
      set(state => {
        const updated = { ...state.streamingMessages }
        delete updated[data.sessionId]
        return { streamingMessages: updated }
      })
    }))

    offs.push(Events.On('agent:session:status', (event: { data: SessionStatusEvent[] }) => {
      const data = event.data[0]
      if (!data) return
      set(state => ({
        conversations: state.conversations.map(c =>
          c.id === data.sessionId ? { ...c, status: data.status } : c
        ),
      }))
    }))

    return () => offs.forEach(off => off())
  },
}))
```

- [ ] **Step 2: Verify frontend compiles**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop/frontend && npx tsc --noEmit
```

Note: This may have type errors until Wails bindings are regenerated. The bindings will be generated after `go build` or `wails dev`. Proceed with remaining frontend tasks and do a full check at the end.

- [ ] **Step 3: Commit**

```bash
git add frontend/src/store/chatStore.ts
git commit -m "feat(frontend): rewrite chatStore with backend integration and Wails events"
```

---

### Task 18: Frontend — MessageItem refactor

**Files:**
- Modify: `frontend/src/components/chat/MessageItem.tsx`

- [ ] **Step 1: Refactor MessageItem to support all content types**

```tsx
// frontend/src/components/chat/MessageItem.tsx
import ReactMarkdown from 'react-markdown'
import remarkGfm from 'remark-gfm'
import rehypeHighlight from 'rehype-highlight'
import 'highlight.js/styles/github-dark-dimmed.min.css'
import { ArrowUpFromLine, ArrowDownToLine } from 'lucide-react'
import { ThinkingBlock } from './ThinkingBlock'
import { ToolCallBlock } from './ToolCallBlock'
import { ToolConfirmCard } from './ToolConfirmCard'
import type { Message } from '@/types'

type Props = {
  message: Message
  isStreaming?: boolean
  streamingContent?: string
  streamingThinking?: string
  pendingConfirm?: {
    requestId: string
    toolName: string
    args: string
    purpose: string
  } | null
  onConfirm?: (message: string) => void
  onReject?: (message: string) => void
}

export function MessageItem({
  message,
  isStreaming,
  streamingContent,
  streamingThinking,
  pendingConfirm,
  onConfirm,
  onReject,
}: Props) {
  const { role, contentType, content, modelName, tokensIn, tokensOut, extra } = message

  if (role === 'user' && contentType === 'text') {
    return (
      <div className="flex justify-end">
        <div className="max-w-[70%] select-text rounded-2xl rounded-br-sm bg-primary text-primary-foreground px-4 py-2.5 text-sm">
          {content}
        </div>
      </div>
    )
  }

  if (role === 'user' && contentType === 'confirm_response') {
    const approved = content.startsWith('approved')
    return (
      <ToolConfirmCard
        toolName=""
        purpose=""
        args=""
        onConfirm={() => {}}
        onReject={() => {}}
        resolved={true}
        approved={approved}
      />
    )
  }

  if (contentType === 'tool_call') {
    let toolName = ''
    let args = content
    try {
      const parsed = JSON.parse(content)
      toolName = parsed.name ?? parsed.ToolName ?? ''
      args = JSON.stringify(parsed.args ?? parsed, null, 2)
    } catch { /* use raw content */ }
    return (
      <ToolCallBlock
        toolName={toolName}
        purpose={extra}
        args={args}
        result=""
        status="completed"
      />
    )
  }

  if (contentType === 'tool_result') {
    return (
      <ToolCallBlock
        toolName=""
        purpose=""
        args=""
        result={content}
        status="completed"
      />
    )
  }

  if (contentType === 'thinking') {
    return <ThinkingBlock content={content} defaultExpanded={false} />
  }

  // Assistant text message — the main case.
  const displayContent = isStreaming && streamingContent !== undefined ? streamingContent : content
  const displayThinking = isStreaming && streamingThinking !== undefined ? streamingThinking : undefined

  return (
    <div className="flex flex-col gap-1">
      {displayThinking && (
        <ThinkingBlock content={displayThinking} defaultExpanded={isStreaming && !displayContent} />
      )}

      {displayContent && (
        <div className="select-text text-sm text-foreground leading-relaxed">
          <ReactMarkdown
            remarkPlugins={[remarkGfm]}
            rehypePlugins={[rehypeHighlight]}
            components={{
              pre: ({ children }) => (
                <pre className="rounded-lg overflow-x-auto my-3">{children}</pre>
              ),
              code: ({ children, className }) => {
                const isBlock = className?.includes('language-')
                return isBlock ? (
                  <code className={className}>{children}</code>
                ) : (
                  <code className="bg-muted px-1 py-0.5 rounded text-xs font-mono">{children}</code>
                )
              },
            }}
          >
            {displayContent}
          </ReactMarkdown>
        </div>
      )}

      {pendingConfirm && onConfirm && onReject && (
        <ToolConfirmCard
          toolName={pendingConfirm.toolName}
          purpose={pendingConfirm.purpose}
          args={pendingConfirm.args}
          onConfirm={onConfirm}
          onReject={onReject}
          resolved={false}
          approved={false}
        />
      )}

      {!isStreaming && (modelName || tokensIn > 0 || tokensOut > 0) && (
        <div className="flex items-center gap-3 text-xs text-muted-foreground mt-1">
          {modelName && <span className="opacity-70">{modelName}</span>}
          {tokensIn > 0 && (
            <span className="flex items-center gap-1">
              <ArrowUpFromLine size={10} />
              {tokensIn}
            </span>
          )}
          {tokensOut > 0 && (
            <span className="flex items-center gap-1">
              <ArrowDownToLine size={10} />
              {tokensOut}
            </span>
          )}
        </div>
      )}
    </div>
  )
}
```

- [ ] **Step 2: Commit**

```bash
git add frontend/src/components/chat/MessageItem.tsx
git commit -m "feat(frontend): refactor MessageItem to support all message content types"
```

---

### Task 19: Frontend — ChatMessages and ChatInput refactor

**Files:**
- Modify: `frontend/src/components/chat/ChatMessages.tsx`
- Modify: `frontend/src/components/chat/ChatInput.tsx`

- [ ] **Step 1: Refactor ChatMessages to use backend data + streaming state**

```tsx
// frontend/src/components/chat/ChatMessages.tsx
import { useEffect } from 'react'
import { ChevronDown, Loader2 } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { cn } from '@/lib/utils'
import { shouldShowTimestamp, formatTimestamp } from '@/lib/utils'
import { useAutoScroll } from '@/hooks/useAutoScroll'
import { useAppStore } from '@/store/appStore'
import { useChatStore } from '@/store/chatStore'
import { MessageItem } from './MessageItem'
import { WelcomeScreen } from './WelcomeScreen'

export function ChatMessages() {
  const { t } = useTranslation()
  const language = useAppStore(s => s.language)
  const {
    currentConversationId,
    messages,
    conversations,
    streamingMessages,
    pendingConfirms,
    respondToConfirm,
  } = useChatStore()

  const currentMessages = currentConversationId
    ? (messages[currentConversationId] ?? [])
    : []

  const currentConv = conversations.find(c => c.id === currentConversationId)
  const isStreaming = currentConv?.status === 'loading'
  const streamingMsg = currentConversationId
    ? streamingMessages[currentConversationId]
    : undefined
  const pendingConfirm = currentConversationId
    ? pendingConfirms[currentConversationId]
    : undefined

  const { containerRef, isAtBottom, scrollToBottom } = useAutoScroll([
    currentMessages.length,
    isStreaming,
    streamingMsg?.content,
    streamingMsg?.thinking,
  ])

  const isNewChat = currentMessages.length === 0 && !isStreaming

  useEffect(() => {
    const frame = requestAnimationFrame(() => scrollToBottom(false))
    return () => cancelAnimationFrame(frame)
  }, [currentConversationId, scrollToBottom])

  const handleConfirm = (message: string) => {
    if (currentConversationId) {
      respondToConfirm(currentConversationId, true, message)
    }
  }

  const handleReject = (message: string) => {
    if (currentConversationId) {
      respondToConfirm(currentConversationId, false, message)
    }
  }

  return (
    <div className="relative flex flex-col flex-1 min-h-0">
      <div
        ref={containerRef}
        className="chat-scroll-area flex-1 overflow-y-auto py-4 scroll-smooth"
      >
        {isNewChat ? (
          <div className="mx-auto flex min-h-full w-full max-w-5xl items-center justify-center px-4">
            <WelcomeScreen />
          </div>
        ) : (
          <div className="mx-auto flex w-full max-w-5xl flex-col gap-4 px-4">
            {currentMessages.map((msg, index) => (
              <div key={msg.id}>
                {shouldShowTimestamp(currentMessages, index) && (
                  <div className="text-center py-2">
                    <span className="text-xs text-muted-foreground px-2">
                      {formatTimestamp(msg.createdAt, language)}
                    </span>
                  </div>
                )}
                <MessageItem message={msg} />
              </div>
            ))}

            {isStreaming && streamingMsg && (
              <MessageItem
                message={{
                  id: 0,
                  sessionId: currentConversationId ?? 0,
                  parentId: null,
                  role: 'assistant',
                  contentType: 'text',
                  content: '',
                  modelName: '',
                  agentName: '',
                  tokensIn: 0,
                  tokensOut: 0,
                  extra: '',
                  createdAt: new Date().toISOString(),
                }}
                isStreaming={true}
                streamingContent={streamingMsg.content}
                streamingThinking={streamingMsg.thinking}
                pendingConfirm={pendingConfirm}
                onConfirm={handleConfirm}
                onReject={handleReject}
              />
            )}

            {isStreaming && !streamingMsg?.content && !streamingMsg?.thinking && (
              <div className="flex items-center gap-2 px-4 text-muted-foreground">
                <Loader2 size={14} className="animate-spin" />
              </div>
            )}
          </div>
        )}
      </div>

      {!isAtBottom && (
        <button
          onClick={() => scrollToBottom(true)}
          className={cn(
            'absolute bottom-4 left-1/2 -translate-x-1/2 z-10',
            'flex items-center gap-1.5 px-3 py-1.5 rounded-full',
            'bg-background border border-border shadow-md',
            'text-xs text-foreground hover:bg-accent transition-colors'
          )}
        >
          <ChevronDown size={14} />
          {t('chat.scrollToBottom')}
        </button>
      )}
    </div>
  )
}
```

- [ ] **Step 2: Update ChatInput — wire real send/stop**

The `handleSend` and `handleStop` functions in `ChatInput.tsx` need to call the backend. The key changes to `frontend/src/components/chat/ChatInput.tsx`:

Replace the `handleSend` function body:

```typescript
  const handleSend = async () => {
    if (!editor || editor.isEmpty || isStreaming) return
    if (!currentConversationId || !selectedModelId) return

    const md = editor.storage.markdown.getMarkdown()
    editor.commands.clearContent()
    setIsEmpty(true)

    const selectedProvider = enabledProviders.find(p =>
      p.models.some(m => m.id === selectedModelId)
    )
    if (!selectedProvider) return

    const pm = providerModels.find(pm => pm.provider.id === selectedProvider.id)
    if (!pm) return

    const enabledUserToolIds = tools.filter(t => t.enabled).map(t => t.id)

    await sendMessage({
      sessionId: currentConversationId,
      content: md,
      baseUrl: pm.provider.base_url,
      apiKey: pm.provider.api_key,
      modelName: providerModels
        .flatMap(pm => pm.models)
        .find(m => m.id === selectedModelId)?.model ?? '',
      enabledUserTools: enabledUserToolIds,
    })
  }
```

Replace the `handleStop` function body:

```typescript
  const handleStop = async () => {
    if (currentConversationId) {
      await stopGeneration(currentConversationId)
    }
  }
```

Add to the destructured imports from `useChatStore`:

```typescript
  const { conversations, currentConversationId, setConversationStatus, sendMessage, stopGeneration } = useChatStore()
```

Remove the `setConversationStatus` usage in `handleStop` (backend handles it now).

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/chat/ChatMessages.tsx frontend/src/components/chat/ChatInput.tsx
git commit -m "feat(frontend): wire ChatMessages and ChatInput to backend"
```

---

### Task 20: Frontend — ConversationList refactor and event init

**Files:**
- Modify: `frontend/src/components/sidebar/ConversationList.tsx`
- Modify: `frontend/src/components/sidebar/ConversationItem.tsx`
- Modify: `frontend/src/components/layout/MainLayout.tsx`

- [ ] **Step 1: Update ConversationList to load from backend**

Key changes to `ConversationList.tsx`:

Replace the destructured store values:

```typescript
  const {
    conversations,
    currentConversationId,
    setCurrentConversation,
    createConversation,
    loadSessions,
    hasMoreSessions,
    sessionsLoading,
  } = useChatStore()
```

Add `useEffect` to load sessions on mount:

```typescript
  useEffect(() => {
    void loadSessions(true)
  }, [loadSessions])
```

Update `handleNewChat`:

```typescript
  const handleNewChat = async () => {
    const id = await createConversation()
    if (id) setCurrentConversation(id)
  }
```

Add a "load more" trigger at the bottom of the list (replace the existing footer):

```tsx
  {hasMoreSessions && (
    <button
      onClick={() => void loadSessions()}
      disabled={sessionsLoading}
      className="w-full py-2 text-center text-xs text-muted-foreground hover:text-foreground transition-colors"
    >
      {sessionsLoading ? t('sidebar.loading') : t('sidebar.loadMore')}
    </button>
  )}
  {!hasMoreSessions && (
    <div className="py-3 text-center text-xs text-muted-foreground">
      {t('sidebar.loadedAll', { count: filtered.length })}
    </div>
  )}
```

- [ ] **Step 2: Update ConversationItem to use async store methods**

In `ConversationItem.tsx`, update the handler calls:

```typescript
  const { deleteConversation, renameConversation, toggleStar } = useChatStore()

  const handleRenameSubmit = async () => {
    if (renameValue.trim()) await renameConversation(id, renameValue.trim())
    setIsRenaming(false)
  }
```

Update the toggleStar call to pass the starred flag:

```tsx
  onClick={async (e) => {
    e.stopPropagation()
    await toggleStar(id, !starred)
    setMenuOpen(false)
  }}
```

Update the delete call:

```tsx
  onClick={async (e) => {
    e.stopPropagation()
    await deleteConversation(id)
    setMenuOpen(false)
  }}
```

- [ ] **Step 3: Initialize event listeners in MainLayout**

Add to `MainLayout.tsx`:

```typescript
import { useChatStore } from '@/store/chatStore'
```

Add `useEffect` inside the `MainLayout` component:

```typescript
  useEffect(() => {
    const cleanup = useChatStore.getState().initEventListeners()
    return cleanup
  }, [])
```

- [ ] **Step 4: Add new i18n keys for loading**

In `zh-CN.ts` sidebar section, add:

```typescript
  loadMore: '加载更多',
  loading: '加载中...',
```

In `en.ts` sidebar section, add:

```typescript
  loadMore: 'Load more',
  loading: 'Loading...',
```

- [ ] **Step 5: Commit**

```bash
git add frontend/src/components/sidebar/ConversationList.tsx frontend/src/components/sidebar/ConversationItem.tsx frontend/src/components/layout/MainLayout.tsx frontend/src/i18n/locales/zh-CN.ts frontend/src/i18n/locales/en.ts
git commit -m "feat(frontend): wire sidebar to backend with pagination and event listeners"
```

---

### Task 21: Generate Wails bindings and full build verification

**Files:**
- Generated: `frontend/bindings/...`

- [ ] **Step 1: Run go mod tidy**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop && go mod tidy
```

- [ ] **Step 2: Generate Wails bindings**

```bash
wails3 generate bindings
```

Or if using the dev command:

```bash
wails3 dev
```

This generates TypeScript bindings for the new `Agent` service under `frontend/bindings/`.

- [ ] **Step 3: Fix any TypeScript errors**

```bash
cd frontend && npx tsc --noEmit
```

Fix any import path or type mismatches based on the generated bindings. The DTO field names in TypeScript will use snake_case (matching the Go json tags).

- [ ] **Step 4: Run the full app**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop && wails3 dev
```

Verify:
1. App starts without errors
2. Can create a new conversation
3. Conversation list loads from backend
4. Can rename/delete/star conversations

- [ ] **Step 5: Commit generated bindings and any fixes**

```bash
git add -A
git commit -m "feat: generate Wails bindings and fix integration issues"
```

---

### Task 22: End-to-end smoke test

- [ ] **Step 1: Test conversation creation**

1. Open the app
2. Click "New Chat" in the sidebar
3. Verify a new conversation appears in the list
4. Verify it becomes the active conversation

- [ ] **Step 2: Test sending a message**

1. Select a model from the model selector
2. Type a message and press Enter
3. Verify the message appears as a user bubble
4. Verify streaming response appears in real-time
5. Verify thinking block shows if model supports it
6. Verify token counts and model name show after completion

- [ ] **Step 3: Test stop generation**

1. Send a long message that will generate a long response
2. Click the stop button during streaming
3. Verify generation stops and status returns to idle

- [ ] **Step 4: Test session management**

1. Create multiple conversations
2. Switch between them — verify messages load correctly
3. Rename a conversation via the context menu
4. Star/unstar a conversation
5. Switch to Favorites tab — verify only starred conversations show
6. Delete a conversation

- [ ] **Step 5: Test background running**

1. Start a generation in one conversation
2. Switch to another conversation while it's running
3. Verify the loading indicator shows in the sidebar for the running conversation
4. Wait for completion — verify the green dot appears
5. Click back — verify messages loaded and dot disappears (marked as read)

- [ ] **Step 6: Test auto-generate title**

1. Create a new conversation and send a message
2. After the response completes, verify the title was auto-generated

- [ ] **Step 7: Commit any fixes from testing**

```bash
git add -A
git commit -m "fix: address issues found during smoke testing"
```
