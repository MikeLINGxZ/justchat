package lifecycle

import (
	"context"
	"math"
	"sync"
	"time"

	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/agents/memory/models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/agents/memory/storage"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"
)

// Manager 管理记忆生命周期：巩固、遗忘、矛盾检测。
type Manager struct {
	storage *storage.Storage
	mu      sync.Mutex
	stopCh  chan struct{}
}

func NewManager(s *storage.Storage) *Manager {
	return &Manager{
		storage: s,
		stopCh:  make(chan struct{}),
	}
}

// Start 启动定时生命周期任务。
func (m *Manager) Start() {
	// 启动时立即执行一次
	go m.runOnce()

	// 巩固：每小时执行
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				m.runOnce()
			case <-m.stopCh:
				return
			}
		}
	}()
}

// Stop 停止定时任务。
func (m *Manager) Stop() {
	close(m.stopCh)
}

func (m *Manager) runOnce() {
	if !m.mu.TryLock() {
		return // 已有任务在执行
	}
	defer m.mu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	m.consolidate(ctx)
	m.detectContradictions(ctx)
	m.forget(ctx)
}

// ---- 巩固 ----

func (m *Manager) consolidate(ctx context.Context) {
	// 1. 提升频繁召回记忆的重要性
	if err := m.boostFrequentMemories(ctx); err != nil {
		logger.Error("consolidate: boost frequent memories error:", err)
	}
}

// boostFrequentMemories 对 RecallCount > 5 且 Importance < 0.9 的记忆提升重要性。
func (m *Manager) boostFrequentMemories(ctx context.Context) error {
	db := m.storage.DB()
	return db.WithContext(ctx).Exec(`
		UPDATE memories
		SET importance = MIN(0.9, importance + 0.05),
		    updated_at = ?
		WHERE recall_count > 5
		  AND importance < 0.9
		  AND is_forgotten = 0
	`, time.Now()).Error
}

// ---- 矛盾检测 ----

// ContradictionPair 表示一对可能矛盾的记忆。
type ContradictionPair struct {
	MemoryA models.Memory
	MemoryB models.Memory
}

func (m *Manager) detectContradictions(ctx context.Context) {
	pairs, err := m.findContradictions(ctx)
	if err != nil {
		logger.Error("detect contradictions error:", err)
		return
	}
	if len(pairs) == 0 {
		return
	}

	// 处理策略：降低较旧记忆的信任分数
	for _, pair := range pairs {
		older := pair.MemoryA
		if pair.MemoryB.CreatedAt.Before(pair.MemoryA.CreatedAt) {
			older = pair.MemoryB
		}
		if err := m.storage.AdjustTrustScore(ctx, older.ID, -0.10); err != nil {
			logger.Error("adjust trust for contradiction error:", err)
		}
	}

	logger.Warm("detected", len(pairs), "potential memory contradictions, adjusted trust scores")
}

// findContradictions 查找共享关键词但内容差异大的记忆对。
// 简化实现：查找 summary 相似（LIKE 前 5 字匹配）但 content 不同的活跃记忆。
func (m *Manager) findContradictions(ctx context.Context) ([]ContradictionPair, error) {
	db := m.storage.DB()

	// 获取最近 30 天的活跃记忆
	var memories []models.Memory
	cutoff := time.Now().AddDate(0, 0, -30)
	if err := db.WithContext(ctx).
		Where("is_forgotten = ? AND created_at >= ?", false, cutoff).
		Order("created_at DESC").
		Limit(200).
		Find(&memories).Error; err != nil {
		return nil, err
	}

	var pairs []ContradictionPair
	for i := 0; i < len(memories); i++ {
		for j := i + 1; j < len(memories); j++ {
			a, b := memories[i], memories[j]
			// 检查 summary 相似性（前 5 字相同）
			runesA := []rune(a.Summary)
			runesB := []rune(b.Summary)
			prefixLen := 5
			if len(runesA) < prefixLen || len(runesB) < prefixLen {
				continue
			}
			if string(runesA[:prefixLen]) != string(runesB[:prefixLen]) {
				continue
			}
			// summary 前缀相同但内容不同 → 可能矛盾
			if a.Content != b.Content {
				pairs = append(pairs, ContradictionPair{MemoryA: a, MemoryB: b})
			}
			if len(pairs) >= 10 {
				return pairs, nil
			}
		}
	}
	return pairs, nil
}

// ---- 遗忘 ----

func (m *Manager) forget(ctx context.Context) {
	if err := m.applyForgetting(ctx); err != nil {
		logger.Error("forget error:", err)
	}
}

// applyForgetting 根据衰减公式标记应遗忘的记忆。
// decay = e^(-0.03 * days_since_last_recall) * importance * trust_score
// 遗忘条件：decay < 0.1 且 RecallCount < 2 且 TrustScore < 0.5
// 保护：Importance >= 0.8 或 TrustScore >= 0.8 或创建 < 7 天不遗忘
func (m *Manager) applyForgetting(ctx context.Context) error {
	db := m.storage.DB()

	// 获取候选记忆
	cutoff := time.Now().AddDate(0, 0, -7)
	var candidates []models.Memory
	if err := db.WithContext(ctx).
		Where("is_forgotten = ?", false).
		Where("importance < 0.8").
		Where("trust_score < 0.8").
		Where("recall_count < 2").
		Where("created_at < ?", cutoff).
		Find(&candidates).Error; err != nil {
		return err
	}

	now := time.Now()
	var forgetIDs []uint
	for _, m := range candidates {
		lastRecall := m.CreatedAt
		if m.LastRecalledAt != nil {
			lastRecall = *m.LastRecalledAt
		}
		daysSince := now.Sub(lastRecall).Hours() / 24.0
		decay := math.Exp(-0.03*daysSince) * m.Importance * m.TrustScore

		if decay < 0.1 {
			forgetIDs = append(forgetIDs, m.ID)
		}
	}

	if len(forgetIDs) == 0 {
		return nil
	}

	logger.Warm("forgetting", len(forgetIDs), "memories due to decay")
	return db.WithContext(ctx).
		Model(&models.Memory{}).
		Where("id IN ?", forgetIDs).
		Updates(map[string]any{
			"is_forgotten": true,
			"updated_at":   now,
		}).Error
}
