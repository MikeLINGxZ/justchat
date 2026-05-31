package agent

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/id/event_id"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/ierror"
	pkgProvider "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/provider"
	agentpkg "trpc.group/trpc-go/trpc-agent-go/agent"
	"trpc.group/trpc-go/trpc-agent-go/event"
	"trpc.group/trpc-go/trpc-agent-go/model"
)

// ChatHandler processes chat requests and bridges runner events to Wails/storage.
type ChatHandler struct {
	manager *Manager
}

// SendMessageParams defines the inputs required to run a chat turn.
type SendMessageParams struct {
	SessionID        uint
	Content          string
	SystemPrompt     string
	PrimedSkill      string
	BaseURL          string
	ApiKey           string
	ModelName        string
	ProviderType     pkgProvider.Type
	EnabledUserTools []string
	Attachments      []Attachment
}

// NewChatHandler creates a chat handler for the given manager.
func NewChatHandler(m *Manager) *ChatHandler {
	return &ChatHandler{manager: m}
}

func (ch *ChatHandler) emitEvent(name string, data any) {
	app := ch.manager.App()
	if app != nil {
		app.Event.Emit(name, data)
	}
}

func (ch *ChatHandler) updateSessionStatus(sessionID uint, status string) {
	_ = ch.manager.Storage().UpdateSessionStatus(sessionID, status)
	ch.emitEvent(event_id.AgentSessionStatus, map[string]any{
		"sessionId": sessionID,
		"status":    status,
	})
}

func (ch *ChatHandler) loadSessionHistory(sessionID uint) ([]model.Message, error) {
	storedMessages, err := ch.manager.Storage().ListMessagesForSession(sessionID, 0, 1000)
	if err != nil {
		return nil, fmt.Errorf("list session messages: %w", err)
	}

	history := make([]model.Message, 0, len(storedMessages))

	for _, stored := range storedMessages {
		msg, ok := toModelMessage(stored)
		if !ok {
			continue
		}
		history = append(history, msg)
	}

	return history, nil
}

func toModelMessage(stored data_models.Message) (model.Message, bool) {
	switch stored.ContentType {
	case "text":
		if stored.Role == "user" {
			// Attachments are loaded lazily from disk on every replay by design:
			// the spec keeps attachments path-only, no caching.
			// Malformed JSON degrades silently to text-only, matching the
			// fallback pattern used by the tool_call arm below.
			atts, err := UnmarshalAttachments(stored.Attachments)
			if err == nil && len(atts) > 0 {
				return BuildUserMessage(stored.Content, atts), true
			}
			return model.NewUserMessage(stored.Content), true
		}
		return model.NewAssistantMessage(stored.Content), true
	case "thinking":
		return model.Message{
			Role:             model.RoleAssistant,
			ReasoningContent: stored.Content,
		}, true
	case "tool_result":
		return model.Message{
			Role:     model.RoleTool,
			ToolName: stored.AgentName,
			Content:  stored.Content,
		}, true
	case "confirm_response":
		return model.NewUserMessage(stored.Content), true
	case "tool_call":
		var payload struct {
			ID   string          `json:"id"`
			Name string          `json:"name"`
			Args json.RawMessage `json:"args"`
		}
		if err := json.Unmarshal([]byte(stored.Content), &payload); err != nil {
			return model.NewAssistantMessage(stored.Content), true
		}
		return model.Message{
			Role: model.RoleAssistant,
			ToolCalls: []model.ToolCall{
				{
					Type: "function",
					ID:   payload.ID,
					Function: model.FunctionDefinitionParam{
						Name:      payload.Name,
						Arguments: payload.Args,
					},
				},
			},
		}, true
	default:
		return model.Message{}, false
	}
}

func normalizeToolArguments(raw json.RawMessage) (string, bool) {
	trimmed := bytes.TrimSpace(raw)
	if len(trimmed) == 0 {
		return "{}", true
	}
	if !json.Valid(trimmed) {
		return "", false
	}

	var compact bytes.Buffer
	if err := json.Compact(&compact, trimmed); err != nil {
		return "", false
	}

	return compact.String(), true
}

func toolCallFingerprint(toolName, args string) string {
	if toolName == "" || args == "" {
		return ""
	}
	return toolName + "\x00" + args
}

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

	session, err := stor.GetSession(params.SessionID)
	if err != nil {
		return fmt.Errorf("get session: %w", err)
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

	r, err := ch.manager.GetOrCreateRunner(
		params.BaseURL,
		params.ApiKey,
		params.ModelName,
		params.ProviderType,
		params.EnabledUserTools,
		session.Kind,
		params.SessionID,
	)
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
		SessionID:    params.SessionID,
		ModelName:    params.ModelName,
		BaseURL:      params.BaseURL,
		APIKey:       params.ApiKey,
		ProviderType: params.ProviderType,
	})

	systemPrompt := ch.memorySystemPrompt(params.SystemPrompt, params.Content)
	userMsg := BuildUserMessage(buildRunContent(systemPrompt, params.PrimedSkill, params.Content), params.Attachments)

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

	go ch.consumeEventsWithParams(streamCtx, params, events)
	return nil
}

// buildRunContent folds hidden-session bootstrap instructions into the first user turn
// without changing the shared runner's global prompt contract.
func buildRunContent(systemPrompt, primedSkill, content string) string {
	systemPrompt = strings.TrimSpace(systemPrompt)
	primedSkill = strings.TrimSpace(primedSkill)
	content = strings.TrimSpace(content)

	parts := make([]string, 0, 3)
	if systemPrompt != "" {
		parts = append(parts, "System instructions:\n"+systemPrompt)
	}
	if primedSkill != "" {
		parts = append(parts, "Before doing any other work, call the Skill tool with {\"name\":\""+primedSkill+"\"}.")
	}
	if content != "" {
		parts = append(parts, content)
	}
	return strings.Join(parts, "\n\n")
}

func (ch *ChatHandler) consumeEvents(ctx context.Context, sessionID uint, modelName string, events <-chan *event.Event) {
	ch.consumeEventsWithParams(ctx, SendMessageParams{SessionID: sessionID, ModelName: modelName}, events)
}

func (ch *ChatHandler) consumeEventsWithParams(ctx context.Context, params SendMessageParams, events <-chan *event.Event) {
	sessionID := params.SessionID
	modelName := params.ModelName
	defer ch.manager.Streams().Remove(sessionID)

	stor := ch.manager.Storage()
	var fullContent string
	var fullThinking string
	var tokensIn int
	var tokensOut int
	var streamSeq int64
	var primaryChoiceIndex *int
	handledToolCalls := make(map[string]struct{})
	lastToolCallFingerprint := ""

	emitChunk := func(delta, contentType, content string) {
		streamSeq++
		ch.emitEvent(event_id.AgentStreamChunk, map[string]any{
			"sessionId":   sessionID,
			"seq":         streamSeq,
			"delta":       delta,
			"content":     content,
			"contentType": contentType,
		})
	}

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
			errContent, _ := json.Marshal(map[string]string{
				"msg":    ierror.ErrAgentStreamError.Msg(),
				"detail": evt.Error.Message,
			})
			_, _ = stor.CreateMessage(data_models.Message{
				SessionID:   sessionID,
				Role:        "assistant",
				ContentType: "error",
				Content:     string(errContent),
				ModelName:   modelName,
			})
			return
		}

		for _, choice := range evt.Choices {
			if primaryChoiceIndex == nil {
				idx := choice.Index
				primaryChoiceIndex = &idx
			}
			if primaryChoiceIndex != nil && choice.Index != *primaryChoiceIndex {
				continue
			}

			delta := choice.Delta
			message := choice.Message
			isToolResultMessage := message.ToolID != ""

			if !isToolResultMessage && delta.ReasoningContent != "" {
				lastToolCallFingerprint = ""
				fullThinking += delta.ReasoningContent
				emitChunk(delta.ReasoningContent, "thinking", fullThinking)
			}
			if !isToolResultMessage && delta.ReasoningContent == "" && message.ReasoningContent != "" && fullThinking == "" {
				lastToolCallFingerprint = ""
				fullThinking = message.ReasoningContent
				emitChunk(message.ReasoningContent, "thinking", fullThinking)
			}

			if !isToolResultMessage && delta.Content != "" {
				lastToolCallFingerprint = ""
				fullContent += delta.Content
				emitChunk(delta.Content, "text", fullContent)
			}
			if !isToolResultMessage && delta.Content == "" && message.Content != "" && fullContent == "" {
				lastToolCallFingerprint = ""
				fullContent = message.Content
				emitChunk(message.Content, "text", fullContent)
			}

			toolCalls := delta.ToolCalls
			if len(toolCalls) == 0 && len(message.ToolCalls) > 0 {
				toolCalls = message.ToolCalls
			}

			for _, tc := range toolCalls {
				normalizedArgs, ok := normalizeToolArguments(tc.Function.Arguments)
				if !ok {
					continue
				}

				if tc.ID != "" {
					if _, exists := handledToolCalls[tc.ID]; exists {
						continue
					}
				}

				toolName := tc.Function.Name
				fingerprint := toolCallFingerprint(toolName, normalizedArgs)
				if fingerprint != "" && fingerprint == lastToolCallFingerprint {
					if tc.ID != "" {
						handledToolCalls[tc.ID] = struct{}{}
					}
					continue
				}
				lastToolCallFingerprint = fingerprint
				if tc.ID != "" {
					handledToolCalls[tc.ID] = struct{}{}
				}

				args := normalizedArgs

				meta, hasMeta := ch.manager.ToolRegistry().Get(toolName)
				purpose := ""
				if hasMeta && meta.FormatPurpose != nil {
					purpose = meta.FormatPurpose(json.RawMessage(args))
				}

				payload, _ := json.Marshal(map[string]any{
					"id":   tc.ID,
					"name": toolName,
					"args": json.RawMessage(args),
				})

				ch.emitEvent(event_id.AgentStreamToolCall, map[string]any{
					"sessionId": sessionID,
					"toolName":  toolName,
					"args":      args,
					"purpose":   purpose,
				})

				_, _ = stor.CreateMessage(data_models.Message{
					SessionID:   sessionID,
					Role:        "assistant",
					ContentType: "tool_call",
					Content:     string(payload),
					ModelName:   modelName,
					Extra:       purpose,
				})

			}

			if message.ToolID != "" {
				lastToolCallFingerprint = ""
				ch.emitEvent(event_id.AgentStreamToolResult, map[string]any{
					"sessionId": sessionID,
					"toolName":  message.ToolName,
					"result":    message.Content,
				})
				_, _ = stor.CreateMessage(data_models.Message{
					SessionID:   sessionID,
					Role:        "tool",
					ContentType: "tool_result",
					Content:     message.Content,
					ModelName:   modelName,
					AgentName:   message.ToolName,
				})
			}
		}

		if evt.Usage != nil {
			tokensIn = evt.Usage.PromptTokens
			tokensOut = evt.Usage.CompletionTokens
		}

		if evt.Done {
			break
		}
	}

	if fullThinking != "" {
		_, _ = stor.CreateMessage(data_models.Message{
			SessionID:   sessionID,
			Role:        "assistant",
			ContentType: "thinking",
			Content:     fullThinking,
			ModelName:   modelName,
		})
	}

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
	ch.enqueueMemoryEncoding(params, fullContent)
}
