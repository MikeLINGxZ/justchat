//go:build !windows

package terminal

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/storage"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// newTerminalTestStorage creates isolated storage for terminal manager tests.
func newTerminalTestStorage(t *testing.T) *storage.Storage {
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
	return stor
}

// writeTerminalTestScript writes an executable script used by PTY tests.
func writeTerminalTestScript(t *testing.T, body string) string {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("PTY test uses POSIX shell script")
	}
	path := filepath.Join(t.TempDir(), "script.sh")
	if err := os.WriteFile(path, []byte("#!/usr/bin/env bash\nset -e\n"+body+"\n"), 0o755); err != nil {
		t.Fatalf("write script: %v", err)
	}
	return path
}

// TestManagerCreateWriteAndPersistOutput verifies PTY output is persisted.
func TestManagerCreateWriteAndPersistOutput(t *testing.T) {
	stor := newTerminalTestStorage(t)
	var emitted []OutputEvent
	mgr := NewManager(stor, func(event OutputEvent) {
		emitted = append(emitted, event)
	})
	script := writeTerminalTestScript(t, `echo "ready"; read line; echo "got:$line"`)

	info, err := mgr.Create(context.Background(), CreateParams{
		SessionID: 7,
		Command:   script,
		Title:     "test terminal",
		Visible:   true,
	})
	if err != nil {
		t.Fatalf("create terminal: %v", err)
	}
	if info.ID == "" {
		t.Fatal("expected terminal id")
	}

	if !eventually(2*time.Second, func() bool {
		chunks, err := stor.ReadTerminalOutput(info.ID, 0)
		return err == nil && strings.Contains(joinChunks(chunks), "ready")
	}) {
		t.Fatalf("expected ready output, events=%+v", emitted)
	}

	if err := mgr.Write(info.ID, "lemon\n"); err != nil {
		t.Fatalf("write terminal: %v", err)
	}

	if !eventually(2*time.Second, func() bool {
		chunks, err := stor.ReadTerminalOutput(info.ID, 0)
		return err == nil && strings.Contains(joinChunks(chunks), "got:lemon")
	}) {
		chunks, _ := stor.ReadTerminalOutput(info.ID, 0)
		t.Fatalf("expected persisted stdin response, got %q events=%+v", joinChunks(chunks), emitted)
	}
}

// eventually polls a condition until it passes or a timeout is reached.
func eventually(timeout time.Duration, fn func() bool) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if fn() {
			return true
		}
		time.Sleep(20 * time.Millisecond)
	}
	return fn()
}

// joinChunks concatenates terminal output chunks for assertions.
func joinChunks(chunks []data_models.TerminalOutputChunk) string {
	var b strings.Builder
	for _, chunk := range chunks {
		b.WriteString(chunk.Data)
	}
	return b.String()
}
