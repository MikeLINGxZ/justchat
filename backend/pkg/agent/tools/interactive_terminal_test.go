package tools

import (
	"context"
	"encoding/json"
	"testing"
)

func TestInteractiveTerminalToolNormalizesStatus(t *testing.T) {
	tool := NewInteractiveTerminalTool()
	raw, err := tool.Call(context.Background(), mustJSON(t, interactiveTerminalInput{
		Status: "waiting",
		Output: "scan to login",
	}))
	if err != nil {
		t.Fatalf("Call: %v", err)
	}
	out, ok := raw.(interactiveTerminalOutput)
	if !ok {
		t.Fatalf("expected interactiveTerminalOutput, got %T", raw)
	}
	if !out.InteractiveTerminal || out.TerminalStatus != "active" || out.TerminalOutput != "scan to login" {
		t.Fatalf("unexpected active output: %+v", out)
	}

	raw, err = tool.Call(context.Background(), mustJSON(t, interactiveTerminalInput{Status: "done"}))
	if err != nil {
		t.Fatalf("Call done: %v", err)
	}
	out, ok = raw.(interactiveTerminalOutput)
	if !ok {
		t.Fatalf("expected interactiveTerminalOutput, got %T", raw)
	}
	if !out.InteractiveTerminal || out.TerminalStatus != "done" || out.TerminalOutput != "" {
		t.Fatalf("unexpected done output: %+v", out)
	}
}

func mustJSON(t *testing.T, value any) []byte {
	t.Helper()
	data, err := json.Marshal(value)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return data
}
