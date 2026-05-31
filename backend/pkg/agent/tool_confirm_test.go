package agent

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/agent/tools"
	"trpc.group/trpc-go/trpc-agent-go/model"
	toolpkg "trpc.group/trpc-go/trpc-agent-go/tool"
)

func TestToolConfirmCallbackWaitsForApprovalBeforeContinuing(t *testing.T) {
	manager := NewManager(newTestStorage(t))

	session, err := manager.Storage().CreateSession(data_models.Session{
		Title:  "Confirm Callback",
		Status: "loading",
	})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	streamCtx, cancel := manager.Streams().Start(context.Background(), session.ID)
	defer cancel()

	ctx := withRunContext(streamCtx, runContext{
		SessionID: session.ID,
		ModelName: "gpt-test",
	})
	callbacks := manager.newToolCallbacks()

	args, err := json.Marshal(map[string]string{
		"command": "mv ~/Downloads/a.pdf ~/Documents/a.pdf",
	})
	if err != nil {
		t.Fatalf("marshal args: %v", err)
	}

	done := make(chan struct{})
	var customResult any
	var callbackErr error
	go func() {
		result, err := callbacks.RunBeforeTool(ctx, &toolpkg.BeforeToolArgs{
			ToolName:  "shell",
			Arguments: args,
		})
		callbackErr = err
		if result != nil {
			customResult = result.CustomResult
		}
		close(done)
	}()

	select {
	case <-done:
		t.Fatal("before-tool callback returned before confirm response")
	case <-time.After(100 * time.Millisecond):
	}

	updatedSession, err := manager.Storage().GetSession(session.ID)
	if err != nil {
		t.Fatalf("get session: %v", err)
	}
	if updatedSession.Status != "waiting-unread" {
		t.Fatalf("expected waiting-unread before approval, got %q", updatedSession.Status)
	}

	if ok := manager.Streams().SendConfirmResponse(session.ID, ConfirmResponse{
		Approved: true,
		Action:   "approve",
	}); !ok {
		t.Fatal("expected confirm response to be accepted")
	}

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("before-tool callback did not return after approval")
	}

	if callbackErr != nil {
		t.Fatalf("unexpected callback error: %v", callbackErr)
	}
	if customResult != nil {
		t.Fatalf("expected nil custom result for approval, got %#v", customResult)
	}

	messages, err := manager.Storage().ListMessagesForSession(session.ID, 0, 10)
	if err != nil {
		t.Fatalf("list messages: %v", err)
	}
	if len(messages) != 1 {
		t.Fatalf("expected 1 confirm response message, got %d", len(messages))
	}
	if messages[0].ContentType != "confirm_response" || messages[0].Content != "approved" {
		t.Fatalf("unexpected confirm message: %+v", messages[0])
	}
}

func TestToolConfirmCallbackSkipsSafeShellCommand(t *testing.T) {
	manager := NewManager(newTestStorage(t))

	session, err := manager.Storage().CreateSession(data_models.Session{
		Title:  "Safe Shell",
		Status: "loading",
	})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	streamCtx, cancel := manager.Streams().Start(context.Background(), session.ID)
	defer cancel()

	ctx := withRunContext(streamCtx, runContext{
		SessionID: session.ID,
		ModelName: "gpt-test",
	})
	callbacks := manager.newToolCallbacks()

	result, err := callbacks.RunBeforeTool(ctx, &toolpkg.BeforeToolArgs{
		ToolName:  "shell",
		Arguments: []byte(`{"command":"pwd"}`),
	})
	if err != nil {
		t.Fatalf("unexpected callback error: %v", err)
	}
	if result != nil {
		t.Fatalf("expected safe shell command to continue without confirmation, got %#v", result)
	}

	updatedSession, err := manager.Storage().GetSession(session.ID)
	if err != nil {
		t.Fatalf("get session: %v", err)
	}
	if updatedSession.Status != "loading" {
		t.Fatalf("expected session to remain loading, got %q", updatedSession.Status)
	}
}

func TestToolConfirmCallbackDeniesDestructiveShellCommand(t *testing.T) {
	manager := NewManager(newTestStorage(t))

	session, err := manager.Storage().CreateSession(data_models.Session{
		Title:  "Deny Shell",
		Status: "loading",
	})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	streamCtx, cancel := manager.Streams().Start(context.Background(), session.ID)
	defer cancel()

	ctx := withRunContext(streamCtx, runContext{
		SessionID: session.ID,
		ModelName: "gpt-test",
	})
	callbacks := manager.newToolCallbacks()

	result, err := callbacks.RunBeforeTool(ctx, &toolpkg.BeforeToolArgs{
		ToolName:  "shell",
		Arguments: []byte(`{"command":"rm -rf ~/Documents/reports"}`),
	})
	if err != nil {
		t.Fatalf("unexpected callback error: %v", err)
	}
	if result == nil {
		t.Fatal("expected destructive shell command to be denied")
	}
	resultMap, ok := result.CustomResult.(map[string]any)
	if !ok {
		t.Fatalf("expected synthetic deny result, got %#v", result.CustomResult)
	}
	if resultMap["error"] != "tool execution denied by policy: destructive shell command" {
		t.Fatalf("unexpected deny result: %#v", resultMap)
	}

	messages, err := manager.Storage().ListMessagesForSession(session.ID, 0, 10)
	if err != nil {
		t.Fatalf("list messages: %v", err)
	}
	if len(messages) != 0 {
		t.Fatalf("expected no user confirmation messages for policy denial, got %d", len(messages))
	}
}

func TestToolConfirmCallbackDeniesChainedDestructiveShellCommand(t *testing.T) {
	manager := NewManager(newTestStorage(t))

	session, err := manager.Storage().CreateSession(data_models.Session{
		Title:  "Deny Chained Shell",
		Status: "loading",
	})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	streamCtx, cancel := manager.Streams().Start(context.Background(), session.ID)
	defer cancel()

	ctx := withRunContext(streamCtx, runContext{
		SessionID: session.ID,
		ModelName: "gpt-test",
	})
	callbacks := manager.newToolCallbacks()

	result, err := callbacks.RunBeforeTool(ctx, &toolpkg.BeforeToolArgs{
		ToolName:  "shell",
		Arguments: []byte(`{"command":"ls ~/Downloads ; rm -rf ~/Documents/reports"}`),
	})
	if err != nil {
		t.Fatalf("unexpected callback error: %v", err)
	}
	if result == nil {
		t.Fatal("expected chained destructive shell command to be denied")
	}
	resultMap, ok := result.CustomResult.(map[string]any)
	if !ok {
		t.Fatalf("expected synthetic deny result, got %#v", result.CustomResult)
	}
	if resultMap["error"] != "tool execution denied by policy: destructive shell command" {
		t.Fatalf("unexpected deny result: %#v", resultMap)
	}
}

func TestToolConfirmCallbackUsesDynamicDeciderForNonStaticTools(t *testing.T) {
	manager := NewManager(newTestStorage(t))
	manager.ToolRegistry().Register(tools.ToolMeta{
		Name:        "dynamic_tool",
		Description: "test tool",
		Category:    tools.CategoryBuiltin,
	})
	manager.SetToolConfirmDecider(func(ctx context.Context, input ToolConfirmDecisionInput) (ToolConfirmDecision, error) {
		if input.ToolName != "dynamic_tool" {
			t.Fatalf("unexpected tool name: %q", input.ToolName)
		}
		if input.Arguments != `{"target":"prod"}` {
			t.Fatalf("unexpected args: %q", input.Arguments)
		}
		return ToolConfirmDecision{
			RequireConfirm: true,
			Reason:         "dynamic risk",
			Source:         "model",
		}, nil
	})

	session, err := manager.Storage().CreateSession(data_models.Session{
		Title:  "Dynamic Confirm",
		Status: "loading",
	})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	streamCtx, cancel := manager.Streams().Start(context.Background(), session.ID)
	defer cancel()

	ctx := withRunContext(streamCtx, runContext{
		SessionID:    session.ID,
		ModelName:    "gpt-test",
		BaseURL:      "https://api.example.com",
		APIKey:       "test-key",
		ProviderType: "openai_compatibility",
	})
	callbacks := manager.newToolCallbacks()

	done := make(chan struct{})
	go func() {
		_, _ = callbacks.RunBeforeTool(ctx, &toolpkg.BeforeToolArgs{
			ToolName:  "dynamic_tool",
			Arguments: []byte(`{"target":"prod"}`),
		})
		close(done)
	}()

	select {
	case <-done:
		t.Fatal("before-tool callback returned before dynamic confirmation")
	case <-time.After(100 * time.Millisecond):
	}

	updatedSession, err := manager.Storage().GetSession(session.ID)
	if err != nil {
		t.Fatalf("get session: %v", err)
	}
	if updatedSession.Status != "waiting-unread" {
		t.Fatalf("expected waiting-unread before approval, got %q", updatedSession.Status)
	}

	if ok := manager.Streams().SendConfirmResponse(session.ID, ConfirmResponse{
		Approved: true,
		Action:   "approve",
	}); !ok {
		t.Fatal("expected confirm response to be accepted")
	}

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("before-tool callback did not return after dynamic approval")
	}
}

func TestToolConfirmCallbackSkipsWhenDynamicDeciderAllowsExecution(t *testing.T) {
	manager := NewManager(newTestStorage(t))
	manager.ToolRegistry().Register(tools.ToolMeta{
		Name:        "safe_tool",
		Description: "test tool",
		Category:    tools.CategoryBuiltin,
	})
	manager.SetToolConfirmDecider(func(ctx context.Context, input ToolConfirmDecisionInput) (ToolConfirmDecision, error) {
		return ToolConfirmDecision{
			RequireConfirm: false,
			Reason:         "read only",
			Source:         "model",
		}, nil
	})

	session, err := manager.Storage().CreateSession(data_models.Session{
		Title:  "Dynamic Skip",
		Status: "loading",
	})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	streamCtx, cancel := manager.Streams().Start(context.Background(), session.ID)
	defer cancel()

	ctx := withRunContext(streamCtx, runContext{
		SessionID:    session.ID,
		ModelName:    "gpt-test",
		BaseURL:      "https://api.example.com",
		APIKey:       "test-key",
		ProviderType: "openai_compatibility",
	})
	callbacks := manager.newToolCallbacks()

	result, err := callbacks.RunBeforeTool(ctx, &toolpkg.BeforeToolArgs{
		ToolName:  "safe_tool",
		Arguments: []byte(`{"query":"status"}`),
	})
	if err != nil {
		t.Fatalf("unexpected callback error: %v", err)
	}
	if result != nil {
		t.Fatalf("expected nil result when dynamic decider skips confirmation, got %#v", result)
	}

	updatedSession, err := manager.Storage().GetSession(session.ID)
	if err != nil {
		t.Fatalf("get session: %v", err)
	}
	if updatedSession.Status != "loading" {
		t.Fatalf("expected session to remain loading, got %q", updatedSession.Status)
	}
}

func TestToolConfirmCallbackReturnsSyntheticResultOnReject(t *testing.T) {
	manager := NewManager(newTestStorage(t))

	session, err := manager.Storage().CreateSession(data_models.Session{
		Title:  "Reject Callback",
		Status: "loading",
	})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	streamCtx, cancel := manager.Streams().Start(context.Background(), session.ID)
	defer cancel()

	ctx := withRunContext(streamCtx, runContext{
		SessionID: session.ID,
		ModelName: "gpt-test",
	})
	callbacks := manager.newToolCallbacks()

	args, err := json.Marshal(map[string]string{
		"command": "mv ~/Downloads/a.pdf ~/Documents/a.pdf",
	})
	if err != nil {
		t.Fatalf("marshal args: %v", err)
	}

	done := make(chan struct{})
	var customResult any
	var callbackErr error
	go func() {
		result, err := callbacks.RunBeforeTool(ctx, &toolpkg.BeforeToolArgs{
			ToolName:  "shell",
			Arguments: args,
		})
		callbackErr = err
		if result != nil {
			customResult = result.CustomResult
		}
		close(done)
	}()

	if ok := manager.Streams().SendConfirmResponse(session.ID, ConfirmResponse{
		Approved: false,
		Message:  "not now",
		Action:   "reject",
	}); !ok {
		t.Fatal("expected reject response to be accepted")
	}

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("before-tool callback did not return after rejection")
	}

	if callbackErr != nil {
		t.Fatalf("unexpected callback error: %v", callbackErr)
	}

	resultMap, ok := customResult.(map[string]any)
	if !ok {
		t.Fatalf("expected synthetic map result on reject, got %#v", customResult)
	}
	if resultMap["error"] != "tool execution rejected by user: not now" {
		t.Fatalf("unexpected reject result: %#v", resultMap)
	}

	messages, err := manager.Storage().ListMessagesForSession(session.ID, 0, 10)
	if err != nil {
		t.Fatalf("list messages: %v", err)
	}
	if len(messages) != 1 {
		t.Fatalf("expected 1 confirm response message, got %d", len(messages))
	}
	if messages[0].ContentType != "confirm_response" || messages[0].Content != "rejected\ncomment: not now" {
		t.Fatalf("unexpected confirm message: %+v", messages[0])
	}
}

func TestToolConfirmCallbackTurnsApprovedCommentIntoRevisionInstruction(t *testing.T) {
	manager := NewManager(newTestStorage(t))
	callbacks := manager.newToolCallbacks()

	msgs, err := callbacks.RunToolResultMessages(
		withToolConfirmation(context.Background(), toolConfirmation{
			Approved: false,
			Message:  "only read the first 20 lines",
			Action:   "comment",
		}),
		&toolpkg.ToolResultMessagesInput{
			ToolName: "file_read",
			DefaultToolMessage: model.Message{
				Role:     model.RoleTool,
				ToolID:   "call_file_read_1",
				ToolName: "file_read",
				Content:  "{\"content\":\"hello\"}",
			},
		},
	)
	if err != nil {
		t.Fatalf("run tool result messages: %v", err)
	}

	msg, ok := msgs.(model.Message)
	if !ok {
		t.Fatalf("expected model.Message, got %#v", msgs)
	}
	if msg.ToolID != "call_file_read_1" {
		t.Fatalf("expected tool id to be preserved, got %#v", msg)
	}
	if msg.Content != "Tool execution paused. The user submitted this comment: only read the first 20 lines. Do not execute the previous tool call unchanged. Revise the tool call to incorporate the comment, then continue." {
		t.Fatalf("unexpected tool content: %q", msg.Content)
	}
}

func TestToolConfirmCallbackSkipsExecutionWhenApprovedWithComment(t *testing.T) {
	manager := NewManager(newTestStorage(t))

	session, err := manager.Storage().CreateSession(data_models.Session{
		Title:  "Approved With Comment",
		Status: "loading",
	})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}

	streamCtx, cancel := manager.Streams().Start(context.Background(), session.ID)
	defer cancel()

	ctx := withRunContext(streamCtx, runContext{
		SessionID: session.ID,
		ModelName: "gpt-test",
	})
	callbacks := manager.newToolCallbacks()

	args := []byte(`{"command":"mv ~/Downloads/a.pdf ~/Documents/a.pdf"}`)

	done := make(chan struct{})
	var customResult any
	var callbackErr error
	var resultCtx context.Context
	go func() {
		result, err := callbacks.RunBeforeTool(ctx, &toolpkg.BeforeToolArgs{
			ToolName:  "shell",
			Arguments: args,
		})
		callbackErr = err
		if result != nil {
			customResult = result.CustomResult
			resultCtx = result.Context
		}
		close(done)
	}()

	if ok := manager.Streams().SendConfirmResponse(session.ID, ConfirmResponse{
		Approved: false,
		Message:  "只列出 markdown 文件",
		Action:   "comment",
	}); !ok {
		t.Fatal("expected confirm response to be accepted")
	}

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("before-tool callback did not return after approval with comment")
	}

	if callbackErr != nil {
		t.Fatalf("unexpected callback error: %v", callbackErr)
	}
	if customResult == nil {
		t.Fatal("expected custom result when approval includes comment")
	}

	resultMap, ok := customResult.(map[string]any)
	if !ok {
		t.Fatalf("expected synthetic map result, got %#v", customResult)
	}
	if resultMap["status"] != "user_comment" {
		t.Fatalf("unexpected custom result: %#v", resultMap)
	}

	messages, err := manager.Storage().ListMessagesForSession(session.ID, 0, 10)
	if err != nil {
		t.Fatalf("list messages: %v", err)
	}
	if len(messages) != 1 {
		t.Fatalf("expected 1 confirm response message, got %d", len(messages))
	}
	if messages[0].ContentType != "confirm_response" || messages[0].Content != "commented\ncomment: 只列出 markdown 文件" {
		t.Fatalf("unexpected confirm message: %+v", messages[0])
	}

	msgs, err := callbacks.RunToolResultMessages(resultCtx, &toolpkg.ToolResultMessagesInput{
		ToolName: "shell",
		DefaultToolMessage: model.Message{
			Role:     model.RoleTool,
			ToolID:   "call_shell_1",
			ToolName: "shell",
			Content:  `{"status":"user_comment"}`,
		},
	})
	if err != nil {
		t.Fatalf("run tool result messages: %v", err)
	}

	msg, ok := msgs.(model.Message)
	if !ok {
		t.Fatalf("expected model.Message, got %#v", msgs)
	}
	if msg.Content != "Tool execution paused. The user submitted this comment: 只列出 markdown 文件. Do not execute the previous tool call unchanged. Revise the tool call to incorporate the comment, then continue." {
		t.Fatalf("unexpected tool content: %q", msg.Content)
	}
}
