// MemoryAgent - 记忆 Agent，基于 RAG + 向量数据库（BoltDB 持久化）实现记录和语义检索日常

package agents

import (
	"bufio"
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/embedding"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"go.etcd.io/bbolt"
)

const (
	dailyVectorDB   = "daily_memory.bolt"
	dailyLogFileOld = "daily_memory.jsonl"
	entriesBucket   = "entries"
	metaBucket      = "meta"
	seqKey          = "seq"
)

const memoryAgentInstruction = `你是一个帮助用户记录和回顾日常的助手。
记录：当用户提及发生的事、消费、安排（如「昨天买了一台电脑花了5200」「今天去了公园」「记一下午饭吃了面」）时，必须使用 record_daily。提取标题（简短概括）和内容（完整描述，含金额、物品等）。若涉及「昨天」「前天」，date 填 "yesterday" 或 "2 days ago"。
回忆：当用户询问过去或今日的安排/消费/做过的事（如「昨天我花了多少钱」「今天要做什么」「昨天做了什么」）时，必须使用 read_daily。date 填 "today"、"yesterday" 或具体 YYYY-MM-DD，query 可填语义描述（如「花了多少钱」「买了什么」）。有记录则如实回复；无记录则告知暂无。
用简洁友好的中文回复。`

// DailyEntry 单条日常记录
type DailyEntry struct {
	ID        int       `json:"id"`
	Date      string    `json:"date"`       // YYYY-MM-DD
	Title     string    `json:"title"`      // 简短标题
	Content   string    `json:"content"`    // 详细内容
	CreatedAt time.Time `json:"created_at"` // 记录时间
	Embedding []float64 `json:"embedding,omitempty"`
}

// VectorStore 向量数据库，使用 BoltDB 持久化
type VectorStore struct {
	db   *bbolt.DB
	path string
	mu   sync.RWMutex
}

// NewVectorStore 创建并打开向量数据库
func NewVectorStore(path string) (*VectorStore, error) {
	if path == "" {
		path = dailyVectorDB
	}
	db, err := bbolt.Open(path, 0600, nil)
	if err != nil {
		return nil, err
	}
	err = db.Update(func(tx *bbolt.Tx) error {
		if _, err := tx.CreateBucketIfNotExists([]byte(entriesBucket)); err != nil {
			return err
		}
		if _, err := tx.CreateBucketIfNotExists([]byte(metaBucket)); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		db.Close()
		return nil, err
	}
	vs := &VectorStore{db: db, path: path}
	vs.migrateFromJSONL()
	return vs, nil
}

// migrateFromJSONL 从旧的 JSONL 文件一次性迁移到 BoltDB
func (s *VectorStore) migrateFromJSONL() {
	n, _ := s.Count()
	if n > 0 {
		return
	}
	f, err := os.Open(dailyLogFileOld)
	if err != nil {
		return
	}
	defer f.Close()
	var migrated int
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		var e DailyEntry
		if err := json.Unmarshal([]byte(line), &e); err != nil {
			continue
		}
		e.ID = 0
		if err := s.Add(&e); err != nil {
			continue
		}
		migrated++
	}
	if migrated > 0 {
		fmt.Fprintf(os.Stderr, "[MemoryAgent] 已从 %s 迁移 %d 条记录到向量数据库\n", dailyLogFileOld, migrated)
	}
}

// Close 关闭数据库
func (s *VectorStore) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db != nil {
		err := s.db.Close()
		s.db = nil
		return err
	}
	return nil
}

// Add 添加一条记录
func (s *VectorStore) Add(entry *DailyEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.db.Update(func(tx *bbolt.Tx) error {
		ent := tx.Bucket([]byte(entriesBucket))
		meta := tx.Bucket([]byte(metaBucket))

		var seq uint64
		if v := meta.Get([]byte(seqKey)); len(v) == 8 {
			seq = binary.BigEndian.Uint64(v)
		}
		seq++
		entry.ID = int(seq)

		seqBytes := make([]byte, 8)
		binary.BigEndian.PutUint64(seqBytes, seq)
		if err := meta.Put([]byte(seqKey), seqBytes); err != nil {
			return err
		}

		key := make([]byte, 8)
		binary.BigEndian.PutUint64(key, seq)
		data, err := json.Marshal(entry)
		if err != nil {
			return err
		}
		return ent.Put(key, data)
	})
}

// SearchByVector 按向量相似度检索 top-k
func (s *VectorStore) SearchByVector(queryVec []float64, dateFilter string, topK int) ([]DailyEntry, error) {
	if topK <= 0 {
		topK = 5
	}

	var scoredList []scored
	s.mu.RLock()
	err := s.db.View(func(tx *bbolt.Tx) error {
		ent := tx.Bucket([]byte(entriesBucket))
		if ent == nil {
			return nil
		}
		return ent.ForEach(func(k, v []byte) error {
			var e DailyEntry
			if err := json.Unmarshal(v, &e); err != nil {
				return nil
			}
			if len(e.Embedding) == 0 {
				return nil
			}
			if dateFilter != "" && e.Date != dateFilter {
				return nil
			}
			sim := cosineSimilarity(queryVec, e.Embedding)
			scoredList = append(scoredList, scored{entry: e, score: sim})
			return nil
		})
	})
	s.mu.RUnlock()

	if err != nil {
		return nil, err
	}

	sort.Slice(scoredList, func(i, j int) bool { return scoredList[i].score > scoredList[j].score })
	if len(scoredList) > topK {
		scoredList = scoredList[:topK]
	}

	result := make([]DailyEntry, 0, len(scoredList))
	for _, sc := range scoredList {
		result = append(result, sc.entry)
	}
	return result, nil
}

// SearchByFilter 按日期/关键词筛选（无语义检索时）
func (s *VectorStore) SearchByFilter(dateFilter, keyword string, limit int) ([]DailyEntry, error) {
	if limit <= 0 {
		limit = 10
	}
	var entries []DailyEntry

	s.mu.RLock()
	err := s.db.View(func(tx *bbolt.Tx) error {
		ent := tx.Bucket([]byte(entriesBucket))
		if ent == nil {
			return nil
		}
		return ent.ForEach(func(k, v []byte) error {
			var e DailyEntry
			if err := json.Unmarshal(v, &e); err != nil {
				return nil
			}
			if dateFilter != "" && e.Date != dateFilter {
				return nil
			}
			if keyword != "" {
				if !strings.Contains(e.Title, keyword) && !strings.Contains(e.Content, keyword) {
					return nil
				}
			}
			entries = append(entries, e)
			return nil
		})
	})
	s.mu.RUnlock()

	if err != nil {
		return nil, err
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Date != entries[j].Date {
			return entries[i].Date > entries[j].Date
		}
		return entries[i].CreatedAt.After(entries[j].CreatedAt)
	})

	if len(entries) > limit {
		entries = entries[:limit]
	}
	return entries, nil
}

// Count 返回记录总数
func (s *VectorStore) Count() (int, error) {
	var n int
	s.mu.RLock()
	err := s.db.View(func(tx *bbolt.Tx) error {
		ent := tx.Bucket([]byte(entriesBucket))
		if ent == nil {
			return nil
		}
		n = ent.Stats().KeyN
		return nil
	})
	s.mu.RUnlock()
	return n, err
}

type scored struct {
	entry DailyEntry
	score float64
}

// resolveDate 将 today/yesterday/昨日/今日 等解析为 YYYY-MM-DD
func resolveDate(s string) (string, bool) {
	s = strings.TrimSpace(strings.ToLower(s))
	if s == "" {
		return "", false
	}
	now := time.Now()
	switch s {
	case "today", "今日", "今天":
		return now.Format("2006-01-02"), true
	case "yesterday", "昨日", "昨天":
		return now.AddDate(0, 0, -1).Format("2006-01-02"), true
	case "2 days ago", "前天":
		return now.AddDate(0, 0, -2).Format("2006-01-02"), true
	case "3 days ago":
		return now.AddDate(0, 0, -3).Format("2006-01-02"), true
	}
	return "", false
}

func cosineSimilarity(a, b []float64) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var dot, normA, normB float64
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}

// record_daily 工具
type recordDailyParams struct {
	Title   string `json:"title" jsonschema:"title=标题,description=日常的简短标题，例如「今天去公园」,required=true"`
	Content string `json:"content" jsonschema:"title=内容,description=日常的详细描述，可以包括做了什么、心情、见闻等,required=true"`
	Date    string `json:"date,omitempty" jsonschema:"title=日期,description=日常发生的日期，YYYY-MM-DD 或 today/yesterday/昨日，不填默认今天"`
}

// read_daily 工具
type readDailyParams struct {
	Query   string `json:"query,omitempty" jsonschema:"title=语义查询,description=用自然语言描述想回忆的内容，如「去公园的那次」「心情不好的时候」。优先使用语义检索"`
	Date    string `json:"date,omitempty" jsonschema:"title=日期,description=要查询的日期，YYYY-MM-DD 或 today/yesterday/昨日，可与语义查询组合筛选"`
	Keyword string `json:"keyword,omitempty" jsonschema:"title=关键词,description=在标题或内容中匹配的关键词，当无语义查询时使用"`
	TopK    int    `json:"top_k,omitempty" jsonschema:"title=返回条数,description=语义检索时返回的最相关条数，默认 5"`
}

func getMemoryTools(store *VectorStore, embedder embedding.Embedder) []tool.BaseTool {
	recordTool, _ := utils.InferTool(
		"record_daily",
		"记录用户的日常、消费、安排。需要 title 和 content，date 可选：默认今天，或填 today/yesterday/昨日。会自动生成向量嵌入并持久化。",
		func(ctx context.Context, params recordDailyParams) (string, error) {
			title := strings.TrimSpace(params.Title)
			content := strings.TrimSpace(params.Content)
			if title == "" {
				return "请提供日常的标题。", nil
			}
			if content == "" {
				return "请提供日常的内容。", nil
			}
			date := strings.TrimSpace(params.Date)
			if date == "" {
				date = time.Now().Format("2006-01-02")
			} else if resolved, ok := resolveDate(date); ok {
				date = resolved
			} else if _, err := time.Parse("2006-01-02", date); err != nil {
				return "日期格式错误，请使用 YYYY-MM-DD，或 today/yesterday。", nil
			}

			entry := &DailyEntry{
				Date:      date,
				Title:     title,
				Content:   content,
				CreatedAt: time.Now(),
			}

			if embedder != nil {
				textToEmbed := title + " " + content
				vecs, err := embedder.EmbedStrings(ctx, []string{textToEmbed})
				if err == nil && len(vecs) > 0 {
					entry.Embedding = vecs[0]
				}
			}

			if err := store.Add(entry); err != nil {
				fmt.Fprintf(os.Stderr, "[MemoryAgent] 写入向量数据库失败: %v\n", err)
				return "记录失败，请稍后重试。", nil
			}

			return fmt.Sprintf("已记录日常：%s（%s）", title, date), nil
		},
	)

	readTool, _ := utils.InferTool(
		"read_daily",
		"读取用户的日常/消费记录。date 可为 YYYY-MM-DD，或 today/yesterday/昨日；query 为语义描述（如「花了多少钱」「买了什么」）；keyword 为关键词。top_k 默认 5。",
		func(ctx context.Context, params readDailyParams) (string, error) {
			n, err := store.Count()
			if err != nil {
				return "读取失败，请稍后重试。", nil
			}
			if n == 0 {
				return "暂无日常记录，可以说「帮我记一下今天...」来记录。", nil
			}

			query := strings.TrimSpace(params.Query)
			date := strings.TrimSpace(params.Date)
			if resolved, ok := resolveDate(date); ok {
				date = resolved
			}
			keyword := strings.TrimSpace(params.Keyword)
			topK := params.TopK
			if topK <= 0 {
				topK = 5
			}

			var result []DailyEntry

			if query != "" && embedder != nil {
				vecs, err := embedder.EmbedStrings(ctx, []string{query})
				if err == nil && len(vecs) > 0 {
					result, err = store.SearchByVector(vecs[0], date, topK)
					if err != nil {
						fmt.Fprintf(os.Stderr, "[MemoryAgent] 向量检索失败: %v\n", err)
					}
				}
			}

			if len(result) == 0 {
				result, err = store.SearchByFilter(date, keyword, topK)
				if err != nil {
					return "读取失败，请稍后重试。", nil
				}
			}

			if len(result) == 0 {
				if query != "" {
					return fmt.Sprintf("未找到与「%s」相关的日常记录。", query), nil
				}
				if date != "" && keyword != "" {
					return fmt.Sprintf("未找到 %s 包含「%s」的日常记录。", date, keyword), nil
				}
				if date != "" {
					return fmt.Sprintf("未找到 %s 的日常记录。", date), nil
				}
				if keyword != "" {
					return fmt.Sprintf("未找到包含「%s」的日常记录。", keyword), nil
				}
				return "暂无日常记录。", nil
			}

			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("找到 %d 条日常：\n\n", len(result)))
			for _, e := range result {
				sb.WriteString(fmt.Sprintf("【%s】%s\n", e.Date, e.Title))
				sb.WriteString(e.Content)
				sb.WriteString("\n\n")
			}
			return strings.TrimRight(sb.String(), "\n"), nil
		},
	)

	return []tool.BaseTool{recordTool, readTool}
}

// NewMemoryAgent 创建记忆 Agent（RAG + 向量数据库），需传入 Embedder 以启用语义检索
func NewMemoryAgent(ctx context.Context, chatModel model.ToolCallingChatModel, embedder embedding.Embedder) (adk.Agent, error) {
	storePath := os.Getenv("DAILY_VECTOR_DB")
	if storePath == "" {
		storePath = dailyVectorDB
	}
	store, err := NewVectorStore(storePath)
	if err != nil {
		return nil, fmt.Errorf("初始化向量数据库失败: %w", err)
	}
	tools := getMemoryTools(store, embedder)

	return adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        "MemoryAgent",
		Description: "可以记录和读取用户日常、消费、安排的助手。当用户提及发生的事/消费（如「昨天买了电脑花了5200」）要记录，或询问过去/今日安排/消费（如「昨天花了多少钱」「今天要做什么」）时，必须转交此 Agent。支持 RAG 语义检索和按日期查询。",
		Instruction: memoryAgentInstruction,
		Model:       chatModel,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: tools,
			},
		},
	})
}
