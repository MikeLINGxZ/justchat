package tools

import (
	"context"
	"strings"
	"testing"
)

func TestQuestionToolReturnsStructuredQuestion(t *testing.T) {
	out, err := questionFunc(context.Background(), questionInput{
		Question: "How should I organize these files?",
		Options: []questionOption{
			{Label: "By month", Description: "Group files by modified month"},
			{Label: "By type", Description: "Group files by extension"},
		},
	})
	if err != nil {
		t.Fatalf("question func: %v", err)
	}
	if out.Question != "How should I organize these files?" {
		t.Fatalf("unexpected question: %+v", out)
	}
	if len(out.Options) != 2 {
		t.Fatalf("expected options to be preserved, got %+v", out.Options)
	}
	if !strings.Contains(out.Instruction, "Ask the user") {
		t.Fatalf("expected user-facing instruction, got %q", out.Instruction)
	}
}

func TestQuestionToolRequiresQuestionText(t *testing.T) {
	_, err := questionFunc(context.Background(), questionInput{})
	if err == nil {
		t.Fatal("expected question text to be required")
	}
}
