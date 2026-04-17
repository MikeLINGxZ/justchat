package search

import (
	"context"
	"math"
	"sort"
	"sync"
	"time"

	"github.com/cloudwego/eino/components/embedding"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/agents/memory/models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/agents/memory/storage"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"
)

const minTrustThreshold = 0.3

// ScoredMemory 带综合打分的记忆。
type ScoredMemory struct {
	Memory     models.Memory
	FTSScore   float64
	VecScore   float64
	Recency    float64
	FinalScore float64
}

// HybridSearcher 混合检索引擎（FTS5 + 向量 + 元数据）。
type HybridSearcher struct {
	storage  *storage.Storage
	embedder embedding.Embedder // 可为 nil（向量搜索不可用时）

	// 嵌入缓存
	cacheMu    sync.RWMutex
	embCache   []storage.EmbeddingEntry
	cacheReady bool
}

// NewHybridSearcher 创建混合检索引擎。embedder 为 nil 时仅用 FTS5/LIKE。
func NewHybridSearcher(s *storage.Storage, embedder embedding.Embedder) *HybridSearcher {
	hs := &HybridSearcher{
		storage:  s,
		embedder: embedder,
	}
	// 异步加载嵌入缓存
	if embedder != nil {
		go hs.refreshCache()
	}
	return hs
}

// Search 执行混合检索，返回按综合得分排序的 Top-K 结果。
func (hs *HybridSearcher) Search(ctx context.Context, keywords []string, queryText string, topK int) []ScoredMemory {
	// 1. FTS5/LIKE 检索
	ftsResults, err := hs.storage.FTSSearch(ctx, keywords, topK*3)
	if err != nil {
		logger.Error("hybrid search FTS error:", err)
		return nil
	}

	// 构建 FTS 得分映射（rank-based）
	ftsScores := make(map[uint]float64)
	for i, m := range ftsResults {
		ftsScores[m.ID] = 1.0 / float64(i+1)
	}

	// 2. 向量检索（可选）
	vecScores := make(map[uint]float64)
	vectorEnabled := false
	if hs.embedder != nil && hs.isCacheReady() {
		vectorEnabled = true
		vecResults := hs.vectorSearch(ctx, queryText, topK*3)
		for _, vs := range vecResults {
			vecScores[vs.MemoryID] = vs.Score
		}
	}

	// 3. 合并去重
	allIDs := make(map[uint]bool)
	for id := range ftsScores {
		allIDs[id] = true
	}
	for id := range vecScores {
		allIDs[id] = true
	}

	// 加载记忆详情
	var scored []ScoredMemory
	now := time.Now()
	for id := range allIDs {
		m, loadErr := hs.storage.GetMemoryByID(ctx, id)
		if loadErr != nil || m == nil {
			continue
		}

		// 信任过滤
		if m.TrustScore < minTrustThreshold {
			continue
		}

		sm := ScoredMemory{
			Memory:   *m,
			FTSScore: ftsScores[id],
			VecScore: vecScores[id],
			Recency:  math.Exp(-0.05 * now.Sub(m.CreatedAt).Hours() / 24.0),
		}

		// 加权打分
		if vectorEnabled {
			sm.FinalScore = (0.30*sm.FTSScore +
				0.30*sm.VecScore +
				0.10*sm.Recency +
				0.15*m.Importance) * m.TrustScore
		} else {
			sm.FinalScore = (0.45*sm.FTSScore +
				0.15*sm.Recency +
				0.20*m.Importance) * m.TrustScore
		}

		scored = append(scored, sm)
	}

	// 排序
	sort.Slice(scored, func(i, j int) bool {
		return scored[i].FinalScore > scored[j].FinalScore
	})

	if len(scored) > topK {
		scored = scored[:topK]
	}
	return scored
}

// VectorEnabled 返回向量搜索是否可用。
func (hs *HybridSearcher) VectorEnabled() bool {
	return hs.embedder != nil && hs.isCacheReady()
}

// RefreshCache 刷新嵌入缓存。
func (hs *HybridSearcher) RefreshCache() {
	hs.refreshCache()
}

// ---- 内部方法 ----

type vecResult struct {
	MemoryID uint
	Score    float64
}

func (hs *HybridSearcher) vectorSearch(ctx context.Context, queryText string, topK int) []vecResult {
	if queryText == "" {
		return nil
	}

	// 向量化查询文本
	vectors, err := hs.embedder.EmbedStrings(ctx, []string{queryText})
	if err != nil || len(vectors) == 0 {
		return nil
	}

	// float64 → float32
	queryVec := make([]float32, len(vectors[0]))
	for i, v := range vectors[0] {
		queryVec[i] = float32(v)
	}

	// 读取缓存
	hs.cacheMu.RLock()
	cache := hs.embCache
	hs.cacheMu.RUnlock()

	// 计算相似度
	results := make([]vecResult, 0, len(cache))
	for _, entry := range cache {
		sim := storage.CosineSimilarity(queryVec, entry.Vector)
		if sim > 0.1 { // 过滤过低相似度
			results = append(results, vecResult{MemoryID: entry.MemoryID, Score: sim})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})
	if len(results) > topK {
		results = results[:topK]
	}
	return results
}

func (hs *HybridSearcher) refreshCache() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	entries, err := hs.storage.LoadAllEmbeddings(ctx)
	if err != nil {
		logger.Error("refresh embedding cache error:", err)
		return
	}

	hs.cacheMu.Lock()
	hs.embCache = entries
	hs.cacheReady = true
	hs.cacheMu.Unlock()
}

func (hs *HybridSearcher) isCacheReady() bool {
	hs.cacheMu.RLock()
	defer hs.cacheMu.RUnlock()
	return hs.cacheReady
}

// EmbedAndStore 为指定记忆生成嵌入并存储。
func (hs *HybridSearcher) EmbedAndStore(ctx context.Context, memoryID uint, text string, modelName string) error {
	if hs.embedder == nil {
		return nil
	}
	vectors, err := hs.embedder.EmbedStrings(ctx, []string{text})
	if err != nil {
		return err
	}
	if len(vectors) == 0 {
		return nil
	}

	vec32 := make([]float32, len(vectors[0]))
	for i, v := range vectors[0] {
		vec32[i] = float32(v)
	}

	return hs.storage.SaveEmbedding(ctx, memoryID, vec32, modelName)
}

// BackfillEmbeddings 为缺少嵌入的记忆批量生成向量。
func (hs *HybridSearcher) BackfillEmbeddings(ctx context.Context, modelName string, batchSize int) (int, error) {
	if hs.embedder == nil {
		return 0, nil
	}

	ids, err := hs.storage.GetMemoryIDsWithoutEmbedding(ctx, batchSize)
	if err != nil {
		return 0, err
	}

	count := 0
	for _, id := range ids {
		m, loadErr := hs.storage.GetMemoryByID(ctx, id)
		if loadErr != nil || m == nil {
			continue
		}
		text := m.Summary + " " + m.Content
		if embedErr := hs.EmbedAndStore(ctx, id, text, modelName); embedErr != nil {
			logger.Error("backfill embedding error for memory", id, ":", embedErr)
			continue
		}
		count++
	}

	if count > 0 {
		hs.refreshCache()
	}
	return count, nil
}
