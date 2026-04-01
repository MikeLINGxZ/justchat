package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/i18n"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/tool_approval"
)

type FileTool struct{}

type fileToolParams struct {
	Operation string `json:"operation"`
	Path      string `json:"path"`
	Content   string `json:"content"`
}

func (f *FileTool) Id() string {
	return "file_tool"
}

func (f *FileTool) Name() string {
	return i18n.TCurrent("tool.file.name", nil)
}

func (f *FileTool) Description() string {
	return i18n.TCurrent("tool.file.description", nil)
}

func (f *FileTool) RequireConfirmation() bool { return true }

func (f *FileTool) Tool() tool.BaseTool {
	return utils.NewTool(
		&schema.ToolInfo{
			Name: "file_tool",
			Desc: i18n.TCurrent("tool.file.description", nil),
			ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
				"operation": {
					Type:     schema.String,
					Desc:     "文件操作类型，可选值：read、write、delete",
					Enum:     []string{"read", "write", "delete"},
					Required: true,
				},
				"path": {
					Type:     schema.String,
					Desc:     "目标文件路径。相对路径会基于当前工作区根目录解析。",
					Required: true,
				},
				"content": {
					Type: schema.String,
					Desc: "写入内容。仅在 operation=write 时使用。",
				},
			}),
		},
		func(ctx context.Context, params fileToolParams) (string, error) {
			targetPath, err := tool_approval.ResolvePath(params.Path)
			if err != nil {
				return "", err
			}

			switch strings.ToLower(strings.TrimSpace(params.Operation)) {
			case "read":
				content, err := os.ReadFile(targetPath)
				if err != nil {
					return "", err
				}
				return mustMarshal(map[string]interface{}{
					"success":   true,
					"operation": "read",
					"path":      targetPath,
					"bytes":     len(content),
					"content":   string(content),
				}), nil
			case "write":
				if err := os.WriteFile(targetPath, []byte(params.Content), 0o644); err != nil {
					return "", err
				}
				return mustMarshal(map[string]interface{}{
					"success":   true,
					"operation": "write",
					"path":      targetPath,
					"bytes":     len([]byte(params.Content)),
				}), nil
			case "delete":
				info, err := os.Stat(targetPath)
				if err != nil {
					return "", err
				}
				if info.IsDir() {
					return "", fmt.Errorf("delete only supports files, got directory: %s", targetPath)
				}
				if err := os.Remove(targetPath); err != nil {
					return "", err
				}
				return mustMarshal(map[string]interface{}{
					"success":   true,
					"operation": "delete",
					"path":      targetPath,
				}), nil
			default:
				return "", fmt.Errorf("unsupported file operation: %s", params.Operation)
			}
		},
	)
}

func (f *FileTool) BuildApprovalPrompt(ctx context.Context, argumentsJSON string) (*tool_approval.ApprovalPrompt, error) {
	var params fileToolParams
	if err := json.Unmarshal([]byte(argumentsJSON), &params); err != nil {
		return nil, err
	}

	targetPath, err := tool_approval.ResolvePath(params.Path)
	if err != nil {
		return nil, err
	}

	operation := strings.ToLower(strings.TrimSpace(params.Operation))
	scope := tool_approval.DescribeScope(targetPath)
	label := fileOperationLabel(operation)
	title := i18n.TCurrent("tool.file.approval.title", map[string]string{"operation": label})
	message := i18n.TCurrent("tool.file.approval.message", map[string]string{
		"operation": label,
		"path":      targetPath,
		"scope":     scope,
	})
	if operation == "write" && strings.TrimSpace(params.Content) != "" {
		preview := params.Content
		if len([]rune(preview)) > 240 {
			preview = string([]rune(preview)[:240]) + "..."
		}
		message += i18n.TCurrent("tool.file.approval.preview", map[string]string{"content": preview})
	}

	return &tool_approval.ApprovalPrompt{
		Title:   title,
		Message: message,
		Scope:   fmt.Sprintf("%s: %s", scope, filepath.Clean(targetPath)),
	}, nil
}

func fileOperationLabel(operation string) string {
	switch operation {
	case "read":
		return i18n.TCurrent("tool.file.operation.read", nil)
	case "write":
		return i18n.TCurrent("tool.file.operation.write", nil)
	case "delete":
		return i18n.TCurrent("tool.file.operation.delete", nil)
	default:
		return operation
	}
}

func mustMarshal(payload interface{}) string {
	bytes, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"success":false,"error":%q}`, err.Error())
	}
	return string(bytes)
}
