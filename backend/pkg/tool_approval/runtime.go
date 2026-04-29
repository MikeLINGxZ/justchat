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

type ApprovalCallback func(result WaitResult)

type Runtime struct {
	mu           sync.Mutex
	waiters      map[string]chan WaitResult
	asyncWaiters map[string]ApprovalCallback
}

func NewRuntime() *Runtime {
	return &Runtime{
		waiters:      map[string]chan WaitResult{},
		asyncWaiters: map[string]ApprovalCallback{},
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
	delete(r.asyncWaiters, approvalID)
}

func (r *Runtime) RegisterAsync(approvalID string, callback ApprovalCallback) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if approvalID == "" {
		return fmt.Errorf("approval id is empty")
	}
	if _, exists := r.waiters[approvalID]; exists {
		return fmt.Errorf("approval %s is already waiting", approvalID)
	}
	if _, exists := r.asyncWaiters[approvalID]; exists {
		return fmt.Errorf("approval %s is already waiting (async)", approvalID)
	}
	r.asyncWaiters[approvalID] = callback
	return nil
}

func (r *Runtime) WaitMulti(ctx context.Context, approvalIDs []string) (map[string]WaitResult, error) {
	if len(approvalIDs) == 0 {
		return nil, nil
	}

	results := make(map[string]WaitResult, len(approvalIDs))
	resultCh := make(chan WaitResult, len(approvalIDs))

	for _, id := range approvalIDs {
		go func(approvalID string) {
			result, err := r.Wait(ctx, approvalID)
			if err != nil {
				result = WaitResult{ApprovalID: approvalID}
			}
			resultCh <- result
		}(id)
	}

	remaining := len(approvalIDs)
	for remaining > 0 {
		select {
		case result := <-resultCh:
			results[result.ApprovalID] = result
			remaining--
		case <-ctx.Done():
			return results, ctx.Err()
		}
	}

	return results, nil
}

func (r *Runtime) ResolveAsync(result WaitResult) error {
	r.mu.Lock()
	callback, ok := r.asyncWaiters[result.ApprovalID]
	if ok {
		delete(r.asyncWaiters, result.ApprovalID)
	}
	waitCh, syncOk := r.waiters[result.ApprovalID]
	r.mu.Unlock()

	if syncOk {
		select {
		case waitCh <- result:
			return nil
		default:
			if callback == nil {
				return fmt.Errorf("approval %s waiter is not ready", result.ApprovalID)
			}
		}
	}

	if callback != nil {
		callback(result)
		return nil
	}

	if !syncOk {
		return fmt.Errorf("approval %s has no active waiter", result.ApprovalID)
	}
	return nil
}

func (r *Runtime) RegisterWithCallback(ctx context.Context, approvalID string) (<-chan WaitResult, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if approvalID == "" {
		return nil, fmt.Errorf("approval id is empty")
	}
	if _, exists := r.waiters[approvalID]; exists {
		return nil, fmt.Errorf("approval %s is already waiting", approvalID)
	}
	if _, exists := r.asyncWaiters[approvalID]; exists {
		return nil, fmt.Errorf("approval %s is already waiting (async)", approvalID)
	}

	ch := make(chan WaitResult, 1)
	r.waiters[approvalID] = ch
	return ch, nil
}
