package tools

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strings"
	"testing"

	einotool "github.com/cloudwego/eino/components/tool"
)

func TestFileToolReadWriteDelete(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, "sample.txt")
	fileTool := &FileTool{}
	invokable, ok := fileTool.Tool().(einotool.InvokableTool)
	if !ok {
		t.Fatal("file tool is not invokable")
	}

	writeResult, err := invokable.InvokableRun(context.Background(), `{"operation":"write","path":"`+target+`","content":"hello"}`)
	if err != nil {
		t.Fatalf("write error = %v", err)
	}
	if !strings.Contains(writeResult, `"operation": "write"`) {
		t.Fatalf("write result = %s, want write operation", writeResult)
	}

	readResult, err := invokable.InvokableRun(context.Background(), `{"operation":"read","path":"`+target+`"}`)
	if err != nil {
		t.Fatalf("read error = %v", err)
	}
	if !strings.Contains(readResult, `"content": "hello"`) {
		t.Fatalf("read result = %s, want content hello", readResult)
	}

	deleteResult, err := invokable.InvokableRun(context.Background(), `{"operation":"delete","path":"`+target+`"}`)
	if err != nil {
		t.Fatalf("delete error = %v", err)
	}
	if !strings.Contains(deleteResult, `"operation": "delete"`) {
		t.Fatalf("delete result = %s, want delete operation", deleteResult)
	}
}

func TestFileToolApprovalPromptIncludesScope(t *testing.T) {
	fileTool := &FileTool{}
	prompt, err := fileTool.BuildApprovalPrompt(context.Background(), `{"operation":"read","path":"README.md"}`)
	if err != nil {
		t.Fatalf("BuildApprovalPrompt() error = %v", err)
	}
	if !strings.Contains(prompt.Message, "工作区内") {
		t.Fatalf("prompt message = %q, want workspace scope", prompt.Message)
	}
	if !strings.Contains(prompt.Message, "1. 允许 2. 拒绝 3. 告诉ai应该怎么做") {
		t.Fatalf("prompt message = %q, want action choices", prompt.Message)
	}
}

func TestShellToolRunsCommand(t *testing.T) {
	shellTool := &ShellTool{}
	invokable, ok := shellTool.Tool().(einotool.InvokableTool)
	if !ok {
		t.Fatal("shell tool is not invokable")
	}

	result, err := invokable.InvokableRun(context.Background(), `{"command":"printf 'hello-shell'","timeout_seconds":5}`)
	if err != nil {
		t.Fatalf("InvokableRun() error = %v", err)
	}

	var payload map[string]interface{}
	if err := json.Unmarshal([]byte(result), &payload); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	if payload["stdout"] != "hello-shell" {
		t.Fatalf("stdout = %#v, want hello-shell", payload["stdout"])
	}
	if int(payload["exit_code"].(float64)) != 0 {
		t.Fatalf("exit_code = %#v, want 0", payload["exit_code"])
	}
}

func TestShellToolApprovalPromptIncludesCommand(t *testing.T) {
	shellTool := &ShellTool{}
	prompt, err := shellTool.BuildApprovalPrompt(context.Background(), `{"command":"pwd","working_directory":".."}`)
	if err != nil {
		t.Fatalf("BuildApprovalPrompt() error = %v", err)
	}
	if !strings.Contains(prompt.Message, "`pwd`") {
		t.Fatalf("prompt message = %q, want command preview", prompt.Message)
	}
	if !strings.Contains(prompt.Message, "1. 允许 2. 拒绝 3. 告诉ai应该怎么做") {
		t.Fatalf("prompt message = %q, want action choices", prompt.Message)
	}
}
