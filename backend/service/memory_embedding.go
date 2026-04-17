package service

import (
	"context"
	"fmt"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/agents/memory/search"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"
)

// ConfigureEmbedding 配置嵌入引擎。前端调用此方法传入嵌入配置后，
// 重建 HybridSearcher 以启用向量搜索。
func (s *Service) ConfigureEmbedding(provider string, baseURL string, apiKey string, model string) error {
	if s.memoryStorage == nil {
		return fmt.Errorf("memory system not initialized")
	}

	cfg := search.EmbeddingConfig{
		Provider: search.EmbeddingProvider(provider),
		BaseURL:  baseURL,
		APIKey:   apiKey,
		Model:    model,
	}

	ctx := context.Background()
	embedder, err := search.NewEmbedder(ctx, cfg)
	if err != nil {
		return fmt.Errorf("failed to create embedder: %w", err)
	}

	// 重建混合检索引擎
	s.memorySearcher = search.NewHybridSearcher(s.memoryStorage, embedder)

	logger.Warm("embedding engine configured:", provider, model)

	// 异步补填存量记忆的嵌入
	go func() {
		bgCtx := context.Background()
		count, backfillErr := s.memorySearcher.BackfillEmbeddings(bgCtx, model, 100)
		if backfillErr != nil {
			logger.Error("backfill embeddings error:", backfillErr)
		} else if count > 0 {
			logger.Warm("backfilled embeddings for", count, "memories")
		}
	}()

	return nil
}

// DisableEmbedding 禁用向量搜索，回退到纯 FTS5/LIKE。
func (s *Service) DisableEmbedding() {
	if s.memoryStorage == nil {
		return
	}
	s.memorySearcher = search.NewHybridSearcher(s.memoryStorage, nil)
	logger.Warm("embedding engine disabled, falling back to FTS5/LIKE")
}

// GetEmbeddingProviders 返回支持的嵌入引擎列表。
func (s *Service) GetEmbeddingProviders() []map[string]string {
	return []map[string]string{
		{
			"id":            string(search.EmbeddingProviderOllama),
			"name":          "Ollama",
			"description":   "本地运行的 Ollama 嵌入模型",
			"default_url":   "http://localhost:11434",
			"default_model": "bge-m3",
			"need_api_key":  "false",
		},
		{
			"id":            string(search.EmbeddingProviderOpenAICompat),
			"name":          "OpenAI 兼容",
			"description":   "支持 DeepSeek、通义千问、OpenRouter 等 OpenAI 兼容的嵌入 API",
			"default_url":   "",
			"default_model": "",
			"need_api_key":  "true",
		},
	}
}
