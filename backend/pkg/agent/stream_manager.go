package agent

import (
	"context"
	"sync"
)

type activeStream struct {
	cancel    context.CancelFunc
	confirmCh chan ConfirmResponse
}

// ConfirmResponse carries the user's response to a tool confirmation request.
type ConfirmResponse struct {
	Approved bool
	Message  string
	Action   string
}

type StreamManager struct {
	mu      sync.RWMutex
	streams map[uint]*activeStream
}

func NewStreamManager() *StreamManager {
	return &StreamManager{
		streams: make(map[uint]*activeStream),
	}
}

func (sm *StreamManager) Start(parentCtx context.Context, sessionID uint) (context.Context, context.CancelFunc) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if existing, ok := sm.streams[sessionID]; ok {
		existing.cancel()
	}

	ctx, cancel := context.WithCancel(parentCtx)
	sm.streams[sessionID] = &activeStream{
		cancel:    cancel,
		confirmCh: make(chan ConfirmResponse, 1),
	}
	return ctx, cancel
}

func (sm *StreamManager) Stop(sessionID uint) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if s, ok := sm.streams[sessionID]; ok {
		s.cancel()
		close(s.confirmCh)
		delete(sm.streams, sessionID)
	}
}

func (sm *StreamManager) Remove(sessionID uint) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.streams, sessionID)
}

func (sm *StreamManager) IsActive(sessionID uint) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	_, ok := sm.streams[sessionID]
	return ok
}

func (sm *StreamManager) SendConfirmResponse(sessionID uint, resp ConfirmResponse) bool {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	s, ok := sm.streams[sessionID]
	if !ok {
		return false
	}
	select {
	case s.confirmCh <- resp:
		return true
	default:
		return false
	}
}

func (sm *StreamManager) WaitForConfirm(ctx context.Context, sessionID uint) (ConfirmResponse, error) {
	sm.mu.RLock()
	s, ok := sm.streams[sessionID]
	sm.mu.RUnlock()

	if !ok {
		return ConfirmResponse{}, context.Canceled
	}

	select {
	case <-ctx.Done():
		return ConfirmResponse{}, ctx.Err()
	case resp, chanOpen := <-s.confirmCh:
		if !chanOpen {
			return ConfirmResponse{}, context.Canceled
		}
		return resp, nil
	}
}
