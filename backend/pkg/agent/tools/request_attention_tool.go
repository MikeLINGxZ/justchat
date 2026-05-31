package tools

import (
	"context"
	"encoding/json"
	"errors"

	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

// AttentionRequester creates and waits on user-attention notifications.
type AttentionRequester interface {
	NotifyAttention(ctx context.Context, sessionID uint, title, message string) (uint, error)
	WaitForResolution(ctx context.Context, notificationID uint) error
}

// RequestAttentionToolName is the builtin tool exposed to task sessions.
const RequestAttentionToolName = "RequestUserAttention"

type requestAttentionInput struct {
	Title   string `json:"title" jsonschema:"description=Short headline shown to the user,required"`
	Message string `json:"message" jsonschema:"description=Detailed request for the user,required"`
}

// BuildRequestAttentionTool returns registry metadata for the attention tool.
func BuildRequestAttentionTool() ToolMeta {
	return ToolMeta{
		Name:        RequestAttentionToolName,
		Description: "Pause the task and ask the user for help or clarification.",
		Category:    CategoryBuiltin,
		FormatPurpose: func(args json.RawMessage) string {
			var input requestAttentionInput
			_ = json.Unmarshal(args, &input)
			return "Ask user for help: " + input.Title
		},
	}
}

// InvokeRequestAttention creates a notification and blocks until the user resolves it.
func InvokeRequestAttention(ctx context.Context, requester AttentionRequester, sessionID uint, args json.RawMessage) (string, error) {
	var input requestAttentionInput
	if err := json.Unmarshal(args, &input); err != nil {
		return "", err
	}
	if input.Title == "" {
		return "", errors.New("title is required")
	}
	notificationID, err := requester.NotifyAttention(ctx, sessionID, input.Title, input.Message)
	if err != nil {
		return "", err
	}
	if err := requester.WaitForResolution(ctx, notificationID); err != nil {
		return "", err
	}
	return "user replied; resume from the latest message in this session", nil
}

// NewRequestAttentionTool creates the function tool bound to one task session.
func NewRequestAttentionTool(requester AttentionRequester, sessionID uint) *function.FunctionTool[requestAttentionInput, string] {
	meta := BuildRequestAttentionTool()
	return function.NewFunctionTool(
		func(ctx context.Context, input requestAttentionInput) (string, error) {
			payload, err := json.Marshal(input)
			if err != nil {
				return "", err
			}
			return InvokeRequestAttention(ctx, requester, sessionID, payload)
		},
		function.WithName(RequestAttentionToolName),
		function.WithDescription(meta.Description),
	)
}
