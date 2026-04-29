package service

import (
	"context"
	"time"

	"github.com/cloudwego/eino/compose"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/llm_provider/tools"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/tool_approval"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/storage"
)

type ToolHandler func(ctx context.Context, callID, toolName, toolArgs string) (string, error)

type ToolMiddleware func(next ToolHandler) ToolHandler

func ChainToolMiddleware(middlewares ...ToolMiddleware) ToolMiddleware {
	return func(next ToolHandler) ToolHandler {
		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}
		return next
	}
}

type TracingConfig struct {
	StartTrace  func(ctx context.Context, callID, toolName, toolArgs string) error
	FinishTrace func(ctx context.Context, callID, toolName, result string, runErr error) error
}

func WithTracing(config TracingConfig) ToolMiddleware {
	return func(next ToolHandler) ToolHandler {
		return func(ctx context.Context, callID, toolName, toolArgs string) (string, error) {
			if config.StartTrace != nil {
				if err := config.StartTrace(ctx, callID, toolName, toolArgs); err != nil {
					return "", err
				}
			}
			result, err := next(ctx, callID, toolName, toolArgs)
			if config.FinishTrace != nil {
				_ = config.FinishTrace(ctx, callID, toolName, result, err)
			}
			return result, err
		}
	}
}

type ApprovalConfig struct {
	Manager        *tool_approval.Runtime
	Storage        *storage.Storage
	CreateApproval func(ctx context.Context, callID, toolName, toolArgs string) (*ApprovalHandle, error)
	HandleDecision func(ctx context.Context, callID, toolName string, decision tool_approval.WaitResult) (string, error)
}

type ApprovalHandle struct {
	ApprovalID string
	ToolCallID string
	ToolName   string
}

func WithApproval(config ApprovalConfig) ToolMiddleware {
	return func(next ToolHandler) ToolHandler {
		return func(ctx context.Context, callID, toolName, toolArgs string) (string, error) {
			registeredTool, ok := tools.ToolRouter.GetToolByID(toolName)
			if ok && registeredTool.RequireConfirmation() {
				handle, err := config.CreateApproval(ctx, callID, toolName, toolArgs)
				if err != nil {
					return "", err
				}
				result, waitErr := config.Manager.Wait(ctx, handle.ApprovalID)
				if waitErr != nil {
					return "", waitErr
				}
				return config.HandleDecision(ctx, callID, toolName, result)
			}
			return next(ctx, callID, toolName, toolArgs)
		}
	}
}

func WithTimeout(timeout time.Duration) ToolMiddleware {
	return func(next ToolHandler) ToolHandler {
		return func(ctx context.Context, callID, toolName, toolArgs string) (string, error) {
			ctx, cancel := context.WithTimeout(ctx, timeout)
			defer cancel()
			return next(ctx, callID, toolName, toolArgs)
		}
	}
}

func WithRetry(maxRetries int, backoff time.Duration) ToolMiddleware {
	return func(next ToolHandler) ToolHandler {
		return func(ctx context.Context, callID, toolName, toolArgs string) (result string, err error) {
			for attempt := 0; attempt <= maxRetries; attempt++ {
				if attempt > 0 {
					select {
					case <-ctx.Done():
						return "", ctx.Err()
					case <-time.After(backoff * time.Duration(1<<(attempt-1))):
					}
				}
				result, err = next(ctx, callID, toolName, toolArgs)
				if err == nil {
					return result, nil
				}
				if !isRetryable(err) {
					return "", err
				}
			}
			return
		}
	}
}

func isRetryable(err error) bool {
	if err == nil {
		return false
	}
	msg := err.Error()
	retryable := []string{"timeout", "temporary", "connection", "EOF", "reset"}
	for _, keyword := range retryable {
		if len(msg) > 0 && containsSubstring(msg, keyword) {
			return true
		}
	}
	return false
}

func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && searchSubstring(s, substr)
}

func searchSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func createApprovalMiddleware(
	runner *completionRunner,
	mgr *tool_approval.Runtime,
	st *storage.Storage,
) ToolMiddleware {
	config := ApprovalConfig{
		Manager: mgr,
		Storage: st,
		CreateApproval: func(ctx context.Context, callID, toolName, toolArgs string) (*ApprovalHandle, error) {
			return runner.createToolApproval(ctx, callID, toolName, toolArgs)
		},
		HandleDecision: func(ctx context.Context, callID, toolName string, decision tool_approval.WaitResult) (string, error) {
			return runner.handleApprovalDecision(ctx, callID, toolName, decision)
		},
	}
	return WithApproval(config)
}

func BuildStandardMiddlewareChain(
	runner *completionRunner,
	approvalMgr *tool_approval.Runtime,
	storage *storage.Storage,
	toolTimeout time.Duration,
) compose.ToolMiddleware {
	return compose.ToolMiddleware{}
}

func createTraceMiddleware(runner *completionRunner) ToolMiddleware {
	return WithTracing(TracingConfig{
		StartTrace: func(ctx context.Context, callID, toolName, toolArgs string) error {
			return nil
		},
		FinishTrace: func(ctx context.Context, callID, toolName, result string, runErr error) error {
			return nil
		},
	})
}
