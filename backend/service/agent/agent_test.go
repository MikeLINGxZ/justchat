package agent

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/id/prompt_id"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/ierror"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/prompt"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/agent/agent_dto"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/storage"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/model"
)

func newTestStorage(t *testing.T) *storage.Storage {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	stor, err := storage.NewStorageFromDB(db)
	if err != nil {
		t.Fatalf("new storage: %v", err)
	}
	return stor
}

func TestCreateSessionAndListSessions(t *testing.T) {
	svc := NewAgent(newTestStorage(t))

	created, err := svc.CreateSession(context.Background(), agent_dto.CreateSessionInput{})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	if created.SessionID == 0 {
		t.Fatal("expected session id to be set")
	}
	if created.Title != "New Chat" {
		t.Fatalf("expected default title, got %q", created.Title)
	}

	listed, err := svc.ListSessions(context.Background(), agent_dto.ListSessionsInput{Limit: 10})
	if err != nil {
		t.Fatalf("list sessions: %v", err)
	}
	if len(listed.Sessions) != 1 {
		t.Fatalf("expected 1 session, got %d", len(listed.Sessions))
	}
	if listed.Sessions[0].ID != created.SessionID {
		t.Fatalf("expected listed session id %d, got %d", created.SessionID, listed.Sessions[0].ID)
	}
}

func TestSpawnTaskSessionCreatesTaskSession(t *testing.T) {
	svc := NewAgent(newTestStorage(t))

	out, err := svc.SpawnTaskSession(context.Background(), agent_dto.SpawnTaskSessionInput{
		Title:       "bg test",
		UserMessage: "hi",
	})
	if err != nil {
		t.Fatal(err)
	}
	if out.SessionID == 0 {
		t.Fatal("expected non-zero session id")
	}

	session, err := svc.manager.Storage().GetSession(out.SessionID)
	if err != nil {
		t.Fatal(err)
	}
	if session.Kind != "task" {
		t.Fatalf("expected task kind, got %q", session.Kind)
	}
}

func TestLoadSessionMessagesMapsStorageMessages(t *testing.T) {
	stor := newTestStorage(t)
	session, err := stor.CreateSession(data_models.Session{Title: "Test", Status: "idle"})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	_, err = stor.CreateMessage(data_models.Message{
		SessionID:   session.ID,
		Role:        "assistant",
		ContentType: "text",
		Content:     "hello",
		ModelName:   "gpt-test",
		TokensIn:    10,
		TokensOut:   20,
		Extra:       "extra",
	})
	if err != nil {
		t.Fatalf("create message: %v", err)
	}

	svc := NewAgent(stor)
	result, err := svc.LoadSessionMessages(context.Background(), agent_dto.LoadSessionMessagesInput{
		SessionID: session.ID,
		Offset:    0,
		Limit:     10,
	})
	if err != nil {
		t.Fatalf("load session messages: %v", err)
	}

	if result.Total != 1 {
		t.Fatalf("expected total 1, got %d", result.Total)
	}
	if len(result.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(result.Messages))
	}
	if result.Messages[0].Content != "hello" {
		t.Fatalf("expected content hello, got %q", result.Messages[0].Content)
	}
}

func TestMarkSessionReadPersistsIdleStatus(t *testing.T) {
	stor := newTestStorage(t)
	session, err := stor.CreateSession(data_models.Session{Title: "Test", Status: "done-unread"})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	svc := NewAgent(stor)
	if _, err := svc.MarkSessionRead(context.Background(), agent_dto.MarkSessionReadInput{SessionID: session.ID}); err != nil {
		t.Fatalf("mark session read: %v", err)
	}

	updated, err := stor.GetSession(session.ID)
	if err != nil {
		t.Fatalf("get session: %v", err)
	}
	if updated.Status != "idle" {
		t.Fatalf("expected idle status, got %q", updated.Status)
	}
}

func TestHelpersConvertModelsToDTOs(t *testing.T) {
	now := time.Date(2026, 5, 15, 12, 0, 0, 0, time.UTC)
	session := data_models.Session{
		OrmModel: data_models.OrmModel{
			ID:        7,
			CreatedAt: now,
			UpdatedAt: now.Add(time.Minute),
		},
		Title:   "Tea",
		Starred: true,
		Status:  "idle",
	}
	msg := data_models.Message{
		OrmModel: data_models.OrmModel{
			ID:        9,
			CreatedAt: now,
		},
		SessionID:   7,
		Role:        "assistant",
		ContentType: "text",
		Content:     "reply",
		ModelName:   "gpt-test",
		AgentName:   "main",
		TokensIn:    1,
		TokensOut:   2,
		Extra:       "meta",
	}

	sessionItem := toSessionItem(session)
	if sessionItem.ID != 7 || sessionItem.Title != "Tea" || !sessionItem.Starred {
		t.Fatalf("unexpected session item: %+v", sessionItem)
	}

	messageItem := toMessageItem(msg)
	if messageItem.ID != 9 || messageItem.Content != "reply" || messageItem.ModelName != "gpt-test" {
		t.Fatalf("unexpected message item: %+v", messageItem)
	}
}

func TestBuildTitlePromptUsesRegisteredPromptID(t *testing.T) {
	prompt.Register(prompt_id.GenChatTitle, "CUSTOM TITLE PROMPT")

	got := buildTitlePrompt("user: hello")
	want := "CUSTOM TITLE PROMPT\n\nuser: hello"

	if got != want {
		t.Fatalf("unexpected title prompt: got %q want %q", got, want)
	}
}

func TestCollectTitleFromEventsPrefersStreamedDeltaWithoutDuplicatingFinalMessage(t *testing.T) {
	events := make(chan *event.Event, 2)
	events <- &event.Event{
		Response: &model.Response{
			Choices: []model.Choice{{
				Delta: model.Message{Content: "简单"},
			}},
		},
	}
	events <- &event.Event{
		Response: &model.Response{
			Choices: []model.Choice{{
				Delta:   model.Message{Content: "问候"},
				Message: model.Message{Content: "简单问候"},
			}},
			Done: true,
		},
	}
	close(events)

	if got := collectTitleFromEvents(events); got != "简单问候" {
		t.Fatalf("expected deduplicated streamed title, got %q", got)
	}
}

func TestCollectTitleFromEventsFallsBackToFinalMessageWhenNoDeltaArrives(t *testing.T) {
	events := make(chan *event.Event, 1)
	events <- &event.Event{
		Response: &model.Response{
			Choices: []model.Choice{{
				Message: model.Message{Content: "Hello"},
			}},
			Done: true,
		},
	}
	close(events)

	if got := collectTitleFromEvents(events); got != "Hello" {
		t.Fatalf("expected final message fallback title, got %q", got)
	}
}

func TestConvertAttachmentsRejectsTooMany(t *testing.T) {
	in := make([]agent_dto.AttachmentInput, maxAttachmentCount+1)
	for i := range in {
		in[i] = agent_dto.AttachmentInput{Path: "/tmp/x"}
	}
	_, err := convertAttachments(in)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
	if !strings.Contains(err.Error(), ierror.ErrAgentTooManyAttachments.Msg()) {
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
	if !strings.Contains(err.Error(), ierror.ErrAgentAttachmentNotFound.Msg()) {
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
	if !strings.Contains(err.Error(), ierror.ErrAgentAttachmentSize.Msg()) {
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
