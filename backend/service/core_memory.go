package service

import (
	"context"
	"strings"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/prompts"
)

func (s *Service) RenderCoreMemorySnapshot(ctx context.Context) string {
	if s.memoryStorage == nil {
		return ""
	}
	snapshot, err := s.memoryStorage.RenderCoreMemorySnapshot(ctx)
	if err != nil {
		logger.Error("render core memory snapshot error:", err)
		return ""
	}
	return strings.TrimSpace(snapshot)
}

func (s *Service) promptSetWithCoreMemory(ctx context.Context, base prompts.PromptSet) prompts.PromptSet {
	prefs, err := s.loadAppPreferences(ctx)
	if err != nil || !prefs.MemorySystemEnabled {
		return base
	}
	snapshot := s.RenderCoreMemorySnapshot(ctx)
	if snapshot == "" {
		return base
	}
	base.MainAgentSystem = strings.TrimSpace(base.MainAgentSystem) + "\n\n" + snapshot
	return base
}
