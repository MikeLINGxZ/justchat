package service

import (
	"fmt"
	"sync/atomic"
)

type TaskState string

const (
	StateIdle              TaskState = "idle"
	StatePreparing         TaskState = "preparing"
	StateRunning           TaskState = "running"
	StateStreaming         TaskState = "streaming"
	StateAwaitingApproval  TaskState = "awaiting_approval"
	StateExecutingTools    TaskState = "executing_tools"
	StateWorkflowPlanning  TaskState = "workflow_planning"
	StateWorkflowExecuting TaskState = "workflow_executing"
	StateWorkflowReviewing TaskState = "workflow_reviewing"
	StateCompleted         TaskState = "completed"
	StateStopped           TaskState = "stopped"
	StateFailed            TaskState = "failed"
)

type StateTransition struct {
	From TaskState
	To   TaskState
}

type transitionMetadata struct {
	Reason string
	Error  string
}

type TaskFSM struct {
	state        atomic.Value
	onTransition func(from, to TaskState, metadata *transitionMetadata)
}

var validTransitions = map[TaskState][]TaskState{
	StateIdle:              {StatePreparing, StateStopped},
	StatePreparing:         {StateRunning, StateFailed, StateStopped},
	StateRunning:           {StateStreaming, StateExecutingTools, StateAwaitingApproval, StateWorkflowPlanning, StateCompleted, StateFailed, StateStopped},
	StateStreaming:         {StateRunning, StateAwaitingApproval, StateExecutingTools, StateWorkflowPlanning, StateCompleted, StateFailed, StateStopped},
	StateAwaitingApproval:  {StateRunning, StateExecutingTools, StateCompleted, StateFailed, StateStopped},
	StateExecutingTools:    {StateRunning, StateStreaming, StateCompleted, StateFailed, StateStopped},
	StateWorkflowPlanning:  {StateWorkflowExecuting, StateFailed, StateStopped},
	StateWorkflowExecuting: {StateWorkflowReviewing, StateFailed, StateStopped},
	StateWorkflowReviewing: {StateWorkflowExecuting, StateCompleted, StateFailed, StateStopped},
	StateCompleted:         {},
	StateStopped:           {},
	StateFailed:            {},
}

var terminalStates = map[TaskState]bool{
	StateCompleted: true,
	StateStopped:   true,
	StateFailed:    true,
}

func NewTaskFSM(initialState TaskState) *TaskFSM {
	fsm := &TaskFSM{}
	fsm.state.Store(initialState)
	return fsm
}

func (fsm *TaskFSM) SetOnTransition(fn func(from, to TaskState, metadata *transitionMetadata)) {
	fsm.onTransition = fn
}

func (fsm *TaskFSM) Current() TaskState {
	return fsm.state.Load().(TaskState)
}

func (fsm *TaskFSM) IsTerminal() bool {
	return terminalStates[fsm.Current()]
}

func (fsm *TaskFSM) CanTransition(target TaskState) bool {
	current := fsm.Current()
	for _, valid := range validTransitions[current] {
		if valid == target {
			return true
		}
	}
	return false
}

func (fsm *TaskFSM) Transition(target TaskState, reason string, errMsg string) error {
	current := fsm.Current()
	if terminalStates[current] {
		return fmt.Errorf("cannot transition from terminal state %q to %q", current, target)
	}
	if !fsm.CanTransition(target) {
		return fmt.Errorf("invalid transition from %q to %q", current, target)
	}
	fsm.state.Store(target)
	meta := &transitionMetadata{Reason: reason, Error: errMsg}
	if fsm.onTransition != nil {
		fsm.onTransition(current, target, meta)
	}
	return nil
}

func (fsm *TaskFSM) TransitionToTerminal(target TaskState, reason string, errMsg string) error {
	if !terminalStates[target] {
		return fmt.Errorf("target %q is not a terminal state", target)
	}
	return fsm.Transition(target, reason, errMsg)
}
