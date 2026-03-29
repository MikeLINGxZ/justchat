package service

import (
	"strings"
	"testing"

	"github.com/cloudwego/eino/schema"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
)

func TestShouldForceWorkflow(t *testing.T) {
	tests := []struct {
		name       string
		message    data_models.Message
		wantForce  bool
		wantRule   string
		wantReason string
	}{
		{
			name: "attached files do not force workflow",
			message: data_models.Message{
				Content: "你好",
				UserMessageExtra: &data_models.UserMessageExtra{
					Files: []data_models.File{{Name: "a.pdf", Path: "/tmp/a.pdf"}},
				},
			},
			wantForce: false,
		},
		{
			name: "selected agents force workflow",
			message: data_models.Message{
				Content: "你好",
				UserMessageExtra: &data_models.UserMessageExtra{
					Agents: []string{"agent-a"},
				},
			},
			wantForce:  true,
			wantRule:   "selected_agents",
			wantReason: "检测到 agent 选择，优先进入任务编排",
		},
		{
			name: "selected tools no longer force workflow",
			message: data_models.Message{
				Content: "帮我阻塞39s",
				UserMessageExtra: &data_models.UserMessageExtra{
					Tools: []string{"block"},
				},
			},
			wantForce: false,
		},
		{
			name: "simple question with selected tools does not force workflow",
			message: data_models.Message{
				Content: "你好，你是什么模型？",
				UserMessageExtra: &data_models.UserMessageExtra{
					Tools: []string{"tool-a"},
				},
			},
			wantForce: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldForceWorkflow(tt.message)
			if got.Force != tt.wantForce {
				t.Fatalf("force = %t, want %t", got.Force, tt.wantForce)
			}
			if got.RuleName != tt.wantRule {
				t.Fatalf("rule = %q, want %q", got.RuleName, tt.wantRule)
			}
			if got.Reason != tt.wantReason {
				t.Fatalf("reason = %q, want %q", got.Reason, tt.wantReason)
			}
		})
	}
}

func TestBuildPlanningPromptRendersDynamicContext(t *testing.T) {
	template := "用户请求：{{user_request}}\n最近上下文：\n{{recent_context}}"
	messages := []schema.Message{
		{Role: schema.User, Content: "你好"},
		{Role: schema.Assistant, Content: "你好呀"},
	}

	got := buildPlanningPrompt(template, "帮我总结", messages)
	if got == "" {
		t.Fatal("buildPlanningPrompt() returned empty string")
	}
	if want := "用户请求：帮我总结"; !strings.Contains(got, want) {
		t.Fatalf("buildPlanningPrompt() = %q, want substring %q", got, want)
	}
	if want := "user: 你好"; !strings.Contains(got, want) {
		t.Fatalf("buildPlanningPrompt() = %q, want substring %q", got, want)
	}
	if want := "assistant: 你好呀"; !strings.Contains(got, want) {
		t.Fatalf("buildPlanningPrompt() = %q, want substring %q", got, want)
	}
}

func TestBuildSynthesisPromptRendersResults(t *testing.T) {
	template := "请求={{user_request}}\n结果=\n{{task_results}}\n反馈={{review_feedback}}"
	plan := workflowPlan{Goal: "完成任务", CompletionCriteria: []string{"准确"}}
	results := map[string]workflowTaskResult{
		"task_1": {TaskID: "task_1", AgentName: "GeneralWorkerAgent", Output: "已完成"},
	}

	got := buildSynthesisPrompt(template, "请总结", plan, results, "")
	if want := "[task_1][GeneralWorkerAgent]\n已完成"; !strings.Contains(got, want) {
		t.Fatalf("buildSynthesisPrompt() = %q, want substring %q", got, want)
	}
	if want := "反馈=无"; !strings.Contains(got, want) {
		t.Fatalf("buildSynthesisPrompt() = %q, want substring %q", got, want)
	}
}

func TestBuildReviewPromptRendersDraft(t *testing.T) {
	template := "目标={{goal}}\n结果=\n{{task_results}}\n草稿={{draft}}"
	plan := workflowPlan{Goal: "完成任务", CompletionCriteria: []string{"准确"}}
	results := map[string]workflowTaskResult{
		"task_1": {TaskID: "task_1", Output: "任务输出"},
	}

	got := buildReviewPrompt(template, "请总结", plan, results, "最终草稿")
	if want := "目标=完成任务"; !strings.Contains(got, want) {
		t.Fatalf("buildReviewPrompt() = %q, want substring %q", got, want)
	}
	if want := "task_1: 任务输出"; !strings.Contains(got, want) {
		t.Fatalf("buildReviewPrompt() = %q, want substring %q", got, want)
	}
	if want := "草稿=最终草稿"; !strings.Contains(got, want) {
		t.Fatalf("buildReviewPrompt() = %q, want substring %q", got, want)
	}
}
