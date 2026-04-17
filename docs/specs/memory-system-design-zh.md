# 记忆系统技术方案

## 概述

为 Lemon Tea Desktop 构建完整的长期记忆能力，使 AI 助手能够跨会话记住用户的经历、偏好、计划和情感状态，提供个性化、有温度的对话体验。

### 设计原则

- **隐形运作**：记忆的读取和写入对用户完全透明，不暴露任何内部机制
- **桌面优先**：所有数据本地存储，嵌入模型本地运行，零隐私泄露
- **渐进增强**：核心功能仅依赖 SQLite FTS5 全文检索，向量搜索作为可选增强
- **不阻塞交互**：记忆写入完全异步，检索采用预取缓存实现零延迟注入
- **容错隔离**：记忆子系统的任何异常不影响主对话流程
- **质量优先**：通过信任分数和矛盾检测机制保证记忆准确性

---

## 1. 现状分析

### 1.1 已有代码

| 模块 | 路径 | 状态 |
|------|------|------|
| 数据模型 | `backend/agents/memory/models/memory.go` | 完整，含 Summary/Content/Type/Time/Location/Characters/Emotion/Importance 等字段 |
| 工具集 | `backend/agents/memory/tools/memory.go` | 完整，含 write_memory / read_memory / edit_memory / get_current_time |
| 存储层 | `backend/agents/memory/storage/memory.go` | 基本完整，支持多维度 LIKE 查询 |
| 记忆 Agent | `backend/agents/memory/agent.go` | 双实现（ReAct + Graph），未接入主流程 |
| Agent 提示词 | `backend/agents/memory/agent.prompt.v1.md` | 详细的认知型 AI 伙伴提示词 |
| 提示词加载 | `backend/pkg/prompts/prompts.go` | MemorySystem 字段已存在，从 `system.memory.md` 加载 |

### 1.2 缺失部分

1. **主流程集成**：`chat_completion_runner.go` 的 `run()` 方法中没有任何记忆相关逻辑
2. **存储统一**：记忆使用独立的 `agent_memory.db`，与主应用的 `data.db` 分离
3. **全文检索**：当前仅用 LIKE 模糊匹配，无 FTS5 索引
4. **向量搜索**：`EmbeddingID` 字段已预留但未实现嵌入存储和检索
5. **自动触发**：缺少判断何时读取/写入记忆的触发机制
6. **上下文注入**：相关记忆未注入到对话中
7. **前端管理 UI**：没有记忆查看、搜索、编辑界面
8. **生命周期管理**：缺少信任评分、矛盾检测、记忆巩固和遗忘机制

---

## 2. 系统架构

### 2.1 整体位置

```
┌──────────────────────────────────────────────────────────────────┐
│                        Frontend (React)                          │
│  ┌──────────┐  ┌──────────┐  ┌──────────────────────────────┐   │
│  │ Chat UI  │  │ Settings │  │ Memory Management (新增)     │   │
│  └────┬─────┘  └──────────┘  └──────────────┬───────────────┘   │
│       │                                      │                   │
├───────┼──────────────────────────────────────┼───────────────────┤
│       │           Wails RPC Bindings         │                   │
├───────┼──────────────────────────────────────┼───────────────────┤
│       ▼                                      ▼                   │
│  ┌─────────────────────┐   ┌────────────────────────────────┐   │
│  │ Service.Completions │   │ Service.Memory* (新增 API)     │   │
│  └────────┬────────────┘   └────────────────────────────────┘   │
│           ▼                                                      │
│  ┌──────────────────────────────────────────────────────┐       │
│  │              completionRunner.run()                    │       │
│  │                                                        │       │
│  │  ┌──────────────────┐                                 │       │
│  │  │ 注入预取缓存 (新增)│ ← 来自上一轮的异步预取结果     │       │
│  │  └────────┬─────────┘                                 │       │
│  │           ▼                                            │       │
│  │  ┌──────────────────┐                                 │       │
│  │  │ 主执行流程        │ (直答 / 工作流)                 │       │
│  │  └────────┬─────────┘                                 │       │
│  │           ▼                                            │       │
│  │  ┌──────────────────┐  ┌──────────────────────┐      │       │
│  │  │ 记忆编码 (异步)   │  │ 预取下一轮上下文 (异步)│      │       │
│  │  └──────────────────┘  └──────────────────────┘      │       │
│  └──────────────────────────────────────────────────────┘       │
│                                                                  │
│  ┌──────────────────────────────────────────────────────┐       │
│  │                  Memory Module                        │       │
│  │  ┌──────────┐ ┌──────────┐ ┌────────┐ ┌──────────┐  │       │
│  │  │ Storage  │ │ Agent    │ │ Search │ │Lifecycle │  │       │
│  │  │(SQLite+  │ │ (ReAct)  │ │(Hybrid)│ │(Trust/   │  │       │
│  │  │ FTS5)    │ │          │ │        │ │ Forget)  │  │       │
│  │  └──────────┘ └──────────┘ └────────┘ └──────────┘  │       │
│  └──────────────────────────────────────────────────────┘       │
│                                                                  │
│  ┌──────────────────────────────────────────────────────┐       │
│  │        SQLite (data.db) — WAL mode                    │       │
│  │  memories | memory_entities | memory_entity_links     │       │
│  │  memory_embeddings | memories_fts                     │       │
│  └──────────────────────────────────────────────────────┘       │
└──────────────────────────────────────────────────────────────────┘
```

### 2.2 数据流（优化后）

```
用户发送消息
    │
    ▼
① 注入预取缓存（零延迟）
    │  从上一轮异步预取的缓存中取出记忆上下文
    │  用 <memory-context> 围栏标签包裹，注入用户消息
    │  缓存为空时回退为同步检索
    │
    ▼
② 主执行流程（直答/工作流）
    │  LLM 生成回复（上下文中已包含记忆信息）
    │
    ▼（以下两步并行异步执行，不阻塞）
    │
    ├─► ③ 记忆编码
    │      Memory Agent 分析对话 → 决定写入/更新/跳过
    │      仅主 Agent 上下文触发，工作流 Worker 不写入
    │
    └─► ④ 预取下一轮上下文
           基于当前对话内容异步检索相关记忆
           结果缓存，下一轮直接使用
```

**与原方案的关键差异**：
- 检索从「同步阻塞」变为「异步预取 + 缓存」，首轮回退为同步
- 上下文从「追加到 system prompt」变为「围栏标签注入用户消息」
- 新增容错隔离，记忆异常不影响主流程

---

## 3. 存储层

### 3.1 统一到 data.db

将 Memory 表迁入主数据库 `data.db`，重构 Storage 接受外部 `*gorm.DB` 实例：

```go
// 改造后
func NewStorage(db *gorm.DB) (*Storage, error) {
    err := db.AutoMigrate(&models.Memory{}, &models.MemoryEntity{}, &models.MemoryEntityLink{})
    if err != nil {
        return nil, err
    }
    return &Storage{sqliteDb: db}, nil
}
```

### 3.2 数据模型（增强）

在现有 Memory 模型基础上新增字段和关联表：

```go
type Memory struct {
    // ... 已有字段保持不变 ...

    // 新增：信任分数
    TrustScore  float64 `gorm:"default:0.5"`  // [0.0, 1.0]，默认 0.5
    RetrievalCount int  `gorm:"default:0"`    // 被检索命中的次数（含未最终注入的）
}

// 新增：实体表
type MemoryEntity struct {
    ID   uint   `gorm:"primarykey"`
    Name string `gorm:"type:varchar(200);uniqueIndex"`
    Type string `gorm:"type:varchar(50);index"` // person / location / event / topic
}

// 新增：记忆-实体关联表
type MemoryEntityLink struct {
    ID       uint `gorm:"primarykey"`
    MemoryID uint `gorm:"index"`
    EntityID uint `gorm:"index"`
    Role     string `gorm:"type:varchar(50)"` // subject / location / participant / topic
}
```

### 3.3 FTS5 全文检索索引（替代 LIKE）

用 FTS5 替换现有的 LIKE 模糊匹配，检索性能从 O(n) 提升到 O(log n)：

```sql
-- 创建 FTS5 虚拟表，对 summary 和 content 建立全文索引
CREATE VIRTUAL TABLE IF NOT EXISTS memories_fts USING fts5(
    summary,
    content,
    content='memories',
    content_rowid='id',
    tokenize='unicode61'   -- 支持中文 Unicode 分词
);

-- 触发器：自动同步 FTS 索引
CREATE TRIGGER memories_ai AFTER INSERT ON memories BEGIN
    INSERT INTO memories_fts(rowid, summary, content) VALUES (new.id, new.summary, new.content);
END;
CREATE TRIGGER memories_ad AFTER DELETE ON memories BEGIN
    INSERT INTO memories_fts(memories_fts, rowid, summary, content) VALUES('delete', old.id, old.summary, old.content);
END;
CREATE TRIGGER memories_au AFTER UPDATE ON memories BEGIN
    INSERT INTO memories_fts(memories_fts, rowid, summary, content) VALUES('delete', old.id, old.summary, old.content);
    INSERT INTO memories_fts(rowid, summary, content) VALUES (new.id, new.summary, new.content);
END;
```

检索查询：

```go
func (s *Storage) FTSSearch(ctx context.Context, keywords []string, limit int) ([]models.Memory, error) {
    // FTS5 MATCH 查询，3× limit 留出重排空间
    query := strings.Join(keywords, " OR ")
    var memories []models.Memory
    err := s.sqliteDb.WithContext(ctx).
        Raw(`SELECT m.* FROM memories m
             JOIN memories_fts ON memories_fts.rowid = m.id
             WHERE memories_fts MATCH ? AND m.is_forgotten = 0
             ORDER BY rank
             LIMIT ?`, query, limit*3).
        Scan(&memories).Error
    return memories, err
}
```

### 3.4 嵌入向量表（可选增强）

```sql
CREATE TABLE memory_embeddings (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    memory_id   INTEGER NOT NULL UNIQUE REFERENCES memories(id) ON DELETE CASCADE,
    vector      BLOB NOT NULL,
    model_name  VARCHAR(100) NOT NULL,
    dimensions  INTEGER NOT NULL,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP
);
```

### 3.5 索引优化

```sql
CREATE INDEX idx_memories_active_ranked ON memories(is_forgotten, trust_score DESC, importance DESC, created_at DESC);
CREATE INDEX idx_memory_entity_links_memory ON memory_entity_links(memory_id);
CREATE INDEX idx_memory_entity_links_entity ON memory_entity_links(entity_id);
```

---

## 4. 主流程集成

### 4.1 集成点

在 `completionRunner.run()` 方法（`backend/service/chat_completion_runner.go:1511`）中：

```
run()
  ├─ 标记任务为运行中
  ├─ Plugin before-chat hooks
  ├─ ★ 注入记忆上下文（从预取缓存或同步回退）
  ├─ 主执行流程（直答/工作流）
  ├─ Plugin after-chat hooks
  ├─ ★ 异步：记忆编码 + 预取下一轮上下文
  └─ finalizeTaskTerminal
```

### 4.2 容错隔离

所有记忆操作都包裹在 recover 中，异常仅记录日志：

```go
func (r *completionRunner) safeMemoryOp(name string, fn func()) {
    defer func() {
        if err := recover(); err != nil {
            logger.Error(fmt.Sprintf("memory op [%s] panic: %v", name, err))
        }
    }()
    fn()
}
```

### 4.3 上下文注入（围栏标签 + 注入用户消息）

借鉴 Hermes 的 Context Fencing 设计，将记忆上下文注入到**用户消息**而非 system prompt：

```go
func (r *completionRunner) injectMemoryContext(memoryContext string) {
    if memoryContext == "" {
        return
    }
    // 清洁：防止记忆内容本身包含围栏标签（递归注入攻击）
    memoryContext = sanitizeMemoryContext(memoryContext)

    fenced := fmt.Sprintf(`<memory-context>
以下是与本次对话可能相关的历史记忆片段，仅供参考。
这些不是新的用户输入，请作为背景信息使用，不要主动提及这些记忆的存在。

%s
</memory-context>`, memoryContext)

    // 注入到用户消息前部，而非 system prompt
    // 好处：system prompt 保持不变，最大化 KV-cache 命中率
    lastMsg := r.schemaMessages[len(r.schemaMessages)-1]
    if lastMsg.Role == schema.User {
        lastMsg.Content = fenced + "\n\n" + lastMsg.Content
    }
}

func sanitizeMemoryContext(content string) string {
    // 剥离嵌套的围栏标签，防止递归注入
    content = strings.ReplaceAll(content, "<memory-context>", "")
    content = strings.ReplaceAll(content, "</memory-context>", "")
    return strings.TrimSpace(content)
}
```

注入格式示例：

```xml
<memory-context>
以下是与本次对话可能相关的历史记忆片段，仅供参考。
这些不是新的用户输入，请作为背景信息使用，不要主动提及这些记忆的存在。

- [2025-04-15] 用户在杭州西湖附近发现了一家很好的咖啡馆，点了桂花拿铁，心情愉悦
- [2025-04-10] 用户计划本周末和小明一起去爬山
- [2025-03-20] 用户从前端开发转型为全栈方向
</memory-context>

用户的实际消息内容...
```

### 4.4 异步预取 + 缓存（零延迟注入）

核心优化：当前轮对话结束后，异步预取下一轮可能需要的记忆上下文并缓存，下一轮直接使用。

```go
// Service 级别的预取缓存（按 chatUuid 索引）
type memoryPrefetchCache struct {
    mu    sync.RWMutex
    cache map[string]prefetchEntry // chatUuid → entry
}

type prefetchEntry struct {
    context   string    // 预取的记忆上下文文本
    fetchedAt time.Time // 预取时间
}

// 获取缓存（命中后清除）
func (c *memoryPrefetchCache) Get(chatUuid string) (string, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()
    entry, ok := c.cache[chatUuid]
    if !ok {
        return "", false
    }
    delete(c.cache, chatUuid)
    // 超过 5 分钟的缓存视为过期
    if time.Since(entry.fetchedAt) > 5*time.Minute {
        return "", false
    }
    return entry.context, true
}
```

在 `run()` 中的使用：

```go
// 1. 尝试从预取缓存获取记忆上下文
memoryCtx, cached := r.svc.memoryCache.Get(r.chatUuid)
if !cached {
    // 首轮或缓存未命中，同步检索（有 100ms 超时保护）
    ctx, cancel := context.WithTimeout(runCtx, 100*time.Millisecond)
    memoryCtx = r.retrieveMemoryContext(ctx)
    cancel()
}
r.safeMemoryOp("inject", func() {
    r.injectMemoryContext(memoryCtx)
})

// ... 主执行流程 ...

// 2. 对话结束后，异步预取 + 异步编码并行执行
go r.safeMemoryOp("encode", func() {
    r.encodeMemoriesAsync(r.schemaMessages, assistantContent)
})
go r.safeMemoryOp("prefetch", func() {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()
    nextCtx := r.retrieveMemoryContext(ctx)
    r.svc.memoryCache.Set(r.chatUuid, nextCtx)
})
```

### 4.5 记忆编码（异步，含工作流隔离）

借鉴 Hermes 的子 Agent 记忆隔离设计：**仅主 Agent 上下文触发记忆编码，工作流 Worker Agent 不写入**。

```go
func (r *completionRunner) encodeMemoriesAsync(messages []*schema.Message, assistantContent string) {
    userMessage := extractLastUserMessage(messages)
    if !shouldEncodeMemories(userMessage, assistantContent) {
        return
    }

    // 工作流隔离：如果本轮走了工作流路径，不触发记忆编码
    // 避免 Planner/Worker/Synthesizer 的中间过程污染记忆
    if r.getWorkflowHandoff() != nil {
        return
    }

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    _, err := r.svc.memoryAgent.Stream(ctx, messages)
    if err != nil {
        logger.Error("memory encoding error", err)
    }
}
```

---

## 5. 混合检索引擎

### 5.1 架构

```
用户消息
    │
    ├──► ① FTS5 全文检索 ──► 候选集 A（3× limit）
    │
    ├──► ② 实体关联检索 ──► 候选集 B
    │       提取消息中的人名/地名 → 查实体表 → 关联记忆
    │
    ├──► ③ 向量检索（可选）──► 候选集 C
    │
    └──► ④ 元数据过滤（时间/类型）──► 缩小范围
          │
          ▼
    合并去重 → 多维加权打分 → 信任过滤 → Top-K 结果
```

### 5.2 FTS5 全文检索（基础，始终可用）

```go
func (s *HybridSearcher) FTSSearch(ctx context.Context, keywords []string, limit int) ([]ScoredMemory, error) {
    memories, err := s.storage.FTSSearch(ctx, keywords, limit)
    if err != nil {
        return nil, err
    }
    results := make([]ScoredMemory, len(memories))
    for i, m := range memories {
        results[i] = ScoredMemory{
            Memory:      m,
            FTSScore:    1.0 / float64(i+1), // rank-based score
        }
    }
    return results, nil
}
```

### 5.3 实体关联检索（新增）

从用户消息中提取人名、地名等实体，通过实体表反查关联记忆：

```go
func (s *HybridSearcher) EntitySearch(ctx context.Context, message string, limit int) ([]ScoredMemory, error) {
    // 1. 从消息中提取可能的实体名称（简单规则匹配）
    entities := extractEntities(message)

    // 2. 查找匹配的实体
    var entityIDs []uint
    s.db.Model(&MemoryEntity{}).Where("name IN ?", entities).Pluck("id", &entityIDs)

    // 3. 通过关联表找到相关记忆
    var memoryIDs []uint
    s.db.Model(&MemoryEntityLink{}).Where("entity_id IN ?", entityIDs).
        Distinct("memory_id").Pluck("memory_id", &memoryIDs)

    // 4. 加载记忆详情
    ...
}
```

### 5.4 向量检索（可选增强）

通过 Ollama 本地嵌入模型（如 `bge-m3`），使用 Eino `embedding/ollama` 组件。Ollama 不可用时自动跳过。

### 5.5 多维加权打分

```go
type ScoredMemory struct {
    Memory          models.Memory
    FTSScore        float64  // FTS5 rank [0, 1]
    EntityScore     float64  // 实体关联度 [0, 1]
    VectorScore     float64  // 向量相似度 [0, 1]（可选）
    RecencyScore    float64  // 时间新近度 [0, 1]
    FinalScore      float64  // 综合得分
}

func (sm *ScoredMemory) ComputeFinalScore(vectorEnabled bool) {
    trust := sm.Memory.TrustScore
    if trust < minTrustThreshold { // 0.3
        sm.FinalScore = 0 // 信任分过低，排除
        return
    }

    recency := math.Exp(-0.05 * daysSince(sm.Memory.CreatedAt))
    sm.RecencyScore = recency

    if vectorEnabled {
        sm.FinalScore = (0.30*sm.FTSScore +
            0.15*sm.EntityScore +
            0.30*sm.VectorScore +
            0.10*recency +
            0.15*sm.Memory.Importance) * trust
    } else {
        // 无向量时，FTS 和实体检索权重提升
        sm.FinalScore = (0.45*sm.FTSScore +
            0.20*sm.EntityScore +
            0.15*recency +
            0.20*sm.Memory.Importance) * trust
    }
}
```

### 5.6 Token 预算控制

注入的记忆上下文控制在 **500 tokens**（约 1000 汉字）以内：

- Top-K 结果按 FinalScore 降序排列
- 逐条累加文本长度，超出预算时截断
- 每条记忆压缩为一行摘要格式：`[日期] 摘要内容（地点、人物等关键元素）`

---

## 6. 信任系统与矛盾检测

### 6.1 信任分数（借鉴 Hermes Holographic）

每条记忆携带 `TrustScore`，范围 `[0.0, 1.0]`，影响检索排序和注入决策。

**初始值**：0.5

**反馈调整（非对称）**：
- 记忆被检索后 LLM 引用了 → 有用反馈：`+0.05`
- 记忆被检索但 LLM 未引用 → 无用反馈：`-0.02`
- 用户明确否认记忆内容 → 强负反馈：`-0.15`
- 用户确认记忆内容 → 强正反馈：`+0.10`

**最低阈值**：`TrustScore < 0.3` 的记忆不参与检索结果排序。

### 6.2 矛盾检测

巩固阶段自动扫描：找到**实体重叠度高但内容相似度低**的记忆对，标记为潜在矛盾。

```go
func (s *LifecycleManager) DetectContradictions(ctx context.Context) ([]MemoryPair, error) {
    // 1. 按共享实体分组：找到共享 ≥2 个实体的记忆对
    // 2. 计算内容相似度（FTS5 token overlap 或向量余弦）
    // 3. 实体重叠 ≥ 2 且内容相似度 < 0.3 → 标记为矛盾
    // 4. 返回矛盾对列表，交由 Memory Agent 在下次编码时处理
}
```

处理策略：
- 矛盾对中较新的记忆保留，较旧的降低 TrustScore
- 严重矛盾（同一实体完全相反的事实）由 Memory Agent 合并或废弃旧记忆

---

## 7. 记忆生命周期

### 7.1 五阶段模型

```
  编码 (Encoding)       巩固 (Consolidation)     矛盾检测 (Contradiction)
  每轮对话后异步执行     定时（启动 + 每小时）     巩固阶段自动执行
       │                       │                        │
       ▼                       ▼                        ▼
  ┌─────────┐           ┌──────────┐            ┌──────────────┐
  │ 新记忆   │──────────►│ 合并/强化 │────────────│ 发现冲突      │
  │ 写入 DB  │           │ 更新嵌入  │            │ 降低旧信任    │
  │ 提取实体 │           │ 提升召回  │            │ Agent 合并    │
  └─────────┘           └──────────┘            └──────────────┘
                              │
       ┌──────────────────────┘
       ▼                         ▼
  ┌─────────┐              ┌─────────────┐
  │ 检索     │              │ 遗忘         │
  │ 信任加权 │              │ 衰减+软删除  │
  │ 预取缓存 │              │ 可恢复       │
  └─────────┘              └─────────────┘
  检索 (Retrieval)          遗忘 (Forgetting)
  预取缓存/同步回退          定时（启动 + 每日）
```

### 7.2 编码（Encoding）

- **触发时机**：每轮对话完成后异步执行
- **执行者**：Memory Agent（ReAct 模式）
- **工作流隔离**：仅直答路径触发，工作流路径（Planner/Worker/Synthesizer）不触发
- **行为**：
  1. 分析对话内容，提取时间/地点/人物/情绪/重要性
  2. 调用 read_memory 检查是否已有相似记忆
  3. 已有则 edit_memory 更新，否则 write_memory 新建
  4. 写入后触发实体提取：从记忆内容中识别人名、地名、事件名，存入实体表和关联表

### 7.3 巩固（Consolidation）

定时任务，应用启动时 + 每小时执行一次：

1. **碎片合并**：FTS5 token overlap > 70% 的相近记忆合并
2. **嵌入更新**：为缺少嵌入的记忆生成向量（Ollama 可用时）
3. **强度提升**：`RecallCount > 5` 的记忆提升 Importance（上限 0.9）
4. **矛盾检测**：扫描实体重叠但内容矛盾的记忆对
5. **嵌入缓存刷新**

### 7.4 检索（Retrieval）

- **主路径**：从预取缓存获取（零延迟）
- **回退路径**：缓存未命中时同步检索（100ms 超时保护）
- **副作用**：命中的记忆更新 `RecallCount++`、`LastRecalledAt`、`RetrievalCount++`

### 7.5 遗忘（Forgetting）

定时任务，应用启动时 + 每日执行一次。

**衰减公式**：

```
decay = e^(-0.03 * days_since_last_recall) * importance * trust_score
```

**遗忘条件**：`decay < 0.1` 且 `RecallCount < 2` 且 `TrustScore < 0.5`

**保护规则**：
- `Importance >= 0.8` 的记忆永不自动遗忘
- `TrustScore >= 0.8` 的记忆永不自动遗忘
- 最近 7 天内创建的记忆不参与遗忘评估
- 遗忘为软删除（`IsForgotten = true`），用户可在前端恢复

---

## 8. 自动触发策略

### 8.1 检索触发（预取 / 同步回退）

使用轻量级规则引擎，不调用 LLM：

```go
func shouldRetrieveMemories(userMessage string) bool {
    // 必须触发的关键词模式
    mustTriggerPatterns := []string{
        "记得", "还记得", "上次", "之前", "昨天", "前天", "上周",
        "上个月", "去年", "我说过", "我提到", "我做过", "我去过",
        "这个月", "最近", "以前", "那次", "那天",
    }
    for _, pattern := range mustTriggerPatterns {
        if strings.Contains(userMessage, pattern) {
            return true
        }
    }

    // 消息长度超过阈值时默认检索
    if utf8.RuneCountInString(userMessage) > 30 {
        return true
    }

    return false
}
```

### 8.2 编码触发（对话后）

```go
func shouldEncodeMemories(userMessage, assistantReply string) bool {
    // 跳过过短的对话
    if utf8.RuneCountInString(userMessage) < 15 {
        return false
    }

    // 包含事件/经历/计划/偏好描述的关键词
    encodePatterns := []string{
        "今天", "昨天", "明天", "计划", "打算", "决定",
        "去了", "买了", "做了", "学了", "吃了", "看了",
        "喜欢", "讨厌", "害怕", "开心", "难过", "生气",
        "工作", "公司", "换了", "搬了", "分手", "结婚",
        "我是", "我在", "我住", "我的",
    }
    for _, pattern := range encodePatterns {
        if strings.Contains(userMessage, pattern) {
            return true
        }
    }

    return false
}
```

### 8.3 触发频率控制

- 预取：每轮对话结束后执行，3 秒超时
- 同步回退：100ms 硬超时，超时则跳过
- 编码：全局 `sync.Mutex` 限制同时 1 个编码任务，多余排队
- 嵌入生成：带缓冲 channel 限流，避免 Ollama 过载

---

## 9. 前端记忆管理 UI

### 9.1 入口

在设置页侧边栏菜单中新增"记忆管理"项，路由键 `memory`。

### 9.2 页面布局

```
┌───────────────────────────────────────────────────────────────┐
│  记忆管理                                                      │
├───────────┬───────────────────────────────────────────────────┤
│           │  ┌─────────────────────────────────────────────┐  │
│  全部 (42) │  │ 搜索记忆...       [时间] [地点] [人物]      │  │
│  事件 (28) │  ├─────────────────────────────────────────────┤  │
│  技能 (5)  │  │                                             │  │
│  计划 (6)  │  │  ┌─────────────────────────────────────┐   │  │
│  已遗忘(3) │  │  │ 首次探访新咖啡馆           2025-04-15 │   │  │
│           │  │  │ 尝试了桂花风味拿铁，体验愉悦...       │   │  │
│           │  │  │ 小区附近 · 信任 0.7 · 重要 0.5       │   │  │
│           │  │  │                    [编辑] [删除]      │   │  │
│           │  │  └─────────────────────────────────────┘   │  │
│           │  │                                             │  │
│           │  │              [加载更多]                      │  │
│           │  └─────────────────────────────────────────────┘  │
├───────────┴───────────────────────────────────────────────────┤
│  共 42 条 | 本周 +5 | 已遗忘 3 | 矛盾 1 待处理              │
└───────────────────────────────────────────────────────────────┘
```

### 9.3 后端 API（Wails 绑定）

| 方法 | 参数 | 说明 |
|------|------|------|
| `GetMemories(query)` | offset, limit, keyword, type, isForgotten | 分页查询（FTS5） |
| `GetMemoryDetail(id)` | memoryID | 详情 + 关联实体 |
| `UpdateMemory(id, data)` | memoryID, 更新字段 | 编辑 |
| `DeleteMemory(id)` | memoryID | 软删除 |
| `RestoreMemory(id)` | memoryID | 恢复 |
| `GetMemoryStats()` | 无 | 统计（总数/本周/已遗忘/矛盾数） |

### 9.4 设置开关

在"通用设置 > 实验室"分区中：
- **启用记忆系统**（总开关）
- **启用向量搜索**（需要 Ollama 可用，默认关闭）

---

## 10. 性能考量

### 10.1 延迟预算

| 操作 | 目标延迟 | 路径 |
|------|----------|------|
| 注入预取缓存 | < 1ms | 同步（内存读取） |
| 同步回退检索（FTS5） | < 50ms | 同步，100ms 硬超时 |
| 同步回退检索（FTS5 + 向量） | < 100ms | 同步，100ms 硬超时 |
| 异步预取 | < 3s | 异步，不阻塞 |
| 记忆编码（LLM 调用） | 2-10s | 异步，不阻塞 |
| Ollama 嵌入 | 200-500ms | 异步 |

### 10.2 内存占用

| 数据 | 估算 | 说明 |
|------|------|------|
| 预取缓存（全部会话） | < 1 MB | 每会话约几 KB 文本 |
| 10,000 条嵌入缓存 | ~30 MB | 768 维 × 4 字节 × 10000 |
| FTS5 索引 | ~原表大小 × 1.5 | SQLite 自动管理 |

### 10.3 并发控制

- 记忆编码：全局 `sync.Mutex`，同时最多 1 个
- 预取：按 chatUuid 去重，同一会话同时最多 1 个预取
- 嵌入生成：缓冲 channel 限流
- 巩固/遗忘定时任务：`sync.Once` 防重复

### 10.4 SQLite 优化

- 启用 WAL 模式：`PRAGMA journal_mode=WAL`，允许读写并发
- 记忆编码批量写入使用事务
- FTS5 索引通过触发器自动维护，无需手动同步

---

## 11. 新增文件清单

```
backend/
├── agents/memory/
│   ├── search/
│   │   ├── hybrid.go          # 混合检索引擎
│   │   ├── scorer.go          # 多维加权打分
│   │   ├── entity.go          # 实体提取与关联检索
│   │   └── prefetch.go        # 预取缓存管理
│   ├── storage/
│   │   ├── storage.go         # (改造) 接受外部 *gorm.DB
│   │   ├── memory.go          # (改造) FTS5 替换 LIKE
│   │   ├── embedding.go       # (新增) 嵌入向量存储
│   │   ├── entity.go          # (新增) 实体和关联 CRUD
│   │   └── fts.go             # (新增) FTS5 索引管理
│   ├── lifecycle/
│   │   ├── consolidator.go    # 记忆巩固
│   │   ├── forgetter.go       # 记忆遗忘
│   │   ├── contradiction.go   # 矛盾检测
│   │   └── trust.go           # 信任分数管理
│   └── models/
│       └── memory.go          # (改造) 新增 TrustScore、实体模型
├── service/
│   ├── memory.go              # (新增) 前端 API 绑定
│   └── chat_completion_runner.go  # (改造) 预取注入 + 异步编码
│
frontend/src/
├── pages/settings/memory/
│   └── index.tsx              # 记忆管理页面
├── stores/
│   └── memoryStore.ts         # 记忆状态管理
```

---

## 12. 分阶段实施

### 阶段一：MVP（核心集成 + FTS5）

**目标**：记忆系统可用，对话中能读写记忆，FTS5 全文检索。

- [ ] 统一存储：Memory 表迁入 `data.db`，重构 Storage 构造函数
- [ ] FTS5 索引：创建虚拟表和同步触发器，替换 LIKE 查询
- [ ] Service 层初始化 memoryStorage 和 memoryAgent
- [ ] `run()` 中添加同步记忆检索 + 围栏标签注入用户消息
- [ ] `run()` 末尾添加异步记忆编码（含工作流隔离）
- [ ] 容错隔离：所有记忆操作 recover 包裹
- [ ] 规则引擎：`shouldRetrieveMemories` / `shouldEncodeMemories`
- [ ] 设置页实验室分区添加"启用记忆系统"开关

### 阶段二：预取缓存 + 信任系统 + 前端 UI

**目标**：零延迟注入，记忆质量控制，用户可管理记忆。

- [ ] 异步预取 + 缓存机制
- [ ] 信任分数：初始值、非对称反馈调整
- [ ] 实体提取和关联存储
- [ ] 后端 Memory CRUD API
- [ ] 前端记忆管理页面（列表、搜索、编辑、删除/恢复）
- [ ] memoryStore（Zustand）
- [ ] RecallCount / LastRecalledAt / TrustScore 更新

### 阶段三：生命周期管理

**目标**：记忆自动巩固、遗忘、矛盾检测。

- [ ] 记忆巩固定时任务（碎片合并、强度提升）
- [ ] 记忆遗忘定时任务（衰减计算、软删除）
- [ ] 矛盾检测（实体重叠 + 内容分歧 → 标记处理）
- [ ] 前端矛盾提示和处理入口

### 阶段四：向量搜索增强

**目标**：语义级别的记忆检索能力。

- [ ] Ollama 嵌入集成（通过 Eino embedding/ollama 组件）
- [ ] `memory_embeddings` 表和存储层
- [ ] 嵌入缓存（内存 + LRU）
- [ ] 混合检索引擎完整版（FTS5 + 实体 + 向量 + 元数据融合）
- [ ] 设置页"启用向量搜索"开关
- [ ] 降级策略（Ollama 不可用时自动回退）
- [ ] 存量记忆的嵌入批量生成（后台任务）

---

## 13. 验证方案

### 功能验证

1. **基础读写**：发送"我今天去了西湖"，下一轮问"我今天做了什么"，验证正确回忆
2. **围栏隔离**：检查 LLM 不会主动提及"我从记忆中看到..."
3. **预取缓存**：第二轮对话延迟应显著低于首轮（对比日志时间戳）
4. **工作流隔离**：触发工作流路径后检查 Memory 表无新增记录
5. **容错**：模拟 Storage 异常，验证主对话流程不受影响
6. **信任过滤**：手动将某记忆 TrustScore 设为 0.2，验证检索不返回
7. **矛盾检测**：写入两条相同人物但矛盾内容的记忆，验证巩固阶段标记

### 性能验证

1. FTS5 检索 10,000 条记忆 < 50ms
2. 预取缓存命中时注入延迟 < 1ms
3. 记忆编码不阻塞用户获取回复

### 边界情况

1. Ollama 不可用时自动降级到 FTS5
2. 记忆数量为 0 时正常运行
3. 并发对话时预取缓存按 chatUuid 隔离
4. 应用重启后定时任务正常恢复

---

## 附录：Hermes Agent 记忆系统调研

> 来源：[github.com/NousResearch/hermes-agent](https://github.com/NousResearch/hermes-agent)
> Hermes Agent 是 Nous Research 开源的 AI Agent 框架，其记忆系统设计成熟且具有多项值得借鉴的设计模式。

### A1. 整体架构：三层记忆

Hermes Agent 采用三层记忆架构，各层职责清晰：

| 层级 | 名称 | 存储方式 | 特点 |
|------|------|----------|------|
| Layer 1 | 内置记忆 | `MEMORY.md` / `USER.md` 磁盘文件 | 始终启用，字符上限硬约束（2200/1375 字符） |
| Layer 2 | 外部记忆提供者 | 插件式，最多 1 个 | 8 种可选（Honcho、Holographic、Mem0 等） |
| Layer 3 | 程序性记忆（Skills） | `~/.hermes/skills/` 目录 | Agent 可自主创建/编辑/删除技能 |

**内置记忆**采用 Markdown 文件存储，操作包括 `add`、`replace`、`remove`，带有重复检测和安全扫描（拦截 prompt 注入、凭据窃取等）。文件级锁保证原子写入。

### A2. MemoryProvider 抽象接口

所有外部记忆提供者实现统一的 ABC 接口 `MemoryProvider`：

```python
class MemoryProvider(ABC):
    def system_prompt_block(self) -> str          # 静态文本注入 system prompt
    def prefetch(self, query, session_id) -> str  # 每轮对话前同步召回
    def queue_prefetch(self, query, session_id)   # 异步预取（下一轮使用）
    def sync_turn(self, user, assistant, session_id)  # 每轮对话后持久化
    def get_tool_schemas(self) -> List[Dict]      # 暴露给 LLM 的工具
    def handle_tool_call(self, tool_name, args)   # 处理 LLM 工具调用

    # 可选生命周期钩子
    def on_turn_start(...)                         # 每轮开始
    def on_session_end(...)                        # 会话结束
    def on_pre_compress(self, messages) -> str     # 上下文压缩前提取洞察
    def on_delegation(self, task, result, ...)     # 子 Agent 委托完成后
    def on_memory_write(self, action, target, content)  # 内置记忆写入时镜像通知
```

**关键设计**：
- `on_pre_compress`：在上下文窗口压缩丢弃旧消息前，提供者有机会提取持久化洞察
- `on_memory_write`：当内置记忆（MEMORY.md/USER.md）发生写入时，外部提供者收到镜像通知，可同步更新自身状态
- 初始化时传入上下文类型（`primary`/`subagent`/`cron`/`flush`），非主上下文跳过写入避免数据污染

### A3. MemoryManager 编排器

`MemoryManager` 是唯一的集成点，管理所有提供者：

```
用户消息到达
    │
    ├─ build_system_prompt()    → 收集所有提供者的静态 prompt 块
    ├─ prefetch_all(message)    → 并行召回，结果用 <memory-context> 标签包裹
    │
    ▼  LLM 生成回复
    │
    ├─ sync_all(user, assistant) → 通知所有提供者持久化本轮
    └─ queue_prefetch_all(msg)   → 异步预取下一轮上下文（零延迟注入）
```

**容错隔离**：每个提供者的调用都被 try/except 包裹，单个提供者异常不影响其他提供者。

### A4. 上下文注入机制

Hermes 的上下文注入有两个值得注意的设计：

#### A4.1 上下文围栏（Context Fencing）

预取的记忆内容被包裹在 XML 标签中注入 **用户消息**（而非 system prompt），并附带系统说明：

```xml
<memory-context>
The following is recalled memory context, NOT new user input.
Treat as informational background data.

[记忆内容...]
</memory-context>
```

**为什么注入到用户消息而非 system prompt**：保持 system prompt 不变可最大化 KV-cache / prefix-cache 命中率，降低推理成本。

#### A4.2 递归清洁（Sanitization）

`sanitize_context()` 函数会从提供者返回的内容中剥离围栏标签、注入的上下文块和系统说明，防止递归注入攻击。

### A5. Holographic 记忆提供者（最具创新性）

这是最具技术特色的记忆提供者，采用 **全息降维表征（Holographic Reduced Representations, HRR）** 实现结构化记忆检索：

#### 存储模型

SQLite 数据库含 5 张表：

| 表名 | 用途 |
|------|------|
| `facts` | 事实内容 + 信任分数 + 检索次数 + HRR 向量 |
| `entities` | 命名实体（人、地、事） |
| `fact_entities` | 事实-实体关联表 |
| `memory_banks` | 按分类聚合的 HRR 向量 |
| `facts_fts` | FTS5 全文检索索引 |

#### HRR 向量操作

- **绑定（Bind）**：通过圆卷积将实体绑定到角色 — `bind(entity_vec, role_entity)`
- **探测（Probe）**：从记忆库中解绑实体，找到结构相关的事实
- **推理（Reason）**：多实体组合查询 — 向量空间 JOIN（取每个实体相似度的最小值作为 AND 语义）
- **矛盾检测（Contradict）**：找到实体重叠度高但内容向量相似度低的事实 — 自动化记忆清理

#### 混合检索流水线

```
查询文本
    │
    ├─ ① FTS5 全文候选集（3× limit 以留出重排空间）
    ├─ ② Jaccard token 重叠度重排
    ├─ ③ HRR 向量相似度
    └─ ④ 信任分数加权
         │
         ▼
    最终得分 = (FTS×0.4 + Jaccard×0.3 + HRR×0.3) × trust_score
              × 可选时间衰减: 0.5^(age_days / half_life)
```

#### 信任系统

- 默认信任分：0.5，范围 0.0-1.0
- **非对称反馈**：有用 +0.05，无用 -0.10（惩罚力度是奖励的 2 倍，偏向准确性）
- 最低信任阈值：0.3，低于此值的事实不参与检索

#### 会话结束时自动提取

使用正则模式从用户消息中检测偏好（"I prefer/like/use..."）和决策（"we decided/agreed..."），自动存储为事实。

### A6. Honcho 用户建模（辩证推理）

Honcho 是最成熟的外部记忆提供者，实现跨会话的 AI 原生用户建模：

#### 双层上下文注入

- **Layer 1（基础上下文）**：会话摘要 + 用户画像 + 用户同伴卡 + AI 自我表征
- **Layer 2（辩证补充）**：多轮 `.chat()` 推理，深度分析用户

#### 辩证深度系统（1-3 轮推理）

| 轮次 | 作用 | 说明 |
|------|------|------|
| Pass 0 | 冷启/热启评估 | 自动检测是否有缓存上下文 |
| Pass 1 | 差距分析与综合 | 条件性跳出（Pass 0 信号足够强时） |
| Pass 2 | 矛盾调和 | 检查并修正冲突 |

#### 成本控制

通过 cadence 参数控制昂贵操作的触发频率：`contextCadence`、`dialecticCadence`、`injectionFrequency`（every-turn / first-turn）、`reasoningLevelCap`。

### A7. 会话搜索（长期情景记忆）

- SQLite + FTS5 索引的会话历史
- 两种模式：列出最近会话（无 LLM）/ 关键词搜索 + LLM 摘要
- 使用辅助模型（Gemini Flash）进行并行异步摘要
- 智能截断：围绕查询词位置取 ~100K 字符窗口

### A8. 值得借鉴的设计模式

| # | 模式 | 说明 | 本方案采纳情况 |
|---|------|------|---------------|
| 1 | **上下文围栏** | XML 标签包裹 + 防注入清洁 | 已采纳：`<memory-context>` 围栏 + `sanitizeMemoryContext` |
| 2 | **注入用户消息** | 保持 system prompt 不变，最大化 KV-cache | 已采纳：注入到用户消息前部 |
| 3 | **异步预取 + 缓存** | 零延迟注入 | 已采纳：`memoryPrefetchCache` 机制 |
| 4 | **容错隔离** | 每个操作独立 recover | 已采纳：`safeMemoryOp` 包裹所有记忆操作 |
| 5 | **非对称信任反馈** | 惩罚 > 奖励 | 已采纳：信任分数系统 |
| 6 | **矛盾检测** | 实体重叠 + 内容分歧 | 已采纳：巩固阶段自动扫描 |
| 7 | **FTS5 全文检索** | 替代 LIKE | 已采纳：FTS5 虚拟表 + 触发器 |
| 8 | **实体结构化存储** | 独立实体表 + 关联表 | 已采纳：MemoryEntity + MemoryEntityLink |
| 9 | **子 Agent 记忆隔离** | 工作流不写入记忆 | 已采纳：工作流路径跳过编码 |
| 10 | **Token 预算控制** | 硬字符上限 | 已采纳：500 tokens 注入预算 |
