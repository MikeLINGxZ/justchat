package terminal

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/ierror"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/terminal/terminal_dto"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/storage"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// newTestTerminalService creates an isolated terminal service for service tests.
func newTestTerminalService(t *testing.T) (*Terminal, *storage.Storage) {
	t.Helper()
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	stor, err := storage.NewStorageFromDB(db)
	if err != nil {
		t.Fatalf("new storage: %v", err)
	}
	return NewTerminal(stor), stor
}

// TestListTerminalsReturnsPersistedItems verifies the Wails method uses its DTO shape.
func TestListTerminalsReturnsPersistedItems(t *testing.T) {
	svc, stor := newTestTerminalService(t)
	_, err := stor.CreateTerminal(data_models.Terminal{
		TerminalID: "term_service",
		SessionID:  12,
		Title:      "Service Terminal",
		Command:    "sh",
		Status:     "done",
	})
	if err != nil {
		t.Fatalf("create terminal: %v", err)
	}

	out, err := svc.ListTerminals(context.Background(), terminal_dto.ListTerminalsInput{SessionID: 12})
	if err != nil {
		t.Fatalf("list terminals: %v", err)
	}
	if len(out.Items) != 1 || out.Items[0].ID != "term_service" {
		t.Fatalf("unexpected items: %+v", out.Items)
	}
}

// TestWriteTerminalInputWrapsInactiveTerminalError verifies public errors use ierror.
func TestWriteTerminalInputWrapsInactiveTerminalError(t *testing.T) {
	svc, _ := newTestTerminalService(t)

	_, err := svc.WriteTerminalInput(context.Background(), terminal_dto.WriteTerminalInputInput{
		TerminalID: "missing",
		Data:       "hello\n",
	})
	if err == nil {
		t.Fatal("expected error")
	}

	var iErr *ierror.IError
	if !errors.As(err, &iErr) {
		t.Fatalf("expected ierror, got %T", err)
	}
	if iErr.Msg != ierror.ErrTerminalWriteInput.Msg() {
		t.Fatalf("unexpected ierror message: %q", iErr.Msg)
	}
}
