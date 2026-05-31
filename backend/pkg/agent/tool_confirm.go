package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/agent/tools"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/id/event_id"
	"trpc.group/trpc-go/trpc-agent-go/model"
	toolpkg "trpc.group/trpc-go/trpc-agent-go/tool"
)

func (m *Manager) newToolCallbacks() *toolpkg.Callbacks {
	return toolpkg.NewCallbacks().
		RegisterBeforeTool(func(
			ctx context.Context,
			args *toolpkg.BeforeToolArgs,
		) (*toolpkg.BeforeToolResult, error) {
			meta, ok := m.lookupToolMeta(args.ToolName)
			if !ok {
				return nil, nil
			}

			runInfo, ok := getRunContext(ctx)
			if !ok {
				return nil, nil
			}

			purpose := ""
			if meta.FormatPurpose != nil && len(args.Arguments) > 0 {
				purpose = meta.FormatPurpose(json.RawMessage(args.Arguments))
			}

			argsText := "{}"
			if len(args.Arguments) > 0 {
				argsText = string(args.Arguments)
			}

			decision, err := m.decideToolConfirmation(ctx, ToolConfirmDecisionInput{
				ToolName:     args.ToolName,
				Arguments:    argsText,
				Purpose:      purpose,
				Meta:         meta,
				ModelName:    runInfo.ModelName,
				BaseURL:      runInfo.BaseURL,
				APIKey:       runInfo.APIKey,
				ProviderType: runInfo.ProviderType,
				SessionID:    runInfo.SessionID,
			})
			if err != nil {
				decision = ToolConfirmDecision{RequireConfirm: true, Reason: "fallback after policy error", Source: "fallback"}
			}
			if decision.Deny {
				return &toolpkg.BeforeToolResult{
					CustomResult: map[string]any{
						"error":  buildDeniedToolMessage(decision.Reason),
						"reason": decision.Reason,
					},
				}, nil
			}
			shouldConfirm := decision.RequireConfirm
			if !shouldConfirm {
				return nil, nil
			}

			m.emitToolConfirmRequest(runInfo.SessionID, args.ToolCallID, args.ToolName, argsText, purpose)
			m.updateSessionStatus(runInfo.SessionID, "waiting-unread")

			resp, err := m.Streams().WaitForConfirm(ctx, runInfo.SessionID)
			if err != nil {
				m.updateSessionStatus(runInfo.SessionID, "idle")
				return nil, err
			}

			action := normalizeConfirmAction(resp)
			comment := strings.TrimSpace(resp.Message)
			confirmContent := formatConfirmContent(action, resp.Message)

			_, _ = m.Storage().CreateMessage(data_models.Message{
				SessionID:   runInfo.SessionID,
				Role:        "user",
				ContentType: "confirm_response",
				Content:     confirmContent,
				Extra:       marshalConfirmMetadata(action, comment),
			})
			m.updateSessionStatus(runInfo.SessionID, "loading")

			if action == "comment" {
				return &toolpkg.BeforeToolResult{
					Context: withToolConfirmation(ctx, toolConfirmation{
						Approved: false,
						Message:  comment,
						Action:   action,
					}),
					CustomResult: map[string]any{
						"status":  "user_comment",
						"comment": comment,
					},
				}, nil
			}

			if action == "reject" {
				return &toolpkg.BeforeToolResult{
					Context: withToolConfirmation(ctx, toolConfirmation{
						Approved: false,
						Message:  comment,
						Action:   action,
					}),
					CustomResult: map[string]any{
						"error":   buildRejectedToolMessage(resp.Message),
						"comment": comment,
					},
				}, nil
			}

			return &toolpkg.BeforeToolResult{
				Context: withToolConfirmation(ctx, toolConfirmation{
					Approved: true,
					Message:  comment,
					Action:   action,
				}),
			}, nil
		}).
		RegisterToolResultMessages(func(
			ctx context.Context,
			in *toolpkg.ToolResultMessagesInput,
		) (any, error) {
			confirmation, ok := getToolConfirmation(ctx)
			if !ok {
				return nil, nil
			}

			defaultMsg, ok := in.DefaultToolMessage.(model.Message)
			if !ok {
				return nil, nil
			}

			note := formatToolConfirmationNote(confirmation)
			if note == "" {
				return nil, nil
			}

			if confirmation.Action == "comment" && strings.TrimSpace(confirmation.Message) != "" {
				defaultMsg.Content = note
				return defaultMsg, nil
			}

			if strings.TrimSpace(defaultMsg.Content) == "" {
				defaultMsg.Content = note
			} else {
				defaultMsg.Content = fmt.Sprintf("%s\n\n%s", note, defaultMsg.Content)
			}

			return defaultMsg, nil
		})
}

func (m *Manager) decideToolConfirmation(ctx context.Context, input ToolConfirmDecisionInput) (ToolConfirmDecision, error) {
	m.mu.RLock()
	decider := m.toolConfirmDecider
	m.mu.RUnlock()
	if decider == nil {
		return ToolConfirmDecision{RequireConfirm: true, Reason: "missing decider", Source: "fallback"}, nil
	}
	return decider(ctx, input)
}

func (m *Manager) lookupToolMeta(toolName string) (tools.ToolMeta, bool) {
	if meta, ok := m.ToolRegistry().Get(toolName); ok {
		return meta, true
	}
	config, err := loadAgentConfig()
	if err != nil {
		return tools.ToolMeta{}, false
	}
	for _, item := range config.Extensions {
		for _, current := range item.Tools {
			if current.ToolID == toolName {
				return tools.ToolMeta{
					Name:            current.ToolID,
					Description:     current.Description,
					Category:        tools.CategoryUser,
					RequiresConfirm: current.RequiresConfirm,
					FormatPurpose: func(args json.RawMessage) string {
						return fmt.Sprintf("Run %s with arguments %s", current.Name, string(args))
					},
				}, true
			}
		}
	}
	return tools.ToolMeta{}, false
}

func (m *Manager) emitToolConfirmRequest(sessionID uint, requestID, toolName, args, purpose string) {
	app := m.App()
	if app == nil {
		return
	}

	app.Event.Emit(event_id.AgentStreamConfirmRequest, map[string]any{
		"sessionId": sessionID,
		"requestId": requestID,
		"toolName":  toolName,
		"args":      args,
		"purpose":   purpose,
	})
}

func (m *Manager) updateSessionStatus(sessionID uint, status string) {
	_ = m.Storage().UpdateSessionStatus(sessionID, status)

	app := m.App()
	if app == nil {
		return
	}

	app.Event.Emit(event_id.AgentSessionStatus, map[string]any{
		"sessionId": sessionID,
		"status":    status,
	})
}

func buildRejectedToolMessage(message string) string {
	if message == "" {
		return "tool execution rejected by user"
	}
	return fmt.Sprintf("tool execution rejected by user: %s", message)
}

func buildDeniedToolMessage(reason string) string {
	reason = strings.TrimSpace(reason)
	if reason == "" {
		return "tool execution denied by policy"
	}
	return fmt.Sprintf("tool execution denied by policy: %s", reason)
}

func normalizeConfirmAction(resp ConfirmResponse) string {
	action := strings.TrimSpace(resp.Action)
	switch action {
	case "approve", "reject", "comment":
		return action
	}
	if resp.Approved {
		return "approve"
	}
	return "reject"
}

func formatConfirmContent(action, message string) string {
	message = strings.TrimSpace(message)
	status := action
	switch action {
	case "approve":
		status = "approved"
	case "reject":
		status = "rejected"
	case "comment":
		status = "commented"
	}
	if message == "" {
		return status
	}
	return fmt.Sprintf("%s\ncomment: %s", status, message)
}

func formatToolConfirmationNote(confirmation toolConfirmation) string {
	message := strings.TrimSpace(confirmation.Message)

	if confirmation.Action == "comment" && message != "" {
		return fmt.Sprintf(
			"Tool execution paused. The user submitted this comment: %s. Do not execute the previous tool call unchanged. Revise the tool call to incorporate the comment, then continue.",
			message,
		)
	}

	status := "Tool confirmation: approved."
	if !confirmation.Approved {
		status = "Tool confirmation: rejected."
	}
	if message == "" {
		return status
	}
	return fmt.Sprintf("%s User comment: %s", status, message)
}

func marshalConfirmMetadata(action, comment string) string {
	payload, err := json.Marshal(map[string]string{
		"action":  action,
		"comment": comment,
	})
	if err != nil {
		return ""
	}
	return string(payload)
}
