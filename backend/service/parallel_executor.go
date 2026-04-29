package service

import (
	"context"
	"strings"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"
)

type toolCallRequest struct {
	CallID    string
	Name      string
	Arguments string
}

type toolCallResult struct {
	CallID string
	Name   string
	Result string
	Err    error
}

type ParallelExecutorConfig struct {
	MaxParallel int
	Timeout     time.Duration
}

func DefaultParallelExecutorConfig() ParallelExecutorConfig {
	return ParallelExecutorConfig{
		MaxParallel: 5,
		Timeout:     2 * time.Minute,
	}
}

type ParallelToolExecutor struct {
	config ParallelExecutorConfig
}

func NewParallelToolExecutor(config ParallelExecutorConfig) *ParallelToolExecutor {
	if config.MaxParallel <= 0 {
		config.MaxParallel = 5
	}
	if config.Timeout <= 0 {
		config.Timeout = 2 * time.Minute
	}
	return &ParallelToolExecutor{config: config}
}

func (e *ParallelToolExecutor) ExecuteParallel(
	ctx context.Context,
	toolCalls []toolCallRequest,
	executor func(ctx context.Context, call toolCallRequest) (string, error),
) ([]toolCallResult, error) {
	batches := batchByDependencies(toolCalls)
	allResults := make([]toolCallResult, 0, len(toolCalls))
	resultsByCallID := make(map[string]string)

	for _, batch := range batches {
		if len(batch) == 1 {
			call := batch[0]
			result, err := executor(ctx, call)
			allResults = append(allResults, toolCallResult{
				CallID: call.CallID,
				Name:   call.Name,
				Result: result,
				Err:    err,
			})
			if err != nil {
				result = err.Error()
			}
			resultsByCallID[call.CallID] = result
			if err != nil {
				return allResults, err
			}
			continue
		}

		g, gCtx := errgroup.WithContext(ctx)
		g.SetLimit(e.config.MaxParallel)
		mu := sync.Mutex{}

		batchResults := make([]toolCallResult, len(batch))
		for i, call := range batch {
			i := i
			call := call
			g.Go(func() error {
				result, err := executor(gCtx, call)
				mu.Lock()
				batchResults[i] = toolCallResult{
					CallID: call.CallID,
					Name:   call.Name,
					Result: result,
					Err:    err,
				}
				if err != nil {
					result = err.Error()
				}
				resultsByCallID[call.CallID] = result
				mu.Unlock()
				return err
			})
		}

		err := g.Wait()
		allResults = append(allResults, batchResults...)
		if err != nil {
			return allResults, err
		}
	}

	return allResults, nil
}

func batchByDependencies(toolCalls []toolCallRequest) [][]toolCallRequest {
	batches := [][]toolCallRequest{}
	completed := make(map[string]bool)

	for len(completed) < len(toolCalls) {
		batch := []toolCallRequest{}
		for _, call := range toolCalls {
			if completed[call.CallID] {
				continue
			}
			if hasUnmetDependencies(call, toolCalls, completed) {
				continue
			}
			batch = append(batch, call)
		}

		if len(batch) == 0 {
			for _, call := range toolCalls {
				if !completed[call.CallID] {
					batch = append(batch, call)
				}
			}
		}

		batches = append(batches, batch)
		for _, call := range batch {
			completed[call.CallID] = true
		}
	}

	return batches
}

func hasUnmetDependencies(call toolCallRequest, allCalls []toolCallRequest, completed map[string]bool) bool {
	for _, other := range allCalls {
		if completed[other.CallID] {
			continue
		}
		if other.CallID == call.CallID {
			continue
		}
		if referenceCallID(call.Arguments, other.CallID) {
			return true
		}
	}
	return false
}

func referenceCallID(args, callID string) bool {
	return strings.Contains(args, callID)
}

func (e *ParallelToolExecutor) AnalyzeDependencies(calls []toolCallRequest) (independent [][]toolCallRequest, dependent []toolCallRequest) {
	for _, call := range calls {
		hasDep := false
		for _, other := range calls {
			if other.CallID == call.CallID {
				continue
			}
			if referenceCallID(call.Arguments, other.CallID) {
				hasDep = true
				break
			}
		}
		if hasDep {
			dependent = append(dependent, call)
		} else {
			independent = append(independent, []toolCallRequest{call})
		}
	}
	if len(independent) > 0 && len(independent[0]) > 1 {
		independent = batchByDependencies(independent[0])
	}
	return
}

func (e *ParallelToolExecutor) HasIndependentCalls(calls []toolCallRequest) bool {
	independent, _ := e.AnalyzeDependencies(calls)
	totalIndependent := 0
	for _, batch := range independent {
		totalIndependent += len(batch)
	}
	return totalIndependent > 1
}
