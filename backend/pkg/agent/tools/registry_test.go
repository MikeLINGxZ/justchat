// backend/pkg/agent/tools/registry_test.go
package tools

import (
	"encoding/json"
	"testing"
)

func TestRegistry_RegisterAndGet(t *testing.T) {
	r := NewRegistry()
	meta := ToolMeta{
		Name:            "test_tool",
		Description:     "A test tool",
		Category:        CategoryBuiltin,
		RequiresConfirm: false,
		FormatPurpose:   func(args json.RawMessage) string { return "testing" },
	}
	r.Register(meta)

	got, ok := r.Get("test_tool")
	if !ok {
		t.Fatal("expected tool to be found")
	}
	if got.Name != "test_tool" {
		t.Fatalf("expected name test_tool, got %s", got.Name)
	}
}

func TestRegistry_BuiltinAndUserTools(t *testing.T) {
	r := NewRegistry()
	r.Register(ToolMeta{Name: "builtin1", Category: CategoryBuiltin})
	r.Register(ToolMeta{Name: "user1", Category: CategoryUser})
	r.Register(ToolMeta{Name: "user2", Category: CategoryUser})

	builtins := r.BuiltinTools()
	if len(builtins) != 1 {
		t.Fatalf("expected 1 builtin, got %d", len(builtins))
	}

	userTools := r.UserTools()
	if len(userTools) != 2 {
		t.Fatalf("expected 2 user tools, got %d", len(userTools))
	}
}

func TestRegistry_EnabledTools(t *testing.T) {
	r := NewRegistry()
	r.Register(ToolMeta{Name: "builtin1", Category: CategoryBuiltin})
	r.Register(ToolMeta{Name: "user1", Category: CategoryUser})
	r.Register(ToolMeta{Name: "user2", Category: CategoryUser})

	enabled := r.EnabledTools([]string{"user1"})
	if len(enabled) != 2 {
		t.Fatalf("expected 2 (1 builtin + 1 user enabled), got %d", len(enabled))
	}

	names := make(map[string]bool)
	for _, m := range enabled {
		names[m.Name] = true
	}
	if !names["builtin1"] || !names["user1"] {
		t.Fatalf("expected builtin1 and user1, got %v", names)
	}
}
