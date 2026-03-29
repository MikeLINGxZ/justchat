package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/llm_provider"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/llm_provider/agents"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/prompts"
)

type routeSource string

const (
	routeSourceGuardRule routeSource = "guard_rule"
	routeSourceMainModel routeSource = "main_model"
)

type workflowGuardDecision struct {
	Force    bool
	Reason   string
	RuleName string
}

type workflowPlan struct {
	Goal               string             `json:"goal"`
	CompletionCriteria []string           `json:"completion_criteria"`
	Tasks              []workflowPlanTask `json:"tasks"`
}

type workflowPlanTask struct {
	ID             string   `json:"id"`
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	Dependencies   []string `json:"dependencies"`
	SuggestedAgent string   `json:"suggested_agent"`
	RequiredTools  []string `json:"required_tools"`
	ExpectedOutput string   `json:"expected_output"`
}

type workflowTaskResult struct {
	TaskID    string   `json:"task_id"`
	Title     string   `json:"title"`
	AgentName string   `json:"agent_name"`
	Output    string   `json:"output"`
	UsedTools []string `json:"used_tools"`
}

type reviewDecision struct {
	Approved          bool     `json:"approved"`
	Issues            []string `json:"issues"`
	RetryInstructions string   `json:"retry_instructions"`
	AffectedTaskIDs   []string `json:"affected_task_ids"`
}

type traceContextKey string

const (
	traceParentStepIDContextKey traceContextKey = "trace_parent_step_id"
	traceAgentNameContextKey    traceContextKey = "trace_agent_name"
)

func shouldForceWorkflow(message data_models.Message) workflowGuardDecision {
	if message.UserMessageExtra != nil {
		if len(message.UserMessageExtra.Agents) > 0 {
			return workflowGuardDecision{
				Force:    true,
				Reason:   "检测到 agent 选择，优先进入任务编排",
				RuleName: "selected_agents",
			}
		}
	}

	content := strings.TrimSpace(message.Content)
	if content == "" {
		return workflowGuardDecision{}
	}

	return workflowGuardDecision{}
}

func generateWorkflowPlan(ctx context.Context, provider *llm_provider.Provider, userRequest string, messages []schema.Message, tools []tool.BaseTool) (workflowPlan, error) {
	toolNames := make([]string, 0, len(tools))
	for _, item := range tools {
		info, infoErr := item.Info(ctx)
		if infoErr != nil || info == nil {
			continue
		}
		toolNames = append(toolNames, info.Name)
	}
	sort.Strings(toolNames)
	systemPrompt := prompts.Render(provider.Prompts().PlannerSystem, map[string]string{
		"tool_names": strings.Join(toolNames, ", "),
	})

	plannerMessages := append([]schema.Message{{Role: schema.System, Content: systemPrompt}}, buildPlannerMessages(provider.Prompts().PlanningUser, userRequest, messages)...)
	resp, err := provider.Generate(ctx, plannerMessages)
	if err != nil {
		return workflowPlan{}, err
	}
	var plan workflowPlan
	if err := unmarshalJSONResponse(resp.Content, &plan); err != nil {
		return workflowPlan{}, err
	}
	if len(plan.Tasks) == 0 {
		plan.Tasks = []workflowPlanTask{{
			ID:             "task_1",
			Title:          "完成用户请求",
			Description:    userRequest,
			SuggestedAgent: "GeneralWorkerAgent",
			ExpectedOutput: "完整回答用户请求",
		}}
	}
	for idx := range plan.Tasks {
		if strings.TrimSpace(plan.Tasks[idx].ID) == "" {
			plan.Tasks[idx].ID = fmt.Sprintf("task_%d", idx+1)
		}
		if strings.TrimSpace(plan.Tasks[idx].SuggestedAgent) == "" {
			plan.Tasks[idx].SuggestedAgent = "GeneralWorkerAgent"
		}
	}
	if strings.TrimSpace(plan.Goal) == "" {
		plan.Goal = userRequest
	}
	if len(plan.CompletionCriteria) == 0 {
		plan.CompletionCriteria = []string{"回答覆盖用户目标", "内容清晰可执行"}
	}
	return plan, nil
}

func buildPlannerMessages(promptTemplate string, userRequest string, messages []schema.Message) []schema.Message {
	contextMessages := collectRecentConversationMessages(messages, 8)
	contextMessages = append(contextMessages, schema.Message{
		Role:    schema.User,
		Content: buildPlanningPrompt(promptTemplate, userRequest, messages),
	})
	return contextMessages
}

func collectRecentConversationMessages(messages []schema.Message, maxCount int) []schema.Message {
	if maxCount <= 0 {
		return nil
	}

	conversation := make([]schema.Message, 0, len(messages))
	for _, msg := range messages {
		if msg.Role != schema.User && msg.Role != schema.Assistant {
			continue
		}
		conversation = append(conversation, cloneWorkflowContextMessage(msg))
	}

	if len(conversation) > maxCount {
		conversation = append([]schema.Message(nil), conversation[len(conversation)-maxCount:]...)
	}
	return conversation
}

func findLatestUserContextMessage(messages []schema.Message) *schema.Message {
	for i := len(messages) - 1; i >= 0; i-- {
		msg := messages[i]
		if msg.Role != schema.User {
			continue
		}
		if len(msg.UserInputMultiContent) == 0 && strings.TrimSpace(msg.Content) == "" {
			continue
		}
		cloned := cloneWorkflowContextMessage(msg)
		return &cloned
	}
	return nil
}

func cloneWorkflowContextMessage(msg schema.Message) schema.Message {
	cloned := msg
	if len(msg.UserInputMultiContent) > 0 {
		cloned.UserInputMultiContent = append([]schema.MessageInputPart(nil), msg.UserInputMultiContent...)
	}
	if len(msg.ToolCalls) > 0 {
		cloned.ToolCalls = append([]schema.ToolCall(nil), msg.ToolCalls...)
	}
	return cloned
}

func buildPlanningPrompt(promptTemplate string, userRequest string, messages []schema.Message) string {
	var recent []string
	for _, msg := range messages {
		if msg.Role != schema.User && msg.Role != schema.Assistant {
			continue
		}
		recent = append(recent, fmt.Sprintf("%s: %s", msg.Role, compactText(msg.Content, 180)))
	}
	if len(recent) > 8 {
		recent = recent[len(recent)-8:]
	}
	recentContext := strings.Join(recent, "\n")
	if recentContext == "" {
		recentContext = "(无)"
	}
	return prompts.Render(promptTemplate, map[string]string{
		"user_request":   userRequest,
		"recent_context": recentContext,
	})
}

func executePlanTask(ctx context.Context, provider *llm_provider.Provider, task workflowPlanTask, plan workflowPlan, priorResults map[string]workflowTaskResult, originalUserMessage *schema.Message, tools []tool.BaseTool, toolMiddleware compose.ToolMiddleware, parentStepID string) (workflowTaskResult, error) {
	taskPrompt := buildWorkerPrompt(plan, task, priorResults)
	agentName := task.SuggestedAgent
	if agentName == "" {
		agentName = "GeneralWorkerAgent"
	}
	instruction := provider.Prompts().WorkerGeneralSystem
	if agentName == "ToolSpecialistAgent" {
		instruction = provider.Prompts().WorkerToolSystem
	}
	agent, err := agents.NewRoleAgent(ctx, provider.ToolCallingModel(), agentName, "工作子代理", instruction, tools, toolMiddleware)
	if err != nil {
		return workflowTaskResult{}, err
	}
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent:           agent,
		EnableStreaming: true,
	})
	runCtx := context.WithValue(ctx, traceParentStepIDContextKey, parentStepID)
	runCtx = context.WithValue(runCtx, traceAgentNameContextKey, agentName)
	workerMessages := buildWorkerMessages(taskPrompt, originalUserMessage)
	messagePointers := make([]*schema.Message, 0, len(workerMessages))
	for i := range workerMessages {
		messagePointers = append(messagePointers, &workerMessages[i])
	}
	iter := runner.Run(runCtx, messagePointers)

	var contentBuilder strings.Builder
	var usedTools []string
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			return workflowTaskResult{}, event.Err
		}
		if event.Output == nil || event.Output.MessageOutput == nil {
			continue
		}
		mo := event.Output.MessageOutput
		if mo.Role == schema.Tool && mo.Message != nil && mo.Message.ToolName != "" {
			usedTools = append(usedTools, mo.Message.ToolName)
			continue
		}
		if mo.Role != schema.Assistant {
			continue
		}
		if mo.IsStreaming && mo.MessageStream != nil {
			for {
				msg, err := mo.MessageStream.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					return workflowTaskResult{}, err
				}
				if msg != nil && msg.Content != "" {
					contentBuilder.WriteString(msg.Content)
				}
			}
			mo.MessageStream.Close()
			continue
		}
		if mo.Message != nil && mo.Message.Content != "" {
			contentBuilder.WriteString(mo.Message.Content)
		}
	}

	return workflowTaskResult{
		TaskID:    task.ID,
		Title:     task.Title,
		AgentName: agentName,
		Output:    strings.TrimSpace(contentBuilder.String()),
		UsedTools: dedupeStrings(usedTools),
	}, nil
}

func buildWorkerMessages(taskPrompt string, originalUserMessage *schema.Message) []schema.Message {
	messages := make([]schema.Message, 0, 2)
	if originalUserMessage != nil {
		messages = append(messages, cloneWorkflowContextMessage(*originalUserMessage))
	}
	messages = append(messages, schema.Message{
		Role:    schema.User,
		Content: taskPrompt,
	})
	return messages
}

func buildWorkerPrompt(plan workflowPlan, task workflowPlanTask, priorResults map[string]workflowTaskResult) string {
	dependencySummaries := make([]string, 0, len(task.Dependencies))
	for _, dep := range task.Dependencies {
		if result, ok := priorResults[dep]; ok {
			dependencySummaries = append(dependencySummaries, fmt.Sprintf("%s: %s", dep, compactText(result.Output, 500)))
		}
	}
	return fmt.Sprintf(`整体目标：%s

完成标准：
%s

当前任务：
- id: %s
- 标题: %s
- 描述: %s
- 期望输出: %s

依赖结果：
%s

请完成当前任务，并输出简洁、可靠、可供后续汇总的结果。`,
		plan.Goal,
		strings.Join(plan.CompletionCriteria, "\n"),
		task.ID,
		task.Title,
		task.Description,
		task.ExpectedOutput,
		strings.Join(dependencySummaries, "\n"),
	)
}

func synthesizeWorkflowAnswer(ctx context.Context, provider *llm_provider.Provider, userRequest string, originalUserMessage *schema.Message, plan workflowPlan, results map[string]workflowTaskResult, reviewFeedback string) (string, error) {
	synthesisMessages := []schema.Message{{Role: schema.System, Content: provider.Prompts().SynthesizerSystem}}
	if originalUserMessage != nil {
		synthesisMessages = append(synthesisMessages, cloneWorkflowContextMessage(*originalUserMessage))
	}
	synthesisMessages = append(synthesisMessages, schema.Message{
		Role:    schema.User,
		Content: buildSynthesisPrompt(provider.Prompts().SynthesisUser, userRequest, plan, results, reviewFeedback),
	})
	resp, err := provider.Generate(ctx, synthesisMessages)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(resp.Content), nil
}

func buildSynthesisPrompt(promptTemplate string, userRequest string, plan workflowPlan, results map[string]workflowTaskResult, reviewFeedback string) string {
	taskIDs := make([]string, 0, len(results))
	for taskID := range results {
		taskIDs = append(taskIDs, taskID)
	}
	sort.Strings(taskIDs)
	parts := make([]string, 0, len(taskIDs))
	for _, taskID := range taskIDs {
		result := results[taskID]
		parts = append(parts, fmt.Sprintf("[%s][%s]\n%s", result.TaskID, result.AgentName, result.Output))
	}
	feedbackBlock := "无"
	if strings.TrimSpace(reviewFeedback) != "" {
		feedbackBlock = reviewFeedback
	}
	taskResults := strings.Join(parts, "\n\n")
	if taskResults == "" {
		taskResults = "(无)"
	}
	return prompts.Render(promptTemplate, map[string]string{
		"user_request":        userRequest,
		"goal":                plan.Goal,
		"completion_criteria": strings.Join(plan.CompletionCriteria, "；"),
		"task_results":        taskResults,
		"review_feedback":     feedbackBlock,
	})
}

func reviewWorkflowAnswer(ctx context.Context, provider *llm_provider.Provider, userRequest string, originalUserMessage *schema.Message, plan workflowPlan, results map[string]workflowTaskResult, draft string) (reviewDecision, error) {
	reviewMessages := []schema.Message{{Role: schema.System, Content: provider.Prompts().ReviewerSystem}}
	if originalUserMessage != nil {
		reviewMessages = append(reviewMessages, cloneWorkflowContextMessage(*originalUserMessage))
	}
	reviewMessages = append(reviewMessages, schema.Message{
		Role:    schema.User,
		Content: buildReviewPrompt(provider.Prompts().ReviewUser, userRequest, plan, results, draft),
	})
	resp, err := provider.Generate(ctx, reviewMessages)
	if err != nil {
		return reviewDecision{}, err
	}
	var result reviewDecision
	if err := unmarshalJSONResponse(resp.Content, &result); err != nil {
		return reviewDecision{}, err
	}
	return result, nil
}

func buildReviewPrompt(promptTemplate string, userRequest string, plan workflowPlan, results map[string]workflowTaskResult, draft string) string {
	taskIDs := make([]string, 0, len(results))
	for taskID := range results {
		taskIDs = append(taskIDs, taskID)
	}
	sort.Strings(taskIDs)
	resultLines := make([]string, 0, len(taskIDs))
	for _, taskID := range taskIDs {
		result := results[taskID]
		resultLines = append(resultLines, fmt.Sprintf("%s: %s", taskID, compactText(result.Output, 240)))
	}
	taskResults := strings.Join(resultLines, "\n")
	if taskResults == "" {
		taskResults = "(无)"
	}
	return prompts.Render(promptTemplate, map[string]string{
		"user_request":        userRequest,
		"goal":                plan.Goal,
		"completion_criteria": strings.Join(plan.CompletionCriteria, "；"),
		"task_results":        taskResults,
		"draft":               draft,
	})
}

func batchTasksByDependencies(tasks []workflowPlanTask, filterIDs map[string]struct{}) [][]workflowPlanTask {
	taskMap := make(map[string]workflowPlanTask, len(tasks))
	inDegree := make(map[string]int, len(tasks))
	adj := make(map[string][]string, len(tasks))
	includeTask := func(task workflowPlanTask) bool {
		if filterIDs == nil {
			return true
		}
		_, ok := filterIDs[task.ID]
		return ok
	}

	for _, task := range tasks {
		if !includeTask(task) {
			continue
		}
		taskMap[task.ID] = task
		inDegree[task.ID] = 0
	}
	for _, task := range taskMap {
		for _, dep := range task.Dependencies {
			if _, ok := taskMap[dep]; !ok {
				continue
			}
			adj[dep] = append(adj[dep], task.ID)
			inDegree[task.ID]++
		}
	}
	var batches [][]workflowPlanTask
	var current []string
	for taskID, degree := range inDegree {
		if degree == 0 {
			current = append(current, taskID)
		}
	}
	sort.Strings(current)
	for len(current) > 0 {
		nextIDs := []string{}
		batch := make([]workflowPlanTask, 0, len(current))
		for _, taskID := range current {
			batch = append(batch, taskMap[taskID])
			for _, child := range adj[taskID] {
				inDegree[child]--
				if inDegree[child] == 0 {
					nextIDs = append(nextIDs, child)
				}
			}
		}
		sort.Strings(nextIDs)
		batches = append(batches, batch)
		current = nextIDs
	}
	if len(batches) == 0 && len(taskMap) > 0 {
		fallback := make([]workflowPlanTask, 0, len(taskMap))
		for _, task := range taskMap {
			fallback = append(fallback, task)
		}
		sort.Slice(fallback, func(i, j int) bool { return fallback[i].ID < fallback[j].ID })
		batches = append(batches, fallback)
	}
	return batches
}

func compactText(value string, maxLen int) string {
	value = strings.TrimSpace(value)
	runes := []rune(value)
	if len(runes) <= maxLen {
		return value
	}
	return string(runes[:maxLen]) + "..."
}

func dedupeStrings(values []string) []string {
	seen := map[string]struct{}{}
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

func unmarshalJSONResponse(raw string, target interface{}) error {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	raw = strings.TrimSpace(raw)
	return json.Unmarshal([]byte(raw), target)
}

func formatWorkflowPlanForTrace(plan workflowPlan) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("## 目标\n%s\n\n", plan.Goal))
	if len(plan.CompletionCriteria) > 0 {
		builder.WriteString("## 完成标准\n")
		for _, item := range plan.CompletionCriteria {
			builder.WriteString(fmt.Sprintf("- %s\n", item))
		}
		builder.WriteString("\n")
	}
	builder.WriteString("## 子任务\n")
	for _, task := range plan.Tasks {
		builder.WriteString(fmt.Sprintf("### %s %s\n", task.ID, task.Title))
		builder.WriteString(fmt.Sprintf("- 描述：%s\n", task.Description))
		builder.WriteString(fmt.Sprintf("- 依赖：%s\n", strings.Join(task.Dependencies, ", ")))
		builder.WriteString(fmt.Sprintf("- Agent：%s\n", task.SuggestedAgent))
		builder.WriteString(fmt.Sprintf("- 工具：%s\n", strings.Join(task.RequiredTools, ", ")))
		builder.WriteString(fmt.Sprintf("- 期望输出：%s\n\n", task.ExpectedOutput))
	}
	return strings.TrimSpace(builder.String())
}

func formatDispatchedTaskForTrace(task workflowPlanTask, batchNo int) string {
	return strings.TrimSpace(fmt.Sprintf("## 批次\n%d\n\n## 任务\n- ID：%s\n- 标题：%s\n- 描述：%s\n- 依赖：%s\n- 期望输出：%s",
		batchNo,
		task.ID,
		task.Title,
		task.Description,
		strings.Join(task.Dependencies, ", "),
		task.ExpectedOutput,
	))
}

func buildAgentTraceDetails(plan workflowPlan, task workflowPlanTask, priorResults map[string]workflowTaskResult, retryInstructions string) []data_models.TraceDetailBlock {
	blocks := []data_models.TraceDetailBlock{
		{
			Kind:    "input",
			Title:   "子任务输入",
			Content: buildWorkerPrompt(plan, task, priorResults),
			Format:  data_models.TraceDetailFormatMarkdown,
		},
	}
	if len(task.Dependencies) > 0 {
		var builder strings.Builder
		for _, dep := range task.Dependencies {
			if result, ok := priorResults[dep]; ok {
				builder.WriteString(fmt.Sprintf("## %s\n%s\n\n", dep, result.Output))
			}
		}
		if builder.Len() > 0 {
			blocks = append(blocks, data_models.TraceDetailBlock{
				Kind:    "dependency",
				Title:   "依赖结果",
				Content: strings.TrimSpace(builder.String()),
				Format:  data_models.TraceDetailFormatMarkdown,
			})
		}
	}
	if strings.TrimSpace(retryInstructions) != "" {
		blocks = append(blocks, data_models.TraceDetailBlock{
			Kind:    "retry",
			Title:   "修正要求",
			Content: retryInstructions,
			Format:  data_models.TraceDetailFormatText,
		})
	}
	return blocks
}

func buildAgentResultTraceDetails(plan workflowPlan, task workflowPlanTask, priorResults map[string]workflowTaskResult, retryInstructions string, result workflowTaskResult, execErr error) []data_models.TraceDetailBlock {
	blocks := buildAgentTraceDetails(plan, task, priorResults, retryInstructions)
	if execErr != nil {
		blocks = append(blocks, data_models.TraceDetailBlock{
			Kind:    "output",
			Title:   "执行错误",
			Content: execErr.Error(),
			Format:  data_models.TraceDetailFormatText,
		})
		return blocks
	}
	blocks = append(blocks, data_models.TraceDetailBlock{
		Kind:    "output",
		Title:   "子任务输出",
		Content: result.Output,
		Format:  data_models.TraceDetailFormatMarkdown,
	})
	if len(result.UsedTools) > 0 {
		blocks = append(blocks, data_models.TraceDetailBlock{
			Kind:    "tool_result",
			Title:   "使用工具",
			Content: strings.Join(result.UsedTools, ", "),
			Format:  data_models.TraceDetailFormatText,
		})
	}
	return blocks
}

func formatReviewDecisionForTrace(review reviewDecision) string {
	payload, err := json.MarshalIndent(review, "", "  ")
	if err != nil {
		return fmt.Sprintf(`{"approved":%t}`, review.Approved)
	}
	return string(payload)
}

func formatRetrySummaryForTrace(review reviewDecision) string {
	return strings.TrimSpace(fmt.Sprintf("## 问题\n%s\n\n## 修正指令\n%s\n\n## 受影响任务\n%s",
		strings.Join(review.Issues, "\n"),
		review.RetryInstructions,
		strings.Join(review.AffectedTaskIDs, ", "),
	))
}
