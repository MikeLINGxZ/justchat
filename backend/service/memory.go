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

// ---- 异步编码 ----

// encodeMemoriesAsync 异步调用 Memory Agent 对本轮对话进行记忆编码。
// 编码前先检查是否已有高度相似的记忆，避免重复存储。
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

	// ---- 编码前去重：检查是否已有高度相似的记忆 ----
	if s.hasSimilarMemory(userMessage) {
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
	// 避免 Memory Agent 看到旧消息重复编码
	recentMessages := extractRecentTurn(messages)
	var msgPtrs []*schema.Message
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

// hasSimilarMemory 检查是否已有与用户消息高度相似的记忆。
// 只做精确度高的比较（summary 包含检查），避免误拦新记忆。
func (s *Service) hasSimilarMemory(userMessage string) bool {
	userMessage = sanitizeMemoryContext(userMessage)
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	// 直接用用户消息的前 20 字做 LIKE 搜索，检查是否已有几乎相同的记忆
	runes := []rune(strings.TrimSpace(userMessage))
	if len(runes) < 5 {
		return false
	}
	searchText := string(runes)
	if len(runes) > 20 {
		searchText = string(runes[:20])
	}

	var count int64
	s.memoryStorage.DB().WithContext(ctx).
		Model(&memory_models.Memory{}).
		Where("is_forgotten = ? AND (summary LIKE ? OR content LIKE ?)", false,
			"%"+searchText+"%", "%"+searchText+"%").
		Count(&count)

	return count > 0
}

// jaccardSimilarity 计算两个文本的字符 bigram Jaccard 相似度。
func jaccardSimilarity(a, b []rune) float64 {
	if len(a) < 2 || len(b) < 2 {
		return 0
	}
	setA := make(map[string]bool)
	for i := 0; i+2 <= len(a); i++ {
		setA[string(a[i:i+2])] = true
	}
	setB := make(map[string]bool)
	for i := 0; i+2 <= len(b); i++ {
		setB[string(b[i:i+2])] = true
	}
	intersection := 0
	for k := range setA {
		if setB[k] {
			intersection++
		}
	}
	union := len(setA) + len(setB) - intersection
	if union == 0 {
		return 0
	}
	return float64(intersection) / float64(union)
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
