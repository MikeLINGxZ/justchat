package search

import (
	"context"
	"testing"

	"github.com/cloudwego/eino/components/embedding"
	memory_models "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/agents/memory/models"
	memory_storage "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/agents/memory/storage"
	app_storage "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/storage"
)

type fakeEmbedder struct{}

func (fakeEmbedder) EmbedStrings(ctx context.Context, texts []string, opts ...embedding.Option) ([][]float64, error) {
	vectors := make([][]float64, 0, len(texts))
	for range texts {
		vectors = append(vectors, []float64{0.1, 0.2, 0.3})
	}
	return vectors, nil
}

func TestBackfillEmbeddingsWorksBeforeCacheReadyAndUsesEmbeddingModel(t *testing.T) {
	t.Setenv("LEMONTEA_DATA_PATH", t.TempDir())

	st, err := app_storage.NewStorage()
	if err != nil {
		t.Fatalf("NewStorage() error = %v", err)
	}
	memStorage, err := memory_storage.NewStorage(st.DB())
	if err != nil {
		t.Fatalf("memory_storage.NewStorage() error = %v", err)
	}
	if err := memStorage.AutoMigrateEmbeddings(); err != nil {
		t.Fatalf("AutoMigrateEmbeddings() error = %v", err)
	}

	id, err := memStorage.WriterMemory(context.Background(), memory_models.Memory{
		Summary:    "去东京看樱花",
		Content:    "我想在春天安排一次东京赏樱旅行。",
		TrustScore: 0.8,
		Importance: 0.9,
	})
	if err != nil {
		t.Fatalf("WriterMemory() error = %v", err)
	}

	searcher := NewHybridSearcher(memStorage, fakeEmbedder{}, "test-embedding-model")
	if !searcher.HasEmbedder() {
		t.Fatal("HasEmbedder() = false, want true")
	}

	count, err := searcher.BackfillEmbeddings(context.Background(), 10)
	if err != nil {
		t.Fatalf("BackfillEmbeddings() error = %v", err)
	}
	if count != 1 {
		t.Fatalf("BackfillEmbeddings() count = %d, want 1", count)
	}

	var embedding memory_storage.MemoryEmbedding
	if err := memStorage.DB().WithContext(context.Background()).
		Where("memory_id = ?", id).
		First(&embedding).Error; err != nil {
		t.Fatalf("load embedding error = %v", err)
	}
	if embedding.ModelName != "test-embedding-model" {
		t.Fatalf("embedding model = %q, want %q", embedding.ModelName, "test-embedding-model")
	}

	memory, err := memStorage.GetMemoryByID(context.Background(), id)
	if err != nil {
		t.Fatalf("GetMemoryByID() error = %v", err)
	}
	if memory == nil || memory.EmbeddingID == nil {
		t.Fatalf("memory embedding_id = %#v, want non-nil", memory)
	}
}
