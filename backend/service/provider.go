package service

import (
	"context"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/llm_provider"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/storage"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/utils/ierror"
)

// GetProviders 获取所有供应商
func (s *Service) GetProviders(ctx context.Context) ([]view_models.Provider, error) {
	providers, err := s.storage.GetProviders(ctx)
	if err != nil {
		return nil, ierror.NewError(err)
	}

	// 获取全部供应商id
	var providerIds []uint
	for _, provider := range providers {
		providerIds = append(providerIds, provider.ID)
	}

	// 获取各个供应商模型列表
	providerIds2Models, err := s.storage.GetModelsByProviderIds(ctx, providerIds)
	if err != nil {
		return nil, ierror.NewError(err)
	}
	providerIds2ModelsVD := make(map[uint][]view_models.Model)
	for providerId, models := range providerIds2Models {
		var modelsVD []view_models.Model
		for _, model := range models {
			modelsVD = append(modelsVD, view_models.Model{
				ID:       model.ID,
				Enable:   model.Enable,
				Alias:    model.Alias,
				Model:    model.Model,
				OwnedBy:  model.OwnedBy,
				Object:   model.Object,
				IsCustom: model.IsCustom,
			})
		}
		providerIds2ModelsVD[providerId] = modelsVD
	}

	// 获取各个供应商默认模型
	providerId2DefaultModel, err := s.storage.GetProviderDefaultModelByProviderIds(ctx, providerIds)
	if err != nil {
		return nil, ierror.NewError(err)
	}
	providerIds2DefaultModelIdVD := make(map[uint]*uint)
	for providerId, model := range providerId2DefaultModel {
		providerIds2DefaultModelIdVD[providerId] = &model.ID
	}

	res := make([]view_models.Provider, len(providers))
	for i, provider := range providers {
		res[i] = view_models.Provider{
			ID:                provider.ID,
			ApiKey:            provider.ApiKey,
			BaseUrl:           provider.BaseUrl,
			FileUploadBaseUrl: provider.FileUploadBaseUrl,
			Enable:            provider.Enable,
			IsDefault:         provider.IsDefault,
			ProviderName:      provider.ProviderName,
			ProviderType:      provider.ProviderType,
			Models:            providerIds2ModelsVD[provider.ID],
			DefaultModelId:    providerIds2DefaultModelIdVD[provider.ID],
		}
	}
	return res, nil
}

// AddProvider 添加供应商
func (s *Service) AddProvider(ctx context.Context, provider view_models.Provider) error {
	providerId, err := s.storage.AddProvider(context.Background(), data_models.Provider{
		ProviderName:      provider.ProviderName,
		ProviderType:      provider.ProviderType,
		BaseUrl:           provider.BaseUrl,
		FileUploadBaseUrl: provider.FileUploadBaseUrl,
		ApiKey:            provider.ApiKey,
		Enable:            provider.Enable,
		IsDefault:         provider.IsDefault,
	})
	if err != nil {
		return ierror.NewError(err)
	}

	// 更新模型信息
	err = s.updateProviderModel(ctx, providerId)
	if err != nil {
		// todo
	}

	return nil
}

// UpdateProvider 更新供应商
func (s *Service) UpdateProvider(ctx context.Context, id uint, provider *view_models.Provider) error {
	// provider为空，代表更新模型信息
	if provider == nil {
		return nil
	}

	provider.ID = id
	err := s.storage.NewFnTransaction(ctx, func(ctx context.Context, tx *storage.Storage) error {
		existingProvider, err := tx.GetProviderByID(ctx, id)
		if err != nil {
			return err
		}
		if existingProvider == nil {
			return ierror.New(ierror.ErrCodeProviderNotFound)
		}

		err = tx.UpdateProvider(ctx, data_models.Provider{
			OrmModel: data_models.OrmModel{
				ID: id,
			},
			ProviderName:      provider.ProviderName,
			ProviderType:      provider.ProviderType,
			BaseUrl:           provider.BaseUrl,
			FileUploadBaseUrl: provider.FileUploadBaseUrl,
			ApiKey:            provider.ApiKey,
			Enable:            provider.Enable,
			IsDefault:         existingProvider.IsDefault,
		})
		if err != nil {
			return err
		}

		if provider.DefaultModelId != nil {
			err = tx.UpdateProviderDefaultModel(ctx, provider.ID, *provider.DefaultModelId)
			if err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return ierror.NewError(err)
	}

	return nil
}

// SetDefaultProvider 设置默认供应商，并返回用于聊天默认选择的模型
func (s *Service) SetDefaultProvider(ctx context.Context, providerId uint) (*view_models.Model, error) {
	var defaultModel *data_models.Model
	err := s.storage.NewFnTransaction(ctx, func(ctx context.Context, tx *storage.Storage) error {
		provider, err := tx.GetProviderByID(ctx, providerId)
		if err != nil {
			return err
		}
		if provider == nil {
			return ierror.New(ierror.ErrCodeProviderNotFound)
		}

		if err := tx.ClearDefaultProviders(ctx); err != nil {
			return err
		}
		if err := tx.SetProviderDefault(ctx, providerId); err != nil {
			return err
		}

		model, err := tx.GetProviderDefaultModel(ctx, providerId)
		if err != nil {
			return err
		}
		if model != nil {
			defaultModel = model
			return nil
		}

		models, err := tx.GetModelsByProviderIdSorted(ctx, providerId)
		if err != nil {
			return err
		}
		if len(models) == 0 {
			return nil
		}

		firstModel := models[0]
		if err := tx.UpdateProviderDefaultModel(ctx, providerId, firstModel.ID); err != nil {
			return err
		}
		defaultModel = &firstModel
		return nil
	})
	if err != nil {
		return nil, ierror.NewError(err)
	}

	if defaultModel == nil {
		return nil, nil
	}
	return &view_models.Model{
		ID:         defaultModel.ID,
		CreatedAt:  defaultModel.CreatedAt,
		UpdatedAt:  defaultModel.UpdatedAt,
		DeletedAt:  defaultModel.DeletedAt,
		ProviderId: defaultModel.ProviderId,
		Model:      defaultModel.Model,
		OwnedBy:    defaultModel.OwnedBy,
		Object:     defaultModel.Object,
		Enable:     defaultModel.Enable,
		Alias:      defaultModel.Alias,
		IsCustom:   defaultModel.IsCustom,
	}, nil
}

// UpdateProviderModels 更新供应商
func (s *Service) UpdateProviderModels(ctx context.Context, providerId uint) error {
	err := s.updateProviderModel(ctx, providerId)
	if err != nil {
		return ierror.NewError(err)
	}
	return nil
}

// DeleteProvider 删除供应商
func (s *Service) DeleteProvider(ctx context.Context, providerId uint) error {
	return s.storage.DeleteProvider(ctx, providerId)
}

// GetProviderModels 获取供应商模型信息
func (s *Service) GetProviderModels(ctx context.Context, provider view_models.Provider) ([]view_models.Model, error) {
	providerModels, err := llm_provider.GetModels(provider.BaseUrl, provider.ApiKey)
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

// AddProviderCustomModel 添加自定义模型
func (s *Service) AddProviderCustomModel(ctx context.Context, providerId uint, modelName string) error {
	providers, err := s.storage.GetProviders(ctx)
	if err != nil {
		return err
	}
	var provider *data_models.Provider
	for _, item := range providers {
		if item.ID == providerId {
			provider = &item
			break
		}
	}
	if provider == nil {
		return ierror.New(ierror.ErrCodeProviderNotFound)
	}

	return s.storage.AddProviderModel(ctx, data_models.Model{
		ProviderId: provider.ID,
		Model:      modelName,
		OwnedBy:    provider.ProviderName,
		Object:     "",
		Enable:     true,
		Alias:      nil,
		IsCustom:   true,
	})
}

func (s *Service) DeleteProviderCustomModel(ctx context.Context, providerId uint, modelName string) error {
	providers, err := s.storage.GetProviders(ctx)
	if err != nil {
		return err
	}
	var provider *data_models.Provider
	for _, item := range providers {
		if item.ID == providerId {
			provider = &item
			break
		}
	}
	if provider == nil {
		return ierror.New(ierror.ErrCodeProviderNotFound)
	}

	err = s.storage.DeleteProviderModel(ctx, providerId, modelName)
	if err != nil {
		return ierror.NewError(err)
	}

	return nil
}

// updateProviderModel 更新供应商模型列表
func (s *Service) updateProviderModel(ctx context.Context, providerId uint) error {
	// 获取供应商信息
	provider, err := s.storage.GetProviderByID(ctx, providerId)
	if err != nil {
		return err
	}
	if provider == nil {
		return nil
	}

	previousDefaultModel, err := s.storage.GetProviderDefaultModel(ctx, providerId)
	if err != nil {
		return err
	}
	previousDefaultModelName := ""
	hadDefaultModel := false
	if previousDefaultModel != nil {
		previousDefaultModelName = previousDefaultModel.Model
		hadDefaultModel = true
	}

	// 获取模型信息
	providerModels, err := llm_provider.GetModels(provider.BaseUrl, provider.ApiKey)
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
			IsCustom:   false,
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

		if !hadDefaultModel {
			return nil
		}

		models, err := s.GetModelsByProviderIdSorted(ctx, providerId)
		if err != nil {
			return err
		}
		if len(models) == 0 {
			return nil
		}

		targetModelID := models[0].ID
		for _, model := range models {
			if model.Model == previousDefaultModelName {
				targetModelID = model.ID
				break
			}
		}

		return s.UpdateProviderDefaultModel(ctx, providerId, targetModelID)
	})
	if err != nil {
		return err
	}

	return nil
}
