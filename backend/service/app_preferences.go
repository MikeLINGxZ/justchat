package service

import (
	"context"
	"errors"

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
		Language: prefs.Language,
		Region:   prefs.Region,
	}
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
		SingletonID: data_models.AppPreferencesSingletonID,
		Language:    language,
		Region:      region,
	}
	if err := s.storage.SaveAppPreferences(context.Background(), prefs); err != nil {
		return nil, ierror.NewError(err)
	}
	i18n.SetCurrentLocale(string(language))
	return toAppPreferencesViewModel(&prefs), nil
}
