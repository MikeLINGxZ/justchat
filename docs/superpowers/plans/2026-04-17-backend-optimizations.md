# Backend Optimizations Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Implement 16 backend optimizations covering workflow efficiency, memory system, tool resolution, concurrency safety, and database performance.

**Architecture:** Optimizations are grouped into 8 independent tasks by file/subsystem affinity. Each task is self-contained and can be committed independently. No new dependencies are introduced.

**Tech Stack:** Go, GORM/SQLite, Eino ADK, MCP-Go

---

### Task 1: Single-Task Workflow Skip Synthesis+Review

**Files:**
- Modify: `backend/service/chat_completion_runner.go:1260-1370`

- [ ] **Step 1: Add single-task fast path before the synthesis/review loop**

In `executeWorkflow()`, after `executeBatches()` (line 1258), insert early exit logic:

```go
// ---- 单任务快速路径：跳过综合+审核 ----
if len(plan.Tasks) == 1 && len(results) == 1 {
    var singleResult workflowTaskResult
    for _, v := range results {
        singleResult = v
    }
    draft := strings.TrimSpace(singleResult.Output)
    if draft == "" {
        draft = "抱歉，子任务未生成有效内容，请重试。"
    }

    r.mu.Lock()
    err := r.startTraceStepLocked("finalize_workflow", "", data_models.TraceStepTypeFinalize, i18n.TCurrent("chat.trace.finalize_title", nil), "单任务直出", draft, "chat.stage.finished", "MainRouterAgent", []data_models.TraceDetailBlock{
        {Kind: "output", Title: i18n.TCurrent("chat.trace.final_draft", nil), Content: draft, Format: data_models.TraceDetailFormatMarkdown},
    }, nil)
    if err == nil {
        r.assistantMessage.Content = draft
        err = r.persistSnapshotLocked(true)
    }
    if err == nil {
        err = r.finishTraceStepLocked("finalize_workflow", "单任务直出", compactText(draft, 240), "chat.stage.finished", "MainRouterAgent", data_models.TraceStepStatusDone, []data_models.TraceDetailBlock{
            {Kind: "output", Title: i18n.TCurrent("chat.trace.final_answer", nil), Content: draft, Format: data_models.TraceDetailFormatMarkdown},
        }, map[string]interface{}{
            "approved": true,
        })
    }
    r.mu.Unlock()
    return err
}
```

Insert this block at line 1259 (right after `executeBatches` returns), **before** the `for attempt := 0; attempt < 2;` loop.

- [ ] **Step 2: Run tests**

Run: `cd /Users/linhuafeng/Work/lemon_tea_desktop && go test ./backend/service/... -v -run TestWorkflow -count=1`
Expected: Existing tests pass (no workflow tests currently break)

- [ ] **Step 3: Build to verify compilation**

Run: `cd /Users/linhuafeng/Work/lemon_tea_desktop && go build ./...`
Expected: BUILD SUCCESS

- [ ] **Step 4: Commit**

```bash
git add backend/service/chat_completion_runner.go
git commit -m "perf(workflow): skip synthesis+review for single-task workflows"
```

---

### Task 2: Batch Results Map Race Condition Fix + Task Manager Concurrency Limit

**Files:**
- Modify: `backend/service/chat_completion_runner.go:1398-1500` (executeBatches)
- Modify: `backend/pkg/tasker/manager.go`

- [ ] **Step 1: Fix race condition in executeBatches — protect results map writes**

In `executeBatches()` around line 1398, add a mutex for results and wrap the write at line ~1497:

Find the section where batch results are collected (after the WaitGroup completes, results are written). The pattern is:

```go
// In the goroutine that runs each task:
// After executePlanTask returns, write result to shared map
```

Add a `resultsMu` local to `executeBatches`:

```go
func (r *completionRunner) executeBatches(runCtx context.Context, plan workflowPlan, results map[string]workflowTaskResult, filterIDs map[string]struct{}, retryInstructions string, originalUserMessage *schema.Message) error {
	batches := batchTasksByDependencies(plan.Tasks, filterIDs)
	toolMiddleware := r.buildToolMiddleware()

	r.mu.Lock()
	retryCount := r.assistantMessage.AssistantMessageExtra.RetryCount
	r.mu.Unlock()

	var resultsMu sync.Mutex // 保护 results map 的并发写入
```

Then wherever `results[item.task.ID] = item.result` is written inside a goroutine, wrap it:

```go
resultsMu.Lock()
results[item.task.ID] = item.result
resultsMu.Unlock()
```

- [ ] **Step 2: Add concurrency limit to Task Manager**

Replace `backend/pkg/tasker/manager.go` content:

```go
package tasker

import "sync"

const defaultMaxConcurrent = 10

// Manager 管理按 taskUuid 区分的后台任务。
var Manager *manager

func init() {
	Manager = &manager{
		tasks:     make(map[string]*Runtime),
		semaphore: make(chan struct{}, defaultMaxConcurrent),
	}
}

type Runtime struct {
	TaskUUID             string
	ChatUUID             string
	AssistantMessageUUID string
	EventKey             string
	stopCh               chan struct{}
}

type manager struct {
	mu        sync.Mutex
	tasks     map[string]*Runtime
	semaphore chan struct{}
}

func (m *manager) StartTask(task Runtime, fn func(stop <-chan struct{})) {
	userStop := make(chan struct{}, 1)
	task.stopCh = userStop

	m.mu.Lock()
	m.tasks[task.TaskUUID] = &task
	m.mu.Unlock()

	go func() {
		// 获取并发槽位
		m.semaphore <- struct{}{}
		defer func() {
			<-m.semaphore // 释放槽位
			m.mu.Lock()
			delete(m.tasks, task.TaskUUID)
			m.mu.Unlock()
			close(userStop)
		}()
		fn(userStop)
	}()
}

func (m *manager) StopTask(taskUUID string) {
	m.mu.Lock()
	task, ok := m.tasks[taskUUID]
	m.mu.Unlock()
	if !ok || task == nil || task.stopCh == nil {
		return
	}
	select {
	case task.stopCh <- struct{}{}:
	default:
	}
}

func (m *manager) StopByEventKey(eventKey string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, task := range m.tasks {
		if task == nil || task.EventKey != eventKey || task.stopCh == nil {
			continue
		}
		select {
		case task.stopCh <- struct{}{}:
		default:
		}
		return
	}
}

func (m *manager) GetTaskRuntime(taskUUID string) (*Runtime, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	task, ok := m.tasks[taskUUID]
	if !ok || task == nil {
		return nil, false
	}
	copyTask := *task
	copyTask.stopCh = nil
	return &copyTask, true
}

func (m *manager) ListRunningTasks() []Runtime {
	m.mu.Lock()
	defer m.mu.Unlock()
	res := make([]Runtime, 0, len(m.tasks))
	for _, task := range m.tasks {
		if task == nil {
			continue
		}
		copyTask := *task
		copyTask.stopCh = nil
		res = append(res, copyTask)
	}
	return res
}
```

- [ ] **Step 3: Build and run tests**

Run: `cd /Users/linhuafeng/Work/lemon_tea_desktop && go build ./... && go test ./backend/pkg/tasker/... -v -count=1`
Expected: BUILD SUCCESS, tests pass

- [ ] **Step 4: Commit**

```bash
git add backend/service/chat_completion_runner.go backend/pkg/tasker/manager.go
git commit -m "fix(concurrency): protect batch results map + add task manager concurrency limit"
```

---

### Task 3: Database Composite Indexes + Merge Count Query

**Files:**
- Modify: `backend/models/data_models/message.go:13-25`
- Modify: `backend/agents/memory/models/memory.go:19-37`
- Modify: `backend/storage/chat_message.go:54-67`

- [ ] **Step 1: Add composite index to Message model**

In `backend/models/data_models/message.go`, modify the `ChatUuid` field tag to include a composite index:

```go
type Message struct {
	OrmModel
	ChatUuid                     string                 `gorm:"index:idx_msg_chat_created,priority:1" json:"chat_uuid"`
	MessageUuid                  string                 `grom:"unique;index" json:"message_uuid"`
	Role                         schema.RoleType        `json:"role"`
```

Note: Remove the standalone `gorm:"index"` from `Role` — it's never queried alone and wastes write performance.

- [ ] **Step 2: Add composite index to Memory model**

In `backend/agents/memory/models/memory.go`, add composite index for `IsForgotten + CreatedAt`:

```go
type Memory struct {
	data_models.OrmModel
	Summary       string     `gorm:"type:varchar(500)"`
	Content       string     `gorm:"type:text;not null"`
	Type          MemoryType `gorm:"type:varchar(50);index"`
	TimeRangStart *time.Time `gorm:"index"`
	TimeRangeEnd  *time.Time `gorm:"index"`
	Location      *string    `gorm:"type:varchar(500)"`
	Characters    *string    `gorm:"type:varchar(500)"`
	Context       *string    `gorm:"type:json"`

	EmbeddingID      *uint      `gorm:"index"`
	Importance       float64    `gorm:"default:0.5;index"`
	EmotionalValence float64    `gorm:"default:0.0"`
	TrustScore       float64    `gorm:"default:0.5"`
	IsForgotten      bool       `gorm:"default:false;index:idx_mem_forgotten_created,priority:1"`
	RecallCount      int        `gorm:"default:0"`
	LastRecalledAt   *time.Time `gorm:"column:last_recalled_at"`
}
```

Also add `CreatedAt` composite participation. Since `OrmModel` embeds `CreatedAt`, GORM AutoMigrate will pick up the index from the `IsForgotten` tag. Add a standalone index on `Importance` for sort queries.

- [ ] **Step 3: Merge GetMessage count into single query path**

Replace `GetMessage` in `backend/storage/chat_message.go`:

```go
func (s *Storage) GetMessage(ctx context.Context, chatUuid string, offset, limit int) ([]data_models.Message, int, error) {
	var messages []data_models.Message
	var total int64

	db := s.sqliteDB.Model(&data_models.Message{}).Where("chat_uuid = ?", chatUuid)

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := db.Order("created_at asc").Offset(offset).Limit(limit).Find(&messages).Error; err != nil {
		return nil, 0, err
	}

	return messages, int(total), nil
}
```

This reuses the same `db` query builder, avoiding redundant WHERE clause evaluation.

- [ ] **Step 4: Build and run tests**

Run: `cd /Users/linhuafeng/Work/lemon_tea_desktop && go build ./... && go test ./backend/... -v -count=1`
Expected: BUILD SUCCESS, all tests pass

- [ ] **Step 5: Commit**

```bash
git add backend/models/data_models/message.go backend/agents/memory/models/memory.go backend/storage/chat_message.go
git commit -m "perf(db): add composite indexes + merge count query"
```

---

### Task 4: Agent Registry O(1) Lookup

**Files:**
- Modify: `backend/pkg/llm_provider/agents/agent_registry.go`

- [ ] **Step 1: Add name index map and update all functions**

Replace `backend/pkg/llm_provider/agents/agent_registry.go`:

```go
package agents

import "sync"

var (
	mu        sync.RWMutex
	registry  []IAgent
	nameIndex map[string]IAgent
)

func init() {
	nameIndex = make(map[string]IAgent)
}

// RegisterAgent 注册一个 Agent 到全局注册表（通常在 init() 中调用）。
func RegisterAgent(a IAgent) {
	mu.Lock()
	defer mu.Unlock()
	registry = append(registry, a)
	nameIndex[a.Name()] = a
}

// AllAgents 返回所有已注册的 Agent。
func AllAgents() []IAgent {
	mu.RLock()
	defer mu.RUnlock()
	result := make([]IAgent, len(registry))
	copy(result, registry)
	return result
}

// FindAgent 按名称查找 Agent（O(1)）。
func FindAgent(name string) (IAgent, bool) {
	mu.RLock()
	defer mu.RUnlock()
	a, ok := nameIndex[name]
	return a, ok
}

// AgentOwnedPromptNames 返回所有 Agent 拥有的提示词名称集合。
func AgentOwnedPromptNames() map[string]struct{} {
	mu.RLock()
	defer mu.RUnlock()
	owned := make(map[string]struct{})
	for _, a := range registry {
		for _, name := range a.PromptNames() {
			owned[name] = struct{}{}
		}
	}
	return owned
}

// UnregisterAgentsByType removes all agents of the given type from the registry.
func UnregisterAgentsByType(agentType AgentType) {
	mu.Lock()
	defer mu.Unlock()
	filtered := make([]IAgent, 0, len(registry))
	for _, a := range registry {
		if a.Type() != agentType {
			filtered = append(filtered, a)
		} else {
			delete(nameIndex, a.Name())
		}
	}
	registry = filtered
}

// UnregisterAgentByName removes an agent by its name from the registry.
func UnregisterAgentByName(name string) {
	mu.Lock()
	defer mu.Unlock()
	filtered := make([]IAgent, 0, len(registry))
	for _, a := range registry {
		if a.Name() != name {
			filtered = append(filtered, a)
		}
	}
	registry = filtered
	delete(nameIndex, name)
}

// AgentsByType 按类型过滤 Agent。
func AgentsByType(agentType AgentType) []IAgent {
	mu.RLock()
	defer mu.RUnlock()
	var result []IAgent
	for _, a := range registry {
		if a.Type() == agentType {
			result = append(result, a)
		}
	}
	return result
}

// AgentsByRole 按角色过滤 Agent。
func AgentsByRole(role AgentRole) []IAgent {
	mu.RLock()
	defer mu.RUnlock()
	var result []IAgent
	for _, a := range registry {
		if a.Role() == role {
			result = append(result, a)
		}
	}
	return result
}
```

- [ ] **Step 2: Build and run tests**

Run: `cd /Users/linhuafeng/Work/lemon_tea_desktop && go build ./... && go test ./backend/... -v -count=1`
Expected: BUILD SUCCESS, all tests pass

- [ ] **Step 3: Commit**

```bash
git add backend/pkg/llm_provider/agents/agent_registry.go
git commit -m "perf(agents): use map index for O(1) FindAgent lookup"
```

---

### Task 5: Memory System — Keyword Stopwords + Prefetch Cache Bound + Encoding Lock

**Files:**
- Modify: `backend/service/memory.go`

- [ ] **Step 1: Add Chinese/English stopwords filter to extractKeywords**

After `extractKeywords` splits segments (around line 162), add stopword filtering before the dedup step:

```go
// 停用词过滤（高频无意义词）
var stopWords = map[string]bool{
	// 中文
	"的": true, "了": true, "是": true, "在": true, "我": true, "有": true,
	"和": true, "就": true, "不": true, "人": true, "都": true, "一": true,
	"一个": true, "上": true, "也": true, "很": true, "到": true, "说": true,
	"要": true, "去": true, "你": true, "会": true, "着": true, "没有": true,
	"看": true, "好": true, "自己": true, "这": true, "他": true, "她": true,
	"那": true, "它": true, "吗": true, "吧": true, "呢": true, "啊": true,
	"嗯": true, "哦": true, "哈": true, "呀": true, "把": true, "被": true,
	"让": true, "给": true, "从": true, "对": true, "但": true, "而": true,
	"还": true, "这个": true, "那个": true, "什么": true, "怎么": true, "可以": true,
	// 英文
	"the": true, "a": true, "an": true, "is": true, "are": true, "was": true,
	"were": true, "be": true, "been": true, "being": true, "have": true, "has": true,
	"had": true, "do": true, "does": true, "did": true, "will": true, "would": true,
	"could": true, "should": true, "may": true, "might": true, "shall": true,
	"can": true, "to": true, "of": true, "in": true, "for": true, "on": true,
	"with": true, "at": true, "by": true, "from": true, "it": true, "this": true,
	"that": true, "and": true, "or": true, "but": true, "not": true, "no": true,
	"i": true, "me": true, "my": true, "you": true, "your": true, "he": true,
	"she": true, "we": true, "they": true,
}
```

Declare this as a package-level `var`. Then in `extractKeywords`, after the segment loop and before dedup:

```go
var keywords []string
for _, seg := range segments {
    seg = strings.TrimSpace(seg)
    if utf8.RuneCountInString(seg) >= 2 && !stopWords[strings.ToLower(seg)] {
        keywords = append(keywords, seg)
    }
}
```

- [ ] **Step 2: Add max size limit to prefetch cache**

Add a `maxSize` constant and eviction to `memoryPrefetchCache.Set`:

```go
const prefetchCacheMaxSize = 100

func (c *memoryPrefetchCache) Set(chatUuid, memCtx string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	// 超过容量时清除最旧的条目
	if len(c.cache) >= prefetchCacheMaxSize {
		var oldestKey string
		var oldestTime time.Time
		for k, v := range c.cache {
			if oldestKey == "" || v.fetchedAt.Before(oldestTime) {
				oldestKey = k
				oldestTime = v.fetchedAt
			}
		}
		if oldestKey != "" {
			delete(c.cache, oldestKey)
		}
	}
	c.cache[chatUuid] = prefetchEntry{context: memCtx, fetchedAt: time.Now()}
}
```

- [ ] **Step 3: Replace global mutex with TryLock + channel queue for encoding**

The current `memoryEncodeMu` (line 19) uses `TryLock` which silently drops encoding if locked. Replace with a buffered channel acting as a queue:

```go
// 记忆编码队列（容量 3，超出丢弃）
var memoryEncodeQueue = make(chan struct{}, 3)
```

Then in `encodeMemoriesAsync`, replace the `TryLock` pattern:

```go
select {
case memoryEncodeQueue <- struct{}{}:
    defer func() { <-memoryEncodeQueue }()
default:
    // 队列已满，跳过本次编码
    return
}
```

Remove the `memoryEncodeMu` variable.

- [ ] **Step 4: Build and run tests**

Run: `cd /Users/linhuafeng/Work/lemon_tea_desktop && go build ./... && go test ./backend/service/... -v -count=1`
Expected: BUILD SUCCESS

- [ ] **Step 5: Commit**

```bash
git add backend/service/memory.go
git commit -m "perf(memory): add stopwords filter, bound prefetch cache, queue-based encoding"
```

---

### Task 6: Hybrid Search — Incremental Cache + GetMemoryByID Batch

**Files:**
- Modify: `backend/agents/memory/search/hybrid.go`

- [ ] **Step 1: Add incremental cache refresh**

Add a `lastCacheID` field and change `refreshCache` to only load new entries:

```go
type HybridSearcher struct {
	storage  *storage.Storage
	embedder embedding.Embedder

	cacheMu     sync.RWMutex
	embCache    []storage.EmbeddingEntry
	cacheReady  bool
	lastCacheID uint // 上次缓存的最大 memory ID
}
```

Update `refreshCache`:

```go
func (hs *HybridSearcher) refreshCache() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	hs.cacheMu.RLock()
	lastID := hs.lastCacheID
	hs.cacheMu.RUnlock()

	var entries []storage.EmbeddingEntry
	var err error
	if lastID == 0 {
		// 首次加载：全量
		entries, err = hs.storage.LoadAllEmbeddings(ctx)
	} else {
		// 增量加载：只取新增
		entries, err = hs.storage.LoadEmbeddingsSince(ctx, lastID)
	}
	if err != nil {
		logger.Error("refresh embedding cache error:", err)
		return
	}

	hs.cacheMu.Lock()
	if lastID == 0 {
		hs.embCache = entries
	} else {
		hs.embCache = append(hs.embCache, entries...)
	}
	// 更新 lastCacheID
	for _, e := range entries {
		if e.MemoryID > hs.lastCacheID {
			hs.lastCacheID = e.MemoryID
		}
	}
	hs.cacheReady = true
	hs.cacheMu.Unlock()
}
```

- [ ] **Step 2: Add LoadEmbeddingsSince to storage**

In `backend/agents/memory/storage/embedding.go`, add:

```go
// LoadEmbeddingsSince 增量加载指定 ID 之后的嵌入。
func (s *Storage) LoadEmbeddingsSince(ctx context.Context, sinceMemoryID uint) ([]EmbeddingEntry, error) {
	var rows []MemoryEmbedding
	if err := s.sqliteDb.WithContext(ctx).
		Joins("JOIN memories ON memories.id = memory_embeddings.memory_id AND memories.is_forgotten = 0").
		Where("memory_embeddings.memory_id > ?", sinceMemoryID).
		Find(&rows).Error; err != nil {
		return nil, err
	}
	entries := make([]EmbeddingEntry, 0, len(rows))
	for _, row := range rows {
		vec := bytesToFloat32s(row.Vector)
		if len(vec) == 0 {
			continue
		}
		entries = append(entries, EmbeddingEntry{MemoryID: row.MemoryID, Vector: vec})
	}
	return entries, nil
}
```

- [ ] **Step 3: Batch GetMemoryByID in Search**

In `hybrid.go` `Search()`, replace the per-ID `GetMemoryByID` loop (lines 89-93) with a batch query:

```go
// 批量加载记忆详情
allIDSlice := make([]uint, 0, len(allIDs))
for id := range allIDs {
    allIDSlice = append(allIDSlice, id)
}
memoryMap, loadErr := hs.storage.GetMemoriesByIDs(ctx, allIDSlice)
if loadErr != nil {
    logger.Error("hybrid search batch load error:", loadErr)
    return nil
}

var scored []ScoredMemory
now := time.Now()
for id := range allIDs {
    m, ok := memoryMap[id]
    if !ok {
        continue
    }
    if m.TrustScore < minTrustThreshold {
        continue
    }
    // ... scoring logic unchanged
}
```

- [ ] **Step 4: Add GetMemoriesByIDs to storage**

In `backend/agents/memory/storage/memory.go` (or appropriate file), add:

```go
// GetMemoriesByIDs 批量查询记忆。
func (s *Storage) GetMemoriesByIDs(ctx context.Context, ids []uint) (map[uint]*models.Memory, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var memories []models.Memory
	if err := s.sqliteDb.WithContext(ctx).Where("id IN ?", ids).Find(&memories).Error; err != nil {
		return nil, err
	}
	result := make(map[uint]*models.Memory, len(memories))
	for i := range memories {
		result[memories[i].ID] = &memories[i]
	}
	return result, nil
}
```

- [ ] **Step 5: Build and run tests**

Run: `cd /Users/linhuafeng/Work/lemon_tea_desktop && go build ./... && go test ./backend/... -v -count=1`
Expected: BUILD SUCCESS

- [ ] **Step 6: Commit**

```bash
git add backend/agents/memory/search/hybrid.go backend/agents/memory/storage/embedding.go backend/agents/memory/storage/memory.go
git commit -m "perf(memory): incremental embedding cache + batch memory loading"
```

---

### Task 7: MCP Tool Resolution — Group by Server + Configurable Timeouts

**Files:**
- Modify: `backend/service/mcp.go:199-255`
- Modify: `backend/plugin/rpc.go:70-115`

- [ ] **Step 1: Group MCP tools by server in resolveSelectedTools**

Replace the per-tool loop in `resolveSelectedTools` (lines 210-252) with a grouped approach:

```go
func (s *Service) resolveSelectedTools(ctx context.Context, toolIDs []string) ([]einotool.BaseTool, map[string]toolMeta, func(), error) {
	var result []einotool.BaseTool
	metaMap := make(map[string]toolMeta)
	var cleanupFns []func()

	cleanup := func() {
		for _, fn := range cleanupFns {
			fn()
		}
	}

	// 按 server 分组非内置工具
	type mcpGroup struct {
		server *data_models.CustomMCPServer
		config *mcpServerConfig
	}
	mcpGroups := make(map[string]*mcpGroup) // serverToolID -> group

	for _, toolID := range toolIDs {
		if builtinTool, ok := llmtools.ToolRouter.GetToolByID(toolID); ok {
			result = append(result, builtinTool.Tool())
			metaMap[toolID] = toolMeta{
				ID:          builtinTool.Id(),
				Name:        builtinTool.Name(),
				Description: builtinTool.Description(),
			}
			continue
		}

		server, err := s.storage.GetCustomMCPServerByToolID(ctx, toolID)
		if err != nil {
			cleanup()
			return nil, nil, nil, err
		}
		if server == nil {
			cleanup()
			return nil, nil, nil, fmt.Errorf("tool %s not found", toolID)
		}
		if !server.Enabled {
			continue
		}

		if _, exists := mcpGroups[server.ToolID]; !exists {
			_, _, _, serverConfig, err := s.parseMCPFolder(server.SourcePath)
			if err != nil {
				cleanup()
				return nil, nil, nil, err
			}
			mcpGroups[server.ToolID] = &mcpGroup{server: server, config: serverConfig}
		}
	}

	// 每个 server 只启动一次
	for _, group := range mcpGroups {
		tools, serverMeta, closeFn, err := s.loadMCPServerTools(ctx, *group.server, group.config)
		if err != nil {
			cleanup()
			return nil, nil, nil, err
		}
		if closeFn != nil {
			cleanupFns = append(cleanupFns, closeFn)
		}
		result = append(result, tools...)
		for k, v := range serverMeta {
			metaMap[k] = v
		}
	}

	return result, metaMap, cleanup, nil
}
```

- [ ] **Step 2: Add context timeout to MCP initialization**

In `loadMCPServerTools` (line 257), add timeout context for Initialize:

```go
func (s *Service) loadMCPServerTools(ctx context.Context, server data_models.CustomMCPServer, configOverride *mcpServerConfig) ([]einotool.BaseTool, map[string]toolMeta, func(), error) {
	config := configOverride
	if config == nil {
		_, _, _, parsedConfig, err := s.parseMCPFolder(server.SourcePath)
		if err != nil {
			return nil, nil, nil, err
		}
		config = parsedConfig
	}

	mcpClient, err := client.NewStdioMCPClient(config.Command, buildMCPEnv(config.Env), config.Args...)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("start mcp server failed: %w", err)
	}

	closeFn := func() {
		_ = mcpClient.Close()
	}

	// 初始化使用超时上下文
	initCtx, initCancel := context.WithTimeout(ctx, 30*time.Second)
	defer initCancel()

	initReq := mcpproto.InitializeRequest{}
	initReq.Params.ProtocolVersion = mcpproto.LATEST_PROTOCOL_VERSION
	initReq.Params.ClientInfo = mcpproto.Implementation{
		Name:    "lemon-tea-desktop",
		Version: "1.0.0",
	}
	if _, err := mcpClient.Initialize(initCtx, initReq); err != nil {
		closeFn()
		return nil, nil, nil, fmt.Errorf("initialize mcp server failed: %w", err)
	}
	// ... rest unchanged
```

- [ ] **Step 3: Make RPC timeout configurable per method**

In `backend/plugin/rpc.go`, replace hardcoded 30s with method-aware timeout:

```go
var rpcTimeouts = map[string]time.Duration{
	"tool/execute":      60 * time.Second,
	"hook/onBeforeChat": 10 * time.Second,
	"hook/onAfterChat":  10 * time.Second,
}

const defaultRPCTimeout = 30 * time.Second

func rpcTimeoutFor(method string) time.Duration {
	if t, ok := rpcTimeouts[method]; ok {
		return t
	}
	return defaultRPCTimeout
}
```

Then in `Call()` method, replace `time.After(30 * time.Second)` with:

```go
case <-time.After(rpcTimeoutFor(method)):
    c.mu.Lock()
    delete(c.pending, id)
    c.mu.Unlock()
    return nil, fmt.Errorf("rpc call %s timed out after %v", method, rpcTimeoutFor(method))
```

- [ ] **Step 4: Build and run tests**

Run: `cd /Users/linhuafeng/Work/lemon_tea_desktop && go build ./... && go test ./backend/... -v -count=1`
Expected: BUILD SUCCESS

- [ ] **Step 5: Commit**

```bash
git add backend/service/mcp.go backend/plugin/rpc.go
git commit -m "perf(mcp): group tools by server + add init timeout + configurable RPC timeouts"
```

---

### Task 8: Workflow JSON Error Handling + Hybrid Search Configurable Weights

**Files:**
- Modify: `backend/service/chat_orchestration.go:86-133, 376-393, 497-520`
- Modify: `backend/agents/memory/search/hybrid.go:107-117`

- [ ] **Step 1: Improve workflow JSON parsing with retry and logging**

Replace `unmarshalJSONResponse` in `chat_orchestration.go`:

```go
func unmarshalJSONResponse(raw string, target interface{}) error {
	raw = strings.TrimSpace(raw)
	// 尝试提取 JSON 块（处理 ```json 包裹）
	raw = strings.TrimPrefix(raw, "```json")
	raw = strings.TrimPrefix(raw, "```JSON")
	raw = strings.TrimPrefix(raw, "```")
	raw = strings.TrimSuffix(raw, "```")
	raw = strings.TrimSpace(raw)

	// 尝试匹配最外层 { ... } 或 [ ... ]
	if idx := strings.IndexAny(raw, "{["); idx > 0 {
		raw = raw[idx:]
	}
	if lastBrace := strings.LastIndex(raw, "}"); lastBrace >= 0 {
		if lastBracket := strings.LastIndex(raw, "]"); lastBracket > lastBrace {
			raw = raw[:lastBracket+1]
		} else {
			raw = raw[:lastBrace+1]
		}
	}

	if err := json.Unmarshal([]byte(raw), target); err != nil {
		logger.Error("workflow JSON parse failed",
			"raw_prefix", compactText(raw, 500),
			"error", err)
		return fmt.Errorf("invalid JSON response: %w", err)
	}
	return nil
}
```

- [ ] **Step 2: Add JSON retry wrapper for generateWorkflowPlan**

In `generateWorkflowPlan`, wrap the Generate+Unmarshal with a single retry:

```go
func generateWorkflowPlan(ctx context.Context, provider *llm_provider.Provider, userRequest string, messages []schema.Message, tools []tool.BaseTool) (workflowPlan, error) {
	toolNames := make([]string, 0, len(tools))
	for _, item := range tools {
		info, infoErr := item.Info(ctx)
		if infoErr != nil || info == nil {
			continue
		}
		toolNames = append(toolNames, info.Name)
	}
	sort.Strings(toolNames)
	systemPrompt := prompts.Render(provider.Prompts().PlannerSystem, map[string]string{
		"tool_names": strings.Join(toolNames, ", "),
	})

	plannerMessages := append([]schema.Message{{Role: schema.System, Content: systemPrompt}}, buildPlannerMessages(provider.Prompts().PlanningUser, userRequest, messages)...)

	var plan workflowPlan
	var lastErr error
	for attempt := 0; attempt < 2; attempt++ {
		resp, err := provider.Generate(ctx, plannerMessages)
		if err != nil {
			return workflowPlan{}, err
		}
		if err := unmarshalJSONResponse(resp.Content, &plan); err != nil {
			lastErr = err
			// 第一次失败：追加纠正消息再试一次
			if attempt == 0 {
				plannerMessages = append(plannerMessages,
					schema.Message{Role: schema.Assistant, Content: resp.Content},
					schema.Message{Role: schema.User, Content: "你的输出格式不正确，请严格按照 JSON 格式输出，不要包含任何额外文字。"},
				)
				continue
			}
			return workflowPlan{}, lastErr
		}
		break
	}

	// ... existing validation logic unchanged
	if len(plan.Tasks) == 0 {
		plan.Tasks = []workflowPlanTask{{
			ID:             "task_1",
			Title:          "完成用户请求",
			Description:    userRequest,
			SuggestedAgent: "GeneralWorkerAgent",
			ExpectedOutput: "完整回答用户请求",
		}}
	}
	for idx := range plan.Tasks {
		if strings.TrimSpace(plan.Tasks[idx].ID) == "" {
			plan.Tasks[idx].ID = fmt.Sprintf("task_%d", idx+1)
		}
		if strings.TrimSpace(plan.Tasks[idx].SuggestedAgent) == "" {
			plan.Tasks[idx].SuggestedAgent = "GeneralWorkerAgent"
		}
	}
	if strings.TrimSpace(plan.Goal) == "" {
		plan.Goal = userRequest
	}
	if len(plan.CompletionCriteria) == 0 {
		plan.CompletionCriteria = []string{"回答覆盖用户目标", "内容清晰可执行"}
	}
	return plan, nil
}
```

Apply the same retry pattern to `reviewWorkflowAnswer` — wrap its `unmarshalJSONResponse` call with a single retry.

- [ ] **Step 3: Make hybrid search weights configurable**

In `backend/agents/memory/search/hybrid.go`, add configurable weights:

```go
// SearchWeights 混合检索权重配置。
type SearchWeights struct {
	FTS        float64
	Vector     float64
	Recency    float64
	Importance float64
}

// DefaultWeights 默认权重（向量可用时）。
var DefaultWeights = SearchWeights{FTS: 0.30, Vector: 0.30, Recency: 0.10, Importance: 0.15}

// FTSOnlyWeights 仅 FTS 时的权重。
var FTSOnlyWeights = SearchWeights{FTS: 0.45, Vector: 0.0, Recency: 0.15, Importance: 0.20}
```

Add `weights` field to `HybridSearcher`:

```go
type HybridSearcher struct {
	storage  *storage.Storage
	embedder embedding.Embedder
	weights        SearchWeights
	ftsOnlyWeights SearchWeights

	cacheMu     sync.RWMutex
	embCache    []storage.EmbeddingEntry
	cacheReady  bool
	lastCacheID uint
}
```

Update `NewHybridSearcher`:

```go
func NewHybridSearcher(s *storage.Storage, embedder embedding.Embedder) *HybridSearcher {
	hs := &HybridSearcher{
		storage:        s,
		embedder:       embedder,
		weights:        DefaultWeights,
		ftsOnlyWeights: FTSOnlyWeights,
	}
	if embedder != nil {
		go hs.refreshCache()
	}
	return hs
}
```

Update scoring in `Search()`:

```go
if vectorEnabled {
    w := hs.weights
    sm.FinalScore = (w.FTS*sm.FTSScore +
        w.Vector*sm.VecScore +
        w.Recency*sm.Recency +
        w.Importance*m.Importance) * m.TrustScore
} else {
    w := hs.ftsOnlyWeights
    sm.FinalScore = (w.FTS*sm.FTSScore +
        w.Recency*sm.Recency +
        w.Importance*m.Importance) * m.TrustScore
}
```

- [ ] **Step 4: Build and run tests**

Run: `cd /Users/linhuafeng/Work/lemon_tea_desktop && go build ./... && go test ./backend/... -v -count=1`
Expected: BUILD SUCCESS, existing orchestration tests pass

- [ ] **Step 5: Commit**

```bash
git add backend/service/chat_orchestration.go backend/agents/memory/search/hybrid.go
git commit -m "fix(workflow): retry on malformed JSON + configurable search weights"
```
