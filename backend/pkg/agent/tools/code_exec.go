package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

const codeExecTimeout = 30 * time.Second

type codeExecInput struct {
	Language string `json:"language" jsonschema:"description=Programming language (python/javascript/bash),required"`
	Code     string `json:"code" jsonschema:"description=Source code to execute,required"`
}

type codeExecOutput struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exit_code"`
}

func codeExecFunc(ctx context.Context, input codeExecInput) (codeExecOutput, error) {
	runners := map[string]struct {
		cmd string
		ext string
	}{
		"python":     {cmd: "python3", ext: ".py"},
		"javascript": {cmd: "node", ext: ".js"},
		"bash":       {cmd: "bash", ext: ".sh"},
	}

	runner, ok := runners[input.Language]
	if !ok {
		return codeExecOutput{}, fmt.Errorf("unsupported language: %s", input.Language)
	}

	tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("lemontea_exec_%d%s", time.Now().UnixNano(), runner.ext))
	if err := os.WriteFile(tmpFile, []byte(input.Code), 0o644); err != nil {
		return codeExecOutput{}, fmt.Errorf("write temp file: %w", err)
	}
	defer os.Remove(tmpFile)

	ctx, cancel := context.WithTimeout(ctx, codeExecTimeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, runner.cmd, tmpFile)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			return codeExecOutput{}, fmt.Errorf("exec: %w", err)
		}
	}

	return codeExecOutput{
		Stdout:   stdout.String(),
		Stderr:   stderr.String(),
		ExitCode: exitCode,
	}, nil
}

func NewCodeExecTool() *function.FunctionTool[codeExecInput, codeExecOutput] {
	return function.NewFunctionTool(
		codeExecFunc,
		function.WithName("code_exec"),
		function.WithDescription("Execute code in Python, JavaScript, or Bash and return the output"),
	)
}

func CodeExecMeta() ToolMeta {
	return ToolMeta{
		Name:            "code_exec",
		Description:     "Execute code in Python, JavaScript, or Bash",
		Category:        CategoryUser,
		RequiresConfirm: true,
		FormatPurpose: func(args json.RawMessage) string {
			var input codeExecInput
			_ = json.Unmarshal(args, &input)
			preview := input.Code
			if len(preview) > 100 {
				preview = preview[:100] + "..."
			}
			return fmt.Sprintf("Execute %s code:\n%s", input.Language, preview)
		},
	}
}
