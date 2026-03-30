package tool_approval

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
)

type ApprovalPrompt struct {
	Title   string
	Message string
	Scope   string
}

type ApprovalAwareTool interface {
	BuildApprovalPrompt(ctx context.Context, argumentsJSON string) (*ApprovalPrompt, error)
}

type WaitResult struct {
	ApprovalID  string
	Decision    data_models.ToolApprovalDecision
	Comment     string
	RespondedAt time.Time
}

type Runtime struct {
	mu      sync.Mutex
	waiters map[string]chan WaitResult
}

func NewRuntime() *Runtime {
	return &Runtime{
		waiters: map[string]chan WaitResult{},
	}
}

var Manager = NewRuntime()

func (r *Runtime) Register(approvalID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if approvalID == "" {
		return fmt.Errorf("approval id is empty")
	}
	if _, exists := r.waiters[approvalID]; exists {
		return fmt.Errorf("approval %s is already waiting", approvalID)
	}
	r.waiters[approvalID] = make(chan WaitResult, 1)
	return nil
}

func (r *Runtime) HasWaiter(approvalID string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.waiters[approvalID]
	return ok
}

func (r *Runtime) Wait(ctx context.Context, approvalID string) (WaitResult, error) {
	r.mu.Lock()
	waitCh, ok := r.waiters[approvalID]
	r.mu.Unlock()
	if !ok {
		return WaitResult{}, fmt.Errorf("approval %s is not waiting", approvalID)
	}

	select {
	case result := <-waitCh:
		r.mu.Lock()
		delete(r.waiters, approvalID)
		r.mu.Unlock()
		return result, nil
	case <-ctx.Done():
		r.mu.Lock()
		delete(r.waiters, approvalID)
		r.mu.Unlock()
		return WaitResult{}, ctx.Err()
	}
}

func (r *Runtime) Resolve(result WaitResult) error {
	r.mu.Lock()
	waitCh, ok := r.waiters[result.ApprovalID]
	r.mu.Unlock()
	if !ok {
		return fmt.Errorf("approval %s has no active waiter", result.ApprovalID)
	}

	select {
	case waitCh <- result:
		return nil
	default:
		return fmt.Errorf("approval %s waiter is not ready", result.ApprovalID)
	}
}

func (r *Runtime) Cancel(approvalID string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.waiters, approvalID)
}
