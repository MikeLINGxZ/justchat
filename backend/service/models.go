package service

import (
	"context"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils/ierror"
)

// GetModels 获取所有模型
func (s *Service) GetModels(enableProvider, enableModel *bool) ([]view_models.Model, error) {
	ctx := context.Background()

	providers, err := s.storage.GetProviders(ctx)
	if err != nil {
		return nil, ierror.NewError(err)
	}

	// 获取目标供应商
	var providerIds []uint
	for _, provider := range providers {
		if enableProvider != nil && *enableProvider {
			if provider.Enable {
				providerIds = append(providerIds, provider.ID)
			}
			continue
		}
		if enableProvider == nil && !*enableModel {
			if !provider.Enable {
				providerIds = append(providerIds, provider.ID)
			}
			continue
		}
		providerIds = append(providerIds, provider.ID)
	}

	providerId2Models, err := s.storage.GetModelsByProviderIds(ctx, providerIds)
	if err != nil {
		return nil, ierror.NewError(err)
	}

	var res []view_models.Model
	for _, models := range providerId2Models {
		for _, model := range models {
			res = append(res, view_models.Model{
				ID:         model.ID,
				CreatedAt:  model.CreatedAt,
				UpdatedAt:  model.UpdatedAt,
				DeletedAt:  model.DeletedAt,
				ProviderId: model.ProviderId,
				Model:      model.Model,
				OwnedBy:    model.OwnedBy,
				Object:     model.Object,
				Enable:     model.Enable,
				Alias:      model.Alias,
				IsCustom:   model.IsCustom,
			})
		}
	}

	return res, nil
}
