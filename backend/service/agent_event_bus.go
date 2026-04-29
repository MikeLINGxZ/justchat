package service

import (
	"sync"
	"time"
)

type AgentEventType string

const (
	EventAgentThinking        AgentEventType = "agent.thinking"
	EventContentChunk         AgentEventType = "content.chunk"
	EventContentComplete      AgentEventType = "content.complete"
	EventToolCallStarted      AgentEventType = "tool.call_started"
	EventToolCallCompleted    AgentEventType = "tool.call_completed"
	EventToolAwaitingApproval AgentEventType = "tool.awaiting_approval"
	EventToolApprovalResolved AgentEventType = "tool.approval_resolved"
	EventStageChanged         AgentEventType = "stage.changed"
	EventTaskStateChanged     AgentEventType = "task.state_changed"
	EventError                AgentEventType = "error"
	EventWorkflowHandoff      AgentEventType = "workflow.handoff"
	EventWorkflowPlanReady    AgentEventType = "workflow.plan_ready"
	EventWorkflowBatchDone    AgentEventType = "workflow.batch_done"
	EventWorkflowSyncDone     AgentEventType = "workflow.synth_done"
	EventWorkflowReviewDone   AgentEventType = "workflow.review_done"
	EventTerminal             AgentEventType = "terminal"
)

type AgentEvent struct {
	Type      AgentEventType
	Timestamp time.Time
	Payload   interface{}
	Metadata  map[string]interface{}
}

type AgentEventBus struct {
	subscribers map[AgentEventType][]chan AgentEvent
	mu          sync.RWMutex
	bufferSize  int
	closed      bool
}

func NewAgentEventBus() *AgentEventBus {
	return &AgentEventBus{
		subscribers: make(map[AgentEventType][]chan AgentEvent),
		bufferSize:  256,
	}
}

func (bus *AgentEventBus) Subscribe(eventType AgentEventType) <-chan AgentEvent {
	bus.mu.Lock()
	defer bus.mu.Unlock()
	ch := make(chan AgentEvent, bus.bufferSize)
	bus.subscribers[eventType] = append(bus.subscribers[eventType], ch)
	return ch
}

func (bus *AgentEventBus) Unsubscribe(eventType AgentEventType, ch <-chan AgentEvent) {
	bus.mu.Lock()
	defer bus.mu.Unlock()
	subs := bus.subscribers[eventType]
	for i, sub := range subs {
		if sub == ch {
			bus.subscribers[eventType] = append(subs[:i], subs[i+1:]...)
			close(sub)
			return
		}
	}
}

func (bus *AgentEventBus) Publish(event AgentEvent) {
	bus.mu.RLock()
	defer bus.mu.RUnlock()
	if bus.closed {
		return
	}
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now()
	}
	subs := bus.subscribers[event.Type]
	for _, ch := range subs {
		select {
		case ch <- event:
		default:
		}
	}
}

func (bus *AgentEventBus) PublishAsync(event AgentEvent) {
	go bus.Publish(event)
}

func (bus *AgentEventBus) Close() {
	bus.mu.Lock()
	defer bus.mu.Unlock()
	bus.closed = true
	for _, subs := range bus.subscribers {
		for _, ch := range subs {
			close(ch)
		}
	}
	bus.subscribers = make(map[AgentEventType][]chan AgentEvent)
}

type EventCollector struct {
	bus    *AgentEventBus
	events []AgentEvent
	mu     sync.Mutex
	closed bool
}

func NewEventCollector(bus *AgentEventBus) *EventCollector {
	return &EventCollector{bus: bus}
}

func (c *EventCollector) Start(types ...AgentEventType) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, eventType := range types {
		ch := c.bus.Subscribe(eventType)
		go func() {
			for event := range ch {
				c.mu.Lock()
				if !c.closed {
					c.events = append(c.events, event)
				}
				c.mu.Unlock()
			}
		}()
	}
}

func (c *EventCollector) Flush() []AgentEvent {
	c.mu.Lock()
	defer c.mu.Unlock()
	result := c.events
	c.events = nil
	return result
}

func (c *EventCollector) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.closed = true
}
