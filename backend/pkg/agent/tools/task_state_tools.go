package tools

import (
	"context"
	"encoding/json"
	"errors"

	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

// TaskStateStore persists resumable task-scoped state for hidden/automated sessions.
type TaskStateStore interface {
	SaveTaskState(sessionID uint, key, value string) error
	LoadTaskState(sessionID uint, key string) (value string, found bool, err error)
}

const (
	SaveTaskStateToolName = "SaveTaskState"
	LoadTaskStateToolName = "LoadTaskState"
)

type saveTaskStateInput struct {
	Key   string `json:"key" jsonschema:"description=state key,required"`
	Value string `json:"value" jsonschema:"description=state value string,required"`
}

type loadTaskStateInput struct {
	Key string `json:"key" jsonschema:"description=state key,required"`
}

func BuildSaveTaskStateTool() ToolMeta {
	return ToolMeta{
		Name:        SaveTaskStateToolName,
		Description: "Persist resumable state for the current task session. Use for device codes, verification URLs, expiry timestamps, and other cross-turn init state.",
		Category:    CategoryBuiltin,
		FormatPurpose: func(args json.RawMessage) string {
			var input saveTaskStateInput
			_ = json.Unmarshal(args, &input)
			return "Save task state: " + input.Key
		},
	}
}

func BuildLoadTaskStateTool() ToolMeta {
	return ToolMeta{
		Name:        LoadTaskStateToolName,
		Description: "Load previously persisted state for the current task session.",
		Category:    CategoryBuiltin,
		FormatPurpose: func(args json.RawMessage) string {
			var input loadTaskStateInput
			_ = json.Unmarshal(args, &input)
			return "Load task state: " + input.Key
		},
	}
}

func InvokeSaveTaskState(ctx context.Context, store TaskStateStore, sessionID uint, args json.RawMessage) (string, error) {
	var input saveTaskStateInput
	if err := json.Unmarshal(args, &input); err != nil {
		return "", err
	}
	if input.Key == "" {
		return "", errors.New("key is required")
	}
	if err := store.SaveTaskState(sessionID, input.Key, input.Value); err != nil {
		return "", err
	}
	return "saved state for key: " + input.Key, nil
}

func InvokeLoadTaskState(ctx context.Context, store TaskStateStore, sessionID uint, args json.RawMessage) (string, error) {
	var input loadTaskStateInput
	if err := json.Unmarshal(args, &input); err != nil {
		return "", err
	}
	if input.Key == "" {
		return "", errors.New("key is required")
	}
	value, found, err := store.LoadTaskState(sessionID, input.Key)
	if err != nil {
		return "", err
	}
	payload, err := json.Marshal(map[string]any{
		"key":   input.Key,
		"found": found,
		"value": value,
	})
	if err != nil {
		return "", err
	}
	return string(payload), nil
}

func NewSaveTaskStateTool(store TaskStateStore, sessionID uint) *function.FunctionTool[saveTaskStateInput, string] {
	meta := BuildSaveTaskStateTool()
	return function.NewFunctionTool(
		func(ctx context.Context, input saveTaskStateInput) (string, error) {
			payload, err := json.Marshal(input)
			if err != nil {
				return "", err
			}
			return InvokeSaveTaskState(ctx, store, sessionID, payload)
		},
		function.WithName(SaveTaskStateToolName),
		function.WithDescription(meta.Description),
	)
}

func NewLoadTaskStateTool(store TaskStateStore, sessionID uint) *function.FunctionTool[loadTaskStateInput, string] {
	meta := BuildLoadTaskStateTool()
	return function.NewFunctionTool(
		func(ctx context.Context, input loadTaskStateInput) (string, error) {
			payload, err := json.Marshal(input)
			if err != nil {
				return "", err
			}
			return InvokeLoadTaskState(ctx, store, sessionID, payload)
		},
		function.WithName(LoadTaskStateToolName),
		function.WithDescription(meta.Description),
	)
}
