package service

import (
	"context"
	"testing"
	"time"

	"github.com/cloudwego/eino/components/embedding"
	memory_models "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/agents/memory/models"
	memory_search "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/agents/memory/search"
	memory_storage "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/agents/memory/storage"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
)

type fakeMemoryEmbedder struct{}

func (fakeMemoryEmbedder) EmbedStrings(ctx context.Context, texts []string, opts ...embedding.Option) ([][]float64, error) {
	vectors := make([][]float64, 0, len(texts))
	for range texts {
		vectors = append(vectors, []float64{0.5, 0.4, 0.3})
	}
	return vectors, nil
}

func TestUpdateMemoryPersistsEditsAndReembeds(t *testing.T) {
	svc, st := newTaskRecoveryTestService(t)

	memStorage, err := memory_storage.NewStorage(st.DB())
	if err != nil {
		t.Fatalf("memory_storage.NewStorage() error = %v", err)
	}
	if err := memStorage.AutoMigrateEmbeddings(); err != nil {
		t.Fatalf("AutoMigrateEmbeddings() error = %v", err)
	}
	svc.memoryStorage = memStorage
	svc.memorySearcher = memory_search.NewHybridSearcher(memStorage, fakeMemoryEmbedder{}, "memory-edit-model")

	id, err := memStorage.WriterMemory(context.Background(), memory_models.Memory{
		Summary: "旧标题",
		Content: "旧内容",
	})
	if err != nil {
		t.Fatalf("WriterMemory() error = %v", err)
	}

	updated, err := svc.UpdateMemory(id, view_models.MemoryUpdateInput{
		Summary: "新标题",
		Content: "新的记忆内容",
		Type:    "event",
	})
	if err != nil {
		t.Fatalf("UpdateMemory() error = %v", err)
	}
	if updated == nil {
		t.Fatal("UpdateMemory() = nil, want value")
	}
	if updated.Summary != "新标题" || updated.Content != "新的记忆内容" {
		t.Fatalf("updated memory = %+v, want edited content", updated)
	}

	deadline := time.Now().Add(2 * time.Second)
	for {
		memory, getErr := memStorage.GetMemoryByID(context.Background(), id)
		if getErr != nil {
			t.Fatalf("GetMemoryByID() error = %v", getErr)
		}
		if memory != nil && memory.EmbeddingID != nil {
			var row memory_storage.MemoryEmbedding
			if err := memStorage.DB().WithContext(context.Background()).
				Where("memory_id = ?", id).
				First(&row).Error; err != nil {
				t.Fatalf("load embedding error = %v", err)
			}
			if row.ModelName != "memory-edit-model" {
				t.Fatalf("embedding model = %q, want %q", row.ModelName, "memory-edit-model")
			}
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("timed out waiting for memory embedding refresh")
		}
		time.Sleep(20 * time.Millisecond)
	}
}
