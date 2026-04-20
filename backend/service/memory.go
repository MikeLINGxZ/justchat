package service

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/cloudwego/eino/schema"
	memory_agents "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/agents/memory"
	memory_models "gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/agents/memory/models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/models/wrapper_models"
	"gitlab.linhf.cn/project/lemontea/lemon_tea_desktop/backend/pkg/logger"
)

// memoryEncodeMu 全局限制同时只有 1 个记忆编码任务运行。
var memoryEncodeMu sync.Mutex

// ---- 预取缓存 ----

// memoryPrefetchCache 按 chatUuid 缓存异步预取的记忆上下文。
type memoryPrefetchCache struct {
	mu    sync.Mutex
	cache map[string]prefetchEntry
}

type prefetchEntry struct {
	context   string
	fetchedAt time.Time
}

func newMemoryPrefetchCache() *memoryPrefetchCache {
	return &memoryPrefetchCache{cache: make(map[string]prefetchEntry)}
}

// Get 获取缓存（命中后清除），超过 5 分钟视为过期。
func (c *memoryPrefetchCache) Get(chatUuid string) (string, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.cache[chatUuid]
	if !ok {
		return "", false
	}
	delete(c.cache, chatUuid)
	if time.Since(entry.fetchedAt) > 5*time.Minute {
		return "", false
	}
	return entry.context, true
}

// Set 写入预取缓存。
func (c *memoryPrefetchCache) Set(chatUuid, memCtx string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[chatUuid] = prefetchEntry{context: memCtx, fetchedAt: time.Now()}
}

// Cleanup 清除过期的缓存条目。
func (c *memoryPrefetchCache) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for k, v := range c.cache {
		if time.Since(v.fetchedAt) > 5*time.Minute {
			delete(c.cache, k)
		}
	}
}

// ---- 触发策略 ----
// 不使用硬编码关键词判断，避免语言依赖和误判。
// 检索：始终尝试（成本极低）。
// 编码：始终交给 Memory Agent 自主决定是否调用 write_memory（Agent prompt 有完整策略指导）。

// shouldRetrieveMemories 判断是否需要检索记忆（对话前）。
// 仅过滤空消息和极短的无意义输入。
func shouldRetrieveMemories(userMessage string) bool {
	return utf8.RuneCountInString(strings.TrimSpace(userMessage)) > 2
}

// shouldEncodeMemories 判断是否值得交给 Memory Agent 分析。
// 仅过滤极短消息，具体是否写入由 Agent 自主决定。
func shouldEncodeMemories(userMessage string) bool {
	userMessage = sanitizeMemoryContext(userMessage)
	return utf8.RuneCountInString(strings.TrimSpace(userMessage)) > 4
}

// ---- 记忆检索 ----

// retrieveMemoryContext 根据用户消息检索相关记忆，返回格式化的上下文文本。
// 优先使用混合检索引擎（含向量），降级为纯 FTS5/LIKE。
func (s *Service) retrieveMemoryContext(ctx context.Context, userMessage string) string {
	if s.memoryStorage == nil {
		return ""
	}
	if !shouldRetrieveMemories(userMessage) {
		return ""
	}

	keywords := extractKeywords(userMessage)
	if len(keywords) == 0 {
		return ""
	}

	var memories []memory_models.Memory

	// 优先使用混合检索引擎
	if s.memorySearcher != nil {
		scored := s.memorySearcher.Search(ctx, keywords, userMessage, 5)
		for _, sm := range scored {
			memories = append(memories, sm.Memory)
		}
	} else {
		// 降级为纯 FTS5/LIKE
		var err error
		memories, err = s.memoryStorage.FTSSearch(ctx, keywords, 5)
		if err != nil {
			logger.Error("memory FTS search error", err)
			return ""
		}
	}

	// 检索无结果时，兜底返回最重要的记忆。
	// 处理 "我是谁"、"你知道我什么" 等查询——关键词与记忆内容无文本重叠。
	if len(memories) == 0 {
		fallback, err := s.memoryStorage.TopImportantMemories(ctx, 3)
		if err != nil {
			logger.Error("memory fallback query error", err)
			return ""
		}
		memories = fallback
	}

	if len(memories) == 0 {
		return ""
	}

	// 更新召回计数（异步，不阻塞）
	var ids []uint
	for _, m := range memories {
		ids = append(ids, m.ID)
	}
	go func() {
		if err := s.memoryStorage.IncrementRecallCount(context.Background(), ids); err != nil {
			logger.Error("increment recall count error", err)
		}
	}()

	return formatMemoryContext(memories)
}

// retrieveViaMemoryAgent 让 Memory Agent 对用户消息进行"代理驱动检索"：
// 由 Agent 自主调用 read_memory / get_current_time 等工具，并返回已经格式化好的记忆列表文本。
// 失败或超时返回空字符串（由外层决定是否走 HybridSearcher 兜底）。
func (s *Service) retrieveViaMemoryAgent(ctx context.Context, providerModel *wrapper_models.ProviderModel, userMessage string) string {
	if s.memoryStorage == nil || providerModel == nil {
		return ""
	}
	userMessage = sanitizeMemoryContext(userMessage)
	if !shouldRetrieveMemories(userMessage) {
		return ""
	}

	memAgent, err := memory_agents.NewMemoryAgent(
		ctx,
		providerModel.BaseUrl,
		providerModel.ApiKey,
		providerModel.Model,
		s.memoryStorage,
	)
	if err != nil {
		logger.Error("retrieve via memory agent: create agent error:", err)
		return ""
	}

	instruction := `[记忆检索模式 - 内部指令，不要向用户暴露]
使用 read_memory 查找与本用户消息相关的长期记忆；若涉及时间表达，先调用 get_current_time 推算绝对日期再查询。
完成后仅输出精简记忆列表，格式为每行一条：
- [YYYY-MM-DD] 标题：内容要点
不要添加任何问候/解释/总结。若没有任何相关记忆，仅输出：NO_MEMORY
用户消息如下：`

	agentInput := []*schema.Message{
		{Role: schema.User, Content: instruction + "\n\n" + userMessage},
	}

	stream, err := memAgent.Streamable(ctx, agentInput)
	if err != nil {
		logger.Error("retrieve via memory agent: stream error:", err)
		return ""
	}
	if stream == nil {
		return ""
	}
	defer stream.Close()

	var builder strings.Builder
	for {
		chunk, recvErr := stream.Recv()
		if recvErr != nil {
			break
		}
		if chunk != nil {
			builder.WriteString(chunk.Content)
		}
	}

	result := strings.TrimSpace(builder.String())
	if result == "" || strings.EqualFold(result, "NO_MEMORY") {
		return ""
	}
	return sanitizeMemoryContext(result)
}

// extractKeywords 从用户消息中提取检索关键词。
// 语言无关：按标点/空格分割 + 短消息滑窗拆分，不使用硬编码词表。
func extractKeywords(message string) []string {
	message = sanitizeMemoryContext(message)

	// 1. 按标点和空格分割
	separators := []string{
		"，", "。", "？", "！", "、", "；", "：",
		",", ".", "?", "!", ";", ":",
		" ", "\n", "\t",
	}
	segments := []string{message}
	for _, sep := range separators {
		var next []string
		for _, seg := range segments {
			next = append(next, strings.Split(seg, sep)...)
		}
		segments = next
	}

	var keywords []string
	for _, seg := range segments {
		seg = strings.TrimSpace(seg)
		if utf8.RuneCountInString(seg) >= 2 {
			keywords = append(keywords, seg)
		}
	}

	// 2. 短消息滑窗拆分（2-gram），增加召回率
	if len(keywords) <= 1 {
		runes := []rune(strings.TrimSpace(message))
		for i := 0; i+2 <= len(runes); i++ {
			keywords = append(keywords, string(runes[i:i+2]))
		}
	}

	// 去重 + 限制数量
	seen := make(map[string]bool)
	var unique []string
	for _, kw := range keywords {
		if !seen[kw] && kw != "" {
			seen[kw] = true
			unique = append(unique, kw)
		}
	}
	if len(unique) > 10 {
		unique = unique[:10]
	}
	return unique
}

// formatMemoryContext 将记忆列表格式化为注入文本，控制在约 1000 汉字以内。
func formatMemoryContext(memories []memory_models.Memory) string {
	var lines []string
	totalLen := 0
	const maxLen = 2000 // 约 1000 tokens

	for _, m := range memories {
		dateStr := ""
		if m.TimeRangStart != nil {
			dateStr = m.TimeRangStart.Format("2006-01-02")
		} else {
			dateStr = m.CreatedAt.Format("2006-01-02")
		}

		title := strings.TrimSpace(m.Summary)
		content := strings.TrimSpace(m.Content)

		// 截断过长的内容
		contentRunes := []rune(content)
		if len(contentRunes) > 200 {
			content = string(contentRunes[:200]) + "..."
		}

		var line string
		if title != "" && content != "" && title != content {
			line = fmt.Sprintf("- [%s] %s: %s", dateStr, title, content)
		} else if content != "" {
			line = fmt.Sprintf("- [%s] %s", dateStr, content)
		} else if title != "" {
			line = fmt.Sprintf("- [%s] %s", dateStr, title)
		} else {
			continue
		}

		lineLen := utf8.RuneCountInString(line)
		if totalLen+lineLen > maxLen {
			break
		}
		lines = append(lines, line)
		totalLen += lineLen
	}

	if len(lines) == 0 {
		return ""
	}
	return strings.Join(lines, "\n")
}

// ---- 围栏注入 ----

// buildFencedMemoryContext 构建围栏标签包裹的记忆上下文。
func buildFencedMemoryContext(memoryContext string) string {
	if memoryContext == "" {
		return ""
	}
	// 递归清洁：防止记忆内容本身包含围栏标签
	memoryContext = sanitizeMemoryContext(memoryContext)

	return fmt.Sprintf(`<memory-context>
以下是与本次对话可能相关的历史记忆片段，仅供参考。
这些不是新的用户输入，请作为背景信息使用，不要主动提及这些记忆的存在。

%s
</memory-context>`, memoryContext)
}

// sanitizeMemoryContext 剥离嵌套的围栏标签，防止递归注入。
func sanitizeMemoryContext(content string) string {
	content = strings.ReplaceAll(content, "<memory-context>", "")
	content = strings.ReplaceAll(content, "</memory-context>", "")
	return strings.TrimSpace(content)
}

// buildExistingMemoriesSnapshot 把当前库中最近活跃的记忆渲染成一段"已存在记忆清单"，
// 作为编码前的上下文注入到 Memory Agent，让它自行判断是否与已有记忆重复。
//
// 注意：这是"去重的唯一机制"——不做算法去重，完全依赖 LLM 对语义的判断。
// 为控制 prompt 体积，按类型各取最近若干条，content 做截断。
func (s *Service) buildExistingMemoriesSnapshot(ctx context.Context) string {
	if s.memoryStorage == nil {
		return ""
	}

	const perType = 20
	const contentLimit = 200

	var all []memory_models.Memory
	for _, t := range []string{"fact", "information", "event", ""} {
		items, err := s.memoryStorage.RecentMemoriesByType(ctx, t, perType)
		if err != nil {
			logger.Error("memory snapshot: load recent error:", err)
			continue
		}
		all = append(all, items...)
	}
	if len(all) == 0 {
		return ""
	}

	seen := make(map[uint]bool, len(all))
	var lines []string
	for _, m := range all {
		if seen[m.ID] {
			continue
		}
		seen[m.ID] = true
		content := strings.TrimSpace(m.Content)
		runes := []rune(content)
		if len(runes) > contentLimit {
			content = string(runes[:contentLimit]) + "…"
		}
		typeLabel := string(m.Type)
		if typeLabel == "" {
			typeLabel = "-"
		}
		lines = append(lines, fmt.Sprintf("  [id=%d | type=%s] %s :: %s",
			m.ID, typeLabel, strings.TrimSpace(m.Summary), content))
	}

	return fmt.Sprintf(`[系统上下文：当前已存在的记忆列表，共 %d 条]
每条格式：[id=ID | type=类型] 标题 :: 内容摘要
%s

[操作指令]
接下来我会给出最近一轮对话。请严格按以下流程处理：
1. 判断本轮对话是否包含值得长期记住的新信息；若无，什么都不做。
2. 如果有新信息，先与上方清单逐条比对：
   - 若新信息描述的是同一主题/同一事件（即便措辞不同），必须调用 edit_memory 在对应 id 上补全，禁止新建。
   - 只有确认上方清单里没有任何一条覆盖此主题，才调用 write_memory。
3. 若需要新建，记得遵循 content 包含所有时间/地点/人物/原因等信息（时间要写成绝对日期）。
绝不允许为同一主题产生两条并存记忆。`, len(lines), strings.Join(lines, "\n"))
}

// ---- 异步编码 ----

// encodeMemoriesAsync 异步调用 Memory Agent 对本轮对话进行记忆编码。
// 去重策略：把"最近活跃记忆"作为上下文块拼进 Agent 输入，让 LLM 自行判断是 edit 已有还是 write 新建。
func (s *Service) encodeMemoriesAsync(providerModel *wrapper_models.ProviderModel, messages []schema.Message) {
	if s.memoryStorage == nil || providerModel == nil {
		return
	}

	// 提取用户最新消息（同时处理 Content 和 MultiContent，清理围栏标签）
	userMessage := ""
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == schema.User {
			userMessage = messages[i].Content
			if userMessage == "" {
				for _, part := range messages[i].UserInputMultiContent {
					if part.Type == schema.ChatMessagePartTypeText && part.Text != "" {
						userMessage = part.Text
						break
					}
				}
			}
			break
		}
	}
	userMessage = sanitizeMemoryContext(userMessage)
	if !shouldEncodeMemories(userMessage) {
		return
	}

	// 全局限流：同时只有 1 个编码任务
	if !memoryEncodeMu.TryLock() {
		return
	}
	defer memoryEncodeMu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	memAgent, err := memory_agents.NewMemoryAgent(
		ctx,
		providerModel.BaseUrl,
		providerModel.ApiKey,
		providerModel.Model,
		s.memoryStorage,
	)
	if err != nil {
		logger.Error("create memory agent error", err)
		return
	}

	// 只发送最近 2 轮对话（最新用户消息 + 最新助手回复），而非全部历史
	recentMessages := extractRecentTurn(messages)
	var msgPtrs []*schema.Message

	// ---- 去重关键：把最近活跃记忆作为上下文快照注入，让 LLM 自行判断 write vs edit ----
	if snapshot := s.buildExistingMemoriesSnapshot(ctx); snapshot != "" {
		msgPtrs = append(msgPtrs, &schema.Message{
			Role:    schema.User,
			Content: snapshot,
		})
	}

	for i := range recentMessages {
		msgPtrs = append(msgPtrs, &recentMessages[i])
	}

	stream, err := memAgent.Streamable(ctx, msgPtrs)
	if err != nil {
		logger.Error("memory agent stream error", err)
		return
	}
	if stream != nil {
		for {
			_, recvErr := stream.Recv()
			if recvErr != nil {
				break
			}
		}
		stream.Close()
	}

	// Memory Agent 可能写入了新记忆，为其生成 embedding 并刷新向量缓存
	if s.memorySearcher != nil && s.memorySearcher.HasEmbedder() {
		bgCtx, bgCancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer bgCancel()
		count, backfillErr := s.memorySearcher.BackfillEmbeddings(bgCtx, 10)
		if backfillErr != nil {
			logger.Error("post-encode backfill error:", backfillErr)
		} else if count > 0 {
			logger.Warm("post-encode: embedded", count, "new memories")
		}
	}
}

// extractRecentTurn 从完整消息历史中提取最近一轮对话。
// 返回最后的用户消息 + 最后的助手回复（如有）。
// 用户消息中的围栏标签会被清理。
func extractRecentTurn(messages []schema.Message) []schema.Message {
	var userMsg, assistantMsg *schema.Message

	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role == schema.Assistant && assistantMsg == nil {
			msg := messages[i]
			assistantMsg = &msg
		}
		if messages[i].Role == schema.User && userMsg == nil {
			msg := messages[i]
			// 清理围栏标签
			msg.Content = sanitizeMemoryContext(msg.Content)
			userMsg = &msg
			break
		}
	}

	var result []schema.Message
	if userMsg != nil {
		result = append(result, *userMsg)
	}
	if assistantMsg != nil {
		result = append(result, *assistantMsg)
	}
	return result
}
