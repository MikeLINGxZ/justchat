package terminal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"sync"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/storage"
)

type OutputEvent struct {
	TerminalID string `json:"terminal_id"`
	SessionID  uint   `json:"session_id"`
	Cursor     int64  `json:"cursor"`
	Data       string `json:"data"`
	Status     string `json:"status"`
	Visible    bool   `json:"visible"`
}

type StatusEvent struct {
	TerminalID string `json:"terminal_id"`
	SessionID  uint   `json:"session_id"`
	Status     string `json:"status"`
	ExitCode   int    `json:"exit_code"`
	Visible    *bool  `json:"visible,omitempty"`
	Error      string `json:"error,omitempty"`
}

type TerminalOutputChunk = data_models.TerminalOutputChunk

type Info struct {
	ID            string `json:"id"`
	SessionID     uint   `json:"session_id"`
	MessageID     *uint  `json:"message_id,omitempty"`
	ToolCallID    string `json:"tool_call_id,omitempty"`
	Title         string `json:"title"`
	Command       string `json:"command"`
	Args          string `json:"args"`
	Cwd           string `json:"cwd"`
	Status        string `json:"status"`
	Visible       bool   `json:"visible"`
	PID           int    `json:"pid"`
	ExitCode      *int   `json:"exit_code,omitempty"`
	CurrentCursor int64  `json:"current_cursor"`
}

type CreateParams struct {
	SessionID  uint
	MessageID  *uint
	ToolCallID string
	Title      string
	Command    string
	Args       []string
	Env        []string
	Cwd        string
	Visible    bool
	Rows       uint16
	Cols       uint16
}

type session struct {
	info Info
	proc ptyProcess
	mu   sync.Mutex
}

type Manager struct {
	store        *storage.Storage
	emitOutput   func(OutputEvent)
	emitStatus   func(StatusEvent)
	mu           sync.RWMutex
	sessions     map[string]*session
	nextID       func() string
	startProcess func(context.Context, CreateParams) (ptyProcess, int, error)
}

// NewManager creates a PTY terminal manager backed by persistent storage.
func NewManager(store *storage.Storage, emitOutput func(OutputEvent)) *Manager {
	return &Manager{
		store:        store,
		emitOutput:   emitOutput,
		sessions:     make(map[string]*session),
		nextID:       defaultID,
		startProcess: startPTYProcess,
	}
}

// SetOutputEmitter installs the callback used to stream output chunks to the UI.
func (m *Manager) SetOutputEmitter(emit func(OutputEvent)) {
	m.emitOutput = emit
}

// SetStatusEmitter installs the callback used to report terminal status changes.
func (m *Manager) SetStatusEmitter(emit func(StatusEvent)) {
	m.emitStatus = emit
}

// Create starts a new PTY process and persists its terminal metadata.
func (m *Manager) Create(ctx context.Context, params CreateParams) (Info, error) {
	if params.Command == "" {
		return Info{}, errors.New("terminal: command required")
	}
	id := m.nextID()
	proc, pid, err := m.startProcess(ctx, params)
	if err != nil {
		return Info{}, err
	}
	argsJSON, _ := json.Marshal(params.Args)
	title := params.Title
	if title == "" {
		title = params.Command
	}
	record, err := m.store.CreateTerminal(data_models.Terminal{
		TerminalID: id,
		SessionID:  params.SessionID,
		MessageID:  params.MessageID,
		ToolCallID: params.ToolCallID,
		Title:      title,
		Command:    params.Command,
		Args:       string(argsJSON),
		Cwd:        params.Cwd,
		Status:     "active",
		Visible:    params.Visible,
		PID:        pid,
		StartedAt:  time.Now(),
	})
	if err != nil {
		_ = proc.Close()
		return Info{}, err
	}
	info := terminalInfo(*record)
	sess := &session{info: info, proc: proc}

	m.mu.Lock()
	m.sessions[id] = sess
	m.mu.Unlock()

	go m.drain(id, sess)
	return info, nil
}

// CreateTerminal exposes Create using the method name expected by agent tools.
func (m *Manager) CreateTerminal(ctx context.Context, params CreateParams) (Info, error) {
	return m.Create(ctx, params)
}

// Write sends bytes to an active PTY session.
func (m *Manager) Write(terminalID string, data string) error {
	m.mu.RLock()
	sess := m.sessions[terminalID]
	m.mu.RUnlock()
	if sess == nil {
		return errors.New("terminal: session not active")
	}
	sess.mu.Lock()
	defer sess.mu.Unlock()
	return sess.proc.Write([]byte(data))
}

// WriteTerminal exposes Write using the method name expected by Wails services.
func (m *Manager) WriteTerminal(terminalID string, data string) error {
	return m.Write(terminalID, data)
}

// Resize changes the active PTY dimensions.
func (m *Manager) Resize(terminalID string, rows, cols uint16) error {
	m.mu.RLock()
	sess := m.sessions[terminalID]
	m.mu.RUnlock()
	if sess == nil {
		return nil
	}
	return sess.proc.Resize(rows, cols)
}

// SetVisible persists terminal visibility without emitting a runtime event.
func (m *Manager) SetVisible(terminalID string, visible bool, title string) error {
	return m.store.UpdateTerminalVisibility(terminalID, visible, title)
}

// SetTerminalVisible updates terminal visibility and emits the resulting status event.
func (m *Manager) SetTerminalVisible(ctx context.Context, terminalID string, visible bool, title string) error {
	_ = ctx
	if err := m.store.UpdateTerminalVisibility(terminalID, visible, title); err != nil {
		return err
	}
	m.mu.RLock()
	sess := m.sessions[terminalID]
	m.mu.RUnlock()
	if sess != nil {
		sess.mu.Lock()
		sess.info.Visible = visible
		if title != "" {
			sess.info.Title = title
		}
		sess.mu.Unlock()
	}
	record, err := m.store.GetTerminal(terminalID)
	if err != nil {
		return err
	}
	if m.emitStatus != nil {
		exitCode := 0
		if record.ExitCode != nil {
			exitCode = *record.ExitCode
		}
		m.emitStatus(StatusEvent{
			TerminalID: terminalID,
			SessionID:  record.SessionID,
			Status:     record.Status,
			ExitCode:   exitCode,
			Visible:    &visible,
		})
	}
	return nil
}

// ListBySession returns all persisted terminals for a chat session.
func (m *Manager) ListBySession(sessionID uint) ([]Info, error) {
	records, err := m.store.ListTerminalsForSession(sessionID)
	if err != nil {
		return nil, err
	}
	out := make([]Info, 0, len(records))
	for _, record := range records {
		out = append(out, terminalInfo(record))
	}
	return out, nil
}

// ReadOutput returns persisted terminal output after a byte cursor.
func (m *Manager) ReadOutput(terminalID string, cursor int64) ([]data_models.TerminalOutputChunk, error) {
	return m.store.ReadTerminalOutput(terminalID, cursor)
}

// ReadTerminalOutput exposes ReadOutput using the method name expected by agent tools.
func (m *Manager) ReadTerminalOutput(ctx context.Context, terminalID string, cursor int64) ([]data_models.TerminalOutputChunk, error) {
	_ = ctx
	return m.ReadOutput(terminalID, cursor)
}

// Wait blocks until the terminal exits and returns the final persisted metadata.
func (m *Manager) Wait(terminalID string) (Info, error) {
	m.mu.RLock()
	sess := m.sessions[terminalID]
	m.mu.RUnlock()
	if sess == nil {
		record, err := m.store.GetTerminal(terminalID)
		if err != nil {
			return Info{}, err
		}
		return terminalInfo(*record), nil
	}
	<-sess.proc.Done()
	record, err := m.store.GetTerminal(terminalID)
	if err != nil {
		return sess.info, err
	}
	return terminalInfo(*record), nil
}

// WaitTerminal exposes Wait using the method name expected by agent tools.
func (m *Manager) WaitTerminal(ctx context.Context, terminalID string) (Info, error) {
	_ = ctx
	return m.Wait(terminalID)
}

// drain persists PTY output chunks and emits terminal lifecycle events.
func (m *Manager) drain(id string, sess *session) {
	for chunk := range sess.proc.Output() {
		stored, err := m.store.AppendTerminalOutput(id, string(chunk))
		if err != nil {
			continue
		}
		if m.emitOutput != nil {
			// Snapshot mutable session fields under the session lock so event
			// visibility stays consistent with SetTerminalVisible updates.
			sess.mu.Lock()
			sessionID := sess.info.SessionID
			visible := sess.info.Visible
			sess.mu.Unlock()
			m.emitOutput(OutputEvent{
				TerminalID: id,
				SessionID:  sessionID,
				Cursor:     stored.CursorEnd,
				Data:       stored.Data,
				Status:     "active",
				Visible:    visible,
			})
		}
	}
	exitCode, waitErr := sess.proc.Wait()
	status := "done"
	if waitErr != nil || exitCode != 0 {
		status = "error"
	}
	_ = m.store.FinishTerminal(id, status, exitCode)

	m.mu.Lock()
	delete(m.sessions, id)
	m.mu.Unlock()

	if m.emitStatus != nil {
		errString := ""
		if waitErr != nil {
			errString = waitErr.Error()
		}
		m.emitStatus(StatusEvent{
			TerminalID: id,
			SessionID:  sess.info.SessionID,
			Status:     status,
			ExitCode:   exitCode,
			Error:      errString,
		})
	}
}

// terminalInfo converts a storage model into the API-facing terminal metadata.
func terminalInfo(record data_models.Terminal) Info {
	return Info{
		ID:            record.TerminalID,
		SessionID:     record.SessionID,
		MessageID:     record.MessageID,
		ToolCallID:    record.ToolCallID,
		Title:         record.Title,
		Command:       record.Command,
		Args:          record.Args,
		Cwd:           record.Cwd,
		Status:        record.Status,
		Visible:       record.Visible,
		PID:           record.PID,
		ExitCode:      record.ExitCode,
		CurrentCursor: record.CurrentCursor,
	}
}

// defaultID generates a reasonably unique terminal ID for persisted rows.
func defaultID() string {
	return fmt.Sprintf("term_%d", time.Now().UnixNano())
}

type ptyProcess interface {
	Output() <-chan []byte
	Write([]byte) error
	Resize(rows, cols uint16) error
	Close() error
	Wait() (int, error)
	Done() <-chan struct{}
}

// buildCommand creates the exec.Cmd that will run inside a PTY.
func buildCommand(_ context.Context, params CreateParams) *exec.Cmd {
	cmd := exec.Command(params.Command, params.Args...)
	if params.Env != nil {
		cmd.Env = params.Env
	}
	if params.Cwd != "" {
		cmd.Dir = params.Cwd
	}
	return cmd
}
