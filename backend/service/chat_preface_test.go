package service

import (
	"testing"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
)

func TestPreserveWorkflowPrefaceCopiesStreamingOutputOnce(t *testing.T) {
	message := data_models.Message{
		Content:          "先给你一个摘要",
		ReasoningContent: "我先快速分析一下",
		AssistantMessageExtra: &data_models.AssistantMessageExtra{
			RouteType: "",
		},
	}

	preserveWorkflowPreface(&message)

	if message.AssistantMessageExtra == nil {
		t.Fatal("AssistantMessageExtra = nil, want non-nil")
	}
	if message.AssistantMessageExtra.PrefaceContent != "先给你一个摘要" {
		t.Fatalf("PrefaceContent = %q, want copied content", message.AssistantMessageExtra.PrefaceContent)
	}
	if message.AssistantMessageExtra.PrefaceReasoningContent != "我先快速分析一下" {
		t.Fatalf("PrefaceReasoningContent = %q, want copied reasoning", message.AssistantMessageExtra.PrefaceReasoningContent)
	}

	message.Content = "新的中间输出"
	message.ReasoningContent = "新的思考"
	preserveWorkflowPreface(&message)

	if message.AssistantMessageExtra.PrefaceContent != "先给你一个摘要" {
		t.Fatalf("PrefaceContent overwritten to %q, want original value", message.AssistantMessageExtra.PrefaceContent)
	}
	if message.AssistantMessageExtra.PrefaceReasoningContent != "我先快速分析一下" {
		t.Fatalf("PrefaceReasoningContent overwritten to %q, want original value", message.AssistantMessageExtra.PrefaceReasoningContent)
	}
}

func TestResetDirectAssistantStatePreservesPrefaceFields(t *testing.T) {
	message := data_models.Message{
		Content:          "最终要被清空的正文",
		ReasoningContent: "最终要被清空的思考",
		AssistantMessageExtra: &data_models.AssistantMessageExtra{
			RouteType:               data_models.RouteTypeWorkflow,
			CurrentStage:            "任务交付",
			CurrentAgent:            "MainRouterAgent",
			PrefaceContent:          "保留的前置草稿",
			PrefaceReasoningContent: "保留的前置思考",
			FinishError:             "some error",
			ExecutionTrace: data_models.ExecutionTrace{
				Steps: []data_models.TraceStep{{StepID: "step-1"}},
			},
		},
	}

	resetDirectAssistantState(&message)

	if message.Content != "" {
		t.Fatalf("Content = %q, want empty", message.Content)
	}
	if message.ReasoningContent != "" {
		t.Fatalf("ReasoningContent = %q, want empty", message.ReasoningContent)
	}
	if message.AssistantMessageExtra.PrefaceContent != "保留的前置草稿" {
		t.Fatalf("PrefaceContent = %q, want preserved", message.AssistantMessageExtra.PrefaceContent)
	}
	if message.AssistantMessageExtra.PrefaceReasoningContent != "保留的前置思考" {
		t.Fatalf("PrefaceReasoningContent = %q, want preserved", message.AssistantMessageExtra.PrefaceReasoningContent)
	}
	if message.AssistantMessageExtra.CurrentStage != "" || message.AssistantMessageExtra.CurrentAgent != "" {
		t.Fatalf("CurrentStage/CurrentAgent = %q/%q, want empty", message.AssistantMessageExtra.CurrentStage, message.AssistantMessageExtra.CurrentAgent)
	}
	if len(message.AssistantMessageExtra.ExecutionTrace.Steps) != 0 {
		t.Fatalf("ExecutionTrace steps len = %d, want 0", len(message.AssistantMessageExtra.ExecutionTrace.Steps))
	}
	if message.AssistantMessageExtra.FinishError != "" {
		t.Fatalf("FinishError = %q, want empty", message.AssistantMessageExtra.FinishError)
	}
}
