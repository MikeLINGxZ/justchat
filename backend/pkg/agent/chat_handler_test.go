package agent

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

type delayedRunner struct {
	events      chan *event.Event
	lastMessage model.Message
	lastRunOpts agentpkg.RunOptions
}

func newDelayedRunner() *delayedRunner {
	return &delayedRunner{events: make(chan *event.Event, 1)}
}

func (r *delayedRunner) Run(
	ctx context.Context,
	userID string,
	sessionID string,
	message model.Message,
	runOpts ...agentpkg.RunOption,
) (<-chan *event.Event, error) {
	r.lastMessage = message
	var opts agentpkg.RunOptions
	for _, opt := range runOpts {
		opt(&opts)
	}
	r.lastRunOpts = opts
	return r.events, nil
}

func (r *delayedRunner) Close() error {
	return nil
}

func TestConsumeEventsPersistsFinalMessageContent(t *testing.T) {
	manager := NewManager(newTestStorage(t))
	handler := NewChatHandler(manager)

	session, err := manager.Storage().CreateSession(data_models.Session{
		Title:  "Test",
		Status: "loading",
	})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	events := make(chan *event.Event, 1)
	events <- &event.Event{
		Response: &model.Response{
			Choices: []model.Choice{
				{
					Message: model.Message{
						Role:    model.RoleAssistant,
						Content: "final answer",
					},
				},
			},
			Usage: &model.Usage{
				PromptTokens:     11,
				CompletionTokens: 22,
			},
			Done: true,
		},
	}
	close(events)

	handler.consumeEvents(context.Background(), session.ID, "gpt-test", events)

	messages, err := manager.Storage().ListMessagesForSession(session.ID, 0, 10)
	if err != nil {
		t.Fatalf("list messages: %v", err)
	}
	if len(messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(messages))
	}
	if messages[0].Content != "final answer" {
		t.Fatalf("expected final answer, got %q", messages[0].Content)
	}
	if messages[0].TokensIn != 11 || messages[0].TokensOut != 22 {
		t.Fatalf("unexpected token usage: in=%d out=%d", messages[0].TokensIn, messages[0].TokensOut)
	}

	updatedSession, err := manager.Storage().GetSession(session.ID)
	if err != nil {
		t.Fatalf("get session: %v", err)
	}
	if updatedSession.Status != "done-unread" {
		t.Fatalf("expected session status done-unread, got %q", updatedSession.Status)
	}
}

func TestSendMessageContinuesAfterRequestContextIsCanceled(t *testing.T) {
	manager := NewManager(newTestStorage(t))
	handler := NewChatHandler(manager)

	session, err := manager.Storage().CreateSession(data_models.Session{
		Title:  "Test",
		Status: "idle",
	})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	r := newDelayedRunner()
	manager.runners["https://api.example.com/gpt-test"] = r

	requestCtx, cancelRequest := context.WithCancel(context.Background())
	if err := handler.SendMessage(requestCtx, SendMessageParams{
		SessionID:    session.ID,
		Content:      "hello",
		BaseURL:      "https://api.example.com",
		ModelName:    "gpt-test",
		ProviderType: "openai_compatibility",
	}); err != nil {
		t.Fatalf("send message: %v", err)
	}

	cancelRequest()
	r.events <- &event.Event{
		Response: &model.Response{
			Choices: []model.Choice{
				{
					Delta: model.Message{
						Role:    model.RoleAssistant,
						Content: "still here",
					},
				},
			},
			Done: true,
		},
	}
	close(r.events)

	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		messages, err := manager.Storage().ListMessagesForSession(session.ID, 0, 10)
		if err != nil {
			t.Fatalf("list messages: %v", err)
		}
		for _, msg := range messages {
			if msg.Role == "assistant" && msg.Content == "still here" {
				return
			}
		}
		time.Sleep(10 * time.Millisecond)
	}

	messages, err := manager.Storage().ListMessagesForSession(session.ID, 0, 10)
	if err != nil {
		t.Fatalf("list messages: %v", err)
	}
	t.Fatalf("expected assistant reply after request context cancel, got messages: %+v", messages)
}

func TestSendMessageReplaysStoredSessionHistoryWhenSwitchingRunner(t *testing.T) {
	manager := NewManager(newTestStorage(t))
	handler := NewChatHandler(manager)

	session, err := manager.Storage().CreateSession(data_models.Session{
		Title:  "Test",
		Status: "idle",
	})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	_, err = manager.Storage().CreateMessage(data_models.Message{
		SessionID:   session.ID,
		Role:        "user",
		ContentType: "text",
		Content:     "old user question",
	})
	if err != nil {
		t.Fatalf("create user message: %v", err)
	}
	_, err = manager.Storage().CreateMessage(data_models.Message{
		SessionID:   session.ID,
		Role:        "assistant",
		ContentType: "text",
		Content:     "old assistant answer",
		ModelName:   "old-model",
	})
	if err != nil {
		t.Fatalf("create assistant message: %v", err)
	}

	r := newDelayedRunner()
	manager.runners["https://api.example.com/new-model"] = r

	if err := handler.SendMessage(context.Background(), SendMessageParams{
		SessionID:    session.ID,
		Content:      "new user question",
		BaseURL:      "https://api.example.com",
		ModelName:    "new-model",
		ProviderType: "openai_compatibility",
	}); err != nil {
		t.Fatalf("send message: %v", err)
	}

	if r.lastMessage.Content != "new user question" {
		t.Fatalf("expected current message to be passed to runner, got %q", r.lastMessage.Content)
	}

	if len(r.lastRunOpts.Messages) != 2 {
		t.Fatalf("expected 2 replayed history messages, got %d", len(r.lastRunOpts.Messages))
	}

	if r.lastRunOpts.Messages[0].Role != model.RoleUser || r.lastRunOpts.Messages[0].Content != "old user question" {
		t.Fatalf("unexpected first replayed message: %+v", r.lastRunOpts.Messages[0])
	}
	if r.lastRunOpts.Messages[1].Role != model.RoleAssistant || r.lastRunOpts.Messages[1].Content != "old assistant answer" {
		t.Fatalf("unexpected second replayed message: %+v", r.lastRunOpts.Messages[1])
	}
}

func TestConsumeEventsPersistsMessageToolCallWithoutBlockingForConfirm(t *testing.T) {
	manager := NewManager(newTestStorage(t))
	handler := NewChatHandler(manager)

	session, err := manager.Storage().CreateSession(data_models.Session{
		Title:  "Confirm Tool Call",
		Status: "loading",
	})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	args, err := json.Marshal(map[string]string{
		"command": "ls /tmp",
	})
	if err != nil {
		t.Fatalf("marshal args: %v", err)
	}

	events := make(chan *event.Event, 2)
	events <- &event.Event{
		Response: &model.Response{
			Choices: []model.Choice{
				{
					Message: model.Message{
						Role: model.RoleAssistant,
						ToolCalls: []model.ToolCall{
							{
								Type: "function",
								ID:   "call_shell_1",
								Function: model.FunctionDefinitionParam{
									Name:      "shell",
									Arguments: args,
								},
							},
						},
					},
				},
			},
		},
	}
	events <- &event.Event{
		Response: &model.Response{
			Done: true,
		},
	}
	close(events)

	streamCtx, _ := manager.Streams().Start(context.Background(), session.ID)
	handler.consumeEvents(streamCtx, session.ID, "gpt-test", events)

	messages, err := manager.Storage().ListMessagesForSession(session.ID, 0, 10)
	if err != nil {
		t.Fatalf("list messages: %v", err)
	}
	if len(messages) != 1 {
		t.Fatalf("expected only tool_call message, got %d", len(messages))
	}
	if messages[0].ContentType != "tool_call" {
		t.Fatalf("expected first message to be tool_call, got %q", messages[0].ContentType)
	}
}

func TestConsumeEventsDoesNotLeakToolResultIntoAssistantText(t *testing.T) {
	manager := NewManager(newTestStorage(t))
	handler := NewChatHandler(manager)

	session, err := manager.Storage().CreateSession(data_models.Session{
		Title:  "Tool Result Leak",
		Status: "loading",
	})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	events := make(chan *event.Event, 2)
	events <- &event.Event{
		Response: &model.Response{
			Choices: []model.Choice{
				{
					Message: model.Message{
						Role:     model.RoleTool,
						ToolID:   "call_shell_1",
						ToolName: "shell",
						Content:  `{"stdout":"file-a\nfile-b","stderr":"","exit_code":0}`,
					},
				},
			},
		},
	}
	events <- &event.Event{
		Response: &model.Response{
			Choices: []model.Choice{
				{
					Message: model.Message{
						Role:    model.RoleAssistant,
						Content: "目录里有两个文件。",
					},
				},
			},
			Done: true,
		},
	}
	close(events)

	handler.consumeEvents(context.Background(), session.ID, "gpt-test", events)

	messages, err := manager.Storage().ListMessagesForSession(session.ID, 0, 10)
	if err != nil {
		t.Fatalf("list messages: %v", err)
	}
	if len(messages) != 2 {
		t.Fatalf("expected tool_result and assistant text, got %d", len(messages))
	}
	if messages[0].ContentType != "tool_result" {
		t.Fatalf("expected first message to be tool_result, got %q", messages[0].ContentType)
	}
	if messages[1].ContentType != "text" {
		t.Fatalf("expected second message to be assistant text, got %q", messages[1].ContentType)
	}
	if messages[1].Content != "目录里有两个文件。" {
		t.Fatalf("expected assistant text without leaked tool result, got %q", messages[1].Content)
	}
}

func TestConsumeEventsIgnoresNonPrimaryChoiceStreamFragments(t *testing.T) {
	manager := NewManager(newTestStorage(t))
	handler := NewChatHandler(manager)

	session, err := manager.Storage().CreateSession(data_models.Session{
		Title:  "Multi Choice Stream",
		Status: "loading",
	})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	events := make(chan *event.Event, 2)
	events <- &event.Event{
		Response: &model.Response{
			Choices: []model.Choice{
				{
					Index: 0,
					Delta: model.Message{
						Role:    model.RoleAssistant,
						Content: "好的，我先去查看文档，找到飞书 CLI 的 npm 包名。",
					},
				},
				{
					Index: 1,
					Delta: model.Message{
						Role:    model.RoleAssistant,
						Content: "加载安装我先 CLI 的技能。",
					},
				},
			},
		},
	}
	events <- &event.Event{
		Response: &model.Response{
			Choices: []model.Choice{
				{
					Index: 0,
					Delta: model.Message{
						Role:    model.RoleAssistant,
						Content: " 找到了飞书 CLI 的 npm 包名。",
					},
					Message: model.Message{
						Role:    model.RoleAssistant,
						Content: "好的，我先去查看文档，找到飞书 CLI 的 npm 包名。 找到了飞书 CLI 的 npm 包名。",
					},
				},
				{
					Index: 1,
					Delta: model.Message{
						Role:    model.RoleAssistant,
						Content: " 查看好的我先，去，文档。",
					},
				},
			},
			Done: true,
		},
	}
	close(events)

	handler.consumeEvents(context.Background(), session.ID, "gpt-test", events)

	messages, err := manager.Storage().ListMessagesForSession(session.ID, 0, 10)
	if err != nil {
		t.Fatalf("list messages: %v", err)
	}
	if len(messages) != 1 {
		t.Fatalf("expected one assistant message, got %d", len(messages))
	}
	got := messages[0].Content
	if got != "好的，我先去查看文档，找到飞书 CLI 的 npm 包名。 找到了飞书 CLI 的 npm 包名。" {
		t.Fatalf("expected only primary choice content, got %q", got)
	}
}

func TestConsumeEventsDeduplicatesAdjacentToolCallReplays(t *testing.T) {
	manager := NewManager(newTestStorage(t))
	handler := NewChatHandler(manager)

	session, err := manager.Storage().CreateSession(data_models.Session{
		Title:  "Duplicate Tool Call",
		Status: "loading",
	})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	args := json.RawMessage(`{"command":"ls -la /Users/linhuafeng/Temp","work_dir":"/Users/linhuafeng/Temp"}`)

	events := make(chan *event.Event, 3)
	events <- &event.Event{
		Response: &model.Response{
			Choices: []model.Choice{
				{
					Delta: model.Message{
						ToolCalls: []model.ToolCall{
							{
								Type: "function",
								Function: model.FunctionDefinitionParam{
									Name:      "shell",
									Arguments: args,
								},
							},
						},
					},
				},
			},
		},
	}
	events <- &event.Event{
		Response: &model.Response{
			Choices: []model.Choice{
				{
					Message: model.Message{
						Role: model.RoleAssistant,
						ToolCalls: []model.ToolCall{
							{
								Type: "function",
								ID:   "call_shell_1",
								Function: model.FunctionDefinitionParam{
									Name:      "shell",
									Arguments: args,
								},
							},
						},
					},
				},
			},
			Done: true,
		},
	}
	close(events)

	handler.consumeEvents(context.Background(), session.ID, "gpt-test", events)

	messages, err := manager.Storage().ListMessagesForSession(session.ID, 0, 10)
	if err != nil {
		t.Fatalf("list messages: %v", err)
	}
	if len(messages) != 1 {
		t.Fatalf("expected only one tool_call message after dedupe, got %d", len(messages))
	}
	if messages[0].ContentType != "tool_call" {
		t.Fatalf("expected first message to be tool_call, got %q", messages[0].ContentType)
	}
	if messages[0].Content == "" {
		t.Fatal("expected persisted tool_call payload")
	}
}

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
	// Inject runner directly — bypasses the real GetOrCreateRunner.
	// Key format matches manager.GetOrCreateRunner's cache key.
	manager.runners["https://api.example.com/m"] = r

	tmpDir := t.TempDir()
	imgPath := filepath.Join(tmpDir, "tiny.png")
	pngBytes, _ := base64.StdEncoding.DecodeString(minimalPNGBase64)
	if err := os.WriteFile(imgPath, pngBytes, 0o644); err != nil {
		t.Fatalf("write png: %v", err)
	}

	err = handler.SendMessage(context.Background(), SendMessageParams{
		SessionID:    session.ID,
		Content:      "describe",
		BaseURL:      "https://api.example.com",
		ApiKey:       "x",
		ModelName:    "m",
		ProviderType: "openai_compatibility",
		Attachments: []Attachment{
			NormalizeAttachment(Attachment{Path: imgPath}),
		},
	})
	if err != nil {
		t.Fatalf("send: %v", err)
	}

	// 1. Persisted user message: Attachments column is valid JSON containing one image entry.
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

	// 2. Runner received a multimodal user message (one ContentPart of type image).
	if len(r.lastMessage.ContentParts) != 1 ||
		r.lastMessage.ContentParts[0].Type != model.ContentTypeImage {
		t.Fatalf("runner did not receive multimodal user message: %+v", r.lastMessage)
	}
}

func TestConsumeEventsSkipsIncompleteToolCallArgumentsUntilJsonIsComplete(t *testing.T) {
	manager := NewManager(newTestStorage(t))
	handler := NewChatHandler(manager)

	session, err := manager.Storage().CreateSession(data_models.Session{
		Title:  "Partial Tool Args",
		Status: "loading",
	})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	events := make(chan *event.Event, 3)
	events <- &event.Event{
		Response: &model.Response{
			Choices: []model.Choice{
				{
					Delta: model.Message{
						ToolCalls: []model.ToolCall{
							{
								Type: "function",
								ID:   "call_shell_1",
								Function: model.FunctionDefinitionParam{
									Name:      "shell",
									Arguments: json.RawMessage(`{"command":"ls -la /Users/linhuafeng/Temp"`),
								},
							},
						},
					},
				},
			},
		},
	}
	events <- &event.Event{
		Response: &model.Response{
			Choices: []model.Choice{
				{
					Message: model.Message{
						Role: model.RoleAssistant,
						ToolCalls: []model.ToolCall{
							{
								Type: "function",
								ID:   "call_shell_1",
								Function: model.FunctionDefinitionParam{
									Name:      "shell",
									Arguments: json.RawMessage(`{"command":"ls -la /Users/linhuafeng/Temp","work_dir":"/Users/linhuafeng/Temp"}`),
								},
							},
						},
					},
				},
			},
			Done: true,
		},
	}
	close(events)

	handler.consumeEvents(context.Background(), session.ID, "gpt-test", events)

	messages, err := manager.Storage().ListMessagesForSession(session.ID, 0, 10)
	if err != nil {
		t.Fatalf("list messages: %v", err)
	}
	if len(messages) != 1 {
		t.Fatalf("expected one tool_call after complete JSON arrives, got %d", len(messages))
	}
	if messages[0].ContentType != "tool_call" {
		t.Fatalf("expected first message to be tool_call, got %q", messages[0].ContentType)
	}
}
