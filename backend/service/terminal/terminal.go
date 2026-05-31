package terminal

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/ierror"
	pkgterminal "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/terminal"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/terminal/terminal_dto"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/storage"
)

// Terminal exposes interactive terminal operations to the frontend.
type Terminal struct {
	manager  *pkgterminal.Manager
	wailsApp *application.App
}

// NewTerminal creates a terminal service using the shared application storage.
func NewTerminal(store *storage.Storage) *Terminal {
	return NewTerminalWithManager(pkgterminal.NewManager(store, nil))
}

// NewTerminalWithManager creates a terminal service from an existing manager.
func NewTerminalWithManager(manager *pkgterminal.Manager) *Terminal {
	return &Terminal{manager: manager}
}

// ListTerminals returns visible and persisted terminal sessions for a chat session.
func (t *Terminal) ListTerminals(ctx context.Context, input terminal_dto.ListTerminalsInput) (*terminal_dto.ListTerminalsOutput, error) {
	_ = ctx
	items, err := t.manager.ListBySession(input.SessionID)
	if err != nil {
		return nil, ierror.Error(ierror.ErrTerminalList, err)
	}
	return &terminal_dto.ListTerminalsOutput{Items: items}, nil
}

// ReadTerminalOutput returns output chunks after the supplied byte cursor.
func (t *Terminal) ReadTerminalOutput(ctx context.Context, input terminal_dto.ReadTerminalOutputInput) (*terminal_dto.ReadTerminalOutputOutput, error) {
	_ = ctx
	chunks, err := t.manager.ReadOutput(input.TerminalID, input.Cursor)
	if err != nil {
		return nil, ierror.Error(ierror.ErrTerminalReadOutput, err)
	}
	return &terminal_dto.ReadTerminalOutputOutput{Chunks: chunks}, nil
}

// WriteTerminalInput sends user keystrokes or pasted text to an active terminal.
func (t *Terminal) WriteTerminalInput(ctx context.Context, input terminal_dto.WriteTerminalInputInput) (*terminal_dto.WriteTerminalInputOutput, error) {
	_ = ctx
	if err := t.manager.Write(input.TerminalID, input.Data); err != nil {
		return nil, ierror.Error(ierror.ErrTerminalWriteInput, err)
	}
	return &terminal_dto.WriteTerminalInputOutput{}, nil
}

// ResizeTerminal updates the PTY size to match the frontend terminal viewport.
func (t *Terminal) ResizeTerminal(ctx context.Context, input terminal_dto.ResizeTerminalInput) (*terminal_dto.ResizeTerminalOutput, error) {
	_ = ctx
	if err := t.manager.Resize(input.TerminalID, input.Rows, input.Cols); err != nil {
		return nil, ierror.Error(ierror.ErrTerminalResize, err)
	}
	return &terminal_dto.ResizeTerminalOutput{}, nil
}
