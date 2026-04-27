package storage

import (
	"context"
	"strings"
	"testing"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/agents/memory/models"
	app_storage "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/storage"
)

func newTestMemoryStorage(t *testing.T) *Storage {
	t.Helper()
	t.Setenv("LEMONTEA_DATA_PATH", t.TempDir())
	st, err := app_storage.NewStorage()
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}
	memStorage, err := NewStorage(st.DB())
	if err != nil {
		t.Fatalf("NewStorage(memory) error = %v", err)
	}
	return memStorage
}

func TestRenderCoreMemorySnapshotUsesCapacityAndDelimiter(t *testing.T) {
	memStorage := newTestMemoryStorage(t)
	ctx := context.Background()
	if _, err := memStorage.AddCoreMemory(ctx, models.MemoryTargetUser, "User is allergic to peanuts.", "test"); err != nil {
		t.Fatalf("AddCoreMemory(user) error = %v", err)
	}
	if _, err := memStorage.AddCoreMemory(ctx, models.MemoryTargetUser, "User prefers concise Chinese replies.", "test"); err != nil {
		t.Fatalf("AddCoreMemory(user 2) error = %v", err)
	}
	if _, err := memStorage.AddCoreMemory(ctx, models.MemoryTargetAgent, "Project uses Go and Wails.", "test"); err != nil {
		t.Fatalf("AddCoreMemory(memory) error = %v", err)
	}

	snapshot, err := memStorage.RenderCoreMemorySnapshot(ctx)
	if err != nil {
		t.Fatalf("RenderCoreMemorySnapshot() error = %v", err)
	}
	for _, want := range []string{"USER PROFILE", "MEMORY (assistant notes)", "§", "/1375 chars", "/2200 chars"} {
		if !strings.Contains(snapshot, want) {
			t.Fatalf("snapshot missing %q:\n%s", want, snapshot)
		}
	}
}

func TestCoreMemoryReplaceRemoveAndCapacity(t *testing.T) {
	memStorage := newTestMemoryStorage(t)
	ctx := context.Background()
	if _, err := memStorage.AddCoreMemory(ctx, models.MemoryTargetUser, "User prefers tea.", "test"); err != nil {
		t.Fatalf("AddCoreMemory() error = %v", err)
	}
	replaced, err := memStorage.ReplaceCoreMemory(ctx, models.MemoryTargetUser, "prefers tea", "User prefers lemon tea.")
	if err != nil {
		t.Fatalf("ReplaceCoreMemory() error = %v", err)
	}
	if !strings.Contains(replaced.Message, "replaced") {
		t.Fatalf("replace message = %q", replaced.Message)
	}
	removed, err := memStorage.RemoveCoreMemory(ctx, models.MemoryTargetUser, "lemon tea")
	if err != nil {
		t.Fatalf("RemoveCoreMemory() error = %v", err)
	}
	if !strings.Contains(removed.Message, "removed") {
		t.Fatalf("remove message = %q", removed.Message)
	}

	tooLong := strings.Repeat("x", UserMemoryCharLimit+1)
	result, err := memStorage.AddCoreMemory(ctx, models.MemoryTargetUser, tooLong, "test")
	if err != nil {
		t.Fatalf("AddCoreMemory(long) error = %v", err)
	}
	if !strings.Contains(result.Message, "exceed") {
		t.Fatalf("long add message = %q, want exceed", result.Message)
	}
}

func TestCoreMemorySubstringMustBeUnique(t *testing.T) {
	memStorage := newTestMemoryStorage(t)
	ctx := context.Background()
	if _, err := memStorage.AddCoreMemory(ctx, models.MemoryTargetUser, "User likes tea in the morning.", "test"); err != nil {
		t.Fatal(err)
	}
	if _, err := memStorage.AddCoreMemory(ctx, models.MemoryTargetUser, "User likes tea after dinner.", "test"); err != nil {
		t.Fatal(err)
	}
	result, err := memStorage.ReplaceCoreMemory(ctx, models.MemoryTargetUser, "tea", "User likes coffee.")
	if err != nil {
		t.Fatalf("ReplaceCoreMemory() error = %v", err)
	}
	if !strings.Contains(result.Message, "matched 2 entries") {
		t.Fatalf("replace ambiguous message = %q", result.Message)
	}
}
