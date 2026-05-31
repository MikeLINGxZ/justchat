package provider

import (
	"context"
	"fmt"
	"sort"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_model"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/ierror"
	pkgProvider "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/provider"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/provider/provider_dto"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/storage"
)

// Provider manages provider persistence and view-model conversion.
type Provider struct {
	istorage *storage.Storage
}

// NewProvider creates a provider service bound to a storage handle.
func NewProvider(istorage *storage.Storage) *Provider {
	return &Provider{istorage: istorage}
}

// CreateProvider inserts a new provider record, bulk-saves its models, and sets the default model when specified.
func (p *Provider) CreateProvider(ctx context.Context, input provider_dto.CreateProviderInput) (*provider_dto.CreateProviderOutput, error) {
	created, err := p.istorage.CreateProvider(data_models.Provider{
		ProviderName: input.ProviderName,
		ProviderType: input.ProviderType,
		BaseUrl:      input.BaseUrl,
		ApiKey:       input.ApiKey,
		Enable:       input.Enable,
	})
	if err != nil {
		return nil, ierror.Error(ierror.ErrProviderCreate, err)
	}

	// Use provided model list; auto-fetch from API if none were supplied.
	modelsToSave := input.Models
	if len(modelsToSave) == 0 && input.BaseUrl != "" {
		if apiModels, fetchErr := pkgProvider.GetModels(input.BaseUrl, input.ApiKey); fetchErr == nil {
			for _, m := range apiModels {
				modelsToSave = append(modelsToSave, view_model.Model{
					Model:   m.ID,
					OwnedBy: m.OwnedBy,
					Object:  m.Object,
					Enable:  true,
				})
			}
		}
	}

	dbModels := make([]data_models.Model, 0, len(modelsToSave))
	for _, m := range modelsToSave {
		dbModels = append(dbModels, data_models.Model{
			ProviderId: created.ID,
			Model:      m.Model,
			OwnedBy:    m.OwnedBy,
			Object:     m.Object,
			Enable:     true,
			IsCustom:   m.IsCustom,
		})
	}
	if err := p.istorage.CreateModels(&dbModels); err != nil {
		return nil, ierror.Error(ierror.ErrProviderCreateModels, err)
	}

	// Persist the selected default model by name.
	if input.DefaultModel != nil {
		for _, m := range dbModels {
			if m.Model == *input.DefaultModel {
				_ = p.istorage.UpsertDefaultModel(created.ID, m.ID)
				break
			}
		}
	}

	return &provider_dto.CreateProviderOutput{}, nil
}

// ListProviders returns all provider view models with per-provider model lists and default-state annotations.
func (p *Provider) ListProviders(ctx context.Context, input provider_dto.ListProvidersInput) (*provider_dto.ListProvidersOutput, error) {
	providers, err := p.istorage.ListProviders()
	if err != nil {
		return nil, ierror.Error(ierror.ErrProviderListProviders, err)
	}

	defaultProviderID := readDefaultProviderID()

	wrappers := make([]provider_dto.ProviderWrapper, 0, len(providers))
	for _, prov := range providers {
		models, err := p.istorage.ListModelsForProvider(prov.ID)
		if err != nil {
			return nil, ierror.Error(ierror.ErrProviderListModels, err)
		}

		defaultModel, _ := p.istorage.GetDefaultModel(prov.ID)

		vmModels := make([]view_model.Model, 0, len(models))
		for _, m := range models {
			isDefault := defaultModel != nil && defaultModel.ModelId == m.ID
			vmModels = append(vmModels, toModelViewModel(m, isDefault))
		}

		wrappers = append(wrappers, provider_dto.ProviderWrapper{
			Provider: toProviderViewModel(prov, prov.ID == defaultProviderID),
			Models:   vmModels,
		})
	}

	return &provider_dto.ListProvidersOutput{Providers: wrappers}, nil
}

// DeleteProvider removes the provider, all its models, and the default-model record.
func (p *Provider) DeleteProvider(ctx context.Context, input provider_dto.DeleteProviderInput) (*provider_dto.DeleteProviderOutput, error) {
	if err := p.istorage.DeleteProviderCascade(uint(input.ProviderId)); err != nil {
		return nil, ierror.Error(ierror.ErrProviderDelete, err)
	}
	return &provider_dto.DeleteProviderOutput{}, nil
}

// EditProvider updates the provider's mutable fields and creates any new models (id == 0) in the input list.
func (p *Provider) EditProvider(ctx context.Context, input provider_dto.EditProviderInput) (*provider_dto.EditProviderOutput, error) {
	id := uint(input.ProviderId)

	if err := p.istorage.UpdateProvider(data_models.Provider{
		OrmModel:     data_models.OrmModel{ID: id},
		ProviderName: input.ProviderName,
		BaseUrl:      input.BaseUrl,
		ApiKey:       input.ApiKey,
		Enable:       input.Enable,
	}); err != nil {
		return nil, ierror.Error(ierror.ErrProviderUpdate, err)
	}

	// Collect existing model names to avoid duplicates.
	existing, err := p.istorage.ListModelsForProvider(id)
	if err != nil {
		return nil, ierror.Error(ierror.ErrProviderListModels, err)
	}
	existingNames := make(map[string]struct{}, len(existing))
	for _, m := range existing {
		existingNames[m.Model] = struct{}{}
	}

	var newModels []data_models.Model
	for _, m := range input.Models {
		if m.ID == 0 {
			if _, dup := existingNames[m.Model]; !dup {
				newModels = append(newModels, data_models.Model{
					ProviderId: id,
					Model:      m.Model,
					OwnedBy:    m.OwnedBy,
					Object:     m.Object,
					Enable:     true,
					IsCustom:   m.IsCustom,
				})
			}
		}
	}
	if len(newModels) > 0 {
		if err := p.istorage.CreateModels(&newModels); err != nil {
			return nil, ierror.Error(ierror.ErrProviderCreateModels, err)
		}
	}

	// Update default model when specified.
	if input.DefaultModel != nil {
		all, err := p.istorage.ListModelsForProvider(id)
		if err != nil {
			return nil, ierror.Error(ierror.ErrProviderListModels, err)
		}
		for _, m := range all {
			if m.Model == *input.DefaultModel {
				_ = p.istorage.UpsertDefaultModel(id, m.ID)
				break
			}
		}
	}

	return &provider_dto.EditProviderOutput{}, nil
}

// DeleteModel soft-deletes a single model record by its ID.
func (p *Provider) DeleteModel(ctx context.Context, input provider_dto.DeleteModelInput) (*provider_dto.DeleteModelOutput, error) {
	if input.ModelId <= 0 {
		return nil, ierror.Error(ierror.ErrProviderInvalidModel, fmt.Errorf("invalid model_id: %d", input.ModelId))
	}
	if err := p.istorage.DeleteModel(uint(input.ModelId)); err != nil {
		return nil, ierror.Error(ierror.ErrProviderDeleteModel, err)
	}
	return &provider_dto.DeleteModelOutput{}, nil
}

// SetDefault marks a provider as the global default and sets its default model.
// When ModelId is nil the first model (sorted by name) is used.
func (p *Provider) SetDefault(ctx context.Context, input provider_dto.SetDefaultInput) (*provider_dto.SetDefaultOutput, error) {
	id := uint(input.ProviderId)

	if err := writeDefaultProviderID(id); err != nil {
		return nil, ierror.Error(ierror.ErrProviderSetDefault, err)
	}

	if input.ModelId != nil {
		if *input.ModelId <= 0 {
			return nil, ierror.Error(ierror.ErrProviderInvalidModel, fmt.Errorf("invalid model_id: %d", *input.ModelId))
		}
		_ = p.istorage.UpsertDefaultModel(id, uint(*input.ModelId))
	} else {
		models, err := p.istorage.ListModelsForProvider(id)
		if err == nil && len(models) > 0 {
			sort.Slice(models, func(i, j int) bool { return models[i].Model < models[j].Model })
			_ = p.istorage.UpsertDefaultModel(id, models[0].ID)
		}
	}

	return &provider_dto.SetDefaultOutput{}, nil
}

// RequestProviderModelList fetches the live model list from the provider API.
func (p *Provider) RequestProviderModelList(ctx context.Context, input provider_dto.RequestProviderModelListInput) (*provider_dto.RequestProviderModelListOutput, error) {
	apiModels, err := pkgProvider.GetModels(input.BaseUrl, input.ApiKey)
	if err != nil {
		return nil, ierror.Error(ierror.ErrProviderFetchModels, err)
	}

	models := make([]view_model.Model, 0, len(apiModels))
	for _, m := range apiModels {
		models = append(models, view_model.Model{
			Model:   m.ID,
			OwnedBy: m.OwnedBy,
			Object:  m.Object,
			Enable:  true,
		})
	}

	return &provider_dto.RequestProviderModelListOutput{Models: models}, nil
}

// ProviderAndModelList returns all providers with their models, both sorted by name ascending.
func (p *Provider) ProviderAndModelList(ctx context.Context, input provider_dto.ProviderAndModelListInput) (*provider_dto.ProviderAndModelListOutput, error) {
	providers, err := p.istorage.ListProviders()
	if err != nil {
		return nil, ierror.Error(ierror.ErrProviderListProviders, err)
	}

	sort.Slice(providers, func(i, j int) bool {
		return providers[i].ProviderName < providers[j].ProviderName
	})

	defaultProviderID := readDefaultProviderID()

	result := make([]view_model.ProviderModel, 0, len(providers))
	for _, prov := range providers {
		models, err := p.istorage.ListModelsForProvider(prov.ID)
		if err != nil {
			return nil, ierror.Error(ierror.ErrProviderListModels, err)
		}

		sort.Slice(models, func(i, j int) bool {
			return models[i].Model < models[j].Model
		})

		defaultModel, _ := p.istorage.GetDefaultModel(prov.ID)

		vmModels := make([]view_model.Model, 0, len(models))
		for _, m := range models {
			vmModels = append(vmModels, toModelViewModel(m, defaultModel != nil && defaultModel.ModelId == m.ID))
		}

		result = append(result, view_model.ProviderModel{
			Provider: toProviderViewModel(prov, prov.ID == defaultProviderID),
			Models:   vmModels,
		})
	}

	return &provider_dto.ProviderAndModelListOutput{ProviderModels: result}, nil
}
