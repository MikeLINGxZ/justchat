package memory

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/memory/memory_dto"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/storage"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestMemoryService(t *testing.T) *Memory {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:%s?mode=memory&cache=shared", t.Name())), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	stor, err := storage.NewStorageFromDB(db)
	if err != nil {
		t.Fatal(err)
	}
	return NewMemory(stor)
}

func TestMemoryServiceCreateListUpdateForgetRestore(t *testing.T) {
	svc := newTestMemoryService(t)

	created, err := svc.CreateMemory(context.Background(), memory_dto.CreateMemoryInput{
		Summary:    "Working style",
		Content:    "The user prefers implementation suggestions grounded in the current repository.",
		Type:       "information",
		Target:     "user",
		Source:     "manual",
		Importance: 85,
		Confidence: 90,
		Pinned:     true,
	})
	if err != nil {
		t.Fatalf("CreateMemory: %v", err)
	}

	listed, err := svc.ListMemories(context.Background(), memory_dto.ListMemoriesInput{Target: "user", Limit: 20})
	if err != nil {
		t.Fatalf("ListMemories: %v", err)
	}
	if listed.Total != 1 || len(listed.Items) != 1 || listed.Items[0].ID != created.Memory.ID {
		t.Fatalf("unexpected listed memories: %+v", listed)
	}

	updated, err := svc.UpdateMemory(context.Background(), memory_dto.UpdateMemoryInput{
		ID:         created.Memory.ID,
		Summary:    "Working style updated",
		Content:    "The user prefers repository-grounded implementation suggestions and concise tradeoffs.",
		Importance: 95,
		Confidence: 92,
		Pinned:     true,
	})
	if err != nil {
		t.Fatalf("UpdateMemory: %v", err)
	}
	if updated.Memory.Summary != "Working style updated" || updated.Memory.Importance != 95 {
		t.Fatalf("unexpected update result: %+v", updated.Memory)
	}

	if _, err := svc.ForgetMemory(context.Background(), memory_dto.ForgetMemoryInput{ID: created.Memory.ID}); err != nil {
		t.Fatalf("ForgetMemory: %v", err)
	}
	stats, err := svc.GetMemoryStats(context.Background(), memory_dto.GetMemoryStatsInput{})
	if err != nil {
		t.Fatalf("GetMemoryStats: %v", err)
	}
	if stats.Stats.Forgotten != 1 || stats.Stats.Active != 0 {
		t.Fatalf("unexpected stats after forget: %+v", stats.Stats)
	}

	if _, err := svc.RestoreMemory(context.Background(), memory_dto.RestoreMemoryInput{ID: created.Memory.ID}); err != nil {
		t.Fatalf("RestoreMemory: %v", err)
	}
	got, err := svc.GetMemory(context.Background(), memory_dto.GetMemoryInput{ID: created.Memory.ID})
	if err != nil {
		t.Fatalf("GetMemory: %v", err)
	}
	if got.Memory.IsForgotten {
		t.Fatalf("expected restored memory, got %+v", got.Memory)
	}
}

func TestMemoryServiceSettingsRoundTrip(t *testing.T) {
	t.Setenv("LEMONTEA_DATA_DIR", t.TempDir())
	svc := newTestMemoryService(t)

	if err := os.WriteFile(filepath.Join(os.Getenv("LEMONTEA_DATA_DIR"), "config.json"), []byte(`{"language":"zh-CN","memory":{"enabled":false}}`), 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := svc.SaveMemorySettings(context.Background(), memory_dto.SaveMemorySettingsInput{Enabled: true}); err != nil {
		t.Fatalf("SaveMemorySettings: %v", err)
	}
	got, err := svc.GetMemorySettings(context.Background(), memory_dto.GetMemorySettingsInput{})
	if err != nil {
		t.Fatalf("GetMemorySettings: %v", err)
	}
	if !got.Enabled {
		t.Fatalf("expected memory setting to be enabled")
	}
}
