# Hidden Sessions + Notifications Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the "background AI task" framework described in spec `docs/dev/13.plugin_cli.md` §3.5 — a hidden chat-session model, a global notification system, and a popup chat window for user intervention. This is Plan B of three; Plan C (CLI Refactor) consumes it. Plan A (Skills Foundation) must merge first.

**Architecture:**
- Add a `kind` column to `Session` (default `"user"`, hidden sessions store `"background"`). Hidden sessions reuse the existing chat / agent pipeline 1:1 (same SendMessage, same streaming events, same tool dispatch) — they are simply filtered out of the normal chat list and rendered through a different surface.
- A new service `backend/service/notification/` owns lightweight notifications: `Notify`, `ListNotifications`, `DismissNotification`, `ResolveNotification`. Storage is a new SQLite table (`notifications`) so notifications survive restarts. Each notification carries `session_id`, `kind` (`needs_attention` / `info`), `title`, `message`, `created_at`, `resolved_at`.
- A new built-in tool `RequestUserAttention(title, message)` exposed to the agent. When called from inside a hidden session it (a) inserts a `tool_result` placeholder noting "awaiting user", (b) calls `Notification.Notify`, (c) emits the global notification event, and (d) **pauses the agent loop** — the agent does not continue until a user reply lands on the session (via the popup).
- Frontend grows three pieces:
  - `useNotificationsStore` — global Zustand store fed by Wails events and an initial fetch.
  - `NotificationBell` — a tray button (mounted in the chat / main window chrome) showing unresolved count, opening a dropdown.
  - `FloatingChatPanel` — a draggable / dismissable modal that renders a single session's transcript and input, reusing the chat-message components. Opens when the user clicks a notification.
- Hidden sessions are excluded from `ListSessions` by default. A new flag `include_hidden` is added but only used internally / by debug surfaces.
- The "AI-generated skill must be confirmed before persist" requirement from Plan A spec §3.4.3 lands here: a new agent tool `ProposeSkill(name, description, body)` produces a notification with a confirm/reject popup; only on confirm does it call the existing `Skills.CreateSkill` with `source: "ai"`.

**Tech Stack:** Go 1.22+, GORM, Wails v3, React, Zustand, `react-rnd` (for the floating window — optional; can be replaced with plain absolute positioning if you prefer no new dep).

---

## Pre-flight

- [ ] **Confirm Plan A is merged.** Run `git log --oneline | grep -i "skills foundation\|register skills"` — at least one Plan A commit should be present.
- [ ] **Decide draggable dep.** If introducing `react-rnd` is undesirable, skip it: the plan's Task 11 includes both variants.

---

## Section 1 — Hidden session storage + spawn API

### Task 1: Schema: add `kind` to `Session`

**Files:**
- Modify: `backend/models/data_models/session.go`
- Modify: `backend/storage/storage.go` (AutoMigrate is presumably here)

- [ ] **Step 1: Write failing test**

Create `backend/storage/session_test.go` (or append):

```go
package storage

import (
	"testing"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
)

func TestSession_KindFieldPersists(t *testing.T) {
	stor := newTestStorage(t) // helper that opens an in-memory sqlite + AutoMigrate
	created, err := stor.CreateSession(data_models.Session{Title: "x", Kind: "background"})
	if err != nil {
		t.Fatal(err)
	}
	got, err := stor.GetSession(created.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Kind != "background" {
		t.Fatalf("kind not persisted: %q", got.Kind)
	}
}
```

If `newTestStorage` doesn't exist, add it (uses `sqlite::memory:`).

- [ ] **Step 2: Run — must fail**

```
go test ./backend/storage/...
```

- [ ] **Step 3: Add field**

In `backend/models/data_models/session.go`:

```go
type Session struct {
	OrmModel
	Title   string `gorm:"type:varchar(255)" json:"title"`
	Starred bool   `gorm:"type:bool;default:0;index" json:"starred"`
	Status  string `gorm:"type:varchar(50);default:'idle'" json:"status"`
	Kind    string `gorm:"type:varchar(20);default:'user';index" json:"kind"`
}
```

- [ ] **Step 4: Run — must pass**

```
go test ./backend/storage/...
```

GORM auto-migration on next app start adds the column; if `storage.go` calls `AutoMigrate(&data_models.Session{})` you're done. If not, locate the migration setup and ensure Session is listed.

- [ ] **Step 5: Commit**

```bash
git add backend/models/data_models/session.go backend/storage/
git commit -m "feat(session): add kind column for hidden background sessions"
```

---

### Task 2: Filter hidden sessions from `ListSessions`

**Files:**
- Modify: `backend/storage/session.go`

- [ ] **Step 1: Extend ListSessions signature**

Change to accept `includeHidden bool`:

```go
func (s *Storage) ListSessions(cursor uint, limit int, starredOnly bool, includeHidden bool) ([]data_models.Session, error) {
	var sessions []data_models.Session
	q := s.sqliteDB.Order("updated_at DESC")
	if cursor > 0 {
		q = q.Where("id < ?", cursor)
	}
	if starredOnly {
		q = q.Where("starred = ?", true)
	}
	if !includeHidden {
		q = q.Where("kind = ? OR kind = ''", "user")
	}
	if err := q.Limit(limit).Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}
```

- [ ] **Step 2: Update all callers**

```
grep -nR "ListSessions(" backend/ | grep -v _test
```

Pass `false` (excluding hidden) at every existing call site.

- [ ] **Step 3: Add test for filter**

In `session_test.go`:

```go
func TestListSessions_HidesBackgroundByDefault(t *testing.T) {
	stor := newTestStorage(t)
	_, _ = stor.CreateSession(data_models.Session{Title: "user-1"})
	_, _ = stor.CreateSession(data_models.Session{Title: "bg-1", Kind: "background"})
	listed, err := stor.ListSessions(0, 50, false, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(listed) != 1 || listed[0].Title != "user-1" {
		t.Fatalf("expected only user-1, got %+v", listed)
	}
	all, _ := stor.ListSessions(0, 50, false, true)
	if len(all) != 2 {
		t.Fatalf("expected 2 with include, got %d", len(all))
	}
}
```

- [ ] **Step 4: Run all backend tests**

```
go test ./...
```

- [ ] **Step 5: Commit**

```bash
git add backend/storage/session.go backend/storage/session_test.go backend/service/ backend/pkg/
git commit -m "feat(session): exclude hidden sessions from default list"
```

---

### Task 3: `SpawnHiddenSession` Wails method

**Files:**
- Create: `backend/service/agent/agent_dto/spawn_hidden_session.go`
- Modify: `backend/service/agent/agent.go`

- [ ] **Step 1: Create DTO**

```go
package agent_dto

type SpawnHiddenSessionInput struct {
	Title          string `json:"title"`
	SystemPrompt   string `json:"system_prompt"`
	UserMessage    string `json:"user_message"`
	SkillName      string `json:"skill_name"` // optional: pre-arms a skill (Plan A's Skill tool sees it)
}

type SpawnHiddenSessionOutput struct {
	SessionID uint   `json:"session_id"`
	Title     string `json:"title"`
}
```

- [ ] **Step 2: Add Wails method**

In `agent.go`, append after `SendMessage`:

```go
// SpawnHiddenSession creates a kind="background" session and immediately submits a first user message.
// The agent stream proceeds in a goroutine; callers get back the session id for progress subscriptions.
func (a *Agent) SpawnHiddenSession(ctx context.Context, input agent_dto.SpawnHiddenSessionInput) (*agent_dto.SpawnHiddenSessionOutput, error) {
	handler, err := a.handlerFor(ctx)
	if err != nil {
		return nil, ierror.Error(ierror.ErrAgentCreateSession, err)
	}
	session, err := handler.CreateSession(ctx, pkgAgent.CreateSessionParams{
		Title: input.Title,
		Kind:  "background",
	})
	if err != nil {
		return nil, ierror.Error(ierror.ErrAgentCreateSession, err)
	}
	go func() {
		_ = handler.SendMessage(context.Background(), pkgAgent.SendMessageParams{
			SessionID:    session.ID,
			Text:         input.UserMessage,
			SystemPrompt: input.SystemPrompt,
			PrimedSkill:  input.SkillName,
		})
	}()
	return &agent_dto.SpawnHiddenSessionOutput{SessionID: session.ID, Title: session.Title}, nil
}
```

(If `pkgAgent.CreateSessionParams.Kind` or `SendMessageParams.SystemPrompt`/`PrimedSkill` don't exist, add them — they're already required by spec §3.5. Update `pkg/agent/manager.go` accordingly.)

- [ ] **Step 3: Test spawn returns id**

In `agent_test.go`, add (with stubbed handler):

```go
func TestSpawnHiddenSession_CreatesBackgroundSession(t *testing.T) {
	a := newTestAgent(t)
	out, err := a.SpawnHiddenSession(context.Background(), agent_dto.SpawnHiddenSessionInput{
		Title: "bg test", UserMessage: "hi",
	})
	if err != nil {
		t.Fatal(err)
	}
	if out.SessionID == 0 {
		t.Fatal("expected non-zero session id")
	}
}
```

- [ ] **Step 4: Run tests + regenerate bindings**

```
go test ./...
wails generate bindings
```

- [ ] **Step 5: Commit**

```bash
git add backend/service/agent/ backend/pkg/agent/ frontend/bindings/
git commit -m "feat(agent): add SpawnHiddenSession method"
```

---

## Section 2 — Notification service

### Task 4: Storage table + model

**Files:**
- Create: `backend/models/data_models/notification.go`
- Create: `backend/storage/notification.go`
- Modify: `backend/storage/storage.go` (AutoMigrate)

- [ ] **Step 1: Define model**

```go
package data_models

import "time"

type Notification struct {
	OrmModel
	SessionID  uint       `gorm:"index" json:"session_id"`
	Kind       string     `gorm:"type:varchar(40);index" json:"kind"` // "needs_attention" | "info"
	Title      string     `gorm:"type:varchar(255)" json:"title"`
	Message    string     `gorm:"type:text" json:"message"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
}
```

- [ ] **Step 2: Add storage methods**

```go
package storage

import (
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
)

func (s *Storage) CreateNotification(n data_models.Notification) (*data_models.Notification, error) {
	if err := s.sqliteDB.Create(&n).Error; err != nil {
		return nil, err
	}
	return &n, nil
}

func (s *Storage) ListNotifications(includeResolved bool) ([]data_models.Notification, error) {
	var items []data_models.Notification
	q := s.sqliteDB.Order("created_at DESC")
	if !includeResolved {
		q = q.Where("resolved_at IS NULL")
	}
	if err := q.Limit(200).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (s *Storage) ResolveNotification(id uint) error {
	now := time.Now()
	return s.sqliteDB.Model(&data_models.Notification{}).Where("id = ?", id).Update("resolved_at", &now).Error
}

func (s *Storage) DeleteNotification(id uint) error {
	return s.sqliteDB.Delete(&data_models.Notification{}, id).Error
}
```

- [ ] **Step 3: Register in AutoMigrate**

In `storage.go`, add `&data_models.Notification{}` to the AutoMigrate call.

- [ ] **Step 4: Write storage test**

```go
func TestNotification_CRUD(t *testing.T) {
	stor := newTestStorage(t)
	n, err := stor.CreateNotification(data_models.Notification{
		SessionID: 1, Kind: "needs_attention", Title: "T", Message: "M",
	})
	if err != nil { t.Fatal(err) }
	list, _ := stor.ListNotifications(false)
	if len(list) != 1 { t.Fatalf("want 1 got %d", len(list)) }
	if err := stor.ResolveNotification(n.ID); err != nil { t.Fatal(err) }
	list, _ = stor.ListNotifications(false)
	if len(list) != 0 { t.Fatalf("expected resolved out of list") }
}
```

- [ ] **Step 5: Run + commit**

```
go test ./backend/storage/...
git add backend/models/data_models/notification.go backend/storage/
git commit -m "feat(notification): add storage model and CRUD"
```

---

### Task 5: Notification Wails service

**Files:**
- Create: `backend/service/notification/notification.go`
- Create: `backend/service/notification/notification_implement.go`
- Create: `backend/service/notification/notification_internal.go`
- Create: `backend/service/notification/notification_dto/notify.go`
- Create: `backend/service/notification/notification_dto/list_notifications.go`
- Create: `backend/service/notification/notification_dto/resolve_notification.go`
- Create: `backend/service/notification/notification_dto/dismiss_notification.go`
- Create: `backend/service/notification/notification_test.go`
- Modify: `main.go` to register the service.
- Add ierror codes in `backend/pkg/ierror/error.go`:

```go
	ErrNotificationCreate errorCode = "ierror.notification.create"
	ErrNotificationList   errorCode = "ierror.notification.list"
	ErrNotificationResolve errorCode = "ierror.notification.resolve"
```

- [ ] **Step 1: Create DTOs**

`notify.go`:

```go
package notification_dto

type NotifyInput struct {
	SessionID uint   `json:"session_id"`
	Kind      string `json:"kind"`
	Title     string `json:"title"`
	Message   string `json:"message"`
}

type NotifyOutput struct {
	NotificationID uint `json:"notification_id"`
}
```

`list_notifications.go`:

```go
package notification_dto

import "time"

type ListNotificationsInput struct {
	IncludeResolved bool `json:"include_resolved"`
}

type NotificationItem struct {
	ID         uint       `json:"id"`
	SessionID  uint       `json:"session_id"`
	Kind       string     `json:"kind"`
	Title      string     `json:"title"`
	Message    string     `json:"message"`
	CreatedAt  time.Time  `json:"created_at"`
	ResolvedAt *time.Time `json:"resolved_at,omitempty"`
}

type ListNotificationsOutput struct {
	Items []NotificationItem `json:"items"`
}
```

`resolve_notification.go`:

```go
package notification_dto

type ResolveNotificationInput struct {
	ID uint `json:"id"`
}

type ResolveNotificationOutput struct{}
```

`dismiss_notification.go`:

```go
package notification_dto

type DismissNotificationInput struct {
	ID uint `json:"id"`
}

type DismissNotificationOutput struct{}
```

- [ ] **Step 2: Create `notification.go`**

```go
package notification

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/ierror"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/notification/notification_dto"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/storage"
)

const eventCreated = "notification.created"
const eventResolved = "notification.resolved"

type Notification struct {
	store    *storage.Storage
	wailsApp *application.App
}

func NewNotification(store *storage.Storage) *Notification {
	return &Notification{store: store}
}

// Notify creates a notification and emits a global event.
func (n *Notification) Notify(ctx context.Context, input notification_dto.NotifyInput) (*notification_dto.NotifyOutput, error) {
	created, err := n.store.CreateNotification(data_models.Notification{
		SessionID: input.SessionID,
		Kind:      input.Kind,
		Title:     input.Title,
		Message:   input.Message,
	})
	if err != nil {
		return nil, ierror.Error(ierror.ErrNotificationCreate, err)
	}
	if n.wailsApp != nil {
		n.wailsApp.Event.Emit(eventCreated, dtoFromModel(*created))
	}
	return &notification_dto.NotifyOutput{NotificationID: created.ID}, nil
}

func (n *Notification) ListNotifications(ctx context.Context, input notification_dto.ListNotificationsInput) (*notification_dto.ListNotificationsOutput, error) {
	items, err := n.store.ListNotifications(input.IncludeResolved)
	if err != nil {
		return nil, ierror.Error(ierror.ErrNotificationList, err)
	}
	out := notification_dto.ListNotificationsOutput{}
	for _, it := range items {
		out.Items = append(out.Items, dtoFromModel(it))
	}
	return &out, nil
}

func (n *Notification) ResolveNotification(ctx context.Context, input notification_dto.ResolveNotificationInput) (*notification_dto.ResolveNotificationOutput, error) {
	if err := n.store.ResolveNotification(input.ID); err != nil {
		return nil, ierror.Error(ierror.ErrNotificationResolve, err)
	}
	if n.wailsApp != nil {
		n.wailsApp.Event.Emit(eventResolved, map[string]uint{"id": input.ID})
	}
	return &notification_dto.ResolveNotificationOutput{}, nil
}

func (n *Notification) DismissNotification(ctx context.Context, input notification_dto.DismissNotificationInput) (*notification_dto.DismissNotificationOutput, error) {
	if err := n.store.DeleteNotification(input.ID); err != nil {
		return nil, ierror.Error(ierror.ErrNotificationResolve, err)
	}
	return &notification_dto.DismissNotificationOutput{}, nil
}

func dtoFromModel(m data_models.Notification) notification_dto.NotificationItem {
	return notification_dto.NotificationItem{
		ID:         m.ID,
		SessionID:  m.SessionID,
		Kind:       m.Kind,
		Title:      m.Title,
		Message:    m.Message,
		CreatedAt:  m.CreatedAt,
		ResolvedAt: m.ResolvedAt,
	}
}
```

- [ ] **Step 3: `notification_implement.go`**

```go
package notification

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
)

func (n *Notification) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {
	n.wailsApp = application.Get()
	return nil
}
```

- [ ] **Step 4: Register in main.go**

```go
import notificationSvc "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/notification"

// In services slice:
application.NewService(notificationSvc.NewNotification(istorage)),
```

- [ ] **Step 5: Test + commit**

```
go test ./backend/service/notification/...
wails generate bindings
git add backend/service/notification/ backend/pkg/ierror/error.go main.go frontend/bindings/
git commit -m "feat(notification): add Wails service for global notifications"
```

---

## Section 3 — `RequestUserAttention` agent tool

### Task 6: Built-in tool definition

**Files:**
- Create: `backend/pkg/agent/tools/request_attention_tool.go`
- Create: `backend/pkg/agent/tools/request_attention_tool_test.go`

- [ ] **Step 1: Write `request_attention_tool.go`**

```go
package tools

import (
	"context"
	"encoding/json"
	"errors"
)

// AttentionRequester is satisfied by the notification service (avoids an import cycle).
type AttentionRequester interface {
	NotifyAttention(ctx context.Context, sessionID uint, title, message string) (notificationID uint, err error)
	WaitForResolution(ctx context.Context, notificationID uint) error
}

const RequestAttentionToolName = "RequestUserAttention"

// BuildRequestAttentionTool returns a tool the agent registers per turn.
func BuildRequestAttentionTool() ToolMeta {
	return ToolMeta{
		Name:        RequestAttentionToolName,
		Description: "Pause the background task and ask the user. Call with {\"title\":\"short\",\"message\":\"what you need from them\"}. Returns when the user replies.",
		Category:    CategoryBuiltin,
		FormatPurpose: func(args json.RawMessage) string {
			var parsed struct {
				Title string `json:"title"`
			}
			_ = json.Unmarshal(args, &parsed)
			return "Asking user: " + parsed.Title
		},
	}
}

// InvokeRequestAttention is called by the agent's tool dispatcher.
// It pushes a notification and blocks until the user replies on the session (resolution arrives via session message stream).
func InvokeRequestAttention(ctx context.Context, requester AttentionRequester, sessionID uint, args json.RawMessage) (string, error) {
	var parsed struct {
		Title   string `json:"title"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(args, &parsed); err != nil {
		return "", err
	}
	if parsed.Title == "" {
		return "", errors.New("title is required")
	}
	notifID, err := requester.NotifyAttention(ctx, sessionID, parsed.Title, parsed.Message)
	if err != nil {
		return "", err
	}
	if err := requester.WaitForResolution(ctx, notifID); err != nil {
		return "", err
	}
	return "user replied; resume from their last message in this session", nil
}
```

- [ ] **Step 2: Write a fake requester test**

```go
package tools

import (
	"context"
	"encoding/json"
	"testing"
)

type fakeRequester struct{ id uint; waited bool }

func (f *fakeRequester) NotifyAttention(_ context.Context, _ uint, _, _ string) (uint, error) {
	f.id++
	return f.id, nil
}
func (f *fakeRequester) WaitForResolution(_ context.Context, _ uint) error {
	f.waited = true
	return nil
}

func TestInvokeRequestAttention_ReturnsAfterResolution(t *testing.T) {
	fr := &fakeRequester{}
	args := json.RawMessage(`{"title":"X","message":"Y"}`)
	out, err := InvokeRequestAttention(context.Background(), fr, 42, args)
	if err != nil { t.Fatal(err) }
	if !fr.waited { t.Fatal("did not wait") }
	if out == "" { t.Fatal("empty result") }
}
```

- [ ] **Step 3: Implement `AttentionRequester` on the notification service**

In `backend/service/notification/notification_internal.go`:

```go
package notification

import (
	"context"
	"errors"
	"sync"
	"time"
)

var (
	waitersMu sync.Mutex
	waiters   = map[uint]chan struct{}{}
)

func (n *Notification) NotifyAttention(ctx context.Context, sessionID uint, title, message string) (uint, error) {
	out, err := n.Notify(ctx, notification_dto.NotifyInput{
		SessionID: sessionID, Kind: "needs_attention", Title: title, Message: message,
	})
	if err != nil {
		return 0, err
	}
	waitersMu.Lock()
	waiters[out.NotificationID] = make(chan struct{}, 1)
	waitersMu.Unlock()
	return out.NotificationID, nil
}

func (n *Notification) WaitForResolution(ctx context.Context, notificationID uint) error {
	waitersMu.Lock()
	ch, ok := waiters[notificationID]
	waitersMu.Unlock()
	if !ok {
		return errors.New("no waiter registered")
	}
	select {
	case <-ch:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(30 * time.Minute):
		return errors.New("attention request timed out")
	}
}

// Called from ResolveNotification to wake any waiting goroutine.
func wakeWaiter(id uint) {
	waitersMu.Lock()
	ch, ok := waiters[id]
	if ok {
		delete(waiters, id)
	}
	waitersMu.Unlock()
	if ok {
		close(ch)
	}
}
```

In `notification.go`, modify `ResolveNotification` to call `wakeWaiter(input.ID)` after marking resolved.

- [ ] **Step 4: Wire tool into agent dispatch**

In the agent's tool registry / dispatcher (extends Plan A Task 14):

```go
case tools.RequestAttentionToolName:
    return tools.InvokeRequestAttention(ctx, a.notificationRequester, currentSessionID, call.Args)
```

Add `notificationRequester tools.AttentionRequester` field on the agent struct + setter (`SetAttentionRequester`). In `main.go`:

```go
notificationService := notificationSvc.NewNotification(istorage)
agentService.SetAttentionRequester(notificationService)
```

- [ ] **Step 5: Register the tool in `BuildSkillTool`'s sibling registration step**

Wherever Plan A registered `BuildSkillTool` per turn, also register `BuildRequestAttentionTool()` only when the current session's `Kind == "background"` (we don't want regular chat sessions to see this tool).

- [ ] **Step 6: Run tests**

```
go test ./...
```

- [ ] **Step 7: Commit**

```bash
git add backend/pkg/agent/tools/request_attention_tool.go backend/pkg/agent/tools/request_attention_tool_test.go backend/service/notification/notification_internal.go backend/service/notification/notification.go backend/pkg/agent/ main.go
git commit -m "feat(agent): add RequestUserAttention tool with resolution waiter"
```

---

### Task 7: `ProposeSkill` agent tool (AI-generated skill confirmation)

**Files:**
- Create: `backend/pkg/agent/tools/propose_skill_tool.go`
- Create: `backend/pkg/agent/tools/propose_skill_tool_test.go`

- [ ] **Step 1: Write the tool**

```go
package tools

import (
	"context"
	"encoding/json"
	"errors"

	pkgskills "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/skills"
)

// SkillCreator is satisfied by the skills service.
type SkillCreator interface {
	Create(skill pkgskills.Skill) (pkgskills.Skill, error)
}

// ConfirmationGate is satisfied by the notification service (NotifyAttention + WaitForResolution).
type ConfirmationGate = AttentionRequester

const ProposeSkillToolName = "ProposeSkill"

func BuildProposeSkillTool() ToolMeta {
	return ToolMeta{
		Name:        ProposeSkillToolName,
		Description: "Propose a new skill to persist. Call with {\"name\":\"kebab-case\",\"description\":\"...\",\"body\":\"markdown\"}. User must confirm before it lands on disk.",
		Category:    CategoryBuiltin,
	}
}

func InvokeProposeSkill(ctx context.Context, gate ConfirmationGate, creator SkillCreator, sessionID uint, args json.RawMessage) (string, error) {
	var proposal struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Body        string `json:"body"`
	}
	if err := json.Unmarshal(args, &proposal); err != nil {
		return "", err
	}
	if proposal.Name == "" || proposal.Description == "" || proposal.Body == "" {
		return "", errors.New("name, description, body all required")
	}
	notifID, err := gate.NotifyAttention(ctx, sessionID,
		"Confirm new skill: "+proposal.Name,
		proposal.Description)
	if err != nil {
		return "", err
	}
	if err := gate.WaitForResolution(ctx, notifID); err != nil {
		return "", err
	}
	if _, err := creator.Create(pkgskills.Skill{
		Name: proposal.Name, Description: proposal.Description, Body: proposal.Body,
		Source: pkgskills.SourceAI,
	}); err != nil {
		return "", err
	}
	return "skill " + proposal.Name + " saved as ai-generated", nil
}
```

Note: the `ConfirmationGate` semantics here assume "resolution = confirm". A reject path is added in Task 10 (frontend popup explicitly calls `ResolveNotification` only on accept; on reject it calls a new `RejectNotification` which wakes the waiter with an error). For now, single-channel semantics keep this task small.

- [ ] **Step 2: Add a `Reject` path to notification service**

In `notification.go`:

```go
func (n *Notification) RejectNotification(ctx context.Context, input notification_dto.ResolveNotificationInput) (*notification_dto.ResolveNotificationOutput, error) {
	if err := n.store.ResolveNotification(input.ID); err != nil {
		return nil, ierror.Error(ierror.ErrNotificationResolve, err)
	}
	wakeWaiterWithError(input.ID, errors.New("user rejected"))
	return &notification_dto.ResolveNotificationOutput{}, nil
}
```

Extend `wakeWaiter` into `wakeWaiterWithError(id uint, err error)` that delivers a sentinel via a second channel; update `WaitForResolution` to read from both.

- [ ] **Step 3: Register tool in agent dispatch**

Same as Skill / RequestAttention: append to per-turn tool registration when `session.Kind == "background"`. Wire dispatcher.

- [ ] **Step 4: Run tests**

- [ ] **Step 5: Commit**

```bash
git add backend/pkg/agent/tools/propose_skill_tool.go backend/pkg/agent/tools/propose_skill_tool_test.go backend/service/notification/ backend/pkg/agent/
git commit -m "feat(agent): add ProposeSkill tool gated on user confirmation"
```

---

## Section 4 — Frontend: notification surface + popup chat

### Task 8: Notification store + binding hook

**Files:**
- Create: `frontend/src/store/notificationsStore.ts`
- Create: `frontend/src/hooks/useNotificationsSubscription.ts`
- Create: `frontend/src/types/notifications.ts`

- [ ] **Step 1: Types**

```ts
export type NotificationKind = 'needs_attention' | 'info'

export type NotificationItem = {
  id: number
  session_id: number
  kind: NotificationKind
  title: string
  message: string
  created_at: string
  resolved_at: string | null
}
```

- [ ] **Step 2: Store**

```ts
import { create } from 'zustand'
import type { NotificationItem } from '@/types/notifications'

type State = {
  items: NotificationItem[]
  set: (items: NotificationItem[]) => void
  prepend: (item: NotificationItem) => void
  remove: (id: number) => void
}

export const useNotificationsStore = create<State>((set) => ({
  items: [],
  set: (items) => set({ items }),
  prepend: (item) => set((s) => ({ items: [item, ...s.items.filter((i) => i.id !== item.id)] })),
  remove: (id) => set((s) => ({ items: s.items.filter((i) => i.id !== id) })),
}))
```

- [ ] **Step 3: Subscription hook**

```ts
import { useEffect } from 'react'
import { Events } from '@wailsio/runtime'
import { Notification as NotificationBinding } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/notification'
import { useNotificationsStore } from '@/store/notificationsStore'
import type { NotificationItem } from '@/types/notifications'

export function useNotificationsSubscription() {
  const setItems = useNotificationsStore((s) => s.set)
  const prepend = useNotificationsStore((s) => s.prepend)
  const remove = useNotificationsStore((s) => s.remove)

  useEffect(() => {
    void NotificationBinding.ListNotifications({ include_resolved: false }).then((res) => {
      setItems(((res?.items ?? []) as unknown as NotificationItem[]))
    })
    const offCreate = Events.On('notification.created', (event: { data: NotificationItem }) => prepend(event.data))
    const offResolve = Events.On('notification.resolved', (event: { data: { id: number } }) => remove(event.data.id))
    return () => { offCreate(); offResolve() }
  }, [setItems, prepend, remove])
}
```

- [ ] **Step 4: Commit**

```bash
git add frontend/src/store/notificationsStore.ts frontend/src/hooks/useNotificationsSubscription.ts frontend/src/types/notifications.ts
git commit -m "feat(notification): add frontend store and subscription hook"
```

---

### Task 9: `NotificationBell` UI

**Files:**
- Create: `frontend/src/components/notifications/NotificationBell.tsx`
- Modify: app shell that hosts the chat header (search for `ChatHeader.tsx` or `AppShell.tsx`)

- [ ] **Step 1: Component**

```tsx
import { useState } from 'react'
import { Bell } from 'lucide-react'
import { useNotificationsStore } from '@/store/notificationsStore'
import { FloatingChatPanel } from './FloatingChatPanel'

export function NotificationBell() {
  const items = useNotificationsStore((s) => s.items)
  const [openSessionID, setOpenSessionID] = useState<number | null>(null)
  const [dropdownOpen, setDropdownOpen] = useState(false)

  return (
    <>
      <div className="relative">
        <button type="button" className="rounded-md p-2 hover:bg-muted" onClick={() => setDropdownOpen((v) => !v)}>
          <Bell size={18} />
          {items.length > 0 && (
            <span className="absolute right-0 top-0 inline-flex h-4 min-w-[16px] items-center justify-center rounded-full bg-red-500 px-1 text-[10px] text-white">
              {items.length}
            </span>
          )}
        </button>
        {dropdownOpen && (
          <div className="absolute right-0 z-40 mt-2 w-80 rounded-md border bg-popover p-2 shadow-lg">
            {items.length === 0 ? (
              <p className="px-2 py-4 text-center text-xs opacity-60">无未处理通知</p>
            ) : (
              <ul className="flex flex-col gap-1">
                {items.map((n) => (
                  <li key={n.id}>
                    <button
                      type="button"
                      onClick={() => { setOpenSessionID(n.session_id); setDropdownOpen(false) }}
                      className="flex w-full flex-col items-start gap-1 rounded-md px-2 py-2 text-left hover:bg-muted"
                    >
                      <span className="text-sm font-medium">{n.title}</span>
                      <span className="line-clamp-2 text-xs opacity-70">{n.message}</span>
                    </button>
                  </li>
                ))}
              </ul>
            )}
          </div>
        )}
      </div>
      {openSessionID !== null && (
        <FloatingChatPanel sessionID={openSessionID} onClose={() => setOpenSessionID(null)} />
      )}
    </>
  )
}
```

- [ ] **Step 2: Mount in app shell**

Find the chat / main window header and add `<NotificationBell />` to its right cluster. Also call `useNotificationsSubscription()` from a top-level component (e.g. `App.tsx`) so the store stays current.

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/notifications/NotificationBell.tsx frontend/src/App.tsx frontend/src/components/...header...
git commit -m "feat(notification): add notification bell in app shell"
```

---

### Task 10: `FloatingChatPanel` popup

**Files:**
- Create: `frontend/src/components/notifications/FloatingChatPanel.tsx`

- [ ] **Step 1: Component (uses existing chat message + input components)**

```tsx
import { useEffect } from 'react'
import { X } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { Notification as NotificationBinding } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/notification'
import { ResolveNotificationInput } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/notification/notification_dto/models'
import { useNotificationsStore } from '@/store/notificationsStore'
import { ChatMessages } from '@/components/chat/ChatMessages'
import { ChatInput } from '@/components/chat/ChatInput'

type Props = {
  sessionID: number
  onClose: () => void
}

export function FloatingChatPanel({ sessionID, onClose }: Props) {
  const { t } = useTranslation()
  const items = useNotificationsStore((s) => s.items)
  const related = items.find((n) => n.session_id === sessionID) ?? null

  const handleResolve = async (accept: boolean) => {
    if (!related) return
    const input = new ResolveNotificationInput({ id: related.id })
    if (accept) {
      await NotificationBinding.ResolveNotification(input)
    } else {
      await NotificationBinding.RejectNotification(input)
    }
  }

  useEffect(() => {
    const onEsc = (e: KeyboardEvent) => { if (e.key === 'Escape') onClose() }
    window.addEventListener('keydown', onEsc)
    return () => window.removeEventListener('keydown', onEsc)
  }, [onClose])

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/40">
      <div className="flex h-[600px] max-h-[90vh] w-[640px] max-w-[95vw] flex-col rounded-md bg-card shadow-xl">
        <div className="flex items-center justify-between border-b px-3 py-2">
          <h3 className="text-sm font-semibold">{related?.title ?? `Session ${sessionID}`}</h3>
          <button type="button" onClick={onClose} className="rounded-md p-1 hover:bg-muted"><X size={16} /></button>
        </div>
        <div className="flex-1 overflow-y-auto p-3">
          <ChatMessages sessionId={sessionID} />
        </div>
        {related?.kind === 'needs_attention' ? (
          <div className="flex items-center gap-2 border-t p-3">
            <button type="button" onClick={() => { void handleResolve(true) }} className="rounded-md bg-primary px-3 py-1 text-sm text-primary-foreground">{t('common.confirm') || '确认'}</button>
            <button type="button" onClick={() => { void handleResolve(false) }} className="rounded-md border px-3 py-1 text-sm">{t('common.reject') || '拒绝'}</button>
            <span className="text-xs opacity-60">{related.message}</span>
          </div>
        ) : (
          <div className="border-t p-2">
            <ChatInput sessionId={sessionID} />
          </div>
        )}
      </div>
    </div>
  )
}
```

(If existing `ChatMessages` / `ChatInput` don't accept `sessionId` as a prop, adapt — they may read from store; in that case wrap with a temporary session-context provider.)

- [ ] **Step 2: Add a backend Wails method `RejectNotification`**

Already added in Task 7 Step 2; regenerate bindings.

- [ ] **Step 3: i18n strings**

```ts
'common.confirm': '确认',
'common.reject': '拒绝',
```

- [ ] **Step 4: Smoke test**

Run app. From a backend dev console (or a quick test action), spawn a hidden session + call `Notify`. Verify bell badge increments, click → popup appears, accept resolves and bell decrements.

- [ ] **Step 5: Commit**

```bash
git add frontend/src/components/notifications/FloatingChatPanel.tsx frontend/src/i18n/
git commit -m "feat(notification): add floating popup for session intervention"
```

---

## Section 5 — End-to-end smoke test

### Task 11: Manual E2E

- [ ] **Step 1: Run app**

```
task dev
```

- [ ] **Step 2: Trigger a hidden session manually**

Open the browser devtools console in the app window. Call:

```js
await window.go.agent.Agent.SpawnHiddenSession({
  title: 'manual bg', system_prompt: '', user_message: 'Please ask me for my favorite color via RequestUserAttention', skill_name: ''
})
```

(Adjust path if Wails exposes services under a different namespace.)

- [ ] **Step 3: Verify**

- The session does **not** appear in the chat sidebar list.
- Within ~10s the bell badge increments by 1.
- Click bell → dropdown shows "Confirm new skill: …" or "what your favorite color".
- Click → popup chat opens, shows messages from the hidden session.
- Type a reply in the popup → agent in the hidden session resumes.
- After agent finishes, notification resolves; bell badge goes back to 0.

- [ ] **Step 4: Confirm AI-skill flow**

In a regular chat session, ask the model: `请用 ProposeSkill 工具记录一个新的 skill 'demo-skill'，描述 'demo', 正文 'hello'`. Expect:
- A notification "Confirm new skill: demo-skill" appears.
- Popup shows confirm / reject buttons.
- Click confirm → `Skills` page now lists `demo-skill` with source `ai`.
- Repeat, but click reject → skill is NOT created; the agent's tool result is the error message.

- [ ] **Step 5: Commit any fix-ups discovered**

---

## Self-review

1. **Spec coverage** — spec §3.5 (hidden session, global notification, popup) is covered by Tasks 1–6, 8–10. Spec §3.4.3 AI-skill confirmation (deferred from Plan A) covered by Task 7.

2. **Placeholder scan** — should find none. Several "search for X" prompts exist where the exact file isn't known; they are explicit grep commands, not placeholders.

3. **Type consistency** —
   - Backend `Notification` model fields ↔ DTO `NotificationItem` fields ↔ frontend `NotificationItem` — all named lowercase JSON (`session_id`, `created_at`, `resolved_at`). Verify in code.
   - `Session.Kind` is string with values `"user"` / `"background"` everywhere.
   - Tool name constants (`SkillToolName`, `RequestAttentionToolName`, `ProposeSkillToolName`) are stable string literals matched by both registration and dispatch.

4. **Cross-plan consistency** — `BuildSkillTool` (from Plan A) and the new tools all use the same `ToolMeta` shape. `Skills.Manager()` accessor (Plan A Task 14) is reused in `ProposeSkill` wiring (Task 7).

---

## Execution handoff

Plan complete at `docs/superpowers/plans/2026-05-25-hidden-sessions.md`. Two execution options:

1. **Subagent-Driven (recommended)** — fresh subagent per task, review between tasks.
2. **Inline Execution** — execute in this session with checkpoints.

Plan A must merge first.
