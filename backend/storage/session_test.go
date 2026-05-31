package storage

import (
	"fmt"
	"testing"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// newTestStorage creates an isolated in-memory storage for storage-layer tests.
func newTestStorage(t *testing.T) *Storage {
	t.Helper()

	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	stor, err := NewStorageFromDB(db)
	if err != nil {
		t.Fatalf("new storage: %v", err)
	}
	return stor
}

// TestSessionKindFieldPersists verifies the session kind round-trips through SQLite.
func TestSessionKindFieldPersists(t *testing.T) {
	stor := newTestStorage(t)

	created, err := stor.CreateSession(data_models.Session{
		Title: "x",
		Kind:  "background",
	})
	if err != nil {
		t.Fatal(err)
	}

	got, err := stor.GetSession(created.ID)
	if err != nil {
		t.Fatal(err)
	}

	if got.Kind != "background" {
		t.Fatalf("kind not persisted: %q", got.Kind)
	}
}

// TestListSessionsHidesBackgroundByDefault verifies background sessions stay hidden unless explicitly requested.
func TestListSessionsHidesBackgroundByDefault(t *testing.T) {
	stor := newTestStorage(t)

	_, err := stor.CreateSession(data_models.Session{Title: "user-1"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = stor.CreateSession(data_models.Session{Title: "bg-1", Kind: "background"})
	if err != nil {
		t.Fatal(err)
	}
	_, err = stor.CreateSession(data_models.Session{Title: "task-1", Kind: "task"})
	if err != nil {
		t.Fatal(err)
	}

	listed, err := stor.ListSessions(0, 50, false, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(listed) != 2 {
		t.Fatalf("expected user and task sessions, got %+v", listed)
	}
	for _, session := range listed {
		if session.Kind == "background" {
			t.Fatalf("background session should stay hidden by default: %+v", listed)
		}
	}

	all, err := stor.ListSessions(0, 50, false, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(all) != 3 {
		t.Fatalf("expected 3 with includeHidden, got %d", len(all))
	}
}

func TestTaskStateRoundTripsPerSessionAndKey(t *testing.T) {
	stor := newTestStorage(t)

	session, err := stor.CreateSession(data_models.Session{Title: "task", Kind: "task"})
	if err != nil {
		t.Fatal(err)
	}

	if err := stor.SaveTaskState(session.ID, "device_code", "abc123"); err != nil {
		t.Fatal(err)
	}
	value, found, err := stor.LoadTaskState(session.ID, "device_code")
	if err != nil {
		t.Fatal(err)
	}
	if !found || value != "abc123" {
		t.Fatalf("unexpected task state value found=%v value=%q", found, value)
	}

	if err := stor.SaveTaskState(session.ID, "device_code", "xyz789"); err != nil {
		t.Fatal(err)
	}
	value, found, err = stor.LoadTaskState(session.ID, "device_code")
	if err != nil {
		t.Fatal(err)
	}
	if !found || value != "xyz789" {
		t.Fatalf("unexpected updated task state value found=%v value=%q", found, value)
	}

	if err := stor.DeleteTaskState(session.ID, "device_code"); err != nil {
		t.Fatal(err)
	}
	value, found, err = stor.LoadTaskState(session.ID, "device_code")
	if err != nil {
		t.Fatal(err)
	}
	if found || value != "" {
		t.Fatalf("expected deleted task state, found=%v value=%q", found, value)
	}
}
