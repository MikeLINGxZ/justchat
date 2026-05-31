package storage

import (
	"strings"
	"testing"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
)

func TestMemoryCRUDAndFilters(t *testing.T) {
	stor := newTestStorage(t)

	userMemory, err := stor.CreateMemory(data_models.Memory{
		Summary:    "Tea preference",
		Content:    "The user prefers lemon tea with less sugar.",
		Type:       "information",
		Target:     "user",
		Source:     "manual",
		Importance: 80,
		Confidence: 90,
		Pinned:     true,
	})
	if err != nil {
		t.Fatalf("CreateMemory user: %v", err)
	}
	_, err = stor.CreateMemory(data_models.Memory{
		Summary: "Project rule",
		Content: "For this project, keep Go database operations under backend/storage.",
		Type:    "fact",
		Target:  "memory",
		Source:  "agent",
	})
	if err != nil {
		t.Fatalf("CreateMemory assistant: %v", err)
	}

	listed, total, err := stor.ListMemories(MemoryListFilter{Target: "user", Limit: 20})
	if err != nil {
		t.Fatalf("ListMemories: %v", err)
	}
	if total != 1 || len(listed) != 1 || listed[0].ID != userMemory.ID {
		t.Fatalf("expected only user memory, total=%d listed=%+v", total, listed)
	}
	if listed[0].CharCount != len([]rune(listed[0].Content)) {
		t.Fatalf("expected char count to be maintained, got %d", listed[0].CharCount)
	}

	updated, err := stor.UpdateMemory(userMemory.ID, MemoryUpdate{
		Summary:    stringPtr("Tea preference updated"),
		Content:    stringPtr("The user prefers lemon tea with less sugar and no ice."),
		Importance: intPtr(95),
	})
	if err != nil {
		t.Fatalf("UpdateMemory: %v", err)
	}
	if updated.Summary != "Tea preference updated" || updated.Importance != 95 {
		t.Fatalf("unexpected updated memory: %+v", updated)
	}

	if err := stor.ForgetMemory(userMemory.ID); err != nil {
		t.Fatalf("ForgetMemory: %v", err)
	}
	listed, total, err = stor.ListMemories(MemoryListFilter{IncludeForgotten: false, Limit: 20})
	if err != nil {
		t.Fatalf("ListMemories after forget: %v", err)
	}
	if total != 1 {
		t.Fatalf("expected forgotten memory to be hidden, total=%d listed=%+v", total, listed)
	}
	listed, total, err = stor.ListMemories(MemoryListFilter{IncludeForgotten: true, Limit: 20})
	if err != nil {
		t.Fatalf("ListMemories include forgotten: %v", err)
	}
	if total != 2 {
		t.Fatalf("expected forgotten memory to be visible, total=%d listed=%+v", total, listed)
	}

	if err := stor.RestoreMemory(userMemory.ID); err != nil {
		t.Fatalf("RestoreMemory: %v", err)
	}
	got, err := stor.GetMemory(userMemory.ID)
	if err != nil {
		t.Fatalf("GetMemory: %v", err)
	}
	if got.IsForgotten {
		t.Fatalf("expected restored memory, got %+v", got)
	}
}

func TestMemorySearchStatsAndCoreRendering(t *testing.T) {
	stor := newTestStorage(t)

	_, _ = stor.CreateMemory(data_models.Memory{
		Summary:    "Low priority",
		Content:    "The user occasionally tries experimental UI ideas.",
		Type:       "information",
		Target:     "user",
		Source:     "agent",
		Importance: 10,
		Confidence: 80,
	})
	important, err := stor.CreateMemory(data_models.Memory{
		Summary:    "Coding style",
		Content:    "The user prefers conservative changes that follow existing project patterns.",
		Type:       "fact",
		Target:     "user",
		Source:     "manual",
		Importance: 90,
		Confidence: 95,
		Pinned:     true,
	})
	if err != nil {
		t.Fatalf("CreateMemory important: %v", err)
	}
	_, _ = stor.CreateMemory(data_models.Memory{
		Summary:    "Storage rule",
		Content:    "Database reads and writes should stay inside backend/storage.",
		Type:       "fact",
		Target:     "memory",
		Source:     "agent",
		Importance: 75,
		Confidence: 90,
	})

	results, err := stor.SearchMemories("database storage", 5)
	if err != nil {
		t.Fatalf("SearchMemories: %v", err)
	}
	if len(results) == 0 || !strings.Contains(results[0].Content, "backend/storage") {
		t.Fatalf("expected storage rule search result, got %+v", results)
	}
	if results[0].RecallCount != 1 || results[0].LastRecalledAt == nil {
		t.Fatalf("expected recall metadata to update, got %+v", results[0])
	}

	core, err := stor.RenderCoreMemory(1000, 1000)
	if err != nil {
		t.Fatalf("RenderCoreMemory: %v", err)
	}
	if !strings.Contains(core, "User memories") || !strings.Contains(core, important.Content) {
		t.Fatalf("expected rendered user memory, got %q", core)
	}
	if !strings.Contains(core, "Assistant memories") || !strings.Contains(core, "backend/storage") {
		t.Fatalf("expected rendered assistant memory, got %q", core)
	}

	stats, err := stor.MemoryStats()
	if err != nil {
		t.Fatalf("MemoryStats: %v", err)
	}
	if stats.Total != 3 || stats.Active != 3 || stats.ByTarget["user"] != 2 || stats.ByTarget["memory"] != 1 {
		t.Fatalf("unexpected stats: %+v", stats)
	}
}

func stringPtr(value string) *string {
	return &value
}

func intPtr(value int) *int {
	return &value
}
