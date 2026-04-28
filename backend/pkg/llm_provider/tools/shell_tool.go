package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/i18n"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/tool_approval"
)

type ShellTool struct{}

type shellToolParams struct {
	Command          string `json:"command"`
	WorkingDirectory string `json:"working_directory"`
	TimeoutSeconds   int    `json:"timeout_seconds"`
}

func (s *ShellTool) Id() string {
	return "shell_tool"
}

func (s *ShellTool) Name() string {
	return i18n.TCurrent("tool.shell.name", nil)
}

func (s *ShellTool) Description() string {
	return shellToolDescription()
}

func (s *ShellTool) RequireConfirmation() bool { return true }

func (s *ShellTool) Tool() tool.BaseTool {
	shellSpec := currentShellSpec()
	return utils.NewTool(
		&schema.ToolInfo{
			Name: "shell_tool",
			Desc: shellToolDescription(),
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"command": {
					Type:     schema.String,
					Desc:     shellSpec.CommandDescription,
					Required: true,
				},
				"working_directory": {
					Type: schema.String,
					Desc: "命令执行目录。相对路径会基于当前工作区根目录解析，默认为工作区根目录。",
				},
				"timeout_seconds": {
					Type: schema.Integer,
					Desc: "命令超时时间（秒），默认 30，最大 120。",
				},
			}),
		},
		func(ctx context.Context, params shellToolParams) (string, error) {
			command := strings.TrimSpace(params.Command)
			if command == "" {
				return "", fmt.Errorf("command is empty")
			}

			timeoutSeconds := params.TimeoutSeconds
			if timeoutSeconds <= 0 {
				timeoutSeconds = 30
			}
			if timeoutSeconds > 120 {
				timeoutSeconds = 120
			}

			dir := tool_approval.ResolveWorkingDirectory(params.WorkingDirectory)
			execCtx, cancel := context.WithTimeout(ctx, time.Duration(timeoutSeconds)*time.Second)
			defer cancel()

			cmd := exec.CommandContext(execCtx, shellSpec.Executable, append(shellSpec.Args, command)...)
			cmd.Dir = dir
			var stdout bytes.Buffer
			var stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			exitCode := 0
			err := cmd.Run()
			if err != nil {
				if exitErr, ok := err.(*exec.ExitError); ok {
					exitCode = exitErr.ExitCode()
				} else {
					return "", err
				}
			}

			return mustMarshal(map[string]interface{}{
				"success":           true,
				"command":           command,
				"shell":             shellSpec.Name,
				"working_directory": dir,
				"stdout":            stdout.String(),
				"stderr":            stderr.String(),
				"exit_code":         exitCode,
			}), nil
		},
	)
}

func (s *ShellTool) BuildApprovalPrompt(ctx context.Context, argumentsJSON string) (*tool_approval.ApprovalPrompt, error) {
	var params shellToolParams
	if err := json.Unmarshal([]byte(argumentsJSON), &params); err != nil {
		return nil, err
	}

	command := strings.TrimSpace(params.Command)
	if command == "" {
		return nil, fmt.Errorf("command is empty")
	}

	dir := tool_approval.ResolveWorkingDirectory(params.WorkingDirectory)
	scope := tool_approval.DescribeScope(dir)
	shellSpec := currentShellSpec()
	message := i18n.TCurrent("tool.shell.approval.message", map[string]string{
		"command":   command,
		"directory": dir,
		"scope":     scope,
		"shell":     shellSpec.DisplayName,
	})
	message = strings.TrimRight(message, "\n") + "\n\n" + i18n.TCurrent("tool.approval.actions", nil)

	return &tool_approval.ApprovalPrompt{
		Title:   fmt.Sprintf("%s（%s）", i18n.TCurrent("tool.shell.approval.title", nil), shellSpec.DisplayName),
		Message: message,
		Scope:   fmt.Sprintf("%s: %s", scope, dir),
	}, nil
}

type shellSpec struct {
	Name               string
	DisplayName        string
	Executable         string
	Args               []string
	Description        string
	CommandDescription string
	Instruction        string
}

func currentShellSpec() shellSpec {
	if runtime.GOOS == "windows" {
		return shellSpec{
			Name:               "powershell",
			DisplayName:        "PowerShell",
			Executable:         "powershell.exe",
			Args:               []string{"-NoProfile", "-NonInteractive", "-ExecutionPolicy", "Bypass", "-Command"},
			Description:        "执行一次非交互式 PowerShell 命令。当前运行环境是 Windows，请使用 PowerShell 语法，例如 Get-ChildItem、Select-String、;。",
			CommandDescription: "要执行的 PowerShell 命令。当前运行环境是 Windows，请使用 PowerShell 语法，例如 Get-ChildItem、Select-String、;，不要使用 bash/zsh/POSIX 语法，除非用户明确提供了可用的兼容环境。",
			Instruction:        "Shell 工具规则：调用 shell_tool 时必须遵循工具描述中的当前 OS 命令语法。当前环境是 Windows，使用 PowerShell 语法；不要使用 bash/zsh/POSIX 管道习惯，除非用户明确提供了可用的兼容环境。",
		}
	}
	return shellSpec{
		Name:               "zsh",
		DisplayName:        "zsh",
		Executable:         "zsh",
		Args:               []string{"-lc"},
		Description:        "执行一次非交互式 zsh 命令。当前运行环境是 macOS/Linux，请使用 POSIX/zsh 语法。",
		CommandDescription: "要执行的 zsh 命令。当前运行环境是 macOS/Linux，请使用 POSIX/zsh 语法。",
		Instruction:        "Shell 工具规则：调用 shell_tool 时必须遵循工具描述中的当前 OS 命令语法。当前环境是 macOS/Linux，使用 POSIX/zsh 语法。",
	}
}

func shellToolDescription() string {
	return currentShellSpec().Description
}

func ShellRuntimeInstruction() string {
	return currentShellSpec().Instruction
}
