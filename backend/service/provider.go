package service

import (
	"context"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/storage"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils/ierror"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils/llm"
)

// GetProviders 获取所有供应商
func (s *Service) GetProviders() ([]view_models.Provider, error) {
	providers, err := s.storage.GetProviders(context.Background())
	if err != nil {
		return nil, ierror.NewError(err)
	}
	res := make([]view_models.Provider, len(providers))
	for i, provider := range providers {
		res[i] = view_models.Provider{
			ID:           provider.ID,
			Alias:        provider.Alias,
			ApiKey:       provider.ApiKey,
			BaseUrl:      provider.BaseUrl,
			Enable:       provider.Enable,
			ProviderName: provider.ProviderName,
		}
	}
	return res, nil
}

// AddProvider 添加供应商
func (s *Service) AddProvider(provider view_models.Provider) error {
	providerId, err := s.storage.AddProvider(context.Background(), data_models.Provider{
		ProviderName: provider.ProviderName,
		BaseUrl:      provider.BaseUrl,
		ApiKey:       provider.ApiKey,
		Enable:       provider.Enable,
		Alias:        provider.Alias,
	})
	if err != nil {
		return ierror.NewError(err)
	}

	// 更新模型信息
	err = s.updateProviderModel(providerId)
	if err != nil {
		// todo
	}

	return nil
}

// UpdateProvider 更新供应商
func (s *Service) UpdateProvider(id uint, provider *view_models.Provider) error {
	// provider为空，代表更新模型信息
	if provider == nil {
		return nil
	}
	return s.storage.UpdateProvider(context.Background(), data_models.Provider{
		OrmModel: data_models.OrmModel{
			ID: id,
		},
		ProviderName: provider.ProviderName,
		BaseUrl:      provider.BaseUrl,
		ApiKey:       provider.ApiKey,
		Enable:       provider.Enable,
		Alias:        provider.Alias,
	})
}

// DeleteProvider 删除供应商
func (s *Service) DeleteProvider(providerId uint) error {
	return s.storage.DeleteProvider(context.Background(), providerId)
}

// GetProviderModels 获取供应商模型信息
func (s *Service) GetProviderModels(provider view_models.Provider) ([]view_models.Model, error) {
	providerModels, err := llm.GetModels(provider.BaseUrl, provider.ApiKey)
	if err != nil {
		return nil, ierror.NewError(err)
	}
	res := make([]view_models.Model, len(providerModels))
	for _, item := range providerModels {
		res = append(res, view_models.Model{
			Model:   item.ID,
			OwnedBy: item.OwnedBy,
			Object:  item.Object,
		})
	}
	return res, nil
}

// updateProviderModel 更新供应商模型列表
func (s *Service) updateProviderModel(providerId uint) error {
	ctx := context.Background()

	// 获取供应商信息
	provider, err := s.storage.GetProviderByID(ctx, providerId)
	if err != nil {
		return err
	}
	if provider == nil {
		return nil
	}

	// 获取模型信息
	providerModels, err := llm.GetModels(provider.BaseUrl, provider.ApiKey)
	if err != nil {
		return err
	}

	var newModels []data_models.Model
	for _, item := range providerModels {
		newModels = append(newModels, data_models.Model{
			ProviderId: providerId,
			Model:      item.ID,
			OwnedBy:    item.OwnedBy,
			Object:     item.Object,
		})
	}

	// ...
	err = s.storage.NewFnTransaction(ctx, func(ctx context.Context, s *storage.Storage) error {
		// 删除供应商下所有模型
		err := s.DeleteAllProviderModel(ctx, providerId)
		if err != nil {
			return err
		}
		// 插入新模型
		for _, newModel := range newModels {
			err = s.AddProviderModel(ctx, newModel)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
