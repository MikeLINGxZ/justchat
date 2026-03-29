package view_models

import "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"

type Task = data_models.Task

type TaskList struct {
	Tasks []Task `json:"tasks"`
}

type TaskStreamEvent struct {
	TaskUuid         string                     `json:"task_uuid"`
	ChatUuid         string                     `json:"chat_uuid"`
	EventKey         string                     `json:"event_key"`
	Status           data_models.TaskStatus     `json:"status"`
	FinishReason     string                     `json:"finish_reason"`
	FinishError      string                     `json:"finish_error"`
	ExecutionTrace   data_models.ExecutionTrace `json:"execution_trace"`
	TraceDelta       []data_models.TraceStep    `json:"trace_delta"`
	CurrentStage     string                     `json:"current_stage"`
	CurrentAgent     string                     `json:"current_agent"`
	RetryCount       int                        `json:"retry_count"`
	AssistantMessage Message                    `json:"assistant_message"`
}
