package tools

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
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
	return "Shell 工具"
}

func (s *ShellTool) Description() string {
	return "执行一次非交互式 shell 命令。每次执行前都必须获得用户确认。"
}

func (s *ShellTool) Tool() tool.BaseTool {
	return utils.NewTool(
		&schema.ToolInfo{
			Name: "shell_tool",
			Desc: "执行一次非交互式 shell 命令，返回 stdout、stderr 和 exit_code。调用前系统会请求用户确认。",
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"command": {
					Type:     schema.String,
					Desc:     "要执行的 shell 命令。",
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

			cmd := exec.CommandContext(execCtx, "zsh", "-lc", command)
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
	return &tool_approval.ApprovalPrompt{
		Title: "Shell 工具请求执行命令",
		Message: fmt.Sprintf(
			"工具想要执行以下命令：\n\n`%s`\n\n工作目录：`%s`\n范围：%s\n",
			command,
			dir,
			scope,
		),
		Scope: fmt.Sprintf("%s：%s", scope, dir),
	}, nil
}
