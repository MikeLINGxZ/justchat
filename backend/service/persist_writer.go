package service

import (
	"context"
	"sync"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"
)

const (
	defaultFlushInterval  = 50 * time.Millisecond
	defaultFlushMinDelta  = 200
	defaultContentBufSize = 256
	defaultEventBufSize   = 128
)

type contentChunk struct {
	resultCh chan error
}

type persistEvent struct {
	kind     string
	payload  interface{}
	resultCh chan error
}

type PersistWriter struct {
	svc    *Service
	runner *completionRunner

	contentCh chan contentChunk
	eventCh   chan persistEvent

	flushInterval time.Duration
	flushMinDelta int

	lastPersistAt         time.Time
	lastPersistContentLen int

	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
	closed bool
}

func NewPersistWriter(svc *Service, runner *completionRunner) *PersistWriter {
	return &PersistWriter{
		svc:           svc,
		runner:        runner,
		contentCh:     make(chan contentChunk, defaultContentBufSize),
		eventCh:       make(chan persistEvent, defaultEventBufSize),
		flushInterval: defaultFlushInterval,
		flushMinDelta: defaultFlushMinDelta,
	}
}

func (w *PersistWriter) Start(parentCtx context.Context) {
	w.ctx, w.cancel = context.WithCancel(parentCtx)
	w.wg.Add(1)
	go w.loop()
}

func (w *PersistWriter) Stop() {
	if w.cancel != nil {
		w.cancel()
	}
	w.wg.Wait()
}

func (w *PersistWriter) NotifyContent() error {
	if w.closed {
		return nil
	}
	chunk := contentChunk{
		resultCh: make(chan error, 1),
	}
	select {
	case w.contentCh <- chunk:
		return <-chunk.resultCh
	case <-w.ctx.Done():
		return w.ctx.Err()
	}
}

func (w *PersistWriter) EnqueueEvent(kind string, payload interface{}) error {
	if w.closed {
		return nil
	}
	event := persistEvent{
		kind:     kind,
		payload:  payload,
		resultCh: make(chan error, 1),
	}
	select {
	case w.eventCh <- event:
		return <-event.resultCh
	case <-w.ctx.Done():
		return w.ctx.Err()
	}
}

func (w *PersistWriter) FlushSync() error {
	if w.closed {
		return nil
	}
	chunk := contentChunk{
		resultCh: make(chan error, 1),
	}
	select {
	case w.contentCh <- chunk:
		return <-chunk.resultCh
	case <-w.ctx.Done():
		return w.ctx.Err()
	}
}

func (w *PersistWriter) loop() {
	defer w.wg.Done()
	ticker := time.NewTicker(w.flushInterval)
	defer ticker.Stop()

	for {
		select {
		case <-w.ctx.Done():
			w.flushAll()
			return

		case chunk := <-w.contentCh:
			if chunk.resultCh != nil {
				chunk.resultCh <- nil
			}
			w.checkFlush()

		case event := <-w.eventCh:
			err := w.handleEvent(event)
			if event.resultCh != nil {
				event.resultCh <- err
			}

		case <-ticker.C:
			w.checkFlush()
		}
	}
}

func (w *PersistWriter) checkFlush() {
	w.runner.mu.Lock()
	currentContentLen := len([]rune(w.runner.assistantMessage.Content))
	shouldPersist := w.lastPersistAt.IsZero() ||
		time.Since(w.lastPersistAt) >= w.flushInterval ||
		currentContentLen-w.lastPersistContentLen >= w.flushMinDelta
	if shouldPersist {
		w.lastPersistAt = time.Now()
		w.lastPersistContentLen = currentContentLen
		w.runner.mu.Unlock()
		w.persistSnapshot(false)
		return
	}
	w.runner.mu.Unlock()
	w.emitSnapshotDirect()
}

func (w *PersistWriter) emitSnapshotDirect() {
	w.runner.mu.Lock()
	snapshot := w.runner.cloneAssistantMessageLocked()
	w.runner.mu.Unlock()
	w.emitSnapshotEvent(snapshot)
}

func (w *PersistWriter) handleEvent(event persistEvent) error {
	switch event.kind {
	case "trace_step":
		w.persistSnapshot(false)
	case "stage_change":
		w.persistSnapshot(true)
	case "snapshot":
		w.persistSnapshot(true)
	case "terminal":
		w.persistSnapshot(true)
	}
	return nil
}

func (w *PersistWriter) flushAll() {
	w.runner.mu.Lock()
	w.lastPersistAt = time.Now()
	w.lastPersistContentLen = len([]rune(w.runner.assistantMessage.Content))
	w.runner.mu.Unlock()
	w.persistSnapshot(true)
}

func (w *PersistWriter) persistSnapshot(updateTask bool) {
	w.runner.mu.Lock()
	msgToSave := w.runner.assistantMessage
	taskToSave := w.runner.task

	saveErr := w.svc.storage.SaveOrUpdateMessage(context.Background(), msgToSave)
	if saveErr != nil {
		logger.Error("persistWriter save message error: ", saveErr)
	}
	if updateTask {
		if taskErr := w.svc.storage.SaveTask(context.Background(), taskToSave); taskErr != nil {
			logger.Error("persistWriter save task error: ", taskErr)
		}
	}

	snapshot := w.runner.cloneAssistantMessageLocked()
	pendingDelta := make([]data_models.TraceStep, len(w.runner.pendingTraceDelta))
	copy(pendingDelta, w.runner.pendingTraceDelta)
	w.runner.pendingTraceDelta = nil
	w.runner.mu.Unlock()

	w.svc.emitTaskEvent(taskToSave, snapshot, pendingDelta)
}

func (w *PersistWriter) emitSnapshotEvent(snapshot data_models.Message) {
	w.runner.mu.Lock()
	taskSnapshot := w.runner.task
	pendingDelta := make([]data_models.TraceStep, len(w.runner.pendingTraceDelta))
	copy(pendingDelta, w.runner.pendingTraceDelta)
	w.runner.pendingTraceDelta = nil
	w.runner.mu.Unlock()

	w.svc.emitTaskEvent(taskSnapshot, snapshot, pendingDelta)
}
