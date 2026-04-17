package service

import (
	"context"

	"github.com/wailsapp/wails/v3/pkg/application"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/agents/memory/lifecycle"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/agents/memory/search"
	memory_storage "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/agents/memory/storage"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/i18n"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/llm_provider/agents"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/prompts"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/plugin"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/storage"
)

const (
	WindowNameHome       = "window_home"
	WindowNameOnboarding = "window_onboarding"
	WindowNameSettings   = "window_settings"
)

type Service struct {
	storage         *storage.Storage
	app             *application.App
	prompts         prompts.PromptSet
	pluginManager   *plugin.Manager
	memoryStorage   *memory_storage.Storage
	memoryCache     *memoryPrefetchCache
	memoryLifecycle *lifecycle.Manager
	memorySearcher  *search.HybridSearcher
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) ServiceStartup(ctx context.Context, options application.ServiceOptions) error {

	istorage, err := storage.NewStorage()
	if err != nil {
		return err
	}

	s.storage = istorage
	s.app = application.Get()

	// 初始化记忆系统存储（复用主数据库连接）
	memStorage, memErr := memory_storage.NewStorage(istorage.DB())
	if memErr != nil {
		logger.Warm("memory storage init failed:", memErr)
	} else {
		s.memoryStorage = memStorage
		// 初始化嵌入表
		if embErr := memStorage.AutoMigrateEmbeddings(); embErr != nil {
			logger.Warm("memory embeddings migration failed:", embErr)
		}
		// 创建混合检索引擎（暂无 embedder，后续可通过 Ollama 注入）
		s.memorySearcher = search.NewHybridSearcher(memStorage, nil)
		// 启动记忆生命周期管理（巩固/遗忘/矛盾检测）
		s.memoryLifecycle = lifecycle.NewManager(memStorage)
		s.memoryLifecycle.Start()
	}
	s.memoryCache = newMemoryPrefetchCache()
	if prefs, prefsErr := s.loadAppPreferences(ctx); prefsErr == nil {
		i18n.SetCurrentLocale(string(prefs.Language))
	}
	if err := s.reloadPromptSet(); err != nil {
		logger.Warm("load prompt set fallback:", err)
	}

	agents.SyncCustomAgentsToRegistry()

	if err := s.syncCustomMCPTools(ctx); err != nil {
		return err
	}

	pluginMgr, pluginErr := plugin.NewManager(s.app)
	if pluginErr != nil {
		logger.Warm("plugin manager creation failed:", pluginErr)
	} else {
		s.pluginManager = pluginMgr
		if initErr := pluginMgr.Init(); initErr != nil {
			logger.Warm("plugin system init failed:", initErr)
		}
	}

	if err := s.recoverStaleRunningTasks(ctx); err != nil {
		return err
	}

	return nil
}
