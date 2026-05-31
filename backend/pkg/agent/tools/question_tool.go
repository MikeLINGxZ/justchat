package tools

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"trpc.group/trpc-go/trpc-agent-go/tool/function"
)

const QuestionToolName = "question"

type questionOption struct {
	Label       string `json:"label" jsonschema:"description=Short option label,required"`
	Description string `json:"description" jsonschema:"description=What choosing this option means"`
}

type questionInput struct {
	Question string           `json:"question" jsonschema:"description=Clear user-facing question to ask,required"`
	Context  string           `json:"context" jsonschema:"description=Brief reason why the question is needed"`
	Options  []questionOption `json:"options" jsonschema:"description=Optional choices for the user"`
}

type questionOutput struct {
	Question    string           `json:"question"`
	Context     string           `json:"context,omitempty"`
	Options     []questionOption `json:"options,omitempty"`
	Instruction string           `json:"instruction"`
}

// BuildQuestionTool returns registry metadata for the clarification tool.
func BuildQuestionTool() ToolMeta {
	return ToolMeta{
		Name:        QuestionToolName,
		Description: "Ask the user a clear clarification question when the request is incomplete, risky, or needs a preference choice.",
		Category:    CategoryBuiltin,
		FormatPurpose: func(_ json.RawMessage) string {
			return "Ask the user a clarification question"
		},
	}
}

// questionFunc validates and formats a user-facing clarification request.
func questionFunc(_ context.Context, input questionInput) (questionOutput, error) {
	question := strings.TrimSpace(input.Question)
	if question == "" {
		return questionOutput{}, errors.New("question is required")
	}
	return questionOutput{
		Question:    question,
		Context:     strings.TrimSpace(input.Context),
		Options:     input.Options,
		Instruction: "Ask the user this question directly, then wait for their answer before continuing.",
	}, nil
}

// NewQuestionTool creates the function tool used to structure user clarifications.
func NewQuestionTool() *function.FunctionTool[questionInput, questionOutput] {
	meta := BuildQuestionTool()
	return function.NewFunctionTool(
		questionFunc,
		function.WithName(QuestionToolName),
		function.WithDescription(meta.Description),
	)
}
