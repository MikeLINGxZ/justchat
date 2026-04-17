package service

import (
	"testing"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
)

func TestGetAppPreferencesReturnsEmbeddingDefaults(t *testing.T) {
	svc, _ := newTaskRecoveryTestService(t)

	prefs, err := svc.GetAppPreferences()
	if err != nil {
		t.Fatalf("GetAppPreferences() error = %v", err)
	}
	if prefs == nil {
		t.Fatal("GetAppPreferences() = nil, want value")
	}
	if prefs.Language != data_models.AppLanguageZhCN {
		t.Fatalf("language = %q, want %q", prefs.Language, data_models.AppLanguageZhCN)
	}
	if prefs.Region != data_models.AppRegionAsia {
		t.Fatalf("region = %q, want %q", prefs.Region, data_models.AppRegionAsia)
	}
	if prefs.EmbeddingProvider != "ollama" {
		t.Fatalf("embedding_provider = %q, want ollama", prefs.EmbeddingProvider)
	}
	if prefs.EmbeddingBaseURL != "http://localhost:11434" {
		t.Fatalf("embedding_base_url = %q, want default ollama URL", prefs.EmbeddingBaseURL)
	}
	if prefs.EmbeddingModel != "bge-m3" {
		t.Fatalf("embedding_model = %q, want bge-m3", prefs.EmbeddingModel)
	}
}

func TestUpdateAppPreferencesPersistsMemoryAndEmbeddingSettings(t *testing.T) {
	svc, _ := newTaskRecoveryTestService(t)

	input := view_models.AppPreferences{
		Language:            data_models.AppLanguageEnUS,
		Region:              data_models.AppRegionNorthAmerica,
		MemorySystemEnabled: true,
		VectorSearchEnabled: true,
		EmbeddingProvider:   "openai_compat",
		EmbeddingBaseURL:    "https://example.com/v1",
		EmbeddingAPIKey:     "sk-test",
		EmbeddingModel:      "text-embedding-3-small",
	}

	updated, err := svc.UpdateAppPreferences(input)
	if err != nil {
		t.Fatalf("UpdateAppPreferences() error = %v", err)
	}
	if updated == nil {
		t.Fatal("UpdateAppPreferences() = nil, want value")
	}
	if !updated.MemorySystemEnabled || !updated.VectorSearchEnabled {
		t.Fatalf("updated flags = %+v, want both enabled", updated)
	}
	if updated.EmbeddingProvider != input.EmbeddingProvider ||
		updated.EmbeddingBaseURL != input.EmbeddingBaseURL ||
		updated.EmbeddingAPIKey != input.EmbeddingAPIKey ||
		updated.EmbeddingModel != input.EmbeddingModel {
		t.Fatalf("updated embedding prefs = %+v, want %+v", updated, input)
	}

	got, err := svc.GetAppPreferences()
	if err != nil {
		t.Fatalf("GetAppPreferences() error = %v", err)
	}
	if got == nil {
		t.Fatal("GetAppPreferences() = nil, want value")
	}
	if !got.MemorySystemEnabled || !got.VectorSearchEnabled {
		t.Fatalf("stored flags = %+v, want both enabled", got)
	}
	if got.EmbeddingProvider != input.EmbeddingProvider ||
		got.EmbeddingBaseURL != input.EmbeddingBaseURL ||
		got.EmbeddingAPIKey != input.EmbeddingAPIKey ||
		got.EmbeddingModel != input.EmbeddingModel {
		t.Fatalf("stored embedding prefs = %+v, want %+v", got, input)
	}
}
