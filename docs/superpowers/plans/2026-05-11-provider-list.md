# Provider List Enhancement Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 实现供应商列表增强功能，包括启用/禁用开关、更多菜单（设为默认/删除）、模型搜索、应用按钮脏状态、API 密钥显示切换、模型数量标签。

**Architecture:** 后端实现 EditProvider/DeleteProvider/SetDefault 接口；前端 ProviderListItem 重构为 div（含 toggle + 更多菜单）、ProviderDetailView 增加脏状态跟踪与搜索功能、AddProviderStepForm 同步添加模型搜索。

**Tech Stack:** Go (GORM, Wails v3 bindings), React, TypeScript, Tailwind CSS, Zustand, i18next

---

## File Map

| 操作 | 文件 |
|------|------|
| 修改 | `backend/storage/provider.go` |
| 修改 | `backend/service/provider/provider_dto/edit_provider.go` |
| 修改 | `backend/service/provider/provider_internal.go` |
| 修改 | `backend/service/provider/provider.go` |
| 修改 | `backend/service/provider/provider_test.go` |
| 重新生成 | `frontend/bindings/.../provider_dto/models.ts` (自动) |
| 修改 | `frontend/src/i18n/locales/zh-CN.ts` |
| 修改 | `frontend/src/i18n/locales/en.ts` |
| 修改 | `frontend/src/store/settingsStore.ts` |
| 修改 | `frontend/src/components/settings/providers/ProviderListItem.tsx` |
| 修改 | `frontend/src/components/settings/providers/ProviderList.tsx` |
| 修改 | `frontend/src/components/settings/SettingsApp.tsx` |
| 修改 | `frontend/src/components/settings/providers/ProviderDetailView.tsx` |
| 修改 | `frontend/src/components/settings/providers/AddProviderStepForm.tsx` |

---

## Task 1: Backend Storage — 新增 UpdateProvider / DeleteProvider / DeleteModelsForProvider / DeleteDefaultModelRecord

**Files:**
- Modify: `backend/storage/provider.go`
- Test: `backend/service/provider/provider_test.go`

- [ ] **Step 1: 写失败测试**

在 `provider_test.go` 末尾追加以下两个测试：

```go
// TestDeleteProviderRemovesRecord verifies that DeleteProvider removes the record from ListProviders.
func TestDeleteProviderRemovesRecord(t *testing.T) {
	svc := newTestProviderService(t)
	ctx := context.Background()

	_, err := svc.CreateProvider(ctx, provider_dto.CreateProviderInput{
		ProviderName: "ToDelete",
		ProviderType: pkgProvider.Deepseek,
		Enable:       true,
	})
	if err != nil {
		t.Fatal(err)
	}

	out, err := svc.ListProviders(ctx, provider_dto.ListProvidersInput{})
	if err != nil {
		t.Fatal(err)
	}
	id := out.Providers[0].Provider.ID

	_, err = svc.DeleteProvider(ctx, provider_dto.DeleteProviderInput{ProviderId: int64(id)})
	if err != nil {
		t.Fatalf("DeleteProvider: %v", err)
	}

	out, _ = svc.ListProviders(ctx, provider_dto.ListProvidersInput{})
	if len(out.Providers) != 0 {
		t.Fatalf("expected 0 providers after delete, got %d", len(out.Providers))
	}
}

// TestEditProviderUpdatesFields verifies that EditProvider persists updated provider fields.
func TestEditProviderUpdatesFields(t *testing.T) {
	svc := newTestProviderService(t)
	ctx := context.Background()

	_, err := svc.CreateProvider(ctx, provider_dto.CreateProviderInput{
		ProviderName: "Original",
		ProviderType: pkgProvider.Deepseek,
		BaseUrl:      "http://original.test",
		ApiKey:       "sk-old",
		Enable:       true,
	})
	if err != nil {
		t.Fatal(err)
	}

	out, _ := svc.ListProviders(ctx, provider_dto.ListProvidersInput{})
	id := out.Providers[0].Provider.ID

	_, err = svc.EditProvider(ctx, provider_dto.EditProviderInput{
		ProviderId:   int64(id),
		ProviderName: "Updated",
		BaseUrl:      "http://updated.test",
		ApiKey:       "sk-new",
		Enable:       false,
	})
	if err != nil {
		t.Fatalf("EditProvider: %v", err)
	}

	out, _ = svc.ListProviders(ctx, provider_dto.ListProvidersInput{})
	got := out.Providers[0].Provider
	if got.Name != "Updated" {
		t.Errorf("expected name Updated, got %q", got.Name)
	}
	if got.BaseURL != "http://updated.test" {
		t.Errorf("expected base_url http://updated.test, got %q", got.BaseURL)
	}
	if got.ApiKey != "sk-new" {
		t.Errorf("expected api_key sk-new, got %q", got.ApiKey)
	}
	if got.Enabled {
		t.Error("expected enabled=false, got true")
	}
}
```

- [ ] **Step 2: 运行测试，确认失败**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop
go test ./backend/service/provider/... -run "TestDeleteProvider|TestEditProvider" -v
```

期望：FAIL — "implement me" panic 或编译错误

- [ ] **Step 3: 在 `backend/storage/provider.go` 末尾追加 4 个方法**

```go
// UpdateProvider updates the name, base URL, API key, and enable flag for the provider identified by p.ID.
func (s *Storage) UpdateProvider(p data_models.Provider) error {
	return s.sqliteDB.Model(&p).Updates(map[string]interface{}{
		"provider_name": p.ProviderName,
		"base_url":      p.BaseUrl,
		"api_key":       p.ApiKey,
		"enable":        p.Enable,
	}).Error
}

// DeleteProvider soft-deletes the provider record identified by providerID.
func (s *Storage) DeleteProvider(providerID uint) error {
	return s.sqliteDB.Delete(&data_models.Provider{}, providerID).Error
}

// DeleteModelsForProvider soft-deletes all model records belonging to providerID.
func (s *Storage) DeleteModelsForProvider(providerID uint) error {
	return s.sqliteDB.Where("provider_id = ?", providerID).Delete(&data_models.Model{}).Error
}

// DeleteDefaultModelRecord soft-deletes the default-model record for providerID.
func (s *Storage) DeleteDefaultModelRecord(providerID uint) error {
	return s.sqliteDB.Where("provider_id = ?", providerID).Delete(&data_models.ProviderDefaultModel{}).Error
}
```

- [ ] **Step 4: 运行测试，确认仍然失败（storage 通过，但 service 仍 panic）**

```bash
go test ./backend/service/provider/... -run "TestDeleteProvider|TestEditProvider" -v
```

期望：FAIL，service 层 panic "implement me"

- [ ] **Step 5: commit**

```bash
git add backend/storage/provider.go backend/service/provider/provider_test.go
git commit -m "feat(storage): add UpdateProvider, DeleteProvider, DeleteModels, DeleteDefaultModelRecord"
```

---

## Task 2: Backend DTO + Service 实现 + 重新生成 Bindings

**Files:**
- Modify: `backend/service/provider/provider_dto/edit_provider.go`
- Modify: `backend/service/provider/provider_internal.go`
- Modify: `backend/service/provider/provider.go`
- Auto-regenerate: `frontend/bindings/.../provider_dto/models.ts`

- [ ] **Step 1: 更新 `edit_provider.go`**

替换文件全部内容：

```go
package provider_dto

import "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_model"

// EditProviderInput carries all mutable fields for an existing provider.
type EditProviderInput struct {
	ProviderId   int64              `json:"provider_id"`   // 供应商ID
	ProviderName string             `json:"provider_name"` // 供应商名称
	BaseUrl      string             `json:"base_url"`      // 供应商基础URL
	ApiKey       string             `json:"api_key"`       // 供应商API密钥
	Enable       bool               `json:"enable"`        // 是否启用
	DefaultModel *string            `json:"default_model"` // 默认模型名称，nil 表示不修改
	Models       []view_model.Model `json:"models"`        // 需要新增的模型（id=0 的条目）
}

type EditProviderOutput struct{}
```

- [ ] **Step 2: 在 `provider_internal.go` 末尾追加 `writeDefaultProviderID`**

```go
// writeDefaultProviderID persists DefaultProviderID in the config file; preserves all other config fields.
func writeDefaultProviderID(providerID uint) error {
	dataDir, err := dir.GetDataDir()
	if err != nil {
		return err
	}
	configPath := filepath.Join(dataDir, dir.ConfigFileName)

	var cfg data_models.Config
	if raw, readErr := os.ReadFile(configPath); readErr == nil {
		_ = json.Unmarshal(raw, &cfg)
	}
	cfg.DefaultProviderID = providerID

	content, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, content, 0o644)
}
```

注意 `provider_internal.go` 已经 import 了 `"encoding/json"`, `"os"`, `"path/filepath"`, `"gitlab.linhf.cn/.../pkg/dir"` 和 `data_models` — 无需重复导入。

- [ ] **Step 3: 实现 `provider.go` 中三个 panic 方法**

将 `DeleteProvider` 替换为：

```go
// DeleteProvider removes the provider, all its models, and the default-model record.
func (p *Provider) DeleteProvider(ctx context.Context, input provider_dto.DeleteProviderInput) (*provider_dto.DeleteProviderOutput, error) {
	id := uint(input.ProviderId)
	if err := p.istorage.DeleteModelsForProvider(id); err != nil {
		return nil, fmt.Errorf("delete models for provider: %w", err)
	}
	if err := p.istorage.DeleteDefaultModelRecord(id); err != nil {
		return nil, fmt.Errorf("delete default model record: %w", err)
	}
	if err := p.istorage.DeleteProvider(id); err != nil {
		return nil, fmt.Errorf("delete provider: %w", err)
	}
	return &provider_dto.DeleteProviderOutput{}, nil
}
```

将 `EditProvider` 替换为：

```go
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
		return nil, fmt.Errorf("update provider: %w", err)
	}

	// Collect existing model names to avoid duplicates.
	existing, err := p.istorage.ListModelsForProvider(id)
	if err != nil {
		return nil, fmt.Errorf("list models: %w", err)
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
			return nil, fmt.Errorf("create new models: %w", err)
		}
	}

	// Update default model when specified.
	if input.DefaultModel != nil {
		all, _ := p.istorage.ListModelsForProvider(id)
		for _, m := range all {
			if m.Model == *input.DefaultModel {
				_ = p.istorage.UpsertDefaultModel(id, m.ID)
				break
			}
		}
	}

	return &provider_dto.EditProviderOutput{}, nil
}
```

将 `SetDefault` 替换为：

```go
// SetDefault marks a provider as the global default and sets its default model.
// When ModelId is nil the first model (sorted by name) is used.
func (p *Provider) SetDefault(ctx context.Context, input provider_dto.SetDefaultInput) (*provider_dto.SetDefaultOutput, error) {
	id := uint(input.ProviderId)

	if err := writeDefaultProviderID(id); err != nil {
		return nil, fmt.Errorf("write default provider id: %w", err)
	}

	if input.ModelId != nil {
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
```

在 `provider.go` 顶部 import 中添加 `"sort"`：

```go
import (
	"context"
	"fmt"
	"sort"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/data_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_model"
	pkgProvider "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/provider"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/provider/provider_dto"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/storage"
)
```

- [ ] **Step 4: 运行测试，确认通过**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop
go test ./backend/service/provider/... -v
```

期望：所有测试 PASS（SetDefault 测试会写临时 config，需设置 `LEMONTEA_DATA_DIR`；其他测试全部通过）

若 TestDeleteProvider/TestEditProvider FAIL，检查 storage 方法实现。

- [ ] **Step 5: 重新生成 TypeScript Bindings**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop
wails3 generate bindings -clean=true -ts
```

期望：`frontend/bindings/.../provider_dto/models.ts` 中的 `EditProviderInput` 类新增 `provider_id`, `provider_name`, `base_url`, `api_key`, `enable`, `default_model`, `models` 字段。

- [ ] **Step 6: 验证 bindings 已更新**

```bash
grep -A5 "export class EditProviderInput" frontend/bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/provider/provider_dto/models.ts
```

期望输出含 `"provider_id": number` 等字段。

- [ ] **Step 7: commit**

```bash
git add backend/service/provider/provider_dto/edit_provider.go \
        backend/service/provider/provider_internal.go \
        backend/service/provider/provider.go \
        backend/service/provider/provider_test.go \
        frontend/bindings/
git commit -m "feat(provider): implement EditProvider, DeleteProvider, SetDefault"
```

---

## Task 3: Frontend i18n — 新增翻译键

**Files:**
- Modify: `frontend/src/i18n/locales/zh-CN.ts`
- Modify: `frontend/src/i18n/locales/en.ts`

- [ ] **Step 1: 在 `zh-CN.ts` 的 `providers:` 部分添加新键**

找到：
```typescript
    providers: {
      add: '添加',
    },
```

替换为：
```typescript
    providers: {
      add: '添加',
      setDefault: '设为默认',
      delete: '删除',
      defaultTag: '默认',
      modelsCount: '{{count}} 个模型',
      searchModelPlaceholder: '搜索模型...',
      noSearchResults: '未找到匹配的模型',
      more: '更多',
      showApiKey: '显示密钥',
      hideApiKey: '隐藏密钥',
    },
```

- [ ] **Step 2: 在 `en.ts` 的 `providers:` 部分添加新键**

找到：
```typescript
    providers: {
      add: 'Add',
    },
```

替换为：
```typescript
    providers: {
      add: 'Add',
      setDefault: 'Set as Default',
      delete: 'Delete',
      defaultTag: 'Default',
      modelsCount: '{{count}} models',
      searchModelPlaceholder: 'Search models...',
      noSearchResults: 'No matching models',
      more: 'More',
      showApiKey: 'Show API Key',
      hideApiKey: 'Hide API Key',
    },
```

- [ ] **Step 3: 验证类型检查通过**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop/frontend
npx tsc --noEmit 2>&1 | head -20
```

期望：无报错（或仅有与本次改动无关的既有警告）

- [ ] **Step 4: commit**

```bash
git add frontend/src/i18n/locales/zh-CN.ts frontend/src/i18n/locales/en.ts
git commit -m "feat(i18n): add provider list enhancement translation keys"
```

---

## Task 4: Frontend Store — 新增 updateProvider / deleteProvider / setDefaultProvider

**Files:**
- Modify: `frontend/src/store/settingsStore.ts`

- [ ] **Step 1: 在 `SettingsState` type 中添加三个新 action 签名**

找到：
```typescript
  setProviders: (providers: SettingsBootstrap['providers']) => void
```

替换为：
```typescript
  setProviders: (providers: SettingsBootstrap['providers']) => void
  updateProvider: (provider: ProviderItem) => void
  deleteProvider: (id: number) => void
  setDefaultProvider: (providerId: number) => void
```

需要在文件顶部添加 `ProviderItem` 导入（如果尚未导入）：

找到：
```typescript
import type { GeneralSettingsTab, SettingsBootstrap, SettingsOption, SettingsPrimaryTab } from '@/types/settings'
```

替换为：
```typescript
import type { GeneralSettingsTab, ProviderItem, SettingsBootstrap, SettingsOption, SettingsPrimaryTab } from '@/types/settings'
```

- [ ] **Step 2: 在 `createSettingsState` 实现部分添加三个 action 实现**

找到：
```typescript
    setProviders: (providers) => set({ providers }),
```

替换为：
```typescript
    setProviders: (providers) => set({ providers }),
    updateProvider: (provider) => set((state) => ({
      providers: state.providers.map(p => p.id === provider.id ? provider : p),
    })),
    deleteProvider: (id) => set((state) => {
      const providers = state.providers.filter(p => p.id !== id)
      return {
        providers,
        selectedProviderId: state.selectedProviderId === id
          ? (providers[0]?.id ?? null)
          : state.selectedProviderId,
      }
    }),
    setDefaultProvider: (providerId) => set((state) => ({
      providers: state.providers.map(p => ({ ...p, is_default: p.id === providerId })),
    })),
```

- [ ] **Step 3: 验证类型检查**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop/frontend
npx tsc --noEmit 2>&1 | head -20
```

- [ ] **Step 4: commit**

```bash
git add frontend/src/store/settingsStore.ts
git commit -m "feat(store): add updateProvider, deleteProvider, setDefaultProvider actions"
```

---

## Task 5: Frontend ProviderListItem — 重构为 div，添加 toggle + 更多菜单

**Files:**
- Modify: `frontend/src/components/settings/providers/ProviderListItem.tsx`

- [ ] **Step 1: 完整替换 `ProviderListItem.tsx`**

```tsx
import { useState } from 'react'
import { MoreHorizontal, Star } from 'lucide-react'
import { useTranslation } from 'react-i18next'
import { cn } from '@/lib/utils'
import type { ProviderItem } from '@/types/settings'

export function ProviderListItem(props: {
  item: ProviderItem
  selected: boolean
  onSelect: (id: number) => void
  onToggleEnable: (item: ProviderItem) => void
  onSetDefault: (id: number) => void
  onDelete: (id: number) => void
}) {
  const { t } = useTranslation()
  const [menuOpen, setMenuOpen] = useState(false)
  const { item } = props

  return (
    <div
      className={cn(
        'group relative flex w-full items-center gap-2 rounded-2xl px-3 py-3 transition-colors',
        props.selected ? 'bg-primary/10 text-primary' : 'bg-background hover:bg-accent'
      )}
    >
      {/* 主点击区域：选中供应商 */}
      <div
        role="button"
        tabIndex={0}
        onClick={() => props.onSelect(item.id)}
        onKeyDown={(e) => { if (e.key === 'Enter' || e.key === ' ') props.onSelect(item.id) }}
        className="flex min-w-0 flex-1 cursor-pointer items-start gap-3"
      >
        <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-2xl">
          {item.icon ? (
            <img
              src={item.icon}
              alt={item.provider_name}
              className="h-10 w-10 rounded-2xl object-contain"
              onError={(e) => { (e.currentTarget as HTMLImageElement).style.display = 'none' }}
            />
          ) : (
            <span className="text-lg font-semibold uppercase text-primary">
              {item.provider_name.slice(0, 2)}
            </span>
          )}
        </div>
        <div className="min-w-0 flex-1">
          <div className="flex items-center gap-2">
            <span className="truncate text-sm font-medium">{item.provider_name}</span>
            {item.is_default && (
              <span className="rounded-full bg-primary/10 px-2 py-0.5 text-xs text-primary">
                {t('settingsPage.providers.defaultTag')}
              </span>
            )}
          </div>
          <p className="mt-1 truncate text-xs text-muted-foreground">{item.base_url}</p>
        </div>
      </div>

      {/* 右侧操作区 */}
      <div className="flex shrink-0 items-center gap-1">
        {/* 启用/禁用开关 */}
        <button
          type="button"
          role="switch"
          aria-checked={item.enabled}
          aria-label={t('settingsPage.providers.more')}
          onClick={(e) => { e.stopPropagation(); props.onToggleEnable(item) }}
          className={cn(
            'relative inline-flex h-5 w-9 shrink-0 items-center rounded-full transition-colors',
            item.enabled ? 'bg-primary' : 'bg-muted'
          )}
        >
          <span className={cn(
            'inline-block h-3.5 w-3.5 rounded-full bg-white shadow transition-transform',
            item.enabled ? 'translate-x-[18px]' : 'translate-x-[3px]'
          )} />
        </button>

        {/* 更多按钮 */}
        <button
          type="button"
          aria-label={t('settingsPage.providers.more')}
          onClick={(e) => { e.stopPropagation(); setMenuOpen(v => !v) }}
          className={cn(
            'flex h-7 w-7 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:bg-accent hover:text-foreground',
            'opacity-0 group-hover:opacity-100',
            menuOpen && 'opacity-100'
          )}
        >
          <MoreHorizontal size={14} />
        </button>
      </div>

      {/* 下拉菜单 */}
      {menuOpen && (
        <>
          <div className="fixed inset-0 z-10" onClick={() => setMenuOpen(false)} />
          <div className="absolute right-0 top-full z-20 mt-1 w-36 rounded-lg border border-border bg-popover py-1 shadow-md">
            <button
              type="button"
              className="flex w-full items-center gap-2 px-3 py-1.5 text-left text-sm hover:bg-accent"
              onClick={(e) => { e.stopPropagation(); props.onSetDefault(item.id); setMenuOpen(false) }}
            >
              <Star size={13} />
              {t('settingsPage.providers.setDefault')}
            </button>
            <button
              type="button"
              className="flex w-full items-center gap-2 px-3 py-1.5 text-left text-sm text-destructive hover:bg-accent"
              onClick={(e) => { e.stopPropagation(); props.onDelete(item.id); setMenuOpen(false) }}
            >
              {t('settingsPage.providers.delete')}
            </button>
          </div>
        </>
      )}
    </div>
  )
}
```

- [ ] **Step 2: 验证类型检查**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop/frontend
npx tsc --noEmit 2>&1 | head -30
```

期望：`ProviderList.tsx` 会报 props 不匹配错误（因为 ProviderList 还未传递新 props），这是预期的，将在 Task 6 修复。

- [ ] **Step 3: commit**

```bash
git add frontend/src/components/settings/providers/ProviderListItem.tsx
git commit -m "feat(provider-list-item): add enable toggle and more menu with set-default/delete actions"
```

---

## Task 6: Frontend ProviderList + SettingsApp — 传递回调、加宽侧栏、添加 refreshProviders

**Files:**
- Modify: `frontend/src/components/settings/providers/ProviderList.tsx`
- Modify: `frontend/src/components/settings/SettingsApp.tsx`

- [ ] **Step 1: 更新 `ProviderList.tsx`**

完整替换文件内容：

```tsx
import { PanelLeftOpen, Plus } from 'lucide-react'
import { useContext } from 'react'
import { useTranslation } from 'react-i18next'
import { SettingsMenuContext } from '@/components/settings/SettingsShell'
import { ProviderListItem } from '@/components/settings/providers/ProviderListItem'
import type { ProviderItem } from '@/types/settings'

export function ProviderList(props: {
  items: ProviderItem[]
  selectedId: number | null
  onSelect: (id: number) => void
  onCreate: () => void
  onToggleEnable: (item: ProviderItem) => void
  onSetDefault: (id: number) => void
  onDelete: (id: number) => void
}) {
  const { t } = useTranslation()
  const { onOpenMenu } = useContext(SettingsMenuContext)

  return (
    <div className="flex h-full flex-col">
      <div className="mb-4 flex items-center justify-between gap-3 px-2">
        <h2 className="text-lg font-semibold text-foreground">{t('settingsPage.primary.providers')}</h2>
        <div className="flex items-center gap-1">
          {onOpenMenu && (
            <button
              type="button"
              aria-label="Open settings menu"
              onClick={onOpenMenu}
              className="inline-flex h-8 items-center justify-center rounded-lg px-2 text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
            >
              <PanelLeftOpen size={15} />
            </button>
          )}
          <button
            type="button"
            onClick={props.onCreate}
            className="inline-flex h-8 items-center gap-1.5 rounded-lg px-2 text-xs text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
            aria-label={t('settingsPage.providers.add')}
          >
            <Plus size={14} />
            <span>{t('settingsPage.providers.add')}</span>
          </button>
        </div>
      </div>

      <div className="space-y-2">
        {props.items.map((item) => (
          <ProviderListItem
            key={item.id}
            item={item}
            selected={item.id === props.selectedId}
            onSelect={props.onSelect}
            onToggleEnable={props.onToggleEnable}
            onSetDefault={props.onSetDefault}
            onDelete={props.onDelete}
          />
        ))}
      </div>
    </div>
  )
}
```

- [ ] **Step 2: 更新 `SettingsApp.tsx`**

完整替换文件内容：

```tsx
import { useCallback, useEffect } from 'react'
import type { ReactNode } from 'react'
import { Globe2, Languages, MonitorSmartphone } from 'lucide-react'
import { AboutSettingsView } from '@/components/settings/about/AboutSettingsView'
import { SettingsShell } from '@/components/settings/SettingsShell'
import { SettingsSectionLayout } from '@/components/settings/common/SettingsSectionLayout'
import { SettingsSubmenuList } from '@/components/settings/common/SettingsSubmenuList'
import { GeneralSettingsPanel } from '@/components/settings/general/GeneralSettingsPanel'
import { ProviderDetailView } from '@/components/settings/providers/ProviderDetailView'
import { ProviderList } from '@/components/settings/providers/ProviderList'
import { useSettingsBootstrap } from '@/hooks/useSettingsBootstrap'
import { useSettingsStore } from '@/store/settingsStore'
import { useTranslation } from 'react-i18next'
import { Window } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/window'
import { Provider as ProviderBinding } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/provider'
import { EditProviderInput } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/provider/provider_dto/models'
import type { ProviderWrapper } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/provider/provider_dto/models'
import type { ProviderItem } from '@/types/settings'

function mapProviderItem(w: ProviderWrapper): ProviderItem {
  return {
    id: w.providers.id,
    provider_name: w.providers.provider_name,
    provider_type: w.providers.provider_type,
    base_url: w.providers.base_url,
    api_key: w.providers.api_key,
    enabled: w.providers.enabled,
    is_default: w.providers.is_default,
    model_count: w.models.length,
    icon: w.providers.icon,
    models: w.models.map((m) => ({
      id: m.id,
      provider_id: m.provider_id,
      model: m.model,
      owned_by: m.owned_by,
      object: m.object,
      enable: m.enable,
      alias: m.alias,
      is_custom: m.is_custom,
      is_default: m.is_default,
    })),
  }
}

export function SettingsApp() {
  useSettingsBootstrap()
  const { t } = useTranslation()

  const setProviders = useSettingsStore((state) => state.setProviders)
  const updateProvider = useSettingsStore((state) => state.updateProvider)
  const deleteProvider = useSettingsStore((state) => state.deleteProvider)
  const setDefaultProvider = useSettingsStore((state) => state.setDefaultProvider)

  const refreshProviders = useCallback(async () => {
    try {
      const result = await ProviderBinding.ListProviders({})
      if (result?.providers) {
        setProviders(result.providers.map(mapProviderItem))
      }
    } catch {
      // ignore
    }
  }, [setProviders])

  useEffect(() => {
    window.addEventListener('focus', refreshProviders)
    return () => window.removeEventListener('focus', refreshProviders)
  }, [refreshProviders])

  const handleToggleEnable = useCallback(async (item: ProviderItem) => {
    const updated = { ...item, enabled: !item.enabled }
    updateProvider(updated) // optimistic
    try {
      await ProviderBinding.EditProvider(new EditProviderInput({
        provider_id: item.id,
        provider_name: item.provider_name,
        base_url: item.base_url,
        api_key: item.api_key,
        enable: !item.enabled,
        default_model: null,
        models: [],
      }))
    } catch {
      updateProvider(item) // rollback
    }
  }, [updateProvider])

  const handleSetDefault = useCallback(async (id: number) => {
    try {
      await ProviderBinding.SetDefault({ provider_id: id, model_id: null })
      setDefaultProvider(id)
    } catch {
      // ignore
    }
  }, [setDefaultProvider])

  const handleDeleteFromList = useCallback(async (id: number) => {
    try {
      await ProviderBinding.DeleteProvider({ provider_id: id })
      deleteProvider(id)
    } catch {
      // ignore
    }
  }, [deleteProvider])

  const activeTab = useSettingsStore((state) => state.activeTab)
  const generalTab = useSettingsStore((state) => state.generalTab)
  const setGeneralTab = useSettingsStore((state) => state.setGeneralTab)
  const version = useSettingsStore((state) => state.version)
  const providers = useSettingsStore((state) => state.providers)
  const selectedProviderId = useSettingsStore((state) => state.selectedProviderId)
  const setSelectedProviderId = useSettingsStore((state) => state.setSelectedProviderId)

  const generalItems = [
    { key: 'display' as const, label: t('settingsPage.general.display.title'), icon: <MonitorSmartphone size={16} /> },
    { key: 'locale' as const, label: t('settingsPage.general.locale.title'), icon: <Languages size={16} /> },
    { key: 'file' as const, label: t('settingsPage.general.file.title'), icon: <Globe2 size={16} className="rotate-45" /> },
  ]

  let content: ReactNode

  if (activeTab === 'about') {
    content = (
      <SettingsSectionLayout>
        <AboutSettingsView version={version} />
      </SettingsSectionLayout>
    )
  } else if (activeTab === 'providers') {
    const selectedProvider = providers.find((provider) => provider.id === selectedProviderId) ?? providers[0] ?? null

    content = (
      <SettingsSectionLayout
        sidebarClassName="w-72"
        sidebar={(
          <ProviderList
            items={providers}
            selectedId={selectedProvider?.id ?? null}
            onSelect={setSelectedProviderId}
            onCreate={() => { void Window.OpenAddProvider({}) }}
            onToggleEnable={handleToggleEnable}
            onSetDefault={handleSetDefault}
            onDelete={handleDeleteFromList}
          />
        )}
      >
        <ProviderDetailView
          provider={selectedProvider}
          onUpdated={refreshProviders}
          onDeleted={deleteProvider}
        />
      </SettingsSectionLayout>
    )
  } else {
    content = (
      <SettingsSectionLayout
        sidebar={
          <SettingsSubmenuList
            title="通用"
            items={generalItems}
            value={generalTab}
            onChange={setGeneralTab}
          />
        }
      >
        <GeneralSettingsPanel />
      </SettingsSectionLayout>
    )
  }

  return <SettingsShell>{content}</SettingsShell>
}
```

- [ ] **Step 3: 验证类型检查**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop/frontend
npx tsc --noEmit 2>&1 | head -30
```

期望：`ProviderDetailView.tsx` 会报 props 不匹配错误（onUpdated/onDeleted 还未添加），将在 Task 7 修复。

- [ ] **Step 4: commit**

```bash
git add frontend/src/components/settings/providers/ProviderList.tsx \
        frontend/src/components/settings/SettingsApp.tsx
git commit -m "feat(provider-list): wire toggle/default/delete callbacks, widen sidebar to w-72"
```

---

## Task 7: Frontend ProviderDetailView — 全面增强

**Files:**
- Modify: `frontend/src/components/settings/providers/ProviderDetailView.tsx`

功能清单：
1. 新增 `onUpdated` / `onDeleted` props
2. key 包含 `provider.enabled` + `version`（list toggle 时自动重置）
3. 脏状态跟踪（apply 按钮默认禁用）
4. API 密钥显示/隐藏切换
5. 模型列表标题右侧显示"N 个模型"标签
6. 模型搜索（从右到左展开）
7. 模型"设为默认"调用 SetDefault 即时保存
8. Apply 调用 EditProvider，成功后触发 onUpdated
9. Delete 调用 DeleteProvider，成功后触发 onDeleted

- [ ] **Step 1: 完整替换 `ProviderDetailView.tsx`**

```tsx
import { useRef, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { RefreshCw, Plus, Eye, EyeOff, Search, X } from 'lucide-react'
import { Provider as ProviderBinding } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/provider'
import { EditProviderInput } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/service/provider/provider_dto/models'
import { Model as ViewModel } from '@bindings/gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/view_model/models'
import { SettingsActionBar } from '@/components/settings/common/SettingsActionBar'
import { SettingsContentLayout } from '@/components/settings/common/SettingsContentLayout'
import { SettingsPanelHeader } from '@/components/settings/common/SettingsPanelHeader'
import { cn } from '@/lib/utils'
import type { ProviderItem, ProviderModel } from '@/types/settings'

export function ProviderDetailView(props: {
  provider: ProviderItem | null
  onUpdated: () => void
  onDeleted: (id: number) => void
}) {
  const [version, setVersion] = useState(0)
  const { t } = useTranslation()

  if (!props.provider) {
    return (
      <SettingsContentLayout
        header={
          <SettingsPanelHeader
            title={t('settingsPage.primary.providers')}
            description={t('settingsPage.addProvider.doneMessage')}
          />
        }
      />
    )
  }

  return (
    <ProviderDetailInner
      key={`${props.provider.id}-${props.provider.enabled}-${version}`}
      provider={props.provider}
      onUpdated={() => { props.onUpdated(); setVersion(v => v + 1) }}
      onDeleted={props.onDeleted}
      t={t}
    />
  )
}

function ProviderDetailInner(props: {
  provider: ProviderItem
  onUpdated: () => void
  onDeleted: (id: number) => void
  t: ReturnType<typeof useTranslation>['t']
}) {
  const { provider, t } = props

  const [enable, setEnable] = useState(provider.enabled)
  const [name, setName] = useState(provider.provider_name)
  const [baseUrl, setBaseUrl] = useState(provider.base_url)
  const [apiKey, setApiKey] = useState(provider.api_key)
  const [showApiKey, setShowApiKey] = useState(false)
  const [defaultModel, setDefaultModel] = useState(
    provider.models.find(m => m.is_default)?.model ?? provider.models[0]?.model ?? ''
  )
  const [models, setModels] = useState<ProviderModel[]>(provider.models)
  const [refreshing, setRefreshing] = useState(false)
  const [applying, setApplying] = useState(false)
  const [deleting, setDeleting] = useState(false)
  const [searchOpen, setSearchOpen] = useState(false)
  const [searchQuery, setSearchQuery] = useState('')
  const searchInputRef = useRef<HTMLInputElement>(null)

  // Track initial values to compute dirty state
  const initialRef = useRef({
    enable: provider.enabled,
    name: provider.provider_name,
    baseUrl: provider.base_url,
    apiKey: provider.api_key,
    defaultModel: provider.models.find(m => m.is_default)?.model ?? provider.models[0]?.model ?? '',
    modelCount: provider.models.length,
  })

  const dirty =
    enable !== initialRef.current.enable ||
    name !== initialRef.current.name ||
    baseUrl !== initialRef.current.baseUrl ||
    apiKey !== initialRef.current.apiKey ||
    defaultModel !== initialRef.current.defaultModel ||
    models.some(m => m.id === 0)

  const filteredModels = searchOpen && searchQuery
    ? models.filter(m => m.model.toLowerCase().includes(searchQuery.toLowerCase()))
    : models

  const handleRefreshModels = async () => {
    if (!baseUrl || refreshing) return
    setRefreshing(true)
    try {
      const result = await ProviderBinding.RequestProviderModelList({ base_url: baseUrl, api_key: apiKey })
      if (result?.models) {
        const fetched = result.models.map((m) => ({
          id: 0, provider_id: provider.id, model: m.model,
          owned_by: m.owned_by, object: m.object, enable: true,
          alias: null, is_custom: false, is_default: false,
        }))
        const existing = new Map<string, ProviderModel>()
        for (const m of models) existing.set(m.model, m)
        for (const m of fetched) {
          if (!existing.has(m.model)) existing.set(m.model, m)
        }
        const merged = Array.from(existing.values())
        setModels(merged)
        if (merged.length > 0 && !defaultModel) setDefaultModel(merged[0].model)
      }
    } catch {
      // ignore
    } finally {
      setRefreshing(false)
    }
  }

  const handleSetDefaultModel = async (m: ProviderModel) => {
    const prev = defaultModel
    setDefaultModel(m.model)
    if (m.id > 0) {
      try {
        await ProviderBinding.SetDefault({ provider_id: provider.id, model_id: m.id })
        // Treat as saved: reset initial defaultModel so this alone doesn't keep apply dirty
        initialRef.current = { ...initialRef.current, defaultModel: m.model }
      } catch {
        setDefaultModel(prev)
      }
    }
  }

  const handleApply = async () => {
    if (!dirty || applying) return
    setApplying(true)
    try {
      const newModels = models
        .filter(m => m.id === 0)
        .map(m => new ViewModel({
          id: 0, provider_id: provider.id, model: m.model,
          owned_by: m.owned_by, object: m.object,
          enable: m.enable, alias: m.alias, is_custom: m.is_custom, is_default: false,
        }))

      await ProviderBinding.EditProvider(new EditProviderInput({
        provider_id: provider.id,
        provider_name: name,
        base_url: baseUrl,
        api_key: apiKey,
        enable,
        default_model: defaultModel || null,
        models: newModels,
      }))

      props.onUpdated()
    } catch {
      // keep the form active so user can retry
    } finally {
      setApplying(false)
    }
  }

  const handleDelete = async () => {
    if (deleting) return
    setDeleting(true)
    try {
      await ProviderBinding.DeleteProvider({ provider_id: provider.id })
      props.onDeleted(provider.id)
    } catch {
      setDeleting(false)
    }
  }

  const openSearch = () => {
    setSearchOpen(true)
    requestAnimationFrame(() => searchInputRef.current?.focus())
  }

  const closeSearch = () => {
    setSearchOpen(false)
    setSearchQuery('')
  }

  return (
    <SettingsContentLayout
      noContentScroll
      header={
        <SettingsPanelHeader
          title={name}
          description={t('settingsPage.addProvider.doneMessage')}
        />
      }
      footprint={
        <SettingsActionBar
          primaryLabel={t('settingsPage.actions.apply')}
          primaryDisabled={!dirty || applying}
          onPrimaryClick={() => { void handleApply() }}
          dangerLabel={t('settingsPage.providers.delete')}
          onDangerClick={() => { void handleDelete() }}
        />
      }
    >
      <div className="flex min-h-0 flex-1 flex-col space-y-5 pb-4 pt-2">
        {/* Enable toggle */}
        <div className="shrink-0 flex flex-col gap-2">
          <span className="text-sm font-medium text-foreground">{t('settingsPage.addProvider.form.enable')}</span>
          <button
            type="button"
            role="switch"
            aria-checked={enable}
            onClick={() => setEnable(v => !v)}
            className={cn(
              'relative inline-flex h-6 w-11 shrink-0 items-center rounded-full transition-colors',
              enable ? 'bg-primary' : 'bg-muted'
            )}
          >
            <span className={cn(
              'inline-block h-4 w-4 rounded-full bg-white shadow transition-transform',
              enable ? 'translate-x-6' : 'translate-x-1'
            )} />
          </button>
        </div>

        {/* Provider name */}
        <div className="shrink-0 space-y-1.5">
          <label className="text-sm font-medium text-foreground">
            {t('settingsPage.addProvider.form.name')} <span className="text-destructive">*</span>
          </label>
          <input
            required
            value={name}
            onChange={(e) => setName(e.target.value)}
            className="w-full rounded-xl border border-border bg-background px-3 py-2 text-sm outline-none focus:border-primary"
          />
        </div>

        {/* Base URL */}
        <div className="shrink-0 space-y-1.5">
          <label className="text-sm font-medium text-foreground">
            {t('settingsPage.addProvider.form.baseUrl')} <span className="text-destructive">*</span>
          </label>
          <input
            required
            value={baseUrl}
            onChange={(e) => setBaseUrl(e.target.value)}
            className="w-full rounded-xl border border-border bg-background px-3 py-2 text-sm outline-none focus:border-primary"
          />
        </div>

        {/* API key with show/hide toggle */}
        <div className="shrink-0 space-y-1.5">
          <label className="text-sm font-medium text-foreground">{t('settingsPage.addProvider.form.apiKey')}</label>
          <div className="flex items-center rounded-xl border border-border bg-background px-3 py-2 focus-within:border-primary">
            <input
              type={showApiKey ? 'text' : 'password'}
              value={apiKey}
              onChange={(e) => setApiKey(e.target.value)}
              className="min-w-0 flex-1 bg-transparent text-sm outline-none"
            />
            <button
              type="button"
              onClick={() => setShowApiKey(v => !v)}
              title={showApiKey ? t('settingsPage.providers.hideApiKey') : t('settingsPage.providers.showApiKey')}
              className="ml-2 shrink-0 text-muted-foreground transition-colors hover:text-foreground"
            >
              {showApiKey ? <EyeOff size={14} /> : <Eye size={14} />}
            </button>
          </div>
        </div>

        {/* Model list */}
        <div className="flex min-h-0 flex-1 flex-col space-y-1.5">
          <div className="flex shrink-0 items-center justify-between">
            <div className="flex items-center gap-2">
              <label className="text-sm font-medium text-foreground">{t('settingsPage.addProvider.form.modelList')}</label>
              <span className="rounded-full bg-muted px-2 py-0.5 text-xs text-muted-foreground">
                {t('settingsPage.providers.modelsCount', { count: models.length })}
              </span>
            </div>
            <div className="flex items-center gap-1">
              {/* Search box — expands from right to left */}
              <div className={cn(
                'flex items-center overflow-hidden rounded-lg border border-border bg-background px-2 transition-all duration-200',
                searchOpen ? 'w-36 opacity-100' : 'w-0 border-0 opacity-0'
              )}>
                <input
                  ref={searchInputRef}
                  value={searchQuery}
                  onChange={(e) => setSearchQuery(e.target.value)}
                  placeholder={t('settingsPage.providers.searchModelPlaceholder')}
                  className="min-w-0 flex-1 bg-transparent py-1 text-xs outline-none"
                />
                <button
                  type="button"
                  onClick={closeSearch}
                  className="ml-1 shrink-0 text-muted-foreground hover:text-foreground"
                >
                  <X size={11} />
                </button>
              </div>
              {/* Search icon button */}
              {!searchOpen && (
                <button
                  type="button"
                  onClick={openSearch}
                  title={t('settingsPage.providers.searchModelPlaceholder')}
                  className="flex h-8 w-8 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:bg-accent"
                >
                  <Search size={14} />
                </button>
              )}
              <button
                type="button"
                onClick={() => { void handleRefreshModels() }}
                disabled={refreshing}
                title={t('settingsPage.addProvider.form.refreshModels')}
                className="flex h-8 w-8 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:bg-accent disabled:opacity-40"
              >
                <RefreshCw size={14} className={refreshing ? 'animate-spin' : ''} />
              </button>
              <button
                type="button"
                title={t('settingsPage.addProvider.form.addCustomModel')}
                className="flex h-8 w-8 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:bg-accent"
              >
                <Plus size={14} />
              </button>
            </div>
          </div>

          <div className="scrollbar-track-transparent min-h-0 flex-1 overflow-y-auto rounded-xl border border-border">
            {models.length === 0 && (
              <div className="flex items-center justify-center gap-1.5 px-3 py-4 text-xs text-muted-foreground">
                {refreshing ? (
                  <RefreshCw size={12} className="animate-spin" />
                ) : (
                  <>
                    {t('settingsPage.addProvider.form.noModels')}
                    <span
                      className="cursor-pointer underline hover:text-foreground"
                      onClick={() => { void handleRefreshModels() }}
                    >
                      {t('settingsPage.addProvider.form.refreshModels')}
                    </span>
                  </>
                )}
              </div>
            )}
            {models.length > 0 && filteredModels.length === 0 && searchQuery && (
              <div className="flex items-center justify-center px-3 py-4 text-xs text-muted-foreground">
                {t('settingsPage.providers.noSearchResults')}
              </div>
            )}
            {filteredModels.map((m, idx) => (
              <div
                key={`${m.is_custom ? 'custom' : 'fetched'}-${m.model}-${idx}`}
                className="group flex items-center gap-2 px-3 py-2 hover:bg-accent"
              >
                <span className="min-w-0 flex-1 truncate text-sm">{m.model}</span>
                <span className="flex shrink-0 items-center gap-1.5">
                  {m.is_custom && (
                    <span className="rounded border border-border px-1.5 py-0.5 text-[10px] text-muted-foreground">
                      {t('settingsPage.addProvider.form.customTag')}
                    </span>
                  )}
                  {defaultModel === m.model ? (
                    <span className="rounded bg-primary/10 px-1.5 py-0.5 text-[10px] text-primary">
                      {t('settingsPage.addProvider.form.defaultTag')}
                    </span>
                  ) : (
                    <button
                      type="button"
                      onClick={() => { void handleSetDefaultModel(m) }}
                      className="hidden rounded px-1.5 py-0.5 text-[10px] text-muted-foreground hover:bg-accent hover:text-foreground group-hover:inline-block"
                    >
                      {t('settingsPage.addProvider.form.setDefault')}
                    </button>
                  )}
                </span>
              </div>
            ))}
          </div>
        </div>
      </div>
    </SettingsContentLayout>
  )
}
```

- [ ] **Step 2: 验证类型检查**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop/frontend
npx tsc --noEmit 2>&1 | head -30
```

期望：0 错误

- [ ] **Step 3: commit**

```bash
git add frontend/src/components/settings/providers/ProviderDetailView.tsx
git commit -m "feat(provider-detail): dirty tracking, API key toggle, model count, search, apply/delete wiring"
```

---

## Task 8: Frontend AddProviderStepForm — 添加模型搜索

**Files:**
- Modify: `frontend/src/components/settings/providers/AddProviderStepForm.tsx`

- [ ] **Step 1: 在 `AddProviderStepForm.tsx` 中添加搜索状态与搜索 UI**

在现有 imports 末尾添加 `Search` 和 `X`：

找到：
```typescript
import { RefreshCw, Plus, Check, X, Trash2 } from 'lucide-react'
```

替换为：
```typescript
import { RefreshCw, Plus, Check, X, Trash2, Search } from 'lucide-react'
```

在组件 state 末尾（`customInputRef` 之后）添加搜索 state：

找到：
```typescript
  const customInputRef = useRef<HTMLInputElement>(null)
  const customModelsRef = useRef<FetchedModel[]>([])
```

替换为：
```typescript
  const customInputRef = useRef<HTMLInputElement>(null)
  const customModelsRef = useRef<FetchedModel[]>([])
  const searchInputRef = useRef<HTMLInputElement>(null)
  const [searchOpen, setSearchOpen] = useState(false)
  const [searchQuery, setSearchQuery] = useState('')
```

在 `filteredModels` 计算（在 JSX return 之前）中添加：

找到：
```typescript
  return (
    <form id="add-provider-form"
```

在其之前插入（注意要放在 return 语句之前）：

```typescript
  const filteredModels = searchOpen && searchQuery
    ? models.filter(m => m.model.toLowerCase().includes(searchQuery.toLowerCase()))
    : models

  const openSearch = () => {
    setSearchOpen(true)
    requestAnimationFrame(() => searchInputRef.current?.focus())
  }

  const closeSearch = () => {
    setSearchOpen(false)
    setSearchQuery('')
  }

```

找到模型列表标题行：
```tsx
          <label className="text-sm font-medium text-foreground">{t('settingsPage.addProvider.form.modelList')}</label>
          <div className="flex gap-1">
            <button
              type="button"
              onClick={handleRefreshModels}
```

替换为：
```tsx
          <label className="text-sm font-medium text-foreground">{t('settingsPage.addProvider.form.modelList')}</label>
          <div className="flex items-center gap-1">
            {/* Search box — expands from right to left */}
            <div className={cn(
              'flex items-center overflow-hidden rounded-lg border border-border bg-background px-2 transition-all duration-200',
              searchOpen ? 'w-36 opacity-100' : 'w-0 border-0 opacity-0'
            )}>
              <input
                ref={searchInputRef}
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                placeholder={t('settingsPage.providers.searchModelPlaceholder')}
                className="min-w-0 flex-1 bg-transparent py-1 text-xs outline-none"
              />
              <button
                type="button"
                onClick={closeSearch}
                className="ml-1 shrink-0 text-muted-foreground hover:text-foreground"
              >
                <X size={11} />
              </button>
            </div>
            {!searchOpen && (
              <button
                type="button"
                onClick={openSearch}
                title={t('settingsPage.providers.searchModelPlaceholder')}
                className="flex h-8 w-8 items-center justify-center rounded-lg text-muted-foreground transition-colors hover:bg-accent"
              >
                <Search size={14} />
              </button>
            )}
            <button
              type="button"
              onClick={handleRefreshModels}
```

找到空模型状态（使用原始 `models` 的两处）：
```tsx
          {models.length === 0 && !addingCustom && (
```

和模型列表渲染（`{models.map(...`）：
```tsx
          {models.map((m, idx) => (
```

替换为：
```tsx
          {models.length === 0 && !addingCustom && (
```
（此处不变）

在空模型状态提示之后、模型列表渲染之前插入无搜索结果提示，同时将 `models.map` 改为 `filteredModels.map`：

找到：
```tsx
          {models.map((m, idx) => (
```

替换为：
```tsx
          {models.length > 0 && filteredModels.length === 0 && searchQuery && (
            <div className="flex items-center justify-center px-3 py-4 text-xs text-muted-foreground">
              {t('settingsPage.providers.noSearchResults')}
            </div>
          )}
          {filteredModels.map((m, idx) => (
```

- [ ] **Step 2: 验证类型检查**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop/frontend
npx tsc --noEmit 2>&1 | head -20
```

期望：0 错误

- [ ] **Step 3: 全量 Go 测试**

```bash
cd /Users/linhuafeng/Work/lemon_tea_desktop
go test ./backend/... -v 2>&1 | tail -20
```

期望：所有测试 PASS

- [ ] **Step 4: commit**

```bash
git add frontend/src/components/settings/providers/AddProviderStepForm.tsx
git commit -m "feat(add-provider): add model list search with expand animation"
```

---

## 验收标准

| 功能 | 验证方法 |
|------|---------|
| 供应商列表宽度加宽 | 目视确认 sidebar 宽度（w-72 = 288px） |
| 列表右侧启用开关 | 点击开关，供应商 enabled 状态切换，Detail View 同步重置 |
| 更多菜单 → 设为默认 | 点击后供应商列表中出现默认标签，其他供应商标签消失 |
| 更多菜单 → 删除 | 供应商从列表中消失，Detail View 切换到下一个 |
| 应用按钮默认禁用 | 初始状态下 Apply 按钮灰色不可点击 |
| 修改后应用按钮启用 | 修改任一字段后 Apply 按钮变可用 |
| API 密钥显示切换 | 点击眼睛图标，密钥明文/密文切换 |
| 模型数量标签 | 模型列表标题旁显示"N 个模型" |
| 模型搜索（Detail View）| 点击搜索图标，搜索框展开，输入过滤，无结果时提示，X 关闭搜索 |
| 模型搜索（Add Provider）| 同上 |
| 模型设为默认（Detail View）| 即时调用 SetDefault，默认标签切换 |
