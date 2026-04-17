package search

import (
	"context"
	"fmt"

	ollama_embedding "github.com/cloudwego/eino-ext/components/embedding/ollama"
	openai_embedding "github.com/cloudwego/eino-ext/libs/acl/openai"
	"github.com/cloudwego/eino/components/embedding"
)

// EmbeddingProvider 嵌入引擎类型。
type EmbeddingProvider string

const (
	EmbeddingProviderOllama       EmbeddingProvider = "ollama"
	EmbeddingProviderOpenAICompat EmbeddingProvider = "openai_compat" // 通用 OpenAI 兼容接口
)

// EmbeddingConfig 嵌入引擎配置。
type EmbeddingConfig struct {
	Provider EmbeddingProvider `json:"provider"`
	BaseURL  string            `json:"base_url"` // API 地址，如 http://localhost:11434 或 https://api.deepseek.com/v1
	APIKey   string            `json:"api_key"`  // API 密钥（Ollama 可为空）
	Model    string            `json:"model"`    // 模型名，如 bge-m3、text-embedding-3-small
}

// NewEmbedder 根据配置创建嵌入引擎实例。
func NewEmbedder(ctx context.Context, cfg EmbeddingConfig) (embedding.Embedder, error) {
	switch cfg.Provider {
	case EmbeddingProviderOllama:
		return newOllamaEmbedder(ctx, cfg)
	case EmbeddingProviderOpenAICompat:
		return newOpenAICompatEmbedder(ctx, cfg)
	default:
		return nil, fmt.Errorf("unsupported embedding provider: %s", cfg.Provider)
	}
}

func newOllamaEmbedder(ctx context.Context, cfg EmbeddingConfig) (embedding.Embedder, error) {
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	model := cfg.Model
	if model == "" {
		model = "bge-m3"
	}
	return ollama_embedding.NewEmbedder(ctx, &ollama_embedding.EmbeddingConfig{
		BaseURL: baseURL,
		Model:   model,
	})
}

func newOpenAICompatEmbedder(ctx context.Context, cfg EmbeddingConfig) (embedding.Embedder, error) {
	if cfg.Model == "" {
		return nil, fmt.Errorf("model is required for openai_compat embedding provider")
	}
	return openai_embedding.NewEmbeddingClient(ctx, &openai_embedding.EmbeddingConfig{
		BaseURL: cfg.BaseURL,
		APIKey:  cfg.APIKey,
		Model:   cfg.Model,
	})
}
