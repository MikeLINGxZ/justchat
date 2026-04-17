package storage

import (
	"context"
	"encoding/binary"
	"math"
	"time"

	"gorm.io/gorm"
)

// MemoryEmbedding 存储记忆的嵌入向量。
type MemoryEmbedding struct {
	ID         uint      `gorm:"primarykey"`
	MemoryID   uint      `gorm:"uniqueIndex;not null"`
	Vector     []byte    `gorm:"type:blob;not null"` // float32 数组序列化
	ModelName  string    `gorm:"type:varchar(100);not null"`
	Dimensions int       `gorm:"not null"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

// AutoMigrateEmbeddings 创建嵌入表。
func (s *Storage) AutoMigrateEmbeddings() error {
	return s.sqliteDb.AutoMigrate(&MemoryEmbedding{})
}

// SaveEmbedding 保存或更新一条记忆的嵌入向量。
func (s *Storage) SaveEmbedding(ctx context.Context, memoryID uint, vector []float32, modelName string) error {
	blob := float32sToBytes(vector)
	emb := MemoryEmbedding{
		MemoryID:   memoryID,
		Vector:     blob,
		ModelName:  modelName,
		Dimensions: len(vector),
	}
	return s.sqliteDb.WithContext(ctx).
		Where("memory_id = ?", memoryID).
		Assign(emb).
		FirstOrCreate(&emb).Error
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

// GetMemoryIDsWithoutEmbedding 返回没有嵌入的活跃记忆 ID 列表。
func (s *Storage) GetMemoryIDsWithoutEmbedding(ctx context.Context, limit int) ([]uint, error) {
	var ids []uint
	err := s.sqliteDb.WithContext(ctx).
		Raw(`SELECT m.id FROM memories m
			 LEFT JOIN memory_embeddings e ON e.memory_id = m.id
			 WHERE e.id IS NULL AND m.is_forgotten = 0
			 LIMIT ?`, limit).
		Scan(&ids).Error
	return ids, err
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
