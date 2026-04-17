package service

import (
	"context"
	"errors"
	"strings"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/i18n"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils/ierror"
)

var validRegions = map[data_models.AppRegion]struct{}{
	data_models.AppRegionAsia:         {},
	data_models.AppRegionEurope:       {},
	data_models.AppRegionNorthAmerica: {},
	data_models.AppRegionSouthAmerica: {},
	data_models.AppRegionAfrica:       {},
	data_models.AppRegionOceania:      {},
	data_models.AppRegionAntarctica:   {},
}

func normalizeLanguage(language data_models.AppLanguage) (data_models.AppLanguage, error) {
	switch language {
	case data_models.AppLanguageZhCN, "":
		return data_models.AppLanguageZhCN, nil
	case data_models.AppLanguageEnUS:
		return data_models.AppLanguageEnUS, nil
	default:
		return "", errors.New(i18n.TCurrent("prefs.language.invalid", map[string]string{"value": string(language)}))
	}
}

func normalizeRegion(region data_models.AppRegion) (data_models.AppRegion, error) {
	if region == "" {
		return data_models.AppRegionAsia, nil
	}
	if _, ok := validRegions[region]; ok {
		return region, nil
	}
	return "", errors.New(i18n.TCurrent("prefs.region.invalid", map[string]string{"value": string(region)}))
}

func toAppPreferencesViewModel(prefs *data_models.AppPreferences) *view_models.AppPreferences {
	return &view_models.AppPreferences{
		Language:            prefs.Language,
		Region:              prefs.Region,
		MemorySystemEnabled: prefs.MemorySystemEnabled,
		VectorSearchEnabled: prefs.VectorSearchEnabled,
		EmbeddingProvider:   prefs.EmbeddingProvider,
		EmbeddingBaseURL:    prefs.EmbeddingBaseURL,
		EmbeddingAPIKey:     prefs.EmbeddingAPIKey,
		EmbeddingModel:      prefs.EmbeddingModel,
	}
}

func normalizeEmbeddingProvider(provider string) string {
	provider = strings.TrimSpace(provider)
	if provider == "" {
		return "ollama"
	}
	return provider
}

func normalizeEmbeddingBaseURL(baseURL string, provider string) string {
	baseURL = strings.TrimSpace(baseURL)
	if baseURL == "" && provider == "ollama" {
		return "http://localhost:11434"
	}
	return baseURL
}

func normalizeEmbeddingModel(model string, provider string) string {
	model = strings.TrimSpace(model)
	if model == "" && provider == "ollama" {
		return "bge-m3"
	}
	return model
}

func (s *Service) loadAppPreferences(ctx context.Context) (*data_models.AppPreferences, error) {
	prefs, err := s.storage.GetAppPreferences(ctx)
	if err != nil {
		return nil, err
	}
	language, err := normalizeLanguage(prefs.Language)
	if err != nil {
		return nil, err
	}
	region, err := normalizeRegion(prefs.Region)
	if err != nil {
		return nil, err
	}
	prefs.EmbeddingProvider = normalizeEmbeddingProvider(prefs.EmbeddingProvider)
	prefs.EmbeddingBaseURL = normalizeEmbeddingBaseURL(prefs.EmbeddingBaseURL, prefs.EmbeddingProvider)
	prefs.EmbeddingModel = normalizeEmbeddingModel(prefs.EmbeddingModel, prefs.EmbeddingProvider)
	prefs.Language = language
	prefs.Region = region
	return prefs, nil
}

func (s *Service) GetAppPreferences() (*view_models.AppPreferences, error) {
	prefs, err := s.loadAppPreferences(context.Background())
	if err != nil {
		return nil, ierror.NewError(err)
	}
	return toAppPreferencesViewModel(prefs), nil
}

func (s *Service) UpdateAppPreferences(input view_models.AppPreferences) (*view_models.AppPreferences, error) {
	language, err := normalizeLanguage(input.Language)
	if err != nil {
		return nil, ierror.NewError(err)
	}
	region, err := normalizeRegion(input.Region)
	if err != nil {
		return nil, ierror.NewError(err)
	}

	prefs := data_models.AppPreferences{
		SingletonID:         data_models.AppPreferencesSingletonID,
		Language:            language,
		Region:              region,
		MemorySystemEnabled: input.MemorySystemEnabled,
		VectorSearchEnabled: input.VectorSearchEnabled,
		EmbeddingProvider:   normalizeEmbeddingProvider(input.EmbeddingProvider),
		EmbeddingBaseURL:    normalizeEmbeddingBaseURL(input.EmbeddingBaseURL, normalizeEmbeddingProvider(input.EmbeddingProvider)),
		EmbeddingAPIKey:     strings.TrimSpace(input.EmbeddingAPIKey),
		EmbeddingModel:      normalizeEmbeddingModel(input.EmbeddingModel, normalizeEmbeddingProvider(input.EmbeddingProvider)),
	}
	if err := s.storage.SaveAppPreferences(context.Background(), prefs); err != nil {
		return nil, ierror.NewError(err)
	}
	i18n.SetCurrentLocale(string(language))
	return toAppPreferencesViewModel(&prefs), nil
}
