package storage

import (
	"testing"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
)

// TestTerminalRoundTripAndOutputCursor verifies terminal rows and output cursors persist correctly.
func TestTerminalRoundTripAndOutputCursor(t *testing.T) {
	stor := newTestStorage(t)

	term, err := stor.CreateTerminal(data_models.Terminal{
		TerminalID:    "term_test",
		SessionID:     7,
		ToolCallID:    "call_1",
		Title:         "Waiting",
		Command:       "bash",
		Args:          "[]",
		Cwd:           "/tmp",
		Status:        "active",
		Visible:       true,
		CurrentCursor: 0,
	})
	if err != nil {
		t.Fatal(err)
	}
	if term.ID == 0 {
		t.Fatal("expected persisted terminal id")
	}

	if _, err := stor.AppendTerminalOutput("term_test", "hello"); err != nil {
		t.Fatal(err)
	}
	if _, err := stor.AppendTerminalOutput("term_test", " world"); err != nil {
		t.Fatal(err)
	}

	chunks, err := stor.ReadTerminalOutput("term_test", 5)
	if err != nil {
		t.Fatal(err)
	}
	if len(chunks) != 1 {
		t.Fatalf("expected one chunk after cursor 5, got %d", len(chunks))
	}
	if chunks[0].CursorStart != 5 || chunks[0].CursorEnd != 11 || chunks[0].Data != " world" {
		t.Fatalf("unexpected chunk: %+v", chunks[0])
	}

	listed, err := stor.ListTerminalsForSession(7)
	if err != nil {
		t.Fatal(err)
	}
	if len(listed) != 1 {
		t.Fatalf("expected one terminal, got %d", len(listed))
	}
	if listed[0].CurrentCursor != 11 {
		t.Fatalf("cursor not updated: %+v", listed[0])
	}
}
