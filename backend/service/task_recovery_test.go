package service

import (
	"context"
	"testing"
	"time"

	"github.com/cloudwego/eino/schema"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/tasker"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/storage"
)

func newTaskRecoveryTestService(t *testing.T) (*Service, *storage.Storage) {
	t.Helper()
	t.Setenv("LEMONTEA_DATA_PATH", t.TempDir())

	st, err := storage.NewStorage()
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}
	svc := NewService()
	svc.storage = st
	return svc, st
}

func seedTaskRecoveryFixture(t *testing.T, st *storage.Storage, task data_models.Task, message data_models.Message) {
	t.Helper()
	ctx := context.Background()

	if err := st.CreateChat(ctx, task.ChatUuid, "测试对话"); err != nil {
		t.Fatalf("CreateChat() error = %v", err)
	}
	message.ChatUuid = task.ChatUuid
	if _, err := st.CreateMessage(ctx, task.ChatUuid, message); err != nil {
		t.Fatalf("CreateMessage() error = %v", err)
	}
	if err := st.CreateTask(ctx, task); err != nil {
		t.Fatalf("CreateTask() error = %v", err)
	}
}

func TestRecoverStaleRunningTasksMarksTaskAndMessageFailed(t *testing.T) {
	svc, st := newTaskRecoveryTestService(t)

	startedAt := time.Now().Add(-2 * time.Second)
	task := data_models.Task{
		TaskUuid:             "task-stale-1",
		ChatUuid:             "chat-stale-1",
		AssistantMessageUuid: "assistant-stale-1",
		Status:               data_models.TaskStatusRunning,
		EventKey:             "event:task:stale-1",
		StartedAt:            &startedAt,
	}
	message := data_models.Message{
		MessageUuid: "assistant-stale-1",
		Role:        schema.Assistant,
		AssistantMessageExtra: &data_models.AssistantMessageExtra{
			CurrentStage: "准备执行",
			CurrentAgent: "planner",
			ToolUses: []data_models.ToolUse{
				{CallID: "tool-1", Status: data_models.ToolUseStatusRunning, StartedAt: &startedAt},
			},
			ExecutionTrace: data_models.ExecutionTrace{
				Steps: []data_models.TraceStep{
					{StepID: "step-1", Status: data_models.TraceStepStatusPending, StartedAt: &startedAt},
				},
			},
		},
	}
	seedTaskRecoveryFixture(t, st, task, message)

	if err := svc.recoverStaleRunningTasks(context.Background()); err != nil {
		t.Fatalf("recoverStaleRunningTasks() error = %v", err)
	}

	gotTask, err := st.GetTask(context.Background(), task.TaskUuid)
	if err != nil {
		t.Fatalf("GetTask() error = %v", err)
	}
	if gotTask == nil {
		t.Fatal("GetTask() = nil, want task")
	}
	if gotTask.Status != data_models.TaskStatusFailed {
		t.Fatalf("task status = %q, want %q", gotTask.Status, data_models.TaskStatusFailed)
	}
	if gotTask.FinishReason != "error" {
		t.Fatalf("task finish_reason = %q, want error", gotTask.FinishReason)
	}
	if gotTask.FinishError != interruptedTaskFinishError {
		t.Fatalf("task finish_error = %q, want %q", gotTask.FinishError, interruptedTaskFinishError)
	}
	if gotTask.FinishedAt == nil {
		t.Fatal("task finished_at = nil, want non-nil")
	}

	gotMessage, err := st.GetMessageByUUID(context.Background(), task.AssistantMessageUuid)
	if err != nil {
		t.Fatalf("GetMessageByUUID() error = %v", err)
	}
	if gotMessage == nil || gotMessage.AssistantMessageExtra == nil {
		t.Fatalf("assistant message = %#v, want assistant extra", gotMessage)
	}
	if gotMessage.AssistantMessageExtra.FinishReason != "error" {
		t.Fatalf("message finish_reason = %q, want error", gotMessage.AssistantMessageExtra.FinishReason)
	}
	if gotMessage.AssistantMessageExtra.FinishError != interruptedTaskFinishError {
		t.Fatalf("message finish_error = %q, want %q", gotMessage.AssistantMessageExtra.FinishError, interruptedTaskFinishError)
	}
	if gotMessage.AssistantMessageExtra.CurrentStage != "" || gotMessage.AssistantMessageExtra.CurrentAgent != "" {
		t.Fatalf("current stage/agent = %q/%q, want empty", gotMessage.AssistantMessageExtra.CurrentStage, gotMessage.AssistantMessageExtra.CurrentAgent)
	}
	if gotMessage.AssistantMessageExtra.ToolUses[0].Status != data_models.ToolUseStatusError {
		t.Fatalf("tool status = %q, want %q", gotMessage.AssistantMessageExtra.ToolUses[0].Status, data_models.ToolUseStatusError)
	}
	if gotMessage.AssistantMessageExtra.ToolUses[0].FinishedAt == nil {
		t.Fatal("tool finished_at = nil, want non-nil")
	}
	if gotMessage.AssistantMessageExtra.ExecutionTrace.Steps[0].Status != data_models.TraceStepStatusError {
		t.Fatalf("trace step status = %q, want %q", gotMessage.AssistantMessageExtra.ExecutionTrace.Steps[0].Status, data_models.TraceStepStatusError)
	}
	if gotMessage.AssistantMessageExtra.ExecutionTrace.Steps[0].FinishedAt == nil {
		t.Fatal("trace step finished_at = nil, want non-nil")
	}
}

func TestGetRunningTasksFiltersStaleTasks(t *testing.T) {
	svc, st := newTaskRecoveryTestService(t)

	liveTask := data_models.Task{
		TaskUuid:             "task-live-1",
		ChatUuid:             "chat-live-1",
		AssistantMessageUuid: "assistant-live-1",
		Status:               data_models.TaskStatusRunning,
		EventKey:             "event:task:live-1",
	}
	staleTask := data_models.Task{
		TaskUuid:             "task-stale-2",
		ChatUuid:             "chat-stale-2",
		AssistantMessageUuid: "assistant-stale-2",
		Status:               data_models.TaskStatusPending,
		EventKey:             "event:task:stale-2",
	}
	seedTaskRecoveryFixture(t, st, liveTask, data_models.Message{
		MessageUuid: "assistant-live-1",
		Role:        schema.Assistant,
		AssistantMessageExtra: &data_models.AssistantMessageExtra{
			CurrentStage: "执行中",
		},
	})
	seedTaskRecoveryFixture(t, st, staleTask, data_models.Message{
		MessageUuid: "assistant-stale-2",
		Role:        schema.Assistant,
		AssistantMessageExtra: &data_models.AssistantMessageExtra{
			CurrentStage: "等待执行",
		},
	})

	stopped := make(chan struct{})
	tasker.Manager.StartTask(tasker.Runtime{TaskUUID: liveTask.TaskUuid, EventKey: liveTask.EventKey}, func(stop <-chan struct{}) {
		<-stop
		close(stopped)
	})
	defer func() {
		tasker.Manager.StopTask(liveTask.TaskUuid)
		<-stopped
	}()

	got, err := svc.GetRunningTasks(context.Background())
	if err != nil {
		t.Fatalf("GetRunningTasks() error = %v", err)
	}
	if got == nil {
		t.Fatal("GetRunningTasks() = nil, want task list")
	}
	if len(got.Tasks) != 1 {
		t.Fatalf("GetRunningTasks().Tasks len = %d, want 1", len(got.Tasks))
	}
	if got.Tasks[0].TaskUuid != liveTask.TaskUuid {
		t.Fatalf("GetRunningTasks().Tasks[0] = %q, want %q", got.Tasks[0].TaskUuid, liveTask.TaskUuid)
	}

	staleAfter, err := st.GetTask(context.Background(), staleTask.TaskUuid)
	if err != nil {
		t.Fatalf("GetTask(stale) error = %v", err)
	}
	if staleAfter == nil || staleAfter.Status != data_models.TaskStatusFailed {
		t.Fatalf("stale task status = %#v, want failed", staleAfter)
	}
}

func TestGetChatActiveTaskSkipsStaleAndReturnsLiveTask(t *testing.T) {
	svc, st := newTaskRecoveryTestService(t)

	liveTask := data_models.Task{
		OrmModel: data_models.OrmModel{
			CreatedAt: time.Now().Add(-2 * time.Minute),
			UpdatedAt: time.Now().Add(-2 * time.Minute),
		},
		TaskUuid:             "task-live-2",
		ChatUuid:             "chat-shared-1",
		AssistantMessageUuid: "assistant-live-2",
		Status:               data_models.TaskStatusRunning,
		EventKey:             "event:task:live-2",
	}
	staleTask := data_models.Task{
		OrmModel: data_models.OrmModel{
			CreatedAt: time.Now().Add(-1 * time.Minute),
			UpdatedAt: time.Now().Add(-1 * time.Minute),
		},
		TaskUuid:             "task-stale-3",
		ChatUuid:             "chat-shared-1",
		AssistantMessageUuid: "assistant-stale-3",
		Status:               data_models.TaskStatusRunning,
		EventKey:             "event:task:stale-3",
	}
	seedTaskRecoveryFixture(t, st, liveTask, data_models.Message{
		MessageUuid: "assistant-live-2",
		Role:        schema.Assistant,
		AssistantMessageExtra: &data_models.AssistantMessageExtra{
			CurrentStage: "执行中",
		},
	})
	if err := st.CreateTask(context.Background(), staleTask); err != nil {
		t.Fatalf("CreateTask(stale) error = %v", err)
	}
	if _, err := st.CreateMessage(context.Background(), staleTask.ChatUuid, data_models.Message{
		ChatUuid:              staleTask.ChatUuid,
		MessageUuid:           "assistant-stale-3",
		Role:                  schema.Assistant,
		AssistantMessageExtra: &data_models.AssistantMessageExtra{CurrentStage: "执行中"},
	}); err != nil {
		t.Fatalf("CreateMessage(stale) error = %v", err)
	}

	stopped := make(chan struct{})
	tasker.Manager.StartTask(tasker.Runtime{TaskUUID: liveTask.TaskUuid, EventKey: liveTask.EventKey}, func(stop <-chan struct{}) {
		<-stop
		close(stopped)
	})
	defer func() {
		tasker.Manager.StopTask(liveTask.TaskUuid)
		<-stopped
	}()

	got, err := svc.GetChatActiveTask(context.Background(), liveTask.ChatUuid)
	if err != nil {
		t.Fatalf("GetChatActiveTask() error = %v", err)
	}
	if got == nil {
		t.Fatal("GetChatActiveTask() = nil, want live task")
	}
	if got.TaskUuid != liveTask.TaskUuid {
		t.Fatalf("GetChatActiveTask().TaskUuid = %q, want %q", got.TaskUuid, liveTask.TaskUuid)
	}

	staleAfter, err := st.GetTask(context.Background(), staleTask.TaskUuid)
	if err != nil {
		t.Fatalf("GetTask(stale) error = %v", err)
	}
	if staleAfter == nil || staleAfter.Status != data_models.TaskStatusFailed {
		t.Fatalf("stale task status = %#v, want failed", staleAfter)
	}
}

func TestRecoverWaitingApprovalTaskExpiresPendingApproval(t *testing.T) {
	svc, st := newTaskRecoveryTestService(t)

	startedAt := time.Now().Add(-2 * time.Second)
	task := data_models.Task{
		TaskUuid:             "task-waiting-approval-1",
		ChatUuid:             "chat-waiting-approval-1",
		AssistantMessageUuid: "assistant-waiting-approval-1",
		Status:               data_models.TaskStatusWaitingApproval,
		EventKey:             "event:task:waiting-approval-1",
		StartedAt:            &startedAt,
	}
	message := data_models.Message{
		MessageUuid: "assistant-waiting-approval-1",
		Role:        schema.Assistant,
		AssistantMessageExtra: &data_models.AssistantMessageExtra{
			CurrentStage: "等待用户确认",
			ToolUses: []data_models.ToolUse{
				{CallID: "tool-approval-1", Status: data_models.ToolUseStatusAwaitingApproval, StartedAt: &startedAt},
			},
			PendingApprovals: []data_models.ToolApprovalSummary{
				{ApprovalID: "approval-pending-1", ToolCallID: "tool-approval-1", ToolName: "文件工具", Status: data_models.ToolApprovalStatusPending},
			},
			ExecutionTrace: data_models.ExecutionTrace{
				Steps: []data_models.TraceStep{
					{StepID: "tool-approval-1", Status: data_models.TraceStepStatusAwaitingApproval, StartedAt: &startedAt},
				},
			},
		},
	}
	seedTaskRecoveryFixture(t, st, task, message)
	if err := st.CreateToolApproval(context.Background(), data_models.ToolApproval{
		ApprovalID:           "approval-pending-1",
		TaskUuid:             task.TaskUuid,
		ChatUuid:             task.ChatUuid,
		AssistantMessageUuid: task.AssistantMessageUuid,
		ToolCallID:           "tool-approval-1",
		ToolID:               "file_tool",
		ToolName:             "文件工具",
		Status:               data_models.ToolApprovalStatusPending,
		Title:                "读取文件",
		Message:              "等待确认",
		RequestedAt:          &startedAt,
	}); err != nil {
		t.Fatalf("CreateToolApproval() error = %v", err)
	}

	if err := svc.recoverStaleRunningTasks(context.Background()); err != nil {
		t.Fatalf("recoverStaleRunningTasks() error = %v", err)
	}

	gotTask, err := st.GetTask(context.Background(), task.TaskUuid)
	if err != nil {
		t.Fatalf("GetTask() error = %v", err)
	}
	if gotTask == nil {
		t.Fatal("GetTask() = nil, want task")
	}
	if gotTask.FinishError != expiredApprovalFinishError {
		t.Fatalf("finish_error = %q, want %q", gotTask.FinishError, expiredApprovalFinishError)
	}

	gotMessage, err := st.GetMessageByUUID(context.Background(), task.AssistantMessageUuid)
	if err != nil {
		t.Fatalf("GetMessageByUUID() error = %v", err)
	}
	if gotMessage == nil || gotMessage.AssistantMessageExtra == nil {
		t.Fatalf("assistant message = %#v, want assistant extra", gotMessage)
	}
	if len(gotMessage.AssistantMessageExtra.PendingApprovals) != 0 {
		t.Fatalf("pending approvals len = %d, want 0", len(gotMessage.AssistantMessageExtra.PendingApprovals))
	}
	if gotMessage.AssistantMessageExtra.ToolUses[0].Status != data_models.ToolUseStatusError {
		t.Fatalf("tool status = %q, want error", gotMessage.AssistantMessageExtra.ToolUses[0].Status)
	}
	if gotMessage.AssistantMessageExtra.ExecutionTrace.Steps[0].Status != data_models.TraceStepStatusError {
		t.Fatalf("trace status = %q, want error", gotMessage.AssistantMessageExtra.ExecutionTrace.Steps[0].Status)
	}

	storedApproval, err := st.GetToolApprovalByApprovalID(context.Background(), "approval-pending-1")
	if err != nil {
		t.Fatalf("GetToolApprovalByApprovalID() error = %v", err)
	}
	if storedApproval == nil {
		t.Fatal("stored approval = nil")
	}
	if storedApproval.Status != data_models.ToolApprovalStatusExpired {
		t.Fatalf("approval status = %q, want expired", storedApproval.Status)
	}
	if storedApproval.ResponseComment != expiredApprovalFinishError {
		t.Fatalf("response comment = %q, want %q", storedApproval.ResponseComment, expiredApprovalFinishError)
	}
}
