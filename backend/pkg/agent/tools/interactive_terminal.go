package tools

import (
	"context"
	"encoding/json"
	"strings"

	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

const InteractiveTerminalToolName = "InteractiveTerminal"

type interactiveTerminalInput struct {
	Status string `json:"status" jsonschema:"enum=active,enum=done,description=active shows the terminal panel; done hides it when the user-assisted step is complete,required"`
	Output string `json:"output" jsonschema:"description=terminal text to show when status is active"`
}

type interactiveTerminalOutput struct {
	InteractiveTerminal bool   `json:"interactive_terminal"`
	TerminalStatus      string `json:"terminal_status"`
	TerminalOutput      string `json:"terminal_output,omitempty"`
}

func BuildInteractiveTerminalTool() ToolMeta {
	return ToolMeta{
		Name:        InteractiveTerminalToolName,
		Description: "Control a user-facing interactive terminal panel. Use status=active when the user must scan a QR code, open a login URL, copy a code, or provide input outside chat. Call status=done once that user-assisted step is complete so the panel can hide.",
		Category:    CategoryBuiltin,
		FormatPurpose: func(args json.RawMessage) string {
			var input interactiveTerminalInput
			_ = json.Unmarshal(args, &input)
			if strings.EqualFold(strings.TrimSpace(input.Status), "done") {
				return "Hide interactive terminal"
			}
			return "Show interactive terminal"
		},
	}
}

func NewInteractiveTerminalTool() *function.FunctionTool[interactiveTerminalInput, interactiveTerminalOutput] {
	meta := BuildInteractiveTerminalTool()
	return function.NewFunctionTool(
		func(ctx context.Context, input interactiveTerminalInput) (interactiveTerminalOutput, error) {
			status := strings.ToLower(strings.TrimSpace(input.Status))
			if status != "done" {
				status = "active"
			}
			return interactiveTerminalOutput{
				InteractiveTerminal: true,
				TerminalStatus:      status,
				TerminalOutput:      input.Output,
			}, nil
		},
		function.WithName(InteractiveTerminalToolName),
		function.WithDescription(meta.Description),
	)
}
