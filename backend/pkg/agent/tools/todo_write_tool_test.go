package tools

import (
	"context"
	"testing"
)

func TestTodoWriteToolSummarizesProgress(t *testing.T) {
	out, err := todoWriteFunc(context.Background(), todoWriteInput{
		Items: []todoWriteItem{
			{Title: "Scan folder", Status: "completed"},
			{Title: "Group PDFs", Status: "in_progress"},
			{Title: "Write report", Status: "pending"},
		},
	})
	if err != nil {
		t.Fatalf("todo write func: %v", err)
	}
	if out.Total != 3 || out.Completed != 1 || out.InProgress != 1 || out.Pending != 1 {
		t.Fatalf("unexpected summary: %+v", out)
	}
	if out.Summary == "" {
		t.Fatal("expected user-readable summary")
	}
}

func TestTodoWriteToolRejectsEmptyList(t *testing.T) {
	_, err := todoWriteFunc(context.Background(), todoWriteInput{})
	if err == nil {
		t.Fatal("expected at least one todo item")
	}
}
