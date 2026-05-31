package tools

import (
	"context"
	"encoding/json"
	"errors"

	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

// CliCommandRunner is satisfied by the plugin service.
type CliCommandRunner interface {
	RunCliCommandSync(ctx context.Context, extensionID string, argv []string, outputMode string, timeoutSeconds int, onProgress func(result string)) (resultJSON string, err error)
}

type CliCommandProgressEmitter interface {
	EmitToolResult(sessionID uint, toolName, result string)
}

const RunCliCommandToolName = "RunCliCommand"
const defaultRunCliCommandTimeoutSeconds = 600

func BuildRunCliCommandTool() ToolMeta {
	return ToolMeta{
		Name:        RunCliCommandToolName,
		Description: "Run an arbitrary command against one installed CLI plugin using Lemontea's bundled runtime and isolated env. Call with {\"id\":\"cli:<name>\",\"argv\":[...],\"output_mode\":\"json|text|lines\",\"timeout_seconds\":600}. Returns the CLI run result JSON and may stream intermediate output while the command is running.",
		Category:    CategoryBuiltin,
		FormatPurpose: func(args json.RawMessage) string {
			var parsed struct {
				ID   string   `json:"id"`
				Argv []string `json:"argv"`
			}
			_ = json.Unmarshal(args, &parsed)
			if len(parsed.Argv) == 0 {
				return "Running CLI command for: " + parsed.ID
			}
			return "Running CLI command for: " + parsed.ID + " " + parsed.Argv[0]
		},
	}
}

type runCliCommandInput struct {
	ID             string   `json:"id" jsonschema:"description=extension ID in the form cli:<name>,required"`
	Argv           []string `json:"argv" jsonschema:"description=argv to run against the CLI executable,required"`
	OutputMode     string   `json:"output_mode,omitempty" jsonschema:"description=one of json text lines"`
	TimeoutSeconds int      `json:"timeout_seconds,omitempty" jsonschema:"description=timeout in seconds"`
}

func InvokeRunCliCommand(ctx context.Context, runner CliCommandRunner, args json.RawMessage, onProgress func(result string)) (string, error) {
	var parsed runCliCommandInput
	if err := json.Unmarshal(args, &parsed); err != nil {
		return "", err
	}
	if parsed.ID == "" {
		return "", errors.New("id is required")
	}
	if len(parsed.Argv) == 0 {
		return "", errors.New("argv is required")
	}
	if parsed.TimeoutSeconds <= 0 {
		parsed.TimeoutSeconds = defaultRunCliCommandTimeoutSeconds
	}
	return runner.RunCliCommandSync(ctx, parsed.ID, parsed.Argv, parsed.OutputMode, parsed.TimeoutSeconds, onProgress)
}

func NewRunCliCommandTool(runner CliCommandRunner, emitter CliCommandProgressEmitter, sessionID uint) *function.FunctionTool[runCliCommandInput, string] {
	meta := BuildRunCliCommandTool()
	return function.NewFunctionTool(
		func(ctx context.Context, input runCliCommandInput) (string, error) {
			payload, _ := json.Marshal(input)
			return InvokeRunCliCommand(ctx, runner, payload, func(result string) {
				if emitter != nil && sessionID != 0 {
					emitter.EmitToolResult(sessionID, RunCliCommandToolName, result)
				}
			})
		},
		function.WithName(RunCliCommandToolName),
		function.WithDescription(meta.Description),
	)
}
