package storage

import (
	"context"
	"encoding/binary"
	"math"
	"time"

	"gorm.io/gorm"
)

// CurrentEmbeddingSchemaVersion 当前 embedding 文本拼接规则版本。
// 每次修改 buildEmbedText 规则都要递增该常量，以触发后台重新生成旧向量。
const CurrentEmbeddingSchemaVersion = 2

// MemoryEmbedding 存储记忆的嵌入向量。
type MemoryEmbedding struct {
	ID            uint      `gorm:"primarykey"`
	MemoryID      uint      `gorm:"uniqueIndex;not null"`
	Vector        []byte    `gorm:"type:blob;not null"` // float32 数组序列化
	ModelName     string    `gorm:"type:varchar(100);not null"`
	Dimensions    int       `gorm:"not null"`
	SchemaVersion int       `gorm:"default:1"` // 拼接规则版本
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
}

// AutoMigrateEmbeddings 创建嵌入表。
func (s *Storage) AutoMigrateEmbeddings() error {
	return s.sqliteDb.AutoMigrate(&MemoryEmbedding{})
}

// SaveEmbedding 保存或更新一条记忆的嵌入向量，并回写 Memory.EmbeddingID。
func (s *Storage) SaveEmbedding(ctx context.Context, memoryID uint, vector []float32, modelName string) error {
	blob := float32sToBytes(vector)
	emb := MemoryEmbedding{
		MemoryID:      memoryID,
		Vector:        blob,
		ModelName:     modelName,
		Dimensions:    len(vector),
		SchemaVersion: CurrentEmbeddingSchemaVersion,
	}
	if err := s.sqliteDb.WithContext(ctx).
		Where("memory_id = ?", memoryID).
		Assign(emb).
		FirstOrCreate(&emb).Error; err != nil {
		return err
	}

	// 回写 Memory.EmbeddingID，使前端能判断该记忆已向量化
	return s.sqliteDb.WithContext(ctx).
		Table("memories").
		Where("id = ?", memoryID).
		Update("embedding_id", emb.ID).Error
}

// EmbeddingEntry 用于向量检索的缓存条目。
type EmbeddingEntry struct {
	MemoryID uint
	Vector   []float32
}

// LoadAllEmbeddings 加载所有嵌入向量（用于内存缓存）。
func (s *Storage) LoadAllEmbeddings(ctx context.Context) ([]EmbeddingEntry, error) {
	var rows []MemoryEmbedding
	if err := s.sqliteDb.WithContext(ctx).
		Joins("JOIN memories ON memories.id = memory_embeddings.memory_id AND memories.is_forgotten = 0").
		Find(&rows).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	entries := make([]EmbeddingEntry, 0, len(rows))
	for _, row := range rows {
		vec := bytesToFloat32s(row.Vector)
		if len(vec) > 0 {
			entries = append(entries, EmbeddingEntry{
				MemoryID: row.MemoryID,
				Vector:   vec,
			})
		}
	}
	return entries, nil
}

// CountEmbeddings 返回嵌入数量。
func (s *Storage) CountEmbeddings(ctx context.Context) (int64, error) {
	var count int64
	err := s.sqliteDb.WithContext(ctx).Model(&MemoryEmbedding{}).Count(&count).Error
	return count, err
}

// GetMemoryIDsWithoutEmbedding 返回需要生成/更新嵌入的活跃记忆 ID 列表。
// 触发条件：无嵌入 OR 嵌入的 schema 版本落后于 CurrentEmbeddingSchemaVersion。
func (s *Storage) GetMemoryIDsWithoutEmbedding(ctx context.Context, limit int) ([]uint, error) {
	var ids []uint
	err := s.sqliteDb.WithContext(ctx).
		Raw(`SELECT m.id FROM memories m
			 LEFT JOIN memory_embeddings e ON e.memory_id = m.id
			 WHERE m.is_forgotten = 0
			   AND (e.id IS NULL OR COALESCE(e.schema_version, 1) < ?)
			 LIMIT ?`, CurrentEmbeddingSchemaVersion, limit).
		Scan(&ids).Error
	return ids, err
}

// RepairEmbeddingIDs 修复历史数据：将已有 embedding 但 memories.embedding_id 为空的记忆回写关联。
func (s *Storage) RepairEmbeddingIDs(ctx context.Context) error {
	return s.sqliteDb.WithContext(ctx).Exec(
		`UPDATE memories SET embedding_id = (
			SELECT e.id FROM memory_embeddings e WHERE e.memory_id = memories.id LIMIT 1
		) WHERE embedding_id IS NULL AND EXISTS (
			SELECT 1 FROM memory_embeddings e WHERE e.memory_id = memories.id
		)`).Error
}

// ---- 序列化工具 ----

func float32sToBytes(floats []float32) []byte {
	buf := make([]byte, len(floats)*4)
	for i, f := range floats {
		binary.LittleEndian.PutUint32(buf[i*4:], math.Float32bits(f))
	}
	return buf
}

func bytesToFloat32s(data []byte) []float32 {
	if len(data)%4 != 0 {
		return nil
	}
	floats := make([]float32, len(data)/4)
	for i := range floats {
		floats[i] = math.Float32frombits(binary.LittleEndian.Uint32(data[i*4:]))
	}
	return floats
}

// CosineSimilarity 计算两个向量的余弦相似度。
func CosineSimilarity(a, b []float32) float64 {
	if len(a) != len(b) || len(a) == 0 {
		return 0
	}
	var dot, normA, normB float64
	for i := range a {
		dot += float64(a[i]) * float64(b[i])
		normA += float64(a[i]) * float64(a[i])
		normB += float64(b[i]) * float64(b[i])
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (math.Sqrt(normA) * math.Sqrt(normB))
}
