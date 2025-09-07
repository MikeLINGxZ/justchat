package service

import (
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils/ierror"
)

// GetProviders 获取所有供应商
func (s *Service) GetProviders() ([]view_models.Provider, error) {
	providers, err := s.storage.GetProviders(s.ctx)
	if err != nil {
		return nil, ierror.NewError(err)
	}
	res := make([]view_models.Provider, len(providers))
	for i, provider := range providers {
		res[i] = view_models.Provider{
			Alias:        provider.Alias,
			ApiKey:       provider.ApiKey,
			BaseUrl:      provider.BaseUrl,
			Enable:       provider.Enable,
			ProviderName: provider.ProviderName,
		}
	}
	return res, nil
}
