package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

type fileReadInput struct {
	Path string `json:"path" jsonschema:"description=Absolute file path to read,required"`
}

type fileReadOutput struct {
	Content string `json:"content"`
	Size    int64  `json:"size"`
}

func fileReadFunc(ctx context.Context, input fileReadInput) (fileReadOutput, error) {
	content, err := os.ReadFile(input.Path)
	if err != nil {
		return fileReadOutput{}, fmt.Errorf("read file: %w", err)
	}
	info, _ := os.Stat(input.Path)
	var size int64
	if info != nil {
		size = info.Size()
	}
	return fileReadOutput{Content: string(content), Size: size}, nil
}

func NewFileReadTool() *function.FunctionTool[fileReadInput, fileReadOutput] {
	return function.NewFunctionTool(
		fileReadFunc,
		function.WithName("file_read"),
		function.WithDescription("Read the content of a file at the given path"),
	)
}

type fileWriteInput struct {
	Path    string `json:"path" jsonschema:"description=Absolute file path to write,required"`
	Content string `json:"content" jsonschema:"description=Content to write,required"`
}

type fileWriteOutput struct {
	BytesWritten int `json:"bytes_written"`
}

func fileWriteFunc(ctx context.Context, input fileWriteInput) (fileWriteOutput, error) {
	if err := os.MkdirAll(filepath.Dir(input.Path), 0o755); err != nil {
		return fileWriteOutput{}, fmt.Errorf("create dir: %w", err)
	}
	if err := os.WriteFile(input.Path, []byte(input.Content), 0o644); err != nil {
		return fileWriteOutput{}, fmt.Errorf("write file: %w", err)
	}
	return fileWriteOutput{BytesWritten: len(input.Content)}, nil
}

func NewFileWriteTool() *function.FunctionTool[fileWriteInput, fileWriteOutput] {
	return function.NewFunctionTool(
		fileWriteFunc,
		function.WithName("file_write"),
		function.WithDescription("Write content to a file at the given path"),
	)
}

func FileReadMeta() ToolMeta {
	return ToolMeta{
		Name:            "file_read",
		Description:     "Read file content",
		Category:        CategoryBuiltin,
		RequiresConfirm: true,
		FormatPurpose: func(args json.RawMessage) string {
			var input fileReadInput
			_ = json.Unmarshal(args, &input)
			return fmt.Sprintf("Read file: %s", input.Path)
		},
	}
}

func FileWriteMeta() ToolMeta {
	return ToolMeta{
		Name:            "file_write",
		Description:     "Write content to file",
		Category:        CategoryBuiltin,
		RequiresConfirm: true,
		FormatPurpose: func(args json.RawMessage) string {
			var input fileWriteInput
			_ = json.Unmarshal(args, &input)
			return fmt.Sprintf("Write to file: %s (%d bytes)", input.Path, len(input.Content))
		},
	}
}
