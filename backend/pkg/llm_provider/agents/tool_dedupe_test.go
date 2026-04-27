package agents

import (
	"context"
	"testing"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

type namedTestTool struct {
	name string
}

func (n namedTestTool) Info(context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name:        n.name,
		Desc:        "test tool",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{}),
	}, nil
}

func (n namedTestTool) InvokableRun(context.Context, string, ...tool.Option) (string, error) {
	return "ok", nil
}

func TestUniqueToolsByInfoNameKeepsFirstTool(t *testing.T) {
	tools := []tool.BaseTool{
		namedTestTool{name: "load_skill"},
		namedTestTool{name: "file_tool"},
		namedTestTool{name: "load_skill"},
	}

	unique, err := uniqueToolsByInfoName(context.Background(), tools)
	if err != nil {
		t.Fatalf("uniqueToolsByInfoName returned error: %v", err)
	}

	if len(unique) != 2 {
		t.Fatalf("expected 2 unique tools, got %d", len(unique))
	}

	info, err := unique[0].Info(context.Background())
	if err != nil {
		t.Fatalf("read first tool info: %v", err)
	}
	if info.Name != "load_skill" {
		t.Fatalf("expected first tool to be load_skill, got %s", info.Name)
	}
}
