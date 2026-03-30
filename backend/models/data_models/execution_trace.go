package data_models

import "time"

type RouteType string

const (
	RouteTypeDirectAnswer RouteType = "direct_answer"
	RouteTypeWorkflow     RouteType = "workflow"
	RouteTypeClarify      RouteType = "clarify"
)

type TraceStepType string

const (
	TraceStepTypeClassify   TraceStepType = "classify"
	TraceStepTypePlan       TraceStepType = "plan"
	TraceStepTypeDispatch   TraceStepType = "dispatch"
	TraceStepTypeAgentRun   TraceStepType = "agent_run"
	TraceStepTypeToolCall   TraceStepType = "tool_call"
	TraceStepTypeSynthesize TraceStepType = "synthesize"
	TraceStepTypeReview     TraceStepType = "review"
	TraceStepTypeRetry      TraceStepType = "retry"
	TraceStepTypeFinalize   TraceStepType = "finalize"
)

type TraceStepStatus string

const (
	TraceStepStatusPending          TraceStepStatus = "pending"
	TraceStepStatusRunning          TraceStepStatus = "running"
	TraceStepStatusAwaitingApproval TraceStepStatus = "awaiting_approval"
	TraceStepStatusDone             TraceStepStatus = "done"
	TraceStepStatusRejected         TraceStepStatus = "rejected"
	TraceStepStatusError            TraceStepStatus = "error"
	TraceStepStatusSkipped          TraceStepStatus = "skipped"
)

type ExecutionTrace struct {
	Steps []TraceStep `json:"steps"`
}

type TraceDetailFormat string

const (
	TraceDetailFormatText     TraceDetailFormat = "text"
	TraceDetailFormatMarkdown TraceDetailFormat = "markdown"
	TraceDetailFormatJSON     TraceDetailFormat = "json"
)

type TraceDetailBlock struct {
	Kind    string            `json:"kind"`
	Title   string            `json:"title"`
	Content string            `json:"content"`
	Format  TraceDetailFormat `json:"format"`
}

type TraceStep struct {
	StepID        string                 `json:"step_id"`
	ParentStepID  string                 `json:"parent_step_id"`
	Type          TraceStepType          `json:"type"`
	Title         string                 `json:"title"`
	Summary       string                 `json:"summary"`
	Status        TraceStepStatus        `json:"status"`
	AgentName     string                 `json:"agent_name"`
	ToolName      string                 `json:"tool_name"`
	InputPreview  string                 `json:"input_preview"`
	OutputPreview string                 `json:"output_preview"`
	StartedAt     *time.Time             `json:"started_at"`
	FinishedAt    *time.Time             `json:"finished_at"`
	ElapsedMs     int64                  `json:"elapsed_ms"`
	DetailBlocks  []TraceDetailBlock     `json:"detail_blocks"`
	Metadata      map[string]interface{} `json:"metadata"`
}
