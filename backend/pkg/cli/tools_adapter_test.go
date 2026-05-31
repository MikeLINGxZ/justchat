package cli

import (
	"context"
	"encoding/json"
	"path/filepath"
	"testing"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	toolpkg "trpc.group/trpc-go/trpc-agent-go/tool"
)

type fakeToolRunner struct {
	lastName     string
	lastToolName string
	lastInput    map[string]any
	result       RunResult
	err          error
}

// RunTool records the forwarded call for assertions.
func (f *fakeToolRunner) RunTool(ctx context.Context, name string, toolName string, input map[string]any) (RunResult, error) {
	f.lastName = name
	f.lastToolName = toolName
	f.lastInput = input
	return f.result, f.err
}

// TestBuildToolSetFiltersEnabledTools verifies only manifest-enabled tools that match the conversation filter are exposed.
func TestBuildToolSetFiltersEnabledTools(t *testing.T) {
	t.Helper()

	item := data_models.ExtensionItem{ID: "cli:demo", Name: "demo-cli", RootDir: filepath.Join("/tmp", "demo")}
	manifest := Manifest{
		Executable: "/tmp/demo",
		Tools: []Tool{
			{
				Name:         "list_items",
				Description:  "list items",
				InputSchema:  json.RawMessage(`{"type":"object","properties":{"query":{"type":"string"}}}`),
				ArgvTemplate: []string{"list", "--query", "{query}"},
				OutputMode:   OutputJSON,
				Enabled:      true,
			},
			{
				Name:         "delete_item",
				Description:  "delete item",
				InputSchema:  json.RawMessage(`{"type":"object","properties":{"id":{"type":"string"}}}`),
				ArgvTemplate: []string{"delete", "{id}"},
				OutputMode:   OutputText,
				Enabled:      false,
			},
		},
	}
	runner := &fakeToolRunner{}

	set, err := buildToolSetFromRunner(runner, item, manifest, []string{"cli_demo_cli_list_items"})
	if err != nil {
		t.Fatalf("buildToolSetFromRunner returned error: %v", err)
	}
	defer func() { _ = set.Close() }()

	tools := set.Tools(context.Background())
	if len(tools) != 1 {
		t.Fatalf("expected 1 tool, got %d", len(tools))
	}
	// Declaration name is the raw tool name; NamedToolSet wrapping adds the `<setName>_` prefix at runtime.
	if tools[0].Declaration().Name != "list_items" {
		t.Fatalf("unexpected tool name: %q", tools[0].Declaration().Name)
	}
	if set.Name() != "cli_demo_cli" {
		t.Fatalf("unexpected toolset name: %q", set.Name())
	}
}

// TestBuildToolSetForwardsCalls verifies tool invocation is proxied into cli.Manager.RunTool and returned as JSON.
func TestBuildToolSetForwardsCalls(t *testing.T) {
	t.Helper()

	item := data_models.ExtensionItem{ID: "cli:demo", Name: "demo-cli", RootDir: filepath.Join("/tmp", "demo-runtime")}
	manifest := Manifest{
		Executable: "/tmp/demo",
		Tools: []Tool{
			{
				Name:         "list_items",
				Description:  "list items",
				InputSchema:  json.RawMessage(`{"type":"object","properties":{"query":{"type":"string"}},"required":["query"]}`),
				ArgvTemplate: []string{"list", "--query", "{query}"},
				OutputMode:   OutputJSON,
				Enabled:      true,
			},
		},
	}
	runner := &fakeToolRunner{
		result: RunResult{
			ExitCode:    0,
			Stdout:      `{"ok":true}`,
			Parsed:      json.RawMessage(`{"ok":true}`),
			DurationMS:  18,
			ParsedLines: []string{},
		},
	}

	set, err := buildToolSetFromRunner(runner, item, manifest, []string{"cli:demo"})
	if err != nil {
		t.Fatalf("buildToolSetFromRunner returned error: %v", err)
	}
	defer func() { _ = set.Close() }()

	tools := set.Tools(context.Background())
	callable, ok := tools[0].(toolpkg.CallableTool)
	if !ok {
		t.Fatalf("expected callable tool, got %T", tools[0])
	}
	result, err := callable.Call(context.Background(), []byte(`{"query":"hello"}`))
	if err != nil {
		t.Fatalf("tool call returned error: %v", err)
	}

	if runner.lastName != "demo-runtime" || runner.lastToolName != "list_items" {
		t.Fatalf("unexpected forwarded target: name=%q tool=%q", runner.lastName, runner.lastToolName)
	}
	if runner.lastInput["query"] != "hello" {
		t.Fatalf("unexpected forwarded input: %+v", runner.lastInput)
	}

	bytes, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("marshal tool result: %v", err)
	}
	if string(bytes) == "{}" {
		t.Fatalf("expected structured result payload, got %s", string(bytes))
	}
}
