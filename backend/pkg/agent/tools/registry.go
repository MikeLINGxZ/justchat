// backend/pkg/agent/tools/registry.go
package tools

import (
	"encoding/json"
	"sync"
)

const (
	CategoryBuiltin = "builtin"
	CategoryUser    = "user"
)

type ToolMeta struct {
	Name            string
	Description     string
	Category        string
	RequiresConfirm bool
	FormatPurpose   func(args json.RawMessage) string
}

type Registry struct {
	mu    sync.RWMutex
	tools map[string]ToolMeta
	order []string
}

func NewRegistry() *Registry {
	return &Registry{
		tools: make(map[string]ToolMeta),
	}
}

func (r *Registry) Register(meta ToolMeta) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, exists := r.tools[meta.Name]; !exists {
		r.order = append(r.order, meta.Name)
	}
	r.tools[meta.Name] = meta
}

func (r *Registry) Get(name string) (ToolMeta, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	meta, ok := r.tools[name]
	return meta, ok
}

func (r *Registry) BuiltinTools() []ToolMeta {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []ToolMeta
	for _, name := range r.order {
		if r.tools[name].Category == CategoryBuiltin {
			result = append(result, r.tools[name])
		}
	}
	return result
}

func (r *Registry) UserTools() []ToolMeta {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []ToolMeta
	for _, name := range r.order {
		if r.tools[name].Category == CategoryUser {
			result = append(result, r.tools[name])
		}
	}
	return result
}

func (r *Registry) EnabledTools(enabled []string) []ToolMeta {
	enabledSet := make(map[string]bool, len(enabled))
	for _, name := range enabled {
		enabledSet[name] = true
	}

	r.mu.RLock()
	defer r.mu.RUnlock()
	var result []ToolMeta
	for _, name := range r.order {
		meta := r.tools[name]
		if meta.Category == CategoryBuiltin || enabledSet[name] {
			result = append(result, meta)
		}
	}
	return result
}
