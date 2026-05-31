package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

const TodoWriteToolName = "todo_write"

type todoWriteItem struct {
	Title  string `json:"title" jsonschema:"description=Short task title,required"`
	Status string `json:"status" jsonschema:"description=Task status: pending, in_progress, or completed,required"`
}

type todoWriteInput struct {
	Items []todoWriteItem `json:"items" jsonschema:"description=Current task list,required"`
}

type todoWriteOutput struct {
	Total      int             `json:"total"`
	Completed  int             `json:"completed"`
	InProgress int             `json:"in_progress"`
	Pending    int             `json:"pending"`
	Items      []todoWriteItem `json:"items"`
	Summary    string          `json:"summary"`
}

// BuildTodoWriteTool returns registry metadata for the task progress tool.
func BuildTodoWriteTool() ToolMeta {
	return ToolMeta{
		Name:        TodoWriteToolName,
		Description: "Track progress for a multi-step user task and keep the assistant organized.",
		Category:    CategoryBuiltin,
		FormatPurpose: func(_ json.RawMessage) string {
			return "Update task progress"
		},
	}
}

// todoWriteFunc validates a task list and returns a compact progress summary.
func todoWriteFunc(_ context.Context, input todoWriteInput) (todoWriteOutput, error) {
	if len(input.Items) == 0 {
		return todoWriteOutput{}, errors.New("at least one todo item is required")
	}
	out := todoWriteOutput{Total: len(input.Items), Items: input.Items}
	for _, item := range input.Items {
		switch strings.TrimSpace(item.Status) {
		case "completed":
			out.Completed++
		case "in_progress":
			out.InProgress++
		default:
			out.Pending++
		}
	}
	out.Summary = fmt.Sprintf("Progress: %d completed, %d in progress, %d pending.", out.Completed, out.InProgress, out.Pending)
	return out, nil
}

// NewTodoWriteTool creates the function tool used to track task progress.
func NewTodoWriteTool() *function.FunctionTool[todoWriteInput, todoWriteOutput] {
	meta := BuildTodoWriteTool()
	return function.NewFunctionTool(
		todoWriteFunc,
		function.WithName(TodoWriteToolName),
		function.WithDescription(meta.Description),
	)
}
