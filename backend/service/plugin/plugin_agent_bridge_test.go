package plugin

import (
	"encoding/json"
	"strings"
	"testing"

	pkgcli "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/cli"
)

func TestFormatCLIProgressMarksStreamingOutputAsInteractiveTerminal(t *testing.T) {
	out := formatCLIProgress(pkgcli.RunProgress{
		Stdout: "QR",
		Stderr: "scan to continue",
	})

	var payload struct {
		InteractiveTerminal bool   `json:"interactive_terminal"`
		TerminalStatus      string `json:"terminal_status"`
		TerminalOutput      string `json:"terminal_output"`
	}
	if err := json.Unmarshal([]byte(out), &payload); err != nil {
		t.Fatalf("progress should be JSON: %v; out=%q", err, out)
	}
	if !payload.InteractiveTerminal {
		t.Fatalf("expected interactive terminal payload: %+v", payload)
	}
	if payload.TerminalStatus != "active" {
		t.Fatalf("terminal_status=%q", payload.TerminalStatus)
	}
	if !strings.Contains(payload.TerminalOutput, "QR") || !strings.Contains(payload.TerminalOutput, "scan to continue") {
		t.Fatalf("terminal_output=%q", payload.TerminalOutput)
	}
}
