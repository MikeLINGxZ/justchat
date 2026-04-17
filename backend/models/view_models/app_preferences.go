package view_models

import "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"

type AppPreferences struct {
	Language            data_models.AppLanguage `json:"language"`
	Region              data_models.AppRegion   `json:"region"`
	MemorySystemEnabled bool                    `json:"memory_system_enabled"`
	VectorSearchEnabled bool                    `json:"vector_search_enabled"`
	EmbeddingProvider   string                  `json:"embedding_provider"`
	EmbeddingBaseURL    string                  `json:"embedding_base_url"`
	EmbeddingAPIKey     string                  `json:"embedding_api_key"`
	EmbeddingModel      string                  `json:"embedding_model"`
}
